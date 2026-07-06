# Phase 3: Build and Distribution - Pattern Map

**Mapped:** 2026-07-06
**Files analyzed:** 4 new/modified files (Makefile, .gitignore, README.pt.md, README.md)
**Analogs found:** 3 / 4 (Makefile has no in-repo analog — sourced from README.md aspirational content, STACK.md, and phase plan build invocations)

---

## File Classification

| New/Modified File | Role | Data Flow | Closest Analog | Match Quality |
|-------------------|------|-----------|----------------|---------------|
| `Makefile` | build orchestration | compile (source → binaries) | `README.md` cross-compile section + `.planning/research/STACK.md` | role-match |
| `.gitignore` | VCS hygiene | filter (exclude artifacts) | `CONTEXT.md` §6 proposed structure | role-match |
| `README.pt.md` | user documentation (PT-BR) | publish (instructions → user) | `README.md` (subset) | exact |
| `README.md` | user documentation (EN) | publish (instructions → user) | `README.md` (self — sync) | exact |

---

## Pattern Assignments

### `Makefile` (build orchestration, compile)

**Analog:** `README.md` cross-compile section (lines 77–92) + `.planning/research/STACK.md` (lines 174–178) + Phase 1 build gate (`01-01-PLAN.md`)

**Why this analog:** No Makefile exists yet. The EN README already documents the target artifact names and `make all` invocation; STACK.md documents the verified `GOOS`/`GOARCH` invocations. Phase 1 established `go build -o criminalsay .` as the native build contract (BUILD-01). The Makefile wraps these proven commands — it does not invent new build semantics.

**Module identity pattern** (go.mod lines 1–5):
```go
module github.com/allanmedeiros71/criminalsay

go 1.22

toolchain go1.25.6
```

**Apply:** `MAIN_PKG := .` — build from repo root where `main.go` lives. No `cmd/` subdirectory. Binary name `criminalsay` matches module short name and BUILD-01.

**Native build pattern** (README.md lines 65–69, Phase 1 gate):
```bash
git clone https://github.com/allanmedeiros71/criminalsay.git
cd criminalsay
go build -o criminalsay .
```

**Apply to `build` target (D-09):**
```makefile
build:
	go build -o $(BINARY) $(MAIN_PKG)
```

Where `BINARY := criminalsay`, `MAIN_PKG := .`.

**Embed anchor constraint** (main.go lines 3–4, 14–15 — Makefile must not break this):
```go
import (
	_ "embed"
	// ...
)

//go:embed data/quotes.json
var quoteData []byte
```

**Apply:** `go build` targets use `.` as package path. No build tags, no `-tags`, no source changes. Pure `go build` only (D-10).

**Cross-compile naming pattern** (README.md lines 85–92 — becomes Makefile output contract):
```markdown
| Platform        | Binary                        |
|-----------------|-------------------------------|
| Linux amd64     | `dist/criminalsay-linux-amd64` |
| macOS amd64     | `dist/criminalsay-darwin-amd64` |
| macOS arm64     | `dist/criminalsay-darwin-arm64` |
| Windows amd64   | `dist/criminalsay-windows-amd64.exe` |
```

**Apply:** Lock `-o` paths to these exact names (D-13/D-14). Windows **must** include `.exe` suffix.

**Cross-compile invocation pattern** (STACK.md lines 174–178, CONTEXT.md §specifics):
```bash
GOOS=linux   GOARCH=amd64  go build -o dist/criminalsay-linux-amd64 .
GOOS=darwin  GOARCH=amd64  go build -o dist/criminalsay-darwin-amd64 .
GOOS=darwin  GOARCH=arm64  go build -o dist/criminalsay-darwin-arm64 .
GOOS=windows GOARCH=amd64  go build -o dist/criminalsay-windows-amd64.exe .
```

**Apply with CGO guard** (D-03, Pitfall 8 — host default is `CGO_ENABLED=1`):
```makefile
export CGO_ENABLED=0

BINARY   := criminalsay
DIST     := dist
MAIN_PKG := .

.PHONY: all build test clean linux darwin-amd64 darwin-arm64 windows

all: linux darwin-amd64 darwin-arm64 windows

build:
	go build -o $(BINARY) $(MAIN_PKG)

test:
	go test ./...

clean:
	rm -rf $(DIST) $(BINARY)

linux: $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY)-linux-amd64 $(MAIN_PKG)

darwin-amd64: $(DIST)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST)/$(BINARY)-darwin-amd64 $(MAIN_PKG)

darwin-arm64: $(DIST)
	GOOS=darwin GOARCH=arm64 go build -o $(DIST)/$(BINARY)-darwin-arm64 $(MAIN_PKG)

windows: $(DIST)
	GOOS=windows GOARCH=amd64 go build -o $(DIST)/$(BINARY)-windows-amd64.exe $(MAIN_PKG)

$(DIST):
	mkdir -p $(DIST)
```

**Test gate pattern** (README.md line 128, Phases 1–2 convention):
```bash
go test ./...
```

**Apply to `test` target (D-09/D-11):**
```makefile
test:
	go test ./...
```

`all` does **not** depend on `test` — cross-compile and test are independent (D-11).

**Target inventory (locked D-01–D-12):**

| Target | Recipe | Req IDs |
|--------|--------|---------|
| `all` | Depends on four platform targets | BUILD-02 |
| `build` | `go build -o criminalsay .` | BUILD-01 |
| `test` | `go test ./...` | BUILD-03 |
| `clean` | `rm -rf dist criminalsay` | D-12 |
| `linux` | `GOOS=linux GOARCH=amd64` → `dist/criminalsay-linux-amd64` | BUILD-02 |
| `darwin-amd64` | `GOOS=darwin GOARCH=amd64` → `dist/criminalsay-darwin-amd64` | BUILD-02 |
| `darwin-arm64` | `GOOS=darwin GOARCH=arm64` → `dist/criminalsay-darwin-arm64` | BUILD-02 |
| `windows` | `GOOS=windows GOARCH=amd64` → `dist/criminalsay-windows-amd64.exe` | BUILD-02 |

**Anti-patterns (from 03-RESEARCH.md + CONTEXT.md):**
- Do NOT create `scripts/build.sh` (D-01 — rejected).
- Do NOT add `-ldflags` for version injection (D-10 — deferred to v2).
- Do NOT run `test` inside `all` (D-11).
- Do NOT use `curl`, `wget`, or remote includes in Makefile (ASVS V10).
- Do NOT omit `export CGO_ENABLED=0` (Pitfall 8 — breaks cross-compile on linux/amd64 host).
- Do NOT use `bin/` subdirectory for local binary (D-12 locks clean scope to root `criminalsay`).

---

### `.gitignore` (VCS hygiene, filter)

**Analog:** `CONTEXT.md` §6 proposed structure (line 122)

**Why this analog:** No `.gitignore` exists in the repo. The original project spec proposed ignoring `/dist` and binaries. Phase 3 artifacts (`dist/` from cross-compile, `criminalsay` from `make build`) must never be committed (D-15/D-16).

**Proposed structure reference** (CONTEXT.md line 122):
```
├── .gitignore              # /dist, binários
```

**Apply (D-15):**
```gitignore
dist/
criminalsay
```

**Rationale:**
- `dist/` — four cross-compiled binaries from `make all` (D-13).
- `criminalsay` — local native binary from `make build` or `go build -o criminalsay .` (D-12 clean scope).

**Optional (planner discretion — NOT required):**
- `*.exe` — only if Windows local builds pollute root; cross-compile writes `.exe` to `dist/` per D-14.
- Editor/OS cruft — out of phase scope.

**Verification pattern** (post-implementation):
```bash
make all
git status --porcelain dist/ criminalsay
# Expected: no tracked files (both ignored)
```

**Anti-patterns:**
- Do NOT commit `dist/` artifacts (D-16).
- Do NOT add overly broad patterns that ignore source files.

---

### `README.pt.md` (user documentation PT-BR, publish)

**Analog:** `README.md` (subset — Installation, Usage, Building sections only)

**Why this analog:** D-05/D-07 mandate a separate minimal PT-BR doc mirroring corrected EN install/usage/build content. PROJECT.md and CONTEXT.md lock user docs in Portuguese; code/comments stay English. README.pt.md is DOC-01 — not a full translation of the EN README.

**EN tagline pattern** (README.md lines 1–3):
```markdown
# CriminalSay

A command-line tool that prints a randomly selected *Criminal Minds* quote to your terminal with styled formatting — red sidebar, highlighted text, and muted attribution. One command, one quote, no interactivity. Think `fortune`, but for the BAU.
```

**Apply (PT-BR equivalent, 1–2 sentences):**
```markdown
# CriminalSay

Ferramenta de linha de comando que exibe uma citação aleatória de *Criminal Minds* no terminal, com formatação estilizada — barra lateral vermelha, texto destacado e atribuição discreta. Um comando, uma citação, sem interatividade.
```

**Installation — go install pattern** (README.md lines 24–30, corrected version):
```markdown
### Go install

Requires [Go 1.22+](https://go.dev/dl/).

```bash
go install github.com/allanmedeiros71/criminalsay@latest
```
```

**Apply (PT-BR, D-08 — keep install path, fix version):**
```markdown
## Instalação

### Via Go install

Requer [Go 1.22+](https://go.dev/dl/).

```bash
go install github.com/allanmedeiros71/criminalsay@latest
```
```

**Installation — build from source pattern** (README.md lines 61–69):
```markdown
## Building from source

Clone the repository and build:

```bash
git clone https://github.com/allanmedeiros71/criminalsay.git
cd criminalsay
go build -o criminalsay .
```
```

**Apply (PT-BR, add `make` alternative per D-08):**
```markdown
### Compilar do código-fonte

```bash
git clone https://github.com/allanmedeiros71/criminalsay.git
cd criminalsay
make
```

Ou diretamente:

```bash
go build -o criminalsay .
```
```

**Usage pattern** (README.md lines 36–59):
```markdown
## Usage

```bash
criminalsay
```

### Shell startup hook

Add to your `~/.bashrc` or `~/.zshrc` ...

### Plain output

```bash
NO_COLOR=1 criminalsay
criminalsay | cat
```
```

**Apply (PT-BR):**
```markdown
## Uso

```bash
criminalsay
```

### Hook no shell

Adicione ao seu `~/.bashrc` ou `~/.zshrc` para exibir uma citação ao abrir o terminal:

```bash
criminalsay
```

### Saída sem cor

```bash
NO_COLOR=1 criminalsay
criminalsay | cat
```
```

**Cross-compile / build section pattern** (README.md lines 77–92):
```markdown
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
```

**Apply (PT-BR, include convenience targets per D-09):**
```markdown
## Build

Para compilar binários de todas as plataformas suportadas:

```bash
make all
```

Os artefatos são gravados em `dist/`:

| Plataforma      | Binário                        |
|-----------------|--------------------------------|
| Linux amd64     | `dist/criminalsay-linux-amd64` |
| macOS amd64     | `dist/criminalsay-darwin-amd64` |
| macOS arm64     | `dist/criminalsay-darwin-arm64` |
| Windows amd64   | `dist/criminalsay-windows-amd64.exe` |

Outros alvos úteis:

```bash
make build   # binário local na raiz do repositório
make test    # executa go test ./...
make clean   # remove dist/ e o binário local
```
```

**Omit from README.pt.md (D-07):**
- Features list (README.md lines 13–20)
- Project structure tree (README.md lines 94–107)
- Quote data format / JSON schema (README.md lines 109–122)
- Development section (README.md lines 124–132)
- License / Author sections (optional one-liner OK)

**Omit per D-08:**
- Pre-built binaries / GitHub Releases section (README.md lines 32–34) — no release pipeline in v1.

**Go version source of truth** (go.mod line 3):
```
go 1.22
```

**Apply:** Both READMEs must state **Go 1.22+**, not 1.21+.

---

### `README.md` (user documentation EN, publish — modify)

**Analog:** `README.md` (self — sync aspirational content to match real Makefile)

**Why this analog:** The EN README was written speculatively before Phase 3. Cross-compile section (lines 77–92) is correct and will become accurate once Makefile lands. Other sections contain inaccuracies enumerated in 03-RESEARCH.md.

**Inaccuracy 1 — Go version** (README.md line 26):
```markdown
Requires [Go 1.21+](https://go.dev/dl/).
```

**Fix (D-06, align with go.mod):**
```markdown
Requires [Go 1.22+](https://go.dev/dl/).
```

**Inaccuracy 2 — Pre-built binaries** (README.md lines 32–34):
```markdown
### Pre-built binaries

Download the binary for your platform from the [Releases](https://github.com/allanmedeiros71/criminalsay/releases) page and place it on your `PATH`.
```

**Fix (D-08):** **Remove entire section** — no release pipeline in v1.

**Inaccuracy 3 — Project structure** (README.md lines 104–106):
```
├── scripts/
│   └── build.sh            # Cross-compile helper
└── Makefile                # Build targets
```

**Fix (D-01):** Remove `scripts/build.sh` line. Keep `Makefile` only:
```
└── Makefile                # Build targets (build, test, clean, cross-compile)
```

**Keep unchanged — becomes accurate after Makefile** (README.md lines 77–92):
```markdown
### Cross-compilation

Build binaries for all supported platforms:

```bash
make all
```

Artifacts are written to `dist/`:
...
```

**Keep unchanged — valid install method** (README.md lines 28–30):
```bash
go install github.com/allanmedeiros71/criminalsay@latest
```

**Development section sync** (README.md lines 124–132):
```markdown
## Development

```bash
# Run tests
go test ./...

# Build and run
go build -o criminalsay . && ./criminalsay
```
```

**Apply:** Add `make test` / `make build` as alternatives (optional, low cost):
```markdown
## Development

```bash
# Run tests
make test        # or: go test ./...

# Build and run
make build && ./criminalsay
```
```

**Optional footer link** (D-05 discretion):
```markdown
---

Documentação em português: [README.pt.md](README.pt.md)
```

**Sync checklist (D-06):**

| # | Action | Source |
|---|--------|--------|
| 1 | `Go 1.21+` → `Go 1.22+` | go.mod line 3 |
| 2 | Remove Pre-built binaries / Releases section | D-08 |
| 3 | Remove `scripts/build.sh` from project structure | D-01 |
| 4 | Keep `make all` / `dist/` table (now accurate) | D-13 |
| 5 | Optional: link to README.pt.md | D-05 |
| 6 | Optional: mention `make test` / `make build` in Development | D-09 |

---

## Shared Patterns

### Go Version Floor (go.mod → READMEs)

**Source:** `go.mod` line 3 + `math/rand/v2` requirement

```go
go 1.22
```

**Apply to:** `README.md`, `README.pt.md` — both must state **Go 1.22+**.

### Module Install Path

**Source:** `go.mod` line 1 + README.md line 29

```bash
go install github.com/allanmedeiros71/criminalsay@latest
```

**Apply to:** Both READMEs in Installation section (D-08).

### Test Gate Convention (Phases 1–2 — unchanged)

**Source:** `README.md` line 128, `.planning/ROADMAP.md`

```bash
go test ./...
```

**Apply to:** Makefile `test` target wraps this exactly. No new test framework.

### Binary Naming Convention

**Source:** BUILD-01 (Phase 1) + README.md dist/ table

| Context | Name | Location |
|---------|------|----------|
| Native build | `criminalsay` | repo root |
| Linux cross | `criminalsay-linux-amd64` | `dist/` |
| macOS amd64 | `criminalsay-darwin-amd64` | `dist/` |
| macOS arm64 | `criminalsay-darwin-arm64` | `dist/` |
| Windows amd64 | `criminalsay-windows-amd64.exe` | `dist/` |

### Documentation Language Split (PROJECT.md + CONTEXT.md)

**Source:** CONTEXT.md line 18

| Artifact | Language | Scope |
|----------|----------|-------|
| Go source, comments | English | unchanged in Phase 3 |
| `README.md` | English | full OSS doc (sync fixes) |
| `README.pt.md` | Portuguese (PT-BR) | minimal DOC-01 only |

### Build-Output Git Hygiene

**Source:** D-15/D-16 + CONTEXT.md §6

```
.gitignore: dist/, criminalsay
make clean: rm -rf dist criminalsay
```

Never commit build artifacts.

### Phase 3 Non-Touch List

**Source:** 03-CONTEXT.md code_context

| File/Package | Action |
|--------------|--------|
| `main.go` | Do NOT modify (unless build-tag issue surfaces) |
| `internal/quotes/` | Do NOT modify |
| `internal/render/` | Do NOT modify |
| `go.mod` / `go.sum` | Do NOT modify (no new deps) |

---

## No Analog Found

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| `Makefile` | build orchestration | compile | No Makefile in repo — patterns from README.md aspirational content, STACK.md cross-compile invocations, and Phase 1 `go build` gate |
| `.gitignore` | VCS hygiene | filter | No `.gitignore` in repo — pattern from CONTEXT.md §6 proposal and D-15/D-16 locked decisions |

`README.pt.md` and `README.md` have exact analogs (README.md self / subset). Makefile and `.gitignore` patterns are fully specified in 03-RESEARCH.md Architecture Patterns section — planner must follow those verified structures.

---

## Metadata

**Analog search scope:** `/home/allan/git/criminalsay` — `README.md`, `go.mod`, `main.go`, `.planning/research/STACK.md`, `CONTEXT.md`, Phase 1 build gates
**Files scanned:** 4 source/config files + 2 planning research files
**Pattern extraction date:** 2026-07-06
**Pattern sources:** README.md aspirational build docs (primary for artifact naming), go.mod (version floor), main.go (embed/build constraints), STACK.md (cross-compile invocations), 03-CONTEXT.md + 03-RESEARCH.md (locked decisions D-01–D-16)
