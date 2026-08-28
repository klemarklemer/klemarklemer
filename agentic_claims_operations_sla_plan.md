# Agentic Claims Operations & SLA Platform

## 1. Executive Summary

This proposal describes an **Agentic Claims Operations & SLA Platform** for a General Insurance Company.

The primary objective is to reduce claim processing time while improving:

- SLA compliance
- Workload distribution
- Document completeness
- Claim processing consistency
- Fraud detection and investigation workflow
- Employee performance visibility
- Management visibility into operational bottlenecks

The platform should not initially attempt to replace human claims decision-makers. Instead, AI should automate repetitive work, recommend actions, prioritize work, and surface evidence while keeping high-risk or legally consequential decisions subject to controlled human approval.

---

# 2. Current Claims Process

The current process can be represented as:

```text
1. Customer creates a claim
          |
          v
2. Customer submits required documents
          |
          v
3. Insurance company receives documents
          |
          v
4. Company assigns an employee
          |
          v
5. Survey may be required
   - Property inspection
   - Vehicle inspection
   - Damage verification
          |
          v
6. Employee reviews everything
          |
          +----------------------+
          |                      |
          v                      v
     Incomplete             Not Covered
          |                      |
          v                      v
      Reprocess                Reject
          |
          +----------------------+
          |
          v
      Potential Fraud
          |
          v
       Investigation
          |
          v
        Decision
          |
          +----------------------+
          |          |           |
          v          v           v
       Accept     Reject      Reprocess
```

This process creates several potential bottlenecks:

- Manual document checking
- Manual assignment
- Waiting for surveys
- Lack of proactive SLA reminders
- Uneven employee workloads
- Rework caused by incomplete documents
- Manual reporting
- Limited visibility into why claims are delayed
- Difficulty measuring employee performance consistently

---

# 3. Target State

The target architecture introduces an intelligent orchestration layer around the existing claims operation.

```text
Customer
   |
   v
Claim Creation
   |
   v
Document Submission
   |
   v
+------------------------------+
| Claim Intake Agent           |
| - Validate documents         |
| - Extract information        |
| - Detect missing documents   |
| - Classify claim             |
+--------------+---------------+
               |
               v
        Claim Classification
        +-------+-------+
        |       |       |
        v       v       v
     Simple  Survey   Complex
        |       |       |
        |       v       |
        |   Survey Team |
        |       |       |
        +-------+-------+
                |
                v
       Assignment Agent
                |
                v
       Employee / Team Queue
                |
                v
        Review & Assessment
                |
        +-------+--------+----------------+
        |       |        |                |
        v       v        v                v
   Incomplete  Not      Fraud          Accept
              Covered   Signal
        |       |        |
        v       v        v
    Reprocess Reject  Investigation
                                |
                                v
                        Human Decision
                                |
                                v
                         Final Decision
                                |
                                v
                        Customer Notification
                                |
                                v
                           Claim Closed
```

Every stage has an SLA clock and operational ownership.

---

# 4. Core Business Objectives

## 4.1 Reduce Claim Processing Time

Measure the full lifecycle:

```text
Claim Created
      |
      v
Claim Closed
```

while also measuring each individual stage.

## 4.2 Improve SLA Compliance

Track:

- SLA compliance rate
- SLA breach rate
- Claims at risk
- Average processing time
- Average time per claim stage

## 4.3 Improve Assignment

Automatically determine the most suitable employee based on:

- Skill
- Claim type experience
- Current workload
- SLA performance
- Accuracy
- Availability
- Claim complexity

## 4.4 Reduce Rework

Identify incomplete or inconsistent submissions as early as possible.

## 4.5 Improve Management Visibility

Management should be able to answer:

> Where are claims getting stuck?

> Which teams are overloaded?

> Which employees are consistently exceeding SLA?

> What type of claim generates the most rework?

---

# 5. High-Level Google Cloud Architecture

```text
                    +----------------------+
                    | Customer / Employee  |
                    |      Portal          |
                    +----------+-----------+
                               |
                               v
                    +----------------------+
                    | Cloud Run            |
                    | Claims API / Backend |
                    +----------+-----------+
                               |
             +-----------------+------------------+
             |                 |                  |
             v                 v                  v
      +------------+    +-------------+    +-------------+
      | Cloud      |    | Firestore   |    | Cloud Tasks |
      | Storage    |    |             |    |             |
      | Documents  |    | Claim State |    | Async Jobs  |
      +------------+    +-------------+    +------+------+
                                                 |
                                                 v
                                         +---------------+
                                         | Cloud Run     |
                                         | Workers       |
                                         +-------+-------+
                                                 |
                                                 v
                                         +---------------+
                                         | Google ADK    |
                                         | Agent Runtime |
                                         +-------+-------+
                                                 |
            +----------------+-------------------+----------------+
            |                |                   |                |
            v                v                   v                v
       Intake Agent   Assignment Agent     Survey Agent    Assessment Agent
            |                |                   |                |
            +----------------+-------------------+----------------+
                                    |
                                    v
                            +----------------+
                            | Gemini 3.5     |
                            | Flash          |
                            +----------------+

                         +------------------+
                         | SLA Agent        |
                         +------------------+
                                    |
                                    v
                         Cloud Tasks / Alerts

                         +------------------+
                         | Cloud Workflows  |
                         | Claim Lifecycle  |
                         +------------------+

                         +------------------+
                         | ReportLab        |
                         | PDF Generation   |
                         +------------------+

                         +-----------------------------+
                         | Cloud Logging + OpenTelemetry|
                         | Logs / Traces / Observability|
                         +-----------------------------+
```

---

# 6. Required Technology Responsibilities

## 6.1 Gemini 3.5 Flash

Gemini 3.5 Flash should serve as the main AI reasoning and document-understanding engine.

Potential uses:

### Document understanding

Extract information from:

- Claim forms
- Police reports
- Vehicle documents
- Repair estimates
- Property documents
- Survey reports
- Photos and other supported claim artifacts

Example output:

```json
{
  "claim_id": "CLM-2026-000123",
  "document_status": "INCOMPLETE",
  "missing_documents": [
    "Police Report"
  ],
  "confidence": 0.96
}
```

### Claim classification

Example:

```json
{
  "claim_type": "MOTOR",
  "complexity": "MEDIUM",
  "survey_required": true,
  "confidence": 0.94
}
```

### Assessment support

Gemini can compare evidence from:

```text
Policy
+
Benefit
+
Claim Documents
+
Survey Report
+
Business Rules
+
Claim History
```

and provide an explainable recommendation.

Example:

```json
{
  "coverage": "VALID",
  "fraud_risk": "LOW",
  "recommendation": "APPROVE",
  "confidence": 0.93,
  "reason": [
    "Claim event matches covered benefit",
    "Required documents are complete",
    "Survey confirms reported damage"
  ]
}
```

The AI output should be treated as decision support unless the specific claim category has been explicitly approved for controlled automation.

---

# 7. Google ADK

Google ADK is the agent orchestration layer.

Use multiple specialized agents rather than one large unrestricted agent.

## 7.1 Supervisor / Claims Orchestrator Agent

Responsibilities:

- Understand current claim state
- Decide the next workflow step
- Invoke the appropriate specialized agent
- Coordinate between agents
- Ensure required information is available before proceeding

Core question:

```text
"What should happen next?"
```

## 7.2 Claim Intake Agent

Responsibilities:

- Validate documents
- Extract claim information
- Identify missing documents
- Detect inconsistencies
- Determine claim type
- Determine claim complexity

## 7.3 Assignment Agent

Responsibilities:

- Find eligible employees
- Check workload
- Match employee skill
- Consider employee SLA performance
- Consider claim complexity
- Recommend or perform assignment according to policy

## 7.4 Survey Agent

Responsibilities:

- Determine whether a survey is required
- Create survey task
- Assign surveyor
- Track survey SLA
- Process survey report
- Escalate overdue surveys

## 7.5 Assessment Agent

Responsibilities:

- Compare claim evidence with policy coverage
- Identify possible exclusions
- Summarize evidence
- Identify inconsistencies
- Produce a recommended outcome
- Escalate uncertain or high-risk cases

Possible recommendations:

```text
ACCEPT
REPROCESS
REJECT - NOT COVERED
POTENTIAL FRAUD
MANUAL REVIEW
```

## 7.6 SLA Agent

Responsibilities:

- Monitor deadlines
- Identify claims at risk
- Create reminders
- Escalate overdue claims
- Escalate to team leaders or managers
- Produce SLA operational metrics

---

# 8. Firestore

Firestore should be the operational state store.

Suggested logical collections:

```text
claims/
employees/
teams/
policies/
benefits/
claim_events/
assignments/
surveys/
sla_rules/
performance_scores/
notifications/
ai_assessments/
```

Example `claims` document:

```json
{
  "claim_id": "CLM-2026-000123",
  "policy_id": "POL-873829",
  "customer_id": "CUS-001",
  "claim_type": "MOTOR",
  "status": "ASSESSMENT",
  "priority": "NORMAL",
  "complexity": "MEDIUM",
  "assigned_employee": "EMP-104",
  "assigned_team": "MOTOR-CLAIMS",
  "current_stage": "ASSESSMENT",
  "created_at": "2026-08-28T10:00:00Z",
  "sla_deadline": "2026-09-04T10:00:00Z"
}
```

## Claim Event History

Maintain an immutable operational timeline.

Example:

```text
10:00 Claim submitted
10:01 Documents validated
10:02 Missing document detected
10:30 Customer resubmitted document
10:32 Claim classified
10:35 Assigned to EMP-104
14:00 Survey requested
Next day 11:00 Survey completed
Next day 11:20 Assessment started
Next day 11:32 Recommendation generated
Next day 12:00 Human approval
Next day 12:05 Customer notified
```

This history is essential for auditability and SLA analysis.

---

# 9. Cloud Storage

Cloud Storage should store claim artifacts rather than large document contents directly in Firestore.

Suggested structure:

```text
claims/
  CLM-2026-000123/
    submitted/
      claim-form.pdf
      police-report.pdf
      vehicle-registration.pdf

    survey/
      survey-report.pdf
      photo-001.jpg
      photo-002.jpg

    generated/
      assessment-report.pdf
      decision-letter.pdf
```

Firestore stores the metadata and object references.

---

# 10. Cloud Run

Cloud Run should host the stateless application components.

Potential services:

```text
claims-api
agent-worker
document-worker
notification-worker
report-worker
integration-worker
```

A single Cloud Run service can be used for an initial MVP, while separating workers into services later when scale and operational requirements justify it.

---

# 11. Cloud Tasks

Cloud Tasks should be used for individual asynchronous jobs.

Examples:

```text
Document processing
SLA reminder
Employee reminder
Survey reminder
Manager escalation
Customer notification
Report generation
```

Example:

```text
Claim submitted
      |
      v
Cloud Task
      |
      v
Document Analysis Worker
      |
      v
Gemini
      |
      v
Firestore
```

Cloud Tasks keeps long-running processing out of synchronous API requests.

---

# 12. Cloud Workflows

Cloud Workflows should orchestrate long-running claim processes.

Example:

```text
1. Receive claim
2. Validate policy
3. Validate documents
4. Classify claim
5. Determine survey requirement
6. Assign employee
7. Wait for survey if required
8. Process survey report
9. Run assessment
10. Determine whether human review is required
11. Generate report
12. Final approval
13. Notify customer
14. Close claim
```

Cloud Workflows should coordinate the business process while Cloud Tasks handles individual asynchronous jobs.

---

# 13. SLA Model

SLA should exist at both the **claim level** and the **stage level**.

## Example Claim-Level SLAs

| Claim Type | Target |
|---|---:|
| Simple Claim | 2 business days |
| Standard Claim | 5 business days |
| Survey Claim | 7 business days |
| Complex Claim | 14 business days |

## Example Stage-Level SLAs

| Stage | Example Target |
|---|---:|
| Document Verification | 4 hours |
| Assignment | 2 hours |
| Survey Assignment | 4 hours |
| Survey Completion | 2 days |
| Survey Report Processing | 1 day |
| Assessment | 1 day |
| Final Approval | 4 hours |
| Customer Notification | 2 hours |

The actual values should be configured from the company's approved SLA policy rather than hard-coded.

---

# 14. Automated SLA Reminder Strategy

Example:

```text
T-24h  -> Employee Reminder
T-12h  -> Strong Reminder
T-4h   -> Team Lead Escalation
T-1h   -> Manager Escalation
T+0    -> SLA Breach
T+1h   -> Operations Escalation
```

Example employee notification:

```text
CLAIM SLA ALERT

Claim: CLM-2026-000123
Stage: Document Review
Assigned to: EMP-104

Remaining SLA: 3h 12m

Status: AT RISK
Required action: Complete document review
```

---

# 15. Automatic Assignment System

The Assignment Agent should not simply select the employee with the highest performance score.

A better assignment model balances:

```text
Assignment Score =
    Skill Match
  + Claim Type Experience
  + Current Workload
  + SLA Performance
  + Accuracy
  + Complexity Capacity
  + Availability
```

Example:

| Employee | Skill | Workload | SLA | Accuracy | Result |
|---|---:|---:|---:|---:|---|
| A | 95 | 80 | 98 | 96 | Medium |
| B | 90 | 30 | 94 | 92 | High |
| C | 97 | 90 | 99 | 98 | Low |

Employee C has excellent performance but is overloaded, so Employee B may be the better assignment.

This prevents the highest performers from becoming operational bottlenecks.

---

# 16. Employee Performance Tiering

The company can introduce performance tiers such as:

```text
A - Expert
B - Strong
C - Developing
D - Improvement Required
```

## Tier A

Potential characteristics:

```text
SLA compliance >= 95%
High decision accuracy
Low rework
Low escalation
Strong compliance
```

May handle:

- Complex claims
- High-value claims
- Difficult cases
- Fraud investigations

## Tier B

Handles:

- Normal claims
- Moderate complexity
- Standard survey cases

## Tier C

Handles:

- Simple claims
- Standard workflows
- Cases with additional review where required

## Tier D

Used primarily as an operational improvement indicator.

Potential responses:

- Coaching
- Training
- Mentoring
- Reduced complexity assignment
- Increased QA review

The system should avoid using AI-generated scores as the sole basis for employee disciplinary decisions.

---

# 17. Performance Score

A transparent scoring model is recommended.

Example:

```text
35% SLA Performance
25% Decision Accuracy
15% Rework Rate
10% Productivity
10% Quality
 5% Compliance
```

Example:

```text
Employee: EMP-104

SLA:             96%
Accuracy:        98%
Rework:           3%
Productivity:    91%
Quality:         95%
Compliance:     100%

Overall Score: 96.4
Tier: A
```

The score can be used as an operational KPI and, subject to HR governance, feed into OKR/performance assessment systems.

---

# 18. Fraud Handling

Fraud detection should be treated separately from ordinary claim assessment.

The system should not simply allow an LLM to declare:

```text
"This claim is fraud."
```

Instead:

```text
Potential Fraud Detection
          |
          v
Fraud Signals
          |
          v
Risk Score + Evidence
          |
          v
Fraud Investigator
          |
          v
Human Decision
```

Example:

```json
{
  "risk_score": 0.87,
  "signals": [
    "Repeated similar claims",
    "Inconsistent accident timeline",
    "Document inconsistency"
  ],
  "requires_investigation": true
}
```

All fraud recommendations should be explainable and traceable to supporting evidence.

---

# 19. Human-in-the-Loop Model

A key principle is:

> The more consequential the decision, the stronger the human control should be.

Suggested model:

| Claim Category | AI Recommendation | Human Review |
|---|---|---|
| Low-risk, simple, clearly covered | Yes | Optional / controlled automation |
| Standard claim | Yes | Yes |
| Policy ambiguity | Yes | Mandatory |
| High-value claim | Yes | Mandatory |
| Potential fraud | Yes | Mandatory |
| Regulatory/legal sensitivity | Yes | Mandatory |

This allows the company to automate progressively without giving an AI model unrestricted authority over claims decisions.

---

# 20. ReportLab

ReportLab should be used for Python-based PDF generation.

Possible outputs:

```text
Claim Assessment Report
Survey Report
Internal Investigation Report
Decision Report
Customer Decision Letter
Management SLA Report
Employee Performance Report
```

Example:

```text
CLAIM ASSESSMENT REPORT

Claim ID: CLM-2026-00123
Policy: POL-873829
Claim Type: Motor

Documents:
[OK] Police Report
[OK] Driver License
[OK] Vehicle Registration
[OK] Repair Estimate

Survey:
[OK] Completed

Coverage: VALID
Fraud Risk: LOW

AI Recommendation:
APPROVE

Human Approval:
REQUIRED

Approved By:
Claims Officer
```

Generated PDFs should be stored in Cloud Storage with their metadata recorded in Firestore.

---

# 21. Cloud Logging + OpenTelemetry

Observability is critical for an agentic insurance platform.

The company should be able to answer:

```text
Who performed the action?
When did it happen?
Which agent performed it?
Which model was used?
Which version was used?
What information was considered?
What recommendation was produced?
Who approved the final decision?
How long did each step take?
```

## Example trace

```text
Claim API
   |
   +-- Intake Agent
   |      |
   |      +-- Gemini
   |
   +-- Assignment Agent
   |
   +-- Survey Workflow
   |
   +-- Assessment Agent
          |
          +-- Gemini
```

## Suggested telemetry fields

```text
claim_id
trace_id
span_id
agent_name
employee_id
team_id
action
timestamp
model_name
model_version
prompt_version
confidence
reason
previous_state
new_state
```

OpenTelemetry should provide distributed traces across:

- API
- Workers
- Agents
- Gemini calls
- Workflows
- External services

Cloud Logging should retain structured logs for operational debugging and audit support.

---

# 22. Sample End-to-End Claim

Consider a motor accident claim.

## 10:00

Customer submits claim.

```text
Status = RECEIVED
```

## 10:01

Intake Agent checks documents.

```text
[OK] Insurance policy
[OK] Driver license
[OK] Vehicle registration
[MISSING] Police report
```

The system automatically requests the missing document.

## 10:30

Customer submits the missing document.

Claim classification:

```text
Type: MOTOR
Complexity: MEDIUM
Survey: REQUIRED
SLA: 7 days
```

## 10:35

Assignment Agent scores available employees.

```text
Employee A -> 82
Employee B -> 94
Employee C -> 71
```

Claim assigned to Employee B.

## Day 1

Survey task created and assigned.

## Day 2

Survey completed.

Survey report and photos are uploaded to Cloud Storage.

## Day 2

Assessment Agent evaluates:

```text
Policy
+
Claim
+
Documents
+
Survey Report
```

Output:

```text
Coverage: VALID
Fraud Risk: LOW
Estimated Loss: IDR XX
Recommendation: APPROVE
```

Human claims officer reviews the recommendation.

## Day 2

Decision:

```text
APPROVED
```

ReportLab generates the assessment report.

Customer receives notification.

Claim is closed.

---

# 23. Management Dashboard

The management dashboard should provide a real-time operational view.

Example:

```text
CLAIMS OPERATIONS

Total Claims              12,842
Open Claims                 1,427
SLA Compliance               94.7%
Average Processing Time      2.8 days
SLA Breaches                    76
Claims At Risk                  143

----------------------------------------

TEAM PERFORMANCE

Team              SLA      Avg Time
Motor             96%      2.1d
Property          91%      4.2d
Health            97%      1.8d
Travel            98%      1.2d

----------------------------------------

BOTTLENECKS

Survey             HIGH
Document Review    MEDIUM
Assessment         LOW
Approval           LOW
```

---

# 24. Recommended KPI / OKR Metrics

The initial management KPI set should include:

## Claims Operations

- SLA Compliance Rate
- SLA Breach Rate
- Average Claim Processing Time
- Median Claim Processing Time
- Average Stage Processing Time
- Claims at Risk
- Claims Pending Assignment

## Quality

- First-Time-Right Rate
- Rework Rate
- Decision Accuracy
- Escalation Rate

## Productivity

- Claims per Employee
- Claims Closed per Day
- Average Workload per Employee

## AI

- AI Recommendation Acceptance Rate
- AI Recommendation Override Rate
- AI Confidence Distribution
- AI Processing Time
- Document Automation Rate

## Fraud

- Fraud Alerts
- Fraud Investigation Rate
- Confirmed Fraud Rate
- False Positive Rate

---

# 25. Implementation Roadmap

## Phase 1 — Claims Command Center

Build:

```text
Cloud Run
Firestore
Cloud Storage
Cloud Tasks
Cloud Logging
OpenTelemetry
```

Focus on:

```text
Claim tracking
Claim ownership
Status tracking
SLA tracking
Audit trail
Operational dashboard
```

The first objective is visibility, not AI.

---

## Phase 2 — Intelligent Operations

Add:

```text
Google ADK
Gemini 3.5 Flash
Document analysis
Claim classification
Assignment Agent
SLA Agent
Automated reminders
```

Focus on reducing administrative work.

---

## Phase 3 — Decision Support

Add:

```text
Policy matching
Coverage analysis
Survey analysis
Fraud signals
Claim assessment recommendation
```

Keep humans in the loop for consequential decisions.

---

## Phase 4 — Performance Management

Add:

```text
Employee scoring
A/B/C/D performance tiers
Team performance
Operational OKRs
Training recommendations
```

Use transparent metrics and governance.

---

## Phase 5 — Controlled Automation

Introduce straight-through processing for carefully defined low-risk cases.

Example:

```text
Low Value
+
Complete Documentation
+
Clearly Covered
+
Low Fraud Risk
+
No Survey Required
+
Approved Business Rules
        |
        v
Controlled Automated Processing
```

Complex or high-risk claims continue through mandatory human review.

---

# 26. Recommended Architecture Evolution

The platform should evolve progressively.

```text
Phase 1
Operational Visibility
        |
        v
Phase 2
SLA Automation
        |
        v
Phase 3
Assignment Automation
        |
        v
Phase 4
Document Intelligence
        |
        v
Phase 5
Assessment Recommendation
        |
        v
Phase 6
Controlled Straight-Through Processing
```

This approach reduces implementation risk and allows measurable benefits at every phase.

---

# 27. Governance Principles

The platform should follow these principles:

## Explainability

Every important AI recommendation should have:

- Evidence
- Reason
- Confidence
- Model/version information
- Human decision status

## Auditability

Every material claim state change should be traceable.

## Human Oversight

High-value, ambiguous, fraud-related, and legally sensitive claims should require human review.

## Separation of Concerns

Use deterministic business rules for rules that must be exact, and use Gemini for unstructured document understanding and reasoning support.

For example:

```text
Business Rule Engine
      +
Gemini
      |
      v
Assessment Recommendation
```

rather than asking Gemini to independently enforce every insurance rule.

## No Hidden Employee Scoring

Employee performance calculations should use documented metrics and should be reviewable by management and HR.

---

# 28. Key Business Outcome

The platform should ultimately transform claims operations from:

```text
Customer
   |
Manual Process
   |
Manual Assignment
   |
Manual Follow-up
   |
Manual Assessment
   |
Manual Reporting
```

into:

```text
Customer
   |
Automated Intake
   |
Intelligent Assignment
   |
SLA Monitoring
   |
AI-Assisted Assessment
   |
Human Approval
   |
Automated Reporting
   |
Continuous Performance Measurement
```

The primary business goal is:

> **Reduce claim processing time while improving SLA compliance, workload distribution, decision consistency, and management visibility.**

---

# 29. Final Reference Architecture

```text
                    +----------------------+
                    | Customer / Employee  |
                    |      Portal          |
                    +----------+-----------+
                               |
                               v
                    +----------------------+
                    | Cloud Run            |
                    | Claims API / Backend |
                    +----------+-----------+
                               |
             +-----------------+------------------+
             |                 |                  |
             v                 v                  v
      +------------+    +-------------+    +-------------+
      | Cloud      |    | Firestore   |    | Cloud Tasks |
      | Storage    |    | Claim State |    | Async Jobs  |
      | Documents  |    | + Events    |    | + Reminders |
      +------------+    +-------------+    +------+------+
                                                 |
                                                 v
                                         +---------------+
                                         | Cloud Run     |
                                         | Workers       |
                                         +-------+-------+
                                                 |
                                                 v
                                         +---------------+
                                         | Google ADK    |
                                         | Agent Runtime |
                                         +-------+-------+
                                                 |
             +----------------+------------------+----------------+
             |                |                   |                |
             v                v                   v                v
        Intake Agent   Assignment Agent     Survey Agent    Assessment Agent
             |                |                   |                |
             +----------------+-------------------+----------------+
                                    |
                                    v
                            +----------------+
                            | Gemini 3.5     |
                            | Flash          |
                            +----------------+
                                    |
                                    v
                           Human Approval /
                           Controlled Action
                                    |
                                    v
                         +------------------+
                         | Cloud Workflows  |
                         | Claim Lifecycle  |
                         +---------+--------+
                                   |
                                   v
                           +---------------+
                           | ReportLab     |
                           | PDF Reports   |
                           +-------+-------+
                                   |
                                   v
                            Cloud Storage

           +------------------------------------------------+
           | Cloud Logging + OpenTelemetry                    |
           | Logs / Traces / Metrics / Agent Observability   |
           +------------------------------------------------+
```

---

# 30. Summary

The recommended platform is not simply an AI chatbot for claims.

It is a **Claims Operations Orchestration Platform** combining:

- Gemini 3.5 Flash for AI reasoning and document understanding
- Google ADK for multi-agent orchestration
- Firestore for operational state
- Cloud Run for APIs and workers
- Cloud Tasks for asynchronous jobs and reminders
- Cloud Workflows for long-running claim processes
- Cloud Storage for claim documents and generated artifacts
- ReportLab for PDF reports
- Cloud Logging + OpenTelemetry for observability and auditability

The platform should first optimize the operational workflow, then progressively introduce AI-assisted decisions and eventually controlled automation for low-risk claims.
