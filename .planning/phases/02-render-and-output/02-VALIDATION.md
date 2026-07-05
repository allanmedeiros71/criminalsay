---
phase: 02
slug: render-and-output
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-07-05
---

# Phase 02 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go stdlib `testing` |
| **Config file** | none |
| **Quick run command** | `go test ./internal/render/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~2 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/render/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 02-01-01 | 01 | 1 | REND-01,02,03,06 | — | N/A | unit | `go test ./internal/render/... -run TestQuote` | ❌ W0 | ⬜ pending |
| 02-01-02 | 01 | 1 | REND-04,05 | — | N/A | unit | `go test ./internal/render/... -run TestQuote_attribution` | ❌ W0 | ⬜ pending |
| 02-01-03 | 01 | 1 | COLOR-01,02 | — | N/A | unit | `go test ./internal/render/... -run TestQuote_noColor` | ❌ W0 | ⬜ pending |
| 02-01-04 | 01 | 1 | CORE-02, BUILD-03 | — | N/A | integration | `go build -o criminalsay . && go test ./...` | partial | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/render/render.go` — `Quote()`, attribution builder, lipgloss styles
- [ ] `internal/render/wrap.go` — 52-col word wrap helpers (optional split)
- [ ] `internal/render/render_test.go` — table-driven tests, `colorEnabled` true/false
- [ ] `main.go` update — replace stub, add `colorEnabled()`, import render

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Visual color in truecolor terminal | COLOR-01 | ANSI depth varies by terminal | Run `./criminalsay` in color terminal; confirm red sidebar and yellow author |
| Pipe output readability | COLOR-02 | D-18 uses `\|` not `▌` | Run `criminalsay \| cat`; confirm no ANSI, `\|` sidebar |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
