# 06 — Assessment recommendation, Human approval, PDF, close

**What to build:** Assessment Agent writes an Assessment recommendation (APPROVE, reasons, confidence) from seeded Policy plus documents. It does not close the Claim. The officer must Approve to create a Decision; then a PDF assessment report is stored and the Claim closes. Reject writes a Claim event and leaves the Claim open. Timeline shows assessment and Human approval.

**Blocked by:** 05 — SLA Agent fires an at-risk reminder on a minute-scale clock

**Status:** ready-for-agent

## Agent handover

Fresh session. Read `CONTEXT.md` (Assessment recommendation, Decision, Human approval), ADR 0001 and 0003, `docs/ui/claims-ops-ux.md` (Human approval pattern), spec stories 16–22, then this ticket only.

`/implement` + `/tdd`: GET shows recommendation; POST approve closes with PDF reference; POST reject keeps Claim open with an event. Stub Gemini with fixed APPROVE JSON in tests. ReportLab (or equivalent) for the PDF. Supervisor invokes Assessment only for this step.

When done: check every box, set **Status: done**, `/code-review`, commit.

## Acceptance criteria

- [ ] Assessment recommendation is stored with reasons and confidence; Claim is not closed by the agent.
- [ ] Approve creates a Decision, generates a PDF in object storage, records a reference, closes the Claim.
- [ ] Reject writes a Claim event and the Claim stays open (no rework wizard).
- [ ] Claim event timeline includes assessment and Human approval (and close on approve).
- [ ] Tests stub the model; they assert Claim + events + PDF reference, not prompts.
- [ ] Ops UI can Approve/Reject and show recommendation + timeline if UI already exists from earlier tickets; Approve/Reject follow the UX Human approval pattern (confirm Decision, closed Claim read-only).
