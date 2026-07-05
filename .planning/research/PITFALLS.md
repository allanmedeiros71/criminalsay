# Domain Pitfalls

**Domain:** Go CLI tool with lipgloss terminal styling and go:embed data embedding
**Researched:** 2026-07-05
**Overall confidence:** MEDIUM (websearch LOW + context7 MEDIUM; cross-checked across official Go issue tracker and charmbracelet docs)

---

## Critical Pitfalls

Mistakes that cause garbled output, broken builds, or rewrites.

---

### Pitfall 1: Using lipgloss v1 API with v2 (or vice versa) — the Renderer was removed

**What goes wrong:**
In lipgloss v1, every `Style` carries a `*Renderer` pointer and the package exposes `lipgloss.DefaultRenderer()`, `lipgloss.SetColorProfile(p)`, and `lipgloss.ColorProfile()`. In v2, the `Renderer` type is gone. Color downsampling no longer happens inside `Style.Render()` — it now happens at print time. If you copy v1 snippets into a v2 project (or vice versa), the code does not compile and the mental model for color degradation is wrong.

**Why it happens:**
Most blog posts and Stack Overflow answers target v1. The v2 release announcement is a GitHub Discussion (#506), not a blog post, so it is easy to miss.

**Consequences:**
- Compile errors on `lipgloss.NewRenderer`, `lipgloss.SetColorProfile`, `lipgloss.DefaultRenderer`.
- If you use the v2 compat shim (`lipgloss/compat`), colors may be stripped when stdout is redirected because the shim reads stdin/stdout globally — the same bug v2 was designed to fix.

**Prevention:**
- Lock the import path explicitly: use `github.com/charmbracelet/lipgloss/v2` in `go.mod`.
- Do NOT mix import paths (`lipgloss` and `lipgloss/v2`) in the same binary.
- For a one-shot CLI (no Bubble Tea), use the v2 writer pattern: `lipgloss.NewWriter(os.Stdout)` to get a `Writer` that handles color downsampling automatically.
- Read `UPGRADE_GUIDE_V2.md` in the lipgloss repo before writing any styling code.

**Detection:**
- `go build` will error immediately on missing `Renderer` symbol.
- If colors always render at full truecolor in tests but strip when piped: the compat shim is in use.

**Phase:** Core render package (Phase 1 / Phase 2).

---

### Pitfall 2: Color profile auto-detection fails silently — colors render as garbled escape codes

**What goes wrong:**
lipgloss detects the color profile by inspecting environment variables (`COLORTERM`, `TERM`, `NO_COLOR`) and whether stdout is a TTY. If the detection is wrong (e.g., stdout is a pipe, a CI runner, or Windows Terminal without `WT_SESSION`), one of two bad things happens:
(a) ANSI escape codes are emitted to a pipe/file that cannot render them — garbled output.
(b) Colors are stripped entirely even in a capable terminal.

**Why it happens:**
- Windows Terminal does not set `COLORTERM=truecolor`; it advertises via `WT_SESSION`.
- tmux/screen change `TERM` to `screen` or `tmux-256color`, confusing profile detection.
- When piped into `.bashrc` / `.zshrc`, stdout is not a TTY at the point lipgloss evaluates it.

**Consequences:**
The `▌` sidebar and styled quote bleed raw escape codes into the shell startup output, making every shell open show `\e[31m▌\e[0m` literally — a terrible default experience.

**Prevention (lipgloss v2):**
```go
w := lipgloss.NewWriter(os.Stdout)
// The Writer auto-detects profile on creation and
// strips ANSI when stdout is not a TTY.
```
Also handle `NO_COLOR` explicitly: if `os.Getenv("NO_COLOR") != ""`, render plain text with no ANSI.

**Prevention (lipgloss v1, if chosen):**
```go
if lipgloss.ColorProfile() == termenv.Ascii {
    // fall back to plain rendering
}
```

**Detection (warning signs):**
- Running `criminalsay | cat` produces raw escape codes.
- Users report garbled output in CI or cron jobs.
- `echo $COLORTERM` is empty on the target terminal.

**Phase:** Core render package; also test in CI (where TTY is absent) from the start.

---

### Pitfall 3: go:embed silently excludes files starting with `.` or `_`

**What goes wrong:**
When you write `//go:embed data` (embedding a whole directory), Go excludes any file whose name begins with `.` (dot) or `_` (underscore) — silently, with no warning. If `quotes.json` is named `_quotes.json` or lives next to hidden files you intended to embed, those files are absent at runtime and the program panics or returns an empty quote list.

**Why it happens:**
Go deliberately mirrors its build-file exclusion rules: files starting with `.` or `_` are ignored by the Go toolchain, so `go:embed` follows the same convention.

**Consequences:**
- `json.Unmarshal` receives an empty `[]byte{}` — program either panics or silently has zero quotes.
- The bug only manifests at runtime, not at build time.

**Prevention:**
- Name the data file `quotes.json` (no leading dot or underscore). ✔ The project already does this.
- If you ever need to embed hidden files, use the `all:` prefix: `//go:embed all:data`.
- Add a startup assertion: `if len(quotesData) == 0 { panic("embed failed: quotes.json is empty") }`.

**Detection:**
- `go build` succeeds; `len(quotesData) == 0` at runtime.
- Running `go list -f '{{.EmbedFiles}}' ./...` shows what was actually embedded.

**Phase:** Data / quotes package (Phase 1). Verify with `go list` immediately after writing the directive.

---

### Pitfall 4: go:embed path cannot reference files outside the package directory

**What goes wrong:**
The `//go:embed` directive is resolved relative to the Go source file containing it. Patterns may not use `..` to ascend to a parent directory. If you place the embed directive in `internal/quotes/loader.go` and the JSON file is at `data/quotes.json` (at the module root), the build fails:

```
pattern ../../../data/quotes.json: invalid pattern syntax
```

**Why it happens:**
Go enforces that all embedded files live within the package's own subtree. This is by design to keep builds reproducible and prevent embedding files outside the module.

**Consequences:**
Build breaks. The fix requires moving the data file or the embed directive to align with the restriction.

**Prevention:**
Two valid layouts:
1. Put `data/quotes.json` inside `internal/quotes/data/quotes.json` and embed from `internal/quotes/`.
2. Embed from `main.go` (module root) with `//go:embed data/quotes.json` and pass the `[]byte` into the quotes package.

The project's stated layout (`data/quotes.json` at root, embed in `internal/quotes/`) **will fail**. Choose layout 1 (move data under the package) or layout 2 (embed in main.go).

**Detection:**
- `go build` errors immediately: `pattern data/quotes.json: cannot embed from outside module source dir`.

**Phase:** This must be resolved before writing any code — it is an architectural constraint. Resolve in Phase 1.

---

### Pitfall 5: Line-wrapping with sidebar character — wrapping at wrong width breaks alignment

**What goes wrong:**
You wrap the quote text to terminal width, then prepend `▌ ` to each line. If you wrap at the full terminal width (e.g., 80), the resulting lines + sidebar overshoot the terminal, causing visual wrapping that misaligns the sidebar on the next terminal line.

**Why it happens:**
The sidebar glyph `▌` is 1 cell wide, plus a space = 2 cells. Wrapping at `termWidth` then prepending 2 cells makes each line `termWidth + 2` cells wide, overflowing.

**Consequences:**
- On narrow terminals (60–80 columns), the quote text wraps mid-word with the sidebar hanging orphaned on the previous line.
- On wide terminals, it looks correct — making the bug easy to miss during development.

**Prevention:**
Wrap at `termWidth - sidebarWidth - margin` before building the styled output:
```go
const sidebarWidth = 2 // "▌ "
wrapWidth := termWidth - sidebarWidth - 1 // 1 cell safety margin
wrapped := wordwrap.String(quote.Text, wrapWidth)
lines := strings.Split(wrapped, "\n")
```
Then join with the sidebar prepended to each line using lipgloss styles.

Use `golang.org/x/term` to detect terminal width:
```go
termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
if err != nil || termWidth == 0 {
    termWidth = 80 // safe fallback
}
```

**Detection:**
- Test on a terminal set to 60 columns: `stty cols 60 && criminalsay`.
- Visual misalignment on short lines is the symptom.

**Phase:** Render package. Write a unit test with a fixed width of 40 columns.

---

### Pitfall 6: Terminal width detection returns 0 or error when stdout is not a TTY

**What goes wrong:**
`term.GetSize(int(os.Stdout.Fd()))` returns an error when stdout is a pipe, a file, or a non-TTY (CI runner, cron, `criminalsay > file.txt`). If not handled, `termWidth` is 0 and word-wrap panics or produces a single unbroken line.

**Why it happens:**
`GetSize` calls the OS `ioctl TIOCGWINSZ` which only works on a real TTY file descriptor.

**Consequences:**
- With `termWidth = 0` and no guard, `wordwrap.String(text, 0)` either panics or wraps at every character.
- With `termWidth < 0` (some errors), negative-width styles produce layout corruption.

**Prevention:**
Always guard:
```go
if err != nil || termWidth <= 0 {
    termWidth = 80
}
```
Also: when stdout is not a TTY, color is already stripped (see Pitfall 2), so the wrapping width issue is less visible — but still produces misformatted plain text.

**Detection:**
- `criminalsay | cat` triggers this path. Add this to CI tests.

**Phase:** Render package.

---

## Moderate Pitfalls

---

### Pitfall 7: math/rand/v2 — seeding is automatic; do NOT call Seed()

**What goes wrong:**
Developers familiar with `math/rand` (v1) call `rand.Seed(time.Now().UnixNano())` at startup. In `math/rand/v2` (Go 1.22+), the `Seed()` top-level function does not exist. The global source is always auto-seeded with a cryptographically random value. Calling `rand.Seed` on v2 causes a compile error.

**Consequence of doing nothing (v2):** Works correctly — random is already seeded. No action required.
**Consequence of importing v1 by mistake:** Sequences are deterministic unless you manually seed, making every run return the same quote (the first or a fixed one).

**Prevention:**
- Import `math/rand/v2`, not `math/rand`.
- Do not call `rand.Seed()` — remove it entirely.
- Random quote selection: `quotes[rand.IntN(len(quotes))]` (note `IntN`, not `Intn`).

**Detection:**
- `go build` errors on `rand.Seed` if using v2 (catches it at compile time).
- All runs return the same quote = v1 imported without seeding.

**Phase:** Quotes package.

---

### Pitfall 8: Cross-compile with CGO_ENABLED=1 (default) breaks non-native targets

**What goes wrong:**
By default, `CGO_ENABLED=1`. When cross-compiling (`GOOS=windows GOARCH=amd64 go build`), the Go toolchain tries to link against a C compiler for the target platform. If no cross-compiler (e.g., MinGW for Windows) is in `$PATH`, the build fails with a linker error.

**Why it happens:**
lipgloss and its dependencies are pure Go — but if any transitive dependency uses cgo, it breaks cross-compile. Additionally, `mattn/go-isatty` (used by lipgloss's dependencies) historically had a cgo path on some platforms.

**Consequences:**
- `go build` fails for non-native targets in CI.
- GitHub Actions matrix builds fail on macOS-to-Linux or Linux-to-Windows cross-compile steps.

**Prevention:**
Set `CGO_ENABLED=0` explicitly for all cross-compile targets:
```makefile
build-linux:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/criminalsay-linux-amd64 .

build-windows:
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/criminalsay-windows-amd64.exe .

build-macos-arm:
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dist/criminalsay-darwin-arm64 .
```

**Detection:**
- Build fails with `exec: "x86_64-linux-gnu-gcc": executable file not found`.
- Adding `CGO_ENABLED=0` fixes it immediately.

**Phase:** Build infrastructure (Makefile / CI). Validate all three targets before shipping.

---

### Pitfall 9: Windows ANSI VT processing not enabled by default in legacy consoles

**What goes wrong:**
Windows 10 added VT/ANSI escape sequence support in build 16257 (2017), but it is **opt-in** — a process must call `SetConsoleMode` with `ENABLE_VIRTUAL_TERMINAL_PROCESSING`. If this is not done, raw escape codes appear literally in `cmd.exe`. lipgloss v2 does not enable VT mode on Windows by itself.

**Consequence:**
Users running `criminalsay` in `cmd.exe` or older PowerShell see `\e[31m▌\e[0m` instead of a red sidebar.

**Prevention:**
Use `golang.org/x/sys/windows` or the `mattn/go-colorable` package to enable VT mode on Windows before any rendering:
```go
// In main(), Windows-only build tag:
import "golang.org/x/sys/windows"
windows.SetConsoleMode(windows.Handle(os.Stdout.Fd()),
    windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
```
Alternatively: lipgloss with `NO_COLOR` or non-TTY detection will strip colors in `cmd.exe` (where stdout is not a VT-capable TTY), producing safe plain output automatically — acceptable for v1.

**Detection:**
- Test in `cmd.exe` (not Windows Terminal): raw escape codes appear.
- `WT_SESSION` env var is set in Windows Terminal; absent in `cmd.exe`.

**Phase:** Render package — add Windows-specific VT enablement or verify lipgloss's non-TTY fallback is acceptable.

---

### Pitfall 10: lipgloss.Width() / lipgloss.Height() measure rendered ANSI, not visual width — misuse causes off-by-one padding

**What goes wrong:**
Developers use `lipgloss.Width(someString)` to measure the visual width of already-styled text to compute padding or alignment. However, if `someString` contains ANSI escape sequences, `lipgloss.Width` strips them and returns the visual (printable) width — which is correct. But if you apply this to unstyled text that contains multi-byte Unicode or emoji, the measurement may be off because of CJK double-width characters.

For this project (ASCII quotes, Latin characters), this is a minor risk. However, if any quote contains Unicode punctuation (em-dash `—`, smart quotes `""`), the width is still measured correctly since those are single-cell.

**Prevention:**
- Keep quote text to printable Latin/UTF-8 single-width characters.
- If quotes ever include CJK or emoji, set `RUNEWIDTH_EASTASIAN=0` in the environment or pin `go-runewidth` to a version that handles it.
- Use `lipgloss.Width(renderedText)` not `len([]rune(text))` for alignment math.

**Phase:** Render package — low risk for this project but worth noting.

---

## Minor Pitfalls

---

### Pitfall 11: Blank import `_ "embed"` is required even if you only use `//go:embed` with a string or []byte variable

**What goes wrong:**
When embedding into a `string` or `[]byte` variable, developers sometimes omit the `import _ "embed"` because they think "I'm not using the `embed.FS` type". The compiler requires the blank import regardless:

```go
import _ "embed"

//go:embed data/quotes.json
var quotesData []byte  // requires _ "embed" even though embed.FS is not used
```

Without it: `go build` fails: `//go:embed only allowed in Go files that import "embed"`.

**Prevention:** Always include `import _ "embed"` in any file using the directive.

**Phase:** Quotes package.

---

### Pitfall 12: go:embed directive must immediately precede the variable declaration — no blank lines between

**What goes wrong:**
```go
//go:embed data/quotes.json

var quotesData []byte  // blank line between directive and var = silent failure
```
The compiler ignores the directive if there is a blank line between it and the `var` declaration. The variable is initialized to `nil`/empty. No error is emitted.

**Prevention:** The directive and the `var` must be adjacent (only `//` line comments allowed between them).

**Detection:** `len(quotesData) == 0` at runtime despite a valid file path.

**Phase:** Quotes package.

---

### Pitfall 13: lipgloss Style is a value type — copying styles is safe, mutating the "same" style is not a problem, but accidental sharing of a `*Style` pointer is

**What goes wrong:**
lipgloss `Style` is a struct (value type), so assignment copies it:
```go
base := lipgloss.NewStyle().Foreground(red)
variant := base.Bold(true) // base is unmodified
```
This is safe. The danger is storing `*Style` pointers and mutating them — but lipgloss's fluent API always returns new values, so this is hard to trigger accidentally. The real pitfall is creating styles in a hot loop unnecessarily (each `lipgloss.NewStyle()` call allocates).

**Prevention:** Create `Style` values once at package init or `var` level, not inside the render function that runs per quote.

**Phase:** Render package — negligible for a one-shot CLI but good practice.

---

### Pitfall 14: `os.Exit` in main skips deferred functions — defer is not useful for cleanup

**What goes wrong:**
If an error path calls `os.Exit(1)`, any `defer` statements do not run. For a one-shot CLI that just prints and exits, this is usually fine — but if you defer file handle closes or buffer flushes, they will be skipped.

**Prevention:** For this project, avoid relying on deferred flushes. Call `os.Exit` only after all output is written (lipgloss output goes directly to stdout, no buffering needed).

**Phase:** main.go.

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Package layout + go:embed | Embed path outside package dir (Pitfall 4) | Resolve layout before writing code |
| Data file naming | Dot/underscore exclusion (Pitfall 3) | Use `quotes.json`; verify with `go list -f '{{.EmbedFiles}}'` |
| Blank import | Missing `_ "embed"` (Pitfall 11) | Linter catches it; add immediately |
| Directive syntax | Blank line between directive and var (Pitfall 12) | Keep them adjacent |
| Render package — color | lipgloss v1 vs v2 API mismatch (Pitfall 1) | Pin v2 in go.mod; read upgrade guide |
| Render package — color | Profile detection / piped output (Pitfall 2) | Use `lipgloss.NewWriter`; test with pipe |
| Render package — wrapping | Sidebar width not subtracted from wrap width (Pitfall 5) | Use `termWidth - sidebarWidth` |
| Render package — wrapping | Zero/error terminal width (Pitfall 6) | Guard with `termWidth = 80` fallback |
| Quotes package — randomness | Importing v1 rand instead of v2 (Pitfall 7) | Import `math/rand/v2`; no Seed() call |
| Build / CI — cross-compile | CGO_ENABLED=1 breaks non-native (Pitfall 8) | Set `CGO_ENABLED=0` in all build targets |
| Windows output | ANSI VT not enabled (Pitfall 9) | Enable VT mode or rely on non-TTY degradation |

---

## Sources

- [lipgloss GitHub repo](https://github.com/charmbracelet/lipgloss) (confidence: MEDIUM)
- [Lip Gloss v2: What's New — Discussion #506](https://github.com/charmbracelet/lipgloss/discussions/506) (confidence: MEDIUM)
- [lipgloss UPGRADE_GUIDE_V2.md](https://github.com/charmbracelet/lipgloss/blob/main/UPGRADE_GUIDE_V2.md) (confidence: MEDIUM)
- [embed package — pkg.go.dev](https://pkg.go.dev/embed) (confidence: MEDIUM)
- [go:embed path restrictions — Issue #58519](https://github.com/golang/go/issues/58519) (confidence: MEDIUM)
- [go:embed dot/underscore exclusion — Issue #43854](https://github.com/golang/go/issues/43854) (confidence: MEDIUM)
- [math/rand/v2 blog post — go.dev](https://go.dev/blog/randv2) (confidence: MEDIUM)
- [CGO cross-compile — Issue #5104](https://github.com/golang/go/issues/5104) (confidence: MEDIUM)
- [Emoji/Unicode Width — Issue #562](https://github.com/charmbracelet/lipgloss/issues/562) (confidence: LOW — websearch)
- [ColorProfile not detected — Issue #439](https://github.com/charmbracelet/lipgloss/issues/439) (confidence: LOW — websearch)
- [Windows ANSI VT — go-supportscolor](https://github.com/jwalton/go-supportscolor) (confidence: LOW — websearch)
- [terminal width pattern — lipgloss Discussion #430](https://github.com/charmbracelet/lipgloss/discussions/430) (confidence: LOW — websearch)
