# Phase 3: Build and Distribution — Research

**Researched:** 2026-07-06
**Domain:** Go cross-compilation (Makefile, CGO_ENABLED=0, dist/ artifacts), user documentation (README sync, README.pt.md)
**Confidence:** HIGH (cross-compile verified on all four targets; Makefile pattern is stdlib-only; doc gaps enumerated against live README.md)

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Build Tooling

- **D-01:** **Makefile only** — no `scripts/build.sh`. The Makefile is the sole build interface; aligns with ROADMAP and research recommendations.
- **D-02:** **Explicit per-platform targets** plus aggregate: `linux`, `darwin-amd64`, `darwin-arm64`, `windows`, and `all` (builds all four). Not a single opaque `all`-only Makefile.
- **D-03:** **CGO_ENABLED=0** via `export CGO_ENABLED=0` at the top of the Makefile (Claude discretion — DRY, one declaration covers all targets). Pitfall 8 from research.
- **D-04:** Makefile **auto-creates `dist/`** with `mkdir -p dist` before writing binaries.

### README Strategy

- **D-05:** **Separate PT-BR doc** — keep `README.md` in English; add `README.pt.md` for DOC-01. Do not replace EN with PT or create a bilingual single file.
- **D-06:** **Sync both READMEs** in this phase — update EN to fix inaccuracies (Go 1.22+ not 1.21+, reference real Makefile targets) and add `README.pt.md` with equivalent corrected content.
- **D-07:** **README.pt.md is minimal DOC-01** — installation, usage, and build instructions only. Omit project structure, JSON schema, development section, and other extras present in the EN README.
- **D-08:** **Installation methods documented:** `go install` and build from source (`git clone` + `make` / `go build`). **No GitHub Releases** mention — v1 has no release pipeline.

### Makefile Targets

- **D-09:** **Convenience targets** (Claude discretion): `build` (local `go build -o criminalsay .`), `test` (`go test ./...`), `clean`, plus cross-compile targets and `all`.
- **D-10:** **No version ldflags** in v1 (Claude discretion) — plain `go build` without `-ldflags` injection.
- **D-11:** **`make test` is independent** of `all` (Claude discretion) — cross-compile does not gate on tests; user runs `make test` separately.
- **D-12:** **`make clean`** removes `dist/` and the local `criminalsay` binary at the repository root.

### dist/ Convention

- **D-13:** **Flat prefixed naming** in `dist/` (Claude discretion, consistent with existing README EN):
  - `dist/criminalsay-linux-amd64`
  - `dist/criminalsay-darwin-amd64`
  - `dist/criminalsay-darwin-arm64`
  - `dist/criminalsay-windows-amd64.exe`
- **D-14:** Windows binary **must include `.exe` suffix**.
- **D-15:** **Add `.gitignore`** with `dist/` and `criminalsay` (local root binary). Repo currently has no `.gitignore`.
- **D-16:** **Never commit `dist/` artifacts** — build outputs are local/CI-only.

### Claude's Discretion

- Exact Makefile variable names (`GOOS`, `GOARCH`, `OUTPUT`, `LDFLAGS` if added later).
- Whether `build` target uses `-o criminalsay` at root or a `bin/` subdirectory (default: root per D-12 clean scope).
- README.pt.md tone and section headings (as long as DOC-01 sections are covered).
- Link from README.md to README.pt.md (optional footer note) — not required but acceptable.

### Deferred Ideas (OUT OF SCOPE)

- `scripts/build.sh` — rejected in favor of Makefile-only (D-01)
- GitHub Releases / pre-built binary downloads — no release pipeline in v1 (D-08)
- GoReleaser, Homebrew, AUR — v2 per PROJECT.md and REQUIREMENTS.md
- Version injection via `-ldflags` — deferred (D-10)
- CI/GitHub Actions for automated cross-compile — out of Phase 3 scope
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| BUILD-02 | Cross-compile script (`Makefile` or `scripts/build.sh`) produces binaries for Linux amd64, macOS amd64, macOS arm64, Windows amd64; all with CGO_ENABLED=0 | Makefile-only per D-01; `export CGO_ENABLED=0` per D-03; four explicit targets + `all` per D-02/D-13; cross-compile verified on dev machine [VERIFIED: 2026-07-06] |
| DOC-01 | README in PT-BR with installation instructions, usage, and build instructions | Separate `README.pt.md` per D-05/D-07; minimal three-section doc; EN `README.md` synced per D-06 (Go 1.22+, remove aspirational Releases/scripts references) |
</phase_requirements>

---

## Summary

Phase 3 is mechanical infrastructure: add a Makefile, `.gitignore`, and Portuguese user documentation. No Go code changes are expected — `main.go`, `internal/quotes`, and `internal/render` are complete from Phases 1–2. The binary already builds with `go build -o criminalsay .` (BUILD-01 satisfied). Cross-compilation to all four targets succeeds with `CGO_ENABLED=0` and plain `go build` invocations [VERIFIED: all four targets built and `file`-typed on linux/amd64 host].

The EN `README.md` is **aspirational**: it documents `make all`, `dist/` naming, and a `Makefile` that do not exist yet; it also contains inaccuracies to fix during sync (Go 1.21+ vs 1.22+, `scripts/build.sh` reference, GitHub Releases download section). Phase 3 makes the README truthful.

**Primary recommendation:** Implement a ~40-line Makefile with `export CGO_ENABLED=0`, per-platform targets using `GOOS`/`GOARCH` prefix syntax, `mkdir -p dist` in each cross-compile recipe (or a shared `$(DIST_DIR)` prerequisite), convenience targets `build`/`test`/`clean`, and aggregate `all`. Add `.gitignore` with `dist/` and `criminalsay`. Create minimal `README.pt.md` (Instalação, Uso, Build). Sync `README.md` fixes per D-06/D-08.

---

## Current State Audit

| Artifact | Exists? | Notes |
|----------|---------|-------|
| `Makefile` | ❌ | Referenced in README.md but not in repo |
| `.gitignore` | ❌ | No ignore rules; `dist/` and root binary would be committable |
| `README.pt.md` | ❌ | DOC-01 deliverable |
| `scripts/build.sh` | ❌ | Correctly absent — rejected per D-01; README incorrectly lists it |
| `go.mod` | ✓ | `go 1.22`, `toolchain go1.25.6` |
| `README.md` | ✓ | Needs sync (see below) |
| Cross-compile (manual) | ✓ | All four targets build with `CGO_ENABLED=0` |

### README.md Inaccuracies to Fix (D-06)

| Location | Current | Correct |
|----------|---------|---------|
| Installation → Go version | "Requires Go 1.21+" | "Requires Go 1.22+" (math/rand/v2 floor in go.mod) |
| Pre-built binaries section | Links to GitHub Releases | **Remove** — no release pipeline in v1 (D-08) |
| Project structure | Lists `scripts/build.sh` | Remove; Makefile only (D-01) |
| Cross-compilation | Documents `make all` | Keep — will become accurate once Makefile lands |
| `go install` path | `github.com/allanmedeiros71/criminalsay@latest` | Keep — valid install method per D-08 |

### go.mod vs README Version

```
go.mod:     go 1.22
README.md:  Go 1.21+   ← WRONG
```

`math/rand/v2` requires Go 1.22+. Both READMEs must state **1.22+**.

---

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Local native build | Makefile `build` target | `go build` direct | D-09 convenience wrapper |
| Cross-compile orchestration | Makefile per-platform targets | Go `GOOS`/`GOARCH` env | Built-in Go cross-compile; no external tools |
| CGO disable | Makefile `export CGO_ENABLED=0` | Per-recipe prefix (redundant) | D-03; Pitfall 8 — default is `1` on linux/amd64 [VERIFIED: `go env CGO_ENABLED` → `1`] |
| Artifact output | `dist/` directory | — | D-04/D-13 flat naming |
| VCS hygiene | `.gitignore` | — | D-15/D-16 prevent committing binaries |
| EN user docs | `README.md` | — | Project default language for OSS audience |
| PT-BR user docs | `README.pt.md` | — | DOC-01; D-05 separate file |
| Test gate | Makefile `test` target | `go test ./...` | D-11 independent of `all`; BUILD-03 already green |

---

## Standard Stack

### Core (no new packages)

| Technology | Version | Purpose | Why Standard |
|------------|---------|---------|--------------|
| Go toolchain | 1.22+ (toolchain 1.25.6 installed) | `go build`, cross-compile via `GOOS`/`GOARCH` | Zero external build deps; official cross-compile since Go 1.5+ |
| GNU Make | system default | Build orchestration | Ubiquitous on Linux/macOS dev machines; Windows users use WSL or `go build` directly |
| Makefile + `go build` | — | Sole build interface (D-01) | Simpler than GoReleaser for v1; deferred to v2 |

**Installation:** No new dependencies. Phase 3 adds only Makefile, `.gitignore`, and markdown files.

**Cross-compile verification (2026-07-06, host linux/amd64):**

```bash
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64  go build -o dist/criminalsay-linux-amd64 .
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64  go build -o dist/criminalsay-darwin-amd64 .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64  go build -o dist/criminalsay-darwin-arm64 .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -o dist/criminalsay-windows-amd64.exe .
```

Results:
- `criminalsay-linux-amd64` → ELF 64-bit LSB executable, x86-64, statically linked
- `criminalsay-darwin-amd64` → Mach-O 64-bit x86_64
- `criminalsay-darwin-arm64` → Mach-O 64-bit arm64
- `criminalsay-windows-amd64.exe` → PE32+ executable (console) x86-64

Pure Go — no `import "C"` in module [VERIFIED: ripgrep across `*.go`]. lipgloss v1.1.0 and transitive deps are pure Go; `CGO_ENABLED=0` is defensive per Pitfall 8.

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Makefile | `scripts/build.sh` | Rejected per D-01; shell script duplicates Make semantics |
| Makefile | GoReleaser | Deferred to v2; YAML config + CI overhead not justified |
| `export CGO_ENABLED=0` | Per-target `CGO_ENABLED=0` prefix | Both work; export is DRY per D-03 |
| `README.pt.md` separate | Bilingual single README | Rejected per D-05 |
| `bin/` subdirectory | Root `criminalsay` binary | D-12 locks clean scope to root binary |

---

## Architecture Patterns

### Recommended Makefile Structure

```makefile
# Source: D-01 through D-16, STACK.md cross-compile section, PITFALLS.md Pitfall 8
export CGO_ENABLED=0

BINARY     := criminalsay
DIST       := dist
MAIN_PKG   := .

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

**Key behaviors:**
- `export CGO_ENABLED=0` at file top covers `build`, all cross targets, and any future targets (D-03).
- `GOOS`/`GOARCH` as recipe-prefix env vars — standard Go idiom; overrides host defaults per target.
- `$(DIST)` prerequisite with `mkdir -p` satisfies D-04; each platform target depends on it.
- `all` lists four platform targets explicitly (D-02); does **not** run `test` (D-11).
- `clean` removes both `dist/` and root `criminalsay` (D-12).
- No `-ldflags` (D-10).

### Pattern: .gitignore

```gitignore
# Source: D-15
dist/
criminalsay
```

Optional additions (planner discretion, not required):
- `*.exe` — only if Windows local builds pollute root; cross-compile writes to `dist/` with `.exe` suffix per D-14.
- Editor/OS cruft (`*.swp`, `.DS_Store`) — out of phase scope unless user requests.

### Pattern: README.pt.md Minimal Structure (D-07)

```markdown
# CriminalSay

[1–2 sentence PT-BR description — mirror EN tagline]

## Instalação

### Via Go install
Requer Go 1.22+.
go install github.com/allanmedeiros71/criminalsay@latest

### Compilar do código-fonte
git clone … && cd criminalsay && make

## Uso

criminalsay
~/.bashrc / ~/.zshrc hook
NO_COLOR=1 criminalsay

## Build

git clone + make all
Tabela dist/ (4 plataformas, nomes D-13)
make build / make test / make clean
```

**Omit from README.pt.md:** Features list, project structure, JSON schema, Development section, License/Author (optional one-liner OK).

### Pattern: README.md Sync Checklist

1. Change "Go 1.21+" → "Go 1.22+"
2. Remove "### Pre-built binaries" section and Releases link (D-08)
3. Remove `scripts/build.sh` from project structure tree (D-01)
4. Keep `make all` / `dist/` table — becomes accurate
5. Optional: footer link "Documentação em português: README.pt.md" (D-05 discretion)
6. Verify `make test` / `go test ./...` in Development section still accurate

### Anti-Patterns to Avoid

- **Adding `scripts/build.sh`:** Explicitly rejected (D-01); README currently references it — remove, don't create.
- **Committing `dist/` artifacts:** Violates D-16; `.gitignore` is mandatory.
- **Forgetting `.exe` on Windows output:** D-14; `go build -o dist/criminalsay-windows-amd64` without suffix still runs on Windows but breaks naming convention and `file` expectations.
- **Running tests inside `all` target:** Violates D-11; keeps cross-compile fast and independent.
- **Mentioning GitHub Releases in PT doc:** D-08 forbids for v1.
- **Modifying Go source for build tags:** NFR-2 — no platform build tags; not needed (pure Go).

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Cross-compile orchestration | Custom shell matrix script | Makefile + `GOOS`/`GOARCH` | Go toolchain handles cross-compile natively |
| Version stamping | Manual string injection | Defer to v2 (D-10) | `-ldflags` adds complexity without v1 release pipeline |
| Release publishing | GitHub Actions / GoReleaser | Local `make all` only | Out of Phase 3 scope |
| Platform detection at build time | `runtime.GOOS` conditionals in Makefile | Explicit per-target recipes | D-02 requires named targets, not magic |

---

## Common Pitfalls

### Pitfall 1: CGO_ENABLED defaults to 1 on native linux/amd64

**What goes wrong:** Cross-compile to Windows/macOS fails with C linker errors if a transitive dep enables cgo.

**Why it happens:** `go env CGO_ENABLED` returns `1` on linux/amd64 [VERIFIED].

**How to avoid:** `export CGO_ENABLED=0` at Makefile top (D-03, Pitfall 8).

**Detection:** `make all` fails with `gcc: executable file not found` for non-native targets.

### Pitfall 2: README/doc drift from reality

**What goes wrong:** Users follow README, `make all` fails because Makefile missing; users expect Releases downloads.

**Why it happens:** README was written speculatively before Phase 3 implementation.

**How to avoid:** D-06 sync pass in same phase as Makefile creation.

**Detection:** Fresh clone + follow README instructions fails.

### Pitfall 3: Windows binary without .exe extension

**What goes wrong:** Naming inconsistency; users on Windows may not recognize executable; ROADMAP success criteria table expects `.exe`.

**How to avoid:** D-14 explicit suffix in `-o` path.

### Pitfall 4: `make clean` misses root binary

**What goes wrong:** `criminalsay` at repo root gets committed or confuses `git status`.

**How to avoid:** D-12 — `rm -rf dist criminalsay` in clean target; `.gitignore` entry for `criminalsay`.

### Pitfall 5: Accidentally committing dist/

**What goes wrong:** Large binaries in git history.

**How to avoid:** D-15 `.gitignore` with `dist/`; verify with `git status` after `make all`.

### Pitfall 6: Windows smoke test gap (carry-over from STATE.md)

**What goes wrong:** lipgloss v1 non-TTY fallback on cmd.exe unverified.

**Why it matters:** Phase 3 ROADMAP success criterion #2 requires binaries run on target platform.

**Mitigation:** Manual smoke test on Windows if available; document as UAT item. Linux-hosted cross-compile produces valid PE binary [VERIFIED: `file`]; runtime behavior is Phase 2 concern.

---

## Code Examples

### Verification Script (post-implementation)

```bash
# Phase gate — all must exit 0
make clean
make all
make test

# Artifact inventory
test $(ls dist/ | wc -l) -eq 4
file dist/criminalsay-linux-amd64   | grep -q 'ELF.*x86-64'
file dist/criminalsay-darwin-amd64    | grep -q 'Mach-O.*x86_64'
file dist/criminalsay-darwin-arm64    | grep -q 'Mach-O.*arm64'
file dist/criminalsay-windows-amd64.exe | grep -q 'PE32+'

# Native smoke
./dist/criminalsay-linux-amd64 | head -1

# Git hygiene
git status --porcelain dist/ criminalsay
# Expected: no tracked files (both ignored)
```

### go install Path (document in both READMEs)

```bash
go install github.com/allanmedeiros71/criminalsay@latest
```

Requires Go 1.22+ on the installing machine. Installs to `$GOPATH/bin` or `$HOME/go/bin` depending on `GOBIN`/`GOPATH` settings — standard Go behavior; no phase-specific handling needed.

---

## State of the Art

| Old (README aspirational) | After Phase 3 | Notes |
|---------------------------|---------------|-------|
| Makefile referenced, missing | Makefile present | D-01 |
| `scripts/build.sh` in tree diagram | Removed from docs | Never existed; D-01 |
| Go 1.21+ | Go 1.22+ | Aligns with go.mod |
| GitHub Releases downloads | Removed | D-08 |
| EN-only docs | EN + README.pt.md | D-05 |

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Pure Go deps — CGO_ENABLED=0 always sufficient | Standard Stack | If future dep adds cgo, cross-compile breaks; mitigated by Pitfall 8 export |
| A2 | `make` available on target dev machines | Patterns | Windows-native devs without Make can use documented `go build` commands |
| A3 | No Go source changes needed | Summary | Only if build-tag issue surfaces (unlikely) |
| A4 | README.pt.md minimal scope satisfies DOC-01 | Patterns | REQUIREMENTS only mandate install/usage/build — confirmed |

---

## Open Questions

1. **README.md link to README.pt.md**
   - What we know: Optional per D-05 discretion.
   - Recommendation: Add one-line footer in EN README — low cost, helps PT-BR users discover doc.

2. **Windows manual UAT**
   - What we know: STATE.md flags cmd.exe lipgloss fallback unverified.
   - Recommendation: Add to phase VERIFICATION.md as manual check; not a Makefile blocker.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go toolchain | build, test, cross-compile | ✓ | go1.25.6 linux/amd64 | — |
| `make` | Makefile invocation | ✓ (assumed) | GNU Make | Direct `go build` commands in README |
| `file` command | artifact type verification | ✓ | system | `go version -m` on binary |

**Missing dependencies with no fallback:** None for implementation.

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` (existing) + Makefile targets (new) |
| Config file | none |
| Quick run command | `make test` |
| Full suite command | `make test` (= `go test ./...`) |
| Build gate command | `make all` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| BUILD-02 | `make all` produces 4 binaries in `dist/` | integration | `make clean && make all && test $(ls dist/ \| wc -l) -eq 4` | ❌ Wave 0 (Makefile) |
| BUILD-02 | Linux amd64 ELF artifact | integration | `file dist/criminalsay-linux-amd64 \| grep -q 'ELF.*x86-64'` | ❌ Wave 0 |
| BUILD-02 | macOS amd64 Mach-O artifact | integration | `file dist/criminalsay-darwin-amd64 \| grep -q 'Mach-O.*x86_64'` | ❌ Wave 0 |
| BUILD-02 | macOS arm64 Mach-O artifact | integration | `file dist/criminalsay-darwin-arm64 \| grep -q 'Mach-O.*arm64'` | ❌ Wave 0 |
| BUILD-02 | Windows amd64 PE artifact with .exe | integration | `file dist/criminalsay-windows-amd64.exe \| grep -q 'PE32+'` | ❌ Wave 0 |
| BUILD-02 | CGO_ENABLED=0 for cross targets | integration | `go version -m dist/criminalsay-linux-amd64 2>/dev/null \| grep -v cgo` or build log inspection; Makefile has `export CGO_ENABLED=0` | ❌ Wave 0 |
| BUILD-02 | Native linux binary runs | smoke | `./dist/criminalsay-linux-amd64 \| head -1` | ❌ Wave 0 |
| BUILD-02 | `make build` local binary | integration | `make build && test -x ./criminalsay` | ❌ Wave 0 |
| BUILD-02 | `make clean` removes artifacts | integration | `make all && make clean && test ! -e dist && test ! -e criminalsay` | ❌ Wave 0 |
| BUILD-03 | Tests still pass (regression) | integration | `make test` or `go test ./...` | ✓ (existing tests green) |
| DOC-01 | README.pt.md exists with install section | manual/grep | `test -f README.pt.md && grep -qi 'instala' README.pt.md` | ❌ Wave 0 |
| DOC-01 | README.pt.md has usage section | manual/grep | `grep -qi 'uso\|usage' README.pt.md` | ❌ Wave 0 |
| DOC-01 | README.pt.md has build section | manual/grep | `grep -qi 'build\|compil' README.pt.md` | ❌ Wave 0 |
| DOC-01 | Go 1.22+ stated in PT doc | manual/grep | `grep -q '1.22' README.pt.md` | ❌ Wave 0 |
| DOC-01 | EN README synced (Go version) | manual/grep | `grep -q '1.22' README.md && ! grep -q '1.21' README.md` | ❌ Wave 0 |
| DOC-01 | No GitHub Releases in docs | manual/grep | `! grep -qi 'releases' README.md README.pt.md` | ❌ Wave 0 (README.md currently fails) |

### Sampling Rate

- **Per task commit:** `make test`
- **Per wave merge:** `make all && make test`
- **Phase gate:** `make clean && make all && make test && file dist/*`

### Wave 0 Gaps

- [ ] `Makefile` — all targets per D-01–D-12
- [ ] `.gitignore` — `dist/`, `criminalsay` per D-15
- [ ] `README.pt.md` — minimal DOC-01 per D-07
- [ ] `README.md` — sync fixes per D-06/D-08

---

## Security Domain

> `security_enforcement: true`, `security_asvs_level: 1` per `.planning/config.json`.

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V1 Architecture | Minimal | Build produces static binaries; no runtime services added |
| V2 Authentication | No | Unchanged — CLI tool |
| V3 Session Management | No | Stateless |
| V4 Access Control | No | No user accounts |
| V5 Input Validation | No change | Build phase does not add input surfaces |
| V6 Cryptography | No | No secrets in build |
| V9 Communication | No | No network I/O |
| V10 Malicious Code | Low | Ensure Makefile does not fetch/execute remote code |
| V14 Configuration | Minimal | `.gitignore` prevents accidental binary commit |

### Known Threat Patterns for this Stack

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Makefile curl\|bash supply chain | Tampering | Keep Makefile to local `go build` only — no `curl`, `wget`, or remote includes (D-01 scope) |
| Committed binary artifacts | Information disclosure | `.gitignore` + D-16; binaries may contain embedded data paths |
| Stale install instructions pointing to compromised Releases | Spoofing | D-08 removes Releases references until v2 pipeline exists |
| World-writable dist/ permissions | Elevation | Default umask; not a v1 concern for local dev builds |

**Security summary:** Phase 3 adds no attack surface to the running binary. ASVS L1 satisfied by keeping the Makefile free of remote code execution, ignoring build artifacts in VCS, and not documenting untrusted download channels.

---

## Project Constraints (from .cursor/rules/)

No `.cursor/rules/` directory exists in this repository. No additional Cursor rule constraints beyond user rules and `.planning/` artifacts.

---

## Sources

### Primary (HIGH confidence — verified on dev machine 2026-07-06)

- `go build` with `CGO_ENABLED=0` for linux/darwin-amd64/darwin-arm64/windows-amd64 — all succeed
- `file` output confirms correct executable formats per platform
- `go env CGO_ENABLED` → `1` (confirms need for explicit disable)
- `go.mod` → `go 1.22`
- `go test ./...` → green (quotes + render)
- ripgrep: no `import "C"` in `*.go`

### Secondary (MEDIUM confidence — project locked artifacts)

- `.planning/phases/03-build-and-distribution/03-CONTEXT.md` — D-01 through D-16
- `.planning/research/PITFALLS.md` — Pitfall 8 (CGO cross-compile)
- `.planning/research/STACK.md` — Makefile cross-compile invocations
- `.planning/research/SUMMARY.md` — Phase 3 deliverables ordering
- `.planning/REQUIREMENTS.md` — BUILD-02, DOC-01
- `README.md` — current aspirational content and inaccuracies enumerated

### Tertiary (LOW confidence — reference)

- [Go cross-compilation — go.dev doc](https://go.dev/doc/install/source#environment) — GOOS/GOARCH environment variables
- [CGO_ENABLED — golang/go Issue #5104](https://github.com/golang/go/issues/5104) — cross-compile behavior

---

## Metadata

**Confidence breakdown:**
- Makefile pattern: HIGH — standard Go idiom; verified cross-compile on host
- README sync scope: HIGH — inaccuracies enumerated line-by-line against live README.md
- CGO_ENABLED=0 necessity: HIGH — default is 1; defensive even for pure Go module
- Windows runtime smoke: MEDIUM — PE binary verified; lipgloss cmd.exe behavior unverified (STATE.md carry-over)
- DOC-01 PT content: HIGH — minimal scope locked in D-07

**Research date:** 2026-07-06
**Valid until:** 2026-08-06 (stable build surface; no new deps expected)

---

## RESEARCH COMPLETE

**Phase:** 03 - build-and-distribution
**Confidence:** HIGH

### Key Findings

- Cross-compile to all four targets works today with `CGO_ENABLED=0` and plain `go build` — no Go code changes needed.
- `export CGO_ENABLED=0` at Makefile top is mandatory (host default is `1` on linux/amd64).
- EN README is aspirational: fix Go 1.22+, remove Releases section and `scripts/build.sh` reference when adding real Makefile.
- Deliverables are Makefile + `.gitignore` + `README.pt.md` (minimal) + README.md sync — no new packages.
- `make test` must remain independent of `make all` per D-11.

### File Created

`.planning/phases/03-build-and-distribution/03-RESEARCH.md`

### Ready for Planning

Research complete. Planner can now create PLAN.md files.
