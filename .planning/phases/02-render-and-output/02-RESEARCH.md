# Phase 2: Render and Output — Research

**Researched:** 2026-07-05
**Domain:** Go terminal rendering (`charmbracelet/lipgloss` v1.1.0), word-wrap layout, TTY/NO_COLOR gating, table-driven tests
**Confidence:** HIGH (lipgloss/termenv APIs verified against installed module source and `go doc`; layout rules locked in 02-CONTEXT.md)

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

#### Attribution Line

- **D-01:** When `author ≠ character` (external quote cited by a CM character), format: `<Author>, cited by <Character> · Criminal Minds · S<season>E<episode>`. Example: `Marcus Aurelius, cited by Spencer Reid · Criminal Minds · S1E1`.
- **D-02:** When `author == character` (original line from the character), show only the name once: `<Author> · Criminal Minds · S<season>E<episode>`. Example: `David Rossi · Criminal Minds · S9E13`.
- **D-03:** Include `episodeTitle` in parentheses after the episode code when non-empty: `· S9E13 (The Road Home)`. Omit parentheses entirely when `episodeTitle` is empty.
- **D-04:** Color scope on attribution: only the author portion (everything before the first ` · `) is yellow/highlighted; `· Criminal Minds · SxEy (Title)` is dim/muted. When D-01 applies, yellow covers the full `"<Author>, cited by <Character>"` segment before the first ` · Criminal Minds`.

#### Word Wrap & Layout

- **D-05:** Fixed wrap width of **52 useful columns** for quote text (not adaptive to terminal width). Aligns with REND-02; predictable for `.bashrc` use.
- **D-06:** Two-space left indent before the `▌` sidebar on every quote line.
- **D-07:** Red `▌` (U+258C) prefix on **every** wrapped line of the quote, not just the first.
- **D-08:** Double quotation marks wrap the full quote text; word-wrap breaks occur inside the quoted block (opening `"` on first line, closing `"` on last line).
- **D-09:** Blank line between the quote block and the attribution line (REND-04).
- **D-10:** Long attribution lines wrap at the same 52-column block width; no truncation.

#### Optional Field Handling

- **D-11:** Empty `character`: omit silently; attribution uses `author` only (no "cited by" clause).
- **D-12:** `season` or `episode` is 0: omit the `S<season>E<episode>` segment entirely (and omit episode title parens if episode code is omitted).
- **D-13:** Empty `episodeTitle`: omit `(Title)` parens; keep `SxEy` when season/episode are valid.

#### Color & Compatibility

- **D-14:** `main.go` detects `stdout.IsTerminal()` and passes `colorEnabled bool` to `render.Quote(q, colorEnabled)`. lipgloss handles color depth (truecolor/256/ANSI-16) when enabled.
- **D-15:** Respect `NO_COLOR` environment variable in `main` — when set or stdout is not a TTY, pass `colorEnabled=false`.
- **D-16:** Quote text between double quotes uses white/bright foreground when color is enabled (REND-03).
- **D-17:** Sidebar `▌` uses red foreground when color is enabled (REND-01).
- **D-18:** **User override of COLOR-02:** In no-color mode (`NO_COLOR`, non-TTY, or piped), replace `▌` with ASCII `|` for readability; strip all ANSI sequences; preserve quotation marks and `·` separators. This intentionally deviates from the literal COLOR-02 wording (which preserves `▌`) — user preference for plainer terminals.

#### Module Contract (unchanged from project spec)

- **D-19:** `render.Quote(q quotes.Quote, colorEnabled bool) string` returns the formatted string without printing. `main` owns `fmt.Println`.
- **D-20:** No box border, no header, no footer (REND-06). Estilo 4 sidebar-only presentation.

### Claude's Discretion

- Exact lipgloss style IDs and helper decomposition inside `internal/render` (single file vs split wrap/style helpers).
- Word-wrap algorithm choice (`strings` + manual vs lipgloss `Width`/`Join` patterns) as long as D-05–D-08 are satisfied.
- Test strategy for ANSI stripping and wrap alignment (table-driven with `colorEnabled` true/false cases).

### Deferred Ideas (OUT OF SCOPE)

None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| CORE-02 | Program prints the formatted quote to stdout and exits with code 0 | `main.go` wires `quotes.Load` → `quotes.Random` → `colorEnabled()` → `render.Quote` → `fmt.Println`; render returns string only (D-19) |
| REND-01 | Red `▌` sidebar prefix on every quote line, two-column left indent | Per-line manual prefix `  ▌ ` with `sidebarStyle.Foreground(red)`; D-18 swaps to `\|` when `colorEnabled=false` |
| REND-02 | Quote wrapped at ≤52 useful columns; sidebar repeated per line | Word-wrap at 52 **total display columns** per line (verified against CONTEXT.md §6.1 example); per-line prefix re-applied after wrap |
| REND-03 | Quote text between double quotes in white/bright | `quoteStyle.Foreground(bright)` applied to quoted segment only; opening `"` on first line, closing on last (D-08) |
| REND-04 | Blank line between quote block and attribution | `strings.Join` quote lines + `"\n\n"` + attribution block |
| REND-05 | Attribution: yellow author + dim rest, `·` separators | Build `authorPart` + `restPart` per D-01–D-04; `yellowStyle` + `dimStyle.Faint(true)`; split color at first ` · Criminal Minds` per line after wrap |
| REND-06 | No box border, header, or footer | Do **not** use `Style.Border()` box mode; sidebar is a plain per-line prefix (D-20) |
| COLOR-01 | TTY with color support uses ANSI sequences | When `colorEnabled=true`, lipgloss `DefaultRenderer` + `termenv` auto-select TrueColor/256/ANSI profile from stdout [VERIFIED: lipgloss v1.1.0 renderer.go] |
| COLOR-02 | NO_COLOR / non-TTY / pipe: no ANSI; readable plain text | `main` gates with `colorEnabled=false`; render uses plain strings + `\|` sidebar (D-18); tests assert no `\x1b[` |
| BUILD-03 | `go test ./...` passes covering `internal/quotes` and `internal/render` | Table-driven `render_test` package; no stdout capture — assert on `render.Quote` return value |
</phase_requirements>

---

## Summary

Phase 2 adds `internal/render` and replaces the Phase 1 stub in `main.go`. The render package owns all Estilo 4 layout: per-line sidebar prefix, 52-column word wrap, quoted text styling, attribution formatting (D-01–D-13), and color degradation (D-14–D-18). `main.go` remains the orchestrator: embed/load/randomize quotes, detect color capability, print the rendered string.

**Critical architectural choice:** Do **not** use `lipgloss.Style.Border()` as the primary sidebar mechanism. While `applyBorder` does paint a left border character on every output line [VERIFIED: lipgloss v1.1.0 `borders.go` lines 390–411], Estilo 4 needs (a) a 2-space indent **before** the sidebar (D-06), (b) a space **after** the sidebar before the quote (CONTEXT.md §6.1), (c) asymmetric quote-mark placement across first/last lines (D-08), and (d) no box edges (D-20). **Manual per-line prefix composition** after word-wrap is the standard pattern that satisfies all four without fighting border width math.

For word-wrap, use `github.com/charmbracelet/x/cellbuf.Wrap` (already a transitive dependency via lipgloss) or an equivalent word-boundary loop. lipgloss `Style.Width(n)` also wraps via `cellbuf.Wrap` internally [VERIFIED: lipgloss v1.1.0 `style.go` line 368], but applying `Width` to the whole quote block cannot inject per-line sidebar prefixes — wrap first, decorate after.

Color gating is **explicit** via the `colorEnabled` parameter (D-14/D-15/D-19). `main.go` computes it from `NO_COLOR` + `term.IsTerminal(os.Stdout.Fd())`. The render package must not re-detect TTY independently; when `colorEnabled=false`, skip lipgloss color calls entirely and emit plain text with `|` sidebar (D-18). When `colorEnabled=true`, lipgloss/termenv handle profile selection on stdout automatically.

**Primary recommendation:** Implement `internal/render` as `wrap.go` (52-col word wrap + quote-line assembly) + `render.go` (`Quote`, attribution builder, lipgloss style helpers). Wire `main.go` with `colorEnabled := os.Getenv("NO_COLOR") == "" && term.IsTerminal(os.Stdout.Fd())`. Test via `render_test` external package asserting string content, line prefixes, wrap width, attribution segments, and absence of ANSI when `colorEnabled=false`.

---

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Quote selection / embed | Application (`main.go` + `internal/quotes`) | — | Already implemented in Phase 1; render consumes `quotes.Quote` |
| TTY / NO_COLOR detection | Application entry (`main.go`) | — | D-14/D-15 lock detection in main; passes bool to render |
| Word-wrap layout | Application logic (`internal/render`) | — | Pure string transformation; no I/O |
| ANSI styling (color mode) | Application logic (`internal/render`) | lipgloss/termenv (library) | lipgloss applies escape sequences; render owns when to call it |
| Plain-text degradation | Application logic (`internal/render`) | — | D-18 override: `\|` sidebar, no lipgloss calls |
| stdout printing | Application entry (`main.go`) | — | D-19: render returns string; main prints |
| Terminal color depth | Library (lipgloss → termenv) | OS stdout | Auto TrueColor/256/ANSI when `colorEnabled=true` |

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/charmbracelet/lipgloss` | v1.1.0 | ANSI styling: `Foreground`, `Faint`, `Render` | Locked in PROJECT.md; already in `go.mod`; auto color profile via termenv |
| `github.com/charmbracelet/x/term` | v0.2.1 (transitive) | `term.IsTerminal(fd)` for stdout TTY check | Same Charm ecosystem; already pulled by lipgloss |
| `github.com/charmbracelet/x/cellbuf` | v0.0.13 (transitive) | `cellbuf.Wrap(s, limit, "")` for word-wrap | Same algorithm lipgloss uses internally; Unicode-aware width |
| `github.com/charmbracelet/x/ansi` | v0.8.0 (transitive) | `ansi.StringWidth(s)` for display-width assertions in tests | Correct width measurement for Unicode + ANSI |
| Go stdlib: `strings`, `fmt`, `os` | Go 1.22+ | Line assembly, env check, printing | No extra deps for layout logic |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `github.com/muesli/termenv` | v0.16.0 (transitive) | `EnvNoColor`, `Profile.Ascii` | Reference only — main owns NO_COLOR gate; useful if planner adds renderer-isolated tests with `lipgloss.SetColorProfile(termenv.Ascii)` |
| Go stdlib: `testing` | stdlib | Table-driven render tests | BUILD-03 |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Manual per-line prefix | `lipgloss.Style.Border({Left:"▌"}).BorderLeft(true)` | Border adds per-line left char but cannot express 2-space pre-indent + post-sidebar gap + D-08 quote placement without padding/border width fighting 52-col budget |
| `cellbuf.Wrap` | `lipgloss.Style.Width(n).Render(text)` | `Width` wraps whole block; cannot inject sidebar between wrap passes |
| `colorEnabled` param | lipgloss `termenv.EnvNoColor()` inside render | Violates D-14/D-15 — main owns detection; render must accept bool |
| `fatih/color` | lipgloss | Rejected in project spec; lipgloss already in go.mod |

**Installation:** No new dependencies. Phase 2 adds the first direct import of lipgloss in `internal/render`. Run `go mod tidy` after creating the package (should be a no-op).

**Version verification:**
```
go list -m -json github.com/charmbracelet/lipgloss@v1.1.0
→ Version: v1.1.0, Time: 2025-03-12T18:56:50Z
```
[VERIFIED: go module proxy — `go list -m -json`]

---

## Package Legitimacy Audit

> Go ecosystem — manual verification via Go module proxy (legitimacy seam supports npm/pypi/crates only).

| Package | Registry | Age | Downloads | Source Repo | Verdict | Disposition |
|---------|----------|-----|-----------|-------------|---------|-------------|
| `github.com/charmbracelet/lipgloss` | sum.golang.org | ~4 yrs | Very high (Charm canonical) | github.com/charmbracelet/lipgloss | OK | Approved |
| `github.com/charmbracelet/x/term` | sum.golang.org | Charm org | Transitive via lipgloss | github.com/charmbracelet/x | OK | Approved (transitive) |
| `github.com/charmbracelet/x/cellbuf` | sum.golang.org | Charm org | Transitive via lipgloss | github.com/charmbracelet/x | OK | Approved (transitive) |

**Packages removed due to SLOP verdict:** none

**Packages flagged as suspicious (SUS):** none

---

## Architecture Patterns

### System Architecture Diagram

```
main.go
  │
  ├─ quotes.Load(embedded JSON) ──► []Quote
  ├─ quotes.Random(qs) ───────────► Quote
  ├─ colorEnabled() ────────────────► bool
  │     ├─ os.Getenv("NO_COLOR") != ""  → false
  │     └─ !term.IsTerminal(stdout.Fd()) → false
  │
  ├─ render.Quote(q, colorEnabled) ─► string (no I/O)
  │     ├─ wrapQuoteText(q.Quote, 52) ──► []string lines
  │     ├─ prefix each line: "  ▌ " or "  | " (D-18)
  │     ├─ apply quote styles (D-16) or plain text
  │     ├─ blank line (D-09)
  │     └─ buildAttribution(q) + wrap + color split (D-01–D-05, D-10)
  │
  └─ fmt.Println(rendered) ────────► stdout, exit 0
```

### Recommended Project Structure

```
internal/
├── quotes/           # Phase 1 (unchanged)
│   ├── quotes.go
│   └── quotes_test.go
└── render/           # Phase 2 (new)
    ├── render.go     # Quote(), attribution builder, style helpers
    ├── wrap.go       # 52-col word wrap, quote-line assembly
    └── render_test.go
main.go               # replace stub; add colorEnabled(), render import
```

### Pattern 1: Color Detection in main.go

**What:** Gate color in `main`, pass bool to render.

**When to use:** Always — locked by D-14/D-15.

```go
// Source: D-14/D-15 + charmbracelet/x/term (go doc)
package main

import (
    "os"

    term "github.com/charmbracelet/x/term"
)

func colorEnabled() bool {
    if os.Getenv("NO_COLOR") != "" {
        return false
    }
    return term.IsTerminal(os.Stdout.Fd())
}
```

`termenv.EnvNoColor()` also checks `NO_COLOR` and `CLICOLOR=0` [VERIFIED: termenv v0.16.0 `termenv.go`]. Project scope is **`NO_COLOR` + TTY only** per D-15 — do not add `CLICOLOR` handling unless user expands scope.

### Pattern 2: lipgloss Style Helpers (color mode only)

**What:** Package-level styles created once; applied only when `colorEnabled=true`.

```go
// Source: lipgloss v1.1.0 go doc + pkg.go.dev v1.1.0 color section
import "github.com/charmbracelet/lipgloss"

var (
    sidebarStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
    quoteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")) // bright white
    authorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow
    dimStyle     = lipgloss.NewStyle().Faint(true)                    // muted meta
)

func styleSidebar(ch rune, enabled bool) string {
    s := string(ch)
    if !enabled {
        return s
    }
    return sidebarStyle.Render(s)
}
```

lipgloss `Color("9")` / `Color("15")` / `Color("11")` work across ANSI/256/TrueColor profiles — termenv downgrades automatically [CITED: pkg.go.dev/github.com/charmbracelet/lipgloss@v1.1.0 — Colors section].

### Pattern 3: 52-Column Quote Line Assembly

**What:** Wrap quote text, then prefix each line. **52 = total display width per line** (matches CONTEXT.md §6.1 example: `  ▌ "` + 47 chars = 52).

```go
// Source: D-05–D-08 + CONTEXT.md §6.1 + cellbuf.Wrap
import (
    "strings"

    "github.com/charmbracelet/x/ansi"
    "github.com/charmbracelet/x/cellbuf"
)

const (
    blockWidth    = 52
    quoteIndent   = "  "       // D-06
    sidebarColor  = '▌'        // D-17; D-18 uses '|'
    sidebarPlain  = '|'
    sidebarGap    = " "
)

func wrapQuoteLines(text string) []string {
    prefix := quoteIndent + string(sidebarColor) + sidebarGap // 4 display cols
    // Middle lines get `prefix` + text = 52 → text budget 48
    // First line adds opening quote: budget 47
    // Last line adds closing quote: budget 47
    // Single-line quote: budget 46

    // Recommended: word-split with strings.Fields, greedy-pack per line budgets
    // OR: cellbuf.Wrap on inner text at 48, then adjust first/last for quote marks
    _ = cellbuf.Wrap
    _ = ansi.StringWidth
    _ = prefix
    return nil // planner implements greedy word-wrap
}

func assembleQuoteLine(prefix, segment string, open, close bool) string {
    var b strings.Builder
    b.WriteString(prefix)
    if open {
        b.WriteByte('"')
    }
    b.WriteString(segment)
    if close {
        b.WriteByte('"')
    }
    return b.String()
}
```

**Width verification in tests:** use `ansi.StringWidth(line) <= 52` for every quote line [VERIFIED: charmbracelet/x/ansi].

### Pattern 4: Attribution Builder (D-01–D-04, D-11–D-13)

**What:** Build plain-text attribution, wrap at 52 cols, apply split coloring.

```go
// Source: 02-CONTEXT.md D-01 through D-13
func buildAttribution(q quotes.Quote) (authorPart, metaPart string) {
    // Author segment (yellow when colored)
    switch {
    case q.Character != "" && q.Author != q.Character:
        authorPart = q.Author + ", cited by " + q.Character // D-01
    default:
        authorPart = q.Author // D-02; D-11: empty character → author only
    }

    // Meta segment (dim when colored) — always starts with " · Criminal Minds"
    meta := " · Criminal Minds"
    if q.Season != 0 && q.Episode != 0 { // D-12: omit if either is 0
        meta += fmt.Sprintf(" · S%dE%d", q.Season, q.Episode)
        if q.EpisodeTitle != "" { // D-03/D-13
            meta += " (" + q.EpisodeTitle + ")"
        }
    }
    return authorPart, meta
}

func renderAttributionLine(line string, colorEnabled bool) string {
    const marker = " · Criminal Minds"
    idx := strings.Index(line, marker)
    if !colorEnabled {
        return line
    }
    if idx < 0 {
        return authorStyle.Render(line) // author-only edge case
    }
    return authorStyle.Render(line[:idx]) + dimStyle.Render(line[idx:])
}
```

**Attribution indent:** CONTEXT.md §6.1 reference output uses **7 spaces** before the author name (`       David Rossi · ...`). D-06 only specifies 2-space indent for quote lines. Apply 7-space attribution indent to match the locked visual spec [CITED: CONTEXT.md §6.1].

### Pattern 5: render.Quote Contract

```go
// Source: D-19
func Quote(q quotes.Quote, colorEnabled bool) string {
    sidebar := sidebarColor
    if !colorEnabled {
        sidebar = sidebarPlain // D-18
    }

    quoteLines := buildQuoteBlock(q.Quote, sidebar, colorEnabled)
    author, meta := buildAttribution(q)
    attrLines := wrapAttribution(author+meta, blockWidth, attrIndent)

    var out strings.Builder
    out.WriteString(strings.Join(quoteLines, "\n"))
    out.WriteString("\n\n") // D-09
    // render attribution lines with indent + color split
    return strings.TrimRight(out.String(), "\n") + "\n"
}
```

### Anti-Patterns to Avoid

- **Using `lipgloss.Style.Border()` for Estilo 4 sidebar:** Border is designed for box edges; fights 52-col budget and D-06 pre-indent. Use per-line prefix instead.
- **Calling `lipgloss.SetColorProfile` globally in render:** Mutates package-global renderer; complicates parallel tests. Prefer `colorEnabled` branch that skips styling.
- **Capturing stdout in render tests:** Violates D-19 testability contract. Test `render.Quote` return value only.
- **Re-detecting TTY inside render:** Duplicates D-14/D-15; creates divergence from main.
- **Keeping `▌` in no-color mode:** Violates D-18 user override (even though COLOR-02 wording mentions preserving `▌`).
- **Using `lipgloss.Style.Width(52)` on the full quote block:** Cannot insert per-line sidebar after wrap; produces wrong layout.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| ANSI color sequences | Manual `\x1b[...]` strings | lipgloss `Style.Render` | Profile-aware (16/256/truecolor); termenv handles downgrade |
| Unicode display width | `len(s)` or byte counting | `ansi.StringWidth` | `▌`, smart quotes, combining marks have variable width |
| Word wrap with ANSI | Custom byte splitter | `cellbuf.Wrap` | Handles ANSI state; same as lipgloss internals |
| TTY detection | `stat /proc/self/fd` hacks | `term.IsTerminal(os.Stdout.Fd())` | Portable across Linux/macOS/Windows |
| NO_COLOR parsing | Custom env logic beyond spec | `os.Getenv("NO_COLOR") != ""` | D-15: any non-empty value disables color per no-color.org convention [CITED: termenv v0.16.0] |

**Key insight:** Hand-roll layout (prefix assembly, attribution string building) but never hand-roll ANSI encoding or terminal width measurement.

---

## Common Pitfalls

### Pitfall 1: Treating 52 columns as inner-text-only width

**What goes wrong:** Lines exceed 52 display columns when sidebar + indent + quotes are added.

**Why it happens:** D-05 says "52 useful columns for quote text" but the CONTEXT.md §6.1 example lines measure **52 total** including `  ▌ "`.

**How to avoid:** Target `ansi.StringWidth(line) <= 52` for every quote line in tests. Budget per-line: prefix (4) + optional quote marks.

**Warning signs:** Visual misalignment in terminal; test assertions on character count instead of display width.

### Pitfall 2: Sidebar only on first line

**What goes wrong:** Wrapped lines lack `▌`/`|` prefix; layout breaks REND-01/REND-02.

**Why it happens:** Applying border or style to a single concatenated string without per-line prefix loop.

**How to avoid:** Wrap first → iterate lines → prefix each (D-07).

**Warning signs:** Only first line starts with sidebar in test output.

### Pitfall 3: ANSI escapes in no-color output

**What goes wrong:** `NO_COLOR=1 criminalsay` still emits `\x1b[` sequences; fails COLOR-02/D-18.

**Why it happens:** lipgloss styles applied even when `colorEnabled=false`, or global renderer still has color profile.

**How to avoid:** When `!colorEnabled`, use plain `fmt`-less string assembly with zero `Style.Render` calls. Test: `strings.Contains(out, "\x1b[")` must be false.

**Warning signs:** `render_test` passes for content but `NO_COLOR=1 criminalsay | cat -v` shows escape codes.

### Pitfall 4: Attribution color bleed across wrap boundary

**What goes wrong:** Second line of wrapped attribution is all yellow or all dim.

**Why it happens:** Coloring the full string before wrap instead of per-line split at ` · Criminal Minds` marker.

**How to avoid:** Wrap plain text first, then `renderAttributionLine` per line (Pattern 4).

**Warning signs:** Yellow extends into `· Criminal Minds · S9E13` on wrapped lines.

### Pitfall 5: D-12 mis-implemented as "both must be non-zero"

**What goes wrong:** `S1E0` or `S0E1` rendered when one field is zero.

**Why it happens:** Checking only one field or using `season != 0 || episode != 0`.

**How to avoid:** Emit `S%dE%d` only when `season != 0 && episode != 0`. If omitted, also omit `(episodeTitle)` parens.

**Warning signs:** Quotes with `episode: 0` show `S1E0`.

### Pitfall 6: Importing render from main without direct lipgloss in main

**What goes wrong:** No issue functionally, but adding TTY logic to render breaks tier separation.

**How to avoid:** Keep **all** `term.IsTerminal` / `NO_COLOR` checks in `main.go` only.

**Warning signs:** `internal/render` imports `os` for env/stdout beyond test helpers.

---

## Code Examples

### main.go Wiring (complete Phase 2 entry)

```go
// Source: D-14/D-15/D-19 + Phase 1 main.go pattern
package main

import (
    _ "embed"
    "fmt"
    "os"

    term "github.com/charmbracelet/x/term"

    "github.com/allanmedeiros71/criminalsay/internal/quotes"
    "github.com/allanmedeiros71/criminalsay/internal/render"
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
    fmt.Print(render.Quote(q, colorEnabled()))
}

func colorEnabled() bool {
    if os.Getenv("NO_COLOR") != "" {
        return false
    }
    return term.IsTerminal(os.Stdout.Fd())
}
```

Note: `fmt.Print` not `Println` — `render.Quote` should include trailing newline if desired, or main adds one. Planner should pick one convention and test consistently.

### Table-Driven render_test (no stdout capture)

```go
// Source: CONTEXT.md §7 + 02-CONTEXT.md D-18 + quotes_test pattern
package render_test

import (
    "strings"
    "testing"

    "github.com/charmbracelet/x/ansi"

    "github.com/allanmedeiros71/criminalsay/internal/quotes"
    "github.com/allanmedeiros71/criminalsay/internal/render"
)

func assertNoANSI(t *testing.T, s string) {
    t.Helper()
    if strings.Contains(s, "\x1b[") {
        t.Fatalf("expected no ANSI escapes, got: %q", s)
    }
}

func TestQuote_noColor(t *testing.T) {
    q := quotes.Quote{
        Quote:        "It's alchemy.",
        Author:       "David Rossi",
        Character:    "David Rossi",
        Season:       9,
        Episode:      13,
        EpisodeTitle: "The Road Home",
    }
    out := render.Quote(q, false)
    assertNoANSI(t, out)
    if !strings.Contains(out, "|") {
        t.Error("expected plain | sidebar in no-color mode (D-18)")
    }
    if strings.Contains(out, "▌") {
        t.Error("unexpected ▌ in no-color mode")
    }
    if !strings.Contains(out, "David Rossi") {
        t.Error("missing author in attribution")
    }
    if !strings.Contains(out, "S9E13") {
        t.Error("missing episode code")
    }
}

func TestQuote_wrapWidth(t *testing.T) {
    long := strings.Repeat("word ", 30)
    q := quotes.Quote{Quote: strings.TrimSpace(long), Author: "A", Character: "A", Season: 1, Episode: 1}
    out := render.Quote(q, false)
    for _, line := range strings.Split(out, "\n") {
        if line == "" || strings.HasPrefix(strings.TrimLeft(line, " "), "·") {
            continue // skip blank and attribution lines
        }
        if strings.Contains(line, "|") || strings.Contains(line, "▌") {
            if w := ansi.StringWidth(line); w > 52 {
                t.Errorf("line exceeds 52 display cols (got %d): %q", w, line)
            }
        }
    }
}

func TestQuote_externalAuthor(t *testing.T) {
    q := quotes.Quote{
        Quote: "Test quote.", Author: "Marcus Aurelius", Character: "Spencer Reid",
        Season: 1, Episode: 1, EpisodeTitle: "Extreme Aggressor",
    }
    out := render.Quote(q, false)
    want := "Marcus Aurelius, cited by Spencer Reid"
    if !strings.Contains(out, want) {
        t.Errorf("missing D-01 attribution %q in:\n%s", want, out)
    }
}
```

### lipgloss Width Behavior (reference — not for sidebar layout)

```go
// Source: go doc lipgloss.Style.Width + style.go
// Width wraps at (width - padding). Useful for attribution-only sub-blocks, NOT quote sidebar.
wrapped := lipgloss.NewStyle().
    Width(48).
    Render("long attribution fragment without sidebar concerns")
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| lipgloss v0.x manual color writer | lipgloss v1.1.0 auto-detect via `DefaultRenderer` + termenv | v1.0+ | No manual `termenv.NewOutput` needed when `colorEnabled=true` |
| `golang.org/x/crypto/ssh/terminal` | `charmbracelet/x/term` | Charm ecosystem consolidation | Use `x/term` (already in go.mod transitively) |
| Global `NO_COLOR` via lipgloss only | Explicit `colorEnabled` param from main | Project D-14/D-15 | Render stays pure/testable; main owns policy |
| COLOR-02 preserve `▌` in plain mode | D-18: use `\|` in plain mode | 02-CONTEXT user override | Tests must assert `\|`, not `▌`, when `colorEnabled=false` |

**Deprecated/outdated:**
- lipgloss v2 import path (`charm.land/lipgloss/v2`) — project locked to v1.1.0
- Detecting color only via lipgloss inside render — superseded by D-14/D-15 explicit bool

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | 52 columns = total display width per quote line (not inner text only) | Pattern 3 | Lines too long/short if planner interprets as inner-only |
| A2 | Attribution indent = 7 spaces per CONTEXT.md §6.1 example | Pattern 4 | Visual mismatch if planner uses different indent |
| A3 | `lipgloss.Color("9")` is acceptable red for REND-01 | Pattern 2 | Minor color shade difference across terminals; acceptable for v1 |
| A4 | `fmt.Print` vs `fmt.Println` for trailing newline | main wiring | Double/missing newline if render and main disagree |

**If A2 is wrong:** Discuss-phase did not lock attribution indent numerically — only "indented". Planner should match CONTEXT.md §6.1 example (7 spaces) unless user corrects.

---

## Open Questions

1. **Trailing newline convention**
   - What we know: Phase 1 used `fmt.Println`; D-19 says render returns string, main prints.
   - What's unclear: Whether `render.Quote` includes final `\n` or main adds it.
   - Recommendation: `render.Quote` returns content **without** trailing newline; `main` uses `fmt.Println(render.Quote(...))` for parity with Phase 1 and simpler tests.

2. **Attribution indent exact count**
   - What we know: CONTEXT.md §6.1 shows 7 spaces; D-06 only numbers quote indent (2).
   - What's unclear: Whether 7 is normative or illustrative.
   - Recommendation: Use 7 spaces to match visual spec; verify in human UAT.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | Build + test | ✓ | go1.25.6 linux/amd64 | — |
| lipgloss v1.1.0 (module cache) | `internal/render` styling | ✓ | v1.1.0 | `go mod download` |
| Color terminal (manual UAT) | COLOR-01 visual check | ✓ (assumed dev machine) | — | Manual smoke only |
| `NO_COLOR` shell env | COLOR-02 test | ✓ | — | Set in test env or `colorEnabled=false` unit tests |

**Missing dependencies with no fallback:** None.

**Missing dependencies with fallback:** None.

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` |
| Config file | none |
| Quick run command | `go test ./internal/render/...` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| CORE-02 | Formatted output string produced | unit | `go test ./internal/render/... -run TestQuote` | ❌ Wave 0 |
| REND-01 | Sidebar on every quote line | unit | `go test ./internal/render/... -run TestQuote_sidebar` | ❌ Wave 0 |
| REND-02 | Lines ≤52 display width | unit | `go test ./internal/render/... -run TestQuote_wrapWidth` | ❌ Wave 0 |
| REND-03 | Quoted text with `"` marks | unit | `go test ./internal/render/... -run TestQuote_quotes` | ❌ Wave 0 |
| REND-04 | Blank line before attribution | unit | `go test ./internal/render/... -run TestQuote_spacing` | ❌ Wave 0 |
| REND-05 | Attribution format + yellow/dim split | unit | `go test ./internal/render/... -run TestQuote_attribution` | ❌ Wave 0 |
| REND-06 | No box border characters (┌┐└┘│─) | unit | `go test ./internal/render/... -run TestQuote_noBox` | ❌ Wave 0 |
| COLOR-01 | Color mode emits ANSI (smoke) | unit | `go test ./internal/render/... -run TestQuote_color` | ❌ Wave 0 |
| COLOR-02 | No-color: no ANSI, `\|` sidebar | unit | `go test ./internal/render/... -run TestQuote_noColor` | ❌ Wave 0 |
| BUILD-03 | Full suite green | integration | `go test ./...` | partial (quotes only) |

### Sampling Rate

- **Per task commit:** `go test ./internal/render/...`
- **Per wave merge:** `go test ./...`
- **Phase gate:** `go build -o criminalsay . && go test ./...` both exit 0

### Wave 0 Gaps

- [ ] `internal/render/render.go` — `Quote()`, attribution builder, styles
- [ ] `internal/render/wrap.go` — 52-col word wrap helpers
- [ ] `internal/render/render_test.go` — table-driven tests, `colorEnabled` true/false cases
- [ ] `main.go` update — replace stub, add `colorEnabled()`, import render

---

## Security Domain

> `security_enforcement: true`, `security_asvs_level: 1` per `.planning/config.json`.

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | No | Print-and-exit CLI |
| V3 Session Management | No | Stateless |
| V4 Access Control | No | No user accounts |
| V5 Input Validation | Minimal | Quote data is compile-time embedded; render formats trusted struct fields |
| V6 Cryptography | No | No secrets |
| V9 Communication | No | No network I/O |
| V13 API | No | No HTTP surface |

### Known Threat Patterns for this Stack

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| ANSI injection via quote text | Spoofing | Quote data is curated JSON; `\x1b` in quote strings could inject terminal controls | Data curation; optional `\x1b` strip in render if untrusted input ever added (out of scope v1) |
| Terminal escape sequences in output | Tampering | lipgloss preserves existing ANSI in wrapped strings; embedded quotes are trusted |

**Security summary:** Attack surface unchanged from Phase 1. Trusted embedded data only. ASVS L1 satisfied.

---

## Project Constraints (from .cursor/rules/)

No `.cursor/rules/` directory exists in this repository. No additional Cursor rule constraints beyond user rules and `.planning/` artifacts.

---

## Sources

### Primary (HIGH confidence — verified via installed module source and `go doc`)

- `github.com/charmbracelet/lipgloss@v1.1.0` — `style.go` (Width/wrap via cellbuf), `borders.go` (per-line left border), `renderer.go` (ColorProfile), `set.go` (Width, MaxWidth, Faint) [VERIFIED: local module + go doc]
- `github.com/muesli/termenv@v0.16.0` — `EnvNoColor`, `EnvColorProfile`, `Profile.Ascii` [VERIFIED: local module source]
- `github.com/charmbracelet/x/term@v0.2.1` — `IsTerminal` [VERIFIED: go doc]
- `github.com/charmbracelet/x/cellbuf` — `Wrap` [VERIFIED: local module source]
- `go list -m -json github.com/charmbracelet/lipgloss@v1.1.0` — version timestamp [VERIFIED: go module proxy]

### Secondary (MEDIUM confidence — project locked artifacts)

- `.planning/phases/02-render-and-output/02-CONTEXT.md` — D-01 through D-20
- `CONTEXT.md` §6.1 — Estilo 4 visual spec, reference output, test strategy §7
- `.planning/REQUIREMENTS.md` — CORE-02, REND-01–06, COLOR-01–02, BUILD-03
- `.planning/phases/01-data-foundation/01-RESEARCH.md` — embed anchor, test patterns, module layout

### Tertiary (LOW confidence — illustrative)

- [CITED: pkg.go.dev/github.com/charmbracelet/lipgloss@v1.1.0] — color ID examples (README shows v2 path; color semantics apply to v1)

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — lipgloss v1.1.0 installed and API-verified
- Architecture: HIGH — layout rules locked in 02-CONTEXT.md; wrap/sidebar pattern derived from spec
- Pitfalls: HIGH — grounded in D-decisions and module source behavior
- Attribution indent (7 spaces): MEDIUM — from visual example, not a numbered D-decision

**Research date:** 2026-07-05
**Valid until:** 2026-08-05 (lipgloss v1.x stable; layout locked)

---

## RESEARCH COMPLETE

**Phase:** 02 - Render and Output
**Confidence:** HIGH

### Key Findings

- Use **manual per-line sidebar prefix** after word-wrap — not `lipgloss.Border` — to satisfy D-06/D-07/D-08/D-20 simultaneously.
- **52 columns = total display width** per quote line (matches CONTEXT.md §6.1 example); verify with `ansi.StringWidth` in tests.
- **Color gating is explicit:** `main` checks `NO_COLOR` + `term.IsTerminal(os.Stdout.Fd())`, passes `colorEnabled` to `render.Quote`; no-color path uses `|` sidebar (D-18) and zero lipgloss styling.
- **Attribution:** build plain text per D-01–D-13, wrap at 52 cols, split yellow/dim at ` · Criminal Minds` per line after wrap.
- **Tests:** table-driven `render_test` package on return string; `assertNoANSI` when `colorEnabled=false`; no stdout capture.

### File Created

`.planning/phases/02-render-and-output/02-RESEARCH.md`

### Confidence Assessment

| Area | Level | Reason |
|------|-------|--------|
| Standard Stack | HIGH | lipgloss v1.1.0 verified in go.mod and module cache |
| Architecture | HIGH | Locked D-01–D-20 + verified lipgloss wrap/border behavior |
| Pitfalls | HIGH | Derived from spec + common lipgloss misuse patterns |

### Open Questions

- Trailing newline: recommend `fmt.Println` in main, render returns string without final `\n`.
- Attribution indent (7 spaces): from visual spec example; confirm in UAT.

### Ready for Planning

Research complete. Planner can now create PLAN.md files.
