# GoreeCloud Maps Competitive Objectives

## Objective

GoreeCloud Maps should meet the core usefulness expectations of modern mapping products while differentiating through privacy, self-hostability, provider replaceability, user ownership, offline resilience, collaboration, accessibility, and GoreeCloud-native integration.

Google Maps and Apple Maps are capability and interaction references only. GoreeCloud must not scrape or copy proprietary map, place, review, imagery, routing, navigation, business, interface, or artwork content from them.

## Capability objectives

### Exploration and places

Maps should provide fast vector-map interaction, place/address/category discovery, rich place details, saved content, context-aware actions, and clear provenance when information comes from providers or user contributions.

### Directions and navigation

Route planning should ultimately cover driving, walking, cycling, transit, and multimodal travel with alternatives, useful constraints, understandable tradeoffs, route warnings, and high-quality turn-by-turn guidance where approved data and runtime capabilities exist.

### Offline and degraded operation

Offline regions should be a first-class product area. Maps should preserve useful map data, saved content, search indexes, and routing data to the degree practical for the chosen region/device, with explicit freshness, storage, update, and failure states.

### Collaboration

Shared collections and maps should use explicit roles, durable ownership, revocation, conflict handling, and auditability rather than relying only on ungoverned share links.

## Differentiation objectives

### Privacy-first architecture

- minimize transfer of precise location/search/route data to third parties;
- keep sensitive personal location-history authority in GoreeCloud Location;
- make provider boundaries and data-use behavior understandable;
- avoid advertising profiles and hidden cross-context tracking;
- keep public geographic-data delivery separate from private Maps/user state.

### Self-hostability and portability

- prefer open/self-hostable routing, geocoding, map-data, and tile capabilities where quality and licensing permit;
- normalize providers behind GoreeCloud interfaces;
- preserve user-created state independently from any one map provider;
- maintain export, backup, restore, and portability paths through Everkeep requirements.

### GoreeCloud interoperability

- use GoreeCloud Identity for authentication/principal identity without transferring Maps resource authorization;
- use GoreeCloud Location for sensitive current-position/tracking capabilities;
- interoperate with GoreeCloud Search for wider geographic discovery where contracts require it;
- use GoreeCloud Mesh for governed cross-application coordination where appropriate;
- consume Privacy Shield and Wardveil Security state substantively rather than as decorative branding.

### Original Glaze UI map experience

Maps should be recognizably GoreeCloud rather than a visual clone of Google Maps or Apple Maps. The map Canvas, search/control islands, place sheets, route panels, navigation state, adaptive layouts, motion, materials, accessibility, and resilience behavior should follow the accepted Maps Glaze UI consumer contract.

## Quality objectives

- responsive map interaction on supported hardware;
- predictable touch, pointer, keyboard, and assistive-technology behavior;
- clear structured alternatives to visual-only map information;
- truthful fallback when GPU effects, providers, location, authentication, or network access are unavailable;
- multi-user isolation and authorization tests as release blockers;
- no silent loss of user-created collections or saved places during provider migrations;
- explicit stale, approximate, unavailable, unverified, and degraded states instead of false certainty;
- provider and dataset licensing/provenance as first-class release evidence.

## Capability horizon

Future evaluation may include real-time traffic, crowdsourced incidents, transit vehicle positions, EV routing, terrain, indoor maps, street-level imagery, high-detail city data, 3D landmarks, AR pedestrian guidance, business-owner updates, community map corrections, map embeds, developer APIs, and privacy-preserving personalization.

Each capability must pass applicable licensing, data-quality, privacy, security, accessibility, operational, platform-integration, recovery, and release-readiness review before it is represented as an accepted GoreeCloud feature.

## Current evidence boundary

The executable foundation remains draft PR #1 at validated head `a8aafba65b89dbcd76661368a605762871abbb35`. Its successful CI is valuable candidate evidence but does not satisfy its outstanding independent-review gate, merge governance, production provider/database/Identity acceptance, deployment, release, or Stable qualification.
