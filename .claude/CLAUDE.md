<!-- GSD:project-start source:PROJECT.md -->

## Project

**CriminalSay**

CriminalSay é uma ferramenta de linha de comando (CLI) em Go que exibe uma citação aleatória da série *Criminal Minds* no terminal com formatação visual elaborada — barra lateral vermelha, texto destacado, atribuição dimmed. Um comando, uma citação, uma saída bonita. Ideal para rodar no `.bashrc`/`.zshrc` a cada abertura de shell, no espírito do `fortune`.

**Core Value:** Rodar `criminalsay` e receber imediatamente uma citação memorável do BAU com estilo, sem interatividade e sem dependências externas.

### Constraints

- **Tech stack**: Go + charmbracelet/lipgloss — decisão travada (D2)
- **Dados**: JSON local embutido via go:embed — decisão travada (D5/D6)
- **Interatividade**: Nenhuma — escopo é exibir e sair (D4)
- **Aleatoriedade**: math/rand/v2 — sorteio simples, sem crypto (D7)
- **Plataformas**: Linux, macOS, Windows — cross-compile (D8)
- **Desempenho**: partida a frio <50ms típico (NFR-1)

<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->

## Technology Stack

## Recommended Stack

### Core Language

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go | 1.22+ (use 1.26 toolchain) | Language runtime | math/rand/v2 requires 1.22+; 1.26 is current stable as of Feb 2026 |

### Terminal Styling

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| charmbracelet/lipgloss | v1.1.0 | Sidebar border, colors, layout | v1.1.0 auto-detects TTY/NO_COLOR/color depth from stdout with zero boilerplate. v2 shifts this burden to the caller (requires `lipgloss.Println()` writer pattern). For a single-binary CLI that just prints to stdout, v1 is correct; the added complexity of v2 is unnecessary and its I/O flexibility is irrelevant here. |

### Standard Library (no external deps)

| Package | Version | Purpose | Why |
|---------|---------|---------|-----|
| `embed` | Go 1.16+ (stdlib) | Embed quotes.json into binary | Standard library directive; zero external dependency; binary ships self-contained |
| `math/rand/v2` | Go 1.22+ (stdlib) | Select random quote | Auto-seeded, no `rand.Seed()` call needed, `IntN(n)` replaces deprecated `Intn(n)`; stdlib means no go.mod entry |
| `encoding/json` | stdlib | Decode embedded JSON | Standard; no third-party JSON library needed for a small static dataset |
| `os` | stdlib | Check `os.Stdout`, exit codes | TTY detection falls to lipgloss; `os` only needed for `os.Exit(0)` in main |

### Build Tooling

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go toolchain | 1.26 | `go build`, `go test`, cross-compile | Built-in `GOOS`/`GOARCH` env vars handle all three target platforms without external tools |
| Makefile | — | Developer convenience targets | `make build`, `make test`, `make release` wrap the `go build` invocations with ldflags for version injection. Simple enough to not need GoReleaser for v1. |

## go.mod

## Project Layout

- `main.go` at root (not under `cmd/`) because there is exactly one binary and no exported packages. The official Go docs recommend `cmd/` only when the module has both library code AND binaries, or multiple binaries.
- `internal/` enforced by the Go toolchain — prevents any external code from importing `quotes` or `render` packages. Free refactoring guarantee.
- `data/` next to `internal/` makes the embed path predictable: `//go:embed ../../data/quotes.json` from inside `internal/quotes/`, or move the embed directive to `main.go` — simpler to put the embed in the `quotes` package itself with a relative path.
- No `pkg/`, no `api/`, no `cmd/` — all three would be cargo-culted from server project templates. Reject them.

### go:embed pattern for quotes.json

### math/rand/v2 pattern for selection

### lipgloss v1 sidebar pattern

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Styling | lipgloss v1.1.0 | lipgloss v2 | v2 requires manual writer setup; complexity not justified for simple stdout CLI |
| Styling | lipgloss v1.1.0 | fatih/color | Redundant; lipgloss covers all coloring needs; two color systems would conflict |
| Styling | lipgloss v1.1.0 | termenv directly | lipgloss is built on termenv; using termenv directly is lower-level than needed |
| Randomness | math/rand/v2 (stdlib) | github.com/golang/exp/rand | Unnecessary external dep; v2 is now in stdlib since Go 1.22 |
| Data embed | go:embed (stdlib) | go-bindata, pkger | Obsolete tools; go:embed is the official stdlib solution since Go 1.16 |
| Build | Makefile | GoReleaser | GoReleaser overhead not justified until distribution pipeline needed (v2) |
| Layout | flat (main.go at root) | cmd/criminalsay/main.go | cmd/ nesting adds path depth with no benefit for single-binary module |

## Installation (developer)

# Install lipgloss

# Build current platform

# Cross-compile (Makefile targets)

# Test

## Sources

- [charmbracelet/lipgloss releases](https://github.com/charmbracelet/lipgloss/releases) — v1.1.0 confirmed March 2025, v2.0.5 confirmed July 2024
- [Lip Gloss v2: What's New (Discussion #506)](https://github.com/charmbracelet/lipgloss/discussions/506) — v2 color downsampling is manual; compat package for migration
- [lipgloss pkg.go.dev v1](https://pkg.go.dev/github.com/charmbracelet/lipgloss) — v1.1.0 import path confirmed
- [Evolving the Go Standard Library with math/rand/v2](https://go.dev/blog/randv2) — Go 1.22 required, auto-seeded, IntN() API
- [Organizing a Go module](https://go.dev/doc/modules/layout) — official layout guidance: cmd/ only for multi-binary or lib+binary modules
- [Go 1.26 Release Notes](https://go.dev/doc/go1.26) — current stable toolchain as of February 2026
- [embed package](https://pkg.go.dev/embed) — stdlib since Go 1.16
- [fatih/color](https://github.com/fatih/color) — v1.19.0 current; respects NO_COLOR; lightweight but redundant alongside lipgloss

<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->

## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->

## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->

## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, `.github/skills/`, or `.codex/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->

## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:

- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->

<!-- GSD:profile-start -->

## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
