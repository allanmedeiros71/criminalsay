# CriminalSay — Contexto do Projeto (GSD / SDD)

> Documento de contexto para desenvolvimento spec-driven com **GSD Core** no **Claude Code**.
> Serve de insumo para a fase **Discuss** de `/gsd-new-project`: reúne decisões travadas,
> escopo, especificação funcional e critérios de aceitação **antes** do planejamento.

---

## 1. Visão geral

**CriminalSay** é uma ferramenta de linha de comando (CLI) que exibe uma citação
aleatória da série *Criminal Minds* no terminal, com formatação visual elaborada
(bordas, cores, padding), **sem interatividade**. Um comando, uma citação, uma saída bonita.

- **Autor:** Allan (`github.com/allanmedeiros71`)
- **Repositório alvo:** `github.com/allanmedeiros71/criminalsay`
- **Licença:** MIT
- **Idioma do código/comentários:** inglês · **Idioma da documentação de usuário:** português
- **Status:** greenfield (projeto novo, do zero)

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

## 6. Arquitetura e estrutura proposta

> Proposta inicial; o `/gsd-plan` pode refinar. Layout idiomático Go.

```
criminalsay/
├── main.go                 # ponto de entrada: sorteia e renderiza
├── go.mod
├── go.sum
├── data/
│   └── quotes.json         # base de citações (go:embed)
├── internal/
│   ├── quotes/
│   │   ├── quotes.go       # carga do JSON embutido, tipo Quote, seleção aleatória
│   │   └── quotes_test.go
│   └── render/
│       ├── render.go       # estilo lipgloss, função Render(Quote) string
│       └── render_test.go
├── scripts/
│   └── build.sh            # cross-compile GOOS/GOARCH para os 3 SOs
├── .gitignore              # /dist, binários
├── LICENSE                 # MIT
└── README.md               # PT-BR: instalação, uso, build
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

- [ ] `criminalsay` exibe uma citação aleatória formatada e sai com código 0.
- [ ] Base embutida com ≥15 citações reais, atribuídas corretamente.
- [ ] Estilo 4 (Terminal cru) implementado com lipgloss conforme seção 6.1, com quebra de linha correta e barra lateral alinhada.
- [ ] Degradação graciosa sem cor (`NO_COLOR` / pipe / não-TTY).
- [ ] `go test ./...` passa; cobertura razoável nos dois pacotes internos.
- [ ] Binários cross-compilados para Linux, macOS (amd64+arm64) e Windows.
- [ ] README em PT-BR com instalação, uso e instruções de build.
- [ ] Repositório publicado em `github.com/allanmedeiros71/criminalsay` com licença MIT.

---

## 10. Decisões em aberto (para a fase Discuss confirmar)

1. **Personagem "narrador":** exibir sempre o personagem que citou (quando conhecido) ou só o autor original da frase?
2. **Idioma das citações:** manter em inglês (original da série) ou incluir tradução PT-BR?
3. **Largura fixa vs. adaptativa** ao tamanho do terminal (a spec do Estilo 4 sugere largura fixa; confirmar).

---

## 11. Como conduzir com GSD Core

```bash
# 1. Instalar o GSD Core no runtime (Claude Code)
npx @opengsd/gsd-core@latest

# 2. Dentro do diretório do projeto, iniciar o projeto no Claude Code
/gsd-new-project
#   → forneça este CONTEXT.md como insumo da fase Discuss

# 3. Seguir o loop de fases do GSD:
#    Discuss  → confirmar as "decisões em aberto" da seção 10
#    Plan     → decompor em fases/tarefas (respeitar decisões travadas da seção 2)
#    Execute  → implementar em ondas com contexto limpo
#    Verify   → rodar testes/critérios da seção 7 e 9 antes de declarar done
#    Ship     → abrir PR, arquivar fase, repetir
```

> **Nota:** o instalador do GSD é obrigatório para compatibilidade entre runtimes —
> não copie arquivos de `agents/` ou `commands/` manualmente.
