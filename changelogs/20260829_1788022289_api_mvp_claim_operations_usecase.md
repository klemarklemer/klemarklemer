# [api] core: MVP Claim Operations Usecase & REST Handlers

**Date:** 2026-08-29  
**Timestamp:** 1788022289  
**Type:** `[api]` backend  
**Service:** `apps/api/services/core`

---

## Summary

Completed the full backend MVP for Claim Operations. All usecase business logic, REST handlers, domain fixes, and Candi stub cleanup are now in place. The service compiles cleanly and all tests pass.

---

## Files Changed

### New Usecase Files

| File | Description |
|---|---|
| `claim/usecase/usecase.go` | Expanded `ClaimUsecase` interface with 6 new methods |
| `claim/usecase/get_all_claim.go` | `GetAllClaim` with `candishared.NewMeta` pagination |
| `claim/usecase/get_detail_claim.go` | `GetDetailClaim` full preloaded claim |
| `claim/usecase/create_claim.go` | `CreateClaim`, `UpdateClaim`, `DeleteClaim` |
| `claim/usecase/upload_document.go` | `UploadDocument` auto-triggers EvaluateIntake |
| `claim/usecase/evaluate_intake.go` | `EvaluateIntake` completeness gate + classification |
| `claim/usecase/run_assignment.go` | `RunAssignment` deterministic officer scoring |
| `claim/usecase/run_assessment.go` | `RunAssessment` generates AssessmentRecommendation |
| `claim/usecase/submit_human_approval.go` | `SubmitHumanApproval` records Decision, closes Claim |
| `claim/usecase/reset_demo.go` | `ResetDemo` restores CLM-2026-0042 demo state |

### Updated Delivery

| File | Description |
|---|---|
| `claim/delivery/resthandler/resthandler.go` | Full REST handler with 9 endpoints |

### Domain Fixes

| File | Description |
|---|---|
| `policy/domain/request.go` | Real Policy fields |
| `policy/domain/response.go` | Real Policy fields |
| `officer/domain/request.go` | Real Officer fields |
| `officer/domain/response.go` | Real Officer fields |
| `officer/usecase/update_officer.go` | Real fields, not placeholder |
| `policy/usecase/update_policy.go` | Real fields, not placeholder |

---

## Business Logic

1. **Intake Agent**: POLICE_REPORT + DAMAGE_PHOTO required for COMPLETE
2. **Assignment Agent**: `(10 - workload) * 0.5 + (skill/5) * 10 * 0.5`
3. **Assessment Agent**: APPROVE outcome, 0.94 confidence, structured reasons
4. **Human approval**: `approved_amount = estimated_loss - deductible`, Stage = CLOSED
5. **Demo reset**: Single POST restores demo to incomplete state

---

## Verification

```
go build ./services/core/...  → exit 0
go test ./services/core/...   → all ok
```
