# Backend Architecture & Technical Specification

> **Visual architecture:** [`architecture.html`](architecture.html) renders the system
> topology, the four loops, and the request lifecycle as diagrams. Open it directly in a
> browser — it needs no server.

## 1. System Overview & Architecture (C4 Model)

Taskmaster Claims Operations backend (`apps/api/services/core`) is an enterprise-grade Go service built on top of the **Candi Clean Architecture Framework**. It provides a unified orchestration layer for general-insurance claims lifecycle management, deterministic SLA tracking, autonomous multi-agent choreographies, and mandatory Human-in-the-Loop decision governance.

### 1.1 System Context Diagram (C4 Level 1)

```mermaid
graph TD
    Claimant["👤 Claimant (Policyholder)"]
    Officer["👔 Claims Officer (Human Evaluator)"]
    SupervisorAgent["🤖 Supervisor Agent (Orchestrator)"]

    subgraph Platform ["Taskmaster Claims Operations Platform"]
        Frontend["Web Operations Console (React + Vite)"]
        BackendAPI["Core API Service (Go Candi Monorepo)"]
        PostgresDB[("PostgreSQL 16\n(Relational State + Events)")]
        GCS[("Cloud Storage\n(Document Blobs)")]
        GeminiEngine["Gemini 3.5 Flash Reasoning Engine"]
    end

    Claimant -->|"Submits Claim & Telemetry"| Frontend
    Officer -->|"Performs Human Approval & Reviews"| Frontend
    Frontend -->|"REST / JSON over HTTPS"| BackendAPI
    BackendAPI -->|"ACID Transactions & Audit Events"| PostgresDB
    BackendAPI -->|"Stores & Fetches Document Blobs"| GCS
    BackendAPI -->|"OCR & Assessment Synthesis"| GeminiEngine
    SupervisorAgent -->|"Orchestrates Autonomous Pipeline"| BackendAPI
```

---

### 1.2 Container Architecture (C4 Level 2)

```mermaid
graph LR
    subgraph ClientLayer ["Client Presentation Layer"]
        Console["React 18 SPA (apps/web)"]
    end

    subgraph ServiceLayer ["Backend Core Service (apps/api/services/core)"]
        RESTRouter["HTTP REST Router (Candi/Fiber-compatible)"]
        Middleware["Security & Auth Middleware (JWT/RBAC)"]
        
        subgraph DomainModules ["Domain Modules"]
            ClaimModule["Claim Module (Usecase + Repo)"]
            OfficerModule["Officer Module (Usecase + Repo)"]
            PolicyModule["Policy Module (Usecase + Repo)"]
        end
        
        SharedUOW["Shared Unit of Work / Repository"]
    end

    subgraph DataStorage ["Data & AI Infrastructure"]
        Postgres[("PostgreSQL 16\n(claims, claim_events, officers, policies)")]
        Redis[("Redis 7\n(SLA Cache & Realtime Locks)")]
        CloudStorage[("Google Cloud Storage\n(PDFs, Photos)")]
    end

    Console -->|"HTTP /v1/claim/*"| RESTRouter
    RESTRouter --> Middleware
    Middleware --> ClaimModule
    Middleware --> OfficerModule
    Middleware --> PolicyModule
    ClaimModule --> SharedUOW
    OfficerModule --> SharedUOW
    PolicyModule --> SharedUOW
    SharedUOW -->|"GORM Connection Pool"| Postgres
    SharedUOW -->|"Distributed Cache"| Redis
    ClaimModule -->|"Pre-signed Blob Uploads"| CloudStorage
```

---

### 1.3 Component & Clean Architecture Breakdown (C4 Level 3)

The backend strictly follows Hexagonal / Clean Architecture principles:

```
apps/api/services/core/
├── cmd/
│   └── migration/               # Goose SQL database migrations
├── configs/                     # Environment and database connection configurations
├── internal/modules/
│   ├── claim/
│   │   ├── delivery/resthandler/# HTTP handlers, route registration, and DTO parsing
│   │   ├── domain/              # Request, Response, and Filter DTOs
│   │   ├── repository/          # GORM SQL data access layer and preload logic
│   │   └── usecase/             # Pure business logic and autonomous agent triggers
│   ├── officer/                 # Claims Officer management and workload tracking
│   └── policy/                  # In-force policy contract verification
└── pkg/
    └── shared/
        ├── domain/              # Core domain entities (Claim, Policy, Officer, Event, Assignment)
        ├── repository/          # Centralized repository interface and transaction runner (WithTransaction)
        └── usecase/             # Shared usecase registry and common interfaces
```

---

## 2. Database Schema & Entity-Relationship Diagram (ERD)

```mermaid
erDiagram
    POLICIES ||--o{ CLAIMS : covers
    CLAIMS_OFFICERS ||--o{ CLAIMS : assigned_to
    CLAIMS ||--o{ CLAIM_DOCUMENTS : contains
    CLAIMS ||--o{ CLAIM_EVENTS : appends
    CLAIMS ||--o| ASSIGNMENTS : evaluates
    CLAIMS ||--o| ASSESSMENT_RECOMMENDATIONS : generates
    CLAIMS_OFFICERS ||--o{ ASSIGNMENTS : awarded_to

    POLICIES {
        int id PK
        string policy_number UK
        string policy_holder_name
        string vehicle_plate
        string vehicle_model
        string coverage_type
        decimal max_coverage_amount
        decimal deductible_amount
        timestamp effective_date
        timestamp expiry_date
        string status
    }

    CLAIMS_OFFICERS {
        int id PK
        string name
        string email UK
        string role
        int current_workload
        decimal motor_skill_rating
        boolean is_available
    }

    CLAIMS {
        int id PK
        string claim_number UK
        int policy_id FK
        string stage
        string document_completeness
        boolean survey_required
        string claim_type
        string severity
        text incident_description
        decimal estimated_loss
        decimal approved_amount
        int current_officer_id FK
        timestamp claim_sla_due_at
        timestamp stage_sla_due_at
        string status
        timestamp created_at
        timestamp updated_at
    }

    CLAIM_DOCUMENTS {
        int id PK
        int claim_id FK
        string document_type
        string file_name
        text file_url
        string status
        text extracted_data
        timestamp uploaded_at
    }

    CLAIM_EVENTS {
        int id PK
        int claim_id FK
        string actor_name
        string actor_type
        string action
        string previous_stage
        string new_stage
        text payload
        timestamp created_at
    }

    ASSIGNMENTS {
        int id PK
        int claim_id FK
        int officer_id FK
        decimal workload_score
        decimal skill_score
        decimal total_score
        timestamp assigned_at
    }

    ASSESSMENT_RECOMMENDATIONS {
        int id PK
        int claim_id FK
        string outcome
        decimal confidence
        text reasons
        timestamp generated_at
    }
```

---

## 3. Claim State Machine & Lifecycle Transitions

The Claim progression follows an immutable state machine strictly mapped to the canonical domain language ([`CONTEXT.md`](file:///Users/dimas/Documents/Project/personal/klemarklemer/CONTEXT.md)):

```mermaid
stateDiagram-v2
    [*] --> INTAKE : Claim Created
    INTAKE --> DOCUMENT_VERIFICATION : Missing Required Artifacts (INCOMPLETE)
    DOCUMENT_VERIFICATION --> ASSIGNMENT : Missing Documents Uploaded (COMPLETE)
    INTAKE --> ASSIGNMENT : All Artifacts Present (COMPLETE)
    ASSIGNMENT --> ASSESSMENT : Officer Assigned via Deterministic Scoring
    ASSESSMENT --> DECISION : Assessment Recommendation Generated
    DECISION --> CLOSED : Human Approval (APPROVE) -> Binding Decision
    DECISION --> ASSESSMENT : Human Approval (REJECT) -> Manual Rework
    CLOSED --> [*]
```

### Valid Transition Matrix

| Current Stage | Event / Action Trigger | Target Stage | Precondition Check |
|---|---|---|---|
| `INTAKE` | `CLAIM_CREATED` | `DOCUMENT_VERIFICATION` | Initial intake checks document completeness. |
| `DOCUMENT_VERIFICATION` | `DOCUMENT_UPLOADED` | `ASSIGNMENT` | `POLICE_REPORT` and `DAMAGE_PHOTO` both verified. |
| `ASSIGNMENT` | `OFFICER_ASSIGNED` | `ASSESSMENT` | Deterministic score calculated; officer bound to claim. |
| `ASSESSMENT` | `RECOMMENDATION_GENERATED` | `DECISION` | Coverage limits checked vs estimated loss; recommendation synthesized. |
| `DECISION` | `HUMAN_APPROVAL_RECORDED` | `CLOSED` | Human officer accepts recommendation; binding Decision recorded; settlement calculated. |
| `DECISION` | `HUMAN_APPROVAL_REJECTED` | `ASSESSMENT` | Human officer rejects recommendation for reassessment. |

---

## 4. End-to-End Execution Flows (The 4 Autonomous Claim Loops)

### Loop 1: Ingestion & Autonomous Intake Verification
```mermaid
sequenceDiagram
    autonumber
    actor Officer as Claims Officer
    participant API as Core REST API
    participant IntakeUC as Intake Usecase
    participant DB as PostgreSQL

    Officer->>API: POST /v1/claim/1/documents (police_report.pdf)
    API->>IntakeUC: UploadDocument(claimID, docPayload)
    IntakeUC->>DB: INSERT INTO claim_documents
    IntakeUC->>DB: INSERT INTO claim_events (DOCUMENT_UPLOADED)
    IntakeUC->>IntakeUC: EvaluateIntake(claimID)
    Note over IntakeUC: Checks if POLICE_REPORT and DAMAGE_PHOTO exist
    IntakeUC->>DB: UPDATE claims SET stage='ASSIGNMENT', document_completeness='COMPLETE'
    IntakeUC->>DB: INSERT INTO claim_events (DOCUMENT_VERIFICATION_COMPLETED, CLAIM_CLASSIFIED)
    IntakeUC-->>API: Triggers Loop 2 (RunAssignment)
```

### Loop 2: Deterministic Officer Assignment
```mermaid
sequenceDiagram
    autonumber
    participant AssignUC as Assignment Usecase
    participant DB as PostgreSQL

    AssignUC->>DB: SELECT * FROM claims_officers WHERE is_available = true
    Note over AssignUC: Compute Workload Score: (10 - workload) * 0.5<br/>Compute Skill Score: (motor_skill / 5.0) * 10 * 0.5<br/>Total Score = Workload + Skill
    Note over AssignUC: Alex Rivera scores highest (7.80 pts)
    AssignUC->>DB: INSERT INTO assignments (claim_id, officer_id, scores)
    AssignUC->>DB: UPDATE claims SET current_officer_id=1, stage='ASSESSMENT'
    AssignUC->>DB: UPDATE claims_officers SET current_workload = current_workload + 1
    AssignUC->>DB: INSERT INTO claim_events (OFFICER_ASSIGNED)
    AssignUC-->>AssignUC: Triggers Loop 3 (RunAssessment)
```

### Loop 3: Assessment Recommendation Reasoning
```mermaid
sequenceDiagram
    autonumber
    participant AssessUC as Assessment Usecase
    participant DB as PostgreSQL

    AssessUC->>DB: SELECT * FROM policies WHERE id = claim.policy_id
    Note over AssessUC: Verify Policy #POL-MOTOR-2026-8819 is ACTIVE<br/>Comprehensive coverage ($45,000) > Estimated Loss ($4,200)<br/>Deductible = $500.00
    AssessUC->>DB: INSERT INTO assessment_recommendations (outcome='APPROVE', confidence=0.94)
    AssessUC->>DB: UPDATE claims SET stage='DECISION'
    AssessUC->>DB: INSERT INTO claim_events (RECOMMENDATION_GENERATED)
    AssessUC-->>AssessUC: Halts and awaits Mandatory Human Approval
```

### Loop 4: Human Approval Checkpoint & Binding Decision
```mermaid
sequenceDiagram
    autonumber
    actor Officer as Claims Officer (Alex Rivera)
    participant API as Core REST API
    participant ApprovalUC as SubmitHumanApproval Usecase
    participant DB as PostgreSQL

    Officer->>API: POST /v1/claim/1/approval { action: "APPROVE", notes: "Verified police report KM 42." }
    API->>ApprovalUC: SubmitHumanApproval(claimID, req)
    Note over ApprovalUC: Compute approved_amount = $4,200 - $500 = $3,700.00
    ApprovalUC->>DB: UPDATE claims SET stage='CLOSED', status='CLOSED', approved_amount=3700.00
    ApprovalUC->>DB: UPDATE claims_officers SET current_workload = current_workload - 1
    ApprovalUC->>DB: INSERT INTO claim_events (HUMAN_APPROVAL_RECORDED)
    ApprovalUC->>DB: INSERT INTO claim_events (DECISION_ISSUED)
    ApprovalUC-->>API: Returns Updated Claim in Stage CLOSED
```

---

## 5. Deterministic SLA Clock Formulation

SLA deadlines are calculated deterministically to maintain strict operational compliance:

$$\text{Claim SLA Deadline} = \text{CreatedAt} + 4\text{ hours}$$
$$\text{Stage SLA Deadline} = \text{StageEntryTime} + \Delta T_{\text{stage}}$$

Where $\Delta T_{\text{stage}}$ is allocated as:
- **Intake / Document Verification**: $30\text{ minutes}$
- **Assignment**: $20\text{ minutes}$
- **Assessment**: $15\text{ minutes}$
- **Decision (Human Approval)**: $10\text{ minutes}$

---

## 6. API Error Codes & Standard Response Structure

All API responses follow the standardized Candi JSON envelope:

```json
{
  "code": 200,
  "message": "Success",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "totalRecords": 1,
    "totalPages": 1
  }
}
```

| HTTP Status Code | Description | Example Condition |
|---|---|---|
| `200 OK` | Request succeeded with payload data. | `GET /v1/claim/1`, `POST /v1/claim/1/approval` |
| `201 Created` | Resource created successfully. | `POST /v1/claim` |
| `400 Bad Request` | Malformed JSON or invalid schema parameters. | Missing `action` in Human approval payload. |
| `404 Not Found` | Claim ID, Officer ID, or Policy not found. | `GET /v1/claim/9999` |
| `500 Internal Server Error` | Database transaction failure or unhandled exception. | DB connection timeout during commit. |
