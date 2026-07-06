# Project Retrospective

*A living document updated after each milestone. Lessons feed forward into future planning.*

## Milestone: v1.0 — MVP

**Shipped:** 2026-07-06
**Phases:** 3 | **Plans:** 4

### What Was Built

- Go module with 20 embedded Criminal Minds quotes, `internal/quotes` package, and data-layer tests
- `internal/render` Estilo 4 styling with word-wrap, color degradation, and `main.go` wiring
- Cross-compile Makefile (`CGO_ENABLED=0`) producing four platform binaries in `dist/`
- `README.md` synced and `README.pt.md` with Instalação, Uso, and Build sections

### What Worked

- Vertical MVP phases left a working artifact after each phase (data → render → ship)
- GSD wave execution with plan summaries enabled clean verification and UAT
- Threat models in PLAN.md made security review straightforward at ship time

### What Was Inefficient

- No milestone audit run before close — recommend `/gsd-audit-milestone` for v1.1+
- Default make shell on host required `SHELL := /bin/bash` workaround (undocumented in plan)

### Patterns Established

- `export CGO_ENABLED=0` at Makefile top for all cross-compile targets
- Separate EN (`README.md`) and PT-BR (`README.pt.md`) user documentation
- Phase branches (`gsd/phase-NN-*`) with one PR per phase

### Key Lessons

1. Cross-compile naming contract (`dist/criminalsay-<platform>`) should be locked before writing README tables
2. Human UAT on non-Linux platforms catches doc/reality gaps automated grep cannot
3. Security ship gate requires `*-SECURITY.md` — generate from PLAN threat models at phase end

### Cost Observations

- Timeline: 2026-07-05 → 2026-07-06 (~2 days)
- 4 plans across 3 phases; PR #3 open for Phase 3 merge

---

## Cross-Milestone Trends

### Process Evolution

| Milestone | Phases | Plans | Key Change |
|-----------|--------|-------|------------|
| v1.0 | 3 | 4 | Greenfield GSD workflow from PROJECT.md through ship |

### Cumulative Quality

| Milestone | Tests | Coverage | Zero-Dep Additions |
|-----------|-------|----------|-------------------|
| v1.0 | go test ./... green | quotes + render | lipgloss v1.1.0 only |

### Top Lessons (Verified Across Milestones)

1. Walking-skeleton first (Phase 1) de-risked render and build work
2. `go test ./...` as regression gate across all phases prevented drift
