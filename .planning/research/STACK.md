# Technology Stack

**Project:** CriminalSay
**Researched:** 2026-07-05

## Recommended Stack

### Core Language

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go | 1.22+ (use 1.26 toolchain) | Language runtime | math/rand/v2 requires 1.22+; 1.26 is current stable as of Feb 2026 |

Go 1.22 is the hard floor because `math/rand/v2` entered the standard library there. Targeting Go 1.22 in `go.mod` maximizes compatibility with users who have not yet upgraded to 1.26, while still unlocking all needed stdlib features. The developer toolchain can (and should) be 1.26.

### Terminal Styling

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| charmbracelet/lipgloss | v1.1.0 | Sidebar border, colors, layout | v1.1.0 auto-detects TTY/NO_COLOR/color depth from stdout with zero boilerplate. v2 shifts this burden to the caller (requires `lipgloss.Println()` writer pattern). For a single-binary CLI that just prints to stdout, v1 is correct; the added complexity of v2 is unnecessary and its I/O flexibility is irrelevant here. |

**Do NOT use v2.** The v2 change that made color downsampling manual was motivated by Bubble Tea TUI conflicts — a problem that does not exist for a simple print-and-exit CLI. Using v1.1.0 means: adaptive color out of the box, `NO_COLOR` respected automatically, ANSI stripped when piped, nothing extra to write.

**Do NOT add fatih/color.** The PROJECT.md lists it as a candidate, but it is redundant when lipgloss is already in scope. fatih/color is useful when you need colored `fmt.Printf`-style calls without any layout. Since every styled element in CriminalSay is a lipgloss `Render()` call (sidebar bar + quote text + attribution), fatih/color adds a second color-management system with no benefit. One dependency, one color system.

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

**GoReleaser is explicitly deferred to v2.** For a single-person project with a small binary, GoReleaser adds a YAML config and CI pipeline that is not worth the setup cost in the first release. A Makefile with three `GOOS`/`GOARCH` calls produces identical artifacts.

## go.mod

```
module github.com/allanmedeiros71/criminalsay

go 1.22

require github.com/charmbracelet/lipgloss v1.1.0

require (
    github.com/charmbracelet/colorprofile v0.1.9 // indirect
    github.com/charmbracelet/x/ansi v0.8.0 // indirect
    github.com/charmbracelet/x/cellbuf v0.0.13 // indirect
    github.com/charmbracelet/x/term v0.2.1 // indirect
    github.com/muesli/termenv v0.15.2 // indirect
    github.com/rivo/uniseg v0.4.7 // indirect
)
```

The indirect dependencies are pulled automatically by `go mod tidy` when lipgloss v1.1.0 is added. The exact patch versions in the indirect block will be resolved by the toolchain; the list above is illustrative — run `go mod tidy` after `go get github.com/charmbracelet/lipgloss@v1.1.0`.

## Project Layout

```
criminalsay/
  go.mod
  go.sum
  main.go                  ← package main; calls render.Print(quotes.Pick())
  data/
    quotes.json            ← embedded via //go:embed in internal/quotes/quotes.go
  internal/
    quotes/
      quotes.go            ← Quote struct, embed directive, Pick() func
      quotes_test.go
    render/
      render.go            ← lipgloss styles, Print(q Quote) func
      render_test.go
  Makefile
  README.md
```

**Why this layout:**

- `main.go` at root (not under `cmd/`) because there is exactly one binary and no exported packages. The official Go docs recommend `cmd/` only when the module has both library code AND binaries, or multiple binaries.
- `internal/` enforced by the Go toolchain — prevents any external code from importing `quotes` or `render` packages. Free refactoring guarantee.
- `data/` next to `internal/` makes the embed path predictable: `//go:embed ../../data/quotes.json` from inside `internal/quotes/`, or move the embed directive to `main.go` — simpler to put the embed in the `quotes` package itself with a relative path.
- No `pkg/`, no `api/`, no `cmd/` — all three would be cargo-culted from server project templates. Reject them.

### go:embed pattern for quotes.json

```go
// internal/quotes/quotes.go
package quotes

import (
    _ "embed"
    "encoding/json"
)

//go:embed quotes.json
var raw []byte

type Quote struct {
    Text      string `json:"text"`
    Character string `json:"character"`
}

func Load() ([]Quote, error) {
    var qs []Quote
    return qs, json.Unmarshal(raw, &qs)
}
```

Store `quotes.json` at `internal/quotes/quotes.json` (co-located with the Go file that embeds it). This keeps the embed path to a single filename with no `../` traversal, which is cleaner and avoids the Go embed restriction on paths containing `..`.

### math/rand/v2 pattern for selection

```go
import "math/rand/v2"

func Pick(qs []Quote) Quote {
    return qs[rand.IntN(len(qs))]
}
```

No seed, no global state setup. `rand.IntN` is auto-seeded with a cryptographically random seed at program start.

### lipgloss v1 sidebar pattern

```go
import "github.com/charmbracelet/lipgloss"

var (
    bar = lipgloss.NewStyle().
        Foreground(lipgloss.Color("9")). // ANSI red, degrades to 16-color correctly
        SetString("▌ ")

    quoteStyle = lipgloss.NewStyle().
        PaddingLeft(2)

    attrStyle = lipgloss.NewStyle().
        PaddingLeft(2).
        Faint(true)
)
```

Use ANSI color index `"9"` (bright red) rather than a hex value. lipgloss will downsample `"9"` correctly across truecolor, 256-color, and 16-color terminals, whereas a hex value requires downsampling logic that v1 handles by mapping to the nearest palette entry.

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

```bash
# Install lipgloss
go get github.com/charmbracelet/lipgloss@v1.1.0
go mod tidy

# Build current platform
go build -o criminalsay .

# Cross-compile (Makefile targets)
GOOS=linux   GOARCH=amd64  go build -o dist/criminalsay-linux-amd64 .
GOOS=darwin  GOARCH=amd64  go build -o dist/criminalsay-darwin-amd64 .
GOOS=darwin  GOARCH=arm64  go build -o dist/criminalsay-darwin-arm64 .
GOOS=windows GOARCH=amd64  go build -o dist/criminalsay-windows-amd64.exe .

# Test
go test ./...
```

## Sources

- [charmbracelet/lipgloss releases](https://github.com/charmbracelet/lipgloss/releases) — v1.1.0 confirmed March 2025, v2.0.5 confirmed July 2024
- [Lip Gloss v2: What's New (Discussion #506)](https://github.com/charmbracelet/lipgloss/discussions/506) — v2 color downsampling is manual; compat package for migration
- [lipgloss pkg.go.dev v1](https://pkg.go.dev/github.com/charmbracelet/lipgloss) — v1.1.0 import path confirmed
- [Evolving the Go Standard Library with math/rand/v2](https://go.dev/blog/randv2) — Go 1.22 required, auto-seeded, IntN() API
- [Organizing a Go module](https://go.dev/doc/modules/layout) — official layout guidance: cmd/ only for multi-binary or lib+binary modules
- [Go 1.26 Release Notes](https://go.dev/doc/go1.26) — current stable toolchain as of February 2026
- [embed package](https://pkg.go.dev/embed) — stdlib since Go 1.16
- [fatih/color](https://github.com/fatih/color) — v1.19.0 current; respects NO_COLOR; lightweight but redundant alongside lipgloss
