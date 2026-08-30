# [api] fix: Make demo reset and officer assignment actually work

**Date:** 2026-08-30  
**Timestamp:** 1788060982  
**Type:** `[api]` bug fix  
**Service:** `apps/api/services/core`

---

## Summary

Verifying the web↔API connection (`[web]` connect change) surfaced real backend
bugs that prevented the demo flow from running at all.

1. **Migration duplicate versions (panic).** `cmd/migration/migrations/` had two
   candi-scaffold *stub* files (`..._claims.sql` / `..._policies.sql` at versions
   `...18` and `...20`) that collide with the real DDL at the same versions.
   goose's sorter panics on duplicate versions. Deleted the two stub files; the
   real schema is created by GORM `AutoMigrate` (the `GetMigrateTables` path) and
   the remaining `policies`/`officers`/`claims`/`seed` migrations.
2. **`ResetDemo` broken.** It deleted only `id=1` then `Save`d a struct with
   `ID:1` preset. gorm `Save` with a non-zero PK issues an `UPDATE`, which after
   the delete matches 0 rows → `record not found`. Also the next reset re-inserted
   the same `claim_number` (`UNIQUE` violation). Fixed by dropping the preset PK
   (so `Save` inserts) and deleting the existing demo Claim by `claim_number`
   before re-seeding; child rows now reference `claim.ID` instead of hardcoded `1`.
3. **`FetchAll` returns 0 rows.** Both claim and officer repos applied
   `LIMIT 0` for the default empty filter; GORM v2 treats `LIMIT 0` as "no rows",
   so `RunAssignment` found 0 officers (`no eligible claims officers found`).
   Changed the guard to `if filter.Limit > 0` so internal calls (Limit 0) return
   all rows while query-param pagination still applies.

## Verification

API now boots against local Postgres/Redis and the full demo flow succeeds via curl:
`POST /v1/demo/reset` → `POST /v1/claim/:id/documents` (upload + intake + assignment)
→ `POST /v1/claim/:id/assessment` (recommendation) →
`POST /v1/claim/:id/approval` (APPROVE → stage CLOSED, approved_amount 3700).
Existing claim usecase tests still pass.

## Notes

- Local Postgres on `:5432` lacked the `klemarklemer_db` DB/role; created
  `user`/`klemarklemer_db` to match `.env.sample` DSNs. Redis was started via
  `deployments/compose/docker-compose.dev.yml`.
