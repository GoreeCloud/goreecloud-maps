# GoreeCloud Maps Benefits

## Users

GoreeCloud Maps is intended to provide one coherent place for map exploration, place discovery, directions, saved places, route planning, offline regions, and collaboration without requiring a Google or Apple account.

The multi-user model keeps personal map content private by default while allowing deliberate sharing of collections, annotations, and route plans. Users can revoke access and retain clear ownership instead of relying on opaque possession-based links.

Offline packages and degraded-mode design improve usefulness when connectivity is limited, providers are unavailable, or the device is operating under restrictive network conditions.

## Privacy

Maps is designed around data minimization. Personal location history is not duplicated from GoreeCloud Location merely because Maps can display a current position. Search, route, place, and location-derived data receive separate retention and sharing rules.

A self-hostable provider architecture reduces the need to send precise map queries, origins, destinations, or location coordinates to unrelated third parties. When an external provider is required, the application can disclose and constrain that boundary rather than embedding it throughout the product.

## Security

Explicit server-side ownership and membership checks make the multi-user model auditable. Provider credentials stay server-side, private map resources use distinct caching rules, and sensitive location/search data is excluded from logs unless a documented operational requirement permits it.

Wardveil integration provides a common GoreeCloud security/protection boundary instead of creating a Maps-specific security identity.

## Accessibility and usability

Glaze UI 2.0 provides a common adaptive interaction grammar while still allowing Maps to be map-centric. The product is designed for touch, pointer, keyboard, large text, reduced motion, reduced transparency, increased contrast, and effects-free operation.

The map is not the only representation of essential information. Place results, route steps, saved places, and navigation instructions can be exposed through structured list/detail alternatives for assistive technology and situations where the visual map is not usable.

## Platform

Maps centralizes reusable map capabilities for GoreeCloud while preserving clear boundaries with Location, Search, Identity, Mesh, Privacy Shield, Wardveil Security, and Everkeep.

The provider architecture allows GoreeCloud to improve routing, geocoding, map data, traffic, transit, imagery, and tile delivery independently. A better provider can be adopted without rewriting the entire application or changing user-owned data formats.

## Operators and administrators

Self-hostable spatial data, routing, and geocoding components give GoreeCloud control over retention, observability, caching, capacity planning, regional coverage, incident response, and backup policies.

Source provenance and versioned offline/map datasets improve reproducibility and rollback. Provider health can be monitored independently, so a geocoding failure does not need to make saved collections or cached maps unavailable.

## Long-term resilience

Open formats, PostGIS-backed first-party state, provider adapters, and Everkeep portability reduce the risk of a single vendor discontinuation making GoreeCloud map data unusable. User-created collections and annotations remain portable even if the underlying map renderer or upstream data provider changes.