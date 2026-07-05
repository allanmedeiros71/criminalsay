# Architecture Patterns

**Project:** CriminalSay
**Researched:** 2026-07-05
**Confidence:** MEDIUM (Go stdlib + lipgloss API verified via pkg.go.dev; patterns confirmed against go.dev official layout guide)

---

## Recommended Architecture

A three-layer pipeline: **embed → decode → select → render → print**.

```
criminalsay/
├── main.go                    # package main — orchestrator only
├── go.mod
├── go.sum
├── data/
│   └── quotes.json            # embedded at compile time; never shipped separately
├── internal/
│   ├── quotes/
│   │   ├── quotes.go          # Quote type, Load(), Random()
│   │   └── quotes_test.go
│   └── render/
│       ├── render.go          # Quote() render function
│       └── render_test.go
├── scripts/
│   └── build.sh               # cross-compile for linux/darwin/windows
├── .gitignore
├── LICENSE
└── README.md
```

---

## Component Boundaries

| Component | Package Path | Responsibility | Talks To |
|-----------|-------------|----------------|----------|
| **Embed anchor** | `internal/quotes` | Holds the `//go:embed` directive and raw bytes | Nothing (data only) |
| **quotes** | `internal/quotes` | Defines `Quote` struct; `Load()` parses JSON; `Random()` selects uniformly | `encoding/json`, `math/rand/v2`, `embed` |
| **render** | `internal/render` | Builds lipgloss styles; `Quote()` produces formatted string | `github.com/charmbracelet/lipgloss`, `internal/quotes.Quote` (type only) |
| **main** | `main` (root) | Calls Load, Random, detects color capability, calls render.Quote, prints | `internal/quotes`, `internal/render`, `os`, `fmt` |

**Dependency graph (no cycles):**

```
main
 ├── internal/quotes   (Load, Random, Quote type)
 └── internal/render   (Quote func)
       └── internal/quotes.Quote  (type reference only)
```

`render` depends on the `Quote` type from `quotes` — this is the only cross-internal dependency. It is one-directional and does not create a cycle.

---

## Data Flow

```
compile time
  data/quotes.json
       │  //go:embed data/quotes.json
       ▼
  var quotesData []byte          (in internal/quotes/quotes.go)

runtime
  quotes.Load()
       │  json.Unmarshal(quotesData, &qs)
       ▼
  []Quote                        (slice of structs)
       │
  quotes.Random(qs)
       │  qs[rand.IntN(len(qs))]
       ▼
  Quote                          (single selected struct)
       │
  color detection (main.go)
       │  os.Getenv("NO_COLOR"), os.Stdout fd check
       ▼
  render.Quote(q, colorEnabled)
       │  lipgloss styles applied, wrapping computed
       ▼
  string                         (formatted, ANSI or plain)
       │
  fmt.Println(...)               (in main.go)
       ▼
  stdout → terminal
```

No goroutines, no channels, no shared state. Entirely linear, top-to-bottom.

---

## Contract Signatures

### `internal/quotes/quotes.go`

```go
package quotes

import (
    _ "embed"
    "encoding/json"
    "fmt"
    "math/rand/v2"
)

//go:embed ../../data/quotes.json
var quotesData []byte

// Quote is the canonical data type shared with render.
type Quote struct {
    Quote        string `json:"quote"`
    Author       string `json:"author"`
    Character    string `json:"character"`     // empty string if unknown
    Season       int    `json:"season"`
    Episode      int    `json:"episode"`
    EpisodeTitle string `json:"episodeTitle"`
}

// Load parses the embedded JSON and returns all quotes.
// Returns an error if the JSON is malformed or the slice is empty.
func Load() ([]Quote, error)

// Random returns a uniformly random element from qs.
// Panics if qs is empty — callers must check Load() result length.
func Random(qs []Quote) Quote {
    return qs[rand.IntN(len(qs))]
}
```

**Embed path note:** The `//go:embed` directive in `internal/quotes/quotes.go` must reference the path relative to that file's location. With `data/quotes.json` at the repo root and `quotes.go` two levels down, the directive is `//go:embed ../../data/quotes.json`. Alternatively, place the embed anchor in `main.go` and pass the bytes (or the loaded slice) into `quotes.Load(data []byte)` — this is slightly cleaner because go:embed paths are relative to the file containing the directive.

**Recommended approach:** Place the `//go:embed` directive in `main.go` (or a dedicated `embed.go` at root) and pass raw bytes to `quotes.Load`:

```go
// main.go or embed.go (package main)
//go:embed data/quotes.json
var quotesData []byte
```

```go
// internal/quotes/quotes.go
func Load(data []byte) ([]Quote, error)
```

This keeps the embed path trivial (`data/quotes.json` relative to root) and keeps `internal/quotes` free of filesystem path knowledge.

---

### `internal/render/render.go`

```go
package render

import (
    "github.com/charmbracelet/lipgloss"
    "github.com/allanmedeiros71/criminalsay/internal/quotes"
)

// Quote renders q using the "Terminal cru" Style 4 and returns
// the formatted string. Does not print.
// If colorEnabled is false, no ANSI sequences are emitted.
func Quote(q quotes.Quote, colorEnabled bool) string
```

Internally, `render.Quote` builds two lipgloss styles:

1. **Sidebar style** — left border only using a custom `Border` struct:
   ```go
   sidebarBorder := lipgloss.Border{Left: "▌"}
   sidebarStyle := lipgloss.NewStyle().
       BorderStyle(sidebarBorder).
       BorderLeft(true).
       BorderRight(false).
       BorderTop(false).
       BorderBottom(false).
       BorderForeground(lipgloss.Color("9")).   // ANSI red
       PaddingLeft(1).
       Width(52)   // wraps text at 52 columns
   ```

2. **Attribution style** — dim text, indented:
   ```go
   attrStyle := lipgloss.NewStyle().
       Faint(true).
       PaddingLeft(7)
   ```

When `colorEnabled` is false, pass a `lipgloss.NewRenderer(os.Stdout)` constructed with `renderer.SetHasDarkBackground(false)` or simply skip all color calls — lipgloss respects `NO_COLOR` automatically when using the default renderer, but the `colorEnabled` param allows explicit override for pipes/non-TTY.

---

### `main.go`

```go
func main() {
    qs, err := quotes.Load(quotesData)
    if err != nil || len(qs) == 0 {
        fmt.Fprintln(os.Stderr, "criminalsay: failed to load quotes:", err)
        os.Exit(1)
    }
    q := quotes.Random(qs)
    colorEnabled := isColorTerminal()
    fmt.Println(render.Quote(q, colorEnabled))
}

func isColorTerminal() bool {
    if os.Getenv("NO_COLOR") != "" {
        return false
    }
    fi, err := os.Stdout.Stat()
    if err != nil {
        return false
    }
    return (fi.Mode() & os.ModeCharDevice) != 0
}
```

`main` is intentionally thin: no business logic, no formatting decisions. Its only jobs are orchestration and error handling.

---

## Build Order (Phase Dependencies)

Build in this order — each step depends on the previous:

```
1. data/quotes.json          — define schema, populate ≥15 quotes
2. internal/quotes/quotes.go — Quote type + Load + Random (no external deps beyond stdlib)
3. internal/quotes/quotes_test.go — verify parsing and selection before render exists
4. internal/render/render.go — depends on quotes.Quote type being stable
5. internal/render/render_test.go — verify formatting contract
6. main.go                   — wire everything; thin orchestrator
7. scripts/build.sh          — cross-compile only after all tests pass
```

Do not write render code before the `Quote` struct is frozen. The struct shape (field names and types) is the interface between the two packages.

---

## Testing Approach Per Component

### `internal/quotes` — pure data layer, easy to test

```go
func TestLoad_ParsesWithoutError(t *testing.T)
func TestLoad_MinimumCount(t *testing.T)         // len >= 15
func TestLoad_AllFieldsNonEmpty(t *testing.T)    // spot-check required fields
func TestRandom_ReturnsFromSlice(t *testing.T)   // result is element of input
func TestRandom_Distribution(t *testing.T)       // run 1000 times, no index always chosen
```

No mocking needed — use the real embedded bytes. `quotes.Load` is a pure function of its `[]byte` input, so tests can also supply crafted JSON for edge cases.

### `internal/render` — string output, no I/O

```go
func TestQuote_ContainsText(t *testing.T)           // output contains q.Quote substring
func TestQuote_ContainsAuthor(t *testing.T)         // output contains q.Author
func TestQuote_SidebarPresentOnEveryLine(t *testing.T) // each non-blank line has ▌
func TestQuote_NoANSIWhenColorDisabled(t *testing.T)   // no \x1b[ in plain output
func TestQuote_LineWidthRespected(t *testing.T)        // no line exceeds width + border overhead
```

Since `render.Quote` returns a string, all assertions operate on the string — no stdout capture required. This is why separating production from printing matters.

### `main` — integration smoke test only

No unit tests on `main` itself. The cross-platform smoke test (per CONTEXT.md section 7) is manual: build and run on Linux, macOS, Windows Terminal; confirm border characters and color appear correctly.

---

## Patterns to Follow

### Pattern 1: Embed anchor in package main

**What:** Declare `//go:embed data/quotes.json` in `main.go` (or a dedicated `embed.go` at root), not inside `internal/quotes/`.

**When:** Always for this project.

**Why:** The `//go:embed` path is resolved relative to the source file containing the directive. Placing it at root keeps the path trivially `data/quotes.json`. If placed inside `internal/quotes/quotes.go` the path becomes `../../data/quotes.json` — legal but fragile under directory refactors.

```go
// embed.go  (package main)
package main

import _ "embed"

//go:embed data/quotes.json
var quotesData []byte
```

### Pattern 2: Render returns string, main prints

**What:** `render.Quote` returns `string`, never calls `fmt.Print`.

**When:** All rendering functions.

**Why:** Allows table-driven tests on render output without redirecting `os.Stdout`. Also allows callers to compose output (e.g., add a newline, prefix a header) without touching render internals.

### Pattern 3: Thin main, all logic in internal

**What:** `main()` is at most 15-20 lines. No parsing, styling, or data logic in `main`.

**When:** Always.

**Why:** Go's `package main` cannot be imported by tests. Any logic in `main.go` is untestable. The split between main and internal packages is the primary testability boundary in Go CLI tools.

---

## Anti-Patterns to Avoid

### Anti-Pattern 1: Global `lipgloss.Style` variables initialized at package level

**What:** Declaring `var sidebarStyle = lipgloss.NewStyle().BorderLeft(true)...` as a package-level var.

**Why bad:** lipgloss styles capture the renderer at initialization time. If the renderer's color profile changes (e.g., `NO_COLOR` detected after import), the cached style won't reflect it. Also makes testing with `colorEnabled=false` harder.

**Instead:** Build styles inside `render.Quote` based on the `colorEnabled` argument, or build two style sets and select at render time.

### Anti-Pattern 2: Calling `json.Unmarshal` inside `Random`

**What:** Re-parsing the JSON on every call to avoid storing the slice.

**Why bad:** Unmarshal is not free. For a startup tool the slice should be parsed once in `Load` and passed around.

**Instead:** `Load` parses once, returns `[]Quote`. `Random` receives the slice.

### Anti-Pattern 3: Panicking in `Random` without caller guard in main

**What:** `Random` panics on empty slice (correct behavior for a logic error), but `main` never checks `len(qs) > 0`.

**Why bad:** A corrupted or empty JSON produces a panic with no user-friendly message.

**Instead:** `main` checks `len(qs) == 0` after `Load` and exits with a clear error message.

### Anti-Pattern 4: Using `math/rand` (v1) instead of `math/rand/v2`

**What:** `import "math/rand"` with manual `rand.Seed(time.Now().UnixNano())`.

**Why bad:** v1 requires explicit seeding; forgetting it produces the same sequence every run. v2 auto-seeds and has a cleaner API (`IntN` vs `Intn`).

**Instead:** `import "math/rand/v2"` — no seeding needed.

---

## Scalability Considerations

This tool is not a server. Scalability is not a concern. The relevant NFR is startup latency.

| Concern | Current scale (1 invocation) | Notes |
|---------|------------------------------|-------|
| JSON parse time | <1ms for ~20 quotes | Negligible; embedded bytes, no I/O |
| lipgloss render time | <1ms | String operations only |
| Total cold start | Target <50ms | Go runtime init dominates; no network, no disk |
| Adding quotes | Edit JSON, rebuild | NFR-5: zero code changes for new quotes |

If the quote database grows to thousands of entries, `json.Unmarshal` into `[]Quote` remains fast enough for a CLI startup context. No optimization needed at any realistic scale for this use case.

---

## Sources

- [go.dev — Organizing a Go module](https://go.dev/doc/modules/layout) — official layout guide (LOW confidence, webfetch)
- [pkg.go.dev/embed](https://pkg.go.dev/embed) — embed package API (LOW confidence, webfetch)
- [pkg.go.dev/math/rand/v2](https://pkg.go.dev/math/rand/v2) — rand v2 API (LOW confidence, webfetch)
- [pkg.go.dev/github.com/charmbracelet/lipgloss](https://pkg.go.dev/github.com/charmbracelet/lipgloss) — lipgloss API (LOW confidence, webfetch)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout) — community layout reference (LOW confidence, websearch)
- CONTEXT.md sections 6, 6.1, 7 — project-defined contracts and style spec (authoritative for this project)
