# ADR 0005: Security, RBAC, and Secure SDLC (SSDLC) Architecture

## Status
Accepted

## Context
Taskmaster Claims Operations orchestrates agentic automation alongside human claims officers in an enterprise general-insurance context. Processing claims involves sensitive policyholder information (Personally Identifiable Information - PII, vehicle telematics, damage evidence, police reports) and binding financial payouts.

To satisfy insurance regulatory standards, competition evaluation criteria, and enterprise security principles, the system requires an explicit security architecture covering:
1. Role-Based Access Control (RBAC) with strict Human-In-The-Loop (HITL) enforcement for binding Decisions.
2. Data privacy and PII protection across immutable audit trails.
3. Cryptographic integrity and tamper-evident event logging.
4. An automated DevSecOps pipeline enforcing Continuous Security in CI/CD.

## Decision

### 1. Dual-Tier Security & Strict RBAC Boundary
We enforce two distinct security zones:
- **Public / Claimant Tier**: Rate-limited ingestion endpoints (`POST /v1/claim`, document upload) protected by WAF and Cloudflare Turnstile token validation.
- **Operational Console Tier**: Authenticated via OpenID Connect (OIDC) / JWT Bearer tokens with strict RBAC:
  - `claims_officer`: Can view assigned claims, upload supplemental artifacts, trigger re-evaluation, and submit Human approval Decisions.
  - `senior_officer` / `claims_manager`: Full override authority, workload adjustments, and SLA policy tuning.
  - `system_agent` (Autonomous Agent Worker): Can execute deterministic triage, intake classification, workload scoring, and generate `AssessmentRecommendation`. **Explicitly forbidden** from executing `SubmitHumanApproval` or closing Claims.

### 2. PII Protection and Field-Level Data Masking
- **In-Transit & At-Rest Encryption**: TLS 1.3 for all external/internal REST communications; PostgreSQL tablespace encryption (AES-256) at rest.
- **Event Log Redaction**: `claim_events.payload` stores JSON payloads with masked sensitive fields (e.g., driver license numbers, national IDs, and unredacted customer phone numbers are hashed or tokenized).
- **Secure Object Storage**: Binary documents (photos, PDFs) are stored in private Cloud Storage buckets with time-limited Pre-Signed URLs (TTL: 15 minutes) for authorized claims officer viewing.

### 3. Parameterized Query Security
- All database interactions use GORM ORM with strictly parameterized SQL statements, eliminating SQL Injection (SQLi) vectors across all filter and search queries.

### 4. DevSecOps & Automated Security Gates in CI
The GitHub Actions CI pipeline enforces automated security gates prior to merge:
- **SAST (Static Application Security Testing)**: `govulncheck` (Go official vulnerability scanner) and `gosec` (Go security AST analyzer).
- **Secret Scanning**: Automated scans preventing API keys, credentials, or private certificates from entering git history.
- **Dependency Vulnerability Auditing**: Automated dependency vulnerability scanning.

## Consequences

### Positive
- **Regulatory Compliance**: Aligns with global insurance privacy regulations and zero-trust operational requirements.
- **Auditability**: Tamper-evident `claim_events` timeline provides legal non-repudiation for all automated and human decisions.
- **Error & Fraud Prevention**: Autonomous agents cannot accidentally or maliciously approve settlements without human approval.
- **Competition Readiness**: Demonstrates enterprise-grade security posture and SDLC maturity.

### Negative / Trade-offs
- Local developer experience requires mock officer context headers (`X-Officer-ID`) or development JWT tokens when developing without external identity providers.
