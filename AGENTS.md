# AGENTS.md

## Sobre o repositorio

Markdown-only: definicoes de agentes e skills para OpenCode, Claude Code, Codex e Pi. Nao ha codigo executavel, dependencias, build ou testes. O unico script e `install.sh`.

## Estrutura

```text
agents/    14 subdiretorios — cada agente tem body.md + opencode.yml + claude.yml
skills/    14 subdiretorios — cada skill tem SKILL.md + references/ (opcional)
commands/  5 subdiretorios  — cada command tem body.md + opencode.yml + claude.yml + pi.yml
install.sh                  — monta e copia tudo para ~/.config/opencode/
```

## Convencoes obrigatorias

- **Idioma:** todo conteudo em portugues do Brasil
- **Frontmatter YAML** separado por harness — cada harness usa os campos que entende:
  - `opencode.yml` (OpenCode usa `mode` e `temperature`):
    - Primario: `description`, `mode: primary`, `model` (token `__OPENCODE_MAIN__`, resolvido na instalacao), `temperature`
    - Subagente: `description`, `mode: subagent`, `temperature` — sem `model`
  - `claude.yml` (Claude Code **nao** usa `mode` nem `temperature`):
    - Primario: `name`, `description`, `model: sonnet`
    - Subagente: `name`, `description` — sem `model`
  - Skills (`SKILL.md`): `name` (kebab-case), `description`
- **XML tags** para estruturar conteudo: `<role>`, `<responsibilities>`, `<rules>`, `<workflow>`, `<instructions>`, `<output_format>`, `<checklist>`, `<principles>`, `<criteria>`, `<context>`, `<code_navigation>`
- **Listas de verificacao** (`- [ ]`) nas skills de revisao e nos agentes reviewer
- **Formato de saida** (`<output_format>`) no final de cada arquivo define o contrato de resposta
- **Bloco de diff estruturado** nos reviewers: toda sugestao usa o formato `path > linha > atual > sugerido > motivo`

## Relacao agentes ↔ skills

| Agente | Mode | Skill | Obs |
|---|---|---|---|
| `coder` | primary | `write-code` | Orquestrador — aciona todos os subagentes |
| `lead` | primary | `plan-implementation` | Orquestrador de planejamento |
| `documenter` | primary | `document-plan`, `get-plan` | Publica planos no Confluence |
| `kanban` | primary | `kanban-force` | Gerencia cards e boards |
| `infra` | primary | `query-argocd` | Consulta ArgoCD |
| `mr_reviewer` | primary | `review-mr` | Revisa MRs via glab |
| `analyzer` | subagent | `analyse-code` | Inspeciona codebase |
| `clarifier` | subagent | `clarify-intent` | Formata perguntas de ambiguidade |
| `planner` | subagent | `plan-tasks` | TaskGraph esqueleto |
| `detailer` | subagent | `detail-tasks` | Enriquece tasks |
| `tester` | subagent | `test-code` | TDD: testes antes e depois da implementacao |
| `code_reviewer` | subagent | `review-code` | Camada 1 — tecnica |
| `business_reviewer` | subagent | `review-code` | Camada 2 — negocio/seguranca (mesma skill, papel diferente) |
| `versioner` | subagent | `version-code` | Operacoes Git; herda modelo do chamador |

`review-code` e compartilhada: o agente identifica se e Camada 1 ou 2 pelo papel que o acionou.

`kanban` e um agente primary independente (standalone). Depende do MCP `kanban-force` para todas as operacoes.

Integracao `coder` -> `kanban`: quando a solicitacao tiver ID de card (ex.: `STK-90AB`, `UST-FF51`) ou pedir operacao de board/card (criar, mover, atualizar, comentar, bloquear, arquivar etc.), o `coder` deve delegar ao `kanban`.

## install.sh — como funciona

- Instala para um ou mais harnesses escolhidos antes da instalacao: OpenCode, Claude Code e/ou Codex (menu interativo ou flag `--harness`)
- Monta cada agente/command juntando `<harness>.yml` (`opencode.yml` ou `claude.yml`) + `body.md` e copia para o diretorio nativo do harness:
  - OpenCode: `$OPENCODE_DIR` (default `~/.config/opencode/`) → `agents/`, `skills/`, `commands/`
  - Claude Code: `$CLAUDE_DIR` (default `~/.claude/`) → `agents/`, `skills/`, `commands/`
  - Codex: skills em `$CODEX_SKILLS_DIR` (default `~/.agents/skills/`); commands viram prompts body-only em `~/.codex/prompts/`; `AGENTS.md` copiado para `~/.codex/`; sem agentes nativos
  - Pi: skills em `$PI_SKILLS_DIR` (default `~/.pi/agent/skills/`); commands viram prompts montados (`pi.yml` + body) em `~/.pi/agent/prompts/`; `AGENTS.md` copiado para `~/.pi/agent/`; sem agentes nativos; modelo via settings/provider do Pi
- Copia skills como diretorios completos (`SKILL.md`) — identicas nos tres harnesses
- Modelo: so agentes **primarios** recebem `model`. OpenCode → `openai/gpt-5.5` (default) ou `vendor/main`; Claude Code → `sonnet`; Codex herda. Subagentes nao recebem `model`
- Flags: `--harness <lista>`, `--vendor <nome>`, `--force` (sobrescreve sem perguntar), `--local` (usa arquivos locais em vez de baixar do GitHub)

## Commits

Conventional Commits em ingles: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`, `style:`, `perf:`, `test:`

## Commands (`commands/`)

| Diretorio | Comando | Descricao |
|---|---|---|
| `doc-plan/` | `/doc-plan` | Publica `.coder/plan.md` no Confluence (space CAT, raiz Implementacoes) via MCP `atlassian_local`; ignora se sem diferencas |
| `get-plan/` | `/get-plan` | Baixa o plano do Confluence e salva em `.coder/plan.md`; cria o arquivo se nao existir |
| `kanban-card/` | `/kanban-card <friendlyID>` | Consulta um card pelo friendlyID via MCP kanban-force e carrega no contexto |
| `mr-review/` | `/mr-review` | Aciona o `mr_reviewer` para revisar o MR aberto na branch atual via `glab` |

Cada command e um subdiretorio com `body.md` (instrucoes) + frontmatter por harness (`opencode.yml`, `claude.yml` e `pi.yml`). O nome do diretorio vira o slash command. Argumentos sao acessados via `$ARGUMENTS` (todos) ou `$1`, `$2`... (posicionais).

Instalados em `$OPENCODE_DIR/commands/` (padrao: `~/.config/opencode/commands/`).

## O que nao fazer

- Nao adicionar codigo executavel, dependencias ou config de build — o repo e apenas Markdown + 1 shell script
- Nao alterar a estrutura de XML tags sem verificar todos os arquivos que a usam
- Nao remover `<output_format>` de nenhum arquivo — e o contrato de resposta do agente
- Nao mudar o valor de `model:` manualmente — `install.sh` sobrescreve na instalacao
