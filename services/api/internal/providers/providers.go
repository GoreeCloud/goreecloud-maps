package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotConfigured  = errors.New("provider not configured")
	ErrInvalidRequest = errors.New("invalid provider request")
	ErrUpstream       = errors.New("provider upstream failure")
)

const maxProviderResponseBytes = 2 << 20

type Set struct {
	geocoder *endpoint
	router   *endpoint
	client   *http.Client
	clientID string
}

type Capabilities struct {
	Geocoding bool `json:"geocoding"`
	Routing   bool `json:"routing"`
}

type SearchResult struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Label     string  `json:"label"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Category  string  `json:"category,omitempty"`
	Type      string  `json:"type,omitempty"`
}

type RoutePoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RouteRequest struct {
	Mode      string       `json:"mode"`
	Locations []RoutePoint `json:"locations"`
}

type Route struct {
	Mode            string     `json:"mode"`
	DistanceMeters  float64    `json:"distanceMeters"`
	DurationSeconds float64    `json:"durationSeconds"`
	Legs            []RouteLeg `json:"legs"`
}

type RouteLeg struct {
	DistanceMeters  float64         `json:"distanceMeters"`
	DurationSeconds float64         `json:"durationSeconds"`
	Shape           string          `json:"shape"`
	Maneuvers       []RouteManeuver `json:"maneuvers"`
}

type RouteManeuver struct {
	Type              int      `json:"type"`
	Instruction       string   `json:"instruction"`
	VerbalInstruction string   `json:"verbalInstruction,omitempty"`
	DistanceMeters    float64  `json:"distanceMeters"`
	DurationSeconds   float64  `json:"durationSeconds"`
	BeginShapeIndex   int      `json:"beginShapeIndex"`
	EndShapeIndex     int      `json:"endShapeIndex"`
	StreetNames       []string `json:"streetNames,omitempty"`
}

type endpoint struct {
	base *url.URL
}

type nominatimPlace struct {
	PlaceID     int64  `json:"place_id"`
	OSMType     string `json:"osm_type"`
	OSMID       int64  `json:"osm_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Latitude    string `json:"lat"`
	Longitude   string `json:"lon"`
	Category    string `json:"category"`
	Class       string `json:"class"`
	Type        string `json:"type"`
}

func NewSet(geocoderBaseURL, routerBaseURL, clientID string) (Set, error) {
	geocoder, err := parseEndpoint(geocoderBaseURL)
	if err != nil {
		return Set{}, fmt.Errorf("geocoder base URL: %w", err)
	}
	router, err := parseEndpoint(routerBaseURL)
	if err != nil {
		return Set{}, fmt.Errorf("router base URL: %w", err)
	}
	return Set{
		geocoder: geocoder,
		router:   router,
		client: &http.Client{
			Timeout: 6 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		clientID: strings.TrimSpace(clientID),
	}, nil
}

func (s Set) Capabilities() Capabilities {
	return Capabilities{Geocoding: s.geocoder != nil, Routing: s.router != nil}
}

func (s Set) Search(ctx context.Context, query string, limit int, acceptLanguage string) ([]SearchResult, error) {
	if s.geocoder == nil {
		return nil, ErrNotConfigured
	}
	query = strings.TrimSpace(query)
	if query == "" || len([]rune(query)) > 240 {
		return nil, ErrInvalidRequest
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}

	requestURL := s.geocoder.actionURL("search")
	values := requestURL.Query()
	values.Set("q", query)
	values.Set("format", "jsonv2")
	values.Set("addressdetails", "1")
	values.Set("limit", strconv.Itoa(limit))
	requestURL.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: build geocoder request", ErrInvalidRequest)
	}
	s.applyHeaders(req, acceptLanguage)

	response, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: geocoder request", ErrUpstream)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("%w: geocoder status %d", ErrUpstream, response.StatusCode)
	}

	var upstream []nominatimPlace
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxProviderResponseBytes))
	if err := decoder.Decode(&upstream); err != nil {
		return nil, fmt.Errorf("%w: decode geocoder response", ErrUpstream)
	}

	results := make([]SearchResult, 0, len(upstream))
	for _, item := range upstream {
		result, ok := normalizeNominatimPlace(item)
		if ok {
			results = append(results, result)
		}
	}
	return results, nil
}

func (s Set) Reverse(ctx context.Context, latitude, longitude float64, acceptLanguage string) (SearchResult, error) {
	if s.geocoder == nil {
		return SearchResult{}, ErrNotConfigured
	}
	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		return SearchResult{}, ErrInvalidRequest
	}

	requestURL := s.geocoder.actionURL("reverse")
	values := requestURL.Query()
	values.Set("lat", strconv.FormatFloat(latitude, 'f', -1, 64))
	values.Set("lon", strconv.FormatFloat(longitude, 'f', -1, 64))
	values.Set("format", "jsonv2")
	values.Set("addressdetails", "1")
	requestURL.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return SearchResult{}, fmt.Errorf("%w: build reverse geocoder request", ErrInvalidRequest)
	}
	s.applyHeaders(req, acceptLanguage)

	response, err := s.client.Do(req)
	if err != nil {
		return SearchResult{}, fmt.Errorf("%w: reverse geocoder request", ErrUpstream)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return SearchResult{}, fmt.Errorf("%w: reverse geocoder status %d", ErrUpstream, response.StatusCode)
	}

	var upstream nominatimPlace
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxProviderResponseBytes))
	if err := decoder.Decode(&upstream); err != nil {
		return SearchResult{}, fmt.Errorf("%w: decode reverse geocoder response", ErrUpstream)
	}
	result, ok := normalizeNominatimPlace(upstream)
	if !ok {
		return SearchResult{}, fmt.Errorf("%w: invalid reverse geocoder coordinates", ErrUpstream)
	}
	return result, nil
}

func (s Set) Route(ctx context.Context, input RouteRequest) (Route, error) {
	if s.router == nil {
		return Route{}, ErrNotConfigured
	}
	costing, mode, err := normalizeMode(input.Mode)
	if err != nil {
		return Route{}, err
	}
	if len(input.Locations) < 2 || len(input.Locations) > 25 {
		return Route{}, ErrInvalidRequest
	}
	locations := make([]map[string]float64, 0, len(input.Locations))
	for _, location := range input.Locations {
		if location.Latitude < -90 || location.Latitude > 90 || location.Longitude < -180 || location.Longitude > 180 {
			return Route{}, ErrInvalidRequest
		}
		locations = append(locations, map[string]float64{"lat": location.Latitude, "lon": location.Longitude})
	}
	payload := map[string]any{
		"locations": locations,
		"costing":   costing,
		"units":     "kilometers",
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return Route{}, fmt.Errorf("%w: encode route request", ErrInvalidRequest)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.router.actionURL("route").String(), bytes.NewReader(encoded))
	if err != nil {
		return Route{}, fmt.Errorf("%w: build router request", ErrInvalidRequest)
	}
	req.Header.Set("Content-Type", "application/json")
	s.applyHeaders(req, "")

	response, err := s.client.Do(req)
	if err != nil {
		return Route{}, fmt.Errorf("%w: router request", ErrUpstream)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return Route{}, fmt.Errorf("%w: router status %d", ErrUpstream, response.StatusCode)
	}

	var upstream struct {
		Trip struct {
			Summary struct {
				Time   float64 `json:"time"`
				Length float64 `json:"length"`
			} `json:"summary"`
			Legs []struct {
				Summary struct {
					Time   float64 `json:"time"`
					Length float64 `json:"length"`
				} `json:"summary"`
				Shape     string `json:"shape"`
				Maneuvers []struct {
					Type              int      `json:"type"`
					Instruction       string   `json:"instruction"`
					VerbalInstruction string   `json:"verbal_pre_transition_instruction"`
					Time              float64  `json:"time"`
					Length            float64  `json:"length"`
					BeginShapeIndex   int      `json:"begin_shape_index"`
					EndShapeIndex     int      `json:"end_shape_index"`
					StreetNames       []string `json:"street_names"`
				} `json:"maneuvers"`
			} `json:"legs"`
		} `json:"trip"`
	}
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxProviderResponseBytes))
	if err := decoder.Decode(&upstream); err != nil {
		return Route{}, fmt.Errorf("%w: decode router response", ErrUpstream)
	}

	route := Route{
		Mode:            mode,
		DistanceMeters:  upstream.Trip.Summary.Length * 1000,
		DurationSeconds: upstream.Trip.Summary.Time,
		Legs:            make([]RouteLeg, 0, len(upstream.Trip.Legs)),
	}
	for _, upstreamLeg := range upstream.Trip.Legs {
		leg := RouteLeg{
			DistanceMeters:  upstreamLeg.Summary.Length * 1000,
			DurationSeconds: upstreamLeg.Summary.Time,
			Shape:           upstreamLeg.Shape,
			Maneuvers:       make([]RouteManeuver, 0, len(upstreamLeg.Maneuvers)),
		}
		for _, maneuver := range upstreamLeg.Maneuvers {
			leg.Maneuvers = append(leg.Maneuvers, RouteManeuver{
				Type:              maneuver.Type,
				Instruction:       maneuver.Instruction,
				VerbalInstruction: maneuver.VerbalInstruction,
				DistanceMeters:    maneuver.Length * 1000,
				DurationSeconds:   maneuver.Time,
				BeginShapeIndex:   maneuver.BeginShapeIndex,
				EndShapeIndex:     maneuver.EndShapeIndex,
				StreetNames:       maneuver.StreetNames,
			})
		}
		route.Legs = append(route.Legs, leg)
	}
	return route, nil
}

func (s Set) applyHeaders(req *http.Request, acceptLanguage string) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GoreeCloud-Maps")
	if language := strings.TrimSpace(acceptLanguage); language != "" {
		req.Header.Set("Accept-Language", language)
	}
	if s.clientID != "" {
		req.Header.Set("X-Client-Id", s.clientID)
	}
}

func parseEndpoint(raw string) (*endpoint, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, errors.New("must be an absolute HTTP(S) URL")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("must not contain credentials, query parameters, or fragments")
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	return &endpoint{base: parsed}, nil
}

func (e *endpoint) actionURL(action string) *url.URL {
	copyURL := *e.base
	copyURL.Path = strings.TrimRight(copyURL.Path, "/") + "/" + strings.TrimLeft(action, "/")
	return &copyURL
}

func normalizeMode(mode string) (costing string, normalized string, err error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "drive", "driving", "auto", "":
		return "auto", "drive", nil
	case "walk", "walking", "pedestrian":
		return "pedestrian", "walk", nil
	case "bike", "biking", "bicycle":
		return "bicycle", "bicycle", nil
	case "transit", "multimodal":
		return "multimodal", "transit", nil
	default:
		return "", "", ErrInvalidRequest
	}
}

func normalizeNominatimPlace(item nominatimPlace) (SearchResult, bool) {
	latitude, err := strconv.ParseFloat(item.Latitude, 64)
	if err != nil || latitude < -90 || latitude > 90 {
		return SearchResult{}, false
	}
	longitude, err := strconv.ParseFloat(item.Longitude, 64)
	if err != nil || longitude < -180 || longitude > 180 {
		return SearchResult{}, false
	}
	name := strings.TrimSpace(item.Name)
	if name == "" {
		name = firstLabelPart(item.DisplayName)
	}
	category := strings.TrimSpace(item.Category)
	if category == "" {
		category = strings.TrimSpace(item.Class)
	}
	return SearchResult{
		ID:        nominatimID(item.OSMType, item.OSMID, item.PlaceID),
		Name:      name,
		Label:     strings.TrimSpace(item.DisplayName),
		Latitude:  latitude,
		Longitude: longitude,
		Category:  category,
		Type:      strings.TrimSpace(item.Type),
	}, true
}

func nominatimID(osmType string, osmID, placeID int64) string {
	if osmID != 0 && strings.TrimSpace(osmType) != "" {
		return fmt.Sprintf("nominatim:%s:%d", strings.ToLower(strings.TrimSpace(osmType)), osmID)
	}
	return fmt.Sprintf("nominatim:place:%d", placeID)
}

func firstLabelPart(label string) string {
	label = strings.TrimSpace(label)
	if index := strings.Index(label, ","); index >= 0 {
		return strings.TrimSpace(label[:index])
	}
	return label
}
