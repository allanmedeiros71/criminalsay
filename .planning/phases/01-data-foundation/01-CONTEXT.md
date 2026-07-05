# Phase 1: Data Foundation - Context

**Gathered:** 2026-07-05
**Status:** Ready for planning

<domain>
## Phase Boundary

Build the complete data layer: `go.mod`, project layout, `data/quotes.json` with ≥15 real Criminal Minds quotes, `internal/quotes` package (Load + Random), and data-layer tests. No rendering code. The phase ends when `go test ./internal/quotes/...` passes and `go build -o criminalsay .` produces a working binary.

</domain>

<decisions>
## Implementation Decisions

### Character Field Semantics

- **D-01:** When a CM character originates the quote (not citing an outside author), both `author` and `character` fields are set to the same value (e.g., `author: "David Rossi"`, `character: "David Rossi"`).
- **D-02:** When `author ≠ character` (outside author cited by a CM character), the Phase 2 attribution line renders as `"<Original Author>, cited by <Character> · Criminal Minds · S<N>E<N>"`. The `quotes.json` data must supply both fields accurately for this to work.

### Quotes Curation

- **D-03:** Claude pre-populates `quotes.json` with real Criminal Minds quotes from training knowledge. The user reviews and edits afterward for accuracy. Target: ≥15 quotes, spread across all main characters (Reid, Hotch, Morgan, Rossi, Garcia, Prentiss, JJ).
- **D-04:** All quotes stay in English (original series language). No PT-BR translation field in the struct or JSON.

### main.go Phase 1 Stub

- **D-05:** At the end of Phase 1, `main.go` calls `quotes.Load()` + `quotes.Random()` and outputs the raw quote text as `fmt.Println(q.Quote + " — " + q.Author)`. This is an unformatted placeholder replaced entirely by Phase 2's render package.

### Claude's Discretion

- Exact episode and season numbers for individual quotes: Claude selects best-known episodes where attribution can be confident. User reviews for accuracy.
- Number of quotes above the 15-quote minimum: Claude may include up to 25 for a richer initial set.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements and Scope
- `.planning/REQUIREMENTS.md` — Full REQ-IDs for Phase 1: DAT-01, DAT-02, DAT-03, CORE-01, CORE-03, BUILD-01, BUILD-03
- `.planning/ROADMAP.md` — Phase 1 success criteria (5 items) and dependency chain

### Project Decisions and Constraints
- `.planning/PROJECT.md` — Locked decisions D1–D10, key constraints (Go 1.22+/1.26 toolchain, lipgloss v1.1.0, math/rand/v2, go:embed anchor in main.go)
- `CONTEXT.md` — Original project context doc with module contracts (quotes.Load, quotes.Random API signatures) and visual style spec (section 6.1, Estilo 4)

### Go Standard Library Patterns
- No external specs — patterns are stdlib (go:embed, math/rand/v2, encoding/json). See CLAUDE.md for code pattern guidance.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- None yet — greenfield project. Only `CONTEXT.md` and `LICENSE` exist at repo root.

### Established Patterns
- `go:embed` anchor MUST live in `main.go` (not in `internal/quotes/`) to avoid the `..` path restriction that Go enforces for sub-packages embedding files outside their subtree.
- `math/rand/v2` is auto-seeded since Go 1.22; no `rand.Seed()` call needed. Use `rand.IntN(n)` (not deprecated `rand.Intn`).

### Integration Points
- `internal/quotes` is the sole dependency of `main.go` in Phase 1. Phase 2's `internal/render` will be added as a second import.
- `data/quotes.json` is embedded via `//go:embed data/quotes.json` in `main.go`; the raw bytes are passed to `quotes.Load([]byte)`.

</code_context>

<specifics>
## Specific Ideas

- Module API contract (from `CONTEXT.md` section 6, locked):
  - `quotes.Load(data []byte) ([]Quote, error)` — parses embedded JSON
  - `quotes.Random(qs []Quote) Quote` — uniform random selection
- `Quote` struct field set (locked, from REQUIREMENTS.md DAT-02):
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
- `quotes.json` array format: `[{ "quote": "...", "author": "...", "character": "...", "season": N, "episode": N, "episodeTitle": "..." }, ...]`

</specifics>

<deferred>
## Deferred Ideas

- Quote language field (PT-BR translation) — explicitly out of scope per D4 / REQUIREMENTS.md Out of Scope
- Character filter flag (`--character`) — v2, deferred per REQUIREMENTS.md
- Season filter flag (`--season`) — v2, deferred per REQUIREMENTS.md
- Adaptive terminal width detection — Phase 2 concern; Phase 1 uses fixed width spec

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 1-Data Foundation*
*Context gathered: 2026-07-05*
