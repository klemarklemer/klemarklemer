# Claims ops UI/UX guideline

**Reading this as:** a regulated insurance **operations console** for claims officers (not a landing page, not a portfolio), with a trust-first modern language, leaning toward IBM Carbon density plus owned components (Radix or customised shadcn), not marketing heroes.

`design-taste-frontend` does not own dashboards. This file is the product UI contract. Pull from that skill only: anti-slop (no AI-purple, no Inter-as-default, no three equal marketing cards), WCAG, reduced motion, reserved layout space, one accent, one radius scale. Do not ship a hero, logo wall, marquee, or scroll-hijack on this app.

**Dials (locked):** `DESIGN_VARIANCE: 3` · `MOTION_INTENSITY: 2` · `VISUAL_DENSITY: 7`

Experienced means: the officer never wonders what to do next. Modern means: crisp type, one accent, quiet chrome. No glitch means: layout does not jump, numbers do not flicker, motion does not fight the pointer. UX means: Nielsen heuristics plus insurance HITL, below.

Domain words come from [CONTEXT.md](../../CONTEXT.md). Product scope from [taskmaster-claims-ops.md](../prd/taskmaster-claims-ops.md).

## When an agent opens this file

You are about to change the ops UI. Follow every **Must** in this file. Done means the Pre-flight at the bottom can be ticked honestly.

## One screen

The hackathon UI is a **Claim workspace**, not a command centre and not a customer portal.

```
+-- App bar: product name, seeded officer, notifications bell (count)
+-- Claim identity: id, Policy id, line MOTOR
+-- Stage + SLA clock (stage-level; remaining or at-risk)
+-- Document completeness + upload
+-- Assignment (scores + owner) when present
+-- Assessment recommendation + Approve / Reject when present
+-- Claim event timeline (newest first or chronological; pick one and keep it)
```

Primary path is vertical: identity → blockers (missing docs, SLA) → action (upload or Human approval) → timeline. Do not hide Approve below a fold on desktop if a recommendation exists.

## UX principles (how they bind this UI)

| Heuristic | Must |
|---|---|
| System status | Stage, SLA clock, and Document completeness are always visible when a Claim is open. At-risk is a text + colour state, not only an icon. |
| Real-world language | Labels use glossary terms: Claim, Stage, SLA clock, Assignment, Assessment recommendation, Decision, Human approval, Claim event. |
| User control | Assessment recommendation never auto-closes the Claim. Approve and Reject are distinct, labelled, keyboard-reachable. Reject does not look like Approve. |
| Consistency | Same Claim id format, same timestamp format (ISO-like local or UTC labelled), same Stage names everywhere. |
| Error prevention | Approve that creates a Decision uses an inline confirm (checkbox or second click on the same control), not a surprise modal maze. State what will happen: "Creates a Decision and closes this Claim." |
| Recognition | Timeline is a readable list of Claim events (who/what/when). Do not require the officer to remember chat. |
| Flexibility | Seed / demo reset can live in the app bar. Power users get the same layout as the demo, not a second skin. |
| Minimalism | No marketing copy, no fake KPI wall, no "AI copilot" chat dock. Notifications are in-app only. |
| Recover | Upload and API errors sit next to the control. Failed Approve leaves the Claim open and says so. |
| Access | WCAG 2.2 AA: contrast, visible focus, labels not placeholders, `prefers-reduced-motion`. |

## Experienced claims officer (job to be done)

In the 4-minute demo the officer should be able to:

1. See the Claim is incomplete without reading a paragraph.
2. Upload the police report from the completeness block.
3. See Assignment scores and owner without opening a picker.
4. See SLA at-risk without hunting a toast that already vanished.
5. Read recommendation reasons, then Approve with confirm, or Reject without closing.

If a control does not serve one of those jobs, it does not ship.

## Visual system

**Type:** Geist or IBM Plex Sans (self-hosted or `next/font` / equivalent). Not Inter unless already in the repo. Mono (Geist Mono or IBM Plex Mono) for Claim ids, timestamps, scores, confidence. Numbers in mono.

**Colour:** zinc/slate neutrals, one accent: **teal** (`oklch` or `#0f766e` family) for primary actions and focus. Semantic: amber for SLA at-risk, red for incomplete documents and Reject, teal for Approve and complete. Do not use purple glows. No pure `#000` or `#fff`; use zinc-950 / zinc-50. One theme for the whole app (light default, `dark:` tokens allowed if both modes are complete).

**Shape:** 8px radius on surfaces, 6px on inputs, not mixed pill + sharp. Hairline `border` + tinted shadow, not floating marketing cards.

**Density:** compact padding (`gap-3` / `p-4` on the workspace). SLA remaining is large enough to read at a glance (`text-lg` mono), not a tiny badge only.

**Icons:** one family (Phosphor preferred). No emoji in the chrome.

## Motion and no-glitch contract

Motion exists only as feedback: 150ms opacity/transform on buttons (`active:scale-[0.98]`), no page load choreography, no infinite pulse on SLA.

**No glitch:**

- Reserve min-height for SLA, completeness, and timeline so async refresh does not jump the Approve button.
- Skeleton shapes match the final layout; no centred spinner as the only loading state for the Claim workspace.
- Polling or SSE must not flash empty then full (keep last good Claim on screen; mark stale if needed).
- Images/PDFs: do not layout-shift; use a fixed preview slot or a text link to the object.
- `min-h-[100dvh]` for the shell if you need full height; never `h-screen`.
- Honour `prefers-reduced-motion: reduce` (instant states).
- No `window` scroll listeners. No custom cursor.

## Copy

- UI strings: glossary terms. "Assessment recommendation" not "AI decision". "Human approval" not "HITL".
- No em-dash (`—`) or en-dash (`–`) in visible UI; use hyphen or a new sentence.
- No "Unleash", "Seamless", "Next-gen".
- Buttons: **Approve** and **Reject** (one intent each). Upload: **Upload police report**. Seed: **Seed demo Claim**.
- Empty timeline: "No Claim events yet." plus what will appear after Intake.
- Confidence: show as a percent or 0-1 consistently; label "Confidence", not "Sure".

## Component states (every interactive region)

| Region | Loading | Empty | Error | Success |
|---|---|---|---|---|
| Claim workspace | Skeleton of identity + Stage + completeness | "Seed demo Claim" primary | Inline error + retry | Claim body |
| Completeness | Skeleton list | "No documents yet" | Upload failed, keep prior list | Checklist with missing/ok |
| Notifications | Bell with no flash of 99 | "No notifications" | Ignore transient fetch fail; keep last count | List in a panel, not only toast |
| Timeline | Skeleton rows | Sentence above | "Could not load Claim events" | Stable order |
| Approve | Button disabled + "Working" | Hidden until recommendation exists | Error next to button; Claim still open | Claim closed; Decision visible |

Toasts are for transient confirm ("Police report uploaded"). SLA at-risk and missing documents stay in the layout, not only in a disappearing toast.

## Human approval pattern

- Recommendation block: outcome, reasons (list), confidence, model is optional in a disclosure, never the hero.
- **Approve** is filled teal; **Reject** is outline or danger, equal size, not greyed as secondary if both are valid.
- Confirm copy names Decision and close. Keyboard: Tab order identity → completeness → recommendation → Approve → Reject → timeline.
- After close, the workspace is read-only except timeline and PDF link. Do not leave a live Approve on a closed Claim.

## Accessibility

- Contrast AA for body and buttons (teal on white must be checked; darken accent if it fails).
- Every input has a visible `<label>`. File input: associated label, not icon-only.
- Focus ring 2px teal offset. Do not `outline-none` without a replacement.
- Live region for SLA at-risk and new notifications (`aria-live="polite"`), not a noisy assertive loop.
- Hit target ≥ 24px (44px preferred for Approve/Reject).

## Stack (when building)

React or server-rendered HTML is fine. If React: Tailwind v4, one icon set, Motion only if you can keep intensity 2 (or skip Motion). Do not add GSAP, Three.js, or a chat widget.

## Anti-patterns (fail the UI ticket)

Landing hero, gradient mesh, glassmorphism nav, three feature cards, Inter + purple, fake dashboard widgets, div-based fake "AI thinking" terminals, chat-first layout, auto-approve, emoji status, layout jump on poll, toast-only SLA.

## Pre-flight (UI tickets)

- [ ] Design read treated as ops console, not marketing.
- [ ] Glossary labels only; no em-dash in UI copy.
- [ ] Stage + completeness + (when in scope) SLA visible without scroll on a 1280px desktop for the seeded Claim.
- [ ] Approve/Reject cannot be confused; Decision is confirmed; closed Claim is read-only.
- [ ] Loading/empty/error/success exist for every region this ticket added.
- [ ] No layout jump: reserved space; last-good data while refreshing.
- [ ] `prefers-reduced-motion` respected; no infinite animation.
- [ ] Focus visible; labels present; contrast checked on accent buttons.
- [ ] One accent, one radius scale, mono for ids and numbers.
