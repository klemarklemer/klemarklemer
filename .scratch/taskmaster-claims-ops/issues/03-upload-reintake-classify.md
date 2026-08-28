# 03 — Officer resubmits police report; Claim classifies without Survey

**What to build:** The officer uploads the police report onto the Claim. Intake runs again. Document completeness becomes complete. Classification is MOTOR / MEDIUM and Survey is not required (demo fixture). Claim events show the upload and the second Intake. Uploads live in object storage; the Claim stores a reference.

**Blocked by:** 02 — Intake records incomplete documents and notifies

**Status:** ready-for-agent

## Agent handover

Fresh session. Read `CONTEXT.md`, spec (stories 6–9, 35), RFC state/storage notes, then this ticket only.

`/implement` + `/tdd`: upload then GET Claim shows complete + MOTOR / MEDIUM + `survey_required` false. Fake model allowed. Do not assign an owner or run assessment.

When done: check every box, set **Status: done**, `/code-review`, commit.

## Acceptance criteria

- [ ] Upload attaches the police report; Claim stores a storage reference, not file bytes in the document store.
- [ ] Second Intake runs because of the upload, not because the officer pasted a chat message.
- [ ] GET Claim shows complete documents, MOTOR, MEDIUM, Survey not required.
- [ ] Claim events include upload and re-Intake.
- [ ] Double-friendly: repeating the same upload job does not create a second classification fight (idempotency key or equivalent).
- [ ] API test covers the path with a stubbed model.
