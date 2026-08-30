# Taskmaster Claims Operations — Devpost submission

Track: **The Taskmaster**

---

## What it does

Motor insurance claims stall in predictable places. Documents arrive incomplete.
Assignment is manual, so workloads skew. Assessment waits on a human to open a policy
and do arithmetic. Nobody can see why a claim is late until someone asks.

Taskmaster Claims Operations puts four autonomous loops around that process. They do the
repetitive work and surface the evidence — and the loop that costs money stops at a human.

| Loop | Agent | What it does |
|---|---|---|
| 1 | Intake | Verifies document completeness, classifies the claim |
| 2 | Assignment | Scores every officer on workload and skill, assigns the best fit |
| 3 | Assessment | Weighs policy cover against estimated loss and evidence — **Gemini 3.5 Flash** |
| 4 | Decision gate | A claims officer binds the decision; only this closes a claim |

Loops 1–3 chain autonomously inside a single request. Loop 4 never does. An `APPROVE`
from the Assessment Agent is a proposal, not a settlement.

Every stage transition writes an immutable event naming the actor, the stage change and
the reasoning — including which engine produced each recommendation — so a claim's whole
history is replayable.

## Features

- **Autonomous claim progression.** Uploading the final document runs intake, assignment
  and assessment in one request and returns the claim at the decision gate.
- **Grounded assessment.** Gemini receives the coverage ceiling, deductible, in-force
  window, estimated loss and the documents on file, and returns a schema-constrained
  outcome, confidence and justification citing those exact figures.
- **Deterministic officer scoring.** Workload and skill weighted 50/50, with the score
  recorded on the assignment so a decision can be audited later.
- **Human decision gate.** Settlement is computed as estimated loss minus deductible,
  floored at zero, and recorded with the officer who approved it.
- **Replayable audit trail.** Every transition is an immutable event carrying actor,
  stage change and payload.
- **Honest degradation.** With no model credential the service still boots and applies
  the same underwriting rules deterministically, reporting which engine answered rather
  than imitating a model.

## Technologies used

| Layer | Technology |
|---|---|
| Reasoning | **Gemini 3.5 Flash** via **Vertex AI**, Google **GenAI SDK** for Go |
| Service | Go 1.26, Candi framework, Echo REST |
| Persistence | PostgreSQL 16, GORM, Goose migrations |
| Coordination | Redis 7 |
| Console | React 18, Vite, Tailwind v4 |
| DevSecOps | govulncheck, gosec, GitHub Actions, Cloud Run |
| Observability | OpenTelemetry tracing throughout |

The reasoning layer sits behind an `Assessor` interface, so the backend is selected by
environment alone: Vertex AI via application default credentials, the Gemini API via key,
or deterministic rules. The same binary runs all three ways.

## Data sources

Synthetic demo data only — no real customer or policy records. The seeded scenario is one
motor claim (`CLM-2026-0042`) against a comprehensive policy with a 45,000 ceiling and a
500 deductible, three claims officers with differing workloads and skill ratings, and two
document types the intake loop checks for. `POST /v1/demo/reset` restores it idempotently.

## Findings and learnings

**A green pipeline can hide an empty one.** The service built and tested clean while every
test in the suite was `t.Skip()` — 18 of 21 functions, zero assertions executing. The
assessment logic that most needed cover had no test file at all. CI told us nothing until
we asked what it was actually asserting.

**"Agentic" is easy to claim and easy to check.** The first assessment implementation
loaded the claim and then read no field from it, returning a fixed `APPROVE` at `0.94`
confidence with a justification citing another claim's police report. It looked
convincing in a single demo. Changing one input exposed it instantly. We now test that
two materially different claims cannot return the same recommendation.

**Failures that report success are the expensive ones.** Three separate defects shared
that shape. Migrations collided on a version number, so nothing could build a schema from
scratch. The demo-reset endpoint deleted its claim and "recreated" it with an UPDATE that
matched zero rows and returned nil — destroying the demo on first press. Officer lookup
ran `LIMIT 0` against a populated table, so assignment could never assign anyone. Each
was masked by an error path that collapsed distinct causes into one message.

**Grounding beats prompting.** The reliability gain came from passing the actual policy
ceiling, deductible and loss and constraining the response schema — not from a longer
prompt. Unrecognised outcomes and empty reasoning are refused rather than written to the
claim, so a model that drifts fails loudly instead of quietly recording something wrong.

**Degrading honestly is a design decision.** Without a credential the service applies the
same rules locally and says so in its startup log and in every recorded event. A keyless
CI run or a fresh clone behaves truthfully instead of appearing to be a model.

**Two wallets, one project.** Gemini API calls failed with depleted prepay credits while
the same Google Cloud project had billing enabled. AI Studio prepay and Cloud billing are
separate; moving to Vertex AI used the funded one and put every call in Cloud Console.

## Known limitations

The claims officer console is a **design prototype**. It runs on local fixture data and
makes no HTTP calls, so it does not yet exercise the backend described above. The
architecture diagram marks that edge as unwired rather than implying a live client.

Loop 1 still records some classification fields as constants rather than deriving them
from the claim.
