# Spec: Taskmaster claims operations tracer bullet

Status: ready-for-agent

**Track:** The Taskmaster  
**Glossary:** [CONTEXT.md](../../CONTEXT.md)  
**PRD:** [docs/prd/taskmaster-claims-ops.md](../../docs/prd/taskmaster-claims-ops.md)  
**RFC:** [docs/rfc/0001-agentic-claims-ops-gcp.md](../../docs/rfc/0001-agentic-claims-ops-gcp.md)  
**ADRs:** 0001 multi-agent ADK, 0002 Firestore state, 0003 rules + Gemini

## Problem Statement

A claims officer cannot get a motor Claim from incomplete documents to a recorded Decision without personally chasing completeness, Assignment, SLA, and paperwork. Chat does not close that gap. For the hackathon, we also cannot prove Google Cloud or autonomous action with a five-phase enterprise platform.

## Solution

An event-driven Supervisor moves one seeded motor Claim through Intake, Assignment, SLA reminder, Assessment recommendation, Human approval, PDF, and close. The officer uses a thin ops UI. The company still owns the Decision. Gemini 3.5 Flash and Google ADK run on Cloud Run with Firestore, Cloud Storage, and Cloud Tasks.

## User Stories

1. As a claims officer, I want a seeded motor Claim with a missing police report, so that the demo starts from a known incomplete Claim.
2. As a claims officer, I want to open Claim detail without chatting, so that I see Stage, Document completeness, and owner.
3. As the Supervisor, I want to run Intake when a Claim is created, so that completeness is recorded without a prompt from the officer.
4. As a claims officer, I want to see that the police report is missing, so that I know what to upload.
5. As a claims officer, I want an in-app notification that documents are incomplete, so that I do not poll the Claim.
6. As a claims officer, I want to upload the police report onto the Claim, so that Intake can run again.
7. As the Supervisor, I want Intake to re-run after upload, so that Document completeness can become complete.
8. As a claims officer, I want classification MOTOR / MEDIUM on the complete Claim, so that Assignment has a type and complexity.
9. As a claims officer, I want the seeded Claim to not require a Survey, so that the live path reaches assessment in minutes.
10. As the Assignment Agent, I want three seeded claims officers with different skill and workload, so that the score is visible and not a coin flip.
11. As a claims officer, I want the Claim assigned to the highest Assignment score automatically, so that I do not pick an owner from a dropdown.
12. As a claims officer, I want to see the three scores and the winner, so that Assignment is explainable.
13. As the SLA Agent, I want a stage-level SLA clock measured in minutes for the demo, so that at-risk can appear in the video.
14. As a claims officer, I want an in-app SLA at-risk reminder, so that chasing is automated.
15. As the SLA Agent, I want Claim events for reminders, so that SLA actions are auditable.
16. As the Assessment Agent, I want seeded Policy text plus uploaded documents, so that the Assessment recommendation cites evidence.
17. As a claims officer, I want an Assessment recommendation of APPROVE with reasons and confidence, so that Human approval is informed.
18. As a claims officer, I want Approve to be required before close, so that a Decision is never implied by the model.
19. As a claims officer, I want Reject to write a Claim event and leave the Claim open, so that an override is visible (no full rework wizard).
20. As a claims officer, I want a generated PDF after Approve, so that the file contains an assessment report.
21. As a claims officer, I want the Claim marked closed after Decision and PDF, so that the lifecycle is complete.
22. As a claims officer, I want a Claim event timeline, so that I can see Intake, Assignment, SLA, assessment, and Human approval in order.
23. As a judge, I want the UI served from Cloud Run, so that Google Cloud is visible in the URL or console.
24. As a judge, I want logs containing claim_id and agent_name, so that agent work is not only UI chrome.
25. As a developer, I want README steps for local and Cloud Run, so that the project is reproducible.
26. As operations, I want only synthetic people and documents, so that no real PII is stored.
27. As the Intake Agent, I want structured extraction (Claim identifiers, document checklist), so that completeness is data, not a paragraph.
28. As the Supervisor, I want specialists not to call each other, so that “what next” lives in one place.
29. As a developer, I want Assignment and SLA without Gemini, so that those chores are stable in a live demo.
30. As a developer, I want Gemini 3.5 Flash on Intake and Assessment only, so that we meet the model requirement without putting the model on clocks.
31. As a claims officer, I want idempotent retries of the same job, so that double-click upload does not double-assign.
32. As a developer, I want OpenTelemetry spans around Supervisor, specialists, and Gemini, so that architecture discipline is demonstrable.
33. As a claims officer, I want current Stage on the Claim, so that I know where the SLA clock applies.
34. As a developer, I want Cloud Tasks for intake, SLA tick, assessment, and PDF, so that HTTP requests do not wait on Gemini.
35. As a claims officer, I want object storage for PDFs and uploads, so that Firestore is not a blob store.

## Implementation Decisions

- Single Cloud Run service: HTTP API plus worker entry for Cloud Tasks.
- Google ADK Supervisor plus Intake, Assignment, SLA, Assessment specialists; specialists do not invoke each other.
- Gemini 3.5 Flash via Vertex AI if the billing project has credits, else Gemini API.
- Firestore is the system of record for Claim, Claim events, seeded employees, seeded Policy, notifications, and stored Assessment recommendations.
- Cloud Storage for submitted and generated artifacts; Firestore holds references.
- Cloud Tasks job kinds: document intake, SLA tick, assessment, report PDF. Payload is claim id, job kind, idempotency key.
- Assignment score is deterministic code over seeded employees. SLA remaining and at-risk are deterministic code. Demo stage SLA clocks are minutes, not business days.
- Seeded Claim: motor, police report missing first, `survey_required` false after classification.
- Human approval mandatory; Assessment Agent writes an Assessment recommendation only.
- In-app notifications only.
- Ops UI: Claim detail, timeline, upload, Approve/Reject. No customer portal.
- Observability: structured logs with claim_id, agent_name, action; OTel traces across API, worker, agents, Gemini.
- Auth for the public URL sufficient to stop anonymous credit drain (API key or equivalent).
- ReportLab for the assessment PDF after Decision.

## Testing Decisions

A good test asserts observable Claim behaviour: Stage, Document completeness, Assignment owner, SLA at-risk notification, Assessment recommendation vs Decision, Claim event sequence, PDF reference present. It does not assert prompt text, ADK internals, or Firestore layout.

Seam (one, highest): the **claims HTTP API** as the claims officer would use it (create/seed, get Claim, upload, list Claim events, approve). Agent internals are behind that seam. Assignment winner and SLA at-risk should be testable with Gemini stubbed or not called. Intake/Assessment tests may use a fake model that returns fixed JSON.

There is no prior test suite in this repo; this is the first application code. Prefer a small API-level test (or in-process app test) over UI snapshots.

## Out of Scope

Employee performance tiers and OKRs; real Policy engine; Cloud Workflows; live Survey assignment; fraud investigator product; customer portal; email; GEAP Agent Registry, Identity, Gateway, Model Armor, Memory Bank; Veo/Lyria/Gemma; multi line-of-business; straight-through processing without Human approval.

## Further Notes

Deadline 31 Aug 2026 17:00 PDT. Demo script is in the PRD. Vision document remains the north star, not the ticket. If the API seam is wrong, say so before implement — the alternative would be testing each specialist in isolation, which splits the Taskmaster story.

## Comments
