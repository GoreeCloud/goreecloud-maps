# GoreeCloud Maps Architecture

## Status

This document describes the implemented source foundation and intended service boundaries for GoreeCloud Maps. The application remains Development.

## Runtime shape

```text
                       GoreeCloud Identity
Browser / future      authentication only
native client          ^
      |                | OIDC Authorization Code + PKCE
      |
      +---- same-origin /api/v1 ----> GoreeCloud Maps API
      |                                  /      |       \
      |                                 /       |        \
      |                            Identity  PostGIS   Provider adapters
      |                            validation user data search / routing /
      |                                                 transit / traffic
      |
      +---- public map-data manifest/style/tiles ----> Map-data edge
                                                       read-only release delivery
                                                       (Worker/R2 source contract)
```

The public map-data path and the authenticated Maps API path are deliberately separate. Public release requests are credentialless; private bearer tokens must not be sent to the map-data origin.

GoreeCloud Location is adjacent to this diagram rather than inside the Maps database. Maps may request approved current-location and sharing capabilities from Location, but personal location history, background tracking, Find My, geofences, and location-sharing authority remain owned by Location.

## Web application

`apps/web` is the first executable user interface. It uses MapLibre GL JS as a narrow rendering dependency and applies the GoreeCloud-controlled Glaze UI 2.0 experience around the renderer.

The default `public/map-style.json` has no remote sources. This is intentional: cloning or launching an unconfigured build must not silently disclose the user's IP address, viewport, or map activity to an unrelated public tile service.

### Geographic release loading

`VITE_MAP_DATA_MANIFEST_URL` is the preferred public map-data seam. When configured, the browser remains on the local empty style until it has independently validated the release.

The current release-loading sequence is:

1. request the configured manifest without credentials, with redirects refused;
2. bound and parse the manifest;
3. verify schema version, release/path consistency, bounds/zoom, attribution, source/version/license/provenance records, and `publicGeographicDataOnly`;
4. request the immutable release `style.json` without credentials;
5. enforce map-data schema v1 before MapLibre sees the style;
6. normalize permitted tiles/glyphs/sprites to the configured immutable release origin;
7. inject validated manifest attribution into the first vector source;
8. provide the normalized in-memory style to MapLibre;
9. keep the local fallback and expose a degraded state if any validation/load step fails.

Schema v1 is deliberately vector/MVT-only. It rejects TileJSON indirection, imported styles, external font-face resources, non-vector sources, off-release tile paths, and off-release glyph/sprite resources. Raster/aerial imagery, terrain, PMTiles/other archives, and additional style/source classes require future explicitly versioned contracts.

`VITE_MAP_STYLE_URL` remains a manual/legacy approved style seam only when `VITE_MAP_DATA_MANIFEST_URL` is empty. A configured manifest failure does not fall through to that raw style seam.

The source therefore supports release pointer roll-forward/rollback without rebuilding the application shell, while immutable release objects remain independently cacheable. This does not prove any map-data service or dataset is deployed.

### Maps API client

The web client contains a same-origin Maps API client. `VITE_MAPS_API_BASE_PATH` defaults to `/api/v1` and must remain an absolute path beginning with `/`; an arbitrary external API origin is rejected. The intended deployment boundary is therefore Browser → controlled same-origin reverse proxy → Maps API rather than Browser → configurable third-party private API.

Current source-connected web behavior includes:

- public provider capability-state loading;
- optional GoreeCloud Identity sign-in through Authorization Code + PKCE;
- authenticated place search and selectable map results;
- authenticated directions that geocode endpoints, request a route, render returned polyline6 geometry, and show maneuver summaries;
- authenticated collection listing/creation and collection-item browsing;
- explicit unavailable/degraded states when the API, Identity registration, geocoder, router, or configured map-data release is absent/invalid;
- truthful placeholders for saved-place workflows and GoreeCloud Location current-position integration.

These are source-level application integrations. They do not imply a live provider, production reverse proxy, production GoreeCloud Identity client, map-data release, or deployment acceptance.

## Browser Identity client

The optional browser Identity integration is a public OIDC client:

- issuer/client registration is configuration-driven and blank by default;
- Authorization Code + PKCE (`S256`);
- no browser client secret;
- OIDC discovery only when Identity is configured and needed;
- HTTPS endpoint enforcement outside localhost development;
- same-origin Maps redirect URI enforcement;
- explicit state validation;
- PKCE verifier/state stored transiently in `sessionStorage`;
- bearer access token held in application memory, not persistent browser storage.

The current Maps disconnect action clears the local in-memory token and transient PKCE state. It does not claim to terminate the provider-wide SSO session. Provider logout/session handling is a later application-integration acceptance scope.

See `docs/IDENTITY.md` for the detailed contract.

## API

`services/api` is the protected application API. It requires:

- a reachable PostgreSQL/PostGIS database;
- a non-owner runtime database role without `BYPASSRLS`;
- an OIDC issuer compatible with the GoreeCloud Identity integration contract;
- a Maps OIDC client ID.

The API does not accept a user ID from the client as an authorization authority. Protected routes expect a bearer **access token**. The current verifier validates the signed JWT against the configured issuer and Maps client audience/expiry semantics, requires a non-empty subject, then calls the provider UserInfo endpoint with that bearer and requires the same subject. Only then is the subject resolved to an application-owned Maps user record and passed into the transaction-local PostgreSQL setting used by row-level-security policies.

The UserInfo confirmation is the current source mechanism for requiring a provider-recognized access token rather than treating an ID token as an API bearer. It also makes protected API calls dependent on Identity availability/latency; production acceptance must decide whether this exact runtime validation topology remains appropriate and validate failure behavior.

The service intentionally exits at startup when the configured database role can own/bypass the RLS-protected tables.

Implemented protected collaboration surfaces include collection list/create/update, member list/add/role-change/removal, and collection-item list/create/update/delete. Collection and item mutations use explicit expected revisions for optimistic conflict detection. Member addition accepts an existing Maps user ID; invitation delivery and governed GoreeCloud Identity/directory resolution are separate capabilities rather than hidden account-discovery behavior.

Provider-backed application routes include authenticated forward geocoding, reverse geocoding, and route planning. A public capabilities route reports only whether geocoding/routing are configured; it does not expose provider origins or credentials. Provider-dependent protected routes remain unavailable when no approved endpoint is configured.

## Database and multi-user authorization

`db/migrations/0001_multi_user_foundation.sql` introduces first-party spatial and collaborative state. The initial model includes users, preferences, saved places, collections, collection members, collection items, and security-sensitive audit events.

`db/migrations/0002_audit_event_rls_hardening.sql` tightens audit-event insertion: an authenticated actor may write a collection-scoped audit event only when that actor has access to the referenced collection. This prevents a compromised ordinary runtime session from forging audit records against an unrelated collection merely by knowing its identifier.

Map collections have one immutable owner and optional `editor` or `viewer` members. Owners manage membership; editors may mutate collection content; viewers are read-only and may remove their own membership. Ownership transfer remains deliberately unsupported until a dedicated transaction and audit contract is designed.

Collection and collection-item revisions provide the current concurrency boundary. Stale collaborative updates are rejected rather than silently overwriting newer state. Implemented member and item mutations emit audit events without placing private collection contents, coordinates, notes, or provider payloads into application logs.

The database owner/migration role and API runtime role must be separate. The RLS helper used to evaluate collection membership is a narrowly scoped security-definer function. Production database acceptance must verify that the runtime role neither owns the protected tables nor holds `BYPASSRLS`.

CI applies every ordered SQL migration in `db/migrations` to an ephemeral PostGIS service and exercises the application store through owner/editor/viewer/stranger principals. The current test scope covers collection visibility, role changes, allowed/denied mutations, optimistic conflicts, private saved-place isolation, self-removal, immutable ownership, audit creation, forged-audit rejection, and the runtime-role privilege guard. This is source-level automated acceptance, not production database acceptance.

## Public map-data plane

`mapdata/` defines the versioned publication contract. `services/mapdata-edge` defines the initial read-only public delivery service around an R2 binding.

The release layout separates one mutable pointer from immutable releases:

```text
manifests/current.json
releases/<release-id>/manifest.json
releases/<release-id>/style.json
releases/<release-id>/tiles/<z>/<x>/<y>.pbf
releases/<release-id>/sprites/...
releases/<release-id>/glyphs/...
```

The edge service allowlists those object shapes and implements GET/HEAD/OPTIONS only. It does not expose R2 listing, arbitrary keys, writes, private APIs, Identity/session data, saved places, collections, searches, routes, or location state.

Mutable current-manifest caching and immutable release-object caching are separate. Rollback changes the current pointer to a previously accepted immutable release; it never overwrites the old release.

CI validates the synthetic manifest and style fixtures and dry-builds the Worker. The fixtures contain no approved geographic dataset.

Production still requires an actual reproducible map-data generation pipeline, dataset license/provenance acceptance, attribution verification, quality/rendering tests, production CORS/origin policy, Worker/R2 or equivalent provisioning, runtime cache/header tests, security/privacy evidence, capacity/abuse controls, publication evidence, and rollback acceptance.

## Provider boundary

Provider integration is capability-based. Maps does not leak provider-specific response models through its first-party APIs or persistent user-data schema.

Implemented initial server-side provider adapters are:

- a Nominatim-compatible forward/reverse geocoding adapter using fixed `/search` and `/reverse` actions;
- a Valhalla-compatible route adapter using fixed `/route`, normalized for drive, walk, bicycle, and transit/multimodal modes.

These adapters are configured only by server-side environment values. No provider origin is enabled by default. Base URLs must be absolute HTTP(S) origins/paths without embedded credentials, query strings, or fragments. Requests have bounded timeouts/responses and validation, provider HTTP redirects are refused, and provider failures are logged by operation class without map searches, coordinates, route origins/destinations, or provider URLs.

The map-data plane is also provider-replaceable, but its client contract is release-oriented instead of allowing arbitrary tile origins through style JSON. MapLibre is the renderer, not the geographic-data authority. The current v1 release can be generated from any dataset/pipeline that passes the GoreeCloud provenance, licensing, schema, quality, and operational gates.

This source boundary is not a complete production SSRF/egress/security control. Approved production deployment still requires network-level egress controls for server provider calls, DNS review, public-delivery origin/CORS controls, licensing/provenance acceptance, capacity/rate-limit controls, and Privacy Shield/Wardveil evidence. See `docs/PROVIDERS.md`.

Additional planned adapter/source boundaries include:

- place/POI enrichment beyond geocoding;
- transit feeds/services beyond the current multimodal route seam;
- traffic and incidents;
- terrain/elevation;
- aerial/satellite imagery;
- street-level imagery;
- indoor maps;
- future archive/offline packaging.

PostGIS is the first-party user/spatial state store, not the sole source of public geographic truth. Nominatim- and Valhalla-compatible interfaces and the map-data release format are application boundaries, not permanent provider commitments.

## Identity and invitation boundary

GoreeCloud Identity establishes authentication. Maps owns authorization. Maps therefore stores its own user-resource ownership and collaborative roles rather than treating an Identity group claim as permanent authorization truth for every Maps object.

GoreeCloud Identity is still completing GoreeCloud-wide production SSO acceptance. Maps may develop against its OIDC-compatible integration contract without claiming that production SSO is approved.

Human-friendly invitations require an approved cross-application recipient-discovery/invitation contract. The current evidence does not establish one. Maps therefore must not query the Identity provider's administrative user directory as an accidental consumer directory, expose provider administration identifiers to ordinary clients, or create an unrelated Maps user-search authority. The existing API's internal Maps user-ID membership primitive remains useful for controlled tests/integration only.

A future invitation boundary must define privacy/enumeration controls, permitted identity attributes, invitation delivery/expiry/acceptance/revocation, disabled-account behavior, abuse controls, auditing, and any Mesh/Notify coordination.

## Privacy boundary

Precise user location, origins/destinations, private route plans, saved places, map searches, collection notes, membership relationships, OAuth authorization codes, PKCE verifiers, and bearer tokens can be sensitive. They must not become ordinary logs, analytics payloads, provider debug strings, cache keys exposed across users, or public CDN objects.

Public map resources and private user resources require separate cache and authorization policies. The public map-data bundle is restricted to public geographic resources and release metadata; it must never contain saved places, collection state, private routes/searches, Identity state, or personal location data.

Provider handlers intentionally avoid logging raw search queries, precise coordinates, route waypoints, provider response bodies, and provider endpoint URLs. Collaboration storage errors similarly log only an operation class, not member IDs, collection contents, precise item coordinates, or notes. The web source does not persist the bearer access token. Public map-data fetches are credentialless. These are source-level minimization controls, not substitutes for runtime observability/privacy acceptance.

## Resilience

The architecture separates public map/provider dependencies from durable user-owned data. A provider or map-data outage must not make previously saved collections disappear.

Provider configuration is optional. The service exposes configured/not-configured capability state so clients can distinguish a missing capability from an empty place/route result. Identity registration is also optional in the web source; when absent, protected features present a sign-in/configuration state rather than creating fallback authentication.

The map-data release model adds a fail-closed local renderer fallback plus immutable release rollback semantics. A broken current manifest/style must not erase user data or cause an unapproved third-party basemap fallback.

Offline packages remain planned as versioned resources with integrity, quota, update, rollback, and deletion behavior. The public release manifest supplies some future freshness/provenance metadata but is not itself the offline package implementation.

The current protected API access-token validation includes a UserInfo call, so Identity outages can make protected operations unavailable even when PostGIS is healthy. That behavior must be intentionally tested and either accepted or redesigned before production.

## Deployment boundary

No production deployment is defined by this foundation. Future deployment work must preserve the GoreeCloud platform boundaries for same-origin private API proxying, separate credentialless public map-data delivery, private/public networking, secrets, HTTPS, CSP, headers, CORS, observability, Cloudflare Pages where applicable, Worker/R2 or equivalent geographic release delivery, backups, recovery, Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, and GoreeCloud Identity.
