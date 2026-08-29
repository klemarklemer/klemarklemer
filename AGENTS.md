# Agent notes

**UI/UX:** When building or changing the claims officer console, follow [docs/ui/claims-ops-ux.md](docs/ui/claims-ops-ux.md) (ops console, not a landing page). Domain words: [CONTEXT.md](CONTEXT.md).

**Changelogs & Commits:**
- Every change must be documented in `changelogs/YYYYMMDD{{timeunix}}_title.md` (e.g. `202608291788021435_title.md`).
- Prefix changes and commits with `[api]` for backend and `[web]` for frontend.

**Tagging & Releases:**
- Follow `{{env}}-{{apps}}-{{service name}}-vA.B.C` (e.g. `production-api-core-v1.0.0`, `staging-web-v1.0.0`). Details in [docs/ops/release-and-tagging.md](docs/ops/release-and-tagging.md).
- Versioning: `A` = Huge change/refactor, `B` = Major change/feature, `C` = Minor change/fix.

