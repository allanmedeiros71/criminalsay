# Walking Skeleton — CriminalSay

**Phase:** 1
**Generated:** 2026-07-05

## Capability Proven End-to-End

Running `./criminalsay` prints one real Criminal Minds quote — selected at random from a compile-time embedded dataset — to stdout and exits 0.

## Architectural Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Language | Go 1.22+ (toolchain go1.25.6) | Single static binary, trivial cross-compile, no runtime required on target; math/rand/v2 requires 1.22+ minimum |
| Terminal styling | charmbracelet/lipgloss v1.1.0 | Auto-detects TTY/NO_COLOR/color depth from stdout with zero boilerplate; v1 is correct for a simple stdout CLI (v2 adds manual writer complexity that is irrelevant here) |
| Data layer | `data/quotes.json` compiled via `go:embed` | Offline-first by design; binary ships self-contained with no filesystem dependency at runtime; zero latency for data access |
| Random selection | `math/rand/v2` (stdlib) | Auto-seeded since Go 1.22 — no `rand.Seed()` call needed; `rand.IntN(n)` replaces deprecated `rand.Intn`; display use does not require crypto randomness |
| JSON parsing | `encoding/json` (stdlib) | Quote data is small and static; `json.Unmarshal` into `[]Quote` is correct and sufficient — no streaming or third-party parser needed |
| Binary entry point | `main.go` at repo root (not under `cmd/`) | Single binary with no exported packages; cmd/ nesting adds path depth with no benefit; official Go layout guide recommends cmd/ only for multi-binary or lib+binary modules |
| Directory layout | `main.go` + `internal/quotes/` + `internal/render/` + `data/` | `internal/` enforced by Go toolchain — prevents external imports; `data/` at root makes embed path `data/quotes.json` valid in `main.go` (Go forbids `..` in embed patterns, so data cannot be embedded from within `internal/`) |
| Module path | `github.com/allanmedeiros71/criminalsay` | Matches GitHub repository; standard Go module naming |
| Deployment | Local binary (`go build -o criminalsay .`); cross-compile in Phase 3 via Makefile | No server or container needed — CLI tool runs directly on user machine |

## Stack Touched in Phase 1

- [x] Project scaffold — `go.mod`, `go.sum`, module path, lipgloss dependency declared for Phase 2
- [x] Data layer — `data/quotes.json` (20 quotes) embedded at compile time; `quotes.Load()` and `quotes.Random()` implemented and tested
- [x] Entry point — `main.go` with embed anchor, empty-slice guard, Phase 1 stub output
- [x] Tests — `internal/quotes/quotes_test.go` with table-driven tests for Load and Random (DAT-01, DAT-02, DAT-03, CORE-01, BUILD-03)
- [ ] Rendering — Not in Phase 1; lipgloss styling added in Phase 2
- [ ] Distribution — Cross-compile Makefile and PT-BR README in Phase 3

## Out of Scope (Deferred to Later Slices)

- Terminal styling, color, sidebar layout (Phase 2 — `internal/render/`)
- Word-wrap and attribution line (Phase 2)
- `NO_COLOR` / non-TTY degradation (Phase 2 — lipgloss handles automatically)
- Cross-compile Makefile and `dist/` artifacts (Phase 3)
- PT-BR README (Phase 3)
- CLI flags (`--character`, `--season`) — deferred to v2
- Multiple visual themes — deferred to v2

## Subsequent Slice Plan

Each later phase adds one vertical slice on top of this skeleton without altering its architectural decisions:

- Phase 2: Running `criminalsay` prints a styled quote (red `▌` sidebar, white quote text, dim attribution line) and exits 0 — adds `internal/render/` using lipgloss v1.1.0; replaces `main.go` stub output with `render.Quote(q)`
- Phase 3: `make all` cross-compiles binaries for Linux amd64, macOS amd64/arm64, and Windows amd64 into `dist/`; adds PT-BR README with installation and usage instructions
