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

	return server.securityHeaders(server.mux)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if err := s.store.Ping(r.Context()); err != nil {
		s.logger.Error("database readiness check failed", "error", err)
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

	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload providers.RouteRequest
	if err := decoder.Decode(&payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "A valid route payload is required.")
		return
	}
	if err := ensureJSONEOF(decoder); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Only one JSON object is allowed.")
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
		s.logger.Error("collection list failed", "error", err)
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "Collections could not be loaded.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"collections": collections})
}

func (s *Server) createCollection(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 32<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := decoder.Decode(&payload); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "A valid collection payload is required.")
		return
	}
	if err := ensureJSONEOF(decoder); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "Only one JSON object is allowed.")
		return
	}

	collection, err := s.store.CreateCollection(r.Context(), user.ID, store.CreateCollectionInput{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		if strings.Contains(err.Error(), "collection name") {
			writeAPIError(w, http.StatusUnprocessableEntity, "invalid_collection", err.Error())
			return
		}
		s.logger.Error("collection creation failed", "error", err)
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "The collection could not be created.")
		return
	}

	writeJSON(w, http.StatusCreated, collection)
}

func (s *Server) authenticatedUser(w http.ResponseWriter, r *http.Request) (store.User, bool) {
	subject, err := s.verifier.Subject(r)
	if err != nil {
		if errors.Is(err, auth.ErrMissingBearerToken) || errors.Is(err, auth.ErrInvalidBearerToken) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="GoreeCloud Maps"`)
			writeAPIError(w, http.StatusUnauthorized, "unauthorized", "A valid GoreeCloud Identity token is required.")
			return store.User{}, false
		}
		s.logger.Error("identity verification failed", "error", err)
		writeAPIError(w, http.StatusUnauthorized, "unauthorized", "Identity verification failed.")
		return store.User{}, false
	}

	user, err := s.store.ResolveUser(r.Context(), subject)
	if err != nil {
		s.logger.Error("identity subject mapping failed", "error", err)
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
