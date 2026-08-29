package httpapi

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/auth"
	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/store"
)

// NewSavedPlaces exposes the owner-scoped saved-place surface as a composable
// subrouter. This keeps the feature isolated from unrelated provider routes
// while sharing the same Identity verifier and RLS-backed store.
func NewSavedPlaces(logger *slog.Logger, dataStore *store.Store, verifier *auth.Verifier) http.Handler {
	server := &Server{
		logger:   logger,
		store:    dataStore,
		verifier: verifier,
		mux:      http.NewServeMux(),
	}
	server.registerSavedPlaceRoutes()
	return server.securityHeaders(server.mux)
}

func (s *Server) registerSavedPlaceRoutes() {
	s.mux.HandleFunc("GET /api/v1/saved-places", s.listSavedPlaces)
	s.mux.HandleFunc("POST /api/v1/saved-places", s.createSavedPlace)
	s.mux.HandleFunc("PATCH /api/v1/saved-places/{savedPlaceID}", s.updateSavedPlace)
	s.mux.HandleFunc("DELETE /api/v1/saved-places/{savedPlaceID}", s.deleteSavedPlace)
}

func (s *Server) listSavedPlaces(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
	if !ok {
		return
	}
	places, err := s.store.ListSavedPlaces(r.Context(), user.ID)
	if err != nil {
		s.writeStoreError(w, err, "list saved places")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"savedPlaces": places})
}

func (s *Server) createSavedPlace(w http.ResponseWriter, r *http.Request) {
	user, ok := s.authenticatedUser(w, r)
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
	}
	if !decodeJSONBody(w, r, 64<<10, &payload, "A valid saved place payload is required.") {
		return
	}
	place, err := s.store.CreateSavedPlace(r.Context(), user.ID, store.CreateSavedPlaceInput{
		Provider:        payload.Provider,
		ProviderPlaceID: payload.ProviderPlaceID,
		Name:            payload.Name,
		Address:         payload.Address,
		Latitude:        payload.Latitude,
		Longitude:       payload.Longitude,
		Note:            payload.Note,
	})
	if err != nil {
		s.writeStoreError(w, err, "create saved place")
		return
	}
	writeJSON(w, http.StatusCreated, place)
}

func (s *Server) updateSavedPlace(w http.ResponseWriter, r *http.Request) {
	user, savedPlaceID, ok := s.savedPlaceRequest(w, r)
	if !ok {
		return
	}
	var payload struct {
		Name             *string  `json:"name"`
		Address          *string  `json:"address"`
		Latitude         *float64 `json:"latitude"`
		Longitude        *float64 `json:"longitude"`
		Note             *string  `json:"note"`
		ExpectedRevision int64    `json:"expectedRevision"`
	}
	if !decodeJSONBody(w, r, 64<<10, &payload, "A valid saved place update is required.") {
		return
	}
	place, err := s.store.UpdateSavedPlace(r.Context(), user.ID, savedPlaceID, store.UpdateSavedPlaceInput{
		Name:             payload.Name,
		Address:          payload.Address,
		Latitude:         payload.Latitude,
		Longitude:        payload.Longitude,
		Note:             payload.Note,
		ExpectedRevision: payload.ExpectedRevision,
	})
	if err != nil {
		s.writeStoreError(w, err, "update saved place")
		return
	}
	writeJSON(w, http.StatusOK, place)
}

func (s *Server) deleteSavedPlace(w http.ResponseWriter, r *http.Request) {
	user, savedPlaceID, ok := s.savedPlaceRequest(w, r)
	if !ok {
		return
	}
	expectedRevision, err := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("expectedRevision")), 10, 64)
	if err != nil || expectedRevision < 1 {
		writeAPIError(w, http.StatusBadRequest, "invalid_revision", "A positive expectedRevision query parameter is required.")
		return
	}
	if err := s.store.DeleteSavedPlace(r.Context(), user.ID, savedPlaceID, expectedRevision); err != nil {
		s.writeStoreError(w, err, "delete saved place")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) savedPlaceRequest(w http.ResponseWriter, r *http.Request) (store.User, string, bool) {
	user, ok := s.authenticatedUser(w, r)
	if !ok {
		return store.User{}, "", false
	}
	savedPlaceID := strings.TrimSpace(r.PathValue("savedPlaceID"))
	if !validUUIDString(savedPlaceID) {
		writeAPIError(w, http.StatusBadRequest, "invalid_resource_id", "A valid saved place ID is required.")
		return store.User{}, "", false
	}
	return user, savedPlaceID, true
}
