# Taskmaster Claims Operations

An agentic claims operations and SLA platform for general insurance, built for the
All Things Agentic Hackathon (The Taskmaster track).

Motor insurance claims stall in predictable places: documents arrive incomplete,
assignment is manual and uneven, assessment waits on a human to read a policy, and
nobody can see why a claim is late. This platform puts four autonomous loops around
that process. They handle the repetitive work and surface evidence — but the loop
that costs money stops at a human.

## The four loops

| Loop | Agent | What it does |
|---|---|---|
| 1 · Intake | `IntakeAgent` | Names the documents actually missing, then has Gemini judge severity and whether a surveyor must inspect |
| 2 · Assignment | `AssignmentAgent` | Scores officers on workload and skill, assigns the best fit |
| 2b · Survey | `AssignmentAgent` | When Loop 1 calls for an inspection, routes the claim to a free surveyor of the fitting specialty |
| 3 · Assessment | `AssessmentAgent` | Weighs policy coverage against estimated loss and evidence, recommends an outcome |
| 4 · Decision gate | Human | A claims officer binds the decision; only this closes a claim |

Loops 1–3 chain autonomously. Loop 4 never does — an `APPROVE` from the Assessment
Agent is a proposal, not a settlement.

Every transition writes an immutable `claim_events` row naming the actor, the stage
change, and the reasoning, so a claim's whole history is replayable.

## Architecture

```
React SPA (apps/web)  ──HTTP──►  Candi Go service (apps/api/services/core)
                                       │
                          ┌────────────┼────────────┐
                          ▼            ▼            ▼
                     PostgreSQL      Redis      Gemini 3.5 Flash
                     claims,         locks      via Vertex AI or
                     events,                    the Gemini API
                     policies
```

Clean architecture per module: `domain` → `repository` → `usecase` → `delivery`.
The reasoning behind Loop 3 sits behind an `Assessor` interface in
`services/core/pkg/shared/gemini`, so the model is swappable and testable.

## Quick start

**Prerequisites:** Go 1.26+, Docker, Node 20+, and the Candi CLI:

```bash
go install github.com/golangid/candi/cmd/candi@latest
```

**1 · Start Postgres and Redis**

```bash
make deps
```

**2 · Configure**

```bash
cp apps/api/services/core/.env.sample apps/api/services/core/.env
```

Set an assessment backend in that file (see [Configuration](#configuration)). With
neither set the service still runs, on deterministic underwriting rules.

**3 · Migrate and seed**

```bash
make migrate
```

**4 · Run**

```bash
make up          # API on :8000, web console on :5173
```

Or separately: `make api` and `make web`.

**5 · Confirm which engine is live**

The service prints its assessment backend at startup:

```
claim usecase: assessment agent engine -> gemini-3.5-flash via vertex-ai
```

If that says `deterministic-rules`, your credential did not load.

## Configuration

Loop 3 picks its backend from the environment alone — the same binary runs all three ways.

**Vertex AI** — calls appear in Google Cloud Console, which doubles as deployment proof:

```bash
gcloud auth application-default login
gcloud services enable aiplatform.googleapis.com --project klemarklemer
```

```env
GOOGLE_GENAI_USE_VERTEXAI=true
GOOGLE_CLOUD_PROJECT=klemarklemer
GOOGLE_CLOUD_LOCATION=asia-southeast1
```

**Gemini API** — simplest for local work. Get a key at https://ai.google.dev:

```env
GEMINI_API_KEY=your-key
```

**Neither** — deterministic underwriting rules. The service boots and still reasons
over the real claim record; it just does so without a model. This keeps CI and a
fresh clone honest rather than imitating one.

`GEMINI_MODEL` overrides the default `gemini-3.5-flash`.

## API

Base URL `http://localhost:8000/v1`.

| Method | Path | |
|---|---|---|
| `GET` | `/claim/` | List claims |
| `GET` | `/claim/:id` | Claim detail with events, documents, recommendation |
| `POST` | `/claim/` | Create a claim |
| `POST` | `/claim/:id/documents` | Attach a document |
| `POST` | `/claim/:id/intake` | Run Loop 1 (chains into 2 and 3) |
| `POST` | `/claim/:id/assignment` | Run Loop 2 |
| `POST` | `/claim/:id/assessment` | Run Loop 3 |
| `POST` | `/claim/:id/approval` | Loop 4 — bind the human decision |
| `POST` | `/demo/reset` | Reset to seed state |

Full OpenAPI 3.1 spec and a Postman collection are in [`docs/api/`](docs/api/).

### Seeing it reason

Create a claim whose estimated loss exceeds the policy ceiling and run assessment.
The agent returns `MANUAL_REVIEW` citing the actual ceiling and loss — not an
approval. Lower the loss below the deductible and it returns `REJECT`, because
nothing is payable. The recommendation follows the record, not a script.

## Deploying to Google Cloud

The API needs a reachable Postgres and Redis; it initialises both at boot.

```bash
gcloud run deploy klemarklemer-api \
  --source apps/api \
  --region asia-southeast1 \
  --project klemarklemer \
  --set-env-vars GOOGLE_GENAI_USE_VERTEXAI=true,GOOGLE_CLOUD_PROJECT=klemarklemer
```

Supply `SQL_DB_READ_DSN`, `SQL_DB_WRITE_DSN`, `REDIS_READ_DSN`, and
`REDIS_WRITE_DSN` pointing at managed instances — Cloud SQL for Postgres, and
Memorystore or any TLS-reachable Redis.

Tagged pushes also trigger `.github/workflows/cd.yml`, which requires the
`GCP_SA_KEY`, `GCP_PROJECT_ID`, and `GCP_REGION` repository secrets.

## Tests

```bash
make test-api      # Go build and unit tests
make test-web      # typecheck and build
make test          # both
```

## Documentation

| | |
|---|---|
| [`docs/architecture/`](docs/architecture/) | C4 model, ERD, state machines, SLA formulas |
| [`docs/architecture/architecture.html`](docs/architecture/architecture.html) | Visual system topology, the four loops, request lifecycle |
| [`docs/security/`](docs/security/) | STRIDE threat model, RBAC, PII handling, SSDLC |
| [`docs/api/`](docs/api/) | API catalog, OpenAPI 3.1, Postman collection |
| [`docs/adr/`](docs/adr/) | Architectural decision records |
| [`agentic_claims_operations_sla_plan.md`](agentic_claims_operations_sla_plan.md) | The originating product plan |

## Layout

```
apps/api/    Candi Go monorepo — services/core holds the claims service
apps/web/    React + Vite claims officer console
deployments/ Docker Compose for local Postgres and Redis
docs/        Architecture, security, API, and decision records
```
