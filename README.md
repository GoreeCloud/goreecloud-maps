# GoreeCloud Maps

GoreeCloud Maps is the first-party GoreeCloud mapping, place-discovery, directions, navigation, saved-place, offline-map, and collaborative-map application and service.

**Lifecycle:** Development

## Current repository state

The accepted `main` baseline is currently a governance, licensing, and product-identity foundation. It contains the approved GoreeCloud Maps branding contract and AGPL-3.0 license, but it does **not** yet contain the executable Maps application foundation.

The executable web/API/PostGIS foundation remains in draft pull request #1 (`agent/maps-foundation`). Its exact validated head `a8aafba65b89dbcd76661368a605762871abbb35` passed CI run `33251905892`, but the pull request remains unmerged because its independent-review gate is still outstanding. Capabilities implemented only in that candidate must not be represented as accepted `main` behavior.

## Product boundary

Maps owns map exploration, place and address discovery, directions and route planning, navigation presentation, saved places, collections, collaborative maps, offline map experiences, map layers, and map-provider orchestration.

GoreeCloud Location remains authoritative for device/current-position services, personal location history, background tracking, Find My capabilities, geofences, and location-sharing permission state. Maps consumes approved Location capabilities rather than creating a competing sensitive-location authority.

## Architecture direction

The validated foundation candidate establishes the intended architecture without making it accepted runtime state:

- TypeScript/Vite and MapLibre GL JS for the web map experience;
- Go for the Maps application API;
- PostgreSQL with PostGIS for first-party spatial/user state;
- GoreeCloud Identity-compatible Authorization Code + PKCE for the browser public-client boundary;
- application-owned owner/editor/viewer authorization plus PostgreSQL row-level security;
- replaceable Nominatim-compatible geocoding and Valhalla-compatible routing adapters, disabled unless explicitly configured;
- a versioned public geographic-data release contract with no private Maps/user data in the public release plane.

No live basemap, geocoder, router, production Identity client, production PostGIS deployment, Cloudflare application deployment, or production dataset is accepted merely because the candidate source supports the relevant interface.

## GoreeCloud platform requirements

Stable qualification requires substantive, evidence-backed integration with the applicable current GoreeCloud platform systems, including:

- Glaze UI for presentation, interaction, accessibility, adaptive behavior, and supported form factors;
- Privacy Shield for privacy controls and truthful data-use/privacy state;
- Wardveil Security for security and protection evidence;
- Everkeep for backup, recovery, preservation, portability, and continuity of eligible user-owned Maps data;
- GoreeCloud Mesh for governed interoperability where required;
- GoreeCloud Identity for authentication/principal identity while Maps retains resource authorization;
- GoreeCloud Location for sensitive device/location capabilities.

Branding or documentation alone does not satisfy those integrations.

## Canonical identity

Branding authority is `GoreeCloud/goreecloud-branding-assets`. The approved Maps product identity is `products/maps/app-icon.svg`, canonical Git blob `07b6e52e04c95e1ec9f703a9d323cf799481351c`.

The folded-map/route identity is intentionally distinct from GoreeCloud Location's positioning/pin identity. Current `main` has no merged executable/PWA launcher surface, so no platform-specific icon derivative is fabricated yet.

See [BRANDING.md](BRANDING.md) for the complete identity and derivative rules.

## Repository governance records

- [SPECIFICATIONS.md](SPECIFICATIONS.md) — durable product, architecture, privacy/security, integration, and acceptance contract.
- [FEATURES.md](FEATURES.md) — evidence-scoped feature inventory.
- [BENEFITS.md](BENEFITS.md) — intended user, privacy, operational, accessibility, and resilience benefits.
- [COMPETITIVE-OBJECTIVES.md](COMPETITIVE-OBJECTIVES.md) — capability and differentiation objectives.
- [BRANDING.md](BRANDING.md) — canonical visual-identity consumer contract.

The authoritative project record is `GoreeCloud/Projects/Project Specification — Maps`; chronological implementation evidence is recorded in `GoreeCloud/Changelogs/Change Log — Maps`.

## License

Unless otherwise noted, GoreeCloud Maps repository source is licensed under the GNU Affero General Public License version 3. Third-party dependencies, map datasets, geographic providers, imagery, transit feeds, and other incorporated material retain their own applicable licenses and attribution requirements.
