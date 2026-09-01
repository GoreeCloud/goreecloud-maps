# GoreeCloud Maps Features

Status vocabulary: **Accepted main**, **Validated candidate**, **Planned**, or **Blocked by prerequisite**. Candidate status is not a production, release, or Stable claim.

## Accepted main

| Feature / record | Status | Evidence boundary |
|---|---|---|
| Original GoreeCloud Maps project identity | Accepted main | Canonical branding consumer contract is merged. |
| GNU AGPL v3 repository license material | Accepted main | Root `LICENSE` exists; third-party data/dependencies retain separate terms. |
| Product boundary vs GoreeCloud Location | Accepted main | Repository/canonical documentation separates Maps from sensitive Location authority. |
| Mandatory repository governance records | Accepted main after this governance change | README, specifications, features, benefits, competitive objectives, and branding are repository records; they do not imply executable acceptance. |

## Executable foundation candidate — PR #1

The following capabilities exist in draft PR #1 at validated head `a8aafba65b89dbcd76661368a605762871abbb35`. CI run `33251905892` passed, but the candidate remains unmerged pending its independent-review gate.

| Capability | Status |
|---|---|
| TypeScript/Vite/MapLibre web shell | Validated candidate |
| Map-as-Canvas responsive web composition | Validated candidate |
| Go Maps API | Validated candidate |
| PostgreSQL/PostGIS data foundation | Validated candidate |
| Owner/editor/viewer authorization model | Validated candidate |
| PostgreSQL row-level-security tests | Validated candidate |
| GoreeCloud Identity-compatible Authorization Code + PKCE browser boundary | Validated candidate |
| Access-token plus UserInfo subject verification | Validated candidate |
| Same-origin Maps API client | Validated candidate |
| Nominatim-compatible forward/reverse geocoding adapter | Validated candidate |
| Valhalla-compatible route adapter | Validated candidate |
| Provider capability status with explicit unconfigured state | Validated candidate |
| Shared collection/member/item API primitives | Validated candidate |
| Owner-scoped Saved Places API/client contract | Validated candidate |
| Versioned public geographic-data release contract | Validated candidate |
| Read-only Cloudflare Worker/R2-oriented edge source contract | Validated candidate |
| Privacy-safe local empty map style/no-provider fallback | Validated candidate |

These items are not accepted `main` capabilities until the candidate passes the remaining review/merge gates.

## Planned product capabilities

- approved live vector-tile/map-style infrastructure;
- rich place/POI data and dedicated nearby discovery;
- approved live geocoder and geographic-quality acceptance;
- approved live routing engine and route-quality acceptance;
- route alternatives, advanced constraints, and traffic-aware estimates;
- turn-by-turn navigation and rerouting;
- GoreeCloud Location current-position integration;
- traffic, closures, incidents, transit, terrain, 3D/globe, imagery, indoor mapping, and EV routing where licensed/approved data exists;
- offline regions, offline search/routing data, package freshness/integrity, quota management, and updates;
- complete Saved Places visual experience and cross-device synchronization;
- full collaborative membership/invitation UI, share links, ownership transfer, annotations, and shared route plans;
- GoreeCloud Search geographic-discovery interoperability;
- representative native/mobile/tablet/desktop/other form-factor clients where justified.

## Blocked by prerequisite

- production GoreeCloud Identity client/session acceptance;
- human-friendly invitation lookup until an approved Identity consumer-directory integration is accepted;
- production geographic providers until licensing, attribution, data quality, privacy, security, network, monitoring, and operational gates pass;
- Stable qualification until current Glaze UI, Privacy Shield, Wardveil Security, Everkeep, Mesh, Identity, Location, accessibility, recovery, deployment, and release gates are satisfied as applicable.

## Evidence rule

A feature listed here may be described as implemented only at the evidence level shown. Repository documentation must never turn an unmerged candidate, provider interface, branding asset, or planned platform integration into an accepted runtime claim.
