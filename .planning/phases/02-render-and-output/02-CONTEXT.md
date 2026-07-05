# Phase 2: Render and Output - Context

**Gathered:** 2026-07-05
**Status:** Ready for planning

<domain>
## Phase Boundary

Replace the Phase 1 stub output in `main.go` with styled terminal rendering via a new `internal/render` package. Running `criminalsay` must print one randomly selected quote using Estilo 4 "Terminal cru" (red sidebar, quoted text, indented attribution), degrade gracefully when color is unavailable, and pass `go test ./...` covering both `internal/quotes` and `internal/render`.

No CLI flags, no alternate themes, no cross-compile Makefile (Phase 3), no README changes beyond what Phase 3 owns.

</domain>

<decisions>
## Implementation Decisions

### Attribution Line

- **D-01:** When `author ≠ character` (external quote cited by a CM character), format: `<Author>, cited by <Character> · Criminal Minds · S<season>E<episode>`. Example: `Marcus Aurelius, cited by Spencer Reid · Criminal Minds · S1E1`.
- **D-02:** When `author == character` (original line from the character), show only the name once: `<Author> · Criminal Minds · S<season>E<episode>`. Example: `David Rossi · Criminal Minds · S9E13`.
- **D-03:** Include `episodeTitle` in parentheses after the episode code when non-empty: `· S9E13 (The Road Home)`. Omit parentheses entirely when `episodeTitle` is empty.
- **D-04:** Color scope on attribution: only the author portion (everything before the first ` · `) is yellow/highlighted; `· Criminal Minds · SxEy (Title)` is dim/muted. When D-01 applies, yellow covers the full `"<Author>, cited by <Character>"` segment before the first ` · Criminal Minds`.

### Word Wrap & Layout

- **D-05:** Fixed wrap width of **52 useful columns** for quote text (not adaptive to terminal width). Aligns with REND-02; predictable for `.bashrc` use.
- **D-06:** Two-space left indent before the `▌` sidebar on every quote line.
- **D-07:** Red `▌` (U+258C) prefix on **every** wrapped line of the quote, not just the first.
- **D-08:** Double quotation marks wrap the full quote text; word-wrap breaks occur inside the quoted block (opening `"` on first line, closing `"` on last line).
- **D-09:** Blank line between the quote block and the attribution line (REND-04).
- **D-10:** Long attribution lines wrap at the same 52-column block width; no truncation.

### Optional Field Handling

- **D-11:** Empty `character`: omit silently; attribution uses `author` only (no "cited by" clause).
- **D-12:** `season` or `episode` is 0: omit the `S<season>E<episode>` segment entirely (and omit episode title parens if episode code is omitted).
- **D-13:** Empty `episodeTitle`: omit `(Title)` parens; keep `SxEy` when season/episode are valid.

### Color & Compatibility

- **D-14:** `main.go` detects `stdout.IsTerminal()` and passes `colorEnabled bool` to `render.Quote(q, colorEnabled)`. lipgloss handles color depth (truecolor/256/ANSI-16) when enabled.
- **D-15:** Respect `NO_COLOR` environment variable in `main` — when set or stdout is not a TTY, pass `colorEnabled=false`.
- **D-16:** Quote text between double quotes uses white/bright foreground when color is enabled (REND-03).
- **D-17:** Sidebar `▌` uses red foreground when color is enabled (REND-01).
- **D-18:** **User override of COLOR-02:** In no-color mode (`NO_COLOR`, non-TTY, or piped), replace `▌` with ASCII `|` for readability; strip all ANSI sequences; preserve quotation marks and `·` separators. This intentionally deviates from the literal COLOR-02 wording (which preserves `▌`) — user preference for plainer terminals.

### Module Contract (unchanged from project spec)

- **D-19:** `render.Quote(q quotes.Quote, colorEnabled bool) string` returns the formatted string without printing. `main` owns `fmt.Println`.
- **D-20:** No box border, no header, no footer (REND-06). Estilo 4 sidebar-only presentation.

### Claude's Discretion

- Exact lipgloss style IDs and helper decomposition inside `internal/render` (single file vs split wrap/style helpers).
- Word-wrap algorithm choice (`strings` + manual vs lipgloss `Width`/`Join` patterns) as long as D-05–D-08 are satisfied.
- Test strategy for ANSI stripping and wrap alignment (table-driven with `colorEnabled` true/false cases).

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements and Scope
- `.planning/REQUIREMENTS.md` — Phase 2 REQ-IDs: CORE-02, REND-01 through REND-06, COLOR-01, COLOR-02, BUILD-03
- `.planning/ROADMAP.md` — Phase 2 success criteria (5 items) and UI hint

### Visual Spec (locked)
- `CONTEXT.md` §6.1 — Estilo 4 "Terminal cru" reference output, sidebar rules, indentation, color roles
- `.planning/phases/01-data-foundation/01-CONTEXT.md` — D-01/D-02 character semantics (extends into attribution format here)

### Project Constraints
- `.planning/PROJECT.md` — lipgloss v1.1.0 (not v2), no interactivity, Go 1.22+
- `.planning/phases/01-data-foundation/01-VERIFICATION.md` — Phase 1 baseline; `main.go` embed anchor and `internal/quotes` API verified

### Library
- `go.mod` — `charmbracelet/lipgloss v1.1.0` already declared (indirect); Phase 2 adds first import in `internal/render`

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/quotes` — `Quote` struct, `Load`, `Random` (tested, 20-quote dataset)
- `main.go` — embed anchor `//go:embed data/quotes.json`, error handling pattern (stderr + exit 1), stub output to replace
- `go.mod` / `go.sum` — lipgloss v1.1.0 pre-declared

### Established Patterns
- Embed anchor MUST stay in `main.go` (not `internal/`)
- Render produces string; main prints (testability — no stdout capture in render tests)
- Table-driven external tests (`quotes_test` package pattern)
- Two-tier errors: internal returns, main prints to stderr and exits

### Integration Points
- `main.go`: `quotes.Load` → guard `len==0` → `quotes.Random` → detect color → `render.Quote` → `fmt.Println`
- `internal/render` imports `internal/quotes` for `Quote` type only
- Phase 3 will add Makefile; Phase 2 must not touch build scripts

</code_context>

<specifics>
## Specific Ideas

Reference output (from `CONTEXT.md` §6.1, with D-03 episode title extension):

```
  ▌ "It's alchemy. Alchemy turns common metals into
  ▌ precious ones. Dreams work the same way. Turning
  ▌ something awful into something better."

       David Rossi · Criminal Minds · S9E13 (The Road Home)
```

External-author example (D-01 + D-03):

```
       Marcus Aurelius, cited by Spencer Reid · Criminal Minds · S1E1 (Extreme Aggressor)
```

No-color example (D-18):

```
  | "Quote text wrapped here..."
```

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 2-Render and Output*
*Context gathered: 2026-07-05*
