# Fase 2 — Instalação multi-harness — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reescrever o `install.sh` para instalar os agentes, skills e commands em qualquer um dos três harnesses (OpenCode, Claude Code, Codex), escolhidos antes de instalar, cada um no diretório e formato nativos, com resolução de modelo ciente do harness.

**Architecture:** O instalador pergunta quais harnesses instalar (um ou vários). Para cada harness, monta agentes (`<harness>.yml` + `body.md`), copia skills (pastas `SKILL.md`) e commands no diretório nativo. OpenCode oferece menu de vendor opcional (default `openai/gpt-5.5`); Claude Code usa alias `sonnet` para principais; Codex herda modelo e recebe skills + prompts + `AGENTS.md`.

**Tech Stack:** Bash (`install.sh`), Markdown/YAML do layout híbrido produzido na Fase 1.

**Pré-requisito:** Fase 1 concluída (layout híbrido em `skills/<n>/SKILL.md`, `agents/<n>/{body.md,opencode.yml,claude.yml}`, `commands/<n>/{...}`).

---

## Mapa de instalação por harness

| Artefato | OpenCode (`~/.config/opencode`) | Claude Code (`~/.claude`) | Codex |
|---|---|---|---|
| Skills | `skills/<n>/SKILL.md` | `skills/<n>/SKILL.md` | `~/.agents/skills/<n>/SKILL.md` |
| Agentes | `agents/<n>.md` (`opencode.yml`) | `agents/<n>.md` (`claude.yml`) | — (só `AGENTS.md`) |
| Commands | `commands/<n>.md` (`opencode.yml`) | `commands/<n>.md` (`claude.yml`) | `~/.codex/prompts/<n>.md` (body-only) |
| Modelo principal | `openai/gpt-5.5` / `vendor/main` | `sonnet` | herda |
| Override de dir | `OPENCODE_DIR` | `CLAUDE_DIR` | `CODEX_SKILLS_DIR` / `CODEX_DIR` |

Codex usa `~/.agents/skills` para skills e `~/.codex/prompts` para prompts. O `AGENTS.md` de orquestração é copiado para `~/.codex/AGENTS.md`.

---

## Task 1: Seleção de harness(es) antes de instalar

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Adicionar a função de seleção múltipla de harness**

```bash
HARNESSES=()   # preenchido pela seleção

select_harness() {
  info "Selecione o(s) harness(es) de destino:"
  echo "        1) opencode"
  echo "        2) claude"
  echo "        3) codex"
  echo "        4) todos"
  local choice
  read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    Números separados por espaço (ex.: 1 2): ")" -a choice </dev/tty
  for c in "${choice[@]}"; do
    case "$c" in
      1) HARNESSES+=(opencode) ;;
      2) HARNESSES+=(claude) ;;
      3) HARNESSES+=(codex) ;;
      4) HARNESSES=(opencode claude codex) ;;
    esac
  done
  [[ ${#HARNESSES[@]} -gt 0 ]] || { echo "Nenhum harness selecionado."; exit 1; }
  # dedup
  HARNESSES=($(printf '%s\n' "${HARNESSES[@]}" | sort -u))
  ok "Harnesses: ${HARNESSES[*]}"
}
```

- [ ] **Step 2 (teste): Selecionar dois harnesses e conferir o array**

Run: `printf '1 2\n' | bash -c 'source install.sh; :' 2>/dev/null || true`
Note: como `install.sh` executa tudo, valide a função isoladamente extraindo-a para `/tmp` ou rodando o install completo no Step da Task 8. Aqui basta inspeção visual de que a seleção `1 2` resulta em `HARNESSES=(opencode claude)`.

- [ ] **Step 3: Commit**

```bash
git add install.sh
git commit -m "feat: add multi-harness selection menu to installer"
```

---

## Task 2: Resolver diretórios de destino por harness

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Adicionar função que define os destinos conforme o harness**

```bash
# variáveis preenchidas por harness_paths()
H_SKILLS=""; H_AGENTS=""; H_COMMANDS=""; H_PROMPTS=""; H_AGENTSMD=""

harness_paths() {
  local h="$1"
  case "$h" in
    opencode)
      local base="${OPENCODE_DIR:-$HOME/.config/opencode}"
      H_SKILLS="$base/skills"; H_AGENTS="$base/agents"; H_COMMANDS="$base/commands"
      H_PROMPTS=""; H_AGENTSMD="" ;;
    claude)
      local base="${CLAUDE_DIR:-$HOME/.claude}"
      H_SKILLS="$base/skills"; H_AGENTS="$base/agents"; H_COMMANDS="$base/commands"
      H_PROMPTS=""; H_AGENTSMD="" ;;
    codex)
      H_SKILLS="${CODEX_SKILLS_DIR:-$HOME/.agents/skills}"
      H_AGENTS=""; H_COMMANDS=""
      H_PROMPTS="${CODEX_DIR:-$HOME/.codex}/prompts"
      H_AGENTSMD="${CODEX_DIR:-$HOME/.codex}/AGENTS.md" ;;
  esac
}
```

- [ ] **Step 2 (teste): Conferir os caminhos resolvidos**

Run:
```bash
source <(awk '/^harness_paths\(\)/,/^}/' install.sh) 2>/dev/null
for h in opencode claude codex; do harness_paths "$h"; echo "$h: skills=$H_SKILLS agents=${H_AGENTS:-—} prompts=${H_PROMPTS:-—}"; done
```
Expected:
```text
opencode: skills=~/.config/opencode/skills agents=~/.config/opencode/agents prompts=—
claude:   skills=~/.claude/skills agents=~/.claude/agents prompts=—
codex:    skills=~/.agents/skills agents=— prompts=~/.codex/prompts
```

- [ ] **Step 3: Commit**

```bash
git add install.sh
git commit -m "feat: resolve native target directories per harness"
```

---

## Task 3: Vendor opcional + resolução de modelo principal por harness

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Tornar a seleção de vendor opcional (só relevante para OpenCode)**

```bash
declare -A MODEL_MAIN=(
  [anthropic]="anthropic/claude-sonnet-4-6"
  [openai]="openai/gpt-5.5"
  [google]="google/gemini-2.5-pro"
  [groq]="groq/llama-3.3-70b-versatile"
  [amazon-bedrock]="amazon-bedrock/amazon.nova-pro-v1:0"
  [github-copilot]="github-copilot/claude-sonnet-4.6"
)
OPENCODE_MAIN="openai/gpt-5.5"   # default

select_vendor_optional() {
  info "Vendor do OpenCode (Enter para usar o default openai/gpt-5.5):"
  local i=1; local names=(anthropic openai google groq amazon-bedrock github-copilot)
  for v in "${names[@]}"; do printf "        %d) %-16s %s\n" "$i" "$v" "${MODEL_MAIN[$v]}"; i=$((i+1)); done
  local choice
  read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    Número (ou Enter): ")" choice </dev/tty
  if [[ "$choice" =~ ^[1-6]$ ]]; then OPENCODE_MAIN="${MODEL_MAIN[${names[$((choice-1))]}]}"; fi
  ok "Modelo principal OpenCode: $OPENCODE_MAIN"
}
```
O menu de vendor só é chamado se `opencode` estiver em `HARNESSES`.

- [ ] **Step 2: Definir a substituição de modelo por harness**

```bash
# aplica a regra de modelo no arquivo de agente já montado
apply_model() {
  local file="$1" harness="$2"
  case "$harness" in
    opencode) sed -i "s|__OPENCODE_MAIN__|$OPENCODE_MAIN|g" "$file" ;;
    claude)   : ;;  # claude.yml já traz `model: sonnet` nos primários; subagentes sem model
    codex)    : ;;  # codex não recebe agentes nesta fase
  esac
}
```

- [ ] **Step 3 (teste): Conferir resolução do default e do vendor**

Run:
```bash
tmp=$(mktemp); echo "model: __OPENCODE_MAIN__" > "$tmp"
OPENCODE_MAIN="openai/gpt-5.5"; sed -i "s|__OPENCODE_MAIN__|$OPENCODE_MAIN|g" "$tmp"; cat "$tmp"
```
Expected: `model: openai/gpt-5.5`.

- [ ] **Step 4: Commit**

```bash
git add install.sh
git commit -m "feat: optional vendor selection and harness-aware model resolution"
```

---

## Task 4: Instalar agentes por harness (OpenCode + Claude Code)

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Generalizar o loop de agentes para usar `<harness>.yml`**

Reusar a função `assemble()` da Fase 1, parametrizada pelo harness:

```bash
install_agents() {
  local harness="$1"
  [[ -n "$H_AGENTS" ]] || return 0   # codex: sem agentes
  mkdir -p "$H_AGENTS"
  for name in "${AGENT_NAMES[@]}"; do
    local dst="$H_AGENTS/$name.md"
    # conflito + confirmação (reaproveitar lógica existente de confirm/skip)
    assemble "$SCRIPT_DIR/agents/$name" "$harness" "$dst"
    apply_model "$dst" "$harness"
  done
}
```
`assemble()` (da Fase 1) escreve `--- + <harness>.yml + --- + body.md`.

- [ ] **Step 2 (teste): Montar agente para Claude Code e validar frontmatter**

Run:
```bash
rm -rf /tmp/cc && CLAUDE_DIR=/tmp/cc ./install.sh --local --force <<< "2" >/dev/null
echo "--- coder (primary) ---"; sed -n '1,5p' /tmp/cc/agents/coder.md
echo "--- tester (subagent) ---"; sed -n '1,5p' /tmp/cc/agents/tester.md
```
Expected:
- `coder.md`: frontmatter com `name: coder`, `description:`, `model: sonnet` (sem `mode`, sem `temperature`).
- `tester.md`: `name: tester`, `description:`, **sem** `model`.

- [ ] **Step 3 (teste): Garantir que OpenCode mantém `mode` e `temperature`**

Run:
```bash
rm -rf /tmp/oc && OPENCODE_DIR=/tmp/oc ./install.sh --local --force <<< $'1\n' >/dev/null
sed -n '1,6p' /tmp/oc/agents/coder.md
```
Expected: frontmatter com `description`, `mode: primary`, `model: openai/gpt-5.5`, `temperature: 0.3`.

- [ ] **Step 4: Commit**

```bash
git add install.sh
git commit -m "feat: install agents per harness with native frontmatter"
```

---

## Task 5: Instalar skills por harness

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Generalizar a cópia de skills**

```bash
install_skills() {
  mkdir -p "$H_SKILLS"
  for name in "${SKILL_NAMES[@]}"; do
    mkdir -p "$H_SKILLS/$name"
    cp -R "$SCRIPT_DIR/skills/$name/." "$H_SKILLS/$name/"
  done
}
```
Skills são **idênticas** nos três harnesses — só muda `$H_SKILLS`.

- [ ] **Step 2 (teste): Skills nos três destinos**

Run:
```bash
rm -rf /tmp/oc /tmp/cc /tmp/cdx
OPENCODE_DIR=/tmp/oc CLAUDE_DIR=/tmp/cc CODEX_SKILLS_DIR=/tmp/cdx/skills CODEX_DIR=/tmp/cdx ./install.sh --local --force <<< "4" >/dev/null
for d in /tmp/oc/skills /tmp/cc/skills /tmp/cdx/skills; do echo "$d: $(ls -d $d/*/ | wc -l) skills"; done
test -f /tmp/cdx/skills/write-code/SKILL.md && echo "codex SKILL.md ok"
```
Expected: cada destino com 14 skills; `codex SKILL.md ok`.

- [ ] **Step 3: Commit**

```bash
git add install.sh
git commit -m "feat: install skill folders to each harness skills dir"
```

---

## Task 6: Instalar commands (OpenCode/Claude) e prompts (Codex)

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Commands para OpenCode e Claude (frontmatter + body)**

```bash
install_commands() {
  local harness="$1"
  if [[ -n "$H_COMMANDS" ]]; then
    mkdir -p "$H_COMMANDS"
    for name in "${COMMAND_NAMES[@]}"; do
      assemble "$SCRIPT_DIR/commands/$name" "$harness" "$H_COMMANDS/$name.md"
    done
  fi
}
```

- [ ] **Step 2: Prompts body-only para o Codex**

```bash
install_codex_prompts() {
  [[ -n "$H_PROMPTS" ]] || return 0
  mkdir -p "$H_PROMPTS"
  for name in "${COMMAND_NAMES[@]}"; do
    cp "$SCRIPT_DIR/commands/$name/body.md" "$H_PROMPTS/$name.md"
  done
}
```
No Codex o prompt é só o corpo (sem frontmatter de `agent:`, que não se aplica).

- [ ] **Step 3 (teste): Commands e prompts corretos**

Run:
```bash
echo "OC: $(ls /tmp/oc/commands/*.md | wc -l) commands"
head -3 /tmp/oc/commands/mr-review.md
echo "Codex prompt (sem frontmatter):"
head -1 /tmp/cdx/prompts/mr-review.md
```
Expected: OC com 4 commands e frontmatter `---`; o prompt do Codex começa direto no corpo (primeira linha **não** é `---`).

- [ ] **Step 4: Commit**

```bash
git add install.sh
git commit -m "feat: install commands for OpenCode/Claude and body-only prompts for Codex"
```

---

## Task 7: Instalar `AGENTS.md` de orquestração no Codex

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Copiar o `AGENTS.md` para o destino do Codex**

```bash
install_codex_agentsmd() {
  [[ -n "$H_AGENTSMD" ]] || return 0
  mkdir -p "$(dirname "$H_AGENTSMD")"
  if [[ -f "$H_AGENTSMD" ]] && ! $FORCE; then
    confirm "Substituir $H_AGENTSMD?" || return 0
  fi
  cp "$SCRIPT_DIR/AGENTS.md" "$H_AGENTSMD"
}
```

- [ ] **Step 2 (teste): AGENTS.md presente no Codex**

Run: `test -f /tmp/cdx/AGENTS.md && echo "ok" || echo "FALTA"`
Expected: `ok`.

- [ ] **Step 3: Commit**

```bash
git add install.sh
git commit -m "feat: install orchestration AGENTS.md for Codex"
```

---

## Task 8: Orquestrar o fluxo principal e o resumo final

**Files:**
- Modify: `install.sh`

- [ ] **Step 1: Montar o fluxo `main` chamando as funções por harness**

```bash
select_harness
case " ${HARNESSES[*]} " in *" opencode "*) select_vendor_optional ;; esac

for h in "${HARNESSES[@]}"; do
  harness_paths "$h"
  info "Instalando para: $h"
  install_agents "$h"
  install_skills
  install_commands "$h"
  if [[ "$h" == "codex" ]]; then install_codex_prompts; install_codex_agentsmd; fi
  ok "Concluído: $h"
done

print_summary   # diretórios usados, modelo aplicado, contagens por harness
```

- [ ] **Step 2: Implementar `print_summary`**

```bash
print_summary() {
  echo ""
  ok "Instalação concluída."
  for h in "${HARNESSES[@]}"; do
    harness_paths "$h"
    echo "  • $h"
    echo "      skills:   $H_SKILLS"
    [[ -n "$H_AGENTS"   ]] && echo "      agents:   $H_AGENTS"
    [[ -n "$H_COMMANDS" ]] && echo "      commands: $H_COMMANDS"
    [[ -n "$H_PROMPTS"  ]] && echo "      prompts:  $H_PROMPTS"
    [[ "$h" == "opencode" ]] && echo "      modelo principal: $OPENCODE_MAIN"
    [[ "$h" == "claude"   ]] && echo "      modelo principal: sonnet"
    [[ "$h" == "codex"    ]] && echo "      modelo: herdado da sessão"
  done
}
```

- [ ] **Step 3 (teste): Instalação completa nos 3 harnesses de uma vez**

Run:
```bash
rm -rf /tmp/oc /tmp/cc /tmp/cdx
OPENCODE_DIR=/tmp/oc CLAUDE_DIR=/tmp/cc CODEX_SKILLS_DIR=/tmp/cdx/skills CODEX_DIR=/tmp/cdx \
  ./install.sh --local --force <<< $'4\n' 
```
Expected (resumo final):
- opencode: 14 agents + 14 skills + 4 commands, modelo `openai/gpt-5.5`.
- claude: 14 agents + 14 skills + 4 commands, modelo `sonnet`.
- codex: 14 skills + 4 prompts + `AGENTS.md`, sem agents.

- [ ] **Step 4 (teste): Nenhum token de modelo vazou em nenhum destino**

Run: `grep -rn '__OPENCODE_MAIN__' /tmp/oc /tmp/cc /tmp/cdx && echo "VAZOU" || echo "ok"`
Expected: `ok`.

- [ ] **Step 5: Commit**

```bash
git add install.sh
git commit -m "feat: orchestrate multi-harness install flow with per-harness summary"
```

---

## Task 9: Atualizar README e ajuda do instalador

**Files:**
- Modify: `README.md`, `install.sh` (texto de `--help`)

- [ ] **Step 1: Documentar a seleção de harness e os diretórios de cada um**

Adicionar ao README: tabela de harnesses × diretórios, menu de seleção, vendor opcional (default `openai/gpt-5.5`), alias `sonnet` no Claude Code, Codex herdando modelo. Documentar os env vars de override (`OPENCODE_DIR`, `CLAUDE_DIR`, `CODEX_DIR`, `CODEX_SKILLS_DIR`).

- [ ] **Step 2: Atualizar o bloco `--help` no install.sh**

Incluir a explicação de seleção múltipla de harness e do vendor opcional.

- [ ] **Step 3 (teste): `--help` reflete o novo fluxo**

Run: `./install.sh --help | grep -iE 'harness|vendor' && echo "ok"`
Expected: linhas mencionando harness e vendor; `ok`.

- [ ] **Step 4: Commit**

```bash
git add README.md install.sh
git commit -m "docs: document multi-harness install flow and overrides"
```

---

## Self-review (cobertura do spec)

- ✅ Seleção de harness antes de instalar, um ou vários (Task 1).
- ✅ Diretórios nativos por harness, com overrides por env var (Task 2).
- ✅ Vendor opcional + `openai/gpt-5.5` default; modelo ciente do harness (Task 3).
- ✅ Agentes com frontmatter nativo; subagentes sem `model` (Task 4).
- ✅ Skills idênticas nos três destinos (Task 5).
- ✅ Commands (OC/CC) + prompts body-only (Codex) (Task 6).
- ✅ `AGENTS.md` de orquestração no Codex (Task 7).
- ✅ Fluxo principal + resumo por harness (Task 8).
- ✅ Documentação e `--help` (Task 9).

## Pontos a confirmar antes/durante a execução

- **`openai/gpt-5.5`**: confirmar o id exato do provider/modelo no OpenCode; trocar o valor em `MODEL_MAIN[openai]` e em `OPENCODE_MAIN` se diferir.
- **Codex custom agents nativos**: fora do escopo desta fase (Codex recebe só skills + prompts + `AGENTS.md`). Se desejado depois, adicionar `agents/<n>/codex.yml` e um `install_codex_agents` após confirmar o formato do arquivo de agente do Codex.
- **Modo remoto (sem `--local`)**: ao baixar via `curl`, baixar explicitamente `body.md`/`opencode.yml`/`claude.yml` por diretório de agente/command e `SKILL.md` por skill (a lista de arquivos é fixa por tipo).
