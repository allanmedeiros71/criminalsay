# CriminalSay

A command-line tool that prints a randomly selected *Criminal Minds* quote to your terminal with styled formatting — red sidebar, highlighted text, and muted attribution. One command, one quote, no interactivity. Think `fortune`, but for the BAU.

```
  ▌ "It's alchemy. Alchemy turns common metals into
  ▌ precious ones. Dreams work the same way. Turning
  ▌ something awful into something better."

       David Rossi · Criminal Minds · S9E13
```

## Features

- **Zero runtime dependencies** — quotes are embedded in the binary at build time (`go:embed`)
- **Offline-first** — no network calls, no external files required
- **Terminal-aware styling** — powered by [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss); adapts to truecolor, 256-color, and ANSI-16 terminals
- **Graceful degradation** — respects `NO_COLOR`, non-TTY output, and pipes; plain text remains readable
- **Cross-platform** — builds for Linux, macOS (amd64 + arm64), and Windows
- **Shell-friendly** — fast cold start (<50ms), ideal for `.bashrc` / `.zshrc`

## Installation

### Go install

Requires [Go 1.22+](https://go.dev/dl/).

```bash
go install github.com/allanmedeiros71/criminalsay@latest
```

## Usage

```bash
criminalsay
```

That's it. The program picks a random quote, prints it to stdout, and exits with code `0`.

### Shell startup hook

Add to your `~/.bashrc` or `~/.zshrc` to greet yourself with a BAU quote on every new terminal:

```bash
criminalsay
```

### Plain output

Color is automatically disabled when output is piped or when `NO_COLOR` is set:

```bash
NO_COLOR=1 criminalsay
criminalsay | cat
```

## Building from source

Clone the repository and build:

```bash
git clone https://github.com/allanmedeiros71/criminalsay.git
cd criminalsay
go build -o criminalsay .
```

Run without installing:

```bash
go run .
```

### Cross-compilation

Build binaries for all supported platforms:

```bash
make all
```

Artifacts are written to `dist/`:

| Platform        | Binary                        |
|-----------------|-------------------------------|
| Linux amd64     | `dist/criminalsay-linux-amd64` |
| macOS amd64     | `dist/criminalsay-darwin-amd64` |
| macOS arm64     | `dist/criminalsay-darwin-arm64` |
| Windows amd64   | `dist/criminalsay-windows-amd64.exe` |

## Project structure

```
criminalsay/
├── main.go                 # Entry point: pick and render a quote
├── data/
│   └── quotes.json         # Quote database (embedded at build time)
├── internal/
│   ├── quotes/             # JSON loading, Quote type, random selection
│   └── render/             # lipgloss styling and word-wrap
└── Makefile                # Build targets (build, test, clean, cross-compile)
```

## Quote data format

Each entry in `data/quotes.json` follows this schema:

| Field          | Type   | Description                                      |
|----------------|--------|--------------------------------------------------|
| `quote`        | string | The quote text                                   |
| `author`       | string | Original author of the quote                     |
| `character`    | string | CM character who said it (optional)              |
| `season`       | int    | Season number                                    |
| `episode`      | int    | Episode number                                   |
| `episodeTitle` | string | Episode title (optional)                         |

To add a new quote, edit `data/quotes.json` and rebuild — no code changes required.

## Development

```bash
# Run tests
make test        # or: go test ./...

# Build and run
make build && ./criminalsay
```

---

Documentação em português: [README.pt.md](README.pt.md)

## License

MIT — see [LICENSE](LICENSE).

## Author

[Allan Medeiros](https://github.com/allanmedeiros71)
