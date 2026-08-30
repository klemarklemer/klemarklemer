# [web] feat: Connect the Claims Ops console to the backend API

**Date:** 2026-08-30  
**Timestamp:** 1788059982  
**Type:** `[web]` feature  
**Service:** `apps/web`

---

## Summary

The web console was a standalone mock driven by synthetic in-memory seed data
(`src/data/claim.ts`). It is now wired to the real Go REST API.

1. Added `src/api/client.ts` — `fetch` wrapper that targets `VITE_API_BASE`
   (default `/api`) and unwraps the candi `HTTPResponse` envelope (`{ success, data }`).
2. Added `src/api/claims.ts` — typed backend models, API calls
   (`resetDemo`, `uploadDocument`, `runAssessment`, `submitApproval`) and a
   `toClaimView` adapter mapping `ResponseClaim` → the frontend `Claim` shape.
3. Rewired `App.tsx`: Seed → `POST /v1/demo/reset`; Upload → `POST /v1/claim/:id/documents`
   then `POST /v1/claim/:id/assessment`; Approve/Reject → `POST /v1/claim/:id/approval`.
4. Added a Vite dev proxy (`/api` → `http://localhost:8000`) to avoid CORS.
5. Corrected a model mismatch: completeness now derives from the backend
   `document_completeness` field (the API only stores present documents, so the
   previous documents-array heuristic would have hidden the upload control).
6. Removed the dead `src/data/claim.ts` mock; updated `App.test.tsx` to mock
   `./api/client` so the suite stays green without a running backend.

## Notes

- Live backend verification was not possible here (Docker unavailable); the
  contract was confirmed against source and the build/typecheck/test suite passes.
- Run the full stack with `make up` (API on `:8000`, web on `:5173`) to exercise it.
