# Ambiente descartável para testar o instalador contra os quatro harnesses de
# verdade, sem sujar a máquina. Nada de volume: o binário é compilado aqui e
# copiado para dentro da imagem, então o que se testa é o conteúdo embutido
# nele — o mesmo que um `go install` entregaria.
#
#   docker build -t coder-test .
#   docker run --rm -it coder-test
#
# O -it não é conforto: o instalador abre /dev/tty para desenhar os
# formulários. Sem terminal ele cai no caminho de script e exige --harness.

# ── estágio de build ────────────────────────────────────────────────────────
# A imagem do Go fica fora da final: o container de teste não precisa de
# toolchain, e arrastá-la triplicaria o tamanho à toa.
FROM golang:1.25-bookworm AS build

WORKDIR /src

# go.mod e go.sum primeiro, sozinhos: as dependências só são baixadas de novo
# quando um deles muda, não a cada alteração no instalador.
COPY go.mod go.sum ./
RUN go mod download

# skills/ e commands/ entram aqui porque o go:embed lê do diretório de build.
# É por isso que o .dockerignore não exclui markdown.
COPY . .
RUN CGO_ENABLED=0 go build -o /out/coder .

# ── imagem de teste ─────────────────────────────────────────────────────────
# Debian, não Alpine: as quatro CLIs trazem binários nativos linkados contra
# glibc. Em musl elas instalam sem reclamar e quebram só na hora de rodar.
#
# Node 22 é exigência do @github/copilot; as outras três aceitam menos.
FROM node:22-slim

# git é o que faz o escopo de projeto existir: sem ele não há .git para o
# instalador encontrar. tree é para conferir o resultado de olho.
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates git tree \
    && rm -rf /var/lib/apt/lists/*

# Os quatro harnesses, sem versão fixada de propósito: o container existe para
# testar contra o que as pessoas instalam hoje, não contra um instantâneo.
#
# bun entra junto porque o bin do omp é um script com `#!/usr/bin/env bun`:
# sem ele o pacote instala e a CLI não sobe. Não é exigência do instalador —
# que só escreve em ~/.agents e nunca invoca a CLI —, é o que separa "o
# harness está instalado" de "o diretório dele existe".
#
# Esta camada vem antes da cópia do binário porque é a cara de construir e a
# que menos muda. Invertida, cada mexida no Go rebaixaria os cinco pacotes.
RUN npm install -g \
        opencode-ai \
        @anthropic-ai/claude-code \
        @oh-my-pi/pi-coding-agent \
        @github/copilot \
        bun \
    && npm cache clean --force

# A última camada, e a única que uma alteração no instalador invalida.
COPY --from=build /out/coder /usr/local/bin/coder

# Repositório de mentira para o escopo de projeto ter onde cair. Como é o
# WORKDIR, o instalador já abre perguntando o escopo, sem preparação nenhuma.
# Para exercitar o caminho global, `cd /tmp`, que não é repositório.
RUN git init -q -b main /work/projeto

WORKDIR /work/projeto
CMD ["bash"]
