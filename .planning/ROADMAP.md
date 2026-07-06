# Roadmap: CriminalSay

## Overview

CriminalSay ships in three phases: first the data foundation (Quote struct, embedded JSON, random selection), then the render and output layer (lipgloss Estilo 4 styling, word-wrap, color degradation, binary wiring), and finally build and distribution (cross-compile Makefile, PT-BR README). Each phase leaves a fully working artifact; later phases depend on earlier ones being stable.

## Phases

**Phase Numbering:**

- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Data Foundation** - go.mod, project layout, embedded quotes JSON, internal/quotes package, data layer tests (completed 2026-07-05)
- [x] **Phase 2: Render and Output** - internal/render with Estilo 4 sidebar, word-wrap, color degradation, main.go wiring; go test ./... green (completed 2026-07-06)
- [ ] **Phase 3: Build and Distribution** - cross-compile Makefile, dist/ artifacts, PT-BR README

## Phase Details

### Phase 1: Data Foundation

**Goal**: The data layer is complete — quotes are embedded at compile time, selectable at random, and fully tested before any rendering code exists.
**Mode:** mvp
**Depends on**: Nothing (first phase)
**Requirements**: DAT-01, DAT-02, DAT-03, CORE-01, CORE-03, BUILD-01, BUILD-03
**Success Criteria** (what must be TRUE):

  1. `go build -o criminalsay .` produces a binary from the repository root (validates go.mod, embed anchor, and package layout)
  2. `internal/quotes` returns a non-empty `Quote` struct with all fields populated when `Load` is called on the embedded bytes
  3. Repeated calls to `quotes.Random` return different quotes over multiple runs, confirming uniform random selection
  4. `go test ./internal/quotes/...` passes with coverage of `Load` and `Random`
  5. `data/quotes.json` contains at least 15 real Criminal Minds quotes, each with `quote`, `author`, `season`, and `episode` fields

**Plans**: 1/1 plans complete

Plans:

- [x] 01-01-PLAN.md — Walking Skeleton: module scaffold, embedded quotes data, internal/quotes package, Phase 1 stub binary

### Phase 2: Render and Output

**Goal**: Running `criminalsay` prints a styled Criminal Minds quote to stdout and exits with code 0; styling degrades cleanly when color is unavailable.
**Mode:** mvp
**Depends on**: Phase 1
**Requirements**: CORE-02, REND-01, REND-02, REND-03, REND-04, REND-05, REND-06, COLOR-01, COLOR-02, BUILD-03
**Success Criteria** (what must be TRUE):

  1. Running `criminalsay` in a color terminal shows a red `▌` sidebar on every line of the quote, quoted text in white/bright, and a dim attribution line formatted as `<author> · Criminal Minds · S<season>E<episode>`
  2. Long quotes wrap at ≤52 useful columns with the `▌` sidebar repeated and aligned on every wrapped line
  3. Running `NO_COLOR=1 criminalsay` or piping output (`criminalsay | cat`) produces readable plain text with `▌`, quotation marks, and `·` separators preserved but no ANSI escape codes
  4. A quote with a missing optional field (empty `character` or `episodeTitle`) renders without layout breakage
  5. `go test ./...` passes green, covering both `internal/quotes` and `internal/render`

**Plans**: 1/1 plans complete

Plans:

- [x] 02-01-PLAN.md — Estilo 4 render package, color degradation, main.go wiring, full test gate

**UI hint**: yes

### Phase 3: Build and Distribution

**Goal**: The project can be cross-compiled to all target platforms from a single `make` command and is documented for Portuguese-speaking users.
**Mode:** mvp
**Depends on**: Phase 2
**Requirements**: BUILD-02, DOC-01
**Success Criteria** (what must be TRUE):

  1. `make all` (or equivalent Makefile target) produces binaries for Linux amd64, macOS amd64, macOS arm64, and Windows amd64 in a `dist/` directory, all built with `CGO_ENABLED=0`
  2. Each binary in `dist/` runs `criminalsay` and outputs a formatted quote without errors on its target platform
  3. `README.md` in PT-BR documents installation (via `go install` or downloading a binary), usage (run `criminalsay`), and build instructions (clone + `make`)

**Plans**: 1/2 plans executed

- [x] 03-01-PLAN.md
- [ ] 03-02-PLAN.md

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Data Foundation | 1/1 | Complete    | 2026-07-05 |
| 2. Render and Output | 1/1 | Complete    | 2026-07-06 |
| 3. Build and Distribution | 1/2 | In Progress|  |
