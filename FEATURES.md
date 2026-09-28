# GoreeCloud Maps Features

Status vocabulary: **Accepted main**, **Stabilization candidate**, **Planned**, or **Blocked by prerequisite**. Candidate status is not a release, production, or Anchor claim.

## Accepted main

| Feature / record | Status |
|---|---|
| GoreeCloud Maps product identity and branding contract | Accepted main |
| GNU AGPL v3 repository license material | Accepted main |
| Maps vs GoreeCloud Location authority boundary | Accepted main |
| Repository governance records and Platform Contract history | Accepted main |

## Current stabilization candidate

| Capability | Status |
|---|---|
| TypeScript/Vite/MapLibre web shell | Stabilization candidate |
| Map-as-canvas responsive composition | Stabilization candidate |
| Go Maps API | Stabilization candidate |
| PostgreSQL/PostGIS data foundation | Stabilization candidate |
| Owner/editor/viewer authorization + RLS tests | Stabilization candidate |
| Identity-compatible Authorization Code + PKCE browser boundary | Stabilization candidate |
| Access-token/UserInfo subject verification | Stabilization candidate |
| Same-origin Maps API client | Stabilization candidate |
| Nominatim-compatible forward/reverse geocoding adapter | Stabilization candidate |
| Valhalla-compatible routing adapter | Stabilization candidate |
| Explicit provider-unconfigured/degraded state | Stabilization candidate |
| Owner-scoped Saved Places API/client workflow | Stabilization candidate |
| Search-result Save action | Stabilization candidate |
| Shared collection/member/item primitives | Stabilization candidate |
| Versioned public geographic-data release contract | Stabilization candidate |
| Read-only map-data edge source contract | Stabilization candidate |
| Privacy-safe local empty map style/no-provider fallback | Stabilization candidate |
| Exact-head CI and Platform Contract 2.0 validation | Stabilization candidate |
| Repository-local GLAZE UI V1.6 / 1.6.0 source mapping | Stabilization candidate / source implemented; downstream acceptance blocked |

## Planned / incomplete product capabilities

- Maps-specific GLAZE UI V1.6 rendered/accessibility/adaptive/performance/browser-device acceptance;
- approved live vector-tile/map-style infrastructure;
- rich POI/place data and nearby discovery;
- approved geocoder and routing quality acceptance;
- route alternatives, advanced constraints, traffic-aware estimates, turn-by-turn navigation, and rerouting;
- current-position integration through GoreeCloud Location;
- traffic, closures, incidents, transit, terrain, 3D/globe, imagery, indoor mapping, and EV routing where approved data exists;
- offline regions, offline search/routing data, freshness/integrity, storage quota, and update controls;
- complete Saved Places and collections UX, synchronization, invitations, ownership transfer, annotations, and shared route plans;
- GoreeCloud Search geographic-discovery interoperability;
- approved native/form-factor clients where justified.

## Blocked prerequisites

Anchor qualification remains blocked by incomplete Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, Mesh, Identity, Policy, and Observability acceptance; live geographic-provider approval; production database/RLS and Identity validation; Location integration where used; accessibility/form-factor/device/GPU evidence; continuity and restore evidence; deployment; rollback; release provenance; and production acceptance.

## Evidence rule

A feature may be described only at its verified evidence level. Repository documentation must never convert a candidate branch, provider seam, branding asset, or planned integration into accepted runtime behavior.
