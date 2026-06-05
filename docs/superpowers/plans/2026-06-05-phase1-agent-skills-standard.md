# Fase 1 — Conformidade com o padrão Agent Skills — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reestruturar o repositório para que as 14 skills sigam o standard agentskills.io (pasta `SKILL.md`, `name` kebab-case) e os agentes/commands adotem o layout híbrido (corpo único + frontmatter por harness), mantendo a instalação OpenCode funcionando.

**Architecture:** Skills viram pastas `<name>/SKILL.md`. Agentes e commands viram `<name>/{body.md, opencode.yml, claude.yml}` — corpo agnóstico + frontmatter nativo por harness. O `install.sh` é ajustado para montar agente = `opencode.yml` + `body.md` e copiar skills como pastas (ainda só OpenCode nesta fase).

**Tech Stack:** Markdown, YAML frontmatter, Bash (`install.sh`), `skills-ref` (validação via `npx`).

---

## Estrutura de arquivos (antes → depois)

```text
ANTES                          DEPOIS
skills/write_code.md      →    skills/write-code/SKILL.md
agents/coder.md           →    agents/coder/{body.md, opencode.yml, claude.yml}
commands/mr_review.md     →    commands/mr-review/{body.md, opencode.yml, claude.yml}
```

**Renomeações snake_case → kebab-case (skills):**

| Atual | Novo | Atual | Novo |
|---|---|---|---|
| `analyse_code` | `analyse-code` | `plan_tasks` | `plan-tasks` |
| `clarify_intent` | `clarify-intent` | `query_argocd` | `query-argocd` |
| `detail_tasks` | `detail-tasks` | `review_code` | `review-code` |
| `document_plan` | `document-plan` | `review_mr` | `review-mr` |
| `get_plan` | `get-plan` | `test_code` | `test-code` |
| `kanban_force` | `kanban-force` | `version_code` | `version-code` |
| `plan_implementation` | `plan-implementation` | `write_code` | `write-code` |

**Renomeações de commands (snake → kebab):** `doc_plan→doc-plan`, `get_plan→get-plan`, `kanban_card→kanban-card`, `mr_review→mr-review`.

**Classificação de agentes:**

- **Primários** (recebem modelo): `coder`, `lead`, `documenter`, `kanban`, `infra`, `mr_reviewer`.
- **Subagentes** (sem `model`): `analyzer`, `clarifier`, `planner`, `detailer`, `tester`, `code_reviewer`, `business_reviewer`, `versioner`.

**Diferenças de frontmatter:**

| Campo | `opencode.yml` | `claude.yml` |
|---|---|---|
| `description` | sim | sim |
| `name` | não | sim (igual ao nome da pasta) |
| `mode` | sim (`primary`/`subagent`) | não |
| `temperature` | sim | não (CC não suporta) |
| `model` | só primários (`__OPENCODE_MAIN__`) | só primários (`sonnet`) |

---

## Task 1: Converter as 14 skills em pastas `SKILL.md` com `name` kebab-case

**Files:**
- Create: `skills/<kebab-name>/SKILL.md` (14 pastas)
- Delete: `skills/<snake_name>.md` (14 arquivos planos)

- [ ] **Step 1: Criar o script de conversão temporário**

Create `/tmp/convert_skills.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

declare -A MAP=(
  [analyse_code]=analyse-code [clarify_intent]=clarify-intent
  [detail_tasks]=detail-tasks [document_plan]=document-plan
  [get_plan]=get-plan [kanban_force]=kanban-force
  [plan_implementation]=plan-implementation [plan_tasks]=plan-tasks
  [query_argocd]=query-argocd [review_code]=review-code
  [review_mr]=review-mr [test_code]=test-code
  [version_code]=version-code [write_code]=write-code
)

for snake in "${!MAP[@]}"; do
  kebab="${MAP[$snake]}"
  src="skills/${snake}.md"
  dir="skills/${kebab}"
  mkdir -p "$dir"
  # injeta a linha `name:` logo após o primeiro `---`
  awk -v name="$kebab" '
    NR==1 && $0=="---" { print; print "name: " name; next }
    { print }
  ' "$src" > "$dir/SKILL.md"
  git rm -q "$src"
  git add "$dir/SKILL.md"
  echo "convertido: $snake -> $dir/SKILL.md"
done
```

- [ ] **Step 2: Executar o script**

Run: `bash /tmp/convert_skills.sh`
Expected: 14 linhas `convertido: ...`. `ls skills/` mostra 14 diretórios e nenhum `.md` solto.

- [ ] **Step 3: Verificar o frontmatter de uma skill convertida**

Run: `sed -n '1,4p' skills/write-code/SKILL.md`
Expected:
```text
---
name: write-code
description: Skill principal do agente coder. ...
---
```

- [ ] **Step 4 (teste): Validar todas as skills com skills-ref**

Run:
```bash
for d in skills/*/; do npx -y skills-ref validate "./$d" || echo "FALHOU: $d"; done
```
Expected: nenhuma linha `FALHOU:`. Cada skill reportada como válida (`name` kebab + `description` presentes).

- [ ] **Step 5: Commit**

```bash
git add skills/
git commit -m "refactor: convert skills to Agent Skills folder layout with kebab-case names"
```

---

## Task 2: Reestruturar os 14 agentes em `agents/<n>/{body.md, opencode.yml, claude.yml}`

**Files:**
- Create: `agents/<name>/body.md`, `agents/<name>/opencode.yml`, `agents/<name>/claude.yml` (14 agentes)
- Delete: `agents/<name>.md` (14 arquivos)

- [ ] **Step 1: Criar o script de split de agentes**

Create `/tmp/split_agents.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

PRIMARY="coder lead documenter kanban infra mr_reviewer"

is_primary() { case " $PRIMARY " in *" $1 "*) return 0;; *) return 1;; esac; }

for src in agents/*.md; do
  name="$(basename "$src" .md)"
  dir="agents/$name"
  mkdir -p "$dir"

  # corpo = tudo após o segundo `---`
  awk '/^---$/{c++; next} c>=2{print}' "$src" > "$dir/body.md"
  # remove linha em branco inicial deixada pelo frontmatter
  sed -i '1{/^$/d}' "$dir/body.md"

  desc="$(awk -F': ' '/^description: /{sub(/^description: /,""); print; exit}' "$src")"
  temp="$(awk -F': ' '/^temperature: /{print $2; exit}' "$src")"
  mode="$(awk -F': ' '/^mode: /{print $2; exit}' "$src")"

  # opencode.yml
  {
    echo "description: $desc"
    echo "mode: $mode"
    if is_primary "$name"; then echo "model: __OPENCODE_MAIN__"; fi
    echo "temperature: $temp"
  } > "$dir/opencode.yml"

  # claude.yml
  {
    echo "name: $name"
    echo "description: $desc"
    if is_primary "$name"; then echo "model: sonnet"; fi
  } > "$dir/claude.yml"

  git rm -q "$src"
  git add "$dir/body.md" "$dir/opencode.yml" "$dir/claude.yml"
  echo "split: $name (primary=$(is_primary "$name" && echo sim || echo nao))"
done
```

- [ ] **Step 2: Executar o script**

Run: `bash /tmp/split_agents.sh`
Expected: 14 linhas `split: ...`, com `primary=sim` para coder, lead, documenter, kanban, infra, mr_reviewer e `primary=nao` para os demais.

- [ ] **Step 3: Verificar um agente PRIMÁRIO (coder)**

Run: `cat agents/coder/opencode.yml agents/coder/claude.yml`
Expected:
```text
description: Agente principal orquestrador de desenvolvimento de software. ...
mode: primary
model: __OPENCODE_MAIN__
temperature: 0.3
name: coder
description: Agente principal orquestrador de desenvolvimento de software. ...
model: sonnet
```

- [ ] **Step 4: Verificar um SUBAGENTE (versioner) — não deve ter `model`**

Run: `cat agents/versioner/opencode.yml agents/versioner/claude.yml`
Expected (sem nenhuma linha `model:`):
```text
description: Subagente especializado em versionamento Git. ...
mode: subagent
temperature: 0.2
name: versioner
description: Subagente especializado em versionamento Git. ...
```

- [ ] **Step 5 (teste): Garantir que o corpo começa com a primeira tag XML**

Run: `head -1 agents/coder/body.md`
Expected: `<role>` (sem frontmatter residual, sem linha em branco no topo).

- [ ] **Step 6 (teste): Garantir que nenhum subagente recebeu `model`**

Run:
```bash
for n in analyzer clarifier planner detailer tester code_reviewer business_reviewer versioner; do
  grep -q '^model:' "agents/$n/opencode.yml" && echo "ERRO model em $n"
done; echo "ok"
```
Expected: apenas `ok`.

- [ ] **Step 7: Commit**

```bash
git add agents/
git commit -m "refactor: split agents into shared body plus per-harness frontmatter"
```

---

## Task 3: Reestruturar os 4 commands em `commands/<n>/{body.md, opencode.yml, claude.yml}`

**Files:**
- Create: `commands/<kebab>/body.md`, `commands/<kebab>/opencode.yml`, `commands/<kebab>/claude.yml`
- Delete: `commands/<snake>.md`

- [ ] **Step 1: Criar o script de split de commands**

Create `/tmp/split_commands.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

declare -A MAP=( [doc_plan]=doc-plan [get_plan]=get-plan [kanban_card]=kanban-card [mr_review]=mr-review )

for snake in "${!MAP[@]}"; do
  kebab="${MAP[$snake]}"
  src="commands/${snake}.md"
  dir="commands/${kebab}"
  mkdir -p "$dir"

  awk '/^---$/{c++; next} c>=2{print}' "$src" > "$dir/body.md"
  sed -i '1{/^$/d}' "$dir/body.md"

  desc="$(awk -F': ' '/^description: /{sub(/^description: /,""); print; exit}' "$src")"
  agent="$(awk -F': ' '/^agent: /{print $2; exit}' "$src")"

  # opencode.yml: description + agent
  { echo "description: $desc"; echo "agent: $agent"; } > "$dir/opencode.yml"
  # claude.yml: apenas description (CC não vincula subagente via frontmatter)
  { echo "description: $desc"; } > "$dir/claude.yml"

  git rm -q "$src"
  git add "$dir"/{body.md,opencode.yml,claude.yml}
  echo "command: $snake -> $kebab (agent=$agent)"
done
```

- [ ] **Step 2: Executar o script**

Run: `bash /tmp/split_commands.sh`
Expected: 4 linhas `command: ...` com `agent=documenter` (x2), `agent=kanban`, `agent=mr_reviewer`.

- [ ] **Step 3: Verificar um command**

Run: `cat commands/mr-review/opencode.yml`
Expected:
```text
description: Revisa um Merge Request do GitLab — ...
agent: mr_reviewer
```

- [ ] **Step 4: Commit**

```bash
git add commands/
git commit -m "refactor: split commands into shared body plus per-harness frontmatter"
```

---

## Task 4: Atualizar todas as referências cruzadas snake → kebab

As referências aparecem em **corpos de skills**, **corpos de agentes**, **corpos de commands** e na **documentação**. Skills se referenciam entre si (ex.: `write-code` cita `analyse_code`).

**Files:**
- Modify: `skills/*/SKILL.md`, `agents/*/body.md`, `commands/*/body.md`, `README.md`, `CLAUDE.md`, `AGENTS.md`

- [ ] **Step 1: Mapear ocorrências atuais (baseline)**

Run:
```bash
grep -rohE 'analyse_code|clarify_intent|detail_tasks|document_plan|get_plan|kanban_force|plan_implementation|plan_tasks|query_argocd|review_code|review_mr|test_code|version_code|write_code' skills agents commands README.md CLAUDE.md AGENTS.md | sort | uniq -c
```
Expected: uma contagem por nome snake_case (baseline a zerar após a substituição).

- [ ] **Step 2: Substituir nomes de skills em corpos e docs**

Run:
```bash
SKILLS='analyse_code clarify_intent detail_tasks document_plan get_plan kanban_force plan_implementation plan_tasks query_argocd review_code review_mr test_code version_code write_code'
FILES=$(grep -rlE "$(echo $SKILLS | tr ' ' '|')" skills agents commands README.md CLAUDE.md AGENTS.md)
for s in $SKILLS; do
  k="${s//_/-}"
  for f in $FILES; do sed -i "s/\b${s}\b/${k}/g" "$f"; done
done
echo "substituído"
```
Note: o `kanban-force` (MCP) e o `atlassian_local` (MCP) **não** estão na lista de skills — não são tocados. Os nomes de **MCP** `kanban-force` e `atlassian_local` permanecem como estão.

- [ ] **Step 3: Substituir nomes de commands nas docs**

Run:
```bash
for pair in "doc_plan:doc-plan" "kanban_card:kanban-card" "mr_review:mr-review"; do
  s="${pair%%:*}"; k="${pair##*:}"
  grep -rl "/$s" README.md CLAUDE.md AGENTS.md commands 2>/dev/null | while read -r f; do
    sed -i "s#/${s}#/${k}#g" "$f"
  done
done
echo "commands atualizados"
```
Note: `get_plan` como command já é coberto pela substituição de skills (vira `get-plan`); confira manualmente que `/get_plan` virou `/get-plan` no README.

- [ ] **Step 4 (teste): Garantir que nenhum nome snake_case de skill sobrou**

Run:
```bash
grep -rnE '\b(analyse_code|clarify_intent|detail_tasks|document_plan|get_plan|kanban_force|plan_implementation|plan_tasks|query_argocd|review_code|review_mr|test_code|version_code|write_code)\b' skills agents commands README.md CLAUDE.md AGENTS.md && echo "AINDA HÁ snake_case" || echo "limpo"
```
Expected: `limpo`.

- [ ] **Step 5 (teste): Reconferir que os diretórios de skills batem com as referências**

Run:
```bash
for k in analyse-code clarify-intent detail-tasks document-plan get-plan kanban-force plan-implementation plan-tasks query-argocd review-code review-mr test-code version-code write-code; do
  test -f "skills/$k/SKILL.md" || echo "FALTA skills/$k/SKILL.md"
done; echo "ok"
```
Expected: apenas `ok`.

- [ ] **Step 6: Commit**

```bash
git add skills agents commands README.md CLAUDE.md AGENTS.md
git commit -m "refactor: update skill and command references to kebab-case names"
```

---

## Task 5: Ajustar o `install.sh` para o novo layout (OpenCode, interino)

Nesta fase o instalador continua **só OpenCode**, mas precisa: (a) copiar skills como **pastas**, (b) montar agentes a partir de `opencode.yml` + `body.md`, (c) montar commands a partir de `opencode.yml` + `body.md`, (d) trocar `__OPENCODE_MAIN__` pelo modelo do vendor. Funções são escritas de forma a serem estendidas na Fase 2.

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Atualizar listas e destino para o novo layout**

Substituir os arrays `AGENTS`/`SKILLS`/`COMMANDS` (caminhos de arquivo) por descoberta de **diretórios** e ajustar destinos para `~/.config/opencode/`:

```bash
OPENCODE_DIR="${OPENCODE_DIR:-$HOME/.config/opencode}"
AGENTS_DST="$OPENCODE_DIR/agents"
SKILLS_DST="$OPENCODE_DIR/skills"
COMMANDS_DST="$OPENCODE_DIR/commands"

AGENT_NAMES=(analyzer business_reviewer clarifier code_reviewer coder detailer documenter infra kanban lead mr_reviewer planner tester versioner)
SKILL_NAMES=(analyse-code clarify-intent detail-tasks document-plan get-plan kanban-force plan-implementation plan-tasks query-argocd review-code review-mr test-code version-code write-code)
COMMAND_NAMES=(doc-plan get-plan kanban-card mr-review)
```

- [ ] **Step 2: Adicionar função de montagem de agente (frontmatter + body)**

```bash
# monta <harness>.yml + body.md => arquivo final com frontmatter YAML
assemble() {
  local dir="$1"      # ex: agents/coder
  local harness="$2"  # opencode
  local dst="$3"      # caminho final .md
  {
    echo "---"
    cat "$dir/${harness}.yml"
    echo "---"
    echo ""
    cat "$dir/body.md"
  } > "$dst"
}

resolve_opencode_model() { # troca __OPENCODE_MAIN__ pelo modelo principal do vendor
  local file="$1" model="$2"
  sed -i "s|__OPENCODE_MAIN__|$model|g" "$file"
}
```

- [ ] **Step 3: Definir o modelo principal do OpenCode (vendor opcional)**

Manter o menu de vendor, porém **opcional** e só com a coluna `main`. Sem escolha → `openai/gpt-5.5`:

```bash
declare -A MODEL_MAIN=(
  [anthropic]="anthropic/claude-sonnet-4-6"
  [openai]="openai/gpt-5.5"
  [google]="google/gemini-2.5-pro"
  [groq]="groq/llama-3.3-70b-versatile"
  [amazon-bedrock]="amazon-bedrock/amazon.nova-pro-v1:0"
  [github-copilot]="github-copilot/claude-sonnet-4.6"
)
OPENCODE_MAIN="openai/gpt-5.5"   # default quando nenhum vendor é escolhido
# se o usuário escolher vendor: OPENCODE_MAIN="${MODEL_MAIN[$VENDOR]}"
```

- [ ] **Step 4: Reescrever os loops de instalação para o novo layout**

```bash
# agentes
for name in "${AGENT_NAMES[@]}"; do
  src_local="$SCRIPT_DIR/agents/$name"   # (modo --local; remoto: baixar os 3 arquivos)
  dst="$AGENTS_DST/$name.md"
  assemble "$src_local" opencode "$dst"
  resolve_opencode_model "$dst" "$OPENCODE_MAIN"   # no-op para subagentes (sem token)
done

# skills (copiar a pasta inteira)
for name in "${SKILL_NAMES[@]}"; do
  mkdir -p "$SKILLS_DST/$name"
  cp -R "$SCRIPT_DIR/skills/$name/." "$SKILLS_DST/$name/"
done

# commands
for name in "${COMMAND_NAMES[@]}"; do
  assemble "$SCRIPT_DIR/commands/$name" opencode "$COMMANDS_DST/$name.md"
done
```
Note: o modo remoto (sem `--local`) deve baixar os arquivos `body.md`/`opencode.yml` de cada diretório via `curl`. Como a lista de arquivos por diretório é fixa (`body.md`, `opencode.yml`, `claude.yml`, e `SKILL.md` para skills), baixe-os explicitamente por nome.

- [ ] **Step 5 (teste): Instalar em diretório temporário e conferir**

Run:
```bash
rm -rf /tmp/oc-test && OPENCODE_DIR=/tmp/oc-test ./install.sh --local --force </dev/null
echo "--- coder (primary) ---"; sed -n '1,6p' /tmp/oc-test/agents/coder.md
echo "--- versioner (subagent) ---"; sed -n '1,6p' /tmp/oc-test/agents/versioner.md
echo "--- skill folder ---"; ls /tmp/oc-test/skills/write-code/
```
Expected:
- `coder.md` começa com frontmatter contendo `model: openai/gpt-5.5` (default) e o corpo `<role>`.
- `versioner.md` **sem** linha `model:`.
- `skills/write-code/` contém `SKILL.md`.

- [ ] **Step 6 (teste): Garantir que nenhum token `__OPENCODE_MAIN__` vazou**

Run: `grep -rn '__OPENCODE_MAIN__' /tmp/oc-test && echo "VAZOU" || echo "ok"`
Expected: `ok`.

- [ ] **Step 7: Commit**

```bash
git add install.sh
git commit -m "feat: install from new hybrid layout (folder skills + assembled agents) for OpenCode"
```

---

## Task 6: Atualizar a documentação (README, CLAUDE.md, AGENTS.md)

**Files:**
- Modify: `README.md`, `CLAUDE.md`, `AGENTS.md`

- [ ] **Step 1: Atualizar a seção de estrutura/diretórios**

Refletir o novo layout em todos os três arquivos:
- `skills/<name>/SKILL.md` com `name` kebab-case.
- `agents/<name>/{body.md, opencode.yml, claude.yml}`.
- `commands/<name>/{body.md, opencode.yml, claude.yml}`.
- Diretório de instalação OpenCode agora é `~/.config/opencode/` (plural: `agents/`, `skills/`, `commands/`).

- [ ] **Step 2: Atualizar a tabela de modelos no README**

Remover a coluna `light`. Documentar a nova regra: principais recebem `openai/gpt-5.5` (default) ou `vendor/main`; subagentes não definem `model`. Remover a menção a `versioner` como "modelo light".

- [ ] **Step 3 (teste): Conferir que não há referência ao layout antigo**

Run:
```bash
grep -rnE 'skills/[a-z_]+\.md|~/\.opencode/|modelo light|version_code' README.md CLAUDE.md AGENTS.md && echo "REVISAR" || echo "ok"
```
Expected: `ok` (ou revisar manualmente os casos legítimos remanescentes).

- [ ] **Step 4: Commit**

```bash
git add README.md CLAUDE.md AGENTS.md
git commit -m "docs: document hybrid layout, Agent Skills folders and new model rule"
```

---

## Task 7: Validação final da Fase 1

- [ ] **Step 1 (teste): Revalidar todas as skills**

Run: `for d in skills/*/; do npx -y skills-ref validate "./$d" || echo "FALHOU: $d"; done`
Expected: nenhum `FALHOU:`.

- [ ] **Step 2 (teste): Instalação limpa em temp e contagem de artefatos**

Run:
```bash
rm -rf /tmp/oc-final && OPENCODE_DIR=/tmp/oc-final ./install.sh --local --force </dev/null >/dev/null
echo "agents: $(ls /tmp/oc-final/agents/*.md | wc -l) (esperado 14)"
echo "skills: $(ls -d /tmp/oc-final/skills/*/ | wc -l) (esperado 14)"
echo "commands: $(ls /tmp/oc-final/commands/*.md | wc -l) (esperado 4)"
```
Expected: 14 / 14 / 4.

- [ ] **Step 3 (teste): Sanidade do frontmatter dos agentes instalados**

Run:
```bash
for f in /tmp/oc-final/agents/*.md; do
  head -1 "$f" | grep -q '^---$' || echo "SEM frontmatter: $f"
done; echo "ok"
```
Expected: apenas `ok`.

- [ ] **Step 4: Commit final (se houver ajustes)**

```bash
git add -A && git commit -m "test: validate phase 1 layout and OpenCode install" || echo "nada a commitar"
```

---

## Self-review (cobertura do spec)

- ✅ Skills → pastas `SKILL.md` + `name` kebab (Task 1) + validação (Task 1/7).
- ✅ Renomeações snake→kebab e referências cruzadas, incluindo skills↔skills (Task 4).
- ✅ Agentes em layout híbrido com regra de modelo (principais com `model`, subagentes sem) (Task 2).
- ✅ Commands em layout híbrido (Task 3).
- ✅ `install.sh` lendo o novo layout, OpenCode com `openai/gpt-5.5` default (Task 5).
- ✅ Docs atualizadas (Task 6).
- A instalação multi-harness (Claude Code/Codex) é a **Fase 2** (plano separado).
