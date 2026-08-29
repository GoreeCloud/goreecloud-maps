# GoreeCloud Maps Specifications

## 1. Status and authority

GoreeCloud Maps is in **Development**. This specification describes the required product and engineering direction; it does not assert that every listed capability is production implemented.

The application targets Glaze UI 2.0.0 Stable and must complete application-specific acceptance before any Stable or production-ready claim.

## 2. Product purpose

Maps provides first-party map exploration, place discovery, geocoding, directions, navigation, saved map content, collaborative map resources, offline maps, and map-related interoperability across GoreeCloud.

Maps does not replace GoreeCloud Location. Location remains authoritative for personal location history, background tracking, device location, geofences, Find My data, and location-sharing permission state.

## 3. Supported users and tenancy

The data model must support multiple authenticated GoreeCloud Identity principals. Every user-owned object must have an explicit owner or tenant boundary.

Required sharing scopes:

- private to one user;
- shared with explicit named members;
- household/group scoped through approved identity/group capabilities;
- link-based sharing only when explicitly enabled and revocable;
- public publishing only as a separately governed future capability.

Collaborative resources use explicit roles: `owner`, `editor`, and `viewer`. Authorization must be checked server-side for every protected read or mutation. Possession of an object identifier must never grant access.

The current source authorization model has automated PostGIS integration coverage for owner/editor/viewer/stranger collection visibility, collection revision conflicts, editor item creation/update/deletion, viewer mutation denial, owner-only membership administration, role changes, member self-removal, private saved-place isolation, collaboration audit records, forged-audit rejection, immutable collection ownership, and refusal of a database-owner/`BYPASSRLS` runtime connection. This is source-level CI acceptance; production database, GoreeCloud Identity SSO, deployment, load, backup, and recovery acceptance remain separate gates.

## 4. Core domains

### Map presentation

The map engine must support smooth pan, zoom, rotate, pitch, user-location display, selectable layers, map style switching, scalable vector rendering, high-DPI output, and graceful degradation when advanced GPU effects are unavailable.

Target capabilities include globe mode, terrain, 3D buildings, indoor layers, accessibility overlays, transit overlays, traffic/incident overlays, and satellite/aerial imagery when an approved provider and licensing model exist.

The current web renderer remains intentionally useful in a data-empty development mode. `VITE_MAP_DATA_MANIFEST_URL` is the preferred live-release seam: when configured, the browser starts with the repository-local empty style, validates the public release manifest and immutable release style, and only then loads the normalized in-memory MapLibre style. `VITE_MAP_STYLE_URL` remains a manual/legacy approved-style seam when the manifest setting is empty. Maps must not silently fall back to a public map provider.

### Public geographic data releases

Public geographic data is a separate delivery domain from private/user Maps state.

A map-data release uses an immutable release identifier with a short-cached mutable current pointer and immutable release objects for manifest, style, vector tiles, glyphs, and sprites. The manifest records release identity/generation time, geographic bounds, zoom range, style/tile paths, attribution, source datasets/versions/licenses/provenance, and a required public-geographic-data-only marker.

Schema v1 is deliberately restricted to MapLibre Style Specification v8 and vector/MVT sources. Producer and browser validation require release-local `tiles/{z}/{x}/{y}.pbf`, disallow TileJSON indirection and style imports, disallow external font-face resources, and constrain glyph/sprite references to the immutable release path. The browser normalizes allowed release resources to the configured release origin and injects validated manifest attribution before MapLibre renders the style.

A malformed, oversized, redirected, off-origin, off-release, unsupported, or unavailable configured release must fail closed to the local empty style. It must not cause automatic use of an unrelated public tile/style provider.

Schema v1 is not the permanent format for every map capability. Raster/aerial imagery, terrain/raster-dem sources, PMTiles or other archive formats, future style imports, and additional source types require explicitly versioned later contracts plus their own licensing, privacy, security, performance, caching, offline, and deployment acceptance.

The current repository contains only synthetic manifest/style fixtures. No generated basemap dataset or live geographic coverage is accepted by this specification state.

### Places and discovery

The application must support forward and reverse geocoding, category search, nearby discovery, place cards, hours/contact metadata when licensed, accessibility metadata, saved places, favorites, user notes, collections, and place-data correction feedback.

Search results must identify provider/source provenance where required and must not silently combine incompatible licensing terms.

The source implementation includes authenticated forward and reverse geocoding API routes backed by a normalized Nominatim-compatible adapter. The web search form and initial category actions call the same-origin Maps API after authentication, render normalized results, and can center/mark a selected result on the map. Category actions currently submit ordinary search text rather than a dedicated nearby/POI ranking contract. No live geocoder endpoint, geographic-quality acceptance, GoreeCloud Search interoperability, rich place-card provider, or production provider approval is claimed.

### Directions and navigation

The routing boundary must support driving, walking, cycling, transit, and multimodal itineraries where data coverage exists. Route responses should support alternatives, ETAs, distance/duration, maneuvers, geometry, accessibility constraints, avoidances, departure/arrival time, and route warnings.

The source implementation includes an authenticated route-planning API and a normalized Valhalla-compatible adapter for drive, walk, bicycle, and transit/multimodal costing. It currently normalizes route/leg distance and duration, encoded route shape, maneuvers, street names, and shape indexes.

The web source includes an initial directions experience that geocodes an origin/destination, requests a route, decodes returned polyline6 geometry, renders the route when the map is ready, fits the viewport, and presents a bounded maneuver summary. This is not turn-by-turn navigation. No approved live router, traffic-aware ETA, route alternatives, advanced constraints, route-quality acceptance, location runtime, or real-device navigation acceptance is claimed.

Turn-by-turn navigation is a separate acceptance boundary requiring runtime GPS integration through GoreeCloud Location, rerouting, off-route detection, maneuver guidance, background/lock-screen behavior where supported, and verified real-device testing.

### Offline

Users must be able to select or download versioned regional map packages. Offline design must account for tiles, styles, glyphs/icons, place indexes, and routing graphs where supported. Packages require integrity metadata, versioning, quota management, update checks, rollback, and clean deletion.

The current public release manifest establishes release identity, generation time, coverage, source/provenance and immutable resource paths that can later contribute to package freshness/integrity decisions, but it is not itself a downloadable/offline-region implementation.

### Collaboration

Users must be able to create collections and shared maps, add places, annotate entries, reorder content, invite/remove members, alter member roles, and revoke sharing. Conflict handling must be deterministic and auditable.

The current source implements authenticated collection updates with expected revisions, member listing/addition/role changes/removal, and collection-item list/create/update/delete. Owner/editor/viewer authorization is enforced by application role checks and PostgreSQL row-level security. Collection and item edits use optimistic revision conflicts, and implemented membership/item mutations emit collection audit events.

The web source can list/create collections and browse their items after authentication. Full membership management, item editing, collection reordering, annotations, share links, ownership transfer, and invitation acceptance UI remain incomplete.

Human-friendly recipient discovery and invitation delivery are **blocked by prerequisite** until GoreeCloud Identity defines an approved cross-application directory/invitation contract. Maps must not use an administrative Identity-provider user-list API as an accidental consumer directory and must not create an unrelated Maps account directory.

## 5. Data architecture

Preferred first-party state store: PostgreSQL with PostGIS.

Minimum entities:

- users / identity subjects;
- preferences;
- saved places;
- collections;
- collection members;
- collection entries;
- route plans;
- map annotations;
- provider-source references;
- offline package manifests;
- contribution/feedback records;
- audit records for security-sensitive sharing changes.

Location history must not be copied into Maps merely for convenience. Maps stores only the minimum user-location-derived state required for map tasks, subject to Privacy Shield rules.

The migration chain currently includes the multi-user/PostGIS foundation plus audit-event RLS hardening. CI must continue applying the complete ordered migration chain rather than testing only the first migration.

Public geographic releases must not be stored as if they were user-owned Maps records. Public tile/style/glyph/sprite distribution and private PostGIS application state require separate authorization, caching, retention, backup, and operational policies.

## 6. Provider architecture

Maps must expose internal interfaces for:

- base map/vector tiles;
- raster/satellite imagery;
- geocoding and reverse geocoding;
- place search and POI enrichment;
- routing;
- transit;
- traffic and incidents;
- elevation/terrain;
- street-level imagery;
- indoor maps.

No provider may be assumed permanent. Provider-specific response types must be normalized at the adapter boundary.

The preferred baseline uses open/self-hostable components and GoreeCloud-controlled delivery. Any external service must be justified by coverage, quality, licensing, privacy, security, operational, and portability requirements.

The first implemented server adapters are Nominatim-compatible forward/reverse geocoding and Valhalla-compatible routing. Their base URLs are optional server-side configuration and are blank by default. Configured URLs must be absolute HTTP(S) URLs without embedded credentials, query parameters, or fragments; Maps appends fixed provider action paths. The implementation bounds HTTP timeout, response size, search limits, waypoint counts, and coordinate/query validation, and refuses upstream HTTP redirects.

The public map-data provider boundary is release-oriented rather than exposing an arbitrary client-configurable tile origin. `mapdata/` defines the v1 release contract; `services/mapdata-edge` defines a read-only Worker/R2 source gateway for allowlisted public release objects. The web client independently validates the manifest/style boundary before rendering.

This application-level URL/resource validation and redirect refusal are not a complete production network defense. Production acceptance must add runtime egress restrictions where relevant, DNS-resolution controls for server providers, production CORS/origin policy for public map data, provider/dataset provenance/license terms, capacity/rate limits, secret handling where applicable, monitoring, Privacy Shield evidence, and Wardveil Security evidence.

Compatibility with an API or release shape does not authorize production use of an arbitrary public Nominatim, Valhalla, map tile, imagery, or other provider instance.

## 7. Platform integration

### GoreeCloud Identity

Identity is the principal/authentication authority. Maps must not create an unrelated permanent account system.

The browser source implements an optional OIDC public client using Authorization Code + PKCE (`S256`) without a client secret. OIDC issuer/client registration is configuration-driven and blank by default. The browser keeps the resulting bearer access token in memory; only transient PKCE verifier/state are placed in `sessionStorage`. Configured redirect URIs must remain on the current Maps origin, and configured Identity endpoints must use HTTPS outside local development.

Protected Maps API routes expect a bearer **access token**. The current verifier checks the signed token against the configured issuer/client audience and standard expiry semantics, requires a subject, then asks the provider UserInfo endpoint to accept the same token and return the same subject. Maps resolves that subject to an internal user and independently performs resource authorization.

The UserInfo confirmation is the current source contract for distinguishing a provider-recognized access token from the earlier ID-token-as-bearer assumption. It introduces a provider availability/latency dependency on protected requests and is not automatically the final production validation architecture.

GoreeCloud Identity is not yet approved for GoreeCloud-wide production SSO. Maps still requires actual client registration and end-to-end evidence for login, callback, user mapping, token expiry, logout/session behavior, disabled accounts, Identity outage behavior, recovery, rollback, privacy, and security before production acceptance.

See `docs/IDENTITY.md` for the detailed source boundary.

### GoreeCloud Location

Maps requests current/approved location capabilities and optional sharing overlays from Location. Permission, precision, tracking, and history truth remain with Location/Privacy Shield.

### GoreeCloud Search

Place and geographic discovery should interoperate with Search through documented APIs/events rather than duplicating universal search responsibilities.

### GoreeCloud Mesh

Cross-application coordination, capability discovery, and governed events should use Mesh when the applicable contracts are available.

### Privacy Shield

Privacy state must be evidence-backed. Requirements include data minimization, explicit sharing, revocation, clear precision controls, private-history options, provider disclosure where needed, and no hidden tracking.

Public map-data releases must contain public geographic data only. Private searches, routes, saved places, collections, Identity state, and precise personal location are prohibited from the public release bundle.

### Wardveil Security

Required controls include authenticated protected APIs, authorization at every resource boundary, rate limiting/abuse controls, input validation, SSRF-resistant provider access, constrained public-release resource loading, secure secret handling, dependency review, and security-event evidence.

### Everkeep

User-owned saved places, collections, annotations, preferences, and other durable personal map data must have export, backup, restore, portability, and deletion semantics appropriate to the data class. Public immutable map releases have a separate continuity/rollback model and must not be confused with user-data backup semantics.

## 8. Glaze UI 2.0 mapping

The map itself is the primary Canvas. Persistent or contextual controls use Soft Glaze/Glaze; expanded menus and sheets use Deep Glaze; active navigation and other ongoing processes may use Live Glaze.

Mobile composition uses a high-information map/viewing zone with primary reachable actions near the lower action zone. Search expands from its invoking control. Place/search results, collections, route details, and status surfaces must use the same Glaze semantic variables and accessibility fallbacks as the rest of Maps rather than introducing a separate design authority.

Desktop uses a denser sidebar/inspector model, keyboard shortcuts, pointer states, resizable panels, and contextual menus. Tablet/foldable layouts use split panes and hinge-aware placement. Reduced transparency/motion and effects-free modes must preserve complete usability.

## 9. Web/API and public map-data integration principles

The initial browser API client uses a same-origin API path, defaulting to `/api/v1`. `VITE_MAPS_API_BASE_PATH` must begin with `/` and must not be an external or scheme-relative origin. This preserves a controlled reverse-proxy boundary and avoids runtime browser configuration that can redirect private bearer/API traffic to arbitrary origins.

The public map-data release is a distinct browser boundary. `VITE_MAP_DATA_MANIFEST_URL` may identify an accepted same-origin or separate public geographic-data origin, but it must use HTTPS outside localhost, cannot contain credentials/query/fragment state, and is fetched without credentials with redirects refused. A configured release origin must not receive private Maps API bearer tokens.

All protected APIs are versioned under `/api/v1/` initially. APIs use stable resource identifiers, explicit pagination where needed, structured errors, idempotency where retries are expected, optimistic concurrency or version fields for collaborative mutations, and server-side authorization.

Provider credentials and internal upstream endpoints must never be exposed to untrusted clients.

Current source routes include:

- `GET /api/v1/me`;
- `GET /api/v1/collections`;
- `POST /api/v1/collections`;
- `PATCH /api/v1/collections/{collectionID}`;
- `GET /api/v1/collections/{collectionID}/members`;
- `POST /api/v1/collections/{collectionID}/members`;
- `PATCH /api/v1/collections/{collectionID}/members/{memberUserID}`;
- `DELETE /api/v1/collections/{collectionID}/members/{memberUserID}`;
- `GET /api/v1/collections/{collectionID}/items`;
- `POST /api/v1/collections/{collectionID}/items`;
- `PATCH /api/v1/collections/{collectionID}/items/{itemID}`;
- `DELETE /api/v1/collections/{collectionID}/items/{itemID}` with a positive `expectedRevision` query parameter;
- `GET /api/v1/search`;
- `GET /api/v1/reverse`;
- `POST /api/v1/routes`;
- public `GET /api/v1/capabilities`, limited to configured/not-configured provider state.

Collection and item updates require expected revisions and return conflict state when the caller is stale. Member roles are limited to editor/viewer because the owner is represented on the collection itself. Member addition currently accepts an existing Maps user ID; this is a controlled API primitive, not a completed invitation experience.

The provider-dependent search/reverse/route routes require authenticated Maps users and return explicit unavailable state when the capability is not configured.

## 10. Security and privacy constraints

Precise location is sensitive. Logs must avoid recording precise coordinates, route origins/destinations, search queries, bearer tokens, OIDC codes, PKCE verifiers, share secrets, or private collection contents unless a specific operational requirement exists and the retention/privacy contract permits it.

Public map tiles and user-private resources must use separate caching rules. Private responses must not become publicly cacheable through CDN configuration. Public map-data release endpoints must never be authorized by private Maps bearer tokens.

The initial provider error path records only an operation class rather than raw query text, coordinates, route waypoints, upstream bodies, or provider URLs. Collaboration storage failures also log the operation class rather than collection contents, member payloads, item coordinates, or notes. The browser source does not persist the bearer access token. Runtime observability must preserve this minimization contract.

No browser client secret is permitted. Production deployment must separately validate CSP, reverse proxy behavior, public map-data CORS/header policy, token exposure surfaces, Identity redirect registrations, Cloudflare/public-delivery configuration, and network boundaries.

## 11. Availability and resilience

The application should remain useful during partial outages. Cached map content, saved places, offline regions, and previously downloaded route/map resources should degrade independently from live provider services.

The UI must distinguish unavailable, stale, offline, delayed, approximate, unauthenticated, Identity-unconfigured, map-release-invalid, and unverified states instead of presenting all failures as empty results.

The provider API exposes configured/not-configured capability state so clients can distinguish a missing capability from a valid empty place/route result. No-provider operation remains an intentional supported development/degraded state.

The public release model supports roll-forward/rollback by changing only the mutable current manifest pointer to an already accepted immutable release. Release objects must never be mutated in place after publication. Runtime cache and rollback behavior still require deployed acceptance evidence.

The current browser disconnect action removes its local in-memory token only and must not be represented as provider-wide SSO logout. Provider logout/session management remains a later acceptance scope.

## 12. Release acceptance

A Stable release requires, as applicable:

- exact-revision CI;
- unit/integration tests;
- API authorization tests;
- multi-user isolation tests;
- database migration validation;
- map-data manifest and style contract validation;
- actual approved geographic dataset ingestion/build provenance;
- actual immutable release publication and rollback acceptance;
- actual GoreeCloud Identity Maps client registration and end-to-end login/session acceptance;
- accessibility and keyboard testing;
- reduced motion/transparency and forced-colors testing;
- responsive/form-factor testing;
- rendered map interaction acceptance;
- provider failure/degraded-mode tests;
- provider egress/SSRF acceptance;
- provider/dataset license/provenance and attribution acceptance;
- route/geocoder geographic-quality acceptance for supported coverage;
- public map-data CORS/cache/security acceptance;
- privacy and security acceptance;
- Everkeep export/restore acceptance for applicable user data;
- Wardveil, Privacy Shield, Glaze UI, Everkeep, and Mesh integration evidence where required;
- native and real-device acceptance for native navigation releases.

## 13. Initial milestones

1. Repository/product foundation and architecture contracts — source foundation established; review/merge remains gated.
2. Web map shell with Glaze UI 2.0 semantics and replaceable map-data/style provider — initial shell, versioned map-data release resolver/validator, and same-origin service surfaces established; approved dataset publication and rendered/live-provider acceptance pending.
3. Identity-backed multi-user API and PostGIS schema for saved places/collections — schema, collaboration APIs, optimistic revision handling, audit events, PKCE browser source, access-token API verification boundary, and automated multi-user/RLS PostGIS acceptance are established; actual Identity client registration, recipient directory/invitations, production database/Identity, load, backup, and deployment acceptance remain pending.
4. Place discovery/geocoding provider and Search interoperability — Nominatim-compatible source adapter/API and initial web results are established; approved provider deployment, quality acceptance, rich place data, and GoreeCloud Search interoperability remain pending.
5. Route planning provider and directions UI — Valhalla-compatible source adapter/API and initial web route rendering/maneuver summary are established; approved router deployment, route-quality acceptance, alternatives, advanced planning controls, and navigation remain pending.
6. Geographic release generation/publication — implement an approved reproducible dataset ingestion/tile generation pipeline, provenance/license evidence, immutable publication, deployed Worker/R2 or accepted equivalent delivery, CORS/cache verification, and rollback acceptance.
7. Offline region/package system.
8. Live navigation and Location integration.
9. Traffic/incidents/transit expansion.
10. Native mobile applications and device acceptance.
