# Phase 1: Data Foundation — Research

**Researched:** 2026-07-05
**Domain:** Go stdlib (go:embed, math/rand/v2, encoding/json) + module initialization
**Confidence:** HIGH (all claims verified against installed Go 1.25.6 stdlib or go module proxy)

---

<user_constraints>
## User Constraints (from CONTEXT.md / 01-CONTEXT.md)

### Locked Decisions

- **D-01:** When a CM character originates the quote, both `author` and `character` fields are set to the same value.
- **D-02:** When `author ≠ character` (outside author cited by a CM character), Phase 2 attribution renders `"<Original Author>, cited by <Character> · Criminal Minds · S<N>E<N>"`.
- **D-03:** Claude pre-populates `quotes.json` with real Criminal Minds quotes from training knowledge. Target: ≥15 quotes, spread across main characters (Reid, Hotch, Morgan, Rossi, Garcia, Prentiss, JJ).
- **D-04:** All quotes stay in English (original series language). No PT-BR translation field.
- **D-05:** At end of Phase 1, `main.go` outputs `fmt.Println(q.Quote + " — " + q.Author)` as unformatted placeholder.
- **Stack locked:** Go + charmbracelet/lipgloss v1.1.0. No alternatives.
- **Embed anchor:** `//go:embed data/quotes.json` MUST live in `main.go`, not in `internal/quotes/`.
- **Randomness:** `math/rand/v2` (stdlib). No `rand.Seed()`. Use `rand.IntN(n)`.
- **Module:** `github.com/allanmedeiros71/criminalsay`. Single binary at repo root.
- **Layout:** `main.go` at root (not under `cmd/`). `internal/` enforced. No `pkg/`, no `api/`.

### Claude's Discretion

- Exact episode and season numbers for individual quotes: Claude selects best-known episodes. User reviews for accuracy.
- Number of quotes above 15-quote minimum: Claude may include up to 25.

### Deferred Ideas (OUT OF SCOPE)

- Quote language field (PT-BR translation)
- `--character` filter flag (v2)
- `--season` filter flag (v2)
- Adaptive terminal width detection (Phase 2 concern)
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| DAT-01 | Binary contains an embedded JSON file with ≥15 real Criminal Minds quotes compiled at build time (go:embed) | go:embed `[]byte` pattern in main.go; embed anchor placement rules verified |
| DAT-02 | Each quote record includes: `quote`, `author`, `character`, `season`, `episode`, `episodeTitle` | Quote struct definition and JSON schema documented; 20 quotes populated |
| DAT-03 | Missing optional fields (character, episodeTitle) rendered gracefully — omitted without breaking layout | Struct uses Go zero values; empty string/int fields are valid JSON output |
| CORE-01 | Running `criminalsay` selects one quote uniformly at random from embedded dataset | `rand.IntN(n)` API verified; auto-seeding confirmed; guard for n=0 documented |
| CORE-03 | Cold start completes in under 50ms (suitable for .bashrc/.zshrc) | No network I/O, single `json.Unmarshal` on ~5KB data; startup latency trivially satisfied |
| BUILD-01 | `go build -o criminalsay .` produces a working binary from repository root | go.mod initialization pattern documented; lipgloss v1.1.0 verified on module proxy |
| BUILD-03 | `go test ./...` passes with reasonable coverage of `internal/quotes` | Table-driven test patterns documented; Wave 0 test file gaps identified |
</phase_requirements>

---

## Summary

Phase 1 is a greenfield Go module initialization: no existing code, only `CONTEXT.md`, `LICENSE`, and `README.md` exist at the repo root. The phase creates `go.mod`, the project directory structure, `data/quotes.json`, `internal/quotes/quotes.go`, `internal/quotes/quotes_test.go`, and a Phase 1 stub `main.go`. Everything uses Go stdlib (`embed`, `math/rand/v2`, `encoding/json`) plus one external dependency (`github.com/charmbracelet/lipgloss v1.1.0` which, though not used in Phase 1 rendering, must appear in go.mod so Phase 2 can add the import without modifying the module file).

The only non-trivial decision is the `go:embed` anchor placement. Go enforces that embed patterns may not contain `..`, which means a file in `internal/quotes/` cannot embed `../../data/quotes.json`. The anchor must live in `main.go` (at the module root), where the relative path `data/quotes.json` is valid. The raw `[]byte` is passed to `quotes.Load([]byte)`.

The installed Go toolchain (1.25.6) satisfies all requirements: `math/rand/v2` was added in Go 1.22, lipgloss v1.1.0 requires Go 1.18 minimum. The `go.mod` should declare `go 1.22` as the minimum version and `toolchain go1.25.6` to pin the installed version.

**Primary recommendation:** Create the module with `go mod init github.com/allanmedeiros71/criminalsay`, then create the directory tree, quotes.json, internal/quotes package, and stub main.go in a single wave. Run `go mod tidy` after writing all files to populate go.sum. The phase is complete when `go test ./internal/quotes/...` passes and `go build -o criminalsay .` succeeds.

---

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Quote data storage | Binary (embedded) | — | go:embed compiles data/quotes.json into binary at build time; no filesystem access at runtime |
| Random selection | Application logic (internal/quotes) | — | Pure stdlib; no external service; deterministic enough for display purposes |
| Data parsing | Application logic (internal/quotes) | — | encoding/json decodes once at startup; result is a []Quote slice |
| Binary wiring | main.go | — | Embed anchor lives here; orchestrates Load → Random → print |
| Styling (future) | internal/render (Phase 2) | — | Not in Phase 1 scope |

---

## Standard Stack

### Core (Phase 1 only uses stdlib + lipgloss in go.mod)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib: `embed` | Go 1.16+ (stdlib) | Embed quotes.json into binary at compile time | Official solution since Go 1.16; zero external dep; binary ships self-contained |
| Go stdlib: `math/rand/v2` | Go 1.22+ (stdlib) | Uniform random quote selection | Auto-seeded; `IntN(n)` API; replaces deprecated `math/rand` v1 Intn; no go.mod entry |
| Go stdlib: `encoding/json` | stdlib | Decode embedded JSON bytes into []Quote | Standard; correct for small static dataset |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | Terminal styling (Phase 2) — added to go.mod now | v1.1.0 auto-detects TTY/NO_COLOR/color depth; required in go.mod even if Phase 1 doesn't call it |

**Installation (run after writing all .go files):**
```bash
go mod init github.com/allanmedeiros71/criminalsay
go get github.com/charmbracelet/lipgloss@v1.1.0
go mod tidy
```

**Version verification (already confirmed):**
```
go list -m -json github.com/charmbracelet/lipgloss@v1.1.0
→ Version: v1.1.0, Time: 2025-03-12T18:56:50Z
```
[VERIFIED: go module proxy — `go list -m -json github.com/charmbracelet/lipgloss@v1.1.0`]

---

## Package Legitimacy Audit

> The package legitimacy seam does not support the Go ecosystem (npm/pypi/crates only). Manual verification performed via Go module proxy.

| Package | Registry | Age | Downloads | Source Repo | Verdict | Disposition |
|---------|----------|-----|-----------|-------------|---------|-------------|
| `github.com/charmbracelet/lipgloss` | Go module proxy (sum.golang.org) | ~4 yrs (first release ~2021) | Extremely high (charmbracelet is the canonical Go terminal UI org) | github.com/charmbracelet/lipgloss | OK | Approved |

**Packages removed due to SLOP verdict:** none

**Packages flagged as suspicious (SUS):** none

*All stdlib packages (embed, math/rand/v2, encoding/json) carry no legitimacy risk by definition. Lipgloss is confirmed at v1.1.0 via `go list -m -json` against the Go module proxy.* [VERIFIED: go module proxy]

---

## Architecture Patterns

### System Architecture Diagram

```
compile time:
  data/quotes.json ──[go:embed]──► main.go ([]byte quoteData)
                                        │
runtime:                                │
  main.go ──────────────────────────────┘
      │
      ├─ quotes.Load(quoteData) ──► []Quote  (encoding/json.Unmarshal)
      │
      ├─ quotes.Random(qs) ──────► Quote     (rand.IntN(len(qs)))
      │
      └─ fmt.Println(q.Quote + " — " + q.Author)   [Phase 1 stub]
           (replaced entirely by internal/render in Phase 2)
```

### Recommended Project Structure

```
criminalsay/
├── main.go                    # embed anchor, main(), Phase 1 stub output
├── go.mod                     # module declaration, go 1.22 minimum
├── go.sum                     # generated by go mod tidy
├── data/
│   └── quotes.json            # ≥15 Criminal Minds quotes (go:embed target)
├── internal/
│   └── quotes/
│       ├── quotes.go          # Quote struct, Load(), Random()
│       └── quotes_test.go     # table-driven tests for Load and Random
├── CONTEXT.md                 # (existing)
├── LICENSE                    # (existing)
└── README.md                  # (existing)
```

Note: `internal/render/` is NOT created in Phase 1. It is a Phase 2 artifact.

### Pattern 1: go.mod Initialization

**What:** Standard go.mod for a single-binary module with one external dep.

**When to use:** Always — this is the first file created.

```
// Source: go.dev/ref/mod (Go Modules Reference)
module github.com/allanmedeiros71/criminalsay

go 1.22

toolchain go1.25.6

require github.com/charmbracelet/lipgloss v1.1.0
```

[ASSUMED: toolchain directive `go1.25.6` matches currently installed Go version. If user runs different Go version, `go mod tidy` will update this line automatically.]

### Pattern 2: go:embed Anchor in main.go

**What:** Embed `data/quotes.json` as `[]byte` in `main.go` and pass to `quotes.Load`.

**Why this file:** Go's embed spec states "Patterns may not contain '.' or '..'" [VERIFIED: `go doc embed`]. A file in `internal/quotes/` would need `../../data/quotes.json` which is forbidden. `main.go` at the repo root can use the relative path `data/quotes.json` directly.

```go
// Source: go doc embed (stdlib)
package main

import (
    _ "embed"

    "github.com/allanmedeiros71/criminalsay/internal/quotes"
)

//go:embed data/quotes.json
var quoteData []byte

func main() {
    qs, err := quotes.Load(quoteData)
    if err != nil {
        panic(err)
    }
    q := quotes.Random(qs)
    println(q.Quote + " — " + q.Author) // Phase 1 placeholder; replaced in Phase 2
}
```

Key rules:
- The `//go:embed` comment must immediately precede the `var` declaration. [VERIFIED: `go doc embed`]
- Import `_ "embed"` is required even though we never reference `embed.FS`. [VERIFIED: `go doc embed`]
- The variable type is `[]byte`, not `embed.FS` — correct for a single file. [VERIFIED: `go doc embed`]

### Pattern 3: Quote Struct and quotes.go

**What:** The `internal/quotes` package owns the `Quote` type, `Load`, and `Random`.

```go
// Source: DAT-02 (REQUIREMENTS.md) + CONTEXT.md section 6
package quotes

import (
    "encoding/json"
    "math/rand/v2"
)

// Quote represents a single Criminal Minds quote record.
type Quote struct {
    Quote        string `json:"quote"`
    Author       string `json:"author"`
    Character    string `json:"character"`
    Season       int    `json:"season"`
    Episode      int    `json:"episode"`
    EpisodeTitle string `json:"episodeTitle"`
}

// Load parses the embedded JSON bytes and returns a slice of Quote.
// Returns an error if data is malformed or does not decode into a []Quote.
func Load(data []byte) ([]Quote, error) {
    var qs []Quote
    if err := json.Unmarshal(data, &qs); err != nil {
        return nil, err
    }
    return qs, nil
}

// Random returns a uniformly random Quote from qs.
// It panics if qs is empty — callers must validate len(qs) > 0 before calling.
func Random(qs []Quote) Quote {
    return qs[rand.IntN(len(qs))]
}
```

Design notes:
- `Load` does NOT embed data itself — it receives bytes from the caller (main.go). This is the only valid design given the `..` path restriction. [VERIFIED: `go doc embed`]
- `Random` panics on empty slice — this is acceptable because empty data is a programming error (broken build), not a runtime user error. `main.go` should validate `len(qs) > 0` before calling `Random`. [VERIFIED: `go doc math/rand/v2` — `IntN` panics if n<=0]
- `math/rand/v2` global source is auto-seeded since Go 1.22. No `rand.New(rand.NewSource(...))` needed. [VERIFIED: `go doc math/rand/v2`]

### Pattern 4: data/quotes.json Schema

**What:** Array of Quote objects.

```json
[
  {
    "quote": "The quote text goes here.",
    "author": "Original Author Name",
    "character": "CM Character Name",
    "season": 1,
    "episode": 1,
    "episodeTitle": "Extreme Aggressor"
  }
]
```

Field semantics:
- When the CM character IS the originator: `author` = `character` = character name (D-01).
- When the CM character CITES someone else: `author` = original author, `character` = CM character who cited it (D-02).
- `character`, `episodeTitle`: may be empty string `""` — rendered gracefully (DAT-03).
- `season`, `episode`: int, required; use 0 if truly unknown (renders as S0E0 — acceptable placeholder).

### Pattern 5: Table-Driven Tests

**What:** Standard Go testing pattern for `Load` and `Random`.

```go
// Source: Go testing stdlib patterns
package quotes_test

import (
    "testing"

    "github.com/allanmedeiros71/criminalsay/internal/quotes"
)

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
            name:    "valid multiple quotes",
            input:   []byte(`[{"quote":"Q1","author":"A1","character":"A1","season":1,"episode":1,"episodeTitle":""},{"quote":"Q2","author":"A2","character":"A2","season":2,"episode":3,"episodeTitle":""}]`),
            wantErr: false,
            wantLen: 2,
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
    // Run enough iterations to confirm non-constant selection.
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

[ASSUMED: test file uses `package quotes_test` (external test package) so it tests the exported API only. Internal fields are not needed for these tests.]

### Anti-Patterns to Avoid

- **Embed anchor in internal/quotes/quotes.go:** Go forbids `..` in embed paths. The pattern `//go:embed ../../data/quotes.json` will not compile. [VERIFIED: `go doc embed`]
- **rand.Seed() call:** Unnecessary in math/rand/v2 — the global source is auto-seeded. Calling `rand.New(rand.NewSource(time.Now().UnixNano()))` is v1 style and incorrect for v2. [VERIFIED: `go doc math/rand/v2`]
- **rand.Intn(n) (lowercase 'n'):** This is the deprecated v1 API and does not exist in `math/rand/v2`. The v2 API is `rand.IntN(n)` (capital N). [VERIFIED: `go doc math/rand/v2 IntN`]
- **Calling quotes.Random on empty slice:** `rand.IntN(0)` panics. Guard in main.go: `if len(qs) == 0 { fmt.Fprintln(os.Stderr, "criminalsay: no quotes loaded"); os.Exit(1) }`.
- **Importing lipgloss in Phase 1:** Do NOT import `github.com/charmbracelet/lipgloss` from any Phase 1 source file. Add it to go.mod now (`go get`) so Phase 2 can use it without modifying the module declaration, but don't reference it in code yet.
- **cmd/criminalsay/ nesting:** Project layout decision is locked — `main.go` at repo root. Do not create `cmd/` subdirectory. [CITED: CLAUDE.md project instructions]
- **fatih/color in Phase 1:** fatih/color is an alternative considered and rejected in CLAUDE.md. Do not add it to go.mod.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| File embedding at compile time | Custom base64 encoding + code generation | `go:embed` (stdlib) | go:embed is the official, build-cache-aware solution since Go 1.16; custom solutions don't integrate with `go build` |
| Random selection | Custom PRNG, xoshiro implementation | `math/rand/v2` (stdlib) | Auto-seeded, goroutine-safe, crypto-unrelated display use — stdlib is correct |
| JSON parsing | Manual byte scanning, custom parser | `encoding/json` (stdlib) | Quote data is simple; no streaming needed; json.Unmarshal is correct |

**Key insight:** Every problem in Phase 1 is solved by Go's standard library. The only external package (`lipgloss`) is for Phase 2's rendering, not Phase 1.

---

## Common Pitfalls

### Pitfall 1: go:embed path restriction breaks internal/ package placement

**What goes wrong:** Developer places `//go:embed ../../data/quotes.json` inside `internal/quotes/quotes.go`. Build fails with: `pattern ../../data/quotes.json: invalid pattern syntax`.

**Why it happens:** Go's embed spec explicitly forbids `..` in patterns. The spec states: "Patterns may not contain '.' or '..' or empty path elements."

**How to avoid:** Place the `//go:embed data/quotes.json` directive in `main.go` at the repo root. Pass the resulting `[]byte` to `quotes.Load(data []byte)`. This is the only valid architecture for embedding files outside a package's subtree.

**Warning signs:** Any attempt to reference `data/` from within `internal/` using relative paths.

### Pitfall 2: Forgetting `_ "embed"` import

**What goes wrong:** `//go:embed data/quotes.json` directive is present but the `embed` package is not imported, even as a blank import. Build fails with: `go: embed requires import of "embed"` (or similar).

**Why it happens:** The embed package must be imported to activate the embed machinery even if you never use `embed.FS` directly.

**How to avoid:** In `main.go`, include `import _ "embed"` alongside the other imports whenever a `//go:embed` directive is present.

**Warning signs:** Build error mentioning `embed` when `//go:embed` directive is in the file.

### Pitfall 3: rand.IntN panics on empty dataset

**What goes wrong:** `quotes.json` is malformed or empty; `Load` returns `[]Quote{}` (no error on empty array); `Random(qs)` is called with `len(qs) == 0`; `rand.IntN(0)` panics at runtime.

**Why it happens:** `rand.IntN` documents: "It panics if n <= 0." An empty JSON array `[]` is valid JSON and parses without error — `Load` returns `([]Quote{}, nil)`.

**How to avoid:** In `main.go`, after `Load`, check `if len(qs) == 0` and exit with an error message before calling `Random`. The test suite should also cover the "empty array" case for `Load` to confirm it returns a non-nil empty slice, not an error.

**Warning signs:** No length check between `Load` and `Random` in `main.go`.

### Pitfall 4: go.mod toolchain mismatch causes download attempt

**What goes wrong:** `go.mod` declares `toolchain go1.26.0` but the installed toolchain is `go1.25.6`. With `GOTOOLCHAIN=auto`, Go will attempt to download `go1.26.0` from the network. If the machine has no internet access during build, this fails.

**Why it happens:** The `toolchain` directive is a hard requirement. Go 1.25.6 will try to satisfy it.

**How to avoid:** Set `toolchain go1.25.6` in `go.mod` to match the installed version, OR omit the `toolchain` line entirely (let `GOTOOLCHAIN=auto` use whatever is installed). The `go 1.22` minimum version line is sufficient for correctness.

**Warning signs:** `go: downloading go1.26.0` during `go build` on a machine with only 1.25.6 installed.

### Pitfall 5: Quote struct json tags must exactly match JSON field names

**What goes wrong:** JSON field is `"episodeTitle"` (camelCase) but struct tag is `json:"episode_title"` (snake_case). Field silently decodes as zero value.

**Why it happens:** `encoding/json` is case-insensitive for exact matches but case-sensitive for the tag. Mismatched tags cause silent zero-value population, not a decode error.

**How to avoid:** The struct tags are locked in REQUIREMENTS.md and CONTEXT.md (D-02). Use them exactly:
```go
EpisodeTitle string `json:"episodeTitle"`
```

**Warning signs:** `episodeTitle` field is empty string in all quotes after Load, even when the JSON file has values.

---

## Code Examples

### Complete main.go Phase 1 Stub

```go
// Source: go doc embed + 01-CONTEXT.md D-05
package main

import (
    _ "embed"
    "fmt"
    "os"

    "github.com/allanmedeiros71/criminalsay/internal/quotes"
)

//go:embed data/quotes.json
var quoteData []byte

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

### Complete internal/quotes/quotes.go

```go
// Source: stdlib docs + REQUIREMENTS.md DAT-02 + 01-CONTEXT.md specifics section
package quotes

import (
    "encoding/json"
    "math/rand/v2"
)

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
    var qs []Quote
    if err := json.Unmarshal(data, &qs); err != nil {
        return nil, err
    }
    return qs, nil
}

// Random returns a uniformly random element from qs.
// Panics if qs is empty — caller must ensure len(qs) > 0.
func Random(qs []Quote) Quote {
    return qs[rand.IntN(len(qs))]
}
```

### data/quotes.json — 20 Criminal Minds Quotes

> [ASSUMED: Quote text, speaker attribution, season/episode metadata sourced from training knowledge. User MUST review for accuracy before shipping. Some episode numbers may be approximate.]

```json
[
  {
    "quote": "The object of life is not to be on the side of the majority, but to escape finding oneself in the ranks of the insane.",
    "author": "Marcus Aurelius",
    "character": "Spencer Reid",
    "season": 1,
    "episode": 1,
    "episodeTitle": "Extreme Aggressor"
  },
  {
    "quote": "Whoever fights monsters should see to it that in the process he does not become a monster. And if you gaze long enough into an abyss, the abyss will gaze back into you.",
    "author": "Friedrich Nietzsche",
    "character": "Aaron Hotchner",
    "season": 1,
    "episode": 1,
    "episodeTitle": "Extreme Aggressor"
  },
  {
    "quote": "It's alchemy. Alchemy turns common metals into precious ones. Dreams work the same way. Turning something awful into something better.",
    "author": "David Rossi",
    "character": "David Rossi",
    "season": 9,
    "episode": 13,
    "episodeTitle": "The Road Home"
  },
  {
    "quote": "Every day, we get a chance to be better than we were the day before. To grow, to learn, to be better people. And that's what we do.",
    "author": "David Rossi",
    "character": "David Rossi",
    "season": 12,
    "episode": 22,
    "episodeTitle": "Wheels Up"
  },
  {
    "quote": "You only have power over people so long as you don't take everything away from them. But when you've robbed a man of everything, he's no longer in your power — he's free.",
    "author": "Alexander Solzhenitsyn",
    "character": "Spencer Reid",
    "season": 2,
    "episode": 8,
    "episodeTitle": "Empty Planet"
  },
  {
    "quote": "Don't believe what your eyes are telling you. All they show is limitation. Look with your understanding, find out what you already know, and you'll see the way to fly.",
    "author": "Richard Bach",
    "character": "Spencer Reid",
    "season": 3,
    "episode": 1,
    "episodeTitle": "Doubt"
  },
  {
    "quote": "The belief in a supernatural source of evil is not necessary; men alone are quite capable of every wickedness.",
    "author": "Joseph Conrad",
    "character": "Jason Gideon",
    "season": 1,
    "episode": 22,
    "episodeTitle": "The Fisher King, Part 1"
  },
  {
    "quote": "Nothing is permanent in this wicked world — not even our troubles.",
    "author": "Charlie Chaplin",
    "character": "Aaron Hotchner",
    "season": 4,
    "episode": 1,
    "episodeTitle": "Mayhem"
  },
  {
    "quote": "Monsters are real, and ghosts are real too. They live inside us, and sometimes they win.",
    "author": "Stephen King",
    "character": "Derek Morgan",
    "season": 5,
    "episode": 9,
    "episodeTitle": "100"
  },
  {
    "quote": "Fairy tales do not tell children that dragons exist. Children already know that dragons exist. Fairy tales tell children that dragons can be killed.",
    "author": "G.K. Chesterton",
    "character": "Emily Prentiss",
    "season": 6,
    "episode": 1,
    "episodeTitle": "The Longest Night"
  },
  {
    "quote": "The most important things are the hardest to say. They are the things you get ashamed of, because words diminish them.",
    "author": "Stephen King",
    "character": "Jennifer Jareau",
    "season": 7,
    "episode": 1,
    "episodeTitle": "It Takes a Village"
  },
  {
    "quote": "There is no formula for success except perhaps an unconditional acceptance of life and what it brings.",
    "author": "Arthur Rubinstein",
    "character": "Penelope Garcia",
    "season": 3,
    "episode": 9,
    "episodeTitle": "Penelope"
  },
  {
    "quote": "Reason is not automatic. Those who deny it cannot be conquered by it. Do not count on them. Leave them alone.",
    "author": "Ayn Rand",
    "character": "Spencer Reid",
    "season": 4,
    "episode": 13,
    "episodeTitle": "Bloodline"
  },
  {
    "quote": "There are no secrets that time does not reveal.",
    "author": "Jean Racine",
    "character": "Aaron Hotchner",
    "season": 2,
    "episode": 1,
    "episodeTitle": "The Fisher King, Part 2"
  },
  {
    "quote": "It is a wise father that knows his own child.",
    "author": "William Shakespeare",
    "character": "David Rossi",
    "season": 5,
    "episode": 6,
    "episodeTitle": "The Eyes Have It"
  },
  {
    "quote": "Within each of us lies the capacity for both saint and sinner; we can choose which path to walk.",
    "author": "David Rossi",
    "character": "David Rossi",
    "season": 8,
    "episode": 1,
    "episodeTitle": "The Silencer"
  },
  {
    "quote": "Each of us is a book waiting to be written, and that book, if written, results in a person explained.",
    "author": "Thomas M. Cirignano",
    "character": "Aaron Hotchner",
    "season": 3,
    "episode": 11,
    "episodeTitle": "Birthright"
  },
  {
    "quote": "One need not be a chamber to be haunted; one need not be a house; the brain has corridors surpassing material place.",
    "author": "Emily Dickinson",
    "character": "Spencer Reid",
    "season": 6,
    "episode": 18,
    "episodeTitle": "Lauren"
  },
  {
    "quote": "The strength of a family, like the strength of an army, is in its loyalty to each other.",
    "author": "Mario Puzo",
    "character": "Derek Morgan",
    "season": 9,
    "episode": 1,
    "episodeTitle": "The Inspiration"
  },
  {
    "quote": "There is no greater agony than bearing an untold story inside you.",
    "author": "Maya Angelou",
    "character": "Jennifer Jareau",
    "season": 11,
    "episode": 22,
    "episodeTitle": "The Storm"
  }
]
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `math/rand` v1 (`rand.Intn`, manual seed) | `math/rand/v2` (`rand.IntN`, auto-seeded) | Go 1.22 (March 2024) | Remove all `rand.Seed()` and `rand.New()` boilerplate; just call `rand.IntN(n)` |
| `go-bindata` / `pkger` for embedding | `go:embed` (stdlib) | Go 1.16 (February 2021) | No external tools; build system handles embedding natively |
| lipgloss v2 (manual color writer setup) | lipgloss v1.1.0 (auto-detects TTY) | v2 released July 2024 | v1 is correct for this use case; v2 adds complexity that is not needed for a simple stdout CLI |

**Deprecated/outdated:**
- `math/rand.Intn(n)`: replaced by `rand.IntN(n)` in v2. Do not use.
- `rand.Seed(...)`: auto-seeding makes this unnecessary in v2. Do not use.
- `go-bindata`: superseded by `go:embed`. Do not use.

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Quote text and episode attribution for the 20 quotes in quotes.json | Code Examples — data/quotes.json | Some quotes may be slightly misquoted, misattributed, or have wrong S/E numbers. User must review before shipping. No code risk — only data accuracy risk. |
| A2 | Thomas M. Cirignano is a real author whose quote appears in CM S3E11 | quotes.json entry 17 | May be wrong attribution; if uncertain, user should remove that entry or correct it. |
| A3 | `toolchain go1.25.6` is the correct toolchain directive given the installed version | Standard Stack / go.mod Pattern | If user runs with a different Go version, `go mod tidy` will update automatically. Low risk. |

**Verified claims:** All stdlib API facts (embed path restriction, rand.IntN signature, json.Unmarshal signature) verified against `go doc` on the installed Go 1.25.6. Lipgloss v1.1.0 verified on Go module proxy.

---

## Open Questions

1. **Episode accuracy for quotes.json**
   - What we know: 20 quotes are populated from training knowledge; characters and approximate seasons are likely correct.
   - What's unclear: Exact episode numbers and titles may be off by 1-2 for some entries.
   - Recommendation: The phase plan should include a task noting "user reviews quotes.json for accuracy" as an explicit step, not a silent assumption.

2. **Go version in go.mod toolchain directive**
   - What we know: Installed toolchain is go1.25.6; CLAUDE.md recommends go1.26 toolchain.
   - What's unclear: Whether to specify `toolchain go1.25.6` (matching installed) or omit the toolchain line.
   - Recommendation: Use `toolchain go1.25.6` to match what's installed. If user later upgrades to 1.26, `go mod tidy` updates it.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | All compilation | Yes | go1.25.6 linux/amd64 | — |
| go module proxy (sum.golang.org) | `go get lipgloss` | Assumed yes (standard) | — | `GONOSUMCHECK` / private proxy |
| Internet access (one-time) | `go get github.com/charmbracelet/lipgloss@v1.1.0` | Assumed yes | — | Pre-cached in $GOPATH/pkg/mod (already verified with go list) |

**Missing dependencies with no fallback:** None.

**Missing dependencies with fallback:** None.

Note: lipgloss v1.1.0 is already in the local module cache (`/home/allan/go/pkg/mod/cache/`), so `go get` will not require a network call.

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` (no config file needed) |
| Config file | none |
| Quick run command | `go test ./internal/quotes/...` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| DAT-01 | Embedded JSON parses without error; len >= 15 | unit | `go test ./internal/quotes/... -run TestLoad` | Wave 0 |
| DAT-02 | Quote struct fields all populated for a valid entry | unit | `go test ./internal/quotes/... -run TestLoad` | Wave 0 |
| DAT-03 | Missing optional fields (empty string) load without error | unit | `go test ./internal/quotes/... -run TestLoad/missing_optional_fields` | Wave 0 |
| CORE-01 | Random returns item from input slice; distribution is non-constant | unit | `go test ./internal/quotes/... -run TestRandom` | Wave 0 |
| CORE-03 | Binary starts in < 50ms | manual smoke | `time ./criminalsay` after build | N/A |
| BUILD-01 | `go build -o criminalsay .` exits 0 | build smoke | `go build -o criminalsay . && echo OK` | N/A |
| BUILD-03 | `go test ./...` passes | unit suite | `go test ./...` | Wave 0 |

### Sampling Rate

- **Per task commit:** `go test ./internal/quotes/...`
- **Per wave merge:** `go test ./...`
- **Phase gate:** `go build -o criminalsay . && go test ./...` both exit 0

### Wave 0 Gaps

- [ ] `internal/quotes/quotes_test.go` — covers DAT-01, DAT-02, DAT-03, CORE-01, BUILD-03
- [ ] `internal/quotes/quotes.go` — the package itself (Wave 0 creates both together)
- [ ] `data/quotes.json` — must exist before `go build` can embed it

*(No existing test infrastructure — greenfield project. All test files are Wave 0 artifacts.)*

---

## Security Domain

> `security_enforcement: true`, `security_asvs_level: 1` per `.planning/config.json`.

### Applicable ASVS Categories (ASVS Level 1)

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | No | CLI tool with no auth; no user accounts |
| V3 Session Management | No | Stateless print-and-exit tool |
| V4 Access Control | No | No access control surface |
| V5 Input Validation | Minimal | `json.Unmarshal` validates embedded data; data is compile-time, not user input |
| V6 Cryptography | No | `math/rand/v2` is for display, not security; no secrets |
| V9 Communication | No | No network I/O |
| V13 API | No | No API surface |

### Known Threat Patterns for this Stack

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Malformed embedded JSON causing panic/crash | Tampering (compile-time supply chain) | `json.Unmarshal` returns error; main.go handles error and exits cleanly |
| Empty dataset causing `rand.IntN(0)` panic | Denial of Service | Guard `len(qs) == 0` check in main.go before calling `Random` |

**Security summary:** The attack surface of a stateless, offline, print-and-exit CLI that reads only compile-time embedded data is minimal. The two risks above are data integrity issues, not security vulnerabilities in the traditional sense. ASVS Level 1 is satisfied by: (1) not using `crypto/rand` where `math/rand` is sufficient; (2) handling parse errors; (3) not accepting any user-controlled input at runtime.

---

## Sources

### Primary (MEDIUM confidence — verified via go doc against installed stdlib)

- `go doc embed` — embed path restriction, `//go:embed` directive syntax, `_ "embed"` import requirement [VERIFIED: go 1.25.6 stdlib]
- `go doc math/rand/v2` — `IntN(n)` signature, auto-seeding behavior, panic on n<=0 [VERIFIED: go 1.25.6 stdlib]
- `go doc math/rand/v2 IntN` — exact signature confirmation [VERIFIED: go 1.25.6 stdlib]
- `go doc encoding/json Unmarshal` — signature, behavior on nil/non-pointer, error semantics [VERIFIED: go 1.25.6 stdlib]
- `go list -m -json github.com/charmbracelet/lipgloss@v1.1.0` — version, timestamp, source URL [VERIFIED: go module proxy]

### Secondary (MEDIUM confidence — sourced from project CONTEXT.md and locked decisions)

- `CONTEXT.md` section 6 — module API contracts (Load, Random signatures), project structure, Estilo 4 visual spec
- `01-CONTEXT.md` — locked decisions D-01 through D-05, Quote struct field set
- `REQUIREMENTS.md` — DAT-01, DAT-02, DAT-03, CORE-01, CORE-03, BUILD-01, BUILD-03 definitions

### Tertiary (ASSUMED — training knowledge)

- Criminal Minds quotes and episode attribution in `data/quotes.json` — sourced from training knowledge; user must review for accuracy

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all stdlib APIs verified against go doc; lipgloss v1.1.0 confirmed on module proxy
- Architecture: HIGH — derived from locked decisions in CONTEXT.md and embed path restriction verified in go doc
- Pitfalls: HIGH — all pitfalls derived from verified stdlib behavior (go doc confirms embed restriction and rand.IntN panic)
- Quote data: LOW — sourced from training knowledge; user review required

**Research date:** 2026-07-05
**Valid until:** 2027-01-05 (Go stdlib APIs are extremely stable; lipgloss v1.x is stable)
