# Phase 3: Build and Distribution - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-07-06
**Phase:** 3-Build and Distribution
**Areas discussed:** Ferramenta de build, Estratégia do README, Targets do Makefile, Convenção dist/

---

## Ferramenta de build

| Option | Description | Selected |
|--------|-------------|----------|
| Só Makefile | Interface principal; alinha com ROADMAP e research | ✓ |
| Só scripts/build.sh | Script POSIX autônomo, sem Make | |
| Ambos | Makefile como interface; build.sh com lógica | |
| Você decide | Claude escolhe | |

**User's choice:** Só Makefile
**Notes:** CONTEXT.md original propunha ambos; usuário optou por simplificar.

| Option | Description | Selected |
|--------|-------------|----------|
| Targets explícitos por plataforma | linux, darwin-amd64, darwin-arm64, windows + all | ✓ |
| Mínimo — só all | Um target que compila os 4 de uma vez | |
| Você decide | | |

**User's choice:** Targets explícitos por plataforma

| Option | Description | Selected |
|--------|-------------|----------|
| CGO_ENABLED=0 por linha | Explícito em cada build | |
| export no topo | DRY, uma declaração | ✓ (Claude discretion) |
| Você decide | | ✓ |

**User's choice:** Você decide → `export CGO_ENABLED=0` no topo

| Option | Description | Selected |
|--------|-------------|----------|
| mkdir -p dist no Makefile | Cria dist/ automaticamente | ✓ |
| Manual | Usuário cria dist/ | |

**User's choice:** Makefile cria dist/ automaticamente

---

## Estratégia do README

| Option | Description | Selected |
|--------|-------------|----------|
| Substituir por PT-BR | README.md inteiro em português | |
| Bilíngue | EN + PT no mesmo arquivo | |
| README.pt.md separado | Manter EN, adicionar PT | ✓ |
| Você decide | | |

**User's choice:** README.pt.md separado

| Option | Description | Selected |
|--------|-------------|----------|
| Atualizar ambos | EN corrige inconsistências; PT espelha | ✓ |
| Só criar PT | EN fica como está | |
| Você decide | | |

**User's choice:** Sincronizar ambos (corrigir Go 1.22+, refs reais ao Makefile)

| Option | Description | Selected |
|--------|-------------|----------|
| Seções completas | Instalação, uso, build, estrutura, dev, licença | |
| Mínimo DOC-01 | Só instalação, uso e build | ✓ |
| Você decide | | |

**User's choice:** README.pt.md mínimo (DOC-01)

| Option | Description | Selected |
|--------|-------------|----------|
| go install + build local | Sem GitHub Releases | ✓ |
| Mencionar Releases | Opção futura/manual | |
| Você decide | | |

**User's choice:** Apenas go install e build local

---

## Targets do Makefile

| Option | Description | Selected |
|--------|-------------|----------|
| build + test + clean + all | Targets de conveniência completos | ✓ (Claude discretion) |
| test + all | Sem target build | |
| Só all + plataformas | Mínimo | |
| Você decide | | ✓ |

**User's choice:** Você decide → build, test, clean + cross-compile targets

| Option | Description | Selected |
|--------|-------------|----------|
| Sem ldflags | Binário simples v1 | ✓ (Claude discretion) |
| git describe versão | -ldflags com tag | |
| Você decide | | ✓ |

**User's choice:** Você decide → sem injeção de versão em v1

| Option | Description | Selected |
|--------|-------------|----------|
| make test = go test ./... | Independente de all | ✓ (Claude discretion) |
| all depende de test | Só compila se testes passam | |
| Você decide | | ✓ |

**User's choice:** Você decide → test independente de all

| Option | Description | Selected |
|--------|-------------|----------|
| clean dist/ + criminalsay | Remove artefatos locais | ✓ |
| clean só dist/ | | |
| Você decide | | |

**User's choice:** clean remove dist/ e binário criminalsay na raiz

---

## Convenção dist/

| Option | Description | Selected |
|--------|-------------|----------|
| Flat com prefixo criminalsay-* | Já documentado no README EN | ✓ (Claude discretion) |
| Flat simples / subpastas | | |
| Você decide | | ✓ |

**User's choice:** Você decide → flat prefixed (criminalsay-linux-amd64, etc.)

| Option | Description | Selected |
|--------|-------------|----------|
| Com .exe | criminalsay-windows-amd64.exe | ✓ |
| Sem .exe | Consistente com outros | |

**User's choice:** Windows com sufixo .exe

| Option | Description | Selected |
|--------|-------------|----------|
| Adicionar .gitignore | dist/ e criminalsay | ✓ |
| Não adicionar | Fora do escopo | |

**User's choice:** Criar .gitignore

| Option | Description | Selected |
|--------|-------------|----------|
| Nunca commitar dist/ | Artefatos locais/CI apenas | ✓ |
| Commitar para releases | | |

**User's choice:** Nunca commitar dist/

---

## Claude's Discretion

- CGO_ENABLED=0 via export no topo do Makefile
- Targets build, test, clean além dos cross-compile
- Sem ldflags de versão em v1
- make test independente de all
- Nomenclatura flat prefixed em dist/ (alinhada ao README EN existente)

## Deferred Ideas

- scripts/build.sh (rejeitado — Makefile único)
- GitHub Releases / downloads de binários pré-compilados
- GoReleaser, Homebrew, AUR (v2)
- Injeção de versão via ldflags
- CI/GitHub Actions automatizado
