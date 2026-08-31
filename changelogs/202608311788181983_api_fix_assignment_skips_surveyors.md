# [api] fix: Stop Loop 2 assigning claims to surveyors

**Date:** 2026-08-31  
**Timestamp:** 1788181983  
**Type:** `[api]` bug fix  
**Service:** `apps/api/services/core`

---

## Summary

`RunAssignment` scored every row `OfficerRepo().FetchAll` returned. Since the
survey change put three surveyors in `claims_officers`, they became candidates for
owning a claim - and on the seeded roster every surveyor outscores every claims
officer, so a MOTOR claim was always assigned to an inspector. Observed live:
Priya Sharma, role Surveyor, specialty Property Inspection, score 9.50.

The scoring is now `selectClaimOwner`, a pure function that skips surveyors and
unavailable officers. Extracting it makes the rule testable without a database.

The old `bestOfficer == nil` fallback picked `officers[0]` with a fabricated score
of 5.0, which would have reintroduced the same defect by selecting an unavailable
officer or a surveyor. Assignment now fails loudly when nobody is eligible, matching
the guard directly above it.

`RoleSurveyor` and `IsSurveyor()` live on the officer domain, replacing the string
literal that `survey.go` was carrying separately.

## Verification

- Five tests pin the rule, including the seeded roster. With the surveyor filter
  removed they fail on "Loop 2 assigned a surveyor to own a claim"; with it they pass.
- `go build` and `go test ./services/core/...` pass.

## Also in this change

`docs/SUBMISSION.md` claimed the console "runs on local fixture data and makes no
HTTP calls". That stopped being true when the frontend was wired to the API. The
limitations section now describes what is actually unfinished: Loop 1's constant
classification fields, and survey handling being reachable only through seeded data.
