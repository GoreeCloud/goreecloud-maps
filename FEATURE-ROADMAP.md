# GoreeCloud Maps — Feature Roadmap

**Status:** Active roadmap control  
**As of:** 2026-09-27  
**Authoritative project record:** Project Specification — Maps  
**Canonical repository:** GoreeCloud/maps  
**Drive control:** `GoreeCloud/Feature Roadmap/GoreeCloud Maps/FEATURE-ROADMAP.docx`

## Purpose

This file is the repository-side feature roadmap control for GoreeCloud Maps. It records current planned and recommended feature work without replacing the authoritative project record, exact implementation evidence, release gates, or GoreeCloud Tasks Management.

## Roadmap

| ID | Feature / obligation | Priority | Current state |
| --- | --- | --- | --- |
| FR-001 | Reconcile current and recommended Maps scope against the authoritative project record and verified repository evidence. | High | Ongoing control |
| FR-002 | Move unfinished actionable obligations into GoreeCloud Tasks Management with priority, dependency, and lifecycle disposition. | High | Ongoing control |
| FR-003 | Do not mark features implemented, complete, cancelled, superseded, Anchor, or production-ready without authoritative evidence. | High | Ongoing control |
| FR-004 | Consolidate the executable MapLibre/PostGIS/API foundation, Saved Places UI, provider-backed search/save, routing, and collaboration work into one current reviewable candidate rather than maintaining a stale stacked-PR chain. | Critical | In progress — stabilization candidate |
| FR-005 | Complete fresh Maps-specific rendered/accessibility/performance/form-factor acceptance for the implemented GLAZE UI V1.6 / 1.6.0 source mapping. | High | Source implemented; acceptance blocked |
| FR-006 | Complete Platform Contract 2.0 nine-system integration work for Manager, Privacy Shield, Wardveil, Everkeep, Glaze UI, Mesh, Identity, Policy, and Observability. | Critical | Blocked / in progress |
| FR-007 | Approve and validate live geographic data, geocoding, and routing providers with licensing, provenance, attribution, quality, egress, privacy, security, and degradation controls. | Critical | Planned / blocked |
| FR-008 | Integrate approved current-position behavior through GoreeCloud Location without creating a second tracking/history authority. | High | Planned / blocked |
| FR-009 | Complete offline-region management, provider degradation/fallback, storage/freshness controls, accessibility, representative-device/browser/GPU validation, deployment, rollback, and Anchor qualification. | High | Planned |

## Current authority boundary

GoreeCloud Maps owns mapping presentation and Maps-owned resources such as geographic search/place metadata, routing presentation, Saved Places, collections, and collaboration where implemented. GoreeCloud Location remains authoritative for current user/device position, background tracking, personal location history, Find My, geofences, and location-sharing permission state.

## Maintenance and synchronization

This roadmap and the corresponding Drive `FEATURE-ROADMAP.docx` must remain materially synchronized with the authoritative project record, live GitHub state, Platform Contract, and GoreeCloud Tasks Management. Missing obligations, stale repository names, stale platform-system counts, obsolete Glaze targets, duplicated work, or undocumented status changes are defects.

Completion and lifecycle claims require the applicable implementation, exact-head validation, review, runtime evidence, release evidence, and production acceptance. A green CI run alone does not establish Anchor.
