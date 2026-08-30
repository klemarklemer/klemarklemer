# [web] feat: Surface Gemini assessment evidence in the Claim workspace

**Date:** 2026-08-30  
**Timestamp:** 1788062982  
**Type:** `[web]` feature  
**Service:** `apps/web`

---

## Summary

The demo was not presentable: it resolved in one magic click and never showed
*why* the AI recommended what it did. Per the plan, the differentiator is the
agentic document understanding + explainable recommendation, not CRUD.

Added an **Assessment evidence** disclosure inside the recommendation block that
makes Gemini's work visible:

1. **Document analysis** — for each present document, parses the backend
   `extracted_data` JSON (detected damage, severity, report ID, reporting
   officer, incident date/location, liability) into labelled fields. This is the
   "AI read the police report and damage photo" story.
2. **Policy check** — shows the policy is in force, coverage type/limit, and
   deductible, so the officer sees the basis for the payout math.

Supporting changes:
- `types.ts`: added `PolicySummary` and `extractedData` on `DocumentItem`,
  `policy`/`fraudSignal` on `Claim`.
- `api/claims.ts`: adapter now maps `policy` and `document.extracted_data`
  from the backend response (both already returned by the API).
- `components/AssessmentEvidence.tsx`: new panel, embedded in
  `RecommendationPanel`.

SLA at-risk is already demonstrated on Seed (resetDemo sets `stage_sla_due_at`
+25m → `atRisk`), so the SLA clock story lands without extra work.

## Verification

`tsc` clean, `vitest` 3/3 pass, `vite build` succeeds. Live API confirmed to
return `policy` (POL-MOTOR-2026-8819) and `extracted_data` on both seeded
documents, so the panel populates on the real flow.

## Next (from the agreed 2/3/4 scope)

- **(2) Fraud scenario:** add a demo claim flagged for investigation with a
  mandatory human-review recommendation (needs a `fraud_signal` field from the
  backend + a second seed).
- **(4) Light triage strip:** a compact "claims needing attention" list
  (stage / SLA / fraud flag), kept inside the workspace per the UX contract
  (no management KPI command centre).
