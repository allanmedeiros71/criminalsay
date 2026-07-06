---
phase: 02-render-and-output
plan: "01"
subsystem: render
status: complete
tags: [go, lipgloss, render, terminal, tdd, estilo-4]
completed_date: "2026-07-06"
duration_minutes: 15
task_count: 4
file_count: 4

dependency_graph:
  requires:
    - phase: 01-data-foundation
      provides: internal/quotes package, embedded JSON, main.go embed anchor
  provides:
    - internal/render package (Quote, wrap, attribution builder, color/plain modes)
    - main.go wired with colorEnabled() and render.Quote
    - render_test.go covering REND/COLOR/BUILD-03 requirements
  affects:
    - Phase 3 (build and distribution — binary now prints styled output)

tech_stack:
  added:
    - charmbracelet/lipgloss v1.1.0 (first direct import in internal/render)
    - charmbracelet/x/term v0.2.1 (TTY detection in main)
    - charmbracelet/x/ansi v0.8.0 (display-width assertions in tests)
  patterns:
    - colorEnabled bool passed from main — render never detects TTY (D-14)
    - wrap-then-decorate — word wrap before lipgloss styling (Pitfall 4)
    - per-line sidebar prefix on every wrapped quote line (D-07)
    - plain mode uses ASCII pipe sidebar, zero ANSI (D-18)

key_files:
  created:
    - internal/render/render.go
    - internal/render/wrap.go
    - internal/render/render_test.go
  modified:
    - main.go

decisions:
  - "lipgloss.SetColorProfile(termenv.TrueColor) when colorEnabled=true so ANSI renders in tests and TTY"
  - "Manual per-line prefix composition after wrap — no Style.Border() for sidebar (D-20)"

requirements-completed: [CORE-02, REND-01, REND-02, REND-03, REND-04, REND-05, REND-06, COLOR-01, COLOR-02, BUILD-03]

coverage:
  - id: D1
    description: "Estilo 4 render with 52-col wrap, sidebar on every line, attribution builder"
    requirement: REND-01
    verification:
      - kind: unit
        ref: "internal/render/render_test.go#TestQuote_sidebar"
        status: pass
    human_judgment: false
  - id: D2
    description: "Plain-text degradation with pipe sidebar and no ANSI when color disabled"
    requirement: COLOR-02
    verification:
      - kind: unit
        ref: "internal/render/render_test.go#TestQuote_noColor"
        status: pass
    human_judgment: false
  - id: D3
    description: "main.go wires colorEnabled detection and styled output"
    requirement: CORE-02
    verification:
      - kind: integration
        ref: "go build -o criminalsay . && NO_COLOR=1 ./criminalsay"
        status: pass
    human_judgment: true
    rationale: "Visual styling in real TTY requires human eye check per ROADMAP UI hint"

metrics:
  duration: 15m
  completed_date: "2026-07-06"
  tasks_completed: 4
  tasks_total: 4
  files_created: 3
  files_modified: 1
---

# Phase 02 Plan 01: Render and Output Summary

**One-liner:** internal/render Estilo 4 package with 52-col word wrap, red ▌ sidebar, yellow/dim attribution, NO_COLOR pipe degradation, and main.go wiring — `go test ./...` green.

## What Was Built

Replaced the Phase 1 stub (`q.Quote + " — " + q.Author`) with a full render layer. Running `./criminalsay` now prints a multi-line Estilo 4 quote with sidebar, quoted text, blank line, and indented attribution. Color mode uses lipgloss; plain mode uses `|` sidebar with zero ANSI sequences.

**Artifacts:**
- `internal/render/wrap.go` — 52-column word wrap with quote-mark assembly (D-05–D-08)
- `internal/render/render.go` — `Quote(q, colorEnabled)`, attribution builder (D-01–D-13), lipgloss styles (D-16–D-17)
- `internal/render/render_test.go` — 12 table-driven tests for layout, wrap width, attribution, color/plain modes
- `main.go` — `colorEnabled()` from NO_COLOR + `term.IsTerminal`, `render.Quote` wiring

## Task Commits

| Task | Name | Focus |
|------|------|-------|
| 1 | Render core + tests | wrap.go, plain-text layout, 7 core tests |
| 2 | Color modes | lipgloss styles, TestQuote_color, TestQuote_noBox |
| 3 | main.go wiring | colorEnabled(), stub replacement |
| 4 | Full suite gate | TestQuote_quotes/attribution/episodeTitle, BUILD-03 |

## Self-Check: PASSED

- `go test ./... -count=1` exits 0
- `go build -o criminalsay .` exits 0
- `NO_COLOR=1 ./criminalsay` uses pipe sidebar, no ANSI
- Phase 1 embed anchor and error handling preserved in main.go
