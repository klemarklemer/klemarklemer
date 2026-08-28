# Claims Operations

The language of a general-insurance claims operation: a Claim moves through Stages under SLA clocks, with Assignments, Assessment recommendations, Human approval, and an immutable Claim event history. This glossary is the domain; it is not a product spec or a cloud architecture.

## Language

**Claim**:
A notified loss against a Policy that the company must process from intake to close.
_Avoid_: Case, ticket, file, FNOL (as a synonym for the whole Claim)

**Policy**:
The in-force contract that defines what a Claim may be covered for.
_Avoid_: Cover note, product (when you mean the contract)

**Stage**:
A named step on the Claim path (for example document verification, assignment, assessment, final approval). A Claim is in exactly one current Stage at a time.
_Avoid_: Status (when you mean the step), phase, workflow node

**SLA clock**:
A deadline attached either to the whole Claim (claim-level) or to the current Stage (stage-level). Both clocks can exist on the same Claim.
_Avoid_: Timer, timeout, due date (use those only as informal speech)

**Assignment**:
The binding of a Claim to a claims officer (and optionally a team) for operational ownership of the current work.
_Avoid_: Routing, allocation, dispatch (when you mean the binding itself)

**Assessment recommendation**:
An explainable proposed outcome for a Claim (for example approve, reject, reprocess, potential fraud, manual review), with evidence, reasons, and confidence. It is not a Decision.
_Avoid_: AI decision, verdict, determination (when you mean the proposal)

**Decision**:
The company's binding outcome on a Claim after Human approval (or, later, after an explicitly authorised straight-through rule). Only a Decision closes the Claim.
_Avoid_: Recommendation, assessment (when you mean the binding outcome)

**Human approval**:
A named claims officer's recorded acceptance or rejection of an Assessment recommendation, which produces a Decision.
_Avoid_: Sign-off, rubber stamp, HITL (in glossary text)

**Claim event**:
An immutable fact that something material happened to a Claim (who/what, when, previous Stage, new Stage). The sequence of Claim events is the operational timeline.
_Avoid_: Log, audit row, history item (when you mean the domain fact)

**Claims officer**:
The human employee who owns operational work on an assigned Claim and who performs Human approval.
_Avoid_: Adjuster (unless the company uses that title), user, agent (humans are not agents)

**Document completeness**:
Whether the Claim has every artifact required for the current Stage (for example a police report on a motor Claim). Incomplete is a Claim fact, not a chat message.
_Avoid_: Missing files (as the concept name)

**Survey**:
A physical or specialist inspection of damage when the Claim requires it. The hackathon demo Claim is seeded so a Survey is not required; the concept still exists in the domain.
_Avoid_: Inspection visit (as the concept name)
