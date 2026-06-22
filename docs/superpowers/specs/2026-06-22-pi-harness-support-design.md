# Design — Suporte ao harness Pi

- **Data:** 2026-06-22
- **Status:** Aprovado (aguardando plano de implementação)
- **Autor:** Paraizo

## Solicitação original

> Quero que adicione suporte ao harness Pi, para que ele também consiga utilizar os agents e skills.

## Contexto

O repositório `coder` é uma coleção **somente Markdown + 1 shell script** (`install.sh`)
que distribui agentes, skills e commands para harnesses de IA. Hoje suporta três:
**OpenCode**, **Claude Code** e **Codex**. O layout é híbrido: o corpo dos agentes e
commands é único e agnóstico de harness; só o frontmatter varia por harness
(`opencode.yml`, `claude.yml`). As skills seguem o padrão aberto
[Agent Skills](https://agentskills.io/specification) e são idênticas nos três harnesses.

O `install.sh` monta cada agente/command juntando `<harness>.yml` + `body.md` e copia
para o diretório nativo de cada harness, aplicando o modelo do vendor escolhido.

## Investigação do harness Pi (`@earendil-works/pi-coding-agent`, v0.79.10)

Pi foi inspecionado localmente (CLI instalada, pacote em `node_modules`, docs bundladas).
Conclusão: em capacidade, **Pi é quase idêntico ao Codex** — skills nativas + prompt
templates + context file `AGENTS.md`, **sem conceito nativo de agente/subagente**
(há apenas o pacote opcional `pi-subagents`, fora de escopo).

Diretórios e formatos nativos do Pi:

| Recurso | Caminho nativo (global) | Formato | Equivale a |
|---|---|---|---|
| **Skills** | `~/.pi/agent/skills/<nome>/SKILL.md` | Agent Skills standard | Codex `~/.agents/skills` |
| **Commands** | `~/.pi/agent/prompts/*.md` (invocados por `/nome`) | Markdown + frontmatter opcional (`description`, `argument-hint`); suporta `$ARGUMENTS`/`$1` | Codex prompts |
| **Orquestração** | `~/.pi/agent/AGENTS.md` | Context file global | Codex `~/.codex/AGENTS.md` |
| **Modelo** | settings/provider do Pi | — | Codex (herdado) |

Fatos confirmados na investigação:

- Pi implementa o padrão Agent Skills e carrega skills de `~/.pi/agent/skills/` (também
  de `~/.agents/skills/`). Cópia direta das skills do repo, sem edição.
- Os `body.md` dos commands **já usam `$ARGUMENTS`**, que é nativo nos prompt templates
  do Pi — compatível sem alteração.
- Pi carrega `AGENTS.md` global de `~/.pi/agent/AGENTS.md` no startup (confirmado em
  `docs/usage.md` do pacote).
- Pi não tem frontmatter de agente; a orquestração dos agentes do repo (coder, lead,
  qa, etc.) é entregue ao Pi pelo `AGENTS.md` — exatamente como já é feito para o Codex.

## Decisões tomadas

1. **Representação dos agentes no Pi:** paridade com Codex. Orquestração via `AGENTS.md`
   global + skills nativas + commands como prompts. Sem expor agentes primários como
   prompt templates (`/coder`, `/lead`…).
2. **Frontmatter dos commands no Pi:** adicionar um `pi.yml` por command
   (`description` + `argument-hint`), montado como `opencode.yml`/`claude.yml`. O Pi
   exibe descrição e dica de argumento no autocomplete do `/`.

## Design da solução

### Arquitetura

Pi entra como **4º harness**, espelhando o Codex. O corpo dos agentes/commands continua
agnóstico; muda o destino e, nos commands, um novo frontmatter `pi.yml`.

Destinos de instalação (com overrides por env var):

```
~/.pi/agent/skills/<nome>/SKILL.md   skills nativas (cópia direta, igual Codex)
~/.pi/agent/prompts/<nome>.md        commands montados (pi.yml + body.md)
~/.pi/agent/AGENTS.md                orquestração global (mesmo arquivo do Codex)
```

- `apply_model` para `pi` → **no-op** (Pi define modelo via settings/provider, não por agente).
- Sem `H_AGENTS` para Pi (não há agente nativo) — orquestração vai pelo `AGENTS.md`.

### Componente 1 — Novos arquivos no repo (5 `pi.yml`)

Um `pi.yml` por command, em `commands/<nome>/pi.yml`. `description` reaproveita o texto
já presente nos `*.yml` existentes; `argument-hint` é adicionado onde há argumento:

| Command | `description` (fonte) | `argument-hint` |
|---|---|---|
| `doc-plan` | reusa existente | — |
| `get-plan` | reusa existente | — |
| `kanban-card` | reusa existente | `"<friendlyID>"` |
| `mr-review` | reusa existente | `"<MR# \| URL>"` |
| `qa` | reusa existente | `"[foco/escopo]"` |

Formato (exemplo `commands/qa/pi.yml`):

```yaml
description: "Valida funcionalmente as modificações da branch atual ..."
argument-hint: "[foco/escopo]"
```

### Componente 2 — Mudanças no `install.sh`

- **Seleção de harness:** adicionar `pi` ao menu interativo (opção 5) e à opção `todos`;
  aceitar `--harness pi` e incluí-lo em `all`. Ordem canônica passa a ser
  `opencode → claude → codex → pi` (dedup e `select_harness`/`resolve_harness_flag`).
- **`harness_paths`:** novo case `pi)`:
  - `H_SKILLS = ${PI_SKILLS_DIR:-${PI_DIR:-$HOME/.pi/agent}/skills}`
  - `H_AGENTS = ""` (sem agentes nativos)
  - `H_COMMANDS = ""` (commands vão para prompts)
  - `H_PROMPTS = ${PI_DIR:-$HOME/.pi/agent}/prompts`
  - `H_AGENTSMD = ${PI_DIR:-$HOME/.pi/agent}/AGENTS.md`
- **`apply_model`:** case `pi)` → no-op (sem token de modelo no `pi.yml`).
- **Instalação dos commands do Pi:** diferente do Codex (body-only via `cp`), o Pi monta
  `pi.yml` + `body.md`. Reutilizar `assemble`/`prepare_assembled_src` com `harness=pi`,
  gravando em `H_PROMPTS`. Implementar como função dedicada `install_pi_prompts`
  (ou generalizar o install de prompts para distinguir "assembled" de "body-only").
- **AGENTS.md do Pi:** generalizar `install_codex_agentsmd` para gravar em `H_AGENTSMD`
  (independente do harness). Fetch remoto **não-fatal** (AGENTS.md é gitignored, igual Codex).
- **Skills:** sem mudança — `install_skills` já é agnóstico, só muda `H_SKILLS`.
- **Loop principal:** para `pi`, chamar `install_skills`, `install_pi_prompts` e o install
  do `AGENTS.md`.
- **`--help` e `print_summary`:** incluir `pi` nos valores aceitos, documentar `PI_DIR` e
  `PI_SKILLS_DIR`, e exibir os caminhos/modelo do Pi no resumo final.

### Componente 3 — Documentação

- `AGENTS.md` — mencionar Pi (estrutura, `pi.yml`, seção do `install.sh`, caminhos).
- `README.md` — Pi nos harnesses, diretórios de instalação, estrutura do repo, requisitos.
- `CLAUDE.md` (projeto) — estrutura com `pi.yml` e o novo destino Pi no resumo do `install.sh`.

## Fluxo de instalação (Pi)

```
install.sh --harness pi [--local]
 ├─ harness_paths pi        → define H_SKILLS/H_PROMPTS/H_AGENTSMD
 ├─ install_skills          → ~/.pi/agent/skills/ (cópia direta dos SKILL.md)
 ├─ install_pi_prompts      → ~/.pi/agent/prompts/ (assemble pi.yml + body.md)
 └─ install_agentsmd (pi)   → ~/.pi/agent/AGENTS.md (cópia/fetch não-fatal)
```

## Tratamento de erros e bordas

- **Modo remoto:** `prepare_assembled_src` baixa `pi.yml` + `body.md` do GitHub; funciona
  porque os `pi.yml` passam a existir no repo. `AGENTS.md` retorna 404 no remoto
  (gitignored) → warn não-fatal, igual ao Codex.
- **Diretórios ausentes:** `mkdir -p` cria `~/.pi/agent/{skills,prompts}` se necessário.
- **Conflitos:** `check_overwrite` já cobre prompts/skills/AGENTS.md (confirmação ou `--force`).
- **`git add -f`:** o `.gitignore` global bloqueia `docs/`, `AGENTS.md` e `CLAUDE.md`;
  versionar esses arquivos exige `git add -f` (relevante para os docs e o AGENTS.md).

## Estratégia de testes

Como o repo não tem suíte automatizada, a validação é manual:

1. `bash install.sh --local --harness pi --force` em ambiente de teste
   (`PI_DIR=$(mktemp -d)/pi/agent`).
2. Conferir que skills, prompts e `AGENTS.md` aparecem nos caminhos esperados.
3. `bash install.sh --local --harness all --force` — confirmar que os quatro harnesses
   instalam sem erro e a ordem canônica inclui `pi` por último.
4. Smoke no Pi real: iniciar `pi`, verificar no header de startup que skills e prompt
   templates do repo são listados, e que `/qa`, `/mr-review` etc. aparecem no autocomplete
   com descrição/argument-hint.

## Critérios de aceite

- [ ] `--harness pi` e `--harness all` instalam Pi nos caminhos corretos.
- [ ] Menu interativo lista `pi` (opção 5) e `todos` inclui Pi.
- [ ] As 15 skills são copiadas para `~/.pi/agent/skills/`.
- [ ] Os 5 commands viram prompts montados (pi.yml + body) em `~/.pi/agent/prompts/`.
- [ ] `AGENTS.md` é instalado em `~/.pi/agent/AGENTS.md` (local), com fetch remoto não-fatal.
- [ ] `apply_model` não altera nada para Pi.
- [ ] `--help`, `print_summary`, `AGENTS.md`, `README.md` e `CLAUDE.md` mencionam Pi.
- [ ] Harnesses existentes (opencode, claude, codex) seguem funcionando sem regressão.

## Fora de escopo (YAGNI)

- Extensão TypeScript (`~/.pi/agent/extensions/`) — violaria a regra "só Markdown + 1 shell script".
- Agentes expostos como prompt templates (`/coder`, `/lead`…) — descartado (paridade com Codex).
- Pacote `pi-subagents` e qualquer dependência executável.
- Modelo por agente no Pi — Pi resolve modelo via settings/provider.
