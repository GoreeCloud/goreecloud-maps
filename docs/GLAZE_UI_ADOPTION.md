# GoreeCloud Maps — GLAZE UI adoption

## Current state

GoreeCloud Maps is a **Forge-stage Development candidate**. The web candidate now activates a repository-local **GLAZE UI V1.6 / 1.6.0** presentation adapter. Source migration to the current shared target is implemented; Maps-specific rendered, accessibility, adaptive, performance, browser/device/GPU, rollback, Human Visual Excellence, release, and production acceptance remain blocked.

This record does not treat a version marker or successful build as application acceptance.

## Current shared authority

Live shared authority verified 2026-09-27 from `GoreeCloud/glaze-ui`:

- official product: **GLAZE UI V1.6**
- approved consumer target: `1.6.0`
- lifecycle: Anchor
- tag: `v1.6.0`
- source qualification anchor: `c7509c79256b04b0aa67cb9dd0737d7588e0ae4a`
- accepted release source: `a7180679ea851389e0f3004515f9a25f420e716d`
- accepted release tree: `9ff0bf7a5f9d64f109d99bf4b76b81bd2a162268`
- shared web material baseline: `css/glaze-v1.4.1.css`
- shared runtime entrypoint: `js/glaze-v1.6.0.mjs`
- immediate known-good rollback runtime: `1.5.1`

V1.6 preserves the accepted V1.4.1 web material baseline while advancing current shared runtime and capability behavior.

## Implemented Maps adapter

The active candidate:

- loads `src/glaze-v1-6.ts` from the web entrypoint;
- marks the application surface with `data-glaze-version="1.6.0"`;
- pins the V1.6 qualification and accepted-release source anchors in the repository-local CSS mapping;
- preserves bounded frosted application chrome around the geographic canvas;
- exposes protected semantic surfaces for status and transient feedback;
- derives accessibility presentation profiles from reduced-motion, reduced-transparency, increased-contrast, and forced-colors preferences;
- adapts input-mode presentation across touch, pointer, and keyboard interaction;
- uses a bounded efficiency profile when the browser reports data-saving preference;
- preserves 48 px target sizing, visible focus, Light/Dark/Deep Dark structure, and no nested blur escalation; and
- keeps all Glaze state presentation-only and caller-truth-owned.

The adapter does not import authority over provider truth, route correctness, private location, privacy/security findings, policy decisions, or operational health.

## Maps / Location authority boundary

Maps does not become authoritative for current user/device location, background tracking, personal location history, Find My, geofences, or location-sharing permission state. Those remain GoreeCloud Location responsibilities.

Maps remains responsible for map presentation and Maps-owned resources such as geographic search/place metadata, routing presentation, Saved Places, collections, and collaboration where implemented and accepted.

## Sensitive context boundary

Raw coordinates, private routes, search context, provider credentials, Identity tokens, security findings, policy decisions, and operational evidence must not become decorative or adaptive presentation inputs unless an approved purpose explicitly authorizes the specific use. Shared Glaze capability does not create collection authority.

## Platform Contract boundary

The candidate adopts **Platform Contract 2.0** with exactly nine Integral Platform Systems. It is classified as `forge`, with `blocked` and `migration-required` flags, Development deployment state, and qualification not started.

The Glaze UI system result is `applicable-blocked`: the current V1.6 source mapping is implemented, but downstream acceptance remains incomplete. GoreeCloud Policy and GoreeCloud Observability are also explicit applicable-but-blocked systems. GoreeCloud Sync remains separately governed and is not a tenth Integral Platform System.

## Acceptance still required

Before Maps can reach Anchor, the exact candidate still requires applicable:

- strict source/build and API/PostGIS/RLS validation;
- rendered map/chrome review and Human Visual Excellence acceptance;
- keyboard and assistive-technology acceptance, including structured non-map alternatives;
- 200% text/reflow, RTL/localization, safe-area, and responsive acceptance;
- Reduced Motion, Reduced Transparency, increased-contrast, and forced-colors acceptance;
- representative browser/device/GPU and performance evidence;
- provider licensing, provenance, attribution, quality, degradation, privacy, and security acceptance;
- Identity, Location, Manager, Privacy Shield, Wardveil Security, Everkeep, Mesh, Policy, and Observability integration evidence;
- deployment, rollback, recovery, release provenance, and production approval.

Passing CI or the Platform Contract validator alone does not establish those outcomes.
