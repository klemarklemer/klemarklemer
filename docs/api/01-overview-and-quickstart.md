# 01 - API Overview & Quickstart

## 1. Overview

The Taskmaster Claims Operations Core API provides a robust, RESTful interface for managing the complete lifecycle of insurance claims, orchestrating autonomous multi-agent operational workflows, tracking deterministic SLA clocks, and enforcing mandatory human approval checkpoints.

---

## 2. Server Base URLs

| Environment | Base URL | Description |
|---|---|---|
| **Local Development** | `http://localhost:8000` | Local Go Core API service |
| **Staging** | `https://staging-api.klemarklemer.internal` | Staging Cloud Run deployment |
| **Production** | `https://api.klemarklemer.internal` | Production Multi-region cluster |

---

## 3. Quickstart: 4-Minute Demo Tracer Bullet Flow

Follow these sequential cURL commands to test the complete end-to-end tracer bullet path:

### Step 1: Reset Demo to Initial Incomplete State
```bash
curl -X POST "http://localhost:8000/v1/demo/reset" \
  -H "Accept: application/json"
```

### Step 2: Fetch Claim Detail (Stage: DOCUMENT_VERIFICATION, INCOMPLETE)
```bash
curl -X GET "http://localhost:8000/v1/claim/1" \
  -H "Accept: application/json"
```

### Step 3: Upload Missing Police Report (Cascades to Stage: DECISION)
Uploading the missing police report causes the Intake Agent to mark completeness as `COMPLETE`, triggers the Assignment Agent to score officers deterministically and bind Alex Rivera, and triggers the Assessment Agent to synthesize an `APPROVE` recommendation.
```bash
curl -X POST "http://localhost:8000/v1/claim/1/documents" \
  -H "Content-Type: application/json" \
  -d '{
    "document_type": "POLICE_REPORT",
    "file_name": "police_report_incident_km42.pdf",
    "file_url": "https://storage.googleapis.com/klemarklemer-claims-docs/claims/CLM-2026-0042/police_report_incident_km42.pdf"
  }'
```

### Step 4: Submit Human Approval (Creates Binding Decision & Closes Claim)
The human claims officer accepts the recommendation, which creates a binding Decision, calculates payout ($4,200 - $500 = $3,700), appends immutable timeline events, and closes the claim.
```bash
curl -X POST "http://localhost:8000/v1/claim/1/approval" \
  -H "Content-Type: application/json" \
  -d '{
    "officer_id": 1,
    "action": "APPROVE",
    "notes": "Verified police report PR-2026-9912 and photo evidence. $500 deductible applied."
  }'
```
