# 07 — Cloud Tasks, Cloud Run proof, logs, README

**What to build:** A judge can believe this runs on Google Cloud: jobs (intake, SLA tick, assessment, PDF) go through Cloud Tasks (or a documented worker URL the same service exposes), structured logs include `claim_id` and `agent_name`, the public URL is not an open credit drain, and README has local plus `gcloud run deploy` spin-up. OTel spans cover Supervisor, specialists, and Gemini when called.

**Blocked by:** 06 — Assessment recommendation, Human approval, PDF, close

**Status:** ready-for-agent

## Agent handover

Fresh session. Read spec stories 23–25, 31–32, 34, RFC cost/threat/README sections, then this ticket only.

`/implement`: wire Cloud Tasks (local emulator or skip-to-in-process behind a flag is acceptable if README says how demo deploy uses real queues). Add API key or equivalent. README: env vars, seed, local tests, deploy, what to show in the 4-minute video (Cloud Run URL, log line). `/code-review`, commit.

Do not add Fleet-track products, email, Survey runtime, or employee tiers.

## Acceptance criteria

- [ ] Worker entry exists for Cloud Tasks payloads: claim id, job kind, idempotency key.
- [ ] Logs on an agent action include `claim_id` and `agent_name`.
- [ ] Public HTTP is gated enough that anonymous strangers cannot burn Gemini quota.
- [ ] README: local run, tests, deploy to Cloud Run, required Google Cloud services (Run, Firestore, Storage, Tasks).
- [ ] README names Gemini 3.5 Flash and Google ADK.
- [ ] Architecture diagram pointer: RFC mermaid is the submission diagram (export/screenshot instructions one paragraph).
