# [api] feat: Give the Intake Agent real reasoning, and revive the survey loop

**Date:** 2026-08-31  
**Timestamp:** 1788183937  
**Type:** `[api]` feature  
**Service:** `apps/api/services/core`, `apps/web`

---

## Summary

Loop 1 was named `IntakeAgent` but reasoned about nothing. It wrote
`severity = "MEDIUM"` and `survey_required = false` for every claim, and reported
`missing_documents: ["POLICE_REPORT"]` as a constant - so a claim missing only the
photo named the wrong document.

`gemini.Classifier` now backs the loop, mirroring `gemini.Assessor`: same credential
resolution, same deterministic fallback, same "which engine answered" reporting. It
judges severity and whether a surveyor must inspect. Which required documents are
absent stays set arithmetic computed in code, so the model cannot name a document
that is on file.

That flag was also the only reason survey handling was dead. Nothing ever set
`survey_required = true`, so `CompleteSurvey` always answered "survey not required"
and the console panel never rendered. With Loop 1 able to call for an inspection, the
whole path is reachable, and `routeToSurvey` hands the claim to a free surveyor whose
specialty fits the claim type - a property inspector on a car crash is the same class
of mistake as a surveyor owning a claim.

Two supporting gaps closed along the way: `ResponseClaim` carried `survey_required`
but none of the surveyor, status, dates or report fields the console reads, and the
repository never preloaded the `Surveyor` relation. The panel would have rendered
blank even once reachable.

## Verification

End to end against a clean database, on deterministic rules:

- Small loss (4,200 of 45,000 cover): `LOW`, no survey, assigned to Elena Rostova,
  Claims Officer. The existing demo path is unchanged apart from a reasoned severity.
- Large loss (30,000 of 45,000): `HIGH`, survey required, stage `SURVEY`, assigned to
  Marcus Webb - Vehicle Inspection, not the less-loaded property inspector.
- `POST /survey/complete` returns 200 and chains to assessment. That call previously
  always failed.
- A claim holding only a damage photo reports `missing_documents: ["POLICE_REPORT"]`;
  the reverse case is pinned by test.

Nine new tests across the classifier and surveyor selection. `go build`,
`go test ./services/core/...`, `npm run typecheck` and `npm run test` all pass.

## Documentation

`docs/adr/0001` described a Supervisor plus four ADK specialists including an SLA
specialist. None of that was built. It now records the architecture that exists, with
a note explaining the correction. The README loop table and the submission's known
limitations are updated to match.
