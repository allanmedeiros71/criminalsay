---
phase: 03
slug: build-and-distribution
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-07-06
---

# Phase 03 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` (existing) + Makefile targets (new) |
| **Config file** | none |
| **Quick run command** | `make test` |
| **Full suite command** | `make test` (= `go test ./...`) |
| **Build gate command** | `make all` |
| **Estimated runtime** | ~5 seconds (tests) + ~30 seconds (cross-compile) |

---

## Sampling Rate

- **After every task commit:** Run `make test`
- **After every plan wave:** Run `make all && make test`
- **Before `/gsd-verify-work`:** `make clean && make all && make test && file dist/*`
- **Max feedback latency:** 35 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 03-01-01 | 01 | 1 | BUILD-02 | T-03-01 | Makefile uses local `go build` only — no remote fetch | integration | `make clean && make all && test $(ls dist/ \| wc -l) -eq 4` | ❌ W0 | ⬜ pending |
| 03-01-02 | 01 | 1 | BUILD-02 | — | N/A | integration | `file dist/criminalsay-linux-amd64 \| grep -q 'ELF.*x86-64'` | ❌ W0 | ⬜ pending |
| 03-01-03 | 01 | 1 | BUILD-02 | — | N/A | integration | `file dist/criminalsay-darwin-amd64 \| grep -q 'Mach-O.*x86_64'` | ❌ W0 | ⬜ pending |
| 03-01-04 | 01 | 1 | BUILD-02 | — | N/A | integration | `file dist/criminalsay-darwin-arm64 \| grep -q 'Mach-O.*arm64'` | ❌ W0 | ⬜ pending |
| 03-01-05 | 01 | 1 | BUILD-02 | — | N/A | integration | `file dist/criminalsay-windows-amd64.exe \| grep -q 'PE32+'` | ❌ W0 | ⬜ pending |
| 03-01-06 | 01 | 1 | BUILD-02 | — | N/A | smoke | `./dist/criminalsay-linux-amd64 \| head -1` | ❌ W0 | ⬜ pending |
| 03-01-07 | 01 | 1 | BUILD-02 | T-03-02 | `.gitignore` prevents binary commit | integration | `make clean && test ! -e dist && test ! -e criminalsay` | ❌ W0 | ⬜ pending |
| 03-01-08 | 01 | 1 | BUILD-03 | — | N/A | integration | `make test` | ✓ | ⬜ pending |
| 03-02-01 | 02 | 2 | DOC-01 | T-03-03 | No GitHub Releases references in docs | manual/grep | `test -f README.pt.md && grep -qi 'instala' README.pt.md` | ❌ W0 | ⬜ pending |
| 03-02-02 | 02 | 2 | DOC-01 | — | N/A | manual/grep | `grep -q '1.22' README.pt.md && grep -q '1.22' README.md && ! grep -q '1.21' README.md` | ❌ W0 | ⬜ pending |
| 03-02-03 | 02 | 2 | DOC-01 | — | N/A | manual/grep | `! grep -qi 'releases' README.md README.pt.md` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `Makefile` — all targets per D-01–D-12
- [ ] `.gitignore` — `dist/`, `criminalsay` per D-15
- [ ] `README.pt.md` — minimal DOC-01 per D-07
- [ ] `README.md` — sync fixes per D-06/D-08

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| macOS binary runs on darwin | BUILD-02 | Cross-compiled Mach-O cannot execute on Linux host | Run `dist/criminalsay-darwin-amd64` on macOS amd64; confirm quote output |
| Windows binary runs on Windows | BUILD-02 | PE binary cannot execute on Linux host | Run `dist/criminalsay-windows-amd64.exe` on Windows; confirm quote output |
| Windows cmd.exe lipgloss fallback | BUILD-02 | STATE.md notes unverified | Run binary in Windows cmd.exe; confirm readable output |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 35s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
