package httpapi

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GoreeCloud/goreecloud-maps/services/api/internal/providers"
)

func TestHealthEndpoint(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := New(logger, nil, nil, providers.Set{})

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if contentType := response.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("expected JSON content type, got %q", contentType)
	}
	if cacheControl := response.Header().Get("Cache-Control"); cacheControl != "no-store" {
		t.Fatalf("expected no-store cache control, got %q", cacheControl)
	}
	if body := response.Body.String(); !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("unexpected health response body: %s", body)
	}
}

func TestCapabilitiesEndpointReportsUnconfiguredProviders(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := New(logger, nil, nil, providers.Set{})

	request := httptest.NewRequest(http.MethodGet, "/api/v1/capabilities", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	body := response.Body.String()
	if !strings.Contains(body, `"geocoding":false`) || !strings.Contains(body, `"routing":false`) {
		t.Fatalf("unexpected capability response body: %s", body)
	}
}
