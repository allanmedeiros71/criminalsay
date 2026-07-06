# CriminalSay

Ferramenta de linha de comando que exibe uma citação aleatória de *Criminal Minds* no terminal, com formatação estilizada — barra lateral vermelha, texto destacado e atribuição discreta. Um comando, uma citação, sem interatividade.

## Instalação

### Via Go install

Requer [Go 1.22+](https://go.dev/dl/).

```bash
go install github.com/allanmedeiros71/criminalsay@latest
```

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

## Uso

```bash
criminalsay
```

O programa escolhe uma citação aleatória, imprime em stdout e encerra com código `0`.

### Hook no shell

Adicione ao seu `~/.bashrc` ou `~/.zshrc` para exibir uma citação ao abrir o terminal:

```bash
criminalsay
```

### Saída sem cor

A cor é desativada automaticamente quando a saída é redirecionada ou quando `NO_COLOR` está definido:

```bash
NO_COLOR=1 criminalsay
criminalsay | cat
```

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
