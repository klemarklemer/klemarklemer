# Changelog: 2026-08-29 - Tagging and Release Rules

## Summary
Defined and implemented the semantic versioning and Git tagging policy (`{{env}}-{{apps}}-{{service name}}-vA.B.C`) and updated the GitHub Actions CD deployment pipeline to trigger on tag pushes.

---

## [ops]
- Created `docs/ops/release-and-tagging.md` documenting the tagging syntax:
  - `{{env}}-{{apps}}-{{service name}}-vA.B.C` (e.g. `production-api-core-v1.0.0`, `staging-web-v1.0.0`).
  - Version increments: `A` = Huge change/refactor, `B` = Major change/feature, `C` = Minor change/fix.
- Updated `AGENTS.md` with tagging and changelog constraints.
- Updated `.github/workflows/cd.yml` to trigger on tags matching `*-api-core-v*` and `*-web-v*`.
