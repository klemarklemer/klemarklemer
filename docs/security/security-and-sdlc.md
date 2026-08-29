# Security Architecture & Secure SDLC (SSDLC) Specification

## 1. Executive Summary

Taskmaster Claims Operations is built upon a **Zero-Trust, Security-by-Design** foundation tailored for enterprise insurance operations. Handling sensitive policyholder telematics, police accident reports, and financial disbursements requires rigorous security controls across all stages of the Software Development Life Cycle (SDLC).

This document outlines the end-to-end security architecture, Threat Modeling (STRIDE), Role-Based Access Control (RBAC), data protection standards, and automated DevSecOps pipelines.

---

## 2. Threat Modeling & Risk Mitigation (STRIDE Analysis)

| Threat Category | Potential Vector | Mitigation & Technical Countermeasure |
|---|---|---|
| **Spoofing** | Rogue actor pretending to be a Claims Officer to approve payouts. | JWT/OIDC authentication with Claims Officer identity binding, cryptographic signature validation, and officer ID claims verification. |
| **Tampering** | Altering damage estimates or manipulating timeline history. | Immutable append-only `claim_events` table; GORM parameterized queries; database transaction rollback on any stage validation failure. |
| **Repudiation** | An officer denying having approved an erroneous or fraudulent claim. | Mandatory `HUMAN_APPROVAL_RECORDED` audit events storing authenticated `officer_name`, `officer_id`, timestamp, and approval rationale. |
| **Information Disclosure** | Unauthorized exposure of policyholder PII or police reports. | Field-level data masking in event payloads, TLS 1.3 encryption in-transit, private Google Cloud Storage with short-lived Pre-Signed URLs (15m TTL). |
| **Denial of Service** | Flooding claim submission or document upload endpoints. | Public gateway rate limiting (100 req/min/IP), payload size limits (15MB for documents), connection pooling in PostgreSQL. |
| **Elevation of Privilege** | Autonomous background agents auto-approving or closing claims. | Hard architectural separation: Agent workers only produce `AssessmentRecommendation` (`APPROVE`/`REJECT`/`MANUAL_REVIEW`); only human credentials can invoke `POST /v1/claim/:id/approval`. |

---

## 3. Role-Based Access Control (RBAC) Matrix

| Resource / Endpoint | Claimant (Public) | Claims Officer | Senior Officer / Manager | System Agent (Worker) |
|---|---|---|---|---|
| `POST /v1/claim` (Submit Claim) | ✅ Allowed | ✅ Allowed | ✅ Allowed | ❌ Denied |
| `GET /v1/claim` (List Claims) | ❌ Denied | ✅ Assigned / Pool | ✅ All Claims | ✅ Read-only |
| `GET /v1/claim/:id` (Claim Detail) | ❌ (Own claim via Token) | ✅ Assigned / Pool | ✅ All Claims | ✅ Read-only |
| `POST /v1/claim/:id/documents` | ✅ (Intake stage) | ✅ Allowed | ✅ Allowed | ❌ Denied |
| `POST /v1/claim/:id/intake` (Evaluate) | ❌ Denied | ✅ Allowed | ✅ Allowed | ✅ Autonomous Trigger |
| `POST /v1/claim/:id/assignment` (Assign) | ❌ Denied | ❌ Denied | ✅ Override | ✅ Autonomous Trigger |
| `POST /v1/claim/:id/assessment` (Assess) | ❌ Denied | ❌ Denied | ✅ Allowed | ✅ Autonomous Trigger |
| `POST /v1/claim/:id/approval` (Decision) | ❌ Denied | ✅ **Human Only** | ✅ **Human Only** | 🛑 **BLOCKED BY DESIGN** |
| `POST /v1/demo/reset` | ❌ Denied (Dev Only) | ❌ Denied | ✅ Dev / Staging | ❌ Denied |

---

## 4. PII Protection & Data Privacy Architecture

```
[ Claimant / Officer ]
         │
         ▼  (TLS 1.3 / HTTPS)
[ API Gateway / Cloud Run ]
         │
    ┌────┴───────────────────────────────────────────────┐
    │ Field-Level PII Masking & Tokenization             │
    │  - Name: J*** D**                                  │
    │  - Plate: B **** UQ                                │
    │  - Event Payload: tokenized reference              │
    └────┬───────────────────────────────────────────────┘
         ├──────────────────────────────┬──────────────────────────────┐
         ▼                              ▼                              ▼
[ PostgreSQL (AES-256) ]      [ Cloud Storage (SSE) ]        [ Gemini / ADK ]
  - Relational claims state      - Raw incident photos          - Unstructured OCR &
  - Immutable claim_events       - Police report PDFs             analysis (Ephemeral,
  - Encrypted tablespaces        - Short-lived signed URLs        no training retention)
```

1. **In-Transit Security**: Strict TLS 1.3 with HSTS headers across all REST endpoints.
2. **At-Rest Encryption**: PostgreSQL database files and Cloud Storage blobs encrypted with customer-managed or Google-managed AES-256 encryption keys.
3. **Pre-Signed Storage URLs**: Document binaries (damage photos, police reports) are never exposed via public URLs. The API dynamically generates 15-minute Pre-Signed URLs upon authenticated request.
4. **Data Minimization in AI Prompts**: Only relevant claim context (damage description, coverage type, deductible) is sent to Gemini reasoning models; personally identifying information like bank account numbers or national IDs are stripped beforehand.

---

## 5. Defense-in-Depth Implementation

### 5.1 SQL Injection Immunity
All database queries in the Candi Go backend utilize GORM's parameterized query engine. Raw string concatenations for SQL generation are strictly prohibited:

```go
// Example: Safe parameterized query in claim repository
db.Where("(claim_number ILIKE '%%' || ? || '%%' OR incident_description ILIKE '%%' || ? || '%%')", filter.Search, filter.Search)
```

### 5.2 Input Validation & JSON Schema Enforcement
Every incoming REST request is unmarshaled into strongly typed Go structs with field-level bounds checking, regex validation for alphanumeric claim numbers, and maximum payload size limits.

### 5.3 Deterministic SLA Clocks & Transaction Isolation
All stage transitions, assignment updates, and event logging execute within ACID-compliant PostgreSQL database transactions (`WithTransaction`). If any event recording fails, the entire stage transition is rolled back to prevent inconsistent states.

---

## 6. Secure SDLC (SSDLC) & DevSecOps Pipeline

The repository integrates automated security gates directly into GitHub Actions (`.github/workflows/ci.yml`):

```
 Git Push / Pull Request
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│ 1. Code Quality & Formatting Gate                           │
│    - go vet, golangci-lint, eslint, tsc --noEmit            │
└────────┬────────────────────────────────────────────────────┘
         ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. SAST (Static Application Security Testing) Gate          │
│    - govulncheck (Official Go Vulnerability Database)       │
│    - gosec (Go Security AST Rule Scanner)                   │
└────────┬────────────────────────────────────────────────────┘
         ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Automated Test Verification & Integration Suite          │
│    - go test ./services/core/... (API)                      │
│    - npm run test (Web Console Vitest + RTL)                │
└────────┬────────────────────────────────────────────────────┘
         ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Build & Artifact Verification                            │
│    - Docker multi-stage non-root container build            │
└─────────────────────────────────────────────────────────────┘
```

### Security Scanning Toolset
- **govulncheck**: Scans Go codebase against the Go Vulnerability Database (Go team maintained).
- **gosec**: Scans Go AST for security vulnerabilities (unhandled errors, SQL injection risks, weak crypto, unsafe memory allocation).
- **Non-Root Container Execution**: Docker container runs as unprivileged user `appuser` (UID: 10001) to prevent container breakout risks.
