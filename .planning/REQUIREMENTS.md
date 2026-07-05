# Requirements — CriminalSay

## v1 Requirements

### Data (DAT)

- [ ] **DAT-01**: Binary contains an embedded JSON file with ≥15 real Criminal Minds quotes compiled at build time (go:embed)
- [ ] **DAT-02**: Each quote record includes: `quote` (text), `author` (original author), `character` (CM character who cited it; may be empty), `season` (int), `episode` (int), `episodeTitle` (string)
- [ ] **DAT-03**: Missing optional fields (character, episodeTitle) are rendered gracefully — omitted without breaking layout

### Core Behavior (CORE)

- [ ] **CORE-01**: Running `criminalsay` with no arguments selects one quote uniformly at random from the embedded dataset
- [ ] **CORE-02**: Program prints the formatted quote to stdout and exits with code 0
- [ ] **CORE-03**: Cold start completes in under 50ms under typical conditions (suitable for .bashrc/.zshrc)

### Rendering (REND)

- [ ] **REND-01**: Output uses Estilo 4 "Terminal cru": a red `▌` (U+258C) sidebar prefix on every line of the quote text, with two columns of left indentation
- [ ] **REND-02**: Quote text is wrapped at a fixed block width (≤52 useful columns), with the `▌` sidebar repeated and aligned on each wrapped line
- [ ] **REND-03**: Quote text is displayed between double quotation marks in a highlighted color (white/bright)
- [ ] **REND-04**: A blank line separates the quote block from the attribution line
- [ ] **REND-05**: Attribution line is indented and formatted as `<author> · Criminal Minds · S<season>E<episode>`, with author in yellow and the rest in dim/muted tone, separated by `·` (U+00B7)
- [ ] **REND-06**: No box border, no header, no footer — lightweight presentation only

### Color & Compatibility (COLOR)

- [ ] **COLOR-01**: When stdout is a TTY with color support, output uses ANSI color sequences (truecolor/256/ANSI-16 auto-detected by lipgloss)
- [ ] **COLOR-02**: When `NO_COLOR` is set, stdout is not a TTY, or output is piped, ANSI sequences are stripped; text remains readable with sidebar `▌`, quotes, and separators preserved

### Build & Distribution (BUILD)

- [ ] **BUILD-01**: `go build -o criminalsay .` produces a working binary from the repository root
- [ ] **BUILD-02**: Cross-compile script (`Makefile` or `scripts/build.sh`) produces binaries for: Linux amd64, macOS amd64, macOS arm64, Windows amd64; all with CGO_ENABLED=0
- [ ] **BUILD-03**: `go test ./...` passes with reasonable coverage of `internal/quotes` and `internal/render`

### Documentation (DOC)

- [ ] **DOC-01**: README in PT-BR with installation instructions, usage, and build instructions

---

## v2 Requirements (Deferred)

- Filter quotes by character (`--character`) — power-user feature; deferred to v2
- Filter quotes by season (`--season`) — era-specific browsing; deferred to v2
- Multiple selectable visual themes (`--style`) — differentiator; deferred to v2
- `--plain` explicit flag — niche need; covered by NO_COLOR; deferred
- `--seed` for deterministic output — useful for debugging; deferred
- Homebrew/AUR/package distribution — discoverability; deferred to v2

---

## Out of Scope

- Interactive TUI / bubbletea — explicitly discarded (D4); escopo é exibir e sair
- Runtime network fetching / scraping — discarded (D5); offline-first by design
- Config files — unnecessary complexity for a print-and-exit tool
- cowsay-style ASCII art framing — aesthetic choice; Estilo 4 is locked for v1 (D10)
- i18n / PT-BR quote translation — original English quotes are a feature, not a limitation

---

## Traceability

| REQ-ID | Phase | Notes |
|--------|-------|-------|
| DAT-01 | Phase 1 | go:embed anchor in main.go; path restriction means data/ at root |
| DAT-02 | Phase 1 | Quote struct definition frozen before render work begins |
| DAT-03 | Phase 1 | Handled by render layer — graceful omission |
| CORE-01 | Phase 1 | quotes.Random([]Quote) using math/rand/v2 |
| CORE-02 | Phase 2 | main.go wiring; fmt.Println(render.Quote(...)) |
| CORE-03 | Phase 1+2 | Embedded data + single JSON unmarshal; no network |
| REND-01 | Phase 2 | lipgloss Border{Left:"▌"} with red BorderForeground |
| REND-02 | Phase 2 | Wrap at (fixedWidth - sidebarWidth); sidebar repeated per line |
| REND-03 | Phase 2 | lipgloss style with white/bright foreground |
| REND-04 | Phase 2 | Blank line between quote block and attribution |
| REND-05 | Phase 2 | Attribution: yellow author + dim rest + · separator |
| REND-06 | Phase 2 | No box border — Estilo 4 is sidebar-only |
| COLOR-01 | Phase 2 | lipgloss auto-detects color depth from stdout |
| COLOR-02 | Phase 2 | lipgloss strips ANSI for non-TTY/NO_COLOR automatically |
| BUILD-01 | Phase 1+2 | Verified at each phase; final binary in Phase 3 |
| BUILD-02 | Phase 3 | Makefile with CGO_ENABLED=0 for all targets |
| BUILD-03 | Phase 1+2 | Tests written alongside implementation per package |
| DOC-01 | Phase 3 | README in PT-BR |
