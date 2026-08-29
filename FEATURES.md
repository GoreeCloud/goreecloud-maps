# GoreeCloud Maps Features

Status vocabulary: **Implemented**, **In progress**, **Planned**, or **Blocked by prerequisite**. A listed feature is not a production claim unless marked Implemented and supported by acceptance evidence.

## Foundation

| Feature | Status | Notes |
|---|---|---|
| Original GoreeCloud application model | Implemented | Repository and architecture prohibit a complete third-party application fork. |
| Glaze UI 2.0.0 target | In progress | Web shell applies the current material/layout/accessibility contract; product-specific rendered acceptance remains required. |
| Replaceable geospatial providers | In progress | Map-style seam plus Nominatim-compatible geocoding and Valhalla-compatible routing adapters are implemented; live provider deployment and acceptance remain pending. |
| Multi-user tenancy model | In progress | Identity-subject mapping, owner/member roles, RLS policies, collaboration APIs, runtime-role fail-closed checks, and automated PostGIS isolation acceptance are implemented; production Identity/database and broader UI acceptance remain pending. |
| Browser OIDC public-client integration | In progress | Authorization Code + PKCE source integration exists with no browser secret and in-memory bearer token handling; actual Maps client registration and end-to-end Identity acceptance remain pending. |
| Same-origin Maps API web client | In progress | Browser source calls a configurable same-origin API path and rejects arbitrary external API origins; deployment/proxy acceptance remains pending. |
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
| Web place-search results | In progress |
| Nearby/category search | In progress |
| Rich place cards | Planned |
| Favorites and saved places | In progress |
| Personal notes | In progress |
| Collections and guides | In progress |
| Recently viewed/search history controls | Planned |
| Place correction/feedback workflow | Planned |
| Search integration with GoreeCloud Search | Planned |

Forward and reverse geocoding have authenticated Maps API routes and a normalized Nominatim-compatible provider adapter. The web search form and category chips now call the same-origin Maps API after authentication, render normalized results, and can move the map/marker to a selected result. No live geocoder endpoint is configured or production-accepted, category behavior is currently ordinary text geocoding rather than a dedicated nearby/POI ranking contract, and GoreeCloud Search interoperability remains pending. Saved-place storage exists; saved-place API/UI workflows remain pending.

## Directions and navigation

| Feature | Status |
|---|---|
| Driving routes | In progress |
| Walking routes | In progress |
| Cycling routes | In progress |
| Transit routes | In progress |
| Multimodal itineraries | In progress |
| Directions web UI | In progress |
| Route geometry rendering | In progress |
| Maneuver summary UI | In progress |
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

The API contains a normalized Valhalla-compatible route adapter for drive, walk, bicycle, and transit/multimodal costing, including distance, duration, encoded shape, and maneuver data. The web source can geocode an origin/destination, request a route, decode/render returned polyline6 geometry, fit the map to it, and show a bounded maneuver summary. No live router endpoint, route alternatives, navigation runtime, traffic-aware ETA, route-quality acceptance, or real-device navigation acceptance exists yet.

## Multi-user and collaboration

| Feature | Status |
|---|---|
| OIDC-authenticated Maps users | In progress |
| Per-user private saved data | In progress |
| Shared collections | In progress |
| Shared collection web browsing | In progress |
| Collection creation web UI | In progress |
| Owner/editor/viewer roles | Implemented |
| PostgreSQL row-level security | Implemented |
| Runtime database privilege guard | Implemented |
| Collection membership API | Implemented |
| Collection item CRUD API | Implemented |
| Optimistic revision conflict handling | Implemented |
| Member removal/revocation | Implemented |
| Invitations and identity-directory resolution | Blocked by prerequisite |
| Full membership-management web UI | Planned |
| Collaborative map annotations | Planned |
| Shared route plans | Planned |
| Household/group collections | Planned |
| Revocable share links | Planned |
| Optional ETA/location overlays through GoreeCloud Location | Planned |
| Audit trail for security-sensitive sharing changes | In progress |

The authenticated collection API supports collection updates, member listing/addition/role changes/removal, collection-item list/create/update/delete, and revision-conflict responses. Membership and item mutations emit collection audit events. The web source can list/create collections and browse collection items after authentication.

Human-friendly recipient discovery and invitation delivery are **blocked by prerequisite** until GoreeCloud Identity defines an approved cross-application directory/invitation contract. Maps does not substitute the Identity provider's administrative user directory or create a parallel account-discovery system.

Automated integration acceptance runs the migration chain against PostGIS and exercises owner/editor/viewer/stranger isolation, private saved-place visibility, collection revision conflicts, editor item creation/update/deletion, viewer mutation denial, owner-only membership administration, role demotion/restoration, member self-removal, immutable collection ownership, collaboration audit records, forged-audit rejection, and the non-owner/no-`BYPASSRLS` runtime-role guard. This validates the current source authorization model in CI; it does not replace production database, GoreeCloud Identity SSO, load, backup, or deployment acceptance.

## Identity integration

| Feature | Status |
|---|---|
| Authorization Code + PKCE browser flow | In progress |
| No browser client secret | Implemented | Source and example configuration contain no client-secret mechanism. |
| Same-origin redirect enforcement | Implemented | Browser source rejects configured redirects outside the current application origin. |
| In-memory access-token handling | Implemented | Access token is not persisted by the source client; PKCE state/verifier are transient in `sessionStorage`. |
| API JWT issuer/audience/expiry verification | In progress |
| Provider UserInfo subject confirmation | In progress |
| Maps Identity client registration | Blocked by prerequisite |
| Production login/logout/session acceptance | Blocked by prerequisite |
| Disabled-account/outage/recovery/rollback acceptance | Blocked by prerequisite |

The source integration expects OAuth/OIDC bearer access tokens rather than using an ID token as an API bearer. Production claims require a real GoreeCloud Identity client registration and end-to-end runtime evidence; repository compilation alone does not prove SSO acceptance.

## Offline and resilience

| Feature | Status |
|---|---|
| Privacy-safe no-provider renderer fallback | Implemented |
| Explicit API/provider unavailable states | In progress |
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
