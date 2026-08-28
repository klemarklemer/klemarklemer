# PRD: Taskmaster Claims Operations (hackathon)

**Track:** All Things Agentic Hackathon — The Taskmaster  
**Deadline:** 31 Aug 2026, 17:00 PDT  
**North star (not this build):** [agentic_claims_operations_sla_plan.md](../../agentic_claims_operations_sla_plan.md)  
**Language:** [CONTEXT.md](../../CONTEXT.md)

## Problem

Claims operations stall on manual document checking, assignment, SLA chasing, and reporting. The officer spends time on chores instead of Decisions. Chatbots that “help write text” do not move a Claim. Judges reward agents that take a goal, plan, and **act**.

## Users

| Actor | Role in this slice |
|---|---|
| Claims officer | Primary demo user. Sees Claim detail, Claim events, uploads the missing document as a stand-in for the customer, performs Human approval. |
| Supervisor (ADK) | Event-driven coordinator. Answers “what should happen next?” and invokes specialists. Not a chatbot. |
| Intake / Assignment / SLA / Assessment agents | Specialists that write Claim state and Claim events. |
| Seeded customer | Not a portal. Represented by the officer’s “resubmit document” action. |

## Goal (hackathon)

Ship a **live tracer bullet**: one seeded motor Claim, incomplete documents, then autonomous intake → Assignment → SLA reminder → Assessment recommendation → Human approval → PDF → Claim closed. Prove it on Google Cloud in an unedited ~4 minute video.

## Success for judges

| Rubric | Evidence in this product |
|---|---|
| Operational utility (40%) | Missing police report detected and requested; Assignment written, not suggested; SLA at-risk reminder fired; Assessment recommendation + PDF produced without the officer driving each step. |
| Architecture (30%) | Specialist agents + supervisor; Claim + Claim events; async jobs; Human approval before Decision; traces on agent/model spans. |
| Demo readiness (30%) | Live path; Cloud Run / console / `.run` URL visible; README to reproduce; architecture diagram from the RFC. |

## In scope (must show live)

1. Seed a motor Claim with incomplete documents (police report missing).
2. Intake Agent extracts/validates, records Document completeness = incomplete, notifies in-app.
3. Officer uploads the missing document (customer resubmit stand-in).
4. Intake re-runs, classifies MOTOR / MEDIUM. **Demo fixture:** `survey_required: false` so the 4-minute path reaches assessment. This is a seeded Claim, not a claim that Surveys do not exist (see CONTEXT.md).
5. Assignment Agent scores three seeded claims officers (skill vs workload) and **writes an Assignment**.
6. SLA Agent, via async job: stage-level SLA clock compressed to minutes for the demo; in-app at-risk reminder.
7. Assessment Agent compares seeded Policy + documents → Assessment recommendation `APPROVE` with reasons and confidence.
8. Human approval is **mandatory**. Officer clicks Approve → Decision → Claim closed.
9. PDF assessment report stored; Claim event timeline visible on one ops screen.

Thin ops UI only: Claim detail, Claim events, upload, Approve. No customer portal.

## Out of scope (Phase 2+ / not weekend tickets)

Employee A–D tiers and OKR scoring; real Policy rule engine; Cloud Workflows; full fraud investigator queue; multi line-of-business; Agent Registry / Model Armor / Memory Bank; Veo / Lyria / Gemma unless leftover time; live Surveyor assignment.

## Demo fixture (honest, not a product lie)

The seeded Claim is constructed so Survey is not required. Product language still distinguishes Survey as a domain concept. Do not tell judges “we never need surveys.”

## Notifications

In-app only for the hackathon (missing-document request, SLA at risk, Assignment made, Assessment recommendation ready). No email/SMS.

## Model access

Gemini 3.5 Flash via Vertex AI if Google Cloud credits are attached to the billing project; otherwise Gemini API. Same model generation either way.

## User stories

1. As a claims officer, I want to open a seeded motor Claim, so that the demo starts from a known incomplete file.
2. As a claims officer, I want to see Document completeness and the missing police report, so that I know why the Claim is stuck.
3. As the Supervisor, I want to run Intake after Claim creation, so that extraction and completeness happen without the officer prompting a chat.
4. As a claims officer, I want an in-app notice that a document is missing, so that I can resubmit on behalf of the customer.
5. As a claims officer, I want to upload the police report onto the Claim, so that Intake can run again.
6. As the Supervisor, I want Intake to classify the Claim after completeness, so that Assignment can proceed.
7. As a claims officer, I want the Claim classified as MOTOR / MEDIUM without a Survey, so that the demo reaches assessment in minutes.
8. As the Assignment Agent, I want to score three eligible claims officers, so that workload is not dumped on the highest performer.
9. As a claims officer, I want the Claim assigned automatically to the winning score, so that I do not pick from a list.
10. As a claims officer, I want to see who owns the Claim, so that Assignment is visible.
11. As the SLA Agent, I want a stage-level SLA clock, so that at-risk work is detectable.
12. As a claims officer, I want an in-app SLA at-risk reminder, so that chasing is not a spreadsheet.
13. As the Assessment Agent, I want Policy + documents as evidence, so that the Assessment recommendation is explainable.
14. As a claims officer, I want reasons and confidence on the Assessment recommendation, so that Human approval is informed.
15. As a claims officer, I want Approve to be required even when the recommendation is APPROVE, so that a Decision is never silent.
16. As a claims officer, I want Reject on the recommendation to keep the Claim open with a Claim event, so that overrides are auditable (minimal path: reject stays in assessment; no full rework UI required).
17. As a claims officer, I want a PDF assessment report after Decision, so that the file is complete.
18. As a claims officer, I want the Claim event timeline, so that I can see who/what moved the Claim.
19. As a judge, I want a Cloud Run URL or console in the video, so that Google Cloud is proven.
20. As a judge, I want a public/private repo with spin-up steps, so that the build is reproducible.
21. As operations, I want synthetic names and documents only, so that no real PII is in the demo.

## KPIs shown on the ops screen (not a full command centre)

- Stage SLA remaining / at risk
- Assignment scores for the three seeded officers
- Assessment recommendation vs Decision (acceptance of this one Claim)

## ~4 minute demo script

| Time | Beat |
|---|---|
| 0:00–0:25 | Problem: claims stall on docs, assignment, SLA. This is Taskmaster, not a chatbot. |
| 0:25–0:45 | Open ops UI on Cloud Run URL (show URL bar). Seeded Claim: police report missing. |
| 0:45–1:20 | Intake already ran (or trigger once). Completeness + in-app missing-doc notice. Upload police report. |
| 1:20–1:50 | Classification + Assignment written; show three scores and winner. |
| 1:50–2:20 | SLA clock compressed; reminder appears. Optional: Cloud Tasks / logs with `claim_id` + `agent_name`. |
| 2:20–3:10 | Assessment recommendation APPROVE + reasons. Officer clicks Approve. PDF + Claim closed + timeline. |
| 3:10–3:40 | GCP proof: Cloud Run service / Vertex or API logs. Architecture diagram one slide. |
| 3:40–4:00 | Value: chores done; Human approval still owns the Decision. |

## Submission checklist (product side)

- Category: The Taskmaster  
- Description: features, tech (Gemini 3.5 Flash, ADK, Cloud Run, Firestore, Cloud Storage, Cloud Tasks), synthetic data sources, learnings  
- Repo + README spin-up  
- Architecture diagram (RFC mermaid exported)  
- Demo video as above  
- Optional later: public write-up + `#AllThingsAgenticHackathon`
