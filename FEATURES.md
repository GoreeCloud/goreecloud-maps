# GoreeCloud Maps Features

Status vocabulary: **Implemented**, **In progress**, **Planned**, or **Blocked by prerequisite**. A listed feature is not a production claim unless marked Implemented and supported by acceptance evidence.

## Foundation

| Feature | Status | Notes |
|---|---|---|
| Original GoreeCloud application model | Implemented | Repository and architecture prohibit a complete third-party application fork. |
| Glaze UI 2.0.0 target | In progress | Design/interaction contract selected; app-specific implementation and acceptance required. |
| Replaceable geospatial providers | In progress | Provider boundaries are defined in specifications; adapters follow in implementation milestones. |
| Multi-user tenancy model | In progress | Owner/member/role model defined; API/database enforcement is the next implementation milestone. |
| GoreeCloud Location boundary | Implemented | Maps and Location responsibilities are explicitly separated. |

## Map experience

| Feature | Status |
|---|---|
| Interactive 2D vector map | Planned |
| Smooth pan, zoom, rotate, pitch | Planned |
| User-location control | Planned |
| Light, dark, and Deep Dark map appearances | Planned |
| Globe view | Planned |
| Terrain/elevation | Planned |
| 3D buildings | Planned |
| Custom layers and overlays | Planned |
| Map style switching | Planned |
| Accessible non-map/list alternatives | Planned |

## Search, places, and discovery

| Feature | Status |
|---|---|
| Forward geocoding | Planned |
| Reverse geocoding | Planned |
| Nearby/category search | Planned |
| Rich place cards | Planned |
| Favorites and saved places | Planned |
| Personal notes | Planned |
| Collections and guides | Planned |
| Recently viewed/search history controls | Planned |
| Place correction/feedback workflow | Planned |
| Search integration with GoreeCloud Search | Planned |

## Directions and navigation

| Feature | Status |
|---|---|
| Driving routes | Planned |
| Walking routes | Planned |
| Cycling routes | Planned |
| Transit routes | Planned |
| Multimodal itineraries | Planned |
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

## Multi-user and collaboration

| Feature | Status |
|---|---|
| Per-user private saved data | In progress |
| Shared collections | Planned |
| Owner/editor/viewer roles | In progress |
| Invitations and revocation | Planned |
| Collaborative map annotations | Planned |
| Shared route plans | Planned |
| Household/group collections | Planned |
| Revocable share links | Planned |
| Optional ETA/location overlays through GoreeCloud Location | Planned |
| Audit trail for security-sensitive sharing changes | Planned |

## Offline and resilience

| Feature | Status |
|---|---|
| Downloadable map regions | Planned |
| Offline place index | Planned |
| Offline routing graphs | Planned |
| Package integrity/version metadata | Planned |
| Background package updates | Planned |
| Storage quotas and cleanup | Planned |
| Stale/offline/degraded state presentation | Planned |

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
| GoreeCloud Identity | Planned |
| GoreeCloud Location | Planned |
| GoreeCloud Search | Planned |
| GoreeCloud Mesh | Planned |
| Privacy Shield | Planned |
| Wardveil Security | Planned |
| Everkeep | Planned |
| Glaze UI 2.0.0 | In progress |

## Form factors

Web is the first implementation target. Mobile, tablet, foldable, desktop, TV, wearable, and spatial interfaces are planned only where the applicable Glaze UI contract and platform-specific acceptance can be satisfied. A responsive web page alone does not constitute native or form-factor acceptance.