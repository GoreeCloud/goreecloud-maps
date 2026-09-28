# GoreeCloud Maps — Notes

## Current working state

- Lifecycle: Forge.
- Deployment state: Development.
- Platform Contract: 2.0.
- Current stabilization candidate: draft PR #11, `stabilization/maps-current-main-20260927`.
- Active Glaze source mapping: GLAZE UI V1.6 / 1.6.0.
- Maps-specific rendered/accessibility/performance acceptance: blocked/pending.
- Live geographic providers: not production accepted.
- GoreeCloud Location current-position integration: pending.
- Anchor qualification: not started.

## Authority notes

`PROJECT-SPECIFICATIONS.md` is the canonical project specification once accepted on `main`. `PROJECT-RECORD.md` preserves significant history. `SPECIFICATIONS.md` remains an implementation-focused companion and must not compete with the project specification.

Feature authority is repository-native through `IMPLEMENTED-FEATURES.md` and `PLANNED-FEATURES.md`. The legacy `FEATURE-ROADMAP.md` is retired after migration verification.

GitHub owns pull-request, review, workflow, branch, commit, and merge state. Google Drive task records track unfinished work but do not mirror PR state.

## Development notes

Maps must keep private user state separate from the public geographic-data release plane and must not absorb GoreeCloud Location tracking/history authority.

A green build, provider-compatible interface, platform badge, or merge does not by itself establish deployment, production, release, or Anchor state.
