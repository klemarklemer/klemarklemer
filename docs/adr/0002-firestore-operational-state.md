# Firestore is the operational state store for Claims and Claim events

The Claim is a long-lived operational record with an append-only timeline. Firestore matches the vision (document per Claim, subcollections or sibling `claim_events`), stays cheap at hackathon scale, and is a listed Google Cloud infrastructure service for submission rules. Cloud Storage holds blobs; Firestore holds metadata and object refs.

**Status:** superseded by [ADR 0004](./0004-postgres-candi-state.md)

## Considered Options

- Cloud SQL / Postgres — stronger relational reporting later; more ops for a weekend tracer bullet.
- Keep all state in ADK session memory — not auditable, not an SLA clock source of truth.
- Firestore (chosen) — Claim + Claim events as the system of record for the demo and for OTel correlation by `claim_id`.

## Consequences

Query patterns stay simple (get Claim, list events by `claim_id`). Management dashboards and SQL analytics are Phase 2+. Emulator or a real project both work for local spin-up; production demo uses a real project.
