# 04 — Assignment Agent writes an owner from scores

**What to build:** After classification, the officer does not pick a person. Assignment Agent scores the three seeded claims officers (skill vs workload at minimum) and writes an Assignment. The officer sees the three scores and who won. No Gemini on this path.

**Blocked by:** 03 — Officer resubmits police report; Claim classifies without Survey

**Status:** ready-for-agent

## Agent handover

Fresh session. Read `CONTEXT.md` (Assignment, Claims officer), ADR 0003, spec stories 10–12 and 29, then this ticket only.

`/implement` + `/tdd`: after the classified Claim exists, GET Claim has an owner; scores are visible on GET or a nested field. Deterministic: fixture employees must yield a stable winner in tests. Supervisor invokes Assignment; Intake is not reused for scoring.

When done: check every box, set **Status: done**, `/code-review`, commit.

## Acceptance criteria

- [ ] Assignment is written automatically after classification is available.
- [ ] GET Claim (or documented subresource) shows three scores and the winning claims officer.
- [ ] The overloaded high-skill officer is not blindly chosen if a lower-load officer scores higher (seed data demonstrates this).
- [ ] Gemini is not called for Assignment (tests fail if a model client is required for this path).
- [ ] Claim events include the Assignment.
- [ ] Idempotent: a second Assignment job does not flip owner without a new rule.
