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

The Worker configuration names the intended R2 resources for deployment, but this source repository does not prove those buckets or a Worker deployment currently exist. Cloudflare provisioning, custom-domain routing, cache behavior, observability, rollback, and production acceptance must be recorded separately when actually performed.

## Client integration

The current web client can consume an approved MapLibre-compatible style through `VITE_MAP_STYLE_URL`. A future client milestone may resolve the active style from the release manifest so roll-forward/rollback can occur without rebuilding the application shell.

Until an approved release is deployed, the repository-local empty style remains the privacy-safe default.
