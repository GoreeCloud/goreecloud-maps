# GoreeCloud Maps

GoreeCloud Maps is the first-party GoreeCloud mapping, place-discovery, directions, saved-place, offline-map, and collaborative-map application and service.

**Lifecycle:** Forge  
**Deployment state:** Development  
**Platform Contract:** 2.0  
**Current approved Glaze UI target:** 1.6.0

## Current repository state

Authoritative `main` remains the accepted governance/licensing/branding baseline until this candidate passes review and protected promotion.

This stabilization branch layers the executable Maps foundation and later Saved Places/search-save work onto current `main` instead of carrying forward the previously diverged stacked PR chain. The candidate includes:

- TypeScript/Vite + MapLibre GL JS web application shell;
- Go Maps API;
- PostgreSQL/PostGIS multi-user state and row-level-security tests;
- GoreeCloud Identity-compatible Authorization Code + PKCE browser boundary;
- owner/editor/viewer resource authorization;
- provider adapters for Nominatim-compatible geocoding and Valhalla-compatible routing;
- owner-scoped Saved Places and shared collection primitives;
- versioned public map-data release validation and a read-only edge-delivery source contract;
- current Platform Contract 2.0 declaration and exact-head CI.

These are candidate capabilities until review and promotion complete. No live provider, production Identity client, production database, deployment, release, or Anchor status is implied.

## Product boundary

Maps owns map exploration, geographic/place search, directions and route presentation, Saved Places, collections, collaborative maps, offline map experiences, map layers, and map-provider orchestration.

GoreeCloud Location remains authoritative for current device/user position, personal location history, background tracking, Find My, geofences, and location-sharing permission state. Maps consumes approved Location capabilities rather than creating a competing sensitive-location authority.

## Architecture

The candidate architecture is provider-replaceable and first-party-state oriented:

- **Web renderer:** MapLibre GL JS.
- **Application API:** Go.
- **Spatial state:** PostgreSQL + PostGIS.
- **Authentication:** GoreeCloud Identity-compatible Authorization Code + PKCE; Maps retains resource authorization.
- **Geocoding:** replaceable Nominatim-compatible adapter.
- **Routing:** replaceable Valhalla-compatible adapter.
- **Public geographic data:** immutable release contracts separated from private Maps/user state.
- **Map-data edge:** read-only delivery boundary for approved public releases.

A fresh development build must not silently fall back to an unrelated public map/geocoder/router provider.

## Platform requirements

Platform Contract 2.0 evaluates exactly nine Integral Platform Systems: GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, GoreeCloud Mesh, GoreeCloud Identity, GoreeCloud Policy, and GoreeCloud Observability.

The web candidate now implements a repository-local GLAZE UI V1.6 / 1.6.0 adapter with protected semantic surfaces, accessibility-profile fallbacks, caller-owned state truth, input-mode adaptation, reduced-motion handling, and bounded performance degradation. Maps-specific rendered, accessibility, adaptive, performance, browser/device/GPU, rollback, and Human Visual Excellence acceptance remains blocked.

GoreeCloud Sync remains separately governed and is not a tenth Integral Platform System.

## Privacy and security boundaries

Private searches, routes, saved resources, collaboration state, Identity state, and precise personal location must not enter the public map-data plane. Provider credentials stay server-side. Authorization is enforced both at the application boundary and, for supported data, with PostgreSQL row-level security.

No source label, branding asset, green CI run, or provider-compatible interface substitutes for Privacy Shield, Wardveil Security, Policy, Observability, continuity, runtime, deployment, or production acceptance.

## Canonical identity

Branding authority is `GoreeCloud/goreecloud-branding-assets`. The approved Maps product identity is `products/maps/app-icon.svg`, canonical Git blob `07b6e52e04c95e1ec9f703a9d323cf799481351c`.

The folded-map/route identity is intentionally distinct from GoreeCloud Location's positioning/pin identity. See [BRANDING.md](BRANDING.md).

## Repository governance records

- [SPECIFICATIONS.md](SPECIFICATIONS.md) — durable product, architecture, privacy/security, integration, and acceptance contract.
- [FEATURES.md](FEATURES.md) — evidence-scoped feature inventory.
- [PROJECT-SPECIFICATIONS.md](PROJECT-SPECIFICATIONS.md) — canonical project requirements, architecture, boundaries, and acceptance gates.
- [PROJECT-RECORD.md](PROJECT-RECORD.md) — significant project history, governance transitions, and evidence.
- [IMPLEMENTED-FEATURES.md](IMPLEMENTED-FEATURES.md) — evidence-scoped implemented capability inventory.
- [PLANNED-FEATURES.md](PLANNED-FEATURES.md) — open, partial, blocked, and planned capability obligations.
- [CHANGELOGS.md](CHANGELOGS.md) — repository-oriented change history.
- [BENEFITS.md](BENEFITS.md) — intended user, privacy, operational, accessibility, and resilience benefits.
- [COMPETITIVE-OBJECTIVES.md](COMPETITIVE-OBJECTIVES.md) — capability and differentiation objectives.
- [BRANDING.md](BRANDING.md) — canonical visual-identity consumer contract.
- [docs/GLAZE_UI_ADOPTION.md](docs/GLAZE_UI_ADOPTION.md) — current implemented/required Glaze boundary.

`PROJECT-SPECIFICATIONS.md` and `PROJECT-RECORD.md` are the repository-local project authority once this migration is accepted on `main`. The former Drive project specification is frozen migration input only until default-branch readback permits its permanent removal. Feature state and changelogs are repository-native; Google Drive must not be maintained as a parallel roadmap or changelog authority.

## Qualification boundary

Maps is not Anchor or production-ready until all applicable exact-head source/build, provider/data, accessibility, database/RLS, Identity, Location, Privacy Shield, Wardveil, Everkeep, Mesh, Policy, Observability, continuity, deployment, rollback, release, and production gates are satisfied.

## License

Unless otherwise noted, GoreeCloud Maps repository source is licensed under GNU AGPL version 3. Third-party dependencies, datasets, providers, imagery, transit feeds, and other incorporated material retain their own licenses and attribution requirements.
