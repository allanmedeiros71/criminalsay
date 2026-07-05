# CriminalSay — Contexto do Projeto (GSD / SDD)

> Documento mestre do projeto **CriminalSay** para humanos e agentes.
> Contém: decisões travadas, especificação funcional, critérios de aceitação **e**
> o fluxo GSD completo adaptado ao estado atual do repositório.

---

## 1. Visão geral

**CriminalSay** é uma ferramenta de linha de comando (CLI) que exibe uma citação
aleatória da série *Criminal Minds* no terminal, com formatação visual elaborada
(bordas, cores, padding), **sem interatividade**. Um comando, uma citação, uma saída bonita.

- **Autor:** Allan (`github.com/allanmedeiros71`)
- **Repositório alvo:** `github.com/allanmedeiros71/criminalsay`
- **Licença:** MIT
- **Idioma do código/comentários:** inglês · **Idioma da documentação de usuário:** português
- **Status:** em desenvolvimento (milestone v1.0 — Fase 01 shipped, Fase 02 em andamento)

### Elevator pitch
Rode `criminalsay` e receba, com estilo, uma citação memorável do BAU — ideal para
rodar no `.bashrc`/`.zshrc` a cada abertura de terminal, no espírito do `fortune`.

---

## 2. Decisões travadas (Discuss → decisões que NÃO devem ser reabertas no planejamento)

Estas decisões já foram discutidas e estão fechadas. O plano deve respeitá-las.

| # | Decisão | Escolha | Justificativa |
|---|---------|---------|---------------|
| D1 | Linguagem | **Go** | Binário único, cross-compile trivial, distribuição simples |
| D2 | Estilização de terminal | **charmbracelet/lipgloss** | Bordas, layout, cores adaptativas a truecolor/256/ANSI-16 |
| D3 | Cor complementar (se necessário) | **fatih/color** | Leve, para casos simples fora do layout do lipgloss |
| D4 | Interatividade | **Nenhuma** | Escopo é exibir e sair; sem TUI, sem bubbletea |
| D5 | Fonte de dados | **JSON local embutido** | Sem API externa, sem scraping em runtime, offline-first |
| D6 | Embutir dados no binário | **`go:embed`** | Binário autocontido, roda sem arquivos externos |
| D7 | Aleatoriedade | **`math/rand/v2`** | Sorteio simples de citação; não precisa de crypto |
| D8 | Plataformas alvo | **Linux, macOS, Windows** | Dev em Linux/macOS, cross-compile para os três |
| D9 | Nome do binário/projeto | **`criminalsay`** | Definido pelo autor |
| D10 | Estilo visual da v1 | **Estilo 4 — Terminal cru** | Leve, sem caixa fechada; ideal para rodar no `.bashrc` a cada abertura de shell |

---

## 3. Escopo

### Dentro do escopo (v1 / MVP)
- Comando `criminalsay` que imprime **uma** citação aleatória formatada e encerra.
- Base de citações em JSON embutido (campos: texto, autor original, personagem que citou, temporada, episódio, título do episódio).
- Um estilo visual padrão de saída (o autor escolherá entre os mockups já produzidos: *Dossiê*, *Minimalista*, *Ficha do perfilador*, *Terminal cru*).
- Detecção/adaptação de capacidade de cor do terminal (via lipgloss).
- Build cross-platform reproduzível (script ou Makefile com GOOS/GOARCH).
- README em português com instruções de instalação, uso e build.

### Fora do escopo (v1) — anotado para não vazar no plano
- Flags de CLI (ex.: filtrar por personagem/temporada) → candidato a v2.
- Múltiplos temas selecionáveis em runtime → candidato a v2.
- Modo interativo / TUI navegável (bubbletea) → explicitamente descartado.
- Busca/scraping de citações em runtime → descartado (ver D5).
- Distribuição via Homebrew/AUR/pacotes → candidato a v2.

---

## 4. Especificação funcional

### FR-1 — Exibir citação aleatória
Ao executar `criminalsay` sem argumentos, o programa seleciona uniformemente ao acaso
uma citação da base embutida e a imprime formatada no `stdout`, encerrando com exit code `0`.

### FR-2 — Estrutura de dados da citação
Cada citação possui: `quote` (texto), `author` (autor original da frase),
`character` (personagem da série que a citou; pode ser vazio), `season` (int),
`episode` (int), `episodeTitle` (string). Campos ausentes são renderizados de forma graciosa (omitidos, sem quebrar o layout).

### FR-3 — Renderização visual (Estilo 4 — Terminal cru)
A saída usa lipgloss no estilo **Terminal cru** (ver seção 6.1): **sem caixa fechada**,
apenas uma barra lateral vermelha (`▌`) marcando as linhas da citação, seguida de uma
linha de atribuição recuada. O texto longo deve quebrar dentro de uma largura definida,
com a barra lateral repetida em cada linha quebrada, sem desalinhar.

### FR-4 — Adaptação de cor
Em terminais sem suporte a cor (ex.: saída redirecionada para arquivo / pipe, `NO_COLOR` setado),
a saída degrada para texto plano legível, sem sequências ANSI cruas.

### FR-5 — Base de dados inicial
A base embutida deve conter, no mínimo, ~15–20 citações reais da série para o MVP,
com atribuição correta de personagem e episódio quando conhecida.

---

## 5. Requisitos não-funcionais

- **NFR-1 (Desempenho):** partida a frio deve ser praticamente instantânea (<50ms típico); adequado para rodar a cada abertura de shell.
- **NFR-2 (Portabilidade):** um único codebase compila para os 3 SOs sem `build tags` específicas de plataforma.
- **NFR-3 (Zero dependências de runtime):** binário autocontido; nenhum arquivo externo necessário.
- **NFR-4 (Legibilidade):** respeitar `NO_COLOR` e detecção de TTY.
- **NFR-5 (Manutenibilidade):** adicionar uma nova citação = editar um JSON, sem tocar em código.

---

## 6. Arquitetura e estrutura do repositório

> Layout atual (pós-Fase 01). Fase 02 adiciona `internal/render/`; Fase 03 adiciona `Makefile`/`dist/`.

```
criminalsay/
├── main.go                 # embed anchor, Load → Random → stub (Fase 02: render)
├── go.mod                  # lipgloss v1.1.0 pre-declarado para Fase 02
├── go.sum
├── data/
│   └── quotes.json         # 20 citações Criminal Minds (go:embed)
├── internal/
│   ├── quotes/             # ✓ Fase 01
│   │   ├── quotes.go       # Quote struct, Load, Random
│   │   └── quotes_test.go
│   └── render/             # Fase 02 (pendente)
│       ├── render.go       # Estilo 4 lipgloss, Render(Quote) string
│       └── render_test.go
├── .planning/              # artefatos GSD (ver seção 11.3)
├── CONTEXT.md              # este documento
├── LICENSE                 # MIT
└── README.md               # EN (Fase 03: README PT-BR)
```

### Contratos de módulo (para orientar planejamento e testes)
- `quotes.Load() ([]Quote, error)` — parseia o JSON embutido.
- `quotes.Random(qs []Quote) Quote` — retorna uma citação uniformemente aleatória.
- `render.Quote(q Quote, colorEnabled bool) string` — produz a saída formatada (sem imprimir).
- `main` orquestra: `Load` → `Random` → detecta cor → `render.Quote` → `fmt.Println`.

Separar **produção da string** da **impressão** é intencional: torna a renderização testável
sem capturar stdout.

### 6.1 Especificação do estilo visual (Estilo 4 — Terminal cru)

Estilo travado para a v1. Referência visual da saída:

```
  ▌ "It's alchemy. Alchemy turns common metals into
  ▌ precious ones. Dreams work the same way. Turning
  ▌ something awful into something better."

       David Rossi · Criminal Minds · S9E13
```

Regras de renderização:
- **Barra lateral:** caractere `▌` (U+258C) em **vermelho**, prefixando cada linha do texto da citação, com um espaço antes e depois.
- **Recuo:** duas colunas de indentação à esquerda da barra.
- **Texto da citação:** entre aspas duplas, cor de destaque (branco/claro), quebrado a uma largura de bloco fixa (sugestão: 52 colúnas úteis; ajustável).
- **Linha em branco** entre a citação e a atribuição.
- **Atribuição:** recuada, formato `<autor> · <série/personagem> · S<temporada>E<episódio>`, com o autor em cor de destaque (amarelo) e o restante em tom apagado (`dim`), separados por `·` (U+00B7).
- **Sem borda de caixa**, sem cabeçalho, sem rodapé — a leveza é intencional.
- Em modo sem cor (FR-4): a barra `▌`, as aspas e os separadores permanecem; apenas as sequências ANSI são omitidas.

---

## 7. Estratégia de testes (o GSD verifica antes de declarar "done")

- **quotes_test.go:** JSON embutido parseia sem erro; base tem N≥15; `Random` só retorna itens da base; distribuição não é trivialmente enviesada (amostragem).
- **render_test.go:** saída contém o texto da citação e o autor; cada linha da citação começa com a barra lateral `▌`; nenhuma linha excede a largura de bloco definida; com `colorEnabled=false` não há códigos ANSI (`\x1b[`) na saída.
- **Fumaça manual multiplataforma:** rodar o binário em Linux/macOS/Windows Terminal e conferir bordas e cores.

Critério de "verde": `go build ./...` e `go test ./...` limpos; binário roda e exibe citação formatada.

---

## 8. Build & distribuição

- **Dev:** `go run .`
- **Build local:** `go build -o criminalsay .`
- **Cross-compile** (`scripts/build.sh`), gerando em `dist/`:
  - `GOOS=linux   GOARCH=amd64`
  - `GOOS=darwin  GOARCH=amd64` e `GOOS=darwin  GOARCH=arm64` (Apple Silicon)
  - `GOOS=windows GOARCH=amd64` (saída `.exe`)
- Uso sugerido no `.zshrc`/`.bashrc`: chamar `criminalsay` no fim do arquivo.

---

## 9. Critérios de aceitação do MVP (Definition of Done)

- [x] Base embutida com ≥15 citações reais, atribuídas corretamente. *(Fase 01 — 20 quotes)*
- [x] `go test ./internal/quotes/...` passa com cobertura adequada. *(Fase 01)*
- [x] `go build -o criminalsay .` produz binário funcional. *(Fase 01 — stub output)*
- [ ] `criminalsay` exibe citação aleatória **formatada** (Estilo 4) e sai com código 0. *(Fase 02)*
- [ ] Estilo 4 (Terminal cru) com lipgloss: barra `▌`, wrap ≤52 col, atribuição dimmed. *(Fase 02)*
- [ ] Degradação graciosa sem cor (`NO_COLOR` / pipe / não-TTY). *(Fase 02)*
- [ ] `go test ./...` passa nos pacotes `quotes` e `render`. *(Fase 02)*
- [ ] Binários cross-compilados para Linux, macOS (amd64+arm64) e Windows. *(Fase 03)*
- [ ] README em PT-BR com instalação, uso e instruções de build. *(Fase 03)*
- [ ] Repositório publicado em `github.com/allanmedeiros71/criminalsay` com licença MIT.

---

## 10. Decisões resolvidas na Fase 01 (referência)

| Tópico | Decisão |
|--------|---------|
| Personagem vs. autor | Campo `author` = autor original; `character` = personagem CM que citou (D-01/D-02) |
| Idioma das citações | Inglês (original da série) — D-04 |
| Largura do bloco | Fixa ≤52 colúteis úteis — definida na Fase 02 (REND-02) |
| Versão lipgloss | **v1.1.0** (não v2) — pre-declarada em `go.mod` |
| Anchor do embed | **`main.go` only** — Go proíbe `..` em patterns de embed |

---

## 11. Fluxo GSD completo — CriminalSay

> **GSD Core** (Git. Ship. Done.) — desenvolvimento spec-driven com planejamento em `.planning/`,
> execução fase a fase, verificação automatizada e ship via PR.
> Runtime: **Cursor** (`~/.cursor/gsd-core/`). Modo deste projeto: **`yolo`**.

### 11.1 Visão geral do loop

```text
┌─────────────────────────────────────────────────────────────────────────┐
│  MILESTONE (v1.0)                                                       │
│                                                                         │
│  /gsd-new-project  ──►  PROJECT · REQUIREMENTS · ROADMAP · STATE        │
│         │                                                               │
│         ▼                                                               │
│  ┌─── POR FASE (1 → 2 → 3) ───────────────────────────────────────┐   │
│  │  discuss ─► plan ─► execute ─► validate ─► verify ─► secure   │   │
│  │     │                                              │            │   │
│  │     └──────────────────────────────────────────────┼──► ship   │   │
│  └────────────────────────────────────────────────────┘            │   │
│         │                                                               │
│         ▼                                                               │
│  /gsd-complete-milestone 1.0.0  ──►  tag · archive · próximo ciclo     │
└─────────────────────────────────────────────────────────────────────────┘
```

Comando de roteamento universal quando não souber qual usar:

```text
/gsd-progress              # onde estou, o que falta
/gsd-progress --do "..."   # roteia intenção em linguagem natural
```

### 11.2 Estado atual deste projeto

| Item | Valor |
|------|-------|
| Milestone | **v1.0** |
| Fases no roadmap | 3 (Data Foundation → Render and Output → Build and Distribution) |
| Fase atual | **2 — Render and Output** (contexto coletado, plano criado) |
| Fase 01 | ✓ Completa · verificada · **shipped — [PR #1](https://github.com/allanmedeiros71/criminalsay/pull/1)** |
| Branch de trabalho | `main` (local) · PR branch: `gsd/phase-01-data-foundation-pr` |
| Retomar sessão | `.planning/phases/02-render-and-output/02-CONTEXT.md` |

**Progresso de requisitos (v1):**

| Grupo | Status |
|-------|--------|
| DAT-01…03, CORE-01, CORE-03, BUILD-01, BUILD-03 | ✓ Fase 01 |
| CORE-02, REND-01…06, COLOR-01…02 | Pendente — Fase 02 |
| BUILD-02, DOC-01 | Pendente — Fase 03 |

### 11.3 Estrutura `.planning/` (mapa de artefatos)

```text
.planning/
├── PROJECT.md          # visão, valor central, decisões-chave
├── REQUIREMENTS.md     # REQ-IDs rastreáveis (DAT, CORE, REND, COLOR, BUILD, DOC)
├── ROADMAP.md          # fases, critérios de sucesso, dependências
├── STATE.md            # memória viva — fase atual, progresso, blockers
├── config.json         # toggles de workflow (modo, segurança, ship, branches)
├── research/           # pesquisa de domínio (criado no new-project)
└── phases/
    └── NN-slug/
        ├── NN-CONTEXT.md       # visão da fase (discuss)
        ├── NN-RESEARCH.md      # pesquisa técnica da fase
        ├── NN-01-PLAN.md       # plano executável (tarefas, must-haves, threat_model)
        ├── NN-01-SUMMARY.md    # o que foi construído (pós-execução)
        ├── NN-VERIFICATION.md  # relatório do verifier (status: passed/blocked)
        ├── NN-VALIDATION.md    # mapa de cobertura de testes (Nyquist)
        ├── NN-SECURITY.md      # threat register — gate de ship (threats_open: 0)
        └── NN-PATTERNS.md      # padrões descobertos na fase
```

**Arquivos estruturais** (viajam no PR): `STATE.md`, `ROADMAP.md`, `PROJECT.md`, `REQUIREMENTS.md`.
**Arquivos transientes** (ruído para review — filtrados por `/gsd-pr-branch`): `phases/**`, `research/**`.

### 11.4 Fase 0 — Inicialização (já concluída)

```bash
# Instalar/atualizar GSD Core
npx @opengsd/gsd-core@latest

# Inicializar projeto (forneceu este CONTEXT.md como insumo)
/gsd-new-project
```

**Entregáveis gerados:** `.planning/PROJECT.md`, `REQUIREMENTS.md`, `ROADMAP.md`, `STATE.md`, `config.json`, pesquisa de domínio.

Para codebases existentes (não é o caso aqui): `/gsd-map-codebase` antes do new-project.

### 11.5 Loop por fase — CriminalSay v1.0

Cada fase do roadmap segue a mesma sequência. Exemplo com as três fases deste projeto:

#### Passo 1 — Discuss (opcional, recomendado)

```text
/gsd-discuss-phase 2
```

- Captura visão, essenciais e limites antes do plano.
- Gera/atualiza `02-CONTEXT.md`.
- **Fase 02:** contexto já coletado — retomar de `02-CONTEXT.md` se necessário.

#### Passo 2 — Plan

```text
/gsd-plan-phase 2
/gsd-plan-phase 2 --research        # força nova pesquisa
/gsd-plan-phase 2 --tdd             # ordem test-first no plano
/gsd-plan-phase 2 --mvp             # fatia vertical MVP
```

- Gera `02-01-PLAN.md` com tarefas, must-haves, `<threat_model>` e critérios de verificação.
- **Fase 01:** plano `01-01-PLAN.md` (Walking Skeleton) — concluído.
- **Fase 02:** plano `02-01-PLAN.md` (Estilo 4 render) — criado, pendente execução.
- **Fase 03:** planos TBD.

#### Passo 3 — Execute

```text
/gsd-execute-phase 2
/gsd-execute-phase 2 --wave 1      # só a onda 1
/gsd-execute-phase 2 --tdd          # enforce RED/GREEN
```

- Executa planos em ondas paralelas; commits atômicos convencionais.
- Atualiza `REQUIREMENTS.md`, `ROADMAP.md`, `STATE.md`.
- **Fase 01 entregou:** `go.mod`, `data/quotes.json`, `internal/quotes/`, `main.go` (stub `<quote> — <author>`).
- **Fase 02 entregará:** `internal/render/`, wiring lipgloss Estilo 4, `go test ./...` green.

#### Passo 4 — Validate (Nyquist)

```text
/gsd-validate-phase 2
```

- Audita lacunas entre plano, testes e requisitos.
- Preenche/atualiza `02-VALIDATION.md`.

#### Passo 5 — Verify (UAT)

```text
/gsd-verify-work 2
```

- Verificação conversacional + evidências automatizadas.
- Gera `02-VERIFICATION.md` com `status: passed` (obrigatório para ship).

Critérios de aceitação deste projeto (seções 7 e 9 deste documento) são a referência humana; o verifier usa must-haves do plano e REQ-IDs.

#### Passo 6 — Secure (gate obrigatório neste projeto)

```text
/gsd-secure-phase 2
```

- **`workflow.security_enforcement: true`** em `.planning/config.json`.
- Exige `NN-SECURITY.md` com **`threats_open: 0`** antes de `/gsd-ship`.
- ASVS Level 1 — threat model do `PLAN.md` + mitigações verificadas.
- **Fase 01:** `01-SECURITY.md` criado (5 ameaças fechadas; T-01-02 mitigada em `main.go`).

#### Passo 7 — Ship

```text
/gsd-ship 2
/gsd-pr-branch 2          # branch limpa (só código + planning estrutural)
```

**Preflight checks (todos obrigatórios):**

1. `VERIFICATION.md` → `status: passed`
2. Working tree limpa (sem changes unstaged)
3. Branch de feature (não `main` direto — `branching_strategy: none` cria branch sob demanda)
4. Remote `origin` + `gh` autenticado
5. Security gate → `threats_open: 0`

**Fluxo de ship deste repositório:**

```text
/gsd-secure-phase N
/gsd-pr-branch N                    # gsd/phase-0N-slug-pr — filtra phases/** do diff
/gsd-ship N                         # push + PR body rico + update STATE.md
```

Templates de branch (`.planning/config.json`):

- Fase: `gsd/phase-{phase}-{slug}` → ex.: `gsd/phase-02-render-and-output`
- PR: `{branch}-pr` → ex.: `gsd/phase-02-render-and-output-pr`

**Fase 01 — referência:** [PR #1](https://github.com/allanmedeiros71/criminalsay/pull/1) · branch `gsd/phase-01-data-foundation-pr`.

#### Passo 8 — Repetir ou fechar milestone

```text
/gsd-execute-phase 3      # última fase
/gsd-complete-milestone 1.0.0
```

Após Fase 3 (Makefile cross-compile + README PT-BR): arquivar milestone, tag git, preparar v2.

### 11.6 Roadmap CriminalSay — comandos por fase

| Fase | Nome | Objetivo | Comando seguinte |
|------|------|----------|------------------|
| **1** | Data Foundation | JSON embed, `internal/quotes`, Walking Skeleton | ✓ Shipped — merge [PR #1](https://github.com/allanmedeiros71/criminalsay/pull/1) |
| **2** | Render and Output | lipgloss Estilo 4, word-wrap, degradacao de cor | `/gsd-execute-phase 2` |
| **3** | Build and Distribution | Makefile, `dist/`, README PT-BR | `/gsd-plan-phase 3` (após Fase 2) |

### 11.7 Comandos auxiliares (uso frequente)

| Comando | Quando usar |
|---------|-------------|
| `/gsd-progress` | Ver barra de progresso, fase atual, próxima ação |
| `/gsd-quick` | Tarefa ad-hoc pequena (`.planning/quick/`) |
| `/gsd-fast "..."` | Fix trivial inline (≤3 arquivos, sem plano) |
| `/gsd-debug "..."` | Debug persistente entre sessões |
| `/gsd-capture` | Guardar ideia, todo ou seed |
| `/gsd-code-review` | Review pós-execução da fase |
| `/gsd-ui-review` | Audit visual (Fase 02 tem UI hint) |
| `/gsd-pause-work` | Handoff ao pausar mid-phase |
| `/gsd-resume-work` | Retomar sessão anterior |
| `/gsd-help --full` | Referência completa de comandos |

### 11.8 Configuração GSD deste projeto

Arquivo: `.planning/config.json`

| Toggle | Valor | Efeito |
|--------|-------|--------|
| `mode` | `yolo` | Execução autônoma, menos prompts |
| `commit_docs` | `true` | Commits automáticos de `.planning/` |
| `branching_strategy` | `none` | Sem branch automática — criar manualmente ou via `/gsd-pr-branch` |
| `workflow.security_enforcement` | `true` | `/gsd-secure-phase` obrigatório antes de ship |
| `workflow.security_asvs_level` | `1` | Verificação grep-level (L1) |
| `workflow.code_review` | `true` | Review oferecido pós-ship |
| `workflow.ui_phase` | `true` | UI-SPEC disponível para Fase 02 |
| `workflow.nyquist_validation` | `true` | `/gsd-validate-phase` ativo |
| `workflow.verifier` | `true` | `/gsd-verify-work` ativo |

### 11.9 Gates de qualidade — CriminalSay

```text
Execute
   │
   ├─► go test ./...          (BUILD-03 — por fase)
   ├─► go build -o criminalsay .
   │
Validate (Nyquist)
   │
Verify (UAT → VERIFICATION.md status: passed)
   │
Secure (SECURITY.md → threats_open: 0)
   │
Ship (push + gh pr create + STATE.md)
```

**Comandos de verificação manual rápida (Fase 02+):**

```bash
go build -o criminalsay .
./criminalsay                        # TTY colorido
NO_COLOR=1 ./criminalsay             # sem ANSI
./criminalsay | cat                  # pipe — texto legível
go test ./... -count=1               # suite completa
```

### 11.10 Próximo passo imediato

```text
/gsd-execute-phase 2
```

Implementar `internal/render/` (Estilo 4), wiring em `main.go`, testes REND/COLOR,
substituir stub `<quote> — <author>` por saída lipgloss formatada.

Se precisar revisar decisões visuais antes de executar:

```text
/gsd-discuss-phase 2 --view    # reler 02-CONTEXT.md interativamente
```

---

## 12. Instalação e manutenção do GSD

```bash
# Instalar/atualizar GSD Core no runtime Cursor
npx @opengsd/gsd-core@latest

# Diagnóstico do diretório de planning
/gsd-health

# Ver todos os comandos
/gsd-help --full
```

> **Nota:** use sempre o instalador oficial — não copie `agents/` ou `commands/` manualmente entre runtimes.
