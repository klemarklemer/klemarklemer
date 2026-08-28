# 02 — Intake records incomplete documents and notifies

**What to build:** After seed (or on create), the Supervisor runs Intake without the officer chatting. The officer sees Document completeness incomplete, police report missing, an in-app notification, and Claim events for Intake. Extraction is structured data. Gemini is faked or stubbed in tests; production may call Gemini 3.5 Flash for Intake only.

**Blocked by:** 01 — Seed Claim visible to the claims officer

**Status:** ready-for-agent

## Agent handover

Fresh session. Read `CONTEXT.md`, `docs/adr/0001-adk-multi-agent-not-one-god-agent.md`, `docs/adr/0003-deterministic-rules-plus-gemini.md`, spec Implementation + Testing Decisions, then this ticket only.

`/implement` + `/tdd` at the claims HTTP API: after seed, GET Claim shows missing police report and a notification exists; Claim events include Intake. Stub the model with fixed JSON. Supervisor may invoke Intake only — no Assignment, SLA, or Assessment specialists yet. Jobs may run in-process (Cloud Tasks come in 07).

When done: check every box, set **Status: done**, `/code-review`, commit.

## Acceptance criteria

- [ ] Seed/create leads to Intake without a chat prompt from the officer.
- [ ] GET Claim shows Document completeness incomplete and police report missing.
- [ ] An in-app notification exists for the missing document.
- [ ] Claim events include an Intake action with actor as the Intake agent name.
- [ ] Specialists do not call each other; Supervisor decides this step.
- [ ] Tests use a fake model; they assert Claim behaviour, not prompt text.
