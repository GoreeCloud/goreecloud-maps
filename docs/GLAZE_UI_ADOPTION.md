# GoreeCloud Maps — GLAZE UI V1.1 Adoption

## Status

Development source migration only. This document does not establish downstream conformance, Stable qualification, deployment, geographic-provider acceptance, or production approval.

## Current authority

- Target: **GLAZE UI V1.1 (`1.1.0`)**
- Stable tag: `v1.1.0`
- Stable release revision: `15cc76d2bcd4065552dc31c77145b63f34d9e7b2`
- Approved visual source: `8ea1f789bbabf943c3359514dc1506b24fa3c51b`
- Optical contract: `contracts/v1.1/optical-refinement.json`
- Atmosphere tokens: `tokens/glaze-v1.1-atmosphere.json`

## Known Stable-line blocker

The published `1.1.0` source remains the current Stable consumer target, but it has a known V1.1 CSS import-closure defect. The governed corrective line is GLAZE UI PR #129 / `1.1.1-rc.1`; that correction is still a Release Candidate with `consumerEligible: false` and is not a corrected immutable Stable release.

Maps may retain this bounded Development source mapping for review, but this `1.1.0` pin must not be used to claim current GLAZE UI conformance, release eligibility, geographic-provider acceptance, deployment acceptance, or production approval. After a corrected immutable Stable release is published, Maps must explicitly re-pin the exact version and revision and repeat all applicable source, rendered map/chrome, accessibility, representative-browser/device/GPU, deployment, and production acceptance.

## Web mapping

Maps keeps the live geographic renderer as the durable working canvas and confines Glaze treatment to bounded application chrome such as search, account, map controls, sheets, cards, directions controls, and transient feedback.

The repository-local V1.1 reconciliation:

- replaces the prior blue presentation identity with neutral-first Deep Teal + Soft Amber atmosphere;
- keeps the default atmosphere to one dominant teal field and one restrained amber counter-light field;
- preserves the existing 48px interaction-target floor;
- aligns current container geometry to V1.1 16 / 24 / 32 references while retaining capsule geometry where appropriate;
- provides explicit keyboard focus treatment;
- preserves Light and Dark system appearance and adds an explicit Deep Dark structural mode;
- removes nested backdrop blur when Glaze surfaces are nested rather than escalating effects;
- falls back to solid raised surfaces under Reduced Transparency;
- removes nonessential transform/motion under Reduced Motion;
- removes custom atmospheric fields and defers color handling to the platform under forced-colors;
- does not require Environmental Color Memory, geographic-content color sampling, or remote color derivation.

## Maps / Location authority boundary

GLAZE UI is presentation only, and Maps remains separate from GoreeCloud Location.

This migration does not request device location, read personal Location history, infer a visit, create a second tracking service, or claim authority over:

- current user/device location;
- personal location history;
- tracking controls;
- Find My;
- geofences;
- location-sharing permission state.

Those remain GoreeCloud Location authority. Maps remains authoritative for map presentation, geographic search/place metadata, routes, Saved Places, collections, and other Maps-owned state where implemented and accepted.

Atmospheric teal or amber never represents security, privacy, Identity authorization, Wardveil findings, Everkeep continuity, provider acceptance, geographic coverage, or route correctness by itself.

## Acceptance still required

Exact Maps revisions still require, as applicable:

- strict web type-check/build evidence;
- representative rendered map/chrome visual review and Human Visual Excellence acceptance;
- keyboard and screen-reader accessibility review, including non-map/list alternatives where required;
- 200% text and responsive reflow;
- RTL/localization;
- Reduced Motion and Reduced Transparency acceptance;
- contrast/high-contrast/forced-colors behavior;
- mobile/desktop/tablet/foldable/safe-area behavior;
- representative browser/GPU performance evidence;
- live approved geographic data/geocoder/router provider evidence;
- production Identity, Privacy Shield, Wardveil, Everkeep, Mesh and Location integrations where applicable;
- deployment, rollback, operational recovery, production signing/distribution where applicable, and production approval.
