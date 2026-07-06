# Phase 3: Build and Distribution - Context

**Gathered:** 2026-07-06
**Status:** Ready for planning

<domain>
## Phase Boundary

Add cross-compile build infrastructure and user-facing documentation so the project can produce binaries for Linux amd64, macOS amd64+arm64, and Windows amd64 from a single `make` invocation, with artifacts in `dist/`, and Portuguese documentation for installation, usage, and build.

No CLI flags, no GoReleaser, no GitHub Actions CI, no Homebrew/AUR packaging, no version ldflags, no scripts/build.sh (Makefile is the sole build interface).

</domain>

<decisions>
## Implementation Decisions

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

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements and Scope
- `.planning/REQUIREMENTS.md` — Phase 3 REQ-IDs: BUILD-02, DOC-01
- `.planning/ROADMAP.md` — Phase 3 success criteria (make all, dist/, PT-BR README)

### Project Constraints
- `.planning/PROJECT.md` — CGO_ENABLED=0, platforms D8, GoReleaser/Homebrew deferred to v2
- `.planning/STATE.md` — Windows cmd.exe lipgloss fallback unverified; note for manual smoke test

### Build Research
- `.planning/research/PITFALLS.md` — Pitfall 8: CGO_ENABLED=1 breaks cross-compile
- `.planning/research/STACK.md` — Makefile target patterns, cross-compile invocations
- `.planning/research/SUMMARY.md` — Phase 3 deliverables and ordering rationale

### Original Project Spec
- `CONTEXT.md` §3 (escopo build/README), §6 (estrutura proposta com Makefile)
- `go.mod` — Go 1.22+ required (README must say 1.22+, not 1.21+)

### Existing Documentation (to sync)
- `README.md` — Current EN README; aspirational `make all` / `dist/` content exists but Makefile does not yet

### Prior Phase Boundaries
- `.planning/phases/02-render-and-output/02-CONTEXT.md` — Phase 2 explicitly deferred Makefile to Phase 3

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `go.mod` / `go.sum` — Go 1.22, toolchain 1.25.6; lipgloss v1.1.0; pure Go, no cgo deps
- `main.go` — complete binary entry point; `go build -o criminalsay .` already works
- `README.md` — EN documentation with target dist/ naming and structure (needs sync with real Makefile)

### Established Patterns
- Embed anchor in `main.go`; no build tags per platform (NFR-2)
- `go test ./...` is the project-wide test gate (Phases 1–2)
- Code/comments in English; user docs in Portuguese (README.pt.md)

### Integration Points
- **New files:** `Makefile`, `.gitignore`, `README.pt.md`
- **Modified files:** `README.md` (sync fixes, optional link to PT)
- **Not created:** `scripts/build.sh` (explicitly rejected per D-01)
- Phase 3 must not modify `internal/quotes`, `internal/render`, or `main.go` logic unless a build-tag issue surfaces

</code_context>

<specifics>
## Specific Ideas

Makefile cross-compile pattern (from research, locked by D-03/D-13):

```makefile
export CGO_ENABLED=0

all: linux darwin-amd64 darwin-arm64 windows

linux:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/criminalsay-linux-amd64 .
```

README.pt.md minimal sections:
1. Instalação (`go install`, build local)
2. Uso (`criminalsay`, `.bashrc`/`.zshrc`, `NO_COLOR`)
3. Build (`git clone`, `make all`, tabela dist/)

</specifics>

<deferred>
## Deferred Ideas

- `scripts/build.sh` — rejected in favor of Makefile-only (D-01)
- GitHub Releases / pre-built binary downloads — no release pipeline in v1 (D-08)
- GoReleaser, Homebrew, AUR — v2 per PROJECT.md and REQUIREMENTS.md
- Version injection via `-ldflags` — deferred (D-10)
- CI/GitHub Actions for automated cross-compile — out of Phase 3 scope

</deferred>

---

*Phase: 3-Build and Distribution*
*Context gathered: 2026-07-06*
