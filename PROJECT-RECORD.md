# GoreeCloud Maps — Project Record

This record preserves significant GoreeCloud Maps history and evidence. It does not replace `CHANGELOGS.md`.

## 2026-09-11 — Drive project specification established

The frozen migration source `Project Specification — Maps.docx` was created in Google Drive. It defined the intended first-party mapping product, Maps/Location authority split, MapLibre/PostGIS direction, provider replaceability, offline/navigation/collaboration scope, privacy/security expectations, and production gates. Its historical repository name and Glaze/lifecycle terminology are superseded by current repository authority.

## 2026-08-29 — Executable foundation validation

Draft foundation PR #1 established the TypeScript/Vite/MapLibre web shell, Go API, PostgreSQL/PostGIS multi-user foundation, Identity-compatible PKCE boundary, provider adapters, public map-data release contracts, and Saved Places service/client work.

Exact head `a8aafba65b89dbcd76661368a605762871abbb35` passed CI run `33251905892`. The PR remained draft and unmerged because its independent review gate was not satisfied.

## 2026-09-01 — Accepted repository-governance baseline

Governance-only PR #4 established the accepted-main README/specification/feature/benefit/competitive/branding/license baseline. The accepted merge did not promote executable PR #1 behavior to main.

## 2026-09-27 — Repository-native feature/changelog migration

Repository-native `IMPLEMENTED-FEATURES.md`, `PLANNED-FEATURES.md`, and `CHANGELOGS.md` were introduced. The legacy Drive roadmap was recorded as retired after parity verification. The later repository migration removes the obsolete `FEATURE-ROADMAP.md` compatibility copy after its current obligations are accounted for.

## 2026-09-27 — Clean current-main stabilization candidate

A first consolidation attempt diverged from newer mainline governance and was closed as superseded. Draft PR #11, branch `stabilization/maps-current-main-20260927`, was rebuilt from current `main`, preserving newer mainline records while layering the executable foundation and later Saved Places/search-save work.

The branch was explicitly synchronized with current `main` using a merge commit rather than pretending content equivalence was ancestry.

## 2026-09-27 — Platform Contract 2.0 reconciliation

The Maps declaration was migrated to Platform Contract 2.0 with lifecycle `forge`, Development deployment state, and exactly nine Integral Platform Systems. GoreeCloud Sync remained separately governed. Policy and Observability were added as explicit platform-system evaluation domains without manufacturing implementation evidence.

## 2026-09-27 — GLAZE UI V1.6 source migration

Maps replaced active repository-local V1.1 activation with a repository-local GLAZE UI V1.6 / 1.6.0 adapter.

The adapter is presentation-only, preserves caller-owned provider/application truth, adds protected semantic surfaces, accessibility-profile fallbacks, input-mode adaptation, Reduced Motion/Transparency behavior, forced-colors handling, and bounded performance degradation, and avoids a new remote browser runtime dependency.

Pre-governance-migration exact head `c3d0183a1db55dc7ea72fb753c97f28e13f561e3` passed:

- Maps CI run `36369007184`;
- Repository Governance run `36369007204`;
- Platform Contract run `36369007699`.

Those runs are historical evidence for that exact head only. Later documentation/governance commits require fresh exact-head validation.

## Current open acceptance boundary

Maps remains Forge/Development. Open gates include independent review/protected promotion of the consolidated candidate, rendered/accessibility/performance/form-factor Glaze acceptance, substantive nine-system integrations, approved live geographic providers, GoreeCloud Location current-position integration, offline-region completeness, production database/Identity acceptance, continuity/restore, deployment, rollback, release provenance, and Anchor qualification.
