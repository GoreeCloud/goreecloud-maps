# GoreeCloud Maps — GLAZE UI adoption

## Current state

GoreeCloud Maps is a **Forge-stage Development candidate**. The consolidated web source currently implements the repository-local **GLAZE UI V1.1 / 1.1.0** presentation mapping. The current approved shared consumer target is **GLAZE UI V1.6 / 1.6.0**, so Maps remains `applicable-migration-required`.

This record does not convert the existing V1.1 bytes into V1.6 by relabeling them. Shared Glaze acceptance is not Maps-specific rendered, accessibility, performance, browser/device, deployment, release, or production acceptance.

## Current shared authority

Live shared authority verified 2026-09-27 from `GoreeCloud/glaze-ui`:

- official product: **GLAZE UI V1.6**
- approved consumer target: `1.6.0`
- lifecycle: Anchor (the compatibility registry retains `status: stable`)
- tag: `v1.6.0`
- source qualification anchor: `c7509c79256b04b0aa67cb9dd0737d7588e0ae4a`
- accepted release source: `a7180679ea851389e0f3004515f9a25f420e716d`
- accepted release tree: `9ff0bf7a5f9d64f109d99bf4b76b81bd2a162268`
- shared web material entrypoint: `css/glaze-v1.4.1.css`
- shared runtime entrypoint: `js/glaze-v1.6.0.mjs`
- immediate known-good rollback runtime: `1.5.1`

V1.6 adds current shared runtime/capability behavior while preserving the accepted V1.4.1 web material baseline. Maps must evaluate applicable V1.5/V1.6 runtime behavior rather than treating a CSS filename as the complete design-system contract.

## Implemented Maps mapping

The current candidate retains its V1.1 source layer because it is real implemented source and historical evidence. It provides:

- neutral-first Deep Teal + Soft Amber product atmosphere;
- a 48 px interaction-target floor;
- explicit keyboard focus treatment;
- Light, Dark, and explicit Deep Dark structural modes;
- bounded frosted application chrome around the geographic canvas;
- no nested blur escalation;
- Reduced Transparency and Reduced Motion fallbacks;
- forced-colors fallback; and
- no environmental/geographic-content sampling as a hidden data-collection mechanism.

These are implemented source properties, not V1.6 conformance claims.

## V1.6 migration requirements

A substantive Maps migration must map the current shared Glaze contract without weakening map usability, attribution, privacy, or provider independence. At minimum it must:

1. keep the geographic renderer as the durable content canvas and apply Glaze to bounded application surfaces;
2. integrate applicable V1.6 runtime/capability behavior through a repository-controlled, auditable boundary;
3. keep provider, route, saved-place, collection, privacy, security, policy, and operational-health authority with their owning systems;
4. keep raw coordinates, private routes, search context, and provider credentials out of decorative/contextual presentation inputs unless a separately approved purpose explicitly authorizes them;
5. preserve Reduced Motion, Reduced Transparency, forced-colors, keyboard, screen-reader/list alternative, 200% text/reflow, safe-area, and responsive behavior; and
6. produce fresh Maps-specific render, accessibility, performance, browser/device/GPU, rollback, and Human Visual Excellence evidence.

## Maps / Location authority boundary

GLAZE UI is presentation only. Maps does not become authoritative for current user/device location, background tracking, personal location history, Find My, geofences, or location-sharing permission state. Those remain GoreeCloud Location responsibilities.

Maps remains responsible for map presentation and Maps-owned resources such as geographic search/place metadata, routing presentation, Saved Places, collections, and collaboration where implemented and accepted.

## Platform Contract boundary

The current candidate adopts **Platform Contract 2.0** with the canonical nine Integral Platform Systems. It is classified as `forge`, with `blocked` and `migration-required` flags, Development deployment state, and qualification not started.

GoreeCloud Policy and GoreeCloud Observability are explicit applicable-but-blocked systems. GoreeCloud Sync remains separately governed and is not a tenth Integral Platform System.

## Acceptance still required

Before Maps can reach Anchor, the exact candidate must satisfy every applicable build, API/database/RLS, rendered Glaze UI, accessibility, supported-platform, security, privacy, provider, Identity, Policy, Observability, continuity, deployment, rollback, release, and production acceptance gate. Passing source CI or the Platform Contract validator alone does not establish those outcomes.
