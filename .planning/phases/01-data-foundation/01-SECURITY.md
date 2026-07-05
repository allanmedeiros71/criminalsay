---
phase: 01
slug: data-foundation
status: verified
threats_open: 0
asvs_level: 1
created: 2026-07-05
---

# Phase 01 — Security

> Per-phase security contract: threat register, accepted risks, and audit trail.

---

## Trust Boundaries

| Boundary | Description | Data Crossing |
|----------|-------------|---------------|
| compile-time → binary | `data/quotes.json` baked into binary at build time; no filesystem access at runtime | Static quote JSON (public TV dialogue) |
| binary → stdout | Output is deterministic quote text; no user-controlled content reaches output path | Quote text + author attribution |

---

## Threat Register

| Threat ID | Category | Component | Severity | Disposition | Mitigation | Status |
|-----------|----------|-----------|----------|-------------|------------|--------|
| T-01-01 | Tampering | `data/quotes.json` (compile-time) | low | accept | Compiled into binary; source under git integrity | closed |
| T-01-02 | Denial of Service | `quotes.Random()` panics on empty slice | medium | mitigate | `main.go` guards `len(qs) == 0` before Random; exits 1 with stderr | closed |
| T-01-03 | Tampering | `go.mod` / `go.sum` supply chain | low | accept | lipgloss v1.1.0 pinned via go.sum on module proxy | closed |
| T-01-04 | Information Disclosure | Embedded binary contains quote JSON | low | accept | Public TV quotes; no private data | closed |
| T-01-SC | Tampering | `go get lipgloss@v1.1.0` install step | low | accept | Version pinned in go.sum; verified via Go module proxy | closed |

---

## Accepted Risks Log

| Risk ID | Threat Ref | Rationale | Accepted By | Date |
|---------|------------|-----------|-------------|------|
| AR-01-01 | T-01-01 | Compile-time embed; runtime tampering impossible | gsd-secure-phase | 2026-07-05 |
| AR-01-03 | T-01-03 | Standard Go module pinning via go.sum | gsd-secure-phase | 2026-07-05 |
| AR-01-04 | T-01-04 | Public domain / published TV dialogue | gsd-secure-phase | 2026-07-05 |
| AR-01-SC | T-01-SC | Go ecosystem module proxy verification | gsd-secure-phase | 2026-07-05 |

---

## Security Audit Trail

| Audit Date | Threats Total | Closed | Open | Run By |
|------------|---------------|--------|------|--------|
| 2026-07-05 | 5 | 5 | 0 | gsd-secure-phase (L1 grep-depth, register from PLAN.md) |

---

## Sign-Off

- [x] All threats have a disposition (mitigate / accept / transfer)
- [x] Accepted risks documented in Accepted Risks Log
- [x] `threats_open: 0` confirmed
- [x] `status: verified` set in frontmatter

**Approval:** verified 2026-07-05
