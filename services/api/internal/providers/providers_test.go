package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchNormalizesNominatimResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			t.Fatalf("expected /search, got %q", r.URL.Path)
		}
		if query := r.URL.Query().Get("q"); query != "Union Station" {
			t.Fatalf("unexpected query %q", query)
		}
		if format := r.URL.Query().Get("format"); format != "jsonv2" {
			t.Fatalf("unexpected format %q", format)
		}
		if language := r.Header.Get("Accept-Language"); language != "en-US" {
			t.Fatalf("unexpected language %q", language)
		}
		if clientID := r.Header.Get("X-Client-Id"); clientID != "maps-test" {
			t.Fatalf("unexpected client id %q", clientID)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"place_id":42,"osm_type":"way","osm_id":99,"name":"Union Station","display_name":"Union Station, Test City","lat":"41.881","lon":"-87.640","category":"railway","type":"station"}]`))
	}))
	defer server.Close()

	providers, err := NewSet(server.URL, "", "maps-test")
	if err != nil {
		t.Fatalf("NewSet failed: %v", err)
	}
	results, err := providers.Search(context.Background(), "Union Station", 5, "en-US")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected one result, got %d", len(results))
	}
	result := results[0]
	if result.ID != "nominatim:way:99" || result.Name != "Union Station" || result.Category != "railway" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if result.Latitude != 41.881 || result.Longitude != -87.640 {
		t.Fatalf("unexpected coordinates: %#v", result)
	}
}

func TestSearchRequiresConfiguredGeocoder(t *testing.T) {
	providers, err := NewSet("", "", "")
	if err != nil {
		t.Fatalf("NewSet failed: %v", err)
	}
	if _, err := providers.Search(context.Background(), "test", 10, ""); err != ErrNotConfigured {
		t.Fatalf("expected ErrNotConfigured, got %v", err)
	}
}

func TestReverseNormalizesNominatimResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reverse" {
			t.Fatalf("expected /reverse, got %q", r.URL.Path)
		}
		if r.URL.Query().Get("lat") != "41.881" || r.URL.Query().Get("lon") != "-87.64" {
			t.Fatalf("unexpected reverse coordinates: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"place_id":7,"osm_type":"node","osm_id":8,"display_name":"100 Main Street, Test City","lat":"41.8811","lon":"-87.6401","category":"place","type":"house"}`))
	}))
	defer server.Close()

	providers, err := NewSet(server.URL, "", "")
	if err != nil {
		t.Fatalf("NewSet failed: %v", err)
	}
	result, err := providers.Reverse(context.Background(), 41.881, -87.64, "")
	if err != nil {
		t.Fatalf("Reverse failed: %v", err)
	}
	if result.ID != "nominatim:node:8" || result.Name != "100 Main Street" {
		t.Fatalf("unexpected reverse result: %#v", result)
	}
}

func TestRouteNormalizesValhallaResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/route" || r.Method != http.MethodPost {
			t.Fatalf("unexpected route request: %s %s", r.Method, r.URL.Path)
		}
		var request struct {
			Costing   string `json:"costing"`
			Locations []struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"locations"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.Costing != "pedestrian" || len(request.Locations) != 2 {
			t.Fatalf("unexpected route request: %#v", request)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"trip":{"summary":{"time":120,"length":1.25},"legs":[{"summary":{"time":120,"length":1.25},"shape":"encoded-shape","maneuvers":[{"type":1,"instruction":"Walk north","verbal_pre_transition_instruction":"Walk north","time":30,"length":0.25,"begin_shape_index":0,"end_shape_index":3,"street_names":["Main Street"]}]}]}}`))
	}))
	defer server.Close()

	providers, err := NewSet("", server.URL, "")
	if err != nil {
		t.Fatalf("NewSet failed: %v", err)
	}
	route, err := providers.Route(context.Background(), RouteRequest{
		Mode: "walk",
		Locations: []RoutePoint{
			{Latitude: 41.88, Longitude: -87.64},
			{Latitude: 41.89, Longitude: -87.63},
		},
	})
	if err != nil {
		t.Fatalf("Route failed: %v", err)
	}
	if route.Mode != "walk" || route.DistanceMeters != 1250 || route.DurationSeconds != 120 {
		t.Fatalf("unexpected route summary: %#v", route)
	}
	if len(route.Legs) != 1 || route.Legs[0].Shape != "encoded-shape" || len(route.Legs[0].Maneuvers) != 1 {
		t.Fatalf("unexpected route legs: %#v", route.Legs)
	}
	if route.Legs[0].Maneuvers[0].DistanceMeters != 250 {
		t.Fatalf("unexpected maneuver: %#v", route.Legs[0].Maneuvers[0])
	}
}

func TestNewSetRejectsCredentialedEndpoint(t *testing.T) {
	_, err := NewSet("https://user:secret@example.com", "", "")
	if err == nil || !strings.Contains(err.Error(), "must not contain credentials") {
		t.Fatalf("expected credential rejection, got %v", err)
	}
}
