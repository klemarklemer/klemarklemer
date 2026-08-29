# Taskmaster Claims Operations API Documentation

## Table of Contents

1. [01 - Overview & Quickstart](./01-overview-and-quickstart.md)
   - System Overview & Autonomous Agent Loops
   - Server Base URLs
   - 4-Minute Demo Tracer Bullet cURL Flow
2. [02 - Claims Operations API](./02-claims-operations-api.md)
   - Claim Ingestion (`POST /v1/claim`)
   - Claim Listing & Filtering (`GET /v1/claim`)
   - Claim Detail Aggregate (`GET /v1/claim/:id`)
   - Document Upload & Cascading Execution (`POST /v1/claim/:id/documents`)
   - Manual Specialist Triggers (`intake`, `assignment`, `assessment`)
   - Human Approval & Binding Decision Gate (`POST /v1/claim/:id/approval`)
   - Demo Reset Utility (`POST /v1/demo/reset`)
3. [03 - Officers & Policies API](./03-officers-and-policies-api.md)
   - Claims Officers & Workload Allocation (`GET /v1/officer`)
   - Active Policy Contracts (`GET /v1/policy`)
4. [04 - Error Handling & Response Envelopes](./04-error-handling-and-envelopes.md)
   - Candi Global JSON Envelope
   - HTTP Status Codes Matrix
   - Validation & Transaction Error Formats

---

## Machine-Readable Specifications

- **OpenAPI 3.1 Specification (YAML):** [`openapi.yaml`](./openapi.yaml)
- **OpenAPI 3.1 Specification (JSON):** [`openapi.json`](./openapi.json)
- **Postman Collection v2.1:** [`postman_collection.json`](./postman_collection.json)
