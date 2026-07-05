# Phase 2: Render and Output - Pattern Map

**Mapped:** 2026-07-05
**Files analyzed:** 3 new/modified files (+ 1 optional helper per RESEARCH discretion)
**Analogs found:** 3 / 3 (all sourced from Phase 1 codebase; lipgloss patterns from RESEARCH.md)

---

## File Classification

| New/Modified File | Role | Data Flow | Closest Analog | Match Quality |
|-------------------|------|-----------|----------------|---------------|
| `internal/render/render.go` | service | transform (Quote → string) | `internal/quotes/quotes.go` | exact |
| `internal/render/render_test.go` | test | transform (assert on return value) | `internal/quotes/quotes_test.go` | exact |
| `main.go` | entrypoint | request-response (stdin→stdout) | `main.go` (self — extend) | exact |
| `internal/render/wrap.go` *(optional)* | utility | transform (word-wrap) | `internal/quotes/quotes.go` | role-match |

---

## Pattern Assignments

### `internal/render/render.go` (service, transform)

**Analog:** `internal/quotes/quotes.go`

**Why this analog:** Same tier (`internal/`), same contract shape — pure functions that accept typed data and return results without I/O, printing, or `os.Exit`. Phase 1 established that internal packages return values; `main` owns side effects.

**Package declaration pattern** (quotes.go lines 1-6):
```go
package quotes

import (
	"encoding/json"
	"math/rand/v2"
)
```

**Apply to render.go:**
```go
package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
)
```

**Public API pattern** (quotes.go lines 8-17, 19-29):
```go
// Quote represents a single Criminal Minds quote record.
// Fields match the data/quotes.json schema exactly.
type Quote struct {
	Quote        string `json:"quote"`
	Author       string `json:"author"`
	Character    string `json:"character"`
	Season       int    `json:"season"`
	Episode      int    `json:"episode"`
	EpisodeTitle string `json:"episodeTitle"`
}

// Load parses JSON-encoded quote data and returns the slice.
// data should be the embedded bytes from main.go.
// Returns a non-nil error if JSON is malformed.
// Returns ([]Quote{}, nil) for a valid empty array — caller must check length.
func Load(data []byte) ([]Quote, error) {
```

**Render equivalent (locked by D-19):**
```go
// Quote formats q as Estilo 4 terminal output and returns the string.
// Does not print or detect TTY — caller passes colorEnabled from main.
// When colorEnabled is false, emits plain text with | sidebar (D-18) and no ANSI.
func Quote(q quotes.Quote, colorEnabled bool) string
```

**Core transform pattern** (quotes.go lines 23-29 — return, no side effects):
```go
func Load(data []byte) ([]Quote, error) {
	var qs []Quote
	if err := json.Unmarshal(data, &qs); err != nil {
		return nil, err
	}
	return qs, nil
}
```

**Apply:** `Quote` returns `string` only. No `fmt.Print*`, no `os` import for stdout/stderr. Unexported helpers (`buildAttribution`, `buildQuoteBlock`, `styleSidebar`) stay private.

**lipgloss style helpers** (no codebase analog — first lipgloss import; from RESEARCH.md Pattern 2):
```go
var (
	sidebarStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red — REND-01
	quoteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")) // bright — REND-03
	authorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow — REND-05
	dimStyle     = lipgloss.NewStyle().Faint(true)                      // muted meta
)

func styleSidebar(ch rune, enabled bool) string {
	s := string(ch)
	if !enabled {
		return s
	}
	return sidebarStyle.Render(s)
}
```

**Color branch pattern** (D-14/D-18 — skip lipgloss when disabled):
```go
sidebar := '▌'
if !colorEnabled {
	sidebar = '|'
}
// When !colorEnabled: zero Style.Render calls on any segment
```

**Attribution builder pattern** (D-01–D-13, from RESEARCH.md Pattern 4):
```go
func buildAttribution(q quotes.Quote) (authorPart, metaPart string) {
	switch {
	case q.Character != "" && q.Author != q.Character:
		authorPart = q.Author + ", cited by " + q.Character
	default:
		authorPart = q.Author
	}
	meta := " · Criminal Minds"
	if q.Season != 0 && q.Episode != 0 {
		meta += fmt.Sprintf(" · S%dE%d", q.Season, q.Episode)
		if q.EpisodeTitle != "" {
			meta += " (" + q.EpisodeTitle + ")"
		}
	}
	return authorPart, meta
}
```

**Per-line color split** (D-04 — color after wrap, not before):
```go
func renderAttributionLine(line string, colorEnabled bool) string {
	const marker = " · Criminal Minds"
	idx := strings.Index(line, marker)
	if !colorEnabled {
		return line
	}
	if idx < 0 {
		return authorStyle.Render(line)
	}
	return authorStyle.Render(line[:idx]) + dimStyle.Render(line[idx:])
}
```

**Assembly pattern** (D-09, D-19):
```go
func Quote(q quotes.Quote, colorEnabled bool) string {
	quoteLines := buildQuoteBlock(q.Quote, colorEnabled)
	author, meta := buildAttribution(q)
	attrLines := wrapAttribution(author+meta, blockWidth, attrIndent)

	var out strings.Builder
	out.WriteString(strings.Join(quoteLines, "\n"))
	out.WriteString("\n\n")
	// append attribution lines with per-line color split
	return out.String() // no trailing \n — main uses fmt.Println
}
```

**Anti-patterns (from 01-PATTERNS.md + RESEARCH):**
- Do NOT use `lipgloss.Style.Border()` for sidebar (D-20).
- Do NOT import `os` for TTY/NO_COLOR detection (D-14/D-15 — main owns that).
- Do NOT call `lipgloss.SetColorProfile` globally.
- Do NOT print from render — testability contract from CONTEXT.md §5.

---

### `internal/render/render_test.go` (test, transform)

**Analog:** `internal/quotes/quotes_test.go`

**Package declaration** (quotes_test.go lines 1-7):
```go
package quotes_test

import (
	"testing"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
)
```

**Apply:**
```go
package render_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
	"github.com/allanmedeiros71/criminalsay/internal/render"
)
```

**Table-driven test skeleton** (quotes_test.go lines 9-57):
```go
func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
		wantLen int
	}{
		{
			name:    "valid single quote",
			input:   []byte(`[{"quote":"Test","author":"A","character":"A","season":1,"episode":1,"episodeTitle":"Pilot"}]`),
			wantErr: false,
			wantLen: 1,
		},
		// ...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := quotes.Load(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && len(got) != tc.wantLen {
				t.Errorf("Load() len = %d, want %d", len(got), tc.wantLen)
			}
		})
	}
}
```

**Apply:** Use `[]struct { name string; q quotes.Quote; colorEnabled bool; wantContains []string; wantNotContains []string }` or separate test functions per requirement. Always `t.Run(tc.name, ...)`.

**Error assertion convention** (quotes_test.go lines 50-51):
```go
if (err != nil) != tc.wantErr {
	t.Fatalf("Load() error = %v, wantErr %v", err, tc.wantErr)
}
```

**Apply:** Use `t.Fatalf` for setup/precondition failures; `t.Errorf` for value/content mismatches.

**Membership check pattern** (quotes_test.go lines 60-77):
```go
func TestRandom_ReturnsItemFromSlice(t *testing.T) {
	qs := []quotes.Quote{
		{Quote: "A", Author: "X"},
		{Quote: "B", Author: "Y"},
		{Quote: "C", Author: "Z"},
	}
	q := quotes.Random(qs)
	found := false
	for _, item := range qs {
		if item.Quote == q.Quote {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Random() returned item not in input slice: %+v", q)
	}
}
```

**Apply to render:** Assert `strings.Contains(out, want)` / `!strings.Contains(out, unwanted)` for attribution segments, sidebar chars, episode codes.

**Test helper for COLOR-02/D-18** (from RESEARCH.md — no stdout capture):
```go
func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("expected no ANSI escapes, got: %q", s)
	}
}
```

**Wrap width assertion** (CONTEXT.md §7 + RESEARCH Pattern 3):
```go
for _, line := range strings.Split(out, "\n") {
	if line == "" {
		continue
	}
	if strings.Contains(line, "|") || strings.Contains(line, "▌") {
		if w := ansi.StringWidth(line); w > 52 {
			t.Errorf("line exceeds 52 display cols (got %d): %q", w, line)
		}
	}
}
```

**Required test cases (from CONTEXT.md §7 + 02-RESEARCH.md):**
| Test | Req IDs | Key assertion |
|------|---------|---------------|
| `TestQuote_noColor` | COLOR-02, D-18 | `assertNoANSI`; contains `\|`; no `▌` |
| `TestQuote_wrapWidth` | REND-02 | `ansi.StringWidth(line) <= 52` on quote lines |
| `TestQuote_externalAuthor` | D-01 | `"Marcus Aurelius, cited by Spencer Reid"` |
| `TestQuote_sameAuthor` | D-02 | author once, no "cited by" |
| `TestQuote_spacing` | REND-04 | `\n\n` between quote block and attribution |
| `TestQuote_sidebar` | REND-01, D-07 | every quote line has sidebar prefix |
| `TestQuote_color` | COLOR-01 | `colorEnabled=true` emits `\x1b[` |

**Critical rule:** Test `render.Quote(q, enabled)` return value only — never capture stdout (D-19, CONTEXT.md §5).

---

### `main.go` (entrypoint, request-response — modify)

**Analog:** `main.go` (self — extend Phase 1 wiring)

**Imports pattern** (main.go lines 3-9):
```go
import (
	_ "embed"
	"fmt"
	"os"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
)
```

**Phase 2 imports (add render + term; keep embed anchor):**
```go
import (
	_ "embed"
	"fmt"
	"os"

	term "github.com/charmbracelet/x/term"

	"github.com/allanmedeiros71/criminalsay/internal/quotes"
	"github.com/allanmedeiros71/criminalsay/internal/render"
)
```

**Embed anchor — DO NOT MOVE** (main.go lines 11-12):
```go
//go:embed data/quotes.json
var quoteData []byte
```

**Error handling pattern — preserve exactly** (main.go lines 14-23):
```go
func main() {
	qs, err := quotes.Load(quoteData)
	if err != nil {
		fmt.Fprintln(os.Stderr, "criminalsay: failed to load quotes:", err)
		os.Exit(1)
	}
	if len(qs) == 0 {
		fmt.Fprintln(os.Stderr, "criminalsay: no quotes available")
		os.Exit(1)
	}
```

**Replace stub only** (main.go lines 24-25):
```go
	// BEFORE (Phase 1 stub — remove):
	q := quotes.Random(qs)
	fmt.Println(q.Quote + " — " + q.Author)

	// AFTER (Phase 2):
	q := quotes.Random(qs)
	fmt.Println(render.Quote(q, colorEnabled()))
```

**New helper — color detection in main only** (D-14/D-15):
```go
func colorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return term.IsTerminal(os.Stdout.Fd())
}
```

**Printing convention:** Use `fmt.Println(render.Quote(...))` — render returns string without trailing `\n`; `Println` adds it (parity with Phase 1 stub).

**Anti-patterns:**
- Do NOT import `lipgloss` in main — styling stays in `internal/render`.
- Do NOT move `//go:embed` to `internal/`.
- Do NOT use builtin `println()`.
- Do NOT add CLI flags (out of scope).

---

### `internal/render/wrap.go` *(optional — Claude's discretion)*

**Analog:** `internal/quotes/quotes.go` (private helpers in same package)

If split from `render.go`, follow same package-level private function pattern as `Load`/`Random` separation of concerns. Use `github.com/charmbracelet/x/cellbuf` for wrap and `github.com/charmbracelet/x/ansi` for width measurement — no new go.mod entries (transitive via lipgloss).

**Constants** (D-05–D-08):
```go
const (
	blockWidth   = 52
	quoteIndent  = "  "
	attrIndent   = "       " // 7 spaces — CONTEXT.md §6.1 example
	sidebarGap   = " "
)
```

**Wrap-then-prefix pattern** (RESEARCH critical architectural choice):
```go
// 1. Word-wrap quote text at per-line budget (account for prefix + quote marks)
// 2. Prefix EACH line: quoteIndent + sidebar + sidebarGap
// 3. Place opening " on first line, closing " on last line (D-08)
```

---

## Shared Patterns

### Two-Tier Error Handling (from 01-PATTERNS.md — unchanged)

**Source:** `internal/quotes/quotes.go` lines 25-27 + `main.go` lines 16-18

**Apply to:** `main.go` only in Phase 2. `internal/render` has no error returns for `Quote` — malformed data is Phase 1 concern; render formats trusted struct fields.

```go
// internal — return errors (quotes only in Phase 2)
if err := json.Unmarshal(data, &qs); err != nil {
	return nil, err
}

// main — print and exit
if err != nil {
	fmt.Fprintln(os.Stderr, "criminalsay: failed to load quotes:", err)
	os.Exit(1)
}
```

### Render Returns String, Main Prints (CONTEXT.md §5, D-19)

**Source:** Phase 1 stub replaced; testability contract

**Apply to:** `internal/render/render.go`, `internal/render/render_test.go`, `main.go`

```go
// render.go — pure string output
func Quote(q quotes.Quote, colorEnabled bool) string { ... }

// main.go — sole stdout writer
fmt.Println(render.Quote(q, colorEnabled()))

// render_test.go — assert on return value, never os.Stdout
out := render.Quote(q, false)
```

### External Test Package (from 01-PATTERNS.md)

**Source:** `internal/quotes/quotes_test.go` line 1

**Apply to:** `internal/render/render_test.go`

```go
package render_test  // NOT package render
```

Tests only the exported `render.Quote` API.

### Table-Driven Tests (from 01-PATTERNS.md)

**Source:** `internal/quotes/quotes_test.go` lines 9-57

**Apply to:** All `render_test.go` cases — `[]struct{}` with `name` field, `t.Run(tc.name, ...)`.

### go:embed Anchor Placement (from 01-PATTERNS.md — unchanged)

**Source:** `main.go` lines 11-12

**Apply to:** `main.go` — do not touch embed anchor in Phase 2 except keeping it in place.

### Module Import Path

**Source:** `go.mod` line 1 + all existing imports

```go
"github.com/allanmedeiros71/criminalsay/internal/quotes"
"github.com/allanmedeiros71/criminalsay/internal/render"
```

### Color Gating Ownership (Phase 2 new — locked D-14/D-15)

**Source:** 02-CONTEXT.md D-14/D-15

| Concern | Owner | Pattern |
|---------|-------|---------|
| `NO_COLOR` env check | `main.go` | `os.Getenv("NO_COLOR") != ""` |
| TTY detection | `main.go` | `term.IsTerminal(os.Stdout.Fd())` |
| Apply/pass bool | `main.go` → `render` | `render.Quote(q, colorEnabled())` |
| ANSI styling | `internal/render` | lipgloss only when `colorEnabled=true` |
| Plain degradation | `internal/render` | `\|` sidebar, no `Style.Render` when false |

---

## No Analog Found

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| lipgloss styling calls | library integration | transform | First direct lipgloss import in codebase — use RESEARCH.md Pattern 2 |
| `ansi.StringWidth` in tests | test utility | transform | No prior usage — transitive dep, pattern from RESEARCH.md |
| `cellbuf.Wrap` in wrap helper | utility | transform | No prior usage — transitive dep, pattern from RESEARCH.md |

All three Phase 2 target files have exact analogs in Phase 1. Library-specific patterns (lipgloss, cellbuf, ansi) have no in-repo analog — planner must follow 02-RESEARCH.md verified API patterns.

---

## Metadata

**Analog search scope:** `/home/allan/git/criminalsay` — `main.go`, `internal/quotes/`, `go.mod`, `.planning/phases/01-data-foundation/01-PATTERNS.md`, `CONTEXT.md`
**Files scanned:** 4 Go source files
**Pattern extraction date:** 2026-07-05
**Pattern sources:** Phase 1 implemented code (primary), 02-RESEARCH.md library patterns (lipgloss/term/cellbuf/ansi), 02-CONTEXT.md locked decisions
