# GoreeCloud Maps

GoreeCloud Maps is the first-party GoreeCloud mapping, place-discovery, directions, navigation, and collaborative map application. It is being built as an original GoreeCloud product with a privacy-first, self-hostable architecture and multi-user support.

**Status:** Development

## Product boundary

GoreeCloud Maps owns the map experience: map rendering, place search and discovery, place cards, directions and route planning, navigation presentation, saved places, collections, collaborative maps, offline map packages, map layers, user-contributed map feedback, and provider orchestration.

GoreeCloud Location remains the authority for personal location history, live device location, tracking, Find My capabilities, geofences, and location-sharing permissions. Maps consumes approved Location capabilities rather than duplicating that sensitive data domain.

## Design and platform requirements

Maps targets the current Stable Glaze UI contract, **Glaze UI 2.0.0**, including Glaze Material semantics, adaptive form-factor transformations, connected transitions, accessible effects fallbacks, reduced-motion behavior, and the Canvas / Surface / Soft Glaze / Glaze / Deep Glaze / Live Glaze hierarchy.

Stable release qualification also requires current accepted integration with:

- Privacy Shield for privacy controls, consent, data minimization, private/approximate location behavior, and truthful privacy state.
- Wardveil Security for security state, abuse resistance, protected service boundaries, and evidence-backed protection experiences.
- Everkeep for backup, recovery, preservation, portability, and continuity of user-owned map data.
- GoreeCloud Mesh for governed interoperability with first-party GoreeCloud services.

Maps must not claim production readiness from design-system alignment alone. Application-specific runtime, accessibility, security, privacy, resilience, platform, and real-device acceptance remain required.

## Architecture direction

The product is designed around replaceable capabilities rather than proprietary platform lock-in:

- **Web renderer:** MapLibre GL JS.
- **Native renderer:** MapLibre Native where a narrow rendering dependency is appropriate.
- **Spatial database:** PostgreSQL + PostGIS.
- **Routing:** provider interface with a Valhalla-compatible server adapter implemented as the initial self-hostable baseline.
- **Geocoding/search:** GoreeCloud-owned provider interface with a Nominatim-compatible forward/reverse adapter implemented, plus planned GoreeCloud Search interoperability.
- **Map data:** open, license-compliant sources such as OpenStreetMap and other approved datasets, processed into versioned public releases and served through GoreeCloud-controlled infrastructure.
- **Tiles/packages:** schema-v1 public releases use vector/MVT ZXY tiles with immutable style/tile/glyph/sprite resources; future offline/archive/imagery schemas remain separate work.
- **Identity:** GoreeCloud Identity is the authentication and principal authority; the web source implements optional Authorization Code + PKCE as a public client and the API expects provider-recognized bearer access tokens.
- **Location:** GoreeCloud Location supplies approved device/user location capabilities.

Provider adapters are mandatory so routing, geocoding, traffic, transit, imagery, and tile delivery can evolve without rewriting the application. The implemented geocoding and routing adapters are configuration-driven and have no provider endpoint enabled by default; they are source foundations, not evidence of live provider coverage or production acceptance.

The web client has source-level same-origin `/api/v1` integration for provider capability status, authenticated place search, provider-backed directions, shared collection listing/creation, and collection-item browsing. Search results can move the map to a returned place, and route results can render returned Valhalla geometry and maneuver summaries. These experiences remain unavailable when the required Maps API, Identity registration, geocoder, or router is not configured; no public provider is silently substituted.

## Public geographic data plane

The repository defines a separate public map-data delivery boundary under `mapdata/` and `services/mapdata-edge`.

- versioned release manifests record release ID, coverage, zoom range, attribution, source dataset versions, licensing/provenance, style path, and tile template;
- immutable objects live beneath `releases/<release-id>/` and must never be replaced in place after publication;
- `manifests/current.json` is the only mutable release pointer and is intentionally short-cached;
- the Cloudflare Worker/R2 source gateway is read-only and allowlists only current/release manifests, styles, vector tiles, sprites, and glyphs;
- the Worker exposes no bucket listing, write surface, private Maps API, Identity state, saved data, searches, routes, collections, or personal location data;
- CI validates both synthetic release-manifest and MapLibre style fixtures and performs a Wrangler dry bundle of the edge service;
- schema v1 permits only vector/MVT release-local source resources and rejects imported styles, TileJSON indirection, external font-face resources, and off-release tile/glyph/sprite resources;
- the browser can use `VITE_MAP_DATA_MANIFEST_URL` to resolve the current release without rebuilding the application shell. It starts on the local empty style, validates the manifest and immutable style, normalizes allowed release resources, injects validated manifest attribution, then supplies the resulting in-memory style to MapLibre;
- a rejected/unavailable release keeps the privacy-safe local fallback instead of silently falling through to an unrelated public style.

`VITE_MAP_STYLE_URL` remains a manual/legacy approved-style seam when no manifest endpoint is configured.

Exact-head repository CI run #78 (run ID `33241003328`) passed on `1655bb5a1dde62964ba2503e37764e4eaa7134af`, including the web production build, manifest and style validation, Worker type/dry-run bundle, and the existing API/PostGIS acceptance suite. This validates the source contracts, not a live geographic-data deployment.

The repository does not contain an approved basemap dataset and does not prove that the intended R2 buckets, Worker, Pages project, routes, or custom domains have been provisioned. Raster/satellite imagery, terrain, PMTiles/other archive formats, and additional source classes require future reviewed schemas. Production publication still requires dataset licensing/provenance, attribution, rendering quality, cache/rollback, CORS/origin policy, Privacy Shield, Wardveil Security, and Cloudflare runtime acceptance.

See [mapdata/README.md](mapdata/README.md), [deployment/cloudflare/README.md](deployment/cloudflare/README.md), [docs/PROVIDERS.md](docs/PROVIDERS.md), and [docs/IDENTITY.md](docs/IDENTITY.md).

## Multi-user model

Maps is designed for individual users, households, teams, and shared collections. User-scoped resources must remain isolated by default. Sharing is explicit and capability-scoped.

Core collaborative resources include saved places, collections/lists, shared maps, route plans, annotations, contributed edits, and optional ETA/location overlays delegated to GoreeCloud Location. Roles are modeled explicitly rather than inferred from possession of a URL.

The current source implements owner/editor/viewer authorization, PostGIS row-level security, collection updates with optimistic revision checks, member list/add/role-change/removal APIs, collection-item list/create/update/delete APIs, and collaboration audit events. Automated CI exercises these paths against a live PostGIS service with owner/editor/viewer/stranger principals.

The web client can list/create authorized collections and browse their items after authentication. Human-friendly invitation discovery is deliberately blocked until GoreeCloud Identity defines an approved cross-application recipient-directory contract. Maps does not use the Identity provider's administrative user directory as a substitute consumer directory. Share links, ownership transfer, full collaboration management UI, and production Identity/database acceptance remain pending.

## Identity source boundary

The optional browser Identity client is configured by issuer/client registration and uses Authorization Code + PKCE (`S256`) without a client secret. PKCE verifier/state are transient, and the bearer access token is kept in memory rather than persistent browser storage. The API validates the configured issuer/audience/expiry and confirms the same subject through the provider UserInfo endpoint before resolving the subject to an internal Maps user.

This is a source integration contract only. GoreeCloud Identity is not yet approved for GoreeCloud-wide production SSO, and Maps does not claim production login/logout/session, disabled-account, outage, recovery, or rollback acceptance.

## Experience direction

The experience is heavily inspired by the useful interaction patterns users expect from leading map products while remaining visually and technically original to GoreeCloud. The roadmap includes:

- fast pan/zoom/rotate/tilt map interaction;
- 2D, globe, terrain, and 3D building views;
- rich place cards and category discovery;
- driving, walking, cycling, transit, and multimodal route planning;
- turn-by-turn navigation with lane, maneuver, speed, incident, and arrival context where supported by verified data;
- offline regions and offline routing where supported;
- saved places, favorites, guides, collections, and collaborative maps;
- privacy-preserving search/history controls;
- live traffic, closures, hazards, and incident reports when reliable data pipelines are available;
- accessible map alternatives, keyboard/pointer/touch navigation, high-contrast modes, reduced transparency, and reduced motion;
- mobile, tablet, foldable, desktop, TV, wearable, and future spatial adaptations where the applicable Glaze UI and platform acceptance requirements are satisfied.

## Repository documents

- [SPECIFICATIONS.md](SPECIFICATIONS.md) — product, architecture, data, API, security, and acceptance requirements.
- [FEATURES.md](FEATURES.md) — feature inventory with implementation status.
- [BENEFITS.md](BENEFITS.md) — user, administrator, privacy, and platform benefits.
- [COMPETITIVE-OBJECTIVES.md](COMPETITIVE-OBJECTIVES.md) — competitive capability targets and differentiation.
- [mapdata/README.md](mapdata/README.md) — public geographic-data release and publication contract.
- [deployment/cloudflare/README.md](deployment/cloudflare/README.md) — intended Pages + Worker/R2 deployment boundary and evidence requirements.
- [docs/PROVIDERS.md](docs/PROVIDERS.md) — provider API contracts, configuration, privacy/security controls, and acceptance gates.
- [docs/IDENTITY.md](docs/IDENTITY.md) — browser OIDC/PKCE, API access-token validation, and invitation-directory boundary.
- `docs/` — architecture, provider, Glaze UI, privacy/security, and integration records as implementation progresses.

## Development rules

GoreeCloud Maps is an original GoreeCloud application. Complete third-party map applications are not used as the product foundation. Narrow foundational dependencies may be adopted when rebuilding them would increase security, protocol, geospatial, rendering, codec, or platform risk.

Map data, imagery, reviews, proprietary place content, and navigation data from Google Maps or Apple Maps must not be copied, scraped, or represented as GoreeCloud-owned data. Inspiration is limited to product capability and interaction objectives.

## License

Unless otherwise noted, repository source is licensed under the GNU Affero General Public License v3.0. Third-party datasets and dependencies retain their own licenses and attribution requirements.
