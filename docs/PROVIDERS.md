# GoreeCloud Maps Provider Contracts

## Status

GoreeCloud Maps now contains initial server-side provider adapters for place geocoding and route planning. These adapters establish replaceable application contracts; they do **not** mean a production geographic-data provider is configured or accepted.

No provider endpoint is configured by default. A clone or ordinary development build therefore does not silently send place searches, coordinates, or route requests to a public third-party service.

## Geocoding contract

The initial geocoding adapter is compatible with a Nominatim-style HTTP API.

Supported upstream actions:

- `/search` for forward geocoding and place/address discovery;
- `/reverse` for coordinate-to-place lookup.

Provider-specific fields are normalized before they cross the Maps API boundary. The first-party result model contains a stable provider reference, display name/label, coordinates, category, and type. Raw upstream payloads are not exposed as the application contract.

## Routing contract

The initial routing adapter is compatible with a Valhalla-style `/route` service.

Supported application modes currently map to:

- drive → `auto`;
- walk → `pedestrian`;
- bicycle → `bicycle`;
- transit → `multimodal`.

The Maps response normalizes route and leg distance in meters, duration in seconds, Valhalla encoded route shape, maneuver instructions, maneuver distances/durations, street names, and shape indexes. The encoded shape remains an interchange detail of the normalized route contract and can be decoded by clients that render route geometry.

## Maps API surface

The initial provider-facing Maps API routes are:

- `GET /api/v1/capabilities` — reports whether geocoding and routing providers are configured;
- `GET /api/v1/search?q=...&limit=...` — authenticated forward geocoding/place search;
- `GET /api/v1/reverse?latitude=...&longitude=...` — authenticated reverse geocoding;
- `POST /api/v1/routes` — authenticated route planning.

The capability endpoint exposes only configured/not-configured state. It does not disclose provider origins or credentials.

## Configuration

The API accepts optional server-side configuration:

- `MAPS_GEOCODER_BASE_URL`
- `MAPS_ROUTER_BASE_URL`
- `MAPS_PROVIDER_CLIENT_ID`

Provider URLs must be absolute HTTP(S) URLs and cannot contain embedded credentials, query strings, or fragments. The API appends only its fixed provider actions. The optional provider client ID is non-secret metadata and is forwarded as `X-Client-Id` when configured.

Provider credentials, should a future approved provider require them, must use a separate server-side secret contract rather than being embedded in a URL or exposed to the browser.

## Initial safety and privacy controls

The implemented source boundary includes:

- no hard-coded public demo-provider origin;
- fixed provider action paths rather than caller-controlled upstream paths;
- provider HTTP redirects are refused rather than automatically followed;
- a six-second provider HTTP timeout;
- a two-MiB maximum decoded provider response stream;
- bounded search-result limits and route waypoint counts;
- coordinate and query validation;
- rejection of credentialed/query-bearing provider base URLs;
- provider origins excluded from public API responses;
- provider failure logs that record the operation class without place queries, coordinates, route endpoints, or provider URLs.

These controls reduce accidental disclosure and configuration abuse, but they are not a complete production egress or SSRF control plane.

## Production acceptance still required

Before a provider can be approved for production Maps traffic, acceptance must cover at least:

- GoreeCloud-controlled or explicitly approved provider endpoints;
- network egress restrictions and DNS behavior appropriate to the runtime;
- source licensing, attribution, provenance, retention, and redistribution terms;
- provider capacity, quotas, rate limiting, abuse controls, and failure behavior;
- secret management when credentials are required;
- observability that avoids sensitive searches and precise coordinates;
- Privacy Shield and Wardveil Security evidence;
- geographic-quality and accessibility validation for intended coverage;
- degraded/offline presentation in the Maps client.

Compatibility with Nominatim or Valhalla does not grant permission to use an arbitrary public instance or demo service for GoreeCloud production traffic.
