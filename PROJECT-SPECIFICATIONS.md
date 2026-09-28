# GoreeCloud Maps — Project Specifications

**Repository:** `GoreeCloud/maps`  
**Project type:** First-party GoreeCloud mapping application and service  
**Lifecycle:** Forge  
**Deployment state:** Development  
**Platform Contract:** 2.0  
**License:** GNU AGPL v3 for GoreeCloud-owned repository source unless separately noted  
**Project authority:** This file is the canonical project specification once accepted on `main`.

## Migration source and precedence

This specification migrates and reconciles the frozen Google Drive source `Project Specification — Maps.docx` (Drive file `1787QSgAEKtngF6rjHmcdvILpoK2wgRlV`, created September 11, 2026) with current repository evidence.

The Drive document is a migration input only. Current verified GitHub implementation evidence controls claims about what is implemented. Historical source statements naming `GoreeCloud/goreecloud-maps`, Glaze UI 2.0, Development/Stable lifecycle terminology, draft PR #1 as the only executable candidate, or Drive changelog/project authority are preserved in Git history and project record as historical context rather than current authority.

## 1. Purpose

GoreeCloud Maps is the first-party GoreeCloud mapping, place-discovery, directions, navigation, saved-place, offline-map, and collaborative-map application and service.

Google Maps and Apple Maps may be used as capability and interaction references only. GoreeCloud must not copy or scrape proprietary map data, place data, reviews, imagery, routing/navigation data, interface assets, wording, or other protected content.

## 2. Product and authority boundary

Maps owns map exploration and presentation, geographic/place discovery, place details, directions and route planning, navigation presentation, Saved Places, collections, collaborative maps, offline map experiences, map layers, and provider orchestration.

GoreeCloud Location remains authoritative for device/current-position services, background tracking, personal location history, Find My capabilities, geofences, location-sharing permission state, and sensitive tracking/sharing policy. Maps may consume approved Location capabilities only through explicit contracts and must not create a competing tracking/history authority.

GoreeCloud Identity authenticates principals. Maps remains responsible for Maps resource authorization. GoreeCloud Search may provide wider discovery capabilities without becoming the owner of Maps resources.

## 3. Required experience

The map is the primary spatial canvas. The product should provide:

- responsive pan, zoom, rotate, tilt, terrain/globe/3D interaction when approved data exists;
- address, place, category, natural-language, and nearby discovery;
- rich place details with provenance and licensed information;
- driving, walking, cycling, transit, and future multimodal route planning;
- route alternatives, distance, estimated arrival, constraints, warnings, and degraded-state truth;
- turn-by-turn navigation and rerouting when accepted runtime capabilities exist;
- traffic, incidents, transit, terrain, imagery, street-level and indoor capabilities only when approved;
- EV routing and charging discovery where trusted data exists;
- Saved Places, favorites, lists, guides, and collaborative collections;
- offline regions with explicit storage, freshness, update, integrity, and failure state;
- privacy-conscious recents/history/personalization that remain user-controlled;
- appropriate web, mobile, tablet, desktop, TV, wearable, spatial, or future form-factor experiences only where implementation and applicable Glaze UI contracts exist.

Planned capabilities are not implementation claims.

## 4. Current candidate architecture

The current stabilization candidate uses:

- TypeScript, Vite, and MapLibre GL JS for the web application;
- Go for the first-party Maps API;
- PostgreSQL with PostGIS for first-party spatial and user-owned state;
- GoreeCloud Identity-compatible Authorization Code + PKCE for the browser public-client boundary;
- application-owned owner/editor/viewer authorization plus PostgreSQL row-level security;
- replaceable Nominatim-compatible forward/reverse geocoding;
- replaceable Valhalla-compatible routing;
- versioned public geographic-data release contracts;
- a read-only edge-delivery boundary for approved public releases;
- a repository-local GLAZE UI V1.6 / 1.6.0 presentation adapter.

A fresh development build must not silently contact an unrelated public tile, geocoder, or router provider.

## 5. Multi-user and authorization requirements

Every user-owned resource requires an explicit ownership/authorization boundary. Shared resources use deliberate roles such as owner, editor, and viewer. Possession of an identifier or URL must not grant access.

Maps must derive principal identity from authenticated credentials, not request-supplied user identifiers. Application authorization and database row-level security should be complementary. Runtime database roles must not own protected tables or possess `BYPASSRLS`.

Human-friendly collaboration lookup must use an approved consumer-facing GoreeCloud Identity contract. An administrative identity directory must not be exposed as a consumer account browser.

## 6. Data and storage

PostgreSQL/PostGIS is the first-party Maps state foundation for eligible user and spatial state. Current foundation domains include Identity-subject mapping, user preferences, Saved Places, collections, memberships, collection items, and audit events.

Private searches, routes, Saved Places, collaboration state, Identity state, and precise personal location must remain separated from the public geographic-data release plane.

User-created state must have defined export, deletion, retention, backup, restore, portability, and migration behavior before Anchor qualification.

## 7. Provider architecture

Provider interfaces must remain replaceable for:

- vector/raster map tiles and styles;
- geocoding/reverse geocoding;
- place/POI discovery and details;
- routing and navigation inputs;
- traffic and incidents;
- transit;
- terrain/elevation;
- imagery and street-level imagery;
- indoor mapping;
- EV/charging information;
- offline packaging and synchronization.

Production provider acceptance requires applicable licensing/provenance, attribution, data-quality/coverage, credential handling, SSRF/egress/DNS controls, bounded requests/timeouts/response sizes/rate limits, privacy/security review, monitoring, incident handling, degradation behavior, and rollback/provider-replacement evidence.

Technical API compatibility is not provider approval.

## 8. Glaze UI and accessibility

The current approved target is GLAZE UI V1.6 / 1.6.0. The stabilization candidate implements a repository-local V1.6 presentation adapter without a remote browser runtime dependency.

Glaze UI remains presentation-only. It must not infer authorization, grant consent, request permissions automatically, manufacture provider/security/privacy truth, navigate automatically from presentation context, or execute consequential/fallback actions automatically.

Maps must provide, as applicable:

- semantic canvas/surface hierarchy rather than indiscriminate transparency;
- protected semantic surfaces for important state;
- touch, pointer, and keyboard interaction;
- visible focus;
- practical interaction targets;
- Light/Dark and supported appearance behavior;
- Reduced Motion and Reduced Transparency;
- increased contrast and forced colors;
- large-text/reflow and RTL/localization resilience;
- structured non-map alternatives for essential information;
- explicit loading, empty, stale, offline, degraded, permission, authentication, warning, privacy, security, and recovery states.

Current source mapping does not establish rendered or production conformance. Maps-specific browser/device/GPU, accessibility, performance, rollback, and Human Visual Excellence evidence remains required.

## 9. Privacy requirements

Location and map activity can reveal sensitive personal behavior. Maps must minimize collection, transfer, duplication, and retention.

Required rules include:

- no background tracking merely for convenience;
- no duplicate precise personal location-history store when GoreeCloud Location owns that data class;
- clear separation between map/search activity and Location tracking history;
- explicit, understandable sharing and history controls;
- no advertising profiles, behavioral advertising, hidden tracking, or unnecessary third-party analytics;
- no remote font requirement or unnecessary external UI dependency;
- no sensitive coordinates, bearer tokens, OIDC codes/verifiers, provider secrets, or private map content in ordinary logs;
- privacy-safe degraded/unavailable state instead of fabricated results.

Privacy Shield remains authoritative for applicable privacy policy/decision state.

## 10. Security requirements

Maps must apply least privilege, explicit authorization, secure session/token handling, input validation, bounded provider requests, dependency/supply-chain controls, secure transport, credential isolation, vulnerability handling, and security-safe diagnostics.

Wardveil Security remains authoritative for shared security/protection truth. Maps UI must not claim protection that the underlying security authority has not established.

## 11. Integral Platform Systems

Platform Contract 2.0 evaluates exactly nine Integral Platform Systems:

1. GoreeCloud Manager
2. Privacy Shield
3. Wardveil Security
4. Everkeep
5. Glaze UI
6. GoreeCloud Mesh
7. GoreeCloud Identity
8. GoreeCloud Policy
9. GoreeCloud Observability

GoreeCloud Sync is separately governed and is not a tenth Integral Platform System. GoreeCloud Location is an application/service dependency for sensitive location capability, not an Integral Platform System.

Each applicable integration must define purpose, interface, authentication, permissions, exchanged data, failure/degraded behavior, compatibility, and evidence. Labels or badges do not establish integration.

## 12. Offline, resilience, and continuity

Offline regions should preserve useful map data, Saved Places, and approved search/routing data to the degree supported by the chosen region/device. Storage use, freshness, integrity, update state, expiry, quota, and failures must be explicit.

Everkeep requirements apply to eligible user-created Maps state. Restore must preserve ownership and authorization and must not resurrect data contrary to deletion/retention policy.

Required and optional dependency failures must be explicit and diagnosable. Local functionality should remain available where safe and technically practical.

## 13. Deployment and operations

Production deployment is separately gated. Runtime configuration must keep secrets out of client bundles and source control. Health/readiness, dependency health, diagnostics, backup/restore, incident handling, rollback, release provenance, and operational runbooks must be validated before Anchor.

Cloudflare-oriented source contracts in the repository are deployment scaffolding only until a specific deployment is verified and accepted.

## 14. Testing and acceptance

Evidence categories remain separate. Source/build tests do not prove rendered accessibility; shared Glaze qualification does not prove Maps conformance; a provider seam does not prove provider acceptance; a merge does not prove deployment or Anchor.

Applicable gates include:

- exact-head/release CI;
- web build and API tests;
- live PostGIS/RLS isolation tests;
- public map-data contract validation;
- provider licensing/quality/privacy/security acceptance;
- production Identity and database acceptance;
- GoreeCloud Location integration acceptance where used;
- substantive nine-system Platform Contract integration evidence;
- rendered Glaze UI and accessibility acceptance;
- representative browser/device/GPU/performance evidence;
- backup/restore and continuity evidence;
- security/privacy review;
- deployment and rollback validation;
- release provenance and explicit Anchor qualification.

## 15. Current implementation status

Authoritative `main` remains the accepted baseline until the active stabilization candidate is reviewed and promoted.

Draft PR #11 consolidates the executable Maps foundation, Saved Places/search-save workflow, provider seams, PostGIS/RLS, public map-data edge contracts, Platform Contract 2.0, and the repository-local GLAZE UI V1.6 source mapping on current-main ancestry.

Candidate status must not be represented as accepted main, deployed, production, or Anchor state.

## 16. Maintenance and retirement

This file must be updated with material scope, architecture, authority, privacy/security, integration, deployment, or acceptance changes. Significant historical events belong in `PROJECT-RECORD.md`; feature state belongs in `IMPLEMENTED-FEATURES.md` and `PLANNED-FEATURES.md`; release/change history belongs in `CHANGELOGS.md`.

Deprecation or retirement must preserve export/migration, historical evidence, security/privacy obligations, and user-data continuity.
