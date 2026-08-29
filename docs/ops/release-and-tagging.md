# Release and Git Tagging Policy

This document outlines the strict semantic versioning, tagging format, and deployment triggers for the repository.

---

## 1. Tagging Format

Every deployment and release tag must strictly follow the syntax:

```
{{env}}-{{app}}-{{service}}-vA.B.C
```

Or for single-service applications (like `apps/web`):
```
{{env}}-{{app}}-vA.B.C
```

### Components:
- **`{{env}}`**: Target environment. One of:
  - `production`
  - `staging`
  - `development`
- **`{{app}}`**: Application directory (`api` or `web`).
- **`{{service}}`**: Service name for backend monorepo (e.g. `core`).
- **`v`**: Literal character prefix for the version number.
- **`A.B.C`**: Semantic version numbers where:
  - **`A`** (Huge Change): Monumental redesign, breaking architecture changes, or complete system refactors.
  - **`B`** (Major Change): Significant new features, new agent specialist workflows, or major domain additions.
  - **`C`** (Minor Change): Small enhancements, UI polish, bug fixes, or performance patches.

---

## 2. Tag Examples

| Target | Environment | Git Tag Example |
|---|---|---|
| Backend `core` service | Development | `development-api-core-v0.1.0` |
| Backend `core` service | Staging | `staging-api-core-v1.0.0` |
| Backend `core` service | Production | `production-api-core-v1.0.0` |
| Frontend `web` console | Staging | `staging-web-v1.0.0` |
| Frontend `web` console | Production | `production-web-v1.0.0` |

---

## 3. Automated CD Trigger Rules

Pushing tags to GitHub triggers automated deployments in [`.github/workflows/cd.yml`](../../.github/workflows/cd.yml):

- **`*-api-core-v*`**: Triggers Docker container build for `apps/api` (`service=core`) and deploys to the corresponding GCP Cloud Run environment (`production`, `staging`, or `development`).
- **`*-web-v*`**: Triggers production build for `apps/web` and deploys to the corresponding frontend hosting target.

---

## 4. Tag Creation Workflow

```bash
# Example: Tagging staging release for API core service
git tag -a staging-api-core-v1.0.0 -m "Release staging-api-core-v1.0.0: Seeded claim ops and REST delivery"
git push origin staging-api-core-v1.0.0
```
