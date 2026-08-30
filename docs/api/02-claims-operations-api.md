# 02 - Claims Operations API Reference

This document provides exhaustive reference documentation for all Claim lifecycle and agentic workflow endpoints.

---

## 1. Endpoints Summary

| Method | Endpoint | Description | Auth / Role |
|---|---|---|---|
| `GET` | `/v1/claim` | List paginated claims with stage & keyword filters | Public / Officer |
| `POST` | `/v1/claim` | Ingest a new claim against an in-force Policy | Public / Officer |
| `GET` | `/v1/claim/:id` | Get full claim aggregate with documents, events, & scores | Officer |
| `POST` | `/v1/claim/:id/documents` | Upload document and trigger autonomous progression | Public / Officer |
| `POST` | `/v1/claim/:id/intake` | Manually evaluate intake document completeness | Officer / Agent |
| `POST` | `/v1/claim/:id/assignment` | Manually execute deterministic officer assignment | Officer / Agent |
| `POST` | `/v1/claim/:id/assessment` | Manually generate assessment recommendation | Officer / Agent |
| `POST` | `/v1/claim/:id/approval` | Submit human approval to record a binding Decision | **Human Officer Only** |
| `POST` | `/v1/demo/reset` | Reset demo scenario to initial state | Dev / Staging |

---

## 2. Detailed Endpoints

### 2.1 List Claims (`GET /v1/claim`)

- **Query Parameters:**
  - `page` *(int, default: 1)*: Page number
  - `limit` *(int, default: 10)*: Page size
  - `stage` *(string)*: `INTAKE`, `DOCUMENT_VERIFICATION`, `ASSIGNMENT`, `ASSESSMENT`, `DECISION`, `CLOSED`
  - `status` *(string)*: `OPEN`, `CLOSED`
  - `search` *(string)*: Substring search across Claim Number and Incident Description

```bash
curl -X GET "http://localhost:8000/v1/claim?page=1&limit=10&status=OPEN"
```

---

### 2.2 Create Claim (`POST /v1/claim`)

- **Request Body:**
```json
{
  "policy_id": 1,
  "claim_type": "MOTOR",
  "incident_description": "Front bumper collision with road barrier during heavy rain on highway KM 42.",
  "estimated_loss": 4200.00
}
```

```bash
curl -X POST "http://localhost:8000/v1/claim" \
  -H "Content-Type: application/json" \
  -d '{"policy_id":1,"claim_type":"MOTOR","incident_description":"Collision KM 42","estimated_loss":4200.00}'
```

---

### 2.3 Get Claim Detail (`GET /v1/claim/:id`)

Fetches complete entity graph: Policy contract, assigned Officer, uploaded Documents, immutable Claim Events, Assignment scoring, and Assessment Recommendation.

```bash
curl -X GET "http://localhost:8000/v1/claim/1"
```

---

### 2.4 Upload Document (`POST /v1/claim/:id/documents`)

Appends a document artifact and cascades through **Intake $\to$ Assignment $\to$ Assessment** in a single autonomous request.

- **Request Body:**
```json
{
  "document_type": "POLICE_REPORT",
  "file_name": "police_report_incident_km42.pdf",
  "file_url": "https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0042/police_report_incident_km42.pdf"
}
```

```bash
curl -X POST "http://localhost:8000/v1/claim/1/documents" \
  -H "Content-Type: application/json" \
  -d '{"document_type":"POLICE_REPORT","file_name":"police_report_incident_km42.pdf"}'
```

---

### 2.5 Submit Human Approval (`POST /v1/claim/:id/approval`)

**Mandatory Governance Gate**: Only human credentials can invoke this endpoint. Approving creates a binding Decision, calculates the settlement (`approved_amount = estimated_loss - deductible`), and closes the claim.

- **Request Body:**
```json
{
  "officer_id": 1,
  "action": "APPROVE",
  "notes": "Verified police report PR-2026-9912. $500 deductible applied."
}
```

```bash
curl -X POST "http://localhost:8000/v1/claim/1/approval" \
  -H "Content-Type: application/json" \
  -d '{"officer_id":1,"action":"APPROVE","notes":"Approved after police report verification."}'
```

## Notes for client developers

These are the details most likely to cost you a wasted round trip or a wrong assumption.

**Claim detail already contains the objects you need.** `GET /v1/claim/{id}` preloads
`policy`, `current_officer`, `documents`, `events`, `assignment` and `recommendation`.
Render straight from one response; there is no need to resolve `policy_id` or
`current_officer_id` yourself.

**Uploading a document runs three loops.** `POST /v1/claim/{id}/documents` completes
intake, assigns an officer and produces the assessment before it returns. Expect the
response to come back with `stage: "DECISION"`, a populated `current_officer`, and a
`recommendation` — not the `DOCUMENT_VERIFICATION` state you posted against. The
individual loop endpoints exist for stepping through the stages one at a time.

**A recommendation is not a decision.** `recommendation.outcome` is `APPROVE`, `REJECT`
or `MANUAL_REVIEW`, and `confidence` runs 0.0–1.0. Nothing is settled until
`POST /v1/claim/{id}/approval` records a human decision; only then do `status` become
`CLOSED` and `approved_amount` become `estimated_loss` minus the policy deductible.
Treat `MANUAL_REVIEW` as a first-class outcome in the UI — it is the agent's honest
answer when the rules disagree, not an error.

**Every recommendation names the engine that produced it.** The
`RECOMMENDATION_GENERATED` event's `payload` is a JSON string carrying `outcome`,
`confidence`, `reasons` and `engine`. `engine` reads `gemini-3.5-flash via vertex-ai`
when a model answered, or `deterministic-rules` when the service is running without a
model credential and fell back to local underwriting rules. Surfacing it is useful:
it tells an operator whether they are looking at model reasoning or the fallback.

**Errors use the same envelope as success.** Check `success`, not the presence of
`data`. A failed loop returns HTTP 400 with a `message` naming the cause, for example
`no claims officers exist to assign claim 1`.
