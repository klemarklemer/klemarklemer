# [api] fix: Make a clean database migrate and boot

**Date:** 2026-08-31  
**Timestamp:** 1788161931  
**Type:** `[api]` bug fix  
**Service:** `apps/api/services/core`

---

## Summary

A clone of `main` could not reach a working state. Two independent defects, both
in the migration set.

1. **`fraud_signal` had no migration.** `shared/domain.Claim` declares
   `fraud_signal` (added with the REJECT_FRAUD work) and no DDL ever creates the
   column, so GORM put it in every INSERT and UPDATE against `claims`. On a
   correctly migrated database both `POST /v1/claim` and `POST /v1/demo/reset`
   returned 400 `pq: column "fraud_signal" of relation "claims" does not exist`.
   Added `20260831061500_add_fraud_signal_to_claims.sql`.

2. **The MVP seed referenced columns added by a later migration.** The survey
   change edited `20260829232721_seed_mvp_data.sql` in place to insert
   `specialty`, `region`, `surveyor_id`, `survey_status` and friends, but those
   columns are created by `20260830235216_add_survey_and_surveyor.sql`, which
   goose runs *after* it. On a fresh database the seed aborted at
   `column "specialty" of relation "claims_officers" does not exist`, leaving
   1 policy, 3 officers and no claims, documents or events at all.

   Fixed forward rather than by renumbering, so already-migrated databases are
   unaffected: the seed is restored to the shape that is valid at its own
   version, and the survey demo rows move to a new
   `20260831061600_seed_survey_demo_claim.sql` that runs after the columns exist.

## Verification

Against a clean PostgreSQL 16:

- `go run cmd/migration/migration.go up` applies all seven migrations, ending at
  `20260831061600`. Seeded state is 1 policy, 6 officers, 2 claims, 3 documents,
  6 events — including the survey-required claim `CLM-2026-0044`. On `main` the
  same run yields 1 / 3 / 0 / 0 / 0.
- `down`, `down`, `up` across the two new migrations is clean.
- Service boots and the two previously failing calls succeed:
  `POST /v1/claim` → 201, `POST /v1/demo/reset` → 200.
- `go build ./services/core/...` and `go test ./services/core/...` pass.

## Notes

Out of scope, found while verifying and worth separate fixes:

- `RunAssignment` scores every row in `claims_officers`, which since the survey
  change includes the three surveyors. They outrank all three claims officers, so
  a MOTOR claim is assigned to a surveyor (observed: Priya Sharma, Surveyor,
  Property Inspection).
- Nothing ever sets `survey_required = true`, and `ResetDemo` seeds both its
  claims with it false, so the survey feature is unreachable once the demo is
  reset. The frontend also calls none of the three survey endpoints.
- `CreateClaim` builds `claim_number` from `time.Now().Unix() % 10000`; two
  claims created in the same second violate the unique constraint.
