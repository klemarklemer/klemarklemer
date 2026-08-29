# [api] fix: Update Go Toolchain to go1.26.7 in go.mod and CI to Resolve Stdlib Vulnerabilities

**Date:** 2026-08-30  
**Timestamp:** 1788023208  
**Type:** `[api]` security fix / CI  
**Service:** `apps/api/services/core`

---

## Summary

Updated `apps/api/go.mod` to specify `toolchain go1.26.7` and `.github/workflows/ci.yml` to use `go-version: 'stable'`. This patches all 26 Go standard library CVEs identified by `govulncheck`, resulting in a clean security audit (0 vulnerabilities).
