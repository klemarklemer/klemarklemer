# 05 — SLA Agent fires an at-risk reminder on a minute-scale clock

**What to build:** The current Stage has an SLA clock measured in minutes for the demo. After a tick, the officer sees an in-app at-risk reminder and a Claim event. Remaining time / at-risk is visible on the Claim. Math is code, not Gemini.

**Blocked by:** 04 — Assignment Agent writes an owner from scores

**Status:** ready-for-agent

## Agent handover

Fresh session. Read `CONTEXT.md` (SLA clock, Stage), ADR 0003, spec stories 13–15, then this ticket only.

`/implement` + `/tdd`: arrange a Stage deadline in the past or near past, run `sla_tick` (in-process is fine), GET Claim / notifications show at-risk. Do not add assessment or PDF.

When done: check every box, set **Status: done**, `/code-review`, commit.

## Acceptance criteria

- [ ] Stage-level SLA clock exists on the Claim (demo scale: minutes).
- [ ] A tick produces an in-app at-risk notification when remaining time is in the at-risk window or elapsed as specified in code.
- [ ] Claim events include the SLA reminder with SLA agent as actor.
- [ ] Gemini is not used to compute remaining time or at-risk.
- [ ] API test does not sleep for real minutes (inject clock or set deadline in the past).
