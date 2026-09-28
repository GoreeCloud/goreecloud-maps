# GoreeCloud Maps — User Manual

## Current scope

This manual describes the current Development stabilization candidate. It does not describe a released or production Maps service, and features that require an unconfigured provider or authentication authority may remain unavailable.

## Opening Maps

The web client opens to the map canvas with search/account controls, an Explore panel or mobile sheet, map controls, primary navigation, and service-state information.

A privacy-safe local map style is used when no approved map-data release or direct style is configured. This prevents a fresh development build from silently contacting an unrelated public tile provider.

## Search

When an approved geocoding provider is configured, use the search field to find places or addresses. Search results can be shown on the map. Provider unavailability is presented explicitly rather than as fabricated empty results.

## Directions

Directions require the routing capability reported by the Maps API. Supported candidate route modes include driving, walking, cycling, and transit/multimodal paths exposed by the configured adapter.

Route geometry and maneuver summaries are source-level Development features. They are not a claim that a production routing provider or navigation-quality acceptance exists.

## Saved Places

Signed-in users can work with owner-scoped Saved Places through the current candidate. Search results may expose an explicit Save action, and Saved Places can be listed or managed through the authenticated Maps boundary.

Saved Places are private Maps-owned state. Saving a provider-returned coordinate does not mean GoreeCloud Location observed or confirmed that the user visited the place.

## Shared collections

The current service foundation supports collection/member/item primitives with owner/editor/viewer authorization. The visual collaboration experience is incomplete. Human-friendly invitation lookup and delivery remain gated on approved Identity integration.

## Account state

Maps uses a GoreeCloud Identity-compatible Authorization Code + PKCE browser boundary when configured. No browser client secret is used. Current source keeps bearer access tokens in memory.

If Identity is not configured, protected features remain unavailable. Maps does not create a fallback account system.

## Current location

The location control is not an independent tracking permission path. Current-position display must use an approved GoreeCloud Location integration. Until that integration is accepted, Maps must not manufacture a location or read private Location history.

## Offline and degraded use

The current candidate provides a privacy-safe local renderer fallback and explicit provider status. Full offline-region download/search/routing behavior is not yet complete.

Unavailable, stale, approximate, offline, or degraded states should remain visibly distinct from current/verified state.

## Accessibility

The current source implements keyboard focus handling, practical control targets, Reduced Motion, Reduced Transparency, increased-contrast/forced-colors behavior, and a GLAZE UI V1.6 presentation mapping. Formal application-specific rendered, assistive-technology, 200% text/reflow, RTL/localization, browser/device/GPU, and performance acceptance remains outstanding.

## Privacy and safety

Maps does not own personal location tracking/history. Do not enter reusable credentials or secrets into place names, notes, search terms, or shared collection content. Private routes, Saved Places, notes, collaboration state, and precise personal coordinates should be treated as sensitive.

## Development limitations

Production geographic providers, production Identity/PostGIS, complete navigation, offline regions, traffic/transit/imagery layers, native clients, complete collaboration, platform-system acceptance, deployment, release, and Anchor qualification remain gated.
