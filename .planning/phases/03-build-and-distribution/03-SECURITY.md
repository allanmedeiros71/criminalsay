---
phase: 03
slug: build-and-distribution
status: verified
threats_open: 0
asvs_level: 1
created: 2026-07-06
---

# Phase 03 — Security

> Per-phase security contract: threat register, accepted risks, and audit trail.

---

## Trust Boundaries

| Boundary | Description | Data Crossing |
|----------|-------------|---------------|
| developer → Makefile | Developer invokes local build recipes; no remote fetch | Local go build commands only |
| build → dist/ | Cross-compiled binaries written to local filesystem | Embedded quote data in binaries |
| dist/ → VCS | Build artifacts must not enter git history | Binary artifacts blocked by .gitignore |
| docs → user | User follows install/build instructions from README files | go install via Go proxy; source build via Makefile |

---

## Threat Register

| Threat ID | Category | Component | Severity | Disposition | Mitigation | Status |
|-----------|----------|-----------|----------|-------------|------------|--------|
| T-03-01 | Tampering | Makefile supply chain via curl/wget/bash | high | mitigate | Makefile uses local `go build` only — no remote fetch or URL includes | closed |
| T-03-02 | Information Disclosure | Committed binary artifacts in git | medium | mitigate | `.gitignore` with `dist/` and `criminalsay`; verified via git check-ignore | closed |
| T-03-03 | Tampering | CGO_ENABLED=1 cross-compile failure | medium | mitigate | `export CGO_ENABLED=0` at Makefile top | closed |
| T-03-04 | Denial of Service | make clean misses artifacts | low | mitigate | `clean` removes both `dist/` and root `criminalsay` | closed |
| T-03-05 | Elevation | World-writable dist/ permissions | low | accept | Default umask; local dev builds only | closed |
| T-03-06 | Spoofing | Stale GitHub Releases install instructions | high | mitigate | Pre-built binaries / Releases section removed from README.md and absent from README.pt.md | closed |
| T-03-07 | Spoofing | scripts/build.sh reference | medium | mitigate | Removed from project structure tree; Makefile-only build path documented | closed |
| T-03-08 | Tampering | Doc drift — README documents make all before Makefile exists | medium | mitigate | 03-02 depends on 03-01; verification confirms make all + dist/ table match | closed |
| T-03-09 | Information Disclosure | Incorrect Go version in docs | low | mitigate | Both READMEs state Go 1.22+ aligned with go.mod | closed |
| T-03-10 | Denial of Service | User follows removed build.sh path | low | mitigate | README.pt.md documents `make` and `go build` only | closed |

---

## Accepted Risks Log

| Risk ID | Threat Ref | Rationale | Accepted By | Date |
|---------|------------|-----------|-------------|------|
| AR-03-05 | T-03-05 | Local dev artifact permissions; not a deployment surface in v1 | gsd-secure-phase | 2026-07-06 |

---

## Security Audit Trail

| Audit Date | Threats Total | Closed | Open | Run By |
|------------|---------------|--------|------|--------|
| 2026-07-06 | 10 | 10 | 0 | gsd-secure-phase (L1, register from 03-01/03-02 PLAN.md threat models) |

---

## Sign-Off

- [x] All threats have a disposition (mitigate / accept / transfer)
- [x] Accepted risks documented in Accepted Risks Log
- [x] `threats_open: 0` confirmed
- [x] `status: verified` set in frontmatter

**Approval:** verified 2026-07-06
