# Phase 1: Data Foundation - Pattern Map

**Mapped:** 2026-07-05
**Files analyzed:** 5 new files
**Analogs found:** 0 / 5 (greenfield — all patterns sourced from Go stdlib and project decisions)

---

## File Classification

| New File | Role | Data Flow | Closest Analog | Match Quality |
|----------|------|-----------|----------------|---------------|
| `go.mod` | config | — | none (greenfield) | no analog |
| `data/quotes.json` | data | batch (compile-time embed) | none (greenfield) | no analog |
| `internal/quotes/quotes.go` | service | transform (JSON → struct) | none (greenfield) | no analog |
| `internal/quotes/quotes_test.go` | test | batch | none (greenfield) | no analog |
| `main.go` | entrypoint | request-response (stdin→stdout) | none (greenfield) | no analog |

---

## Pattern Assignments

### `go.mod` (config)

**Analog:** none — greenfield. Use Go module reference pattern.

**Pattern (go.mod content):**
```
module github.com/allanmedeiros71/criminalsay

go 1.22

toolchain go1.25.6

require github.com/charmbracelet/lipgloss v1.1.0
```

**Key rules:**
- `go 1.22` minimum is required for `math/rand/v2`.
- `toolchain go1.25.6` matches the installed toolchain. If omitted, `go mod tidy` sets it automatically from the running toolchain.
- lipgloss is added here now (Phase 1) even though it is NOT imported by any Phase 1 `.go` file. Phase 2 will import it without modifying go.mod.
- Do NOT add `fatih/color` or any other styling library.

**Bootstrap commands (run in order after all .go files are written):**
```bash
go mod init github.com/allanmedeiros71/criminalsay
go get github.com/charmbracelet/lipgloss@v1.1.0
go mod tidy
```

---

### `data/quotes.json` (data — compile-time embed target)

**Analog:** none — greenfield. Schema is locked by DAT-02 and decisions D-01/D-02.

**JSON array schema:**
```json
[
  {
    "quote": "Quote text.",
    "author": "Original Author Name",
    "character": "CM Character Name",
    "season": 1,
    "episode": 1,
    "episodeTitle": "Episode Title"
  }
]
```

**Field semantics (locked):**
- When CM character IS the originator (D-01): `author` = `character` = character name.
- When CM character CITES someone (D-02): `author` = original author, `character` = CM character who cited it.
- `character` and `episodeTitle` may be empty string `""` — they are optional display fields (DAT-03).
- `season` and `episode` are ints; use `0` only if truly unknown.

**JSON tag names must exactly match the Go struct tags** (camelCase):
- `"quote"`, `"author"`, `"character"`, `"season"`, `"episode"`, `"episodeTitle"`

**Minimum content:** ≥15 quotes; up to 25. Spread across Reid, Hotch, Morgan, Rossi, Garcia, Prentiss, JJ (D-03). Pre-populated 20-quote set is provided in RESEARCH.md Code Examples section — copy it verbatim as the initial dataset. User reviews for accuracy post-phase.

---

### `internal/quotes/quotes.go` (service, transform)

**Analog:** none — greenfield. Pattern from Go stdlib (`encoding/json`, `math/rand/v2`).

**Imports pattern:**
```go
import (
    "encoding/json"
    "math/rand/v2"
)
```

**Struct pattern (locked by DAT-02 / CONTEXT.md section 6):**
```go
type Quote struct {
    Quote        string `json:"quote"`
    Author       string `json:"author"`
    Character    string `json:"character"`
    Season       int    `json:"season"`
    Episode      int    `json:"episode"`
    EpisodeTitle string `json:"episodeTitle"`
}
```

**Load function pattern:**
```go
// Load parses JSON-encoded quote data and returns the slice.
// data should be the embedded bytes from main.go.
// Returns a non-nil error if JSON is malformed.
// Returns ([]Quote{}, nil) for a valid empty array — caller must check length.
func Load(data []byte) ([]Quote, error) {
    var qs []Quote
    if err := json.Unmarshal(data, &qs); err != nil {
        return nil, err
    }
    return qs, nil
}
```

**Random function pattern:**
```go
// Random returns a uniformly random element from qs.
// Panics if qs is empty — caller must ensure len(qs) > 0.
func Random(qs []Quote) Quote {
    return qs[rand.IntN(len(qs))]
}
```

**Critical API constraints:**
- `Load` receives `[]byte` from caller (main.go) — it does NOT embed data itself. The go:embed anchor cannot be placed inside `internal/` due to the `..` path restriction (go doc embed: "Patterns may not contain '.' or '..'").
- `Random` uses `rand.IntN(n)` (capital N) — this is the `math/rand/v2` API. `rand.Intn` (lowercase n) is the deprecated v1 API and does not exist in v2.
- No `rand.Seed()` or `rand.New()` — `math/rand/v2` global source is auto-seeded since Go 1.22.
- `Random` panics on empty slice — this is correct; empty dataset is a programming error, not a user error.

---

### `internal/quotes/quotes_test.go` (test)

**Analog:** none — greenfield. Pattern: Go stdlib `testing`, table-driven tests, external test package.

**Package declaration:**
```go
package quotes_test
```
Use `quotes_test` (external test package) — tests only the exported API.

**Imports pattern:**
```go
import (
    "testing"

    "github.com/allanmedeiros71/criminalsay/internal/quotes"
)
```

**Table-driven test pattern for Load:**
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
        {
            name:    "empty array",
            input:   []byte(`[]`),
            wantErr: false,
            wantLen: 0,
        },
        {
            name:    "malformed JSON",
            input:   []byte(`not json`),
            wantErr: true,
        },
        {
            name:    "missing optional fields",
            input:   []byte(`[{"quote":"Q","author":"A","character":"","season":1,"episode":1,"episodeTitle":""}]`),
            wantErr: false,
            wantLen: 1,
        },
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

**Random function tests — membership and distribution:**
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

func TestRandom_Distribution(t *testing.T) {
    qs := []quotes.Quote{
        {Quote: "A"}, {Quote: "B"}, {Quote: "C"},
        {Quote: "D"}, {Quote: "E"}, {Quote: "F"},
    }
    seen := make(map[string]bool)
    for i := 0; i < 200; i++ {
        q := quotes.Random(qs)
        seen[q.Quote] = true
    }
    if len(seen) < 2 {
        t.Errorf("Random() appears non-random: only %d unique values in 200 calls", len(seen))
    }
}
```

**Requirements coverage:**
- `TestLoad/valid_single_quote` → DAT-01, DAT-02
- `TestLoad/missing_optional_fields` → DAT-03
- `TestLoad/empty_array` → guards Pitfall 3 (rand.IntN(0) panic)
- `TestLoad/malformed_JSON` → error path
- `TestRandom_*` → CORE-01

---

### `main.go` (entrypoint, request-response)

**Analog:** none — greenfield. Pattern from go:embed stdlib + internal/quotes API.

**Imports pattern:**
```go
import (
    _ "embed"
    "fmt"
    "os"

    "github.com/allanmedeiros71/criminalsay/internal/quotes"
)
```

**Embed anchor pattern (MUST appear in main.go, not in internal/):**
```go
//go:embed data/quotes.json
var quoteData []byte
```

Critical: The `//go:embed` comment must immediately precede the `var` declaration with no blank lines between. The blank `_ "embed"` import is mandatory even though `embed.FS` is never referenced.

**Core pattern (Phase 1 stub — replaced entirely in Phase 2):**
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
    q := quotes.Random(qs)
    fmt.Println(q.Quote + " — " + q.Author)
}
```

**Error handling pattern:**
- Load errors → `fmt.Fprintln(os.Stderr, ...)` + `os.Exit(1)`.
- Empty dataset → explicit length guard before calling `Random` (Pitfall 3: `rand.IntN(0)` panics).
- No panics in main() — all error paths exit cleanly with a message.

**Anti-patterns to avoid in main.go:**
- Do NOT import `github.com/charmbracelet/lipgloss` in Phase 1 source (it is in go.mod for Phase 2, not yet used).
- Do NOT use `println()` (builtin) — use `fmt.Println()` for portable output.
- Do NOT place `//go:embed` in any file under `internal/`.

---

## Shared Patterns

### Error Handling (applies to main.go and quotes.go)

**Pattern:** Two-tier error strategy.
- `internal/quotes`: return errors to caller, never print or exit.
- `main.go`: handle all errors from internal packages; print to stderr; exit with code 1.

```go
// In internal/quotes/quotes.go — return, don't print
if err := json.Unmarshal(data, &qs); err != nil {
    return nil, err
}

// In main.go — print and exit
if err != nil {
    fmt.Fprintln(os.Stderr, "criminalsay: failed to load quotes:", err)
    os.Exit(1)
}
```

### go:embed Anchor Placement (applies to main.go only)

**Rule:** `//go:embed` anchors for files outside a package's subtree must live in `main.go`. This is an inviolable Go toolchain constraint (`go doc embed`: "Patterns may not contain '.' or '..'"). The `data/` directory is outside `internal/quotes/`, so the anchor cannot be in `quotes.go`.

### math/rand/v2 API (applies to quotes.go)

**Rule:** Use `rand.IntN(n)` (capital N). No `rand.Seed()`. No `rand.New()`. Auto-seeded since Go 1.22.

### Table-Driven Tests (applies to quotes_test.go)

**Rule:** All test cases in a `[]struct{}` with a `name` field, iterated with `t.Run(tc.name, ...)`. Use `t.Fatalf` for setup failures, `t.Errorf` for value mismatches.

---

## No Analog Found

All files in this phase have no analog — the codebase is greenfield. The planner must use RESEARCH.md patterns directly (reproduced above in full) for all implementation.

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| `go.mod` | config | — | First file in greenfield project |
| `data/quotes.json` | data | compile-time embed | No data files exist yet |
| `internal/quotes/quotes.go` | service | transform | No Go source files exist yet |
| `internal/quotes/quotes_test.go` | test | batch | No test files exist yet |
| `main.go` | entrypoint | request-response | No Go source files exist yet |

---

## Metadata

**Analog search scope:** entire repository (only CONTEXT.md, LICENSE, README.md, and planning documents found — no Go source files)
**Files scanned:** 0 Go source files (greenfield)
**Pattern extraction date:** 2026-07-05
**Pattern sources:** Go stdlib go doc (go 1.25.6), RESEARCH.md verified patterns, locked decisions from CONTEXT.md
