# 01 — Seed Claim visible to the claims officer

**What to build:** A claims officer can seed the demo motor Claim and open it: they see identifiers, current Stage, Document completeness still unknown or empty, no owner yet, and an empty Claim event timeline. Synthetic people and Policy exist so later tickets have data. No Intake Agent yet.

**Blocked by:** None — can start immediately.

**Status:** ready-for-agent

## Agent handover

Fresh session. Read `CONTEXT.md`, `docs/adr/0002-firestore-operational-state.md`, `docs/ui/claims-ops-ux.md`, `.scratch/taskmaster-claims-ops/spec.md` (Testing Decisions), then this ticket only.

`/implement` + `/tdd` at the **claims HTTP API** seam (seed + get Claim + list Claim events). Thin ops UI is in scope if it only displays what the API already returns; it must pass the UX Pre-flight.

When done: check every box, set **Status: done**, `/code-review`, commit. Leave Intake, upload, Assignment, SLA, assessment, Cloud Run README for later tickets.

## Acceptance criteria

- [ ] Seed creates one synthetic motor Claim (police report not yet modelled as missing by an agent).
- [ ] GET Claim returns Stage, identifiers, and no Assignment owner.
- [ ] List Claim events returns an empty or seed-only list without crashing.
- [ ] Three synthetic claims officers and one synthetic Policy are loadable for later tickets.
- [ ] API-level test covers seed + get; no Gemini, no ADK specialists.
- [ ] Glossary terms used in API/UI copy (Claim, Stage, Claim event) match CONTEXT.md.
- [ ] If UI shipped: `docs/ui/claims-ops-ux.md` Pre-flight is honest (ops console, reserved layout, no marketing hero).
