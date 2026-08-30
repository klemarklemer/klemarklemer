# [web] feat: Fraud scenario + Triage strip in the Claim workspace

**Date:** 2026-08-30  
**Timestamp:** 1788064982  
**Type:** `[web]` feature  
**Service:** `apps/web`, `apps/api`

---

## Summary

Two linked slices that make the demo presentable and demonstrate the plan's
agentic orchestration story:

1. **Fraud scenario** — `POST /v1/demo/reset` now seeds **two** demo claims:
   - **CLM-2026-0042** (normal MOTOR, incomplete, missing police report) — the
     original flow the officer works through.
   - **CLM-2026-0043** (fraud-flagged MOTOR, COMPLETE, stage DECISION) — policy
     holder name mismatch, duplicate vehicle claim history, estimated loss at
     policy max. Recommendation = `MANUAL_REVIEW` (not APPROVE), fraud_signal
     text populated. Approve is disabled; Reject is labeled "Reject Claim
     (Fraud Override)". A red fraud banner in the Claim identity shows the
     signal text.

2. **Triage strip** — a compact horizontal list above the workspace showing
   all demo claims with:
   - Claim number, line, stage
   - SLA at-risk amber dot
   - Red "Fraud" badge when `fraud_signal` present
   - Click to load any claim into the workspace (GET `/v1/claim/:id`)

   Lives inside the one-screen workspace (no command centre), per UX contract.

3. **API additions** (`apps/api`):
   - `fraud_signal` column on `claims` table (nullable `varchar(512)`), added
     via GORM AutoMigrate.
   - `ResetDemo` upserts the fraud claim alongside the normal one, with
     `MANUAL_REVIEW` recommendation, two documents (claim form + police
     report), and a `FRAUD_SIGNAL_DETECTED` event.

4. **Frontend additions** (`apps/web`):
   - `api/client.ts`: `apiGet` helper.
   - `api/claims.ts`: `listClaims`, `getClaim`, `toTriageItem`, `fraudSignal`
     on `Claim`, `TriageItem` type.
   - `components/TriageStrip.tsx`: clickable claim cards.
   - `components/ClaimIdentity.tsx`: fraud banner (red, `Warning` icon).
   - `RecommendationPanel.tsx`: handles `MANUAL_REVIEW` — shows mandatory-review
     notice, disables Approve, keeps Reject with "Fraud Override" label.
   - `App.tsx`: loads triage on mount/seed, click handler loads claim detail.

## Verification

- Backend: `go build` + API restart → `POST /demo/reset` returns both claims,
  `GET /v1/claim` lists 2 (one fraud), `GET /v1/claim/:id` returns full
  fraud_signal + MANUAL_REVIEW recommendation + documents + policy.
- Frontend: `tsc` clean, `vitest` 3/3 pass, `vite build` succeeds.

## Next (if needed)

- Persist the fraud claim across reloads (already works — DB-backed).
- Extend the triage to include officer workload / team queues (scope creep).