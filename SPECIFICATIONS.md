# GoreeCloud Maps Specifications

## 1. Status and authority

GoreeCloud Maps is in **Development**. The authoritative product/project record is `GoreeCloud/Projects/Project Specification — Maps`; this repository record is the code-adjacent implementation/governance companion and must remain consistent with that authority.

Accepted `main` currently contains governance, licensing, and approved branding, not the executable Maps application foundation. Draft PR #1 (`agent/maps-foundation`) contains the validated executable candidate at `a8aafba65b89dbcd76661368a605762871abbb35`, with CI run `33251905892` successful. PR #1 is still unmerged because its independent-review gate remains outstanding.

No capability that exists only in PR #1 may be described as accepted `main`, production, released, or Stable behavior.

## 2. Product purpose

Maps is the GoreeCloud-owned mapping, place-discovery, directions, navigation, saved-place, offline-map, and collaborative-map application/service.

Maps may use Google Maps and Apple Maps as capability/interaction references only. It must not copy or scrape proprietary map data, place data, reviews, imagery, routing/navigation data, interface assets, wording, or protected product content.

## 3. Product and authority boundaries

Maps owns:

- interactive map browsing and exploration;
- place/address/category discovery and map-specific search experiences;
- place details and map-context actions;
- directions and route planning;
- turn-by-turn navigation presentation when applicable runtime capabilities are approved;
- saved places, favorites, lists, guides, and collaborative collections;
- offline map-region experiences;
- map layers, terrain, 3D/globe, imagery, traffic, incidents, transit overlays, and indoor mapping when supported by approved data/providers;
- map-specific preferences and history subject to Privacy Shield and retention requirements.

GoreeCloud Location remains authoritative for device/current-position services, background tracking, personal location history, Find My capabilities, geofences, location-sharing permission state, and sensitive tracking/sharing policy. Maps must consume approved Location capabilities instead of duplicating that authority.

GoreeCloud Identity authenticates principals, while Maps remains responsible for Maps resource authorization. GoreeCloud Search remains the wider search authority; Maps may expose geographic/place search contracts without becoming the platform-wide search system.

## 4. Architecture direction

The accepted architecture direction is provider-replaceable and first-party-state oriented.

The validated PR #1 candidate currently demonstrates, but `main` does not yet accept:

- TypeScript/Vite/MapLibre GL JS web client;
- Go application API;
- PostgreSQL/PostGIS spatial and user-owned state;
- GoreeCloud Identity-compatible browser Authorization Code + PKCE with no browser client secret;
- access-token/UserInfo subject validation;
- owner/editor/viewer authorization with PostgreSQL row-level security and a non-owner runtime role;
- same-origin Maps API browser access;
- Nominatim-compatible forward/reverse geocoding adapter;
- Valhalla-compatible route adapter;
- versioned public geographic-data releases and a read-only edge delivery boundary;
- owner-scoped Saved Places and shared collection primitives.

Those candidate capabilities are source evidence only until the independent-review and merge gates are satisfied.

## 5. Data and provider requirements

Provider interfaces must remain replaceable for map tiles/styles, geocoding, place data, routing, transit, traffic/incidents, imagery, terrain/elevation, street-level imagery, indoor mapping, and offline packaging.

Self-hostable/open-data approaches are preferred when quality, licensing, privacy, security, operational reliability, freshness, attribution, and coverage are adequate. Technical API compatibility does not constitute provider approval.

Production provider acceptance requires, as applicable:

- dataset/provider licensing and provenance review;
- attribution correctness;
- data-quality and geographic-coverage validation;
- credential and secret handling;
- SSRF/egress/DNS/network controls;
- bounded request, timeout, response-size, and rate-limit behavior;
- monitoring and incident handling;
- Privacy Shield and Wardveil evidence;
- rollback/recovery and provider-replacement procedures.

Private searches, routes, saved places, collections, Identity state, and precise personal location must never be mixed into a public map-data release plane.

## 6. Multi-user and authorization requirements

Maps is designed for multiple authenticated GoreeCloud users. Every user-owned object requires an explicit ownership/authorization boundary.

Shared resources use deliberate roles such as owner, editor, and viewer. Possession of an identifier or URL must not grant access. Human-friendly recipient lookup/invitations must use an approved GoreeCloud Identity consumer-directory contract when available; Maps must not use an administrative Identity directory as a consumer account browser.

Application authorization and database enforcement should be complementary. The validated PR #1 candidate includes PostGIS row-level-security evidence, but that remains candidate evidence until merge and does not prove production database acceptance.

## 7. Privacy and security requirements

Location and map activity can reveal sensitive personal behavior. Maps must minimize collection and disclosure.

Required principles include:

- no background tracking merely for convenience;
- no duplicate precise location-history store when GoreeCloud Location owns that data class;
- clear separation between map/search history and personal location-tracking history;
- user-controlled sharing, history, personalization, and revocation;
- no third-party advertising profiles, hidden analytics/tracking, remote fonts, or unnecessary browser dependencies;
- no silently enabled public tile/geocoder/router provider in a fresh development build;
- no reusable credentials, bearer tokens, OIDC codes, PKCE verifiers, private map contents, or precise coordinates in source control or ordinary logs;
- explicit authorization at every protected resource boundary;
- privacy-safe degraded/unavailable states instead of fabricated results.

## 8. GoreeCloud platform requirements

Stable qualification requires current accepted, substantive integration with the applicable GoreeCloud platform systems:

- **Glaze UI** — presentation, interaction, accessibility, adaptive layout/form-factor contracts, semantic material/effect behavior, and resilience modes.
- **Privacy Shield** — privacy-control, data-use, minimization, consent/permission, and truthful privacy-state contracts.
- **Wardveil Security** — evidence-backed security, protection, trust, abuse resistance, and response-state contracts.
- **Everkeep** — backup, recovery, preservation, portability, continuity, and deletion/export semantics for eligible user-owned Maps data.
- **GoreeCloud Mesh** — governed interoperability and capability coordination where platform architecture requires it.
- **GoreeCloud Identity** — authentication/principal identity without transferring Maps resource authorization.
- **GoreeCloud Location** — sensitive device/current-location capabilities without Maps creating a competing tracking/history authority.

Decorative names, badges, or artwork do not satisfy integration requirements.

## 9. Glaze UI and form-factor requirements

The current authoritative Maps documentation targets Glaze UI 2.0.0 Stable for the foundation candidate. Any later Glaze UI migration must be an explicit, reviewed Maps consumer change rather than an assumed version bump.

The map is the primary Canvas. Search, inspectors, place cards, navigation controls, sheets, toolbars, and account controls must use appropriate semantic surfaces/material roles rather than indiscriminate transparency.

Maps must support, as applicable, mobile and desktop compositions, touch/pointer/keyboard operation, visible focus, practical target sizes, light/dark appearance, reduced motion, reduced transparency, increased contrast, forced colors, and structured non-map alternatives for essential information.

Responsive web behavior does not by itself establish native mobile/tablet/desktop/TV/wearable/spatial acceptance.

## 10. Branding contract

Branding authority is `GoreeCloud/goreecloud-branding-assets`.

The approved canonical Maps asset is `products/maps/app-icon.svg`, Git blob `07b6e52e04c95e1ec9f703a9d323cf799481351c`. Platform derivatives may be introduced only when a real consumer surface exists and must remain traceable to the approved canonical source.

Brand colors and artwork are identity only. They must not substitute for Glaze UI semantic warning, danger, security, privacy, route-safety, or other state treatment.

## 11. License and third-party material

The repository currently carries GNU Affero General Public License version 3 material in `LICENSE`. Unless separately noted, GoreeCloud-owned repository source is intended to follow that repository license.

Third-party dependencies, datasets, imagery, transit feeds, map providers, geographic sources, fonts, and other incorporated material retain their own licenses/notices. Maps release readiness requires applicable licensing and attribution evidence; the repository AGPL license does not relicense third-party map/data content.

## 12. Stable and production-readiness gates

Maps remains Development until all applicable gates are satisfied together, including:

- executable foundation accepted through required review/merge governance;
- exact-release CI;
- approved live geographic data and provider contracts;
- licensing/provenance/attribution verification;
- production database and authorization/RLS acceptance;
- production GoreeCloud Identity acceptance;
- GoreeCloud Location acceptance where used;
- substantive current Privacy Shield, Wardveil Security, Everkeep, Mesh, and Glaze UI integration evidence;
- accessibility and representative form-factor/device acceptance;
- privacy/security review;
- monitoring, incident, backup/restore, disaster-recovery, and rollback procedures;
- deployment validation;
- explicit release and Stable-promotion authorization.

A green candidate CI run, branding completion, governance documentation, or technical provider compatibility does not independently satisfy these gates.
