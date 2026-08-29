# GEMINI.md - Agent Instructional Context

This file serves as the foundational instruction manual and operational context for Gemini CLI and other AI agents working within this workspace. It maps the project architecture, tech stack, domain language, UI/UX mandates, and development guidelines.

---

## 1. Project Overview

This repository hosts **Taskmaster Claims Operations**, an Agentic Claims Operations & SLA Platform designed for the *All Things Agentic Hackathon — The Taskmaster* (August 2026). 

The platform optimizes general-insurance claims handling by orchestrating an intelligent multi-agent layer around claims operations. Rather than replacing human claims decision-makers, the platform automates repetitive chores (document checking, assignment score calculation, SLA reminders, assessment analysis) while ensuring high-risk or legally consequential decisions require explicit **Human approval** before becoming a binding **Decision**.

### Key Technology Stack
*   **Frontend:** React (v18), TypeScript, Vite, Tailwind CSS (v4), Phosphor Icons, and Vitest for testing.
*   **Backend (Planned/RFC):** Cloud Run service, Firestore (state store), Cloud Storage (document blobs), Cloud Tasks (async jobs), Google ADK (Agent Development Kit) for multi-agent choreography, and Gemini 3.5 Flash for unstructured extraction and assessment reasoning.

---

## 2. Directory Structure

*   `/web/`: Complete React/TS single-page Claims Workspace console.
*   `/api/`: Backend directory placeholder (contains only `.gitkeep`; to be populated per RFC 0001).
*   `/docs/`: Design and architectural system of record.
    *   `/docs/prd/`: Product Requirements Document (`taskmaster-claims-ops.md`).
    *   `/docs/rfc/`: Technical Proposals (`0001-agentic-claims-ops-gcp.md`).
    *   `/docs/adr/`: Architectural Decision Records (ADRs 0001, 0002, 0003).
    *   `/docs/ui/`: Operations UI/UX guidelines (`claims-ops-ux.md`).
*   `CONTEXT.md`: The official Domain Glossary defining the strict language of the insurance operations.
*   `AGENTS.md`: Crucial pointer file referencing UX guidelines and glossary mappings.
*   `agentic_claims_operations_sla_plan.md`: The broader, long-term enterprise platform proposal.

---

## 3. Strict Domain Language & Glossary (Mandatory)

To prevent terminology drift in the codebase, tests, UI, and prompts, agents **MUST** strictly adhere to the glossary defined in `CONTEXT.md`. Do not use terms from the "Avoid" lists:

| Domain Term | Definition | Avoid (Do Not Use) |
|---|---|---|
| **Claim** | A notified loss against a Policy that must be processed from intake to close. | Case, ticket, file, FNOL (as synonym for Claim) |
| **Policy** | The in-force contract defining what a Claim may be covered for. | Cover note, product (when meaning contract) |
| **Stage** | A named step on the Claim path (e.g., Intake, Assessment, Decision, Closed). | Status, phase, workflow node |
| **SLA clock** | A deadline attached to the whole Claim or to the current Stage. | Timer, timeout, due date |
| **Assignment** | The binding of a Claim to a claims officer for operational ownership. | Routing, allocation, dispatch |
| **Assessment recommendation** | An explainable proposed outcome (e.g., APPROVE, REJECT, MANUAL_REVIEW) with evidence and confidence. It is *not* a Decision. | AI decision, verdict, determination |
| **Decision** | The company's binding outcome on a Claim after Human approval. Only a Decision closes a Claim. | Recommendation, assessment (when meaning binding outcome) |
| **Human approval** | A named claims officer's recorded acceptance/rejection of a recommendation. | Sign-off, rubber stamp, HITL (in glossary text) |
| **Claim event** | An immutable fact that something material happened to a Claim (appended to timeline). | Log, audit row, history item |
| **Claims officer** | The human employee who owns operational work and performs Human approval. | Adjuster, user, agent (humans are not agents) |
| **Document completeness** | Whether the Claim has every artifact required for the current Stage. | Missing files |
| **Survey** | A physical or specialist inspection of damage. | Inspection visit |

---

## 4. UI/UX & Visual System Mandates (from `docs/ui/claims-ops-ux.md`)

When modifying or designing the `/web/` frontend console, follow these rigid UI contracts. This is an **operations console**, not a marketing landing page or portfolio.

### UI Dials & Constraints
*   **Locked Dials:** `DESIGN_VARIANCE: 3` · `MOTION_INTENSITY: 2` · `VISUAL_DENSITY: 7`
*   **No "AI Slop":** Ban landing heroes, gradient meshes, glassmorphism navs, purple glows, three equal marketing cards, custom cursors, and "AI copilot" chat docks.
*   **Theme & Accent:** High contrast neutrals (zinc/slate). One accent: **Teal** (`#0f766e` or equivalent) for primary actions and focus. Amber for SLA at-risk, red for incomplete/Reject, teal for Approve/complete.
*   **Typography:** Geist or IBM Plex Sans. **Mono** (Geist Mono or IBM Plex Mono) is **mandatory** for Claim IDs, timestamps, SLA clocks, and scores. Numbers must always be rendered in mono.
*   **Density & Shapes:** Compact padding (`gap-3` / `p-4`). 8px border-radius on surfaces, 6px on inputs. Hairline borders + tinted shadows.

### UX Principles
*   **Real-world language:** Labels must match glossary terms exactly. Never show "AI decision" or "HITL".
*   **Error Prevention:** Approve button uses an inline confirmation checkbox, not a surprise modal. Label says: *"Creates a Decision and closes this Claim."*
*   **No-Glitch Layout:** Reserve min-height for SLA, completeness block, and timeline. Polling/updates must keep the last good data on screen with skeleton loading instead of jumping or flashing empty.
*   **Human-In-The-Loop:** Human approval is mandatory. Recommendations never auto-approve or auto-close.

---

## 5. Architectural Alignment (ADRs & RFC)

*   **ADR 0001 (Multi-Agent ADK):** We run multiple ADK specialist agents (Intake, Assignment, SLA, Assessment) orchestrated by a **Supervisor** that reads Claim state and determines "what next". Specialists do *not* call each other. Gemini is attached *only* to Intake and Assessment.
*   **ADR 0002 (Operational State):** Firestore is the system of record for claims and append-only `claim_events`. Cloud Storage stores raw documents (blobs), and Firestore holds object references.
*   **ADR 0003 (Deterministic Chores vs Gemini Reasoning):**
    *   *Deterministic Code:* SLA clocks and Assignment scores (workload, skill, and criteria) are computed deterministically in code to ensure repeatability and trust.
    *   *Gemini:* Document extraction, classification assistance, and Assessment recommendation narrative/reasons use Gemini 3.5 Flash JSON outputs.
    *   *Human approval:* Decision creation is fully human-driven.

---

## 6. Building, Running, and Testing Commands

Commands are run within the `/web/` directory.

### Setup & Install
```bash
cd web
npm install
```

### Local Development
```bash
npm run dev
```

### Compilation & Type-Checking
```bash
npm run typecheck
# or:
tsc --noEmit
```

### Testing (Vitest + RTL)
```bash
npm run test
# For watching mode:
npm run test:watch
```

### Production Build
```bash
npm run build
```

---

## 7. Development Guidelines & Pre-flight Checklist

Before completing any task or PR:
1.  **Language Check:** Verify that only glossary terms are used. "Assessment recommendation" not "AI decision". "Human approval" not "HITL". No em-dashes (`—`) or en-dashes (`–`) in visible UI copy.
2.  **State Completeness:** Ensure loading, empty, error, and success states exist and are visually correct for any interactive section added.
3.  **No Layout Jumps:** Reserve container layout heights to prevent elements from shifting when async operations complete.
4.  **Accessibility (WCAG AA):** Focus outlines must be visible, contrast must meet AA standard (especially white text on teal), and `prefers-reduced-motion` must be respected.
5.  **Test Verification:** Ensure all edits are covered by matching tests in `App.test.tsx` (or a dedicated new test file). Running `npm run test` must pass cleanly.
6.  **No Unrequested Commits:** Do not stage or commit files unless explicitly directed.
