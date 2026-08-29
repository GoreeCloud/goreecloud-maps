# GoreeCloud Maps Architecture

## Status

This document describes the implemented foundation and intended service boundaries for GoreeCloud Maps. The application remains Development.

## Runtime shape

```text
Browser / future native client
            |
            v
      GoreeCloud Maps API
       /        |        \
      /         |         \
Identity     PostGIS     Provider adapters
  OIDC       user data    tiles / search /
                         routing / transit /
                         traffic / imagery
```

GoreeCloud Location is adjacent to this diagram rather than inside the Maps database. Maps may request approved current-location and sharing capabilities from Location, but personal location history, background tracking, Find My, geofences, and location-sharing authority remain owned by Location.

## Web application

`apps/web` is the first executable user interface. It uses MapLibre GL JS as a narrow rendering dependency and applies the GoreeCloud-controlled Glaze UI 2.0 experience around the renderer.

The default `public/map-style.json` has no remote sources. A live map only appears when `VITE_MAP_STYLE_URL` is set to an approved MapLibre-compatible style endpoint. This is intentional: cloning the repository must not silently disclose the user's IP address, viewport, or map activity to an unrelated public tile service.

## API

`services/api` is the initial protected application API. It requires:

- a reachable PostgreSQL/PostGIS database;
- a non-owner runtime database role without `BYPASSRLS`;
- an OIDC issuer compatible with the GoreeCloud Identity integration contract;
- a Maps OIDC client ID.

The API does not accept a user ID from the client as an authorization authority. It validates the bearer token, resolves the OIDC `sub` claim to an application-owned Maps user record, and then passes that internal user ID into a transaction-local PostgreSQL setting used by row-level-security policies.

The service intentionally exits at startup when the configured database role can own/bypass the RLS-protected tables.

## Database and multi-user authorization

`db/migrations/0001_multi_user_foundation.sql` introduces first-party spatial and collaborative state. The initial model includes users, preferences, saved places, collections, collection members, collection items, and security-sensitive audit events.

Map collections have one immutable owner and optional `editor` or `viewer` members. Ownership transfer is deliberately unsupported in the first migration; it requires a dedicated future transaction and audit contract.

The database owner/migration role and API runtime role must be separate. The RLS helper used to evaluate collection membership is a narrowly scoped security-definer function. Production database acceptance must verify that the runtime role neither owns the protected tables nor holds `BYPASSRLS`.

## Provider boundary

Provider integration is capability-based. Maps will not leak provider-specific response models through its first-party APIs or persistent user-data schema.

Planned adapter boundaries include:

- vector/base-map tiles and styles;
- forward and reverse geocoding;
- place search and POI enrichment;
- directions/routing;
- transit;
- traffic and incidents;
- terrain/elevation;
- aerial/satellite imagery;
- street-level imagery;
- indoor maps.

MapLibre is the selected initial renderer, not the map-data provider. PostGIS is the first-party user/spatial state store, not the sole source of public geographic truth. A self-hosted Valhalla-class router is the preferred initial routing direction, but the application-facing route contract remains replaceable.

## Identity boundary

GoreeCloud Identity establishes authentication. Maps owns authorization. Maps therefore stores its own user-resource ownership and collaborative roles rather than treating an Identity group claim as permanent authorization truth for every Maps object.

GoreeCloud Identity is still completing GoreeCloud-wide production SSO acceptance. Maps may develop against its OIDC-compatible integration contract without claiming that production SSO is already approved.

## Privacy boundary

Precise user location, origins/destinations, private route plans, saved places, and map searches can be sensitive. They must not become ordinary logs, analytics payloads, provider debug strings, cache keys exposed across users, or public CDN objects.

Public map resources and private user resources require separate cache and authorization policies. External provider calls must be documented with data-minimization and provenance requirements before production use.

## Resilience

The architecture separates public map/provider dependencies from durable user-owned data. A provider outage must not make previously saved collections disappear. Offline packages are planned as versioned resources with integrity, quota, update, rollback, and deletion behavior.

## Deployment boundary

No production deployment is defined by this foundation. Future deployment work must preserve the GoreeCloud platform boundaries for private/public networking, secrets, HTTPS, observability, Cloudflare Pages where applicable, backups, recovery, Privacy Shield, Wardveil Security, Everkeep, and GoreeCloud Mesh.
