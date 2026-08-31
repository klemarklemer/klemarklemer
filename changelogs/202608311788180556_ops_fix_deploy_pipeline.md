# [ops] fix: Make the backend deploy pipeline actually run

**Date:** 2026-08-31  
**Timestamp:** 1788180556  
**Type:** `[ops]` bug fix  
**Service:** `apps/api`, `.github/workflows`

---

## Summary

The API had never deployed once. Four defects, each blocking the next.

1. **The workflow was rejected before it started.** Four steps guarded on
   `if: ${{ secrets.WIF_PROVIDER != '' }}`. `secrets` is not an available context
   in a step `if`, so every tag push produced a 0-second failure with no log.
   `actionlint` on the old file: `context "secrets" is not allowed here.
   available contexts are "env", "github", ...`. The secret is now surfaced as a
   job-level `env` and the guards test that instead.

2. **The image had never built.** Three faults, hit in sequence:
   - base image `golang:1.25.0-alpine3.21` against `go.mod` requiring `go 1.26.0`
   - `globalshared/` never copied into the build context, so
     `monorepo/globalshared` was unresolvable
   - `COPY .../.env`, a gitignored file absent from any clean checkout, which
     failed the final stage

3. **No runtime configuration reached the service.** The deploy set only
   `CORS_ALLOW_ORIGINS`; the service needs database and Redis DSNs and panics on
   an empty Redis DSN. It now receives the full candi app config, both DSNs from
   secrets, and the Vertex AI settings so assessment calls appear in Cloud Console.

4. **The default model was below the required floor.** `defaultModel` was
   `gemini-2.5-flash`, so a deployment without `GEMINI_MODEL` set would silently
   run 2.5. Raised to `gemini-3.5-flash` and set explicitly on the deploy.

## Verification

- `actionlint` reports no expression errors on the new workflow; the same tool
  reports the four `secrets`-context errors on the current `main`.
- `docker build --build-arg SERVICE_NAME=core -f apps/api/Dockerfile apps/api`
  succeeds for the first time.
- The built image runs with **environment variables only, no .env file**, and
  serves real rows from a local Postgres and Redis. The missing `.env` is a
  warning from candi's loader, not a failure.
- `go build` and `go test ./services/core/...` pass.

## Secrets this expects

`WIF_PROVIDER`, `WIF_SERVICE_ACCOUNT`, `GCP_PROJECT_ID` and `GCP_REGION` already
exist. Deployment additionally needs:

| Secret | Purpose |
|---|---|
| `SQL_DB_DSN` | Cloud SQL connection string, used for both read and write |
| `REDIS_DSN` | Redis connection string |
| `VPC_CONNECTOR` | optional; set to reach Memorystore on a private address |
| `CLOUDSQL_INSTANCE` | optional; set to attach the instance over a unix socket |

The two optional flags are omitted from the gcloud call when their secret is empty.
