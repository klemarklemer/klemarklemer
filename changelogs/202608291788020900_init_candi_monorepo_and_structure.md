# Changelog: 2026-08-29 - Init Candi Monorepo and Structure

## Summary
Restructured the repository into a monorepo layout containing `apps/web` (React/TS console) and `apps/api` (Golang Candi monorepo with `core` service), added PostgreSQL & Redis docker compose configuration, and formalized workflow specifications and ADR 0004.

---

## [api]
- Initialized Candi monorepo codebase in `apps/api` (`monorepo`).
- Scaffolded service `core` with modules `claim`, `officer`, and `policy`.
- Configured REST delivery via `go-chi` and GORM database repository layer.
- Aligned database connection to PostgreSQL `klemarklemer_db` and Redis in `.env` and `.env.sample`.
- Created `deployments/compose/docker-compose.dev.yml` for local PostgreSQL 16 and Redis 7 containers.
- Created `workflows/claim-ops-loop.md` documenting the autonomous agent loops and mandatory human approval checkpoint.
- Published ADR 0004 for PostgreSQL + GORM + Candi state store architecture (superseding ADR 0002).
- Verified `go build ./services/core/...` compiles cleanly.

---

## [web]
- Moved React/TS claims operations console into `apps/web` within the monorepo structure.
- Preserved all UI components (`ClaimIdentity`, `DocumentCompleteness`, `AssignmentPanel`, `RecommendationPanel`, `TimelinePanel`, `NotificationsPopover`, `Toast`).

---

## [ops]
- Created root `Makefile` coordinating `up`, `dev`, `deps`, `api`, `web`, `migrate`, and `test` targets.
