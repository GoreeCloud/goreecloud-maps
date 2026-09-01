# GoreeCloud Maps Geographic Data Plane

## Status

This directory defines the source contract for versioned public geographic-data releases used by GoreeCloud Maps. It does **not** contain or claim an approved production map dataset.

The current public-delivery model is intentionally separate from private Maps user data. Map tiles, styles, glyphs, sprites, attribution, and public release metadata belong to the public geographic-data plane; saved places, collections, route plans, searches, Identity principals, and precise personal location data do not.

## Release shape

A release uses one immutable release identifier and the following object layout:

```text
manifests/current.json
releases/<release-id>/manifest.json
releases/<release-id>/style.json
releases/<release-id>/tiles/<z>/<x>/<y>.pbf
releases/<release-id>/sprites/<name>.json
releases/<release-id>/sprites/<name>.png
releases/<release-id>/glyphs/<fontstack>/<range>.pbf
```

`manifests/current.json` is the only intentionally mutable pointer. It receives a short cache lifetime. Every object below `releases/<release-id>/` is immutable and must never be replaced in place after publication. Corrections produce a new release identifier.

## Manifest contract

`release-manifest.schema.json` documents the machine-readable contract. `scripts/validate-mapdata-release.mjs` performs the repository acceptance checks currently enforced by CI.

A production release manifest must identify:

- release identifier and generation time;
- style and vector-tile paths;
- zoom and geographic coverage bounds;
- user-visible attribution requirements;
- source dataset names and versions;
- source licensing/provenance records;
- confirmation that the public delivery bundle contains public geographic data only.

The manifest does not authorize a dataset merely because it is syntactically valid. Licensing, attribution, privacy, provenance, update cadence, geographic quality, accessibility, operational capacity, and Wardveil/Privacy Shield acceptance remain separate gates.

## Style contract — schema v1

`mapdata/examples/style.example.json` is a synthetic style fixture and `scripts/validate-mapdata-style.mjs` enforces the initial style boundary in CI.

Map-data schema v1 is intentionally narrow:

- MapLibre Style Specification version 8;
- vector/MVT sources only;
- release-local `tiles/{z}/{x}/{y}.pbf` tile templates only;
- no TileJSON indirection;
- no imported styles;
- no external `font-faces` resources;
- optional glyphs must use the release-local `glyphs/{fontstack}/{range}.pbf` path;
- optional sprite resources must remain inside the immutable release path;
- non-background vector layers must reference a declared source and source-layer.

The web client independently validates this boundary before rendering a release. It fetches the manifest and style itself, bounds their sizes, refuses redirects, validates provenance/public-data state, converts release-local tile/glyph/sprite resources to the configured release origin, rejects off-origin/off-release resources, and injects the manifest’s validated attribution into the rendered vector source.

Raster/satellite imagery, terrain sources, external style imports, remote font-face resources, PMTiles/other archive packaging, and other future source types require an explicit later schema version and their own privacy, licensing, performance, offline, and operational acceptance. They are not silently accepted by schema v1.

## Release generation boundary

A future ingestion/build pipeline may use OpenStreetMap and other approved datasets, but it must not place raw provider secrets, personal user data, private route/search records, or unapproved proprietary content into the public release bundle.

The pipeline must be reproducible enough to record:

1. exact source dataset/version or acquisition snapshot;
2. license and attribution obligations;
3. transformations applied;
4. tile/style build version;
5. validation results;
6. release identifier;
7. publication/rollback status.

No release is considered published merely because files exist in a repository or object store.

## Cloudflare delivery boundary

`services/mapdata-edge` is the initial read-only Cloudflare Worker/R2 delivery implementation. It exposes only the current manifest and allowlisted immutable release object shapes. It does not provide object writes, bucket listing, arbitrary key access, or private-user data access.

The Worker configuration names the intended R2 resources for deployment, but this source repository does not prove those buckets or a Worker deployment currently exist. Cloudflare provisioning, custom-domain routing, cache behavior, observability, rollback, production CORS policy, and production acceptance must be recorded separately when actually performed.

## Client integration

The preferred web configuration is `VITE_MAP_DATA_MANIFEST_URL`. When configured, Maps starts with the repository-local empty style, retrieves and validates the public release manifest and release style, and only then replaces the renderer style with the validated in-memory MapLibre style. This allows a mutable current-manifest pointer to roll forward or back without rebuilding the application shell while immutable release objects remain cacheable.

`VITE_MAP_STYLE_URL` remains a manual/legacy approved-style seam and is used only when the manifest setting is empty. A configured release-manifest failure does not fall through to an unrelated style: Maps keeps the privacy-safe local empty style and exposes a degraded/unavailable state.

Until an approved release is actually deployed and accepted, the repository-local empty style remains the default and no live geographic coverage is claimed.
