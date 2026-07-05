# Feature Landscape

**Domain:** fortune-style CLI quote display tool (Go, Criminal Minds theme)
**Researched:** 2026-07-05
**Overall confidence:** MEDIUM (web sources LOW; lipgloss docs MEDIUM via context7)

---

## Table Stakes

Features users expect from a fortune-style tool. Missing = product feels incomplete or broken.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Random quote selection | Core value proposition — every run produces something different | Low | `math/rand/v2` with simple index pick; already decided (D7) |
| Quote + attribution display | Users expect "who said it" alongside the quote; anonymous quotes feel incomplete | Low | JSON struct with `text` and `speaker` fields |
| Word wrap at terminal width | Raw print without wrap causes ragged terminal breaks; looks amateur | Low-Med | lipgloss `Width()` on style + `TerminalSize(os.Stdout)`; fallback 80 cols |
| Sidebar / visual framing | Distinguishes from plain `echo`; the `▌` sidebar is the signature style of Estilo 4 | Low | lipgloss border or manual `▌` prefix per line |
| Color that degrades gracefully | Fortune tools piped or redirected must not emit ANSI garbage; NO_COLOR and non-TTY must produce clean text | Low | lipgloss does this automatically when not a TTY; NO_COLOR respected natively |
| Embedded quote data (no external deps) | Binary must be self-contained; users expect `go install` and done — no config files | Low | `go:embed` on `data/quotes.json`; already decided (D5/D6) |
| ≥15 real quotes from the series | Enough variety that the tool doesn't feel like a stub; repeat probability <1-in-15 per run | Low | Curation work, not code complexity |
| Exit code 0 on success | Shell scripts and `.bashrc` sourcing depend on clean exit; non-zero would suppress output in many setups | Trivial | `os.Exit(0)` or just `return` from main |
| Cold start <50ms | `.bashrc` usage means every terminal open blocks on this; slow startup = removed from profile | Low | Embed + single JSON parse + print; no I/O, no network — trivially achievable |

---

## Differentiators

Features that set this tool apart. Not expected in v1, but would add real value in v2.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Filter by character (`--character`) | Power users want quotes from Gideon specifically; unlocks persona-driven mood | Med | Requires JSON field `speaker`, filter pass before random selection; v2 candidate |
| Filter by season (`--season`) | Fans want era-specific quotes; useful for episode rewatch periods | Med | Requires `season` field in JSON schema; v2 candidate |
| Multiple visual themes (`--style`) | Different aesthetics for different contexts (minimal vs. boxed vs. colorful) | Med | Parameterize the render package; v2 candidate |
| Pipe-friendly plain output (`--plain`) | Enables chaining with `cowsay`, `lolcat`, or other tools via Unix pipeline | Low | Detect non-TTY (already table stakes) or explicit flag |
| Quote count display (`--count`) | Shows how many quotes are in the database; useful for contributors | Trivial | `len(quotes)` print; nearly free |
| Deterministic selection (`--seed`) | Reproducibility for screenshots, demos, testing without mocking | Low | Pass seed to `rand.New(rand.NewSource(seed))` |
| Distribution via Homebrew/AUR | Discoverability for non-Go users who won't `go install` | Med | Packaging work, not code; already explicitly deferred in PROJECT.md |
| Contributor-friendly quote addition | JSON schema + validation script lowers barrier to community PRs | Low | JSON schema file + CI validation; v2 |

---

## Anti-Features

Features to deliberately NOT build in v1 (and in some cases ever).

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| Interactive TUI / browseable quote list | Defeats the zero-friction goal; bubbletea adds 200ms+ startup and dependency bloat | Display and exit; explicitly out of scope (D4) |
| Runtime scraping / API fetch | Network call in `.bashrc` = broken offline, slow on flaky connections, privacy concern | All data embedded at build time (D5) |
| Multiple simultaneous outputs | Fortune outputs one quote per invocation — that's the contract | Single quote, single render, done |
| strfile / % delimiter format | Legacy format; binary index file means two files, not one; JSON is cleaner for Go embed | JSON with structured fields (`text`, `speaker`) |
| Cowsay-style ASCII art framing | ASCII cows are charming but not on-brand for Criminal Minds; sidebar is cleaner | Sidebar `▌` with lipgloss (Estilo 4 already locked) |
| Rainbow/lolcat coloring | Too playful; contradicts the dark/serious tone of Criminal Minds quotes | Single accent color (red sidebar) per Estilo 4 |
| Config file (`~/.criminalsayrc`) | Config files require migration, documentation, and error handling; zero benefit for a display-only tool | Env vars (`NO_COLOR`) and future flags are sufficient |
| External fortune database compatibility | Supporting `fortune -f` file paths adds complexity; this is a curated Criminal Minds tool, not a general fortune replacement | Curated embedded data only |
| Localization / i18n | Quotes are in English/Portuguese; translation infrastructure is overkill for v1 | Quotes stored verbatim; display as-is |
| Semantic versioning of quote schema | JSON schema versioning adds migration code; quotes don't need migrations | Flat JSON array; schema changes are additive |

---

## Feature Dependencies

```
Random selection → Embedded quote data
  (selection requires the pool to be loaded)

Word wrap → Terminal width detection
  (wrap must know the target column count)

Word wrap → Sidebar alignment
  (each wrapped line must carry the sidebar prefix to keep visual consistency)

Graceful color degradation → TTY detection
  (lipgloss auto-strips; but NO_COLOR env var must also be checked)

Sidebar rendering → Word wrap
  (word wrap must happen BEFORE sidebar prefix is prepended, or wrap budget must account for prefix width)

Filter by character (v2) → JSON schema with `speaker` field
  (table stakes already requires speaker for attribution; v2 filter is additive)

Filter by season (v2) → JSON schema with `season` field
  (season is NOT required for v1; must be added to JSON before v2 filter works)
```

---

## MVP Recommendation

**Build in v1 (all table stakes):**

1. JSON quote loader with `go:embed` — the data foundation everything else rests on
2. Random selection from the loaded pool — the core behavior
3. Terminal width detection with 80-column fallback — must happen before render
4. Word wrap respecting width and sidebar prefix width — wrap budget = termWidth - 2 (for `▌ `)
5. Estilo 4 render: red `▌` sidebar, quoted text, dimmed attribution — the visual identity
6. Graceful color degradation — NO_COLOR + non-TTY strips all ANSI; lipgloss handles this automatically
7. ≥15 real Criminal Minds quotes curated into `data/quotes.json`

**Defer to v2:**

- `--character` / `--season` filter flags (schema already accommodates speaker; season needs adding)
- Multiple themes (`--style`)
- `--plain` explicit flag (non-TTY detection already covers the pipe use case)
- Homebrew/AUR distribution

**Never build:**

- TUI / interactive mode
- Runtime network fetching
- Config files
- Cowsay / lolcat style output

---

## Key Implementation Notes

**Word wrap ordering matters.** Wrap the raw text first (at `termWidth - prefixLen`), then prepend the sidebar character to each line. If you wrap after adding ANSI, the escape sequences inflate the string length and wrapping breaks. lipgloss `Wrap()` is ANSI-aware and handles this correctly if the whole styled block goes through lipgloss.

**Sidebar alignment with multi-line quotes.** Each wrapped line needs the `▌ ` prefix, not just the first. This means splitting the quote into lines after wrapping and iterating — or using lipgloss `BorderLeft` which applies the border to every line automatically.

**Test the 80-column fallback.** The fallback path (no TTY, no terminal size) is exercised in CI where stdout is piped. Tests should assert that the output without a TTY is plain text, no ANSI.

---

## Sources

- [fortune(6) Linux man page](https://linux.die.net/man/6/fortune) — command flags (MEDIUM confidence, cross-verified)
- [fortune (Unix) — Wikipedia](https://en.wikipedia.org/wiki/Fortune_(Unix)) — history and data format (LOW confidence)
- [GitHub — shlomif/fortune-mod](https://github.com/shlomif/fortune-mod) — strfile format reference (LOW confidence)
- [Terminal Joy: Fortune, Cowsay, Figlet, Lolcat — Medium](https://medium.com/@Smyekh/terminal-joy-how-fortune-cowsay-figlet-and-lolcat-add-life-to-my-developer-workflow-b5b1c6b10474) — ecosystem patterns (LOW confidence)
- [charmbracelet/lipgloss — GitHub](https://github.com/charmbracelet/lipgloss) — color degradation, word wrap, terminal size (MEDIUM confidence)
- [fatih/color — Go Packages](https://pkg.go.dev/github.com/fatih/color) — NO_COLOR and TTY detection (LOW confidence)
- [Terminal Color Detection — terminfo.dev](https://terminfo.dev/fundamentals/color-detection) — NO_COLOR standard (LOW confidence)
- [Go CLI tutorial: fortune clone — flaviocopes.com](https://flaviocopes.com/go-tutorial-fortune/) — Go implementation patterns (LOW confidence)
- [go-wordwrap — mitchellh/go-wordwrap](https://github.com/mitchellh/go-wordwrap) — word wrapping library reference (LOW confidence)
