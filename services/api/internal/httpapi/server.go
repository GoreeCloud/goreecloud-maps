package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/auth"
	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/providers"
	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/store"
)

type Server struct {
	logger    *slog.Logger
	store     *store.Store
	verifier  *auth.Verifier
	providers providers.Set
	mux       *http.ServeMux
}

type apiError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func New(logger *slog.Logger, dataStore *store.Store, verifier *auth.Verifier, providerSet providers.Set) http.Handler {
	server := &Server{
		logger:    logger,
		store:     dataStore,
		verifier:  verifier,
		providers: providerSet,
		mux:       http.NewServeMux(),
	}

	server.mux.HandleFunc("GET /healthz", server.health)
	server.mux.HandleFunc("GET /readyz", server.ready)
	server.mux.HandleFunc("GET /api/v1/capabilities", server.capabilities)
	server.mux.HandleFunc("GET /api/v1/me", server.me)
	server.mux.HandleFunc("GET /api/v1/search", server.searchPlaces)
	server.mux.HandleFunc("GET /api/v1/reverse", server.reversePlace)
	server.mux.HandleFunc("POST /api/v1/routes", server.createRoute)
	server.mux.HandleFunc("GET /api/v1/collections", server.listCollections)
	server.mux.HandleFunc("POST /api/v1/collections", server.createCollection)
	server.mux.HandleFunc("PATCH /api/v1/collections/{collectionID}", server.updateCollection)
	server.mux.HandleFunc("GET /api/v1/collections/{collectionID}/members", server.listCollectionMembers)
	server.mux.HandleFunc("POST /api/v1/collections/{collectionID}/members", server.addCollectionMember)
	server.mux.HandleFunc("PATCH /api/v1/collections/{collectionID}/members/{memberUserID}", server.updateCollectionMember)
	server.mux.HandleFunc("DELETE /api/v1/collections/{collectionID}/members/{memberUserID}", server.removeCollectionMember)
	server.mux.HandleFunc("GET /api/v1/collections/{collectionID}/items", server.listCollectionItems)
	server.mux.HandleFunc("POST /api/v1/collections/{collectionID}/items", server.createCollectionItem)
	server.mux.HandleFunc("PATCH /api/v1/collections/{collectionID}/items/{itemID}", server.updateCollectionItem)
	server.mux.HandleFunc("DELETE /api/v1/collections/{collectionID}/items/{itemID}", server.deleteCollectionItem)

	return server.securityHeaders(server.mux)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(r.Context()); err != nil {
		s.logger.Error("database readiness check failed")
		writeAPIError(w, http.StatusServiceUnavailable, "not_ready", "The Maps API is not ready.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) capabilities(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"providers": s.providers.Capabilities(),
	})
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) searchPlaces(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.authenticatedUser(w, r); !ok {
		return
	}

	limit := 10
	if rawLimit := strings.TrimSpace(r.URL.Query().Get("limit")); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit < 1 || parsedLimit > 20 {
			writeAPIError(w, http.StatusBadRequest, "invalid_limit", "Search limit must be between 1 and 20.")
			return
		}
		limit = parsedLimit
	}

	results, err := s.providers.Search(r.Context(), r.URL.Query().Get("q"), limit, r.Header.Get("Accept-Language"))
	if err != nil {
		s.writeProviderError(w, err, "Place search")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": results})
}

func (s *Server) reversePlace(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.authenticatedUser(w, r); !ok {
		return
	}

	latitude, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("latitude")), 64)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_coordinate", "A valid latitude is required.")
		return
	}
	longitude, err := strconv.ParseFloat(strings.TrimSpace(r.URL.Query().Get("longitude")), 64)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_coordinate", "A valid longitude is required.")
		return
	}

	result, err := s.providers.Reverse(r.Context(), latitude, longitude, r.Header.Get("Accept-Language"))
	if err != nil {
		s.writeProviderError(w, err, "Reverse geocoding")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) createRoute(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.authenticatedUser(w, r); !ok {
		return
	}

	var payload providers.RouteRequest
	if !decodeJSONBody(w, r, 64<<10, &payload, "A valid route payload is required.") {
		return
	}

	route, err := s.providers.Route(r.Context(), payload)
	if err != nil {
		s.writeProviderError(w, err, "Route planning")
		return
	}
	writeJSON(w, http.StatusOK, route)
}

func (s *Server) listCollections(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
	if !ok {
		return
	}

	collections, err := s.store.ListCollections(r.Context(), user.ID)
	if err != nil {
		s.writeStoreError(w, err, "list collections")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"collections": collections})
}

func (s *Server) createCollection(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
	if !ok {
		return
	}

	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if !decodeJSONBody(w, r, 32<<10, &payload, "A valid collection payload is required.") {
		return
	}

	collection, err := s.store.CreateCollection(r.Context(), user.ID, store.CreateCollectionInput{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		s.writeStoreError(w, err, "create collection")
		return
	}
	writeJSON(w, http.StatusCreated, collection)
}

func (s *Server) updateCollection(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	var payload struct {
		Name             *string `json:"name"`
		Description      *string `json:"description"`
		ExpectedRevision int64   `json:"expectedRevision"`
	}
	if !decodeJSONBody(w, r, 32<<10, &payload, "A valid collection update is required.") {
		return
	}
	collection, err := s.store.UpdateCollection(r.Context(), user.ID, collectionID, store.UpdateCollectionInput{
		Name:             payload.Name,
		Description:      payload.Description,
		ExpectedRevision: payload.ExpectedRevision,
	})
	if err != nil {
		s.writeStoreError(w, err, "update collection")
		return
	}
	writeJSON(w, http.StatusOK, collection)
}

func (s *Server) listCollectionMembers(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	members, err := s.store.ListCollectionMembers(r.Context(), user.ID, collectionID)
	if err != nil {
		s.writeStoreError(w, err, "list collection members")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"members": members})
}

func (s *Server) addCollectionMember(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	var payload struct {
		UserID string `json:"userId"`
		Role   string `json:"role"`
	}
	if !decodeJSONBody(w, r, 16<<10, &payload, "A valid member payload is required.") {
		return
	}
	if !validUUIDString(payload.UserID) {
		writeAPIError(w, http.StatusUnprocessableEntity, "invalid_request", "A valid member user ID is required.")
		return
	}
	member, err := s.store.AddCollectionMember(r.Context(), user.ID, collectionID, payload.UserID, payload.Role)
	if err != nil {
		s.writeStoreError(w, err, "add collection member")
		return
	}
	writeJSON(w, http.StatusCreated, member)
}

func (s *Server) updateCollectionMember(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	memberUserID := strings.TrimSpace(r.PathValue("memberUserID"))
	if !validUUIDString(memberUserID) {
		writeAPIError(w, http.StatusBadRequest, "invalid_resource_id", "A valid member user ID is required.")
		return
	}
	var payload struct {
		Role string `json:"role"`
	}
	if !decodeJSONBody(w, r, 16<<10, &payload, "A valid member update is required.") {
		return
	}
	member, err := s.store.UpdateCollectionMemberRole(r.Context(), user.ID, collectionID, memberUserID, payload.Role)
	if err != nil {
		s.writeStoreError(w, err, "update collection member")
		return
	}
	writeJSON(w, http.StatusOK, member)
}

func (s *Server) removeCollectionMember(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	memberUserID := strings.TrimSpace(r.PathValue("memberUserID"))
	if !validUUIDString(memberUserID) {
		writeAPIError(w, http.StatusBadRequest, "invalid_resource_id", "A valid member user ID is required.")
		return
	}
	if err := s.store.RemoveCollectionMember(r.Context(), user.ID, collectionID, memberUserID); err != nil {
		s.writeStoreError(w, err, "remove collection member")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listCollectionItems(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	items, err := s.store.ListCollectionItems(r.Context(), user.ID, collectionID)
	if err != nil {
		s.writeStoreError(w, err, "list collection items")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) createCollectionItem(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	var payload struct {
		Provider        string  `json:"provider"`
		ProviderPlaceID string  `json:"providerPlaceId"`
		Name            string  `json:"name"`
		Address         string  `json:"address"`
		Latitude        float64 `json:"latitude"`
		Longitude       float64 `json:"longitude"`
		Note            string  `json:"note"`
		SortKey         int64   `json:"sortKey"`
	}
	if !decodeJSONBody(w, r, 64<<10, &payload, "A valid collection item is required.") {
		return
	}
	item, err := s.store.CreateCollectionItem(r.Context(), user.ID, collectionID, store.CreateCollectionItemInput{
		Provider:        payload.Provider,
		ProviderPlaceID: payload.ProviderPlaceID,
		Name:            payload.Name,
		Address:         payload.Address,
		Latitude:        payload.Latitude,
		Longitude:       payload.Longitude,
		Note:            payload.Note,
		SortKey:         payload.SortKey,
	})
	if err != nil {
		s.writeStoreError(w, err, "create collection item")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (s *Server) updateCollectionItem(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	itemID := strings.TrimSpace(r.PathValue("itemID"))
	if !validUUIDString(itemID) {
		writeAPIError(w, http.StatusBadRequest, "invalid_resource_id", "A valid collection item ID is required.")
		return
	}
	var payload struct {
		Name             *string  `json:"name"`
		Address          *string  `json:"address"`
		Latitude         *float64 `json:"latitude"`
		Longitude        *float64 `json:"longitude"`
		Note             *string  `json:"note"`
		SortKey          *int64   `json:"sortKey"`
		ExpectedRevision int64    `json:"expectedRevision"`
	}
	if !decodeJSONBody(w, r, 64<<10, &payload, "A valid collection item update is required.") {
		return
	}
	item, err := s.store.UpdateCollectionItem(r.Context(), user.ID, collectionID, itemID, store.UpdateCollectionItemInput{
		Name:             payload.Name,
		Address:          payload.Address,
		Latitude:         payload.Latitude,
		Longitude:        payload.Longitude,
		Note:             payload.Note,
		SortKey:          payload.SortKey,
		ExpectedRevision: payload.ExpectedRevision,
	})
	if err != nil {
		s.writeStoreError(w, err, "update collection item")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) deleteCollectionItem(w http.ResponseWriter, r *http.Request) {
	user, collectionID, ok := s.collectionRequest(w, r)
	if !ok {
		return
	}
	itemID := strings.TrimSpace(r.PathValue("itemID"))
	if !validUUIDString(itemID) {
		writeAPIError(w, http.StatusBadRequest, "invalid_resource_id", "A valid collection item ID is required.")
		return
	}
	expectedRevision, err := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("expectedRevision")), 10, 64)
	if err != nil || expectedRevision < 1 {
		writeAPIError(w, http.StatusBadRequest, "invalid_revision", "A positive expectedRevision query parameter is required.")
		return
	}
	if err := s.store.DeleteCollectionItem(r.Context(), user.ID, collectionID, itemID, expectedRevision); err != nil {
		s.writeStoreError(w, err, "delete collection item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) collectionRequest(w http.ResponseWriter, r *http.Request) (store.User, string, bool) {
	user, ok := s.authenticatedUser(w, r)
	if !ok {
		return store.User{}, "", false
	}
	collectionID := strings.TrimSpace(r.PathValue("collectionID"))
	if !validUUIDString(collectionID) {
		writeAPIError(w, http.StatusBadRequest, "invalid_resource_id", "A valid collection ID is required.")
		return store.User{}, "", false
	}
	return user, collectionID, true
}

func (s *Server) authenticatedUser(w http.ResponseWriter, r *http.Request) (store.User, bool) {
	subject, err := s.verifier.Subject(r)
	if err != nil {
		if errors.Is(err, auth.ErrMissingBearerToken) || errors.Is(err, auth.ErrInvalidBearerToken) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="GoreeCloud Maps"`)
			writeAPIError(w, http.StatusUnauthorized, "unauthorized", "A valid GoreeCloud Identity token is required.")
			return store.User{}, false
		}
		s.logger.Error("identity verification failed")
		writeAPIError(w, http.StatusUnauthorized, "unauthorized", "Identity verification failed.")
		return store.User{}, false
	}

	user, err := s.store.ResolveUser(r.Context(), subject)
	if err != nil {
		s.logger.Error("identity subject mapping failed")
		writeAPIError(w, http.StatusInternalServerError, "identity_mapping_failed", "The authenticated account could not be mapped to Maps.")
		return store.User{}, false
	}
	return user, true
}

func (s *Server) writeProviderError(w http.ResponseWriter, err error, operation string) {
	switch {
	case errors.Is(err, providers.ErrNotConfigured):
		writeAPIError(w, http.StatusServiceUnavailable, "provider_not_configured", operation+" is not configured.")
	case errors.Is(err, providers.ErrInvalidRequest):
		writeAPIError(w, http.StatusUnprocessableEntity, "invalid_request", operation+" request is invalid.")
	default:
		s.logger.Warn("map provider operation failed", "operation", operation)
		writeAPIError(w, http.StatusBadGateway, "provider_unavailable", operation+" is temporarily unavailable.")
	}
}

func (s *Server) writeStoreError(w http.ResponseWriter, err error, operation string) {
	switch {
	case errors.Is(err, store.ErrInvalidInput):
		writeAPIError(w, http.StatusUnprocessableEntity, "invalid_request", "The request is invalid.")
	case errors.Is(err, store.ErrNotFound):
		writeAPIError(w, http.StatusNotFound, "not_found", "The requested Maps resource was not found.")
	case errors.Is(err, store.ErrForbidden):
		writeAPIError(w, http.StatusForbidden, "forbidden", "The authenticated user cannot perform this operation.")
	case errors.Is(err, store.ErrConflict):
		writeAPIError(w, http.StatusConflict, "revision_conflict", "The resource changed. Reload it and retry with the current revision.")
	default:
		s.logger.Error("maps storage operation failed", "operation", operation)
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "The Maps operation could not be completed.")
	}
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, maxBytes int64, payload any, invalidMessage string) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", invalidMessage)
		return false
	}
	if err := ensureJSONEOF(decoder); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Only one JSON object is allowed.")
		return false
	}
	return true
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var trailing any
	err := decoder.Decode(&trailing)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err == nil {
		return errors.New("unexpected second JSON value")
	}
	return err
}

func validUUIDString(value string) bool {
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return false
	}
	for index := 0; index < len(value); index++ {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			continue
		}
		character := value[index]
		if !((character >= '0' && character <= '9') || (character >= 'a' && character <= 'f') || (character >= 'A' && character <= 'F')) {
			return false
		}
	}
	return true
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	var payload apiError
	payload.Error.Code = code
	payload.Error.Message = message
	writeJSONStatus(w, status, payload)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	writeJSONStatus(w, status, payload)
}

func writeJSONStatus(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return
	}
}
