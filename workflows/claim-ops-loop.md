# Workflow: Claims Operations & Human Approval Loop

## 1. Overview & Objective

This workflow coordinates an end-to-end motor claim lifecycle from intake to close. Repetitive operational tasks (document completeness checking, assignment scoring, SLA evaluation, policy-vs-loss assessment) are executed autonomously by specialist agents. Legally binding decisions are deferred to a "push right" human checkpoint where a named claims officer reviews a decision-ready brief and issues a binding **Decision**.

---

## 2. Shared Vocabulary & Domain Entities

- **Claim**: The notified loss under review (`INTAKE` -> `DOCUMENT_VERIFICATION` -> `ASSIGNMENT` -> `ASSESSMENT` -> `DECISION` -> `CLOSED`).
- **Policy**: The in-force insurance contract governing coverage and deductibles.
- **Document Completeness**: `INCOMPLETE` (missing police report) or `COMPLETE`.
- **Assignment**: Binding between a Claim and an eligible Claims Officer based on deterministic workload/skill scoring.
- **Assessment Recommendation**: Proposed outcome (`APPROVE`, `REJECT`, `MANUAL_REVIEW`) generated with Gemini reasoning.
- **Human Approval Checkpoint**: The decision checkpoint where a Claims Officer confirms or overrides the recommendation.
- **Decision**: The binding company outcome that closes the Claim.
- **Claim Event**: Immutable event log appended to the Claim's timeline.

---

## 3. Operational Loops

```mermaid
flowchart TD
    subgraph Loop 1: Intake & Classification (Autonomous)
        T1[Trigger: Claim Created / Document Uploaded] --> CheckDocs[Intake Agent: Extract & Verify Documents]
        CheckDocs --> CompleteCheck{All required docs present?}
        CompleteCheck -->|No| SetIncomplete[Set Completeness = INCOMPLETE\nEmit In-App Missing Doc Alert]
        CompleteCheck -->|Yes| SetComplete[Set Completeness = COMPLETE\nClassify MOTOR/MEDIUM\nSurvey Required = false]
    end

    subgraph Loop 2: Assignment (Autonomous)
        SetComplete --> AssignTrigger[Trigger: Document Completeness = COMPLETE]
        AssignTrigger --> ScoreOfficers[Assignment Agent: Calculate Workload vs Skill Score]
        ScoreOfficers --> BindOfficer[Write Binding Assignment\nAppend ASSIGNED Claim Event]
    end

    subgraph Loop 3: Assessment Reasoning (Autonomous)
        BindOfficer --> AssessTrigger[Trigger: Assignment Written]
        AssessTrigger --> RunAssessment[Assessment Agent: Compare Policy vs Claim Details via Gemini]
        RunAssessment --> WriteRec[Write Assessment Recommendation\nOutcome: APPROVE\nReasons + Confidence]
    end

    subgraph Loop 4: Human Approval Checkpoint (Push-Right Checkpoint)
        WriteRec --> Brief[Present Decision-Ready Brief to Claims Officer]
        Brief --> HumanAction{Claims Officer Action}
        HumanAction -->|Approve| MakeDecision[Record Human Approval\nCreate Binding DECISION\nClose Claim\nGenerate PDF]
        HumanAction -->|Reject/Override| ManualRework[Record Rejection Event\nKeep Claim Open for Rework]
    end
```

---

## 4. Checkpoint Details: Claims Officer Approval

- **Type**: Push-Right Checkpoint (Mandatory).
- **Actor**: Named Claims Officer.
- **Brief Delivered to Officer**:
  - Claim Number & Policy details (holder, vehicle plate, coverage).
  - Document completeness verification summary.
  - Assessment recommendation (`APPROVE`), confidence score, and clear narrative justification.
  - Estimated loss vs deductible breakdown.
- **Action**: Explicit approval checkbox & button: *"Creates a Decision and closes this Claim."*
- **Outcome**: Transitions Stage to `CLOSED`, records immutable `HUMAN_APPROVAL` Claim event, and stores assessment report.
