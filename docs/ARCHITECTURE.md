# GoreeCloud Maps Architecture

## Status

This document describes the implemented foundation and intended service boundaries for GoreeCloud Maps. The application remains Development.

## Runtime shape

```text
Browser / future native client
            |
            v
      GoreeCloud Maps API
       /        |        \
      /         |         \
Identity     PostGIS     Provider adapters
  OIDC       user data    tiles / search /
                         routing / transit /
                         traffic / imagery
```

GoreeCloud Location is adjacent to this diagram rather than inside the Maps database. Maps may request approved current-location and sharing capabilities from Location, but personal location history, background tracking, Find My, geofences, and location-sharing authority remain owned by Location.

## Web application

`apps/web` is the first executable user interface. It uses MapLibre GL JS as a narrow rendering dependency and applies the GoreeCloud-controlled Glaze UI 2.0 experience around the renderer.

The default `public/map-style.json` has no remote sources. A live map only appears when `VITE_MAP_STYLE_URL` is set to an approved MapLibre-compatible style endpoint. This is intentional: cloning the repository must not silently disclose the user's IP address, viewport, or map activity to an unrelated public tile service.

The web search, directions, and collaboration controls are not yet connected to the protected API. Server-side provider and collaboration implementations therefore do not yet constitute end-to-end place-search, directions, or shared-map user experiences.

## API

`services/api` is the initial protected application API. It requires:

- a reachable PostgreSQL/PostGIS database;
- a non-owner runtime database role without `BYPASSRLS`;
- an OIDC issuer compatible with the GoreeCloud Identity integration contract;
- a Maps OIDC client ID.

The API does not accept a user ID from the client as an authorization authority. It validates the bearer token, resolves the OIDC `sub` claim to an application-owned Maps user record, and then passes that internal user ID into a transaction-local PostgreSQL setting used by row-level-security policies.

The service intentionally exits at startup when the configured database role can own/bypass the RLS-protected tables.

Implemented protected collaboration surfaces include collection list/create/update, member list/add/role-change/removal, and collection-item list/create/update/delete. Collection and item mutations use explicit expected revisions for optimistic conflict detection. Member addition accepts an existing Maps user ID; invitation delivery and governed GoreeCloud Identity/directory resolution are separate future capabilities rather than hidden account-discovery behavior.

Initial provider-backed application routes include authenticated forward geocoding, reverse geocoding, and route planning. A public capabilities route reports only whether geocoding/routing are configured; it does not expose provider origins or credentials. Provider-dependent protected routes remain unavailable when no approved endpoint is configured.

## Database and multi-user authorization

`db/migrations/0001_multi_user_foundation.sql` introduces first-party spatial and collaborative state. The initial model includes users, preferences, saved places, collections, collection members, collection items, and security-sensitive audit events.

`db/migrations/0002_audit_event_rls_hardening.sql` tightens audit-event insertion: an authenticated actor may write a collection-scoped audit event only when that actor has access to the referenced collection. This prevents a compromised ordinary runtime session from forging audit records against an unrelated collection merely by knowing its identifier.

Map collections have one immutable owner and optional `editor` or `viewer` members. Owners manage membership; editors may mutate collection content; viewers are read-only and may remove their own membership. Ownership transfer remains deliberately unsupported until a dedicated transaction and audit contract is designed.

Collection and collection-item revisions provide the current concurrency boundary. Stale collaborative updates are rejected rather than silently overwriting newer state. Implemented member and item mutations emit audit events without placing private collection contents, coordinates, notes, or provider payloads into application logs.

The database owner/migration role and API runtime role must be separate. The RLS helper used to evaluate collection membership is a narrowly scoped security-definer function. Production database acceptance must verify that the runtime role neither owns the protected tables nor holds `BYPASSRLS`.

CI applies every ordered SQL migration in `db/migrations` to an ephemeral PostGIS service and exercises the application store through owner/editor/viewer/stranger principals. The current test scope covers collection visibility, role changes, allowed/denied mutations, optimistic conflicts, private saved-place isolation, self-removal, immutable ownership, audit creation, forged-audit rejection, and the runtime-role privilege guard. This is source-level automated acceptance, not production database acceptance.

## Provider boundary

Provider integration is capability-based. Maps does not leak provider-specific response models through its first-party APIs or persistent user-data schema.

Implemented initial server-side provider adapters are:

- a Nominatim-compatible forward/reverse geocoding adapter using fixed `/search` and `/reverse` actions;
- a Valhalla-compatible route adapter using fixed `/route`, normalized for drive, walk, bicycle, and transit/multimodal modes.

These adapters are configured only by server-side environment values. No provider origin is enabled by default. Base URLs must be absolute HTTP(S) origins/paths without embedded credentials, query strings, or fragments. Requests have bounded timeouts/responses and validation, provider HTTP redirects are refused, and provider failures are logged by operation class without map searches, coordinates, route origins/destinations, or provider URLs.

This source boundary is not a complete production SSRF/egress control. Approved production deployment still requires network-level egress controls, DNS review, licensing/provenance acceptance, capacity/rate-limit controls, and Privacy Shield/Wardveil evidence. See `docs/PROVIDERS.md`.

Additional planned adapter boundaries include:

- vector/base-map tiles and styles;
- place/POI enrichment beyond geocoding;
- transit feeds/services beyond the current multimodal route seam;
- traffic and incidents;
- terrain/elevation;
- aerial/satellite imagery;
- street-level imagery;
- indoor maps.

MapLibre is the selected initial renderer, not the map-data provider. PostGIS is the first-party user/spatial state store, not the sole source of public geographic truth. Nominatim- and Valhalla-compatible interfaces are replaceable application adapters rather than permanent provider commitments.

## Identity boundary

GoreeCloud Identity establishes authentication. Maps owns authorization. Maps therefore stores its own user-resource ownership and collaborative roles rather than treating an Identity group claim as permanent authorization truth for every Maps object.

GoreeCloud Identity is still completing GoreeCloud-wide production SSO acceptance. Maps may develop against its OIDC-compatible integration contract without claiming that production SSO is already approved.

Future invitations must resolve recipients through an approved Identity/directory capability. Maps must not create an unrelated user-search authority or expose identity-subject identifiers as a substitute for a governed invitation experience.

## Privacy boundary

Precise user location, origins/destinations, private route plans, saved places, map searches, collection notes, and membership relationships can be sensitive. They must not become ordinary logs, analytics payloads, provider debug strings, cache keys exposed across users, or public CDN objects.

Public map resources and private user resources require separate cache and authorization policies. External provider calls must be documented with data-minimization and provenance requirements before production use.

Provider handlers intentionally avoid logging raw search queries, precise coordinates, route waypoints, provider response bodies, and provider endpoint URLs. Collaboration storage errors similarly log only an operation class, not member IDs, collection contents, precise item coordinates, or notes. These are source-level minimization controls, not substitutes for runtime observability/privacy acceptance.

## Resilience

The architecture separates public map/provider dependencies from durable user-owned data. A provider outage must not make previously saved collections disappear. Offline packages are planned as versioned resources with integrity, quota, update, rollback, and deletion behavior.

Provider configuration is optional. The service exposes configured/not-configured capability state so clients can distinguish a missing capability from an empty place/route result.

## Deployment boundary

No production deployment is defined by this foundation. Future deployment work must preserve the GoreeCloud platform boundaries for private/public networking, secrets, HTTPS, observability, Cloudflare Pages where applicable, backups, recovery, Privacy Shield, Wardveil Security, Everkeep, and GoreeCloud Mesh.
