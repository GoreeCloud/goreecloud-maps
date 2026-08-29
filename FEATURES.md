# GoreeCloud Maps Features

Status vocabulary: **Implemented**, **In progress**, **Planned**, or **Blocked by prerequisite**. A listed feature is not a production claim unless marked Implemented and supported by acceptance evidence.

## Foundation

| Feature | Status | Notes |
|---|---|---|
| Original GoreeCloud application model | Implemented | Repository and architecture prohibit a complete third-party application fork. |
| Glaze UI 2.0.0 target | In progress | Web shell applies the current material/layout/accessibility contract; product-specific rendered acceptance remains required. |
| Replaceable geospatial providers | In progress | Map-style seam plus Nominatim-compatible geocoding and Valhalla-compatible routing adapters are implemented; live provider deployment and acceptance remain pending. |
| Multi-user tenancy model | In progress | Identity-subject mapping, owner/member roles, RLS policies, API runtime-role fail-closed checks, collection primitives, and automated PostGIS isolation acceptance are implemented; broader collaboration and production-runtime acceptance remain pending. |
| GoreeCloud Location boundary | Implemented | Maps and Location responsibilities are explicitly separated. |

## Map experience

| Feature | Status |
|---|---|
| Interactive 2D vector map | In progress |
| Smooth pan, zoom, rotate, pitch | In progress |
| User-location control | Planned |
| Light, dark, and Deep Dark map appearances | In progress |
| Globe view | Planned |
| Terrain/elevation | Planned |
| 3D buildings | Planned |
| Custom layers and overlays | Planned |
| Map style switching | Planned |
| Accessible non-map/list alternatives | Planned |

The current renderer shell is functional but uses a local data-empty style unless an approved `VITE_MAP_STYLE_URL` is configured. It therefore does not yet claim live geographic map coverage.

## Search, places, and discovery

| Feature | Status |
|---|---|
| Forward geocoding | In progress |
| Reverse geocoding | In progress |
| Nearby/category search | Planned |
| Rich place cards | Planned |
| Favorites and saved places | In progress |
| Personal notes | In progress |
| Collections and guides | In progress |
| Recently viewed/search history controls | Planned |
| Place correction/feedback workflow | Planned |
| Search integration with GoreeCloud Search | Planned |

Forward and reverse geocoding now have authenticated Maps API routes and a normalized Nominatim-compatible provider adapter. No live geocoder endpoint is configured or production-accepted, and the web search UI is not yet connected to this API. Saved-place and collection storage primitives exist in the first migration; saved-place API/UI workflows remain pending.

## Directions and navigation

| Feature | Status |
|---|---|
| Driving routes | In progress |
| Walking routes | In progress |
| Cycling routes | In progress |
| Transit routes | In progress |
| Multimodal itineraries | In progress |
| Route alternatives | Planned |
| Departure/arrival-time planning | Planned |
| Avoid tolls/highways/ferries and other supported constraints | Planned |
| Elevation-aware cycling/walking | Planned |
| Turn-by-turn navigation | Planned |
| Dynamic rerouting | Planned |
| Lane/maneuver guidance | Planned |
| Speed-limit presentation | Planned |
| Incident/closure warnings | Planned |
| EV charging-aware routing | Planned |

The API now contains a normalized Valhalla-compatible route adapter for drive, walk, bicycle, and transit/multimodal costing, including distance, duration, encoded shape, and maneuver data. No live router endpoint, directions UI, navigation runtime, or route-quality acceptance exists yet.

## Multi-user and collaboration

| Feature | Status |
|---|---|
| OIDC-authenticated Maps users | In progress |
| Per-user private saved data | In progress |
| Shared collections | In progress |
| Owner/editor/viewer roles | In progress |
| PostgreSQL row-level security | Implemented |
| Runtime database privilege guard | Implemented |
| Invitations and revocation | Planned |
| Collaborative map annotations | Planned |
| Shared route plans | Planned |
| Household/group collections | Planned |
| Revocable share links | Planned |
| Optional ETA/location overlays through GoreeCloud Location | Planned |
| Audit trail for security-sensitive sharing changes | In progress |

Automated integration acceptance now runs the migration against PostGIS and exercises owner/editor/viewer/stranger isolation, private saved-place visibility, editor mutation, viewer mutation denial, membership authority, member self-removal, immutable collection ownership, and the non-owner/no-`BYPASSRLS` runtime-role guard. This validates the current source authorization model in CI; it does not replace production database, GoreeCloud Identity SSO, load, backup, or deployment acceptance.

## Offline and resilience

| Feature | Status |
|---|---|
| Privacy-safe no-provider renderer fallback | Implemented |
| Downloadable map regions | Planned |
| Offline place index | Planned |
| Offline routing graphs | Planned |
| Package integrity/version metadata | Planned |
| Background package updates | Planned |
| Storage quotas and cleanup | Planned |
| Stale/offline/degraded state presentation | In progress |

## Data and community

| Feature | Status |
|---|---|
| Open map data ingestion | Planned |
| Source attribution and provenance | Planned |
| User-submitted place feedback | Planned |
| Road closure/hazard reports | Planned |
| Community incident reports | Planned |
| Moderation/abuse controls | Planned |
| Transit feed ingestion | Planned |
| Traffic feed ingestion | Planned |
| Indoor map datasets | Planned |
| Street-level imagery provider | Blocked by prerequisite |
| Satellite/aerial imagery provider | Blocked by prerequisite |

## GoreeCloud platform integration

| Integration | Status |
|---|---|
| GoreeCloud Identity | In progress |
| GoreeCloud Location | Planned |
| GoreeCloud Search | Planned |
| GoreeCloud Mesh | Planned |
| Privacy Shield | Planned |
| Wardveil Security | Planned |
| Everkeep | Planned |
| Glaze UI 2.0.0 | In progress |

## Form factors

The responsive web shell currently covers initial mobile and desktop composition work. Native mobile, tablet, foldable, desktop, TV, wearable, and spatial acceptance remain planned and must satisfy the applicable Glaze UI contract and platform-specific evidence requirements. Responsive web behavior alone does not constitute native or form-factor acceptance.
