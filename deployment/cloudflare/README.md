# GoreeCloud Maps — Cloudflare Deployment Contract

## Status

This document defines the intended Cloudflare deployment shape for GoreeCloud Maps. It is a deployment contract, not evidence that a Maps Cloudflare Pages project, Worker, R2 bucket, route, custom domain, or production deployment currently exists.

## Application shell — Cloudflare Pages

The current web application is a Vite project under `apps/web`.

Intended Pages build settings:

- repository: `GoreeCloud/goreecloud-maps`;
- production branch: `main` after the normal review/merge gate;
- root directory: `apps/web`;
- build command: `npm install --ignore-scripts --no-audit --no-fund && npm run build`;
- build output directory: `dist`;
- preview deployments: pull-request branches when enabled by the connected Pages project.

`apps/web/public/_headers` supplies baseline static-response hardening and immutable caching for Vite fingerprinted assets. The application HTML and the local fallback style remain revalidated rather than treated as immutable.

Pages deployment must not inject a browser client secret. Any GoreeCloud Identity browser integration remains an OIDC public-client registration using Authorization Code + PKCE.

## Public geographic data — Worker + R2

Large versioned geographic assets are kept separate from the application shell.

`services/mapdata-edge` defines a read-only Worker with an R2 binding named `MAPS_DATA`. Its intended responsibilities are limited to:

- current public release-manifest reads;
- immutable versioned release-manifest reads;
- immutable MapLibre style reads;
- vector-tile reads;
- sprite reads;
- glyph reads;
- public CORS/security/cache headers;
- explicit health response.

The Worker does not implement write operations, R2 bucket listing, arbitrary object-key proxying, private Maps API access, Identity/session handling, saved-place access, collection access, search history, or personal location access.

The checked-in Wrangler configuration names intended R2 resources. Those names are desired deployment resources only; their presence in source does not prove provisioning.

## Cache contract

The delivery model uses two cache classes:

1. `manifests/current.json` is mutable and receives a short browser/CDN cache lifetime so a release pointer can roll forward or back.
2. `releases/<release-id>/...` is immutable and receives a one-year immutable cache lifetime. Published objects beneath a release identifier must never be replaced in place.

A rollback changes the current manifest pointer to a previously accepted release; it does not mutate the older release objects.

## Geographic release publication

Before publication, a release must pass the repository manifest validator plus separate operational acceptance for:

- data licensing and redistribution terms;
- required attribution text/link behavior;
- source provenance and dataset version;
- supported geographic bounds and zooms;
- tile/style compatibility;
- map rendering and label quality;
- accessibility considerations;
- abuse/rate/capacity controls;
- cache and rollback behavior;
- public/private-data separation;
- Privacy Shield and Wardveil Security evidence.

The current repository contains only a synthetic validation fixture. It does not contain an approved map dataset.

## Same-origin API boundary

The browser Maps API client currently uses same-origin `/api/v1` paths by default. A future Cloudflare production deployment must preserve that browser contract through a reviewed routing/proxy design or an equivalent accepted same-origin arrangement. The static Pages project must not accidentally rewrite `/api/*` to the SPA shell.

No such production API route is claimed by this source contract.

## Cloudflare acceptance evidence

When the connected Cloudflare control plane is used for Maps, record at minimum:

- Pages project identity and deployment revision;
- Worker deployment revision;
- R2 resource identities;
- routes/custom domains actually configured;
- environment variables and secret *names* without recording reusable secret values;
- deployment/preview status;
- cache/header checks;
- rollback test;
- production smoke test;
- any Cloudflare security or observability controls relied upon.

Deployment success is a separate lifecycle event from source CI success.
