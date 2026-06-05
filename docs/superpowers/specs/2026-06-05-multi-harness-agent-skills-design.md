# Design — Migração para Agent Skills + instalação multi-harness

- **Data:** 2026-06-05
- **Repositório:** `coder` (coleção de agentes/skills/commands em Markdown + `install.sh`)
- **Objetivo:** Converter o projeto para o standard [agentskills.io](https://agentskills.io) e permitir instalar para múltiplos harnesses (OpenCode, Claude Code, Codex), escolhendo o harness antes de instalar.

## Contexto

Hoje o repositório tem:

- `agents/*.md` — 14 agentes (conceito de harness), frontmatter OpenCode (`description`, `mode`, `model`, `temperature`).
- `skills/*.md` — 14 skills **planas** com nomes em snake_case (`analyse_code`), tendo apenas `description` no frontmatter.
- `commands/*.md` — 4 slash commands.
- `install.sh` — instala em `~/.opencode/` e troca `model:` por vendor (menu de 6 vendors, tiers `main`/`light`).

Dois desalinhamentos com o standard:

1. O standard agentskills.io define **apenas skills** (pasta + `SKILL.md` com `name` + `description`). Skills planas e em snake_case são **inválidas** (o `name` só aceita `[a-z0-9-]`, sem `_`).
2. "Agentes" não existem no standard — são um conceito de cada harness, com frontmatter e diretórios diferentes entre OpenCode, Claude Code e Codex.

## Decisões tomadas (brainstorming)

| # | Decisão | Escolha |
|---|---------|---------|
| 1 | Como representar os 14 agentes | **Skills = núcleo portável; agentes nativos por harness.** Skills viram `SKILL.md` (iguais nos 3 harnesses); agentes têm frontmatter nativo por harness. |
| 2 | Layout da fonte no repo | **Híbrido.** Skills e o **corpo** dos agentes em fonte única; só o **frontmatter** do agente é separado por harness. |
| 3 | Resolução de modelo/vendor | **Ciente do harness.** Cada harness recebe o formato de modelo que entende. |
| 4 | Regra de modelo (OpenCode e Claude Code) | **Só agentes principais recebem modelo; subagentes herdam** (sem campo `model`). |

### Achados externos relevantes

- **Codex passou a suportar subagentes** (custom agents com modelo próprio, ligados por padrão). Logo, não fica restrito a "só skills + AGENTS.md" — pode receber agentes nativos numa fase posterior.
- Diretórios reais (global):
  - OpenCode: `~/.config/opencode/skills/<n>/SKILL.md`, `~/.config/opencode/agent/`, `~/.config/opencode/command/` (também lê `~/.claude/skills/` e `~/.agents/skills/` por compatibilidade).
  - Claude Code: `~/.claude/skills/`, `~/.claude/agents/`, `~/.claude/commands/`.
  - Codex: skills em `~/.agents/skills/`; prompts em `~/.codex/prompts/`.

## Estrutura do repositório (padrão de diretórios escolhido)

```text
skills/
  <skill-name>/
    SKILL.md            # name (kebab) + description + corpo (conforme spec)
    references/         # opcional, quando houver material extenso
agents/
  <agent-name>/
    body.md             # corpo compartilhado (XML tags), SEM frontmatter
    opencode.yml        # frontmatter OpenCode
    claude.yml          # frontmatter Claude Code
    # codex.yml         # opcional — Fase 2 (subagentes nativos do Codex)
commands/
  <command-name>/
    body.md             # corpo compartilhado
    opencode.yml        # frontmatter OpenCode
    claude.yml          # frontmatter Claude Code
    # Codex → emitido como prompt body-only em ~/.codex/prompts/<n>.md
AGENTS.md               # doc de orquestração — também serve ao Codex
install.sh
README.md
CLAUDE.md
```

Regras do layout:

- O **corpo** de cada agente fica num único `body.md` (sem duplicação de conteúdo).
- O frontmatter por harness é um arquivo `.yml` curto (poucas linhas), escrito à mão.
- O campo de modelo usa **placeholder** (`__OPENCODE_MAIN__`) trocado por `sed` na instalação — **sem parsing de YAML em bash**.

## Conformidade com o standard agentskills.io (Fase 1)

- Cada skill vira **pasta com `SKILL.md`**, frontmatter com `name` (obrigatório, **kebab-case**, igual ao nome da pasta) + `description`.
- Renomeações snake_case → kebab-case (14 skills):

  | Atual | Novo |
  |---|---|
  | `analyse_code` | `analyse-code` |
  | `clarify_intent` | `clarify-intent` |
  | `detail_tasks` | `detail-tasks` |
  | `document_plan` | `document-plan` |
  | `get_plan` | `get-plan` |
  | `kanban_force` | `kanban-force` |
  | `plan_implementation` | `plan-implementation` |
  | `plan_tasks` | `plan-tasks` |
  | `query_argocd` | `query-argocd` |
  | `review_code` | `review-code` |
  | `review_mr` | `review-mr` |
  | `test_code` | `test-code` |
  | `version_code` | `version-code` |
  | `write_code` | `write-code` |

- **Todas as referências cruzadas** atualizadas para os novos nomes: corpos dos agentes (que citam skills por nome), `install.sh`, `README.md`, `CLAUDE.md`, `AGENTS.md` e os commands.
- Validação obrigatória por skill: `npx -y skills-ref validate ./skills/<nome>`.

## Mapa de instalação por harness

| Artefato | OpenCode | Claude Code | Codex |
|---|---|---|---|
| Skills | `~/.config/opencode/skills/<n>/SKILL.md` | `~/.claude/skills/<n>/SKILL.md` | `~/.agents/skills/<n>/SKILL.md` |
| Agentes | `~/.config/opencode/agents/<n>.md` (`opencode.yml` + `body.md`) | `~/.claude/agents/<n>.md` (`claude.yml` + `body.md`) | Fase 2 (opcional); na Fase 1, só `AGENTS.md` |
| Commands | `~/.config/opencode/commands/<n>.md` | `~/.claude/commands/<n>.md` | `~/.codex/prompts/<n>.md` (body-only) |
| Override de dir | `OPENCODE_DIR` | `CLAUDE_DIR` | `CODEX_DIR` |

## Resolução de modelo

Aplica-se a regra **"só principais recebem modelo; subagentes herdam"**.

Agentes **principais** (recebem modelo): `coder`, `lead`, `documenter`, `kanban`, `infra`, `mr_reviewer`.
Agentes **subagentes** (sem `model`): `analyzer`, `clarifier`, `planner`, `detailer`, `tester`, `code_reviewer`, `business_reviewer`, `versioner`.

| Harness | Principais | Subagentes |
|---|---|---|
| OpenCode | `openai/gpt-5.5` (default, sem vendor) **ou** `vendor/main` (vendor escolhido) | sem campo `model` |
| Claude Code | alias `sonnet` (default, editável) | sem campo `model` (herda) |
| Codex | sem patch (herda da sessão) | sem patch |

Consequências:

- O tier **`light`/`__LIGHT__` é eliminado** — só era usado no `versioner`, que é subagente e agora não recebe modelo.
- O menu de vendor do OpenCode fica reduzido à coluna **`main`** (a coluna `light` some).
- A **seleção de vendor passa a ser opcional**: sem escolha, principais usam `openai/gpt-5.5`.
- `agents/<primary>/opencode.yml` contém `model: __OPENCODE_MAIN__`; `agents/<subagent>/opencode.yml` **omite** a linha `model`.
- `agents/<primary>/claude.yml` contém `model: sonnet`; `agents/<subagent>/claude.yml` **omite** `model`.

## Novo fluxo do `install.sh` (Fase 2)

1. **Escolher harness(es)** — menu permitindo um ou vários (`opencode` / `claude` / `codex`).
2. **Se OpenCode** estiver entre os escolhidos → oferecer menu de vendor (opcional; pular = `openai/gpt-5.5`).
3. Para cada harness escolhido:
   - **Agentes**: monta `frontmatter(<harness>.yml) + body.md`, troca placeholder de modelo conforme a regra do harness, grava no diretório nativo. (Codex: pula na Fase 1.)
   - **Skills**: copia a pasta `SKILL.md` inteira para o diretório de skills do harness.
   - **Commands**: copia para `commands/` (OC) / `commands/` (CC) / `prompts/` (Codex, body-only).
   - **Codex**: além das skills/prompts, instala/aponta o `AGENTS.md` de orquestração.
4. Mantém `--force`, `--local`, `--help` e o gate de confirmação por arquivo já existente.
5. Resumo final por harness instalado (diretórios, modelo aplicado, artefatos).

## Fases de entrega

### Fase 1 — Conformidade com o standard

- Reestruturar as 14 skills em pastas `SKILL.md` (kebab-case + `name`).
- Reorganizar os 14 agentes em `agents/<n>/{body.md, opencode.yml, claude.yml}`.
- Reorganizar os 4 commands em `commands/<n>/{body.md, opencode.yml, claude.yml}`.
- Atualizar todas as referências cruzadas e a documentação (`README`, `CLAUDE.md`, `AGENTS.md`).
- Ajustar o `install.sh` para ler o novo layout mantendo a instalação OpenCode funcionando.
- Validar cada skill com `skills-ref`.

### Fase 2 — Instalação multi-harness

- Reescrever o `install.sh` com seleção de harness, diretórios nativos, montagem de frontmatter por harness e resolução de modelo ciente do harness (incluindo `openai/gpt-5.5` default e vendor opcional).
- Suporte Codex: skills (`~/.agents/skills/`), prompts (`~/.codex/prompts/`) e `AGENTS.md`.
- Opcional: `agents/<n>/codex.yml` para subagentes nativos do Codex (após confirmar o formato).

## Pontos a confirmar durante o plano (não bloqueiam o desenho)

- Formato exato do arquivo de **custom agent do Codex** antes de gerar `codex.yml`.
- ✅ Resolvido: OpenCode usa subdiretórios **no plural** (`agents/`, `commands/`, `skills/`) sob `~/.config/opencode/`.
- Confirmar o id do provider/modelo `gpt-5.5` no formato OpenCode (`openai/gpt-5.5`).
- Frontmatter de command no Claude Code (`argument-hint`, `allowed-tools`) vs OpenCode (`agent`, `model`).

## Fora de escopo (YAGNI)

- Não adicionar código executável ou dependências de build além do `install.sh` (mantém a filosofia "Markdown + 1 shell script").
- Não criar gerador/transformador que parseie YAML em bash (motivo da escolha pelo layout híbrido).
- Não suportar harnesses além de OpenCode, Claude Code e Codex nesta entrega.

## Riscos

- **Sincronização de frontmatter por harness**: cada agente tem 2 (ou 3) arquivos `.yml` curtos; risco de divergência de `description`. Mitigação: arquivos minúsculos e estáveis; checklist de revisão.
- **Diretórios dos harnesses podem mudar** entre versões; mitigar com overrides por env var e confirmação na doc.
- **Quebra de referências** ao renomear skills snake→kebab; mitigar com varredura completa de todas as ocorrências antes do commit.
