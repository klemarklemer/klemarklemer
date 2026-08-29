# ADR 0004: PostgreSQL with GORM and Candi Framework for Backend State and Agent Orchestration

## Context and Problem Statement

The platform requires a robust, relational, and auditable operational state store for Claims, Policies, Claims Officers, and append-only Claim Events. Initially, ADR 0002 selected Firestore. However, managing relational queries, atomic transactions during Stage transitions, and deterministic assignment scoring across multiple entities benefits from a relational schema with strong transactional guarantees. Furthermore, a clean multi-agent architecture in Go needs a unified microservice framework with built-in dependency injection, REST delivery, and async task queue workers.

## Decision Drivers

- **Clean Architecture & Maintainability**: Standardized 4-layer architecture (delivery, domain, repository, usecase) via the Candi framework.
- **Relational Integrity & ACID Transactions**: Atomic stage transitions and claim timeline event logging using PostgreSQL and GORM with `WithTransaction`.
- **Built-in Async Task Queue & Schedulers**: Candi's native Redis-backed Task Queue Worker and Cron Scheduler for SLA clocks and async document extraction.
- **Unified Native Go Multi-Agent Layer**: Native Go implementation of the Supervisor and specialist agents (Intake, Assignment, SLA, Assessment) using Google GenAI SDK.
- **Cloud Run Deployment**: Single portable container deployable to GCP Cloud Run, connecting to Cloud SQL (PostgreSQL), Memorystore (Redis), and Google Cloud Storage (GCS).

## Considered Options

1. **Firestore + Cloud Tasks + Python ADK** (ADR 0002 baseline): Good for noSQL prototyping, but lacks relational joins and requires multi-service choreography.
2. **Go Candi with PostgreSQL, GORM, Redis Task Queue, and GenAI Go SDK** (Chosen): Clean architecture, strong transactions, built-in task queues and schedulers, single container deployment on GCP Cloud Run.

## Decision Outcome

Chosen option: **PostgreSQL + GORM with Candi Framework in Go (`/api`)**. This ADR formally supersedes ADR 0002.

### Positive Consequences

- All domain entities (`Claim`, `Policy`, `ClaimsOfficer`, `ClaimEvent`, `Assignment`, `AssessmentRecommendation`) have clear relational schemas and GORM models.
- Business processes are strictly isolated in Usecases; atomic updates to Claim status and Claim Events are guaranteed with database transactions.
- Async document analysis and SLA tick processing are handled reliably via Candi's Redis Task Queue and Cron Scheduler without requiring external webhook orchestration.
- GCS is used for blob storage (uploaded files and generated PDF reports), storing object URIs in PostgreSQL.

### Negative Consequences / Trade-offs

- Requires running PostgreSQL and Redis instances (or Cloud SQL and Memorystore on GCP / Docker compose locally).

## Links

- Supersedes: [ADR 0002 (Firestore Operational State)](./0002-firestore-operational-state.md)
- Architecture RFC: [RFC 0001](../rfc/0001-agentic-claims-ops-gcp.md)
- Domain Glossary: [CONTEXT.md](../../CONTEXT.md)
