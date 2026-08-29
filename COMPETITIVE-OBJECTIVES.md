# GoreeCloud Maps Competitive Objectives

## Objective

GoreeCloud Maps should meet the core usefulness expectations established by modern mapping products while differentiating through privacy, self-hostability, interoperability, user ownership, offline resilience, and GoreeCloud-native design.

Google Maps and Apple Maps are capability references, not source-code, data, artwork, wording, or visual-design sources. GoreeCloud must not scrape or copy proprietary map, place, review, imagery, navigation, or business data from them.

## Baseline capability objectives

### Exploration

Maps should provide smooth vector-map interaction, fast search, recognizable map styling, place selection, rich place detail, category discovery, favorites, recent context, globe/terrain/3D views where data supports them, and rapid transitions between browsing and directions.

### Places

Place cards should eventually support verified address/contact information, hours, categories, accessibility attributes, photos from licensed sources, user notes, saved state, provenance, related places, and useful actions. Place information must distinguish authoritative/provider data from user-contributed or inferred data.

### Directions

The route planner should compete on driving, walking, cycling, transit, and multimodal usefulness. It should support alternatives, departure/arrival planning, route avoidances, elevation context, accessibility constraints, and understandable route comparison.

### Navigation

The navigation experience should target clear turn-by-turn guidance, rerouting, lane/maneuver context, speed and incident information when supported, arrival context, low-distraction presentation, voice/haptic integration on capable platforms, and Live Glaze ongoing-navigation surfaces.

### Offline

Offline capability should be a first-class product area rather than a degraded afterthought. Regional downloads should preserve useful maps, saved content, search indexes, and routing data to the degree practical for the selected region and device.

### Collaboration

Shared collections and collaborative maps should be stronger than simple share links. GoreeCloud should provide explicit roles, membership, revocation, conflict handling, auditability, and durable ownership.

## Differentiation objectives

### Privacy-first architecture

- Minimize transfer of precise location/search/route data to third parties.
- Keep personal location history in GoreeCloud Location rather than duplicating it.
- Support private search/history modes and transparent provider boundaries.
- Avoid advertising profiles and hidden cross-context tracking.

### Self-hostability and portability

- Prefer open/self-hostable routing, geocoding, map data, and tile pipelines.
- Normalize providers behind GoreeCloud interfaces.
- Preserve user-created data independently from a map provider.
- Maintain export/restore paths through Everkeep.

### GoreeCloud interoperability

- Use Identity rather than a separate permanent account system.
- Use Location for sensitive device/location capabilities.
- Integrate geographic discovery with GoreeCloud Search.
- Use Mesh for governed cross-application coordination as contracts mature.

### Glaze UI 2.0 map-native experience

Maps should not visually imitate Google Maps or Apple Maps. The map canvas, floating search/control islands, connected place sheets, route panels, and Live Glaze navigation state should make the application recognizably GoreeCloud while preserving familiar task concepts.

## Quality objectives

- Fast initial shell and map interaction on supported hardware.
- Predictable map gesture behavior and keyboard/pointer equivalents.
- Clear fallback when GPU effects, live providers, location, or network access are unavailable.
- Multi-user isolation tests as a release blocker.
- No silent loss of user-created collections or annotations during provider migrations.
- Explicit stale/approximate/unverified state instead of false certainty.

## Competitive feature horizon

Future evaluation should include real-time traffic, crowdsourced incidents, transit vehicle positions, EV routing, indoor maps, street-level imagery, high-detail city data, 3D landmarks, AR pedestrian guidance, business-owner updates, community map corrections, map embeds, developer APIs, and privacy-preserving personalization.

Each capability must pass licensing, data-quality, privacy, security, accessibility, operational, and platform-readiness review before it is represented as an implemented GoreeCloud feature.