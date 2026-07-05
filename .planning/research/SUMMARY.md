# Project Research Summary

**Project:** CriminalSay
**Domain:** fortune-style Go CLI tool with terminal styling
**Researched:** 2026-07-05
**Confidence:** MEDIUM

## Executive Summary

CriminalSay is a single-binary Go CLI tool in the tradition of `fortune` — run it, get a random Criminal Minds quote with a styled sidebar, done. The expert pattern for this domain is: embed curated JSON data at compile time, parse once at startup, render with lipgloss, print and exit. The entire program is linear with no goroutines, no network, no filesystem I/O at runtime. Cold start must stay under 50ms for `.bashrc` use; this is trivially achievable with embedded data and a single JSON unmarshal.

The recommended stack is Go 1.22+ with lipgloss v1.1.0 and nothing else external. lipgloss v1.1.0 handles TTY detection, NO_COLOR, and color profile downsampling automatically — exactly what a simple stdout CLI needs. The research explicitly rejects lipgloss v2 (manual color writer setup is unnecessary complexity), fatih/color (redundant alongside lipgloss), and all runtime dependencies beyond stdlib. The entire external dependency footprint is one package.

The critical risks are architectural and easy to avoid if addressed upfront. The go:embed directive cannot reference files outside the embedding file's package subtree — the embed anchor must live in `main.go` and pass raw bytes into `internal/quotes`. Word-wrap must account for sidebar width before prepending the `▌` prefix, and terminal width detection must guard against zero/error returns. These are well-understood patterns with clear solutions; the project is low-risk overall.

## Key Findings

### Recommended Stack

A minimal dependency footprint: one external package (lipgloss v1.1.0) and three stdlib packages (embed, math/rand/v2, encoding/json). Go 1.22 is the hard floor because math/rand/v2 entered the standard library there. lipgloss v1.1.0 is correct for this use case — it auto-detects TTY/NO_COLOR/color depth from stdout with zero boilerplate. Build tooling is a simple Makefile with GOOS/GOARCH cross-compile invocations; GoReleaser is deferred to v2.

**Core technologies:**
- Go 1.22+ (target 1.26 toolchain): Language runtime — math/rand/v2 requires 1.22; single-binary cross-compile is built-in
- charmbracelet/lipgloss v1.1.0: Terminal styling — auto-detects TTY/NO_COLOR/color depth; zero boilerplate for stdout CLI
- embed (stdlib): Compile-time data embedding — binary is self-contained, no runtime file I/O
- math/rand/v2 (stdlib): Random selection — auto-seeded, no Seed() call needed
- encoding/json (stdlib): Quote data decoding — sufficient for a small static dataset

### Expected Features

**Must have (table stakes):**
- Random quote selection — core value proposition; every run produces something different
- Quote + attribution display — users expect speaker credit alongside text
- Word wrap at terminal width — raw unwrapped output looks amateur; 80-column fallback for non-TTY
- Sidebar visual framing — the `▌` red sidebar is the defining Estilo 4 identity
- Graceful color degradation — NO_COLOR and pipe/non-TTY must produce clean plain text
- Embedded quote data — self-contained binary; go install and done, no config files
- 15+ real Criminal Minds quotes — enough variety that the tool doesn't feel like a stub
- Exit code 0 on success — shell scripts and .bashrc depend on clean exit
- Cold start under 50ms — .bashrc usage means every terminal open waits on this

**Should have (v2 differentiators):**
- Filter by character (`--character`) — power users want Gideon-only quotes
- Filter by season (`--season`) — era-specific browsing for rewatch periods
- Multiple visual themes (`--style`) — different aesthetics for different contexts
- Homebrew/AUR distribution — discoverability for non-Go users

**Defer (v2+):** `--plain` flag, `--seed` for deterministic selection, contributor JSON schema validation

**Never build:** Interactive TUI, runtime network fetching, config files

### Architecture Approach

Three-layer linear pipeline: embed at compile time, decode and select at startup, render and print to stdout. No goroutines, no shared state, no channels. Two internal packages — `internal/quotes` (data layer) and `internal/render` (presentation layer) — wired by a thin `main.go` orchestrator. The embed anchor lives in `main.go` to avoid the go:embed `..` path restriction.

**Major components:**
1. `main.go` — orchestrator only; Load, Random, color detection, render.Quote, fmt.Println; 15-20 lines max
2. `internal/quotes` — Quote struct, Load([]byte), Random([]Quote); pure stdlib
3. `internal/render` — lipgloss styles, Quote(q, colorEnabled) returns string (never prints)
4. `data/quotes.json` — flat JSON array embedded at compile time via anchor in main.go

### Critical Pitfalls

1. **go:embed path cannot use `..`** — embed directive in `internal/quotes/` referencing `../../data/quotes.json` breaks the build. Fix: embed anchor in `main.go`, pass `[]byte` to `quotes.Load`.
2. **lipgloss v1 vs v2 API mismatch** — the Renderer type was removed in v2; mixing import paths causes compile errors. Fix: pin `github.com/charmbracelet/lipgloss v1.1.0` in go.mod.
3. **Word-wrap ignores sidebar width** — wrapping at full terminal width then prepending `▌ ` (2 cells) overflows on narrow terminals. Fix: wrap at `termWidth - sidebarWidth - 1`; guard with 80-column fallback.
4. **Color profile detection fails when piped** — garbled ANSI in .bashrc if stdout is not a TTY. Fix: lipgloss v1 strips ANSI automatically when not a TTY; also check NO_COLOR explicitly.
5. **CGO_ENABLED=1 breaks cross-compile** — transitive cgo dependency fails cross-compile in CI. Fix: `CGO_ENABLED=0` in all Makefile build targets.

## Implications for Roadmap

Based on research, suggested phase structure:

### Phase 1: Data Foundation
**Rationale:** The Quote struct is the interface between packages; render cannot be written until it is stable. The go:embed path constraint is an architectural decision that must be resolved before any code is written — it determines where the embed anchor lives.
**Delivers:** go.mod with lipgloss v1.1.0, project layout, `data/quotes.json` with 15+ real quotes, `internal/quotes` package (Quote struct, Load, Random), embed anchor in main.go, full data layer test coverage.
**Addresses:** embedded quote data, 15+ quotes (both table stakes)
**Avoids:** go:embed path restriction (Pitfall 4), dot/underscore exclusion (Pitfall 3), blank import omission (Pitfall 11), directive/var blank line (Pitfall 12), math/rand v1 import (Pitfall 7)

### Phase 2: Render and Output
**Rationale:** With the Quote struct frozen and tested, render has a stable interface. Terminal width detection, word-wrap, and lipgloss styling are coupled and must be built together in one coherent pass.
**Delivers:** `internal/render` with lipgloss Estilo 4 sidebar styling, word-wrap at `termWidth - sidebarWidth`, color degradation, full render tests; thin `main.go` wiring everything; `go test ./...` green.
**Addresses:** word wrap, sidebar framing, color degradation, exit code 0 (all table stakes)
**Avoids:** sidebar width not subtracted (Pitfall 5), terminal width zero/error (Pitfall 6), lipgloss v1/v2 mismatch (Pitfall 1), color profile detection failure (Pitfall 2)

### Phase 3: Build and Distribution
**Rationale:** Cross-compile and Makefile targets are mechanical once the binary works. Platform smoke tests validate the tool end-to-end before shipping.
**Delivers:** Makefile with CGO_ENABLED=0 cross-compile targets (Linux amd64, macOS amd64+arm64, Windows amd64), dist/ artifacts, README in PT-BR.
**Avoids:** CGO_ENABLED=1 cross-compile failure (Pitfall 8), Windows ANSI VT raw escape codes (Pitfall 9)

### Phase Ordering Rationale

- Data before render: Quote struct must be stable; render package imports the type
- Render before build: cross-compile is pointless if the binary misbehaves
- go:embed path constraint resolved in Phase 1 before any other code exists — it is a build-time architectural constraint, not a runtime concern
- All table-stakes features fit within three phases; no artificial splitting needed for a project this size

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | lipgloss v1 vs v2 distinction verified from official release notes and upgrade guide; stdlib choices authoritative |
| Features | MEDIUM | fortune-tool conventions from multiple sources; lipgloss rendering from pkg.go.dev |
| Architecture | MEDIUM | go:embed restriction confirmed from official Go issue tracker; layout from official Go guide |
| Pitfalls | MEDIUM | go:embed pitfalls from official Go issues; lipgloss v1/v2 from charmbracelet upgrade guide |

**Overall confidence:** MEDIUM

### Gaps to Address

- **Quote content curation:** 15+ real Criminal Minds quotes must be sourced and curated — this is human curation work required before Phase 1 is complete.
- **Windows smoke test:** lipgloss v1 non-TTY fallback on cmd.exe has not been verified in this environment; validate during Phase 3.
- **Lipgloss border character rendering:** The `Border{Left: "▌"}` custom border pattern is derived from API docs; run a terminal smoke test early in Phase 2 to confirm correct cell width rendering.

## Sources

### Primary (MEDIUM-HIGH confidence)
- charmbracelet/lipgloss releases + Discussion #506 — v1 vs v2 API differences, color downsampling
- go.dev/blog/randv2 — math/rand/v2 auto-seeding, IntN API, Go 1.22 requirement
- pkg.go.dev/embed — embed package API and path restrictions
- go.dev/doc/modules/layout — official module layout guidance
- golang/go Issue #58519 — go:embed path restriction confirmation

### Secondary (MEDIUM confidence)
- lipgloss UPGRADE_GUIDE_V2.md — v1 to v2 migration details
- golang/go Issue #43854 — dot/underscore embed exclusion
- go.dev/doc/go1.26 — current stable toolchain confirmation

### Tertiary (LOW confidence)
- flaviocopes.com Go fortune clone tutorial — implementation pattern reference
- linux.die.net fortune(6) man page — fortune tool conventions

---
*Research completed: 2026-07-05*
*Ready for roadmap: yes*
