# CriminalSay

## What This Is

CriminalSay é uma ferramenta de linha de comando (CLI) em Go que exibe uma citação aleatória da série *Criminal Minds* no terminal com formatação visual elaborada — barra lateral vermelha, texto destacado, atribuição dimmed. Um comando, uma citação, uma saída bonita. Ideal para rodar no `.bashrc`/`.zshrc` a cada abertura de shell, no espírito do `fortune`.

## Core Value

Rodar `criminalsay` e receber imediatamente uma citação memorável do BAU com estilo, sem interatividade e sem dependências externas.

## Requirements

### Validated

- [x] Base de citações em JSON embutida com ≥15 citações reais (Phase 1)
- [x] Comando `criminalsay` exibe citação aleatória formatada e encerra com exit code 0 (Phase 2)
- [x] Estilo visual Estilo 4: barra lateral `▌`, texto entre aspas, atribuição recuada (Phase 2)
- [x] Quebra de linha com barra lateral alinhada em cada linha (Phase 2)
- [x] Degradação graciosa sem cor — NO_COLOR / pipe / não-TTY (Phase 2)
- [x] `go test ./...` passa em quotes/ e render/ (Phases 1–2)

### Active

- [ ] Build cross-platform: Linux amd64, macOS amd64+arm64, Windows amd64
- [ ] README em PT-BR com instalação, uso e instruções de build

### Out of Scope

- Flags de CLI (filtrar por personagem/temporada) — v2
- Múltiplos temas selecionáveis em runtime — v2
- Modo interativo / TUI navegável (bubbletea) — explicitamente descartado
- Busca/scraping de citações em runtime — descartado (dados offline-first por decisão D5)
- Distribuição via Homebrew/AUR/pacotes — v2

## Context

- Projeto greenfield em Go, autor: Allan (github.com/allanmedeiros71)
- Repositório: github.com/allanmedeiros71/criminalsay
- Licença: MIT
- Código e comentários em inglês; documentação de usuário em português
- Estilo visual já escolhido e travado: Estilo 4 "Terminal cru" (seção 6.1 do CONTEXT.md)
- Estrutura de pacotes já definida: main.go + internal/quotes/ + internal/render/ + data/quotes.json

## Constraints

- **Tech stack**: Go + charmbracelet/lipgloss — decisão travada (D2)
- **Dados**: JSON local embutido via go:embed — decisão travada (D5/D6)
- **Interatividade**: Nenhuma — escopo é exibir e sair (D4)
- **Aleatoriedade**: math/rand/v2 — sorteio simples, sem crypto (D7)
- **Plataformas**: Linux, macOS, Windows — cross-compile (D8)
- **Desempenho**: partida a frio <50ms típico (NFR-1)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Linguagem: Go | Binário único, cross-compile trivial, distribuição simples | — Pending |
| Estilização: lipgloss | Bordas, layout, cores adaptativas a truecolor/256/ANSI-16 | — Pending |
| Cor complementar: fatih/color | Leve, para casos simples fora do layout lipgloss | — Pending |
| Fonte de dados: JSON local | Sem API externa, offline-first | — Pending |
| Embed: go:embed | Binário autocontido, sem arquivos externos | — Pending |
| Estilo visual: Estilo 4 "Terminal cru" | Leve, sem caixa fechada; ideal para .bashrc | — Pending |
| Aleatoriedade: math/rand/v2 | Sorteio simples, sem necessidade de crypto | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-07-05 after initialization*
