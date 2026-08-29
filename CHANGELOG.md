# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and commits/entries follow the prefixes:
- `[api]` for backend changes
- `[web]` for frontend changes

---

## [Unreleased]

### Added
- `[api]` Initialized Candi monorepo codebase in `apps/api` with service `core`.
- `[api]` Generated `claim`, `officer`, and `policy` modules in `services/core` with REST delivery (go-chi) and GORM repository abstractions.
- `[api]` Configured PostgreSQL (`klemarklemer_db`) and Redis DSNs in `.env` and `.env.sample`.
- `[api]` Added `deployments/compose/docker-compose.dev.yml` for local PostgreSQL 16 and Redis 7 containers.
- `[api]` Created `workflows/claim-ops-loop.md` documenting the autonomous agent loops and mandatory human approval checkpoint.
- `[api]` Documented ADR 0004 for PostgreSQL + GORM + Candi state store architecture (superseding ADR 0002).
- `[web]` Moved React/TS claims operations console into `apps/web` within the monorepo structure.
- `[api]` Created root `Makefile` coordinating `up`, `dev`, `deps`, `api`, `web`, `migrate`, and `test` targets across `apps/api` and `apps/web`.
