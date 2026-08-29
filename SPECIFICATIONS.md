# GoreeCloud Maps Specifications

## 1. Status and authority

GoreeCloud Maps is in **Development**. This specification describes the required product and engineering direction; it does not assert that every listed capability is implemented.

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

## 4. Core domains

### Map presentation

The map engine must support smooth pan, zoom, rotate, pitch, user-location display, selectable layers, map style switching, scalable vector rendering, high-DPI output, and graceful degradation when advanced GPU effects are unavailable.

Target capabilities include globe mode, terrain, 3D buildings, indoor layers, accessibility overlays, transit overlays, traffic/incident overlays, and satellite/aerial imagery when an approved provider and licensing model exist.

### Places and discovery

The application must support forward and reverse geocoding, category search, nearby discovery, place cards, hours/contact metadata when licensed, accessibility metadata, saved places, favorites, user notes, collections, and place-data correction feedback.

Search results must identify provider/source provenance where required and must not silently combine incompatible licensing terms.

The initial source implementation includes authenticated forward and reverse geocoding API routes backed by a normalized Nominatim-compatible adapter. This is an **In progress** capability: no live geocoder endpoint, Search interoperability, web-result UI, geographic-quality acceptance, or production provider approval is claimed.

### Directions and navigation

The routing boundary must support driving, walking, cycling, transit, and multimodal itineraries where data coverage exists. Route responses should support alternatives, ETAs, distance/duration, maneuvers, geometry, accessibility constraints, avoidances, departure/arrival time, and route warnings.

The initial source implementation includes an authenticated route-planning API and a normalized Valhalla-compatible adapter for drive, walk, bicycle, and transit/multimodal costing. It currently normalizes route/leg distance and duration, encoded route shape, maneuvers, street names, and shape indexes. Live routing coverage, route-quality acceptance, alternatives, advanced constraints, directions UI, and navigation runtime remain pending.

Turn-by-turn navigation is a separate acceptance boundary requiring runtime GPS integration, rerouting, off-route detection, maneuver guidance, background/lock-screen behavior where supported, and verified real-device testing.

### Offline

Users must be able to select or download versioned regional map packages. Offline design must account for tiles, styles, glyphs/icons, place indexes, and routing graphs where supported. Packages require integrity metadata, versioning, quota management, update checks, rollback, and clean deletion.

### Collaboration

Users must be able to create collections and shared maps, add places, annotate entries, reorder content, invite/remove members, alter member roles, and revoke sharing. Conflict handling must be deterministic and auditable.

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

The first implemented server adapters are Nominatim-compatible forward/reverse geocoding and Valhalla-compatible routing. Their base URLs are optional server-side configuration and are blank by default. Configured URLs must be absolute HTTP(S) URLs without embedded credentials, query parameters, or fragments; Maps appends fixed provider action paths. The implementation bounds HTTP timeout, response size, search limits, waypoint counts, and coordinate/query validation.

This application-level URL validation is not a complete production SSRF defense. Production acceptance must add runtime egress restrictions and review DNS/redirect behavior, provider provenance/license terms, capacity/rate limits, secret handling where applicable, monitoring, Privacy Shield evidence, and Wardveil Security evidence.

Compatibility with an API shape does not authorize production use of an arbitrary public Nominatim, Valhalla, or other provider instance.

## 7. Platform integration

### GoreeCloud Identity

Identity is the principal/authentication authority. Maps must not create an unrelated permanent account system.

### GoreeCloud Location

Maps requests current/approved location capabilities and optional sharing overlays from Location. Permission, precision, tracking, and history truth remain with Location/Privacy Shield.

### GoreeCloud Search

Place and geographic discovery should interoperate with Search through documented APIs/events rather than duplicating universal search responsibilities.

### GoreeCloud Mesh

Cross-application coordination, capability discovery, and governed events should use Mesh when the applicable contracts are available.

### Privacy Shield

Privacy state must be evidence-backed. Requirements include data minimization, explicit sharing, revocation, clear precision controls, private-history options, provider disclosure where needed, and no hidden tracking.

### Wardveil Security

Required controls include authenticated protected APIs, authorization at every resource boundary, rate limiting/abuse controls, input validation, SSRF-resistant provider access, secure secret handling, dependency review, and security-event evidence.

### Everkeep

User-owned saved places, collections, annotations, preferences, and other durable personal map data must have export, backup, restore, portability, and deletion semantics appropriate to the data class.

## 8. Glaze UI 2.0 mapping

The map itself is the primary Canvas. Persistent or contextual controls use Soft Glaze/Glaze; expanded menus and sheets use Deep Glaze; active navigation and other ongoing processes may use Live Glaze.

Mobile composition uses a high-information map/viewing zone with primary reachable actions near the lower action zone. Search expands from its invoking control. Place cards and route details use connected sheet transformations rather than unrelated modal jumps.

Desktop uses a denser sidebar/inspector model, keyboard shortcuts, pointer states, resizable panels, and contextual menus. Tablet/foldable layouts use split panes and hinge-aware placement. Reduced transparency/motion and effects-free modes must preserve complete usability.

## 9. API principles

All protected APIs are versioned under `/api/v1/` initially. APIs use stable resource identifiers, explicit pagination, structured errors, idempotency where retries are expected, optimistic concurrency or version fields for collaborative mutations, and server-side authorization.

Provider credentials and internal upstream endpoints must never be exposed to untrusted clients.

Current source routes include:

- `GET /api/v1/me`;
- `GET /api/v1/collections`;
- `POST /api/v1/collections`;
- `GET /api/v1/search`;
- `GET /api/v1/reverse`;
- `POST /api/v1/routes`;
- public `GET /api/v1/capabilities`, limited to configured/not-configured provider state.

The provider-dependent search/reverse/route routes require authenticated Maps users and return explicit unavailable state when the capability is not configured.

## 10. Security and privacy constraints

Precise location is sensitive. Logs must avoid recording precise coordinates, route origins/destinations, search queries, tokens, or share secrets unless a specific operational requirement exists and the retention/privacy contract permits it.

Public map tiles and user-private resources must use separate caching rules. Private responses must not become publicly cacheable through CDN configuration.

The initial provider error path records only an operation class rather than raw query text, coordinates, route waypoints, upstream bodies, or provider URLs. Runtime observability must preserve this minimization contract.

## 11. Availability and resilience

The application should remain useful during partial outages. Cached map content, saved places, offline regions, and previously downloaded route/map resources should degrade independently from live provider services.

The UI must distinguish unavailable, stale, offline, delayed, approximate, and unverified states instead of presenting all failures as empty results.

The provider API exposes configured/not-configured capability state so clients can distinguish a disabled provider from a valid empty result. No-provider operation remains an intentional supported development/degraded state.

## 12. Release acceptance

A Stable release requires, as applicable:

- exact-revision CI;
- unit/integration tests;
- API authorization tests;
- multi-user isolation tests;
- database migration validation;
- accessibility and keyboard testing;
- reduced motion/transparency and forced-colors testing;
- responsive/form-factor testing;
- rendered map interaction acceptance;
- provider failure/degraded-mode tests;
- provider egress/SSRF acceptance;
- provider license/provenance and attribution acceptance;
- route/geocoder geographic-quality acceptance for supported coverage;
- privacy and security acceptance;
- Everkeep export/restore acceptance;
- Wardveil, Privacy Shield, Glaze UI, and Everkeep integration evidence;
- native and real-device acceptance for native navigation releases.

## 13. Initial milestones

1. Repository/product foundation and architecture contracts — source foundation established; review/merge remains gated.
2. Web map shell with Glaze UI 2.0 semantics and replaceable map-style provider — initial shell established; rendered/provider acceptance pending.
3. Identity-backed multi-user API and PostGIS schema for saved places/collections — initial schema and collection primitives established; two-user runtime acceptance pending.
4. Place discovery/geocoding provider and Search interoperability — Nominatim-compatible source adapter/API established; approved provider deployment, web integration, quality acceptance, and GoreeCloud Search interoperability pending.
5. Route planning provider and directions UI — Valhalla-compatible source adapter/API established; approved router deployment, directions UI, route-quality acceptance, and advanced planning controls pending.
6. Offline region/package system.
7. Live navigation and Location integration.
8. Traffic/incidents/transit expansion.
9. Native mobile applications and device acceptance.
