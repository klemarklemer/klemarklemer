# Use multiple ADK specialist agents, not one unrestricted agent

Claims operations mix exact chores (Assignment scores, SLA clocks) with unstructured reasoning (documents, coverage narrative). A single god agent would bury those seams, make Human approval harder to guarantee, and look like a chatbot with tools. We run a Supervisor that only answers “what should happen next?” from Claim state, and four specialists (Intake, Assignment, SLA, Assessment) that do not call each other. Gemini is attached to Intake and Assessment only.

**Status:** accepted

## Considered Options

- One ADK agent with all tools — faster to stub, worse to test and to demo as Taskmaster architecture.
- Specialists that invoke each other — hidden orchestration, duplicate “what next” logic.
- Supervisor + specialists (chosen) — matches the vision and the judging architecture axis.

## Consequences

New Stage transitions go through the Supervisor. Adding a Survey specialist later is another specialist plus Supervisor rules, not a rewrite of a monolith prompt.
