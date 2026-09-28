# GoreeCloud Maps — Privacy Policy

## Scope

This repository policy describes the intended privacy behavior of GoreeCloud Maps and the current Development candidate. It does not create provider or deployment claims that are not otherwise accepted.

## Privacy principles

Maps is privacy-first and data-minimizing. It must not sell location or map activity, build advertising profiles, include behavioral advertising, or silently introduce third-party tracking.

## Location authority

GoreeCloud Location—not Maps—owns current device/user position, background tracking, personal location history, Find My, geofences, and location-sharing permission state.

Maps may use a current position only through an approved Location contract. Maps must not enable background tracking merely to improve map convenience or create a duplicate precise personal-history store.

## Data Maps may process

Depending on enabled features and user action, Maps may process:

- geographic search/place queries;
- route origins, destinations, waypoints, and route results;
- Saved Places, notes, collections, membership, and collaboration state;
- map/application preferences;
- Identity subject references required for authentication/resource ownership;
- operational/security audit records that avoid unnecessary content.

Provider-returned place and route information remains subject to its applicable provenance/license.

## Provider disclosure

No public tile, geocoder, router, imagery, traffic, transit, or place provider should be contacted silently by a fresh development build.

When an external provider is approved and used, Maps must disclose the applicable data-transfer boundary and minimize what is transmitted. Provider credentials remain server-side.

## Public map-data separation

Public geographic releases must not contain private searches, routes, Saved Places, collaboration content, Identity state, or precise personal-location history.

## Retention, deletion, export, and recovery

User-owned Maps state requires explicit retention, deletion, export/portability, backup, restore, and migration semantics before production acceptance. Restore must respect deletion/retention policy and authorization.

## Logs and diagnostics

Ordinary logs must not contain bearer tokens, OIDC codes/verifiers, reusable credentials, provider secrets, private map contents, or raw precise coordinates unless a separately approved diagnostic path explicitly requires and protects them.

GoreeCloud Observability may receive minimized operational evidence under its authority; missing telemetry must not be presented as healthy state.

## Platform authorities

Privacy Shield is authoritative for shared privacy policy/decision state. Wardveil Security, Identity, Policy, Observability, Everkeep, Mesh, Manager, and Glaze UI retain their respective domains. A UI badge or local setting does not manufacture underlying privacy/security authorization.

## Current status

Maps is Forge/Development. Source implementation and CI do not establish production privacy acceptance. Provider, Identity, data-retention, continuity, platform-system, deployment, and Anchor gates remain open.
