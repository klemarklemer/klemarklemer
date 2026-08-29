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
