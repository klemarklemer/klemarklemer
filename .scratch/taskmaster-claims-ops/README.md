# Taskmaster claims ops — agent handover

Work **one ticket per fresh session**. Read [CONTEXT.md](../../CONTEXT.md), then [spec.md](spec.md), then the ticket. If the ticket includes ops UI, also read [docs/ui/claims-ops-ux.md](../../docs/ui/claims-ops-ux.md) and pass its Pre-flight. Run `/implement` (it drives `/tdd`). Commit on the ticket.

**Frontier:** the lowest-numbered file whose `Blocked by` tickets are all `Status: done` and whose own status is `ready-for-agent`. Claim it by setting `Status: claimed` before coding; set `Status: done` when every acceptance box is checked.

Do not triage these tickets. Do not implement later tickets in the same session.
