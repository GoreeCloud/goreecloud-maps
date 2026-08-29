# GoreeCloud Maps Provider Contracts

## Status

GoreeCloud Maps contains initial server-side provider adapters for place geocoding and route planning plus an initial public geographic-data release/delivery contract. These components establish replaceable application contracts; they do **not** mean a production geographic-data provider is configured or accepted.

No geocoder, router, or production basemap release is configured by default. A clone or ordinary development build therefore does not silently send place searches, coordinates, route requests, or tile requests to a public third-party service.

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

## Server-provider configuration

The API accepts optional server-side configuration:

- `MAPS_GEOCODER_BASE_URL`
- `MAPS_ROUTER_BASE_URL`
- `MAPS_PROVIDER_CLIENT_ID`

Provider URLs must be absolute HTTP(S) URLs and cannot contain embedded credentials, query strings, or fragments. The API appends only its fixed provider actions. The optional provider client ID is non-secret metadata and is forwarded as `X-Client-Id` when configured.

Provider credentials, should a future approved provider require them, must use a separate server-side secret contract rather than being embedded in a URL or exposed to the browser.

## Public map-data release contract

The base-map delivery boundary is separate from geocoding/routing and separate from private user data.

`mapdata/release-manifest.schema.json` and `scripts/validate-mapdata-release.mjs` define the current source contract for a published geographic release. A release identifies:

- immutable release ID;
- generation time;
- MapLibre style path;
- vector-tile path template;
- supported zoom range and geographic bounds;
- user-visible attribution;
- source dataset/version and licensing/provenance information;
- confirmation that the public bundle contains public geographic data only.

The object layout is versioned beneath `releases/<release-id>/`. Published release objects must not be changed in place. `manifests/current.json` is the intentionally mutable pointer used for controlled roll-forward/rollback.

### Style schema v1

`scripts/validate-mapdata-style.mjs` validates the release style before publication and the browser independently repeats the security-relevant resource checks before rendering.

Schema v1 requires:

- MapLibre style version 8;
- vector/MVT sources only;
- exactly the release-local `tiles/{z}/{x}/{y}.pbf` template rather than arbitrary tile URLs;
- no TileJSON `url` indirection;
- no imported styles;
- no external `font-faces` resources;
- release-local glyph and sprite resource paths when those features are used;
- declared vector source/source-layer relationships for non-background layers.

The browser's `VITE_MAP_DATA_MANIFEST_URL` flow fetches the current manifest and immutable style without credentials and refuses redirects. It bounds both documents, verifies manifest source/provenance/public-data state, normalizes allowed tile/glyph/sprite resources to the configured immutable release origin, rejects off-origin/off-release resources, injects validated manifest attribution, and gives MapLibre the normalized in-memory style. A rejected release retains the repository-local empty map.

This v1 contract intentionally does not admit raster/aerial imagery, terrain/raster-dem sources, PMTiles or other archive formats, imported style fragments, or additional source types. Those capabilities require future versioned contracts instead of weakening v1 implicitly.

### Read-only edge delivery

`services/mapdata-edge` defines the first read-only Cloudflare Worker/R2 source implementation. Its request surface is allowlisted to:

- current release manifest;
- immutable release manifest;
- immutable style JSON;
- immutable vector-tile PBF objects;
- immutable sprite JSON/PNG objects;
- immutable glyph PBF objects.

The Worker has no write method, bucket listing, arbitrary object-key proxy, private Maps API, authentication/session handling, search/route proxy, or user-data surface. Versioned release objects receive immutable cache semantics; the current manifest receives a short cache lifetime. The current source sends public CORS headers because these objects are intended public geographic assets, not private user resources; production origin/CORS policy still requires deployed acceptance.

The checked-in R2 names and Worker name are desired deployment resources, not proof that Cloudflare resources have been provisioned. See `deployment/cloudflare/README.md` and `mapdata/README.md`.

## Initial safety and privacy controls

The implemented provider/data-plane source boundaries include:

- no hard-coded public demo-provider origin;
- fixed geocoder/router action paths rather than caller-controlled upstream paths;
- provider HTTP redirects are refused rather than automatically followed;
- a six-second provider HTTP timeout;
- a two-MiB maximum decoded provider response stream;
- bounded search-result limits and route waypoint counts;
- coordinate and query validation;
- rejection of credentialed/query-bearing provider base URLs;
- provider origins excluded from public API responses;
- provider failure logs that record the operation class without place queries, coordinates, route endpoints, or provider URLs;
- no public map-data writes or object listing through the edge gateway;
- strict map-data URL-path shapes and tile-coordinate bounds;
- separate cache classes for current release metadata versus immutable release objects;
- machine-readable attribution/provenance requirements for future map-data releases;
- credentialless browser map-data requests;
- bounded manifest/style responses and redirect refusal;
- producer- and consumer-side style resource constraints that prevent release style JSON from redirecting MapLibre to unrelated tile/font/sprite origins;
- fail-closed local empty-style behavior for invalid/unavailable configured releases.

These controls reduce accidental disclosure and configuration abuse, but they are not a complete production egress, ingestion, CDN, data-license, or SSRF/security control plane.

## Production acceptance still required

Before geographic/provider traffic can be approved for production Maps use, acceptance must cover at least:

- GoreeCloud-controlled or explicitly approved endpoints/resources;
- network egress restrictions and DNS behavior appropriate to server runtimes;
- source licensing, attribution, provenance, retention, and redistribution terms;
- reproducible map-data ingestion/build pipeline and update cadence;
- release-manifest and style validation on the actual generated output;
- geographic coverage, rendering/label quality, and accessibility validation;
- provider/storage capacity, quotas, rate limiting, abuse controls, and failure behavior;
- cache, release-pointer, rollback, and stale-release behavior;
- production CORS/origin behavior for the public map-data service;
- secret management when credentials are required;
- observability that avoids sensitive searches and precise personal coordinates;
- public/private-data separation;
- Privacy Shield and Wardveil Security evidence;
- degraded/offline presentation in the Maps client;
- actual Cloudflare Pages/Worker/R2 deployment and runtime acceptance where those services are used.

Compatibility with Nominatim or Valhalla does not grant permission to use an arbitrary public instance or demo service for GoreeCloud production traffic. Likewise, a syntactically valid map-data release does not authorize publication of an unapproved dataset.
