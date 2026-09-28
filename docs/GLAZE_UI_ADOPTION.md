# GoreeCloud Maps — GLAZE UI adoption

## Current state

GoreeCloud Maps is a **Forge-stage Development candidate**. The consolidated web source now implements a repository-local **GLAZE UI V1.6 / 1.6.0** presentation adapter against the current approved shared target. Glaze is now `applicable-blocked`: source migration is present, but Maps-specific rendered and runtime acceptance is not complete.

The V1.6 adapter is a substantive repository-local source change, not a version-label-only rewrite. It avoids a remote browser runtime dependency, preserves caller-owned provider/application truth, and adds V1.6 semantic-surface, accessibility-profile, input-mode, motion, and constrained-performance behavior. Shared Glaze acceptance still does not equal Maps-specific rendered, accessibility, performance, browser/device, deployment, release, or production acceptance.

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

## Implemented Maps V1.6 mapping

The active candidate now uses `apps/web/src/glaze-v1-6.ts` and `apps/web/src/glaze-v1-6.css`. The older V1.1 source remains inactive historical provenance. The active V1.6 mapping provides:

- neutral-first Deep Teal + Soft Amber product atmosphere;
- a 48 px interaction-target floor;
- explicit keyboard focus treatment;
- Light, Dark, and explicit Deep Dark structural modes;
- bounded frosted application chrome around the geographic canvas;
- no nested blur escalation;
- Reduced Transparency and Reduced Motion fallbacks;
- forced-colors fallback; and
- no environmental/geographic-content sampling as a hidden data-collection mechanism.

These are implemented source properties. They establish a current-target source mapping, not application-specific rendered or production conformance.

## V1.6 acceptance requirements

The source migration now maps the current shared Glaze contract without weakening map usability, attribution, privacy, or provider independence. Remaining acceptance must verify that it:

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

The current candidate adopts **Platform Contract 2.0** with the canonical nine Integral Platform Systems. It is classified as `forge`, with a `blocked` flag, Development deployment state, and qualification not started. Glaze source is current-target V1.6, while downstream acceptance remains blocked.

GoreeCloud Policy and GoreeCloud Observability are explicit applicable-but-blocked systems. GoreeCloud Sync remains separately governed and is not a tenth Integral Platform System.

## Acceptance still required

Before Maps can reach Anchor, the exact candidate must satisfy every applicable build, API/database/RLS, rendered Glaze UI, accessibility, supported-platform, security, privacy, provider, Identity, Policy, Observability, continuity, deployment, rollback, release, and production acceptance gate. Passing source CI or the Platform Contract validator alone does not establish those outcomes.
