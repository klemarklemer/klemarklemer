# Separate specialist agents, not one unrestricted agent

Claims operations mix exact chores (assignment scores, SLA clocks) with unstructured
reasoning (incident narrative, coverage, evidence). A single agent holding every tool
would bury those seams, make the human approval gate harder to guarantee, and amount to
a chatbot with tools. So each loop is its own usecase with its own inputs, and the loops
chain directly rather than through an orchestrator.

Gemini is attached where the work is genuinely a judgement:

| Loop | Reasoning | Why |
|---|---|---|
| 1 · Intake | **Gemini** (`gemini.Classifier`) | Severity and whether an inspection is needed are judgements about the incident narrative |
| 2 · Assignment | deterministic | Workload and skill scoring is arithmetic; a model here would only add variance |
| 3 · Assessment | **Gemini** (`gemini.Assessor`) | Coverage against evidence is unstructured, and the reasoning must be auditable |
| 4 · Decision | human | Only this closes a claim |

Both model-backed loops sit behind an interface with a deterministic implementation, so
an unconfigured environment still reasons over real claim data and the startup log names
which engine answered.

**Status:** accepted

## Considered Options

- One agent with all tools — faster to stub, worse to test, and it hides where a human
  must intervene.
- Specialists that invoke each other freely — hidden orchestration and duplicated
  "what happens next" logic.
- A separate supervisor deciding the next stage — a layer whose only input would be the
  claim's own stage, which each loop already knows.
- Loops that chain explicitly, chosen — `EvaluateIntake` hands to assignment or to survey,
  assignment hands to assessment, assessment stops at the human gate. The whole path is
  readable in one file per loop, and every transition writes a `claim_events` row.

## Consequences

Adding a loop means a new usecase and one explicit hand-off, not a rewrite of a monolith
prompt. The cost is that the hand-offs are compiled in rather than configured; changing
the order is a code change.

## Note on an earlier revision

This record previously described a Supervisor plus four ADK specialists, including an SLA
specialist, and said Gemini was attached to Intake and Assessment. None of the ADK,
supervisor or SLA specialist was built, and until the Intake Agent was given a classifier
that loop wrote its severity and survey flag as constants. The record now describes what
the code does.
