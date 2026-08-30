# [ops] fix: Add pull-requests read permission for paths-filter and update Node.js to 22 LTS

**Date:** 2026-08-30  
**Timestamp:** 1788058086  
**Type:** `[ops]` CI fix  
**Service:** `.github/workflows/ci.yml`, `.github/workflows/cd.yml`

---

## Summary

1. Added `pull-requests: read` and `contents: read` permissions to `.github/workflows/ci.yml` so that `dorny/paths-filter@v3` has access to fetch the pull request changed files list via the GitHub API.
2. Updated Node.js setup in `ci.yml` and `cd.yml` from `20` to `22` (LTS), eliminating GitHub Actions Node 20 deprecation warnings.
