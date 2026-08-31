# [web] docs: Stop the architecture page calling the console unwired

**Date:** 2026-08-31  
**Timestamp:** 1788182968  
**Type:** `[web]` documentation fix  
**Service:** `docs/architecture`

---

## Summary

The architecture page still drew the console as a dangling node - a dotted
`not yet wired` edge to the REST delivery, dashed node styling, and a note reading
"currently a design prototype with local fixture data - it makes no HTTP calls".

None of that has been true since the frontend was connected to the API. The note
also cited `apps/web/src/data/claim.ts`, which was deleted in that same change.

This is the diagram a judge looks at, so it was actively understating the project.
The same claim was corrected in `docs/SUBMISSION.md` separately.

The edge is now solid and labelled `HTTP · /v1`, the gap styling is gone, and the
note points at `apps/web/src/api/claims.ts`, which exists.

## Verification

Confirmed against a clean clone of `main`: the console calls `/v1/claim`,
`/v1/claim/:id/documents`, `/v1/claim/:id/assessment`, `/v1/claim/:id/approval`
and `/v1/demo/reset` through `apps/web/src/api/claims.ts`, and
`apps/web/src/data/claim.ts` no longer exists.
