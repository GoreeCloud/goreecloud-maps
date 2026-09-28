# GoreeCloud Maps Specifications

## 1. Status and authority

GoreeCloud Maps is **Forge / Development**. `PROJECT-SPECIFICATIONS.md` is the canonical project specification once accepted on `main`; this `SPECIFICATIONS.md` file is the implementation-focused code-adjacent companion and must remain consistent with that authority. `PROJECT-RECORD.md` preserves significant project history.

Authoritative `main` remains the accepted governance/licensing/branding baseline until the current stabilization candidate passes review and protected promotion.

The current clean stabilization candidate is layered from current `main` and consolidates the previously stacked executable foundation, Saved Places, search-save, Glaze, and governance work. Candidate capabilities must not be described as accepted `main`, production, released, or Anchor behavior until promotion and all applicable acceptance gates complete.

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

The current stabilization candidate demonstrates, but authoritative `main` does not yet accept:

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

Those candidate capabilities are source evidence only until independent review and protected-promotion gates are satisfied.

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

Anchor qualification under Platform Contract 2.0 requires current accepted, substantive integration with exactly nine Integral Platform Systems:

- **GoreeCloud Manager** — lifecycle, administration, inventory, operational control, approvals, remediation, and management-plane visibility.
- **Privacy Shield** — purpose, minimization, consent/permission, retention, sharing, tracking privacy, and truthful privacy state.
- **Wardveil Security** — security, trust, abuse resistance, threat handling, response, and evidence.
- **Everkeep** — backup, restore, recovery, preservation, portability, migration readiness, and continuity.
- **Glaze UI** — presentation, interaction, accessibility, adaptive behavior, and evidence-state presentation.
- **GoreeCloud Mesh** — governed first-party discovery, dependency awareness, coordination, events, and evidence routing.
- **GoreeCloud Identity** — authentication and principal identity without transferring Maps resource authorization.
- **GoreeCloud Policy** — governed policy representation, evaluation, decisions, enforcement coordination, explanation, freshness, and evidence.
- **GoreeCloud Observability** — health, metrics, logs/events/traces, diagnostics, dependency health, correlation, freshness, and operational evidence.

GoreeCloud Sync remains separately governed and is not a tenth Integral Platform System. GoreeCloud Location remains an application/service dependency for sensitive current-position behavior, not an Integral Platform System.

Decorative names, badges, metadata, or artwork do not satisfy integration requirements.

## 9. Glaze UI and form-factor requirements

The stabilization source now implements a repository-local GLAZE UI V1.6 / 1.6.0 presentation adapter. It uses current shared source/release authority as a reference while keeping runtime behavior local, presentation-only, and caller-truth-owned. This closes the source-version migration gap but does not establish rendered, accessibility, adaptive, browser/device/GPU, performance, rollback, or production acceptance.

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

## 12. Anchor and production-readiness gates

Maps remains Forge/Development until all applicable gates are satisfied together, including:

- executable candidate accepted through required review/protected promotion;
- exact-head CI and current Platform Contract 2.0 validation;
- GLAZE UI V1.6 rendered/accessibility/adaptive/performance/browser-device acceptance for the implemented source mapping;
- approved live geographic data and provider contracts;
- licensing/provenance/attribution verification;
- production database and authorization/RLS acceptance;
- production GoreeCloud Identity acceptance;
- GoreeCloud Location acceptance where used;
- substantive Manager, Privacy Shield, Wardveil Security, Everkeep, Mesh, Policy, and Observability integration evidence;
- accessibility and representative form-factor/device/browser/GPU acceptance;
- privacy/security review;
- monitoring, incident, backup/restore, disaster-recovery, and rollback procedures;
- deployment validation;
- explicit release and Anchor authorization.

A green candidate CI run, branding completion, governance documentation, or technical provider compatibility does not independently satisfy these gates.
