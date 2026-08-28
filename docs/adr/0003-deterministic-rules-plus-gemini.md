# Deterministic rules for SLA and Assignment; Gemini for extraction and reasoning only

SLA clocks and Assignment scores must be repeatable in a demo and explainable to a claims officer. If Gemini computes “who is free” or “when we breach,” numbers drift and Human approval cannot trust the chores. Code owns SLA remaining/at-risk and the Assignment score (skill, workload, and the other documented factors we seed). Gemini 3.5 Flash owns document understanding, classification assistance, and Assessment recommendation narrative with evidence. A Decision is only Human approval (this cut), never the model.

**Status:** accepted

## Considered Options

- Ask Gemini to enforce Policy wording and SLAs — one prompt, unbounded liability and flaky demos.
- Rules engine product (Drools, etc.) — too much for seeded Policy JSON.
- Split (chosen) — small Python (or equivalent) functions for clocks and scores; Gemini JSON for Intake/Assessment; Supervisor sequences them.

## Consequences

Assessment recommendation can be wrong; the officer still decides. Tests can lock Assignment winner and SLA at-risk without calling the model. Prompt/version fields on Claim events still record which model produced the recommendation.
