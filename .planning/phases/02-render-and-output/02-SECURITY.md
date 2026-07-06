---
phase: 02
slug: render-and-output
status: verified
threats_open: 0
asvs_level: 1
created: 2026-07-06
---

# Phase 02 — Security

> Per-phase security contract: threat register, accepted risks, and audit trail.

---

## Trust Boundaries

| Boundary | Description | Data Crossing |
|----------|-------------|---------------|
| compile-time → binary | Embedded quotes.json baked at build; render formats trusted Quote struct fields only | Static quote JSON (public TV dialogue) |
| render → stdout | Formatted string crosses to terminal; quote text from curated embedded data, not user input | Quote text + attribution + ANSI when color enabled |
| lipgloss → ANSI output | When colorEnabled=true, lipgloss emits escape sequences to stdout | Terminal styling only; no secrets |

---

## Threat Register

| Threat ID | Category | Component | Severity | Disposition | Mitigation | Status |
|-----------|----------|-----------|----------|-------------|------------|--------|
| T-02-01 | Spoofing | Quote text with embedded terminal control sequences | low | accept | Curated JSON under source control; no runtime user input | closed |
| T-02-02 | Tampering | lipgloss style output via malformed Quote fields | low | accept | Compile-time embedded fields; string concat only, no eval | closed |
| T-02-03 | Denial of Service | Panic on empty quotes slice in main | medium | mitigate | Phase 1 `len(qs)==0` guard preserved in main.go before Random | closed |
| T-02-04 | Information Disclosure | ANSI sequences leak terminal capabilities | low | accept | Standard terminal styling; no secrets in output | closed |
| T-02-05 | Elevation | lipgloss/x/term transitive deps supply chain | low | accept | Pinned via go.sum; packages audited in 02-RESEARCH.md | closed |
| T-02-SC | Tampering | lipgloss v1.1.0 module integrity | low | accept | Pre-declared Phase 1; checksum in go.sum | closed |

---

## Accepted Risks Log

| Risk ID | Threat Ref | Rationale | Accepted By | Date |
|---------|------------|-----------|-------------|------|
| AR-02-01 | T-02-01 | Trusted curated dataset; v1 offline model | gsd-secure-phase | 2026-07-06 |
| AR-02-02 | T-02-02 | No user-controlled render input | gsd-secure-phase | 2026-07-06 |
| AR-02-04 | T-02-04 | Normal CLI color output | gsd-secure-phase | 2026-07-06 |
| AR-02-05 | T-02-05 | Go module proxy + go.sum pinning | gsd-secure-phase | 2026-07-06 |
| AR-02-SC | T-02-SC | Same lipgloss pin as Phase 1 | gsd-secure-phase | 2026-07-06 |

---

## Security Audit Trail

| Audit Date | Threats Total | Closed | Open | Run By |
|------------|---------------|--------|------|--------|
| 2026-07-06 | 6 | 6 | 0 | gsd-secure-phase (L1, register from 02-01-PLAN.md) |

---

## Sign-Off

- [x] All threats have a disposition (mitigate / accept / transfer)
- [x] Accepted risks documented in Accepted Risks Log
- [x] `threats_open: 0` confirmed
- [x] `status: verified` set in frontmatter

**Approval:** verified 2026-07-06
