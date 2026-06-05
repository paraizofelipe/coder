# coder

Conjunto de agentes e skills para [OpenCode](https://opencode.ai), [Claude Code](https://claude.ai/code) e [Codex](https://github.com/openai/codex) que implementa um fluxo disciplinado de desenvolvimento de software com análise prévia, TDD, revisão técnica, revisão de segurança e versionamento controlado.

## Commands

| Command | Uso | Descrição |
|---|---|---|
| `/doc-plan` | `/doc-plan` | Publica `.coder/plan.md` no Confluence (space CAT, subpágina de Implementações) via MCP `atlassian_local`; ignora se não houver diferenças |
| `/get-plan` | `/get-plan` | Baixa o plano de implementação do Confluence e salva em `.coder/plan.md`; cria o arquivo se não existir |
| `/kanban-card` | `/kanban-card <friendlyID>` | Consulta um card pelo friendlyID via MCP `kanban-force` e carrega todas as informações no contexto (ignora cards arquivados) |
| `/mr-review` | `/mr-review` | Aciona o `mr_reviewer` para revisar o Merge Request aberto na branch atual via `glab` |

## Agentes

| Agente | Mode | Função |
|---|---|---|
| `coder` | primary | Orquestrador principal — coordena subagentes de desenvolvimento e delega operações de card/board ao `kanban` |
| `lead` | primary | Orquestrador de planejamento — gera `.coder/tasks.md` e delega implementação ao `coder` após aprovação |
| `documenter` | primary | Publica planos de implementação no Confluence via MCP `atlassian_local` |
| `kanban` | primary | Gerenciamento de cards e boards via MCP `kanban-force` |
| `infra` | primary | Consulta aplicações no ArgoCD |
| `mr_reviewer` | primary | Revisa Merge Requests do GitLab via CLI `glab` |
| `analyzer` | subagent | Inspeciona a codebase antes de qualquer modificação |
| `clarifier` | subagent | Formata perguntas de ambiguidade com opções e recomendação |
| `planner` | subagent | Produz o TaskGraph esqueleto a partir da intenção esclarecida |
| `detailer` | subagent | Enriquece cada task com preview, testes, critérios e contrato |
| `tester` | subagent | Cria e executa testes com abordagem TDD |
| `code_reviewer` | subagent | Revisão técnica de código — Camada 1 |
| `business_reviewer` | subagent | Revisão de negócio e segurança — Camada 2 |
| `versioner` | subagent | Executa operações Git com confirmação explícita |

## Fluxo de desenvolvimento

```mermaid
flowchart LR

    Coder(["🤖 Coder\n(primary)"])
    Kanban(["🤖 Kanban\n(primary)"])
    Analyzer["🤖 Analyzer\n(subagent)"]
    Tester["🤖 Tester\n(subagent)"]
    Versioner["🤖 Versioner\n(subagent)"]
    CodeReviewer["🤖 CodeReviewer\n(subagent)"]
    BusinessReviewer["🤖 BusinessReviewer\n(subagent)"]

    Coder --> Analyzer
    Coder --> Versioner
    Coder --> Tester
    Coder --> Kanban
    Coder --> CodeReviewer
    Coder --> BusinessReviewer

    Kanban --> kanban-force{{kanban-force}}:::task
    kanban-force --> MCP[(mcp-kanban-force)]:::mcp
    
    Analyzer --> analyse-code{{analyse-code}}:::task

    Tester --> test-code{{test-code}}:::task

    Versioner --> version-code{{version-code}}:::task

    Coder --> write-code{{write-code}}:::task

    CodeReviewer --> review-code{{review-code}}:::task

    BusinessReviewer --> review-code{{review-code}}:::task

    classDef agent fill:#4A9,color:#fff
    classDef subAgent fill:#555,color:#fff
    classDef task fill:#1565C0,color:#fff
    classDef mcp fill:#555,color:#fff

    class Coder,Kanban agent
    class Analyzer,Tester,Versioner,CodeReviewer,BusinessReviewer subAgent
```

Quando a solicitação contiver ID de card ou operação de board/card (criar, mover, atualizar, comentar, bloquear, arquivar etc.), o `coder` delega ao `kanban`, que opera via MCP `kanban-force`. Para solicitações mistas, o `kanban` executa primeiro e o fluxo de código segue depois.

O `tester` é acionado em dois momentos distintos: antes da implementação para criar os testes que devem falhar (fase red do TDD) e depois da implementação para confirmar que todos passam (fase green). Nenhum código é versionado sem o parecer final do `business_reviewer`.

## Instalação

### Via curl (recomendado)

```bash
curl -fsSL https://raw.githubusercontent.com/paraizofelipe/coder/main/install.sh | bash
```

### Via wget

```bash
wget -qO- https://raw.githubusercontent.com/paraizofelipe/coder/main/install.sh | bash
```

### A partir do repositório local

Clone o repositório e execute o script com a flag `--local`:

```bash
git clone https://github.com/paraizofelipe/coder.git
cd coder
./install.sh --local
```

Ao executar o instalador, a primeira etapa é **selecionar o(s) harness(es)** de destino (OpenCode, Claude Code, Codex ou todos). Em seguida, se OpenCode estiver entre os selecionados, há a opção de escolher o vendor de modelos. Ambas as etapas podem ser ignoradas fornecendo as flags `--harness` e `--vendor` diretamente.

## Seleção de harness

No início da instalação, o instalador exibe um menu interativo para escolher para qual(is) harness(es) instalar:

```text
[info]  Selecione o(s) harness(es) de destino:
        1) opencode
        2) claude
        3) codex
        4) todos

[?]    Números separados por espaço (ex.: 1 2):
```

Escolha um ou mais números. Para instalar em todos os harnesses de uma vez, selecione `4`. Para pular o menu, use a flag `--harness`:

```bash
# instalar apenas no Claude Code
./install.sh --local --harness claude

# instalar no OpenCode e no Claude Code
./install.sh --local --harness opencode,claude

# instalar em todos
./install.sh --local --harness all
```

## Seleção de vendor

A seleção de vendor é **opcional** e **relevante apenas para o OpenCode**. Quando o OpenCode está entre os harnesses selecionados, o instalador pergunta o vendor desejado (ou aceita Enter para manter o padrão):

```
[info]  Vendor do OpenCode (Enter para usar o default openai/gpt-5.5):

        1) anthropic        main: anthropic/claude-sonnet-4-6
        2) openai           main: openai/gpt-5.5
        3) google           main: google/gemini-2.5-pro
        4) groq             main: groq/llama-3.3-70b-versatile
        5) amazon-bedrock   main: amazon-bedrock/amazon.nova-pro-v1:0
        6) github-copilot   main: github-copilot/claude-sonnet-4.6

[?]    Número do vendor [1-6] (Enter = padrão):
```

O vendor define o modelo **main** aplicado aos agentes **primários** (`coder`, `lead`, `documenter`, `kanban`, `infra`, `mr_reviewer`) no OpenCode. O padrão é `openai/gpt-5.5`.

Para pular o menu, use a flag `--vendor`:

```bash
./install.sh --local --harness opencode --vendor anthropic
```

**Claude Code** usa `sonnet` como modelo dos primários, independentemente de vendor. **Codex** não recebe `model` por agente — herda o modelo da sessão.

> Para verificar os modelos disponíveis no seu ambiente OpenCode: `opencode models <vendor>`

## Diretórios de instalação

Os artefatos são instalados nos diretórios nativos de cada harness:

### OpenCode

| Tipo | Destino padrão |
|---|---|
| Agentes | `~/.config/opencode/agents/<nome>.md` (arquivo montado com frontmatter) |
| Skills | `~/.config/opencode/skills/<nome>/` (pasta com `SKILL.md`) |
| Commands | `~/.config/opencode/commands/<nome>.md` (arquivo montado com frontmatter) |

Override: `OPENCODE_DIR` (padrão: `~/.config/opencode`)

### Claude Code

| Tipo | Destino padrão |
|---|---|
| Agentes | `~/.claude/agents/<nome>.md` (arquivo montado com frontmatter) |
| Skills | `~/.claude/skills/<nome>/` (pasta com `SKILL.md`) |
| Commands | `~/.claude/commands/<nome>.md` (arquivo montado com frontmatter) |

Override: `CLAUDE_DIR` (padrão: `~/.claude`)

### Codex

| Tipo | Destino padrão |
|---|---|
| Skills | `~/.agents/skills/<nome>/` (pasta com `SKILL.md`) |
| Prompts (commands) | `~/.codex/prompts/<nome>.md` (apenas o corpo, sem frontmatter) |
| `AGENTS.md` | `~/.codex/AGENTS.md` (arquivo de orquestração) |

O Codex **não** recebe arquivos de agente — a orquestração é feita via `AGENTS.md`. Os commands são instalados como prompts body-only (sem frontmatter de `agent:`).

Overrides: `CODEX_DIR` (padrão: `~/.codex`) e `CODEX_SKILLS_DIR` (padrão: `~/.agents/skills`)

> **Atenção (Codex, modo remoto):** o `AGENTS.md` é gitignored e não está disponível no GitHub. Em instalação via `curl | bash` o arquivo será pulado com um aviso. Para instalar o `AGENTS.md`, use `--local`.

### Sobrescrevendo diretórios via variável de ambiente

```bash
OPENCODE_DIR=/caminho/customizado ./install.sh --local --harness opencode
CLAUDE_DIR=/outro/caminho ./install.sh --local --harness claude
CODEX_DIR=~/.meu-codex CODEX_SKILLS_DIR=~/.meu-codex/skills ./install.sh --local --harness codex
```

## Checagem antes de instalar

O instalador verifica, para cada agente e skill, se já existe um arquivo com o mesmo nome no diretório de destino. Quando encontra um conflito, exibe um aviso e pergunta se deve substituir:

```
[warn]  Já existe: /home/user/.claude/agents/coder.md
[?]    Substituir coder.md? [s/N]
```

Responda `s` para substituir ou pressione Enter para pular.

## Opções do instalador

| Flag | Descrição |
|---|---|
| `--harness <lista>` | Harness(es) a instalar sem menu interativo. Valores: `opencode`, `claude`, `codex`, `all` (separados por vírgula ou espaço, ex.: `opencode,claude`) |
| `--vendor <nome>` | Vendor para o OpenCode sem menu interativo. Valores: `anthropic`, `openai`, `google`, `groq`, `amazon-bedrock`, `github-copilot`. Inválido → erro + exit 1. Ignorado se OpenCode não estiver nos harnesses selecionados |
| `--force`, `-f` | Substitui todos os arquivos sem perguntar |
| `--local`, `-l` | Instala a partir dos arquivos locais do repositório clonado |
| `--help`, `-h` | Exibe a ajuda |

### Exemplos

```bash
# instalação interativa (recomendado para primeira vez)
curl -fsSL https://raw.githubusercontent.com/paraizofelipe/coder/main/install.sh | bash

# forçar reinstalação sem confirmações, todos os harnesses, vendor openai
./install.sh --local --force --harness all --vendor openai

# instalar apenas no OpenCode e Claude Code, vendor anthropic
./install.sh --local --harness opencode,claude --vendor anthropic

# instalar apenas no Codex (sem vendor)
./install.sh --local --harness codex

# forçar substituição em modo remoto
curl -fsSL https://raw.githubusercontent.com/paraizofelipe/coder/main/install.sh | bash -s -- --force --harness claude
```

## Requisitos

- Um ou mais harnesses instalados: [OpenCode](https://opencode.ai), [Claude Code](https://claude.ai/code) e/ou [Codex](https://github.com/openai/codex)
- `curl` ou `wget` (para instalação remota)
- `bash` >= 4.0
- MCP `kanban-force` configurado (necessário para operações de board/card com o agente `kanban`)

## Modelos configurados

Os modelos são definidos durante a instalação conforme o harness e o vendor escolhido. Apenas os agentes **primários** recebem `model` no frontmatter. Os **subagentes** não têm `model` definido e herdam o modelo do agente que os aciona.

**Agentes primários:** `coder`, `lead`, `documenter`, `kanban`, `infra`, `mr_reviewer`

**Subagentes (herdam):** `analyzer`, `clarifier`, `planner`, `detailer`, `tester`, `code_reviewer`, `business_reviewer`, `versioner`

| Harness | Vendor | Modelo dos primários |
|---|---|---|
| OpenCode | `anthropic` | `anthropic/claude-sonnet-4-6` |
| OpenCode | `openai` (padrão) | `openai/gpt-5.5` |
| OpenCode | `google` | `google/gemini-2.5-pro` |
| OpenCode | `groq` | `groq/llama-3.3-70b-versatile` |
| OpenCode | `amazon-bedrock` | `amazon-bedrock/amazon.nova-pro-v1:0` |
| OpenCode | `github-copilot` | `github-copilot/claude-sonnet-4.6` |
| Claude Code | — | `sonnet` |
| Codex | — | herdado da sessão |
