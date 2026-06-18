#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────
# coder — multi-harness agents & skills installer
# (instala a partir do layout híbrido:
#   agents/<name>/{opencode.yml,claude.yml,body.md}
#   commands/<name>/{opencode.yml,claude.yml,body.md}
#   skills/<name>/SKILL.md[, references/])
#
# Harnesses suportados: OpenCode, Claude Code, Codex.
# Cada um recebe os artefatos no diretório e formato nativos.
# ─────────────────────────────────────────────

REPO_URL="https://raw.githubusercontent.com/paraizofelipe/coder/main"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ── listas de nomes (não caminhos) ────────────

AGENT_NAMES=(
  analyzer
  business_reviewer
  clarifier
  code_reviewer
  coder
  detailer
  documenter
  infra
  kanban
  lead
  mr_reviewer
  planner
  qa
  tester
  versioner
)

SKILL_NAMES=(
  analyse-code
  clarify-intent
  detail-tasks
  document-plan
  get-plan
  kanban-force
  plan-implementation
  plan-tasks
  query-argocd
  review-code
  review-mr
  test-code
  validate-implementation
  version-code
  write-code
)

COMMAND_NAMES=(
  doc-plan
  get-plan
  kanban-card
  mr-review
  qa
)

# ── mapa de vendors e modelo principal ────────

VENDOR_NAMES=(
  "anthropic"
  "openai"
  "google"
  "groq"
  "amazon-bedrock"
  "github-copilot"
)

declare -A MODEL_MAIN=(
  [anthropic]="anthropic/claude-sonnet-4-6"
  [openai]="openai/gpt-5.5"
  [google]="google/gemini-2.5-pro"
  [groq]="groq/llama-3.3-70b-versatile"
  [amazon-bedrock]="amazon-bedrock/amazon.nova-pro-v1:0"
  [github-copilot]="github-copilot/claude-sonnet-4.6"
)

# modelo principal do OpenCode (default quando nenhum vendor é escolhido)
OPENCODE_MAIN="openai/gpt-5.5"
VENDOR=""

# ── helpers ──────────────────────────────────

RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

info()    { echo -e "${CYAN}${BOLD}[info]${RESET}  $*"; }
ok()      { echo -e "${GREEN}${BOLD}[ok]${RESET}    $*"; }
warn()    { echo -e "${YELLOW}${BOLD}[warn]${RESET}  $*"; }
skip()    { echo -e "        ${YELLOW}↳ pulado${RESET}"; }
installed() { echo -e "        ${GREEN}↳ instalado${RESET}"; }

confirm() {
  local msg="$1"
  local answer
  read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    $msg [s/N] ")" answer </dev/tty
  [[ "$answer" =~ ^[sSyY]$ ]]
}

# ── montagem (agentes e commands) ─────────────
#
# Monta o arquivo final `.md` no formato:
#   ---
#   <conteúdo de <harness>.yml>
#   ---
#   <linha em branco>
#   <conteúdo de body.md>
#
assemble() {
  local dir="$1"       # diretório fonte (ex: agents/coder)
  local harness="$2"   # opencode | claude
  local dst="$3"       # caminho final do .md
  {
    echo "---"
    cat "$dir/${harness}.yml"
    echo "---"
    echo ""
    cat "$dir/body.md"
  } > "$dst"
}

# aplica a regra de modelo no arquivo de agente já montado.
#  opencode → substitui o token __OPENCODE_MAIN__ pelo modelo escolhido.
#  claude   → no-op (claude.yml já traz `model: sonnet` nos primários; subagentes sem model).
#  codex    → no-op (codex não recebe agentes nesta fase).
apply_model() {
  local file="$1" harness="$2"
  case "$harness" in
    opencode) sed -i.bak "s|__OPENCODE_MAIN__|$OPENCODE_MAIN|g" "$file" && rm -f "$file.bak" ;;
    claude)   : ;;
    codex)    : ;;
  esac
}

# verifica conflito + confirmação antes de sobrescrever.
# retorna 0 = pode prosseguir, 1 = pular.
check_overwrite() {
  local dst="$1"
  local label="$2"
  if [[ -e "$dst" ]]; then
    warn "Já existe: $dst"
    if ! $FORCE; then
      if ! confirm "Substituir $label?"; then
        skip
        return 1
      fi
    fi
  fi
  return 0
}

# ── seleção de harness ────────────────────────

HARNESSES=()
FLAG_HARNESS=""   # valor de --harness (se fornecido)

select_harness() {
  info "Selecione o(s) harness(es) de destino:"
  echo "        1) opencode"
  echo "        2) claude"
  echo "        3) codex"
  echo "        4) todos"
  local choice=()
  # lê de /dev/tty para funcionar em curl|bash (stdin = pipe).
  # se /dev/tty não disponível (sem tty E sem flag), aborta com graça.
  if ! read -r -a choice -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    Números separados por espaço (ex.: 1 2): ")" </dev/tty; then
    choice=()
  fi
  local c
  for c in "${choice[@]}"; do
    case "$c" in
      1) HARNESSES+=(opencode) ;;
      2) HARNESSES+=(claude) ;;
      3) HARNESSES+=(codex) ;;
      4) HARNESSES=(opencode claude codex) ;;
    esac
  done
  if [[ ${#HARNESSES[@]} -eq 0 ]]; then
    warn "Nenhum harness selecionado."
    exit 1
  fi
  # dedup preservando ordem canônica: opencode → claude → codex
  local -a canonical=(opencode claude codex)
  local -a deduped=()
  local h
  for h in "${canonical[@]}"; do
    local present
    for present in "${HARNESSES[@]}"; do
      if [[ "$present" == "$h" ]]; then
        deduped+=("$h")
        break
      fi
    done
  done
  HARNESSES=("${deduped[@]}")
  ok "Harnesses: ${HARNESSES[*]}"
}

# resolve --harness flag value into HARNESSES (canonical order, deduped)
resolve_harness_flag() {
  local input="$1"
  # normaliza separadores: vírgula → espaço
  input="${input//,/ }"
  local -a raw=()
  read -r -a raw <<< "$input"
  local -a canonical=(opencode claude codex)
  local h token found
  for token in "${raw[@]}"; do
    if [[ "$token" == "all" ]]; then
      HARNESSES=(opencode claude codex)
      ok "Harnesses: ${HARNESSES[*]}"
      return
    fi
  done
  # para cada token válido, adiciona; inválidos geram erro
  local -a selected=()
  for token in "${raw[@]}"; do
    found=false
    for h in "${canonical[@]}"; do
      if [[ "$token" == "$h" ]]; then
        selected+=("$h")
        found=true
        break
      fi
    done
    if ! $found; then
      echo -e "${RED}${BOLD}[erro]${RESET} harness inválido: '$token'. Use: opencode, claude, codex, all."
      exit 1
    fi
  done
  if [[ ${#selected[@]} -eq 0 ]]; then
    echo -e "${RED}${BOLD}[erro]${RESET} --harness requer ao menos um valor."
    exit 1
  fi
  # dedup preservando ordem canônica
  local -a deduped=()
  for h in "${canonical[@]}"; do
    local p
    for p in "${selected[@]}"; do
      if [[ "$p" == "$h" ]]; then
        deduped+=("$h")
        break
      fi
    done
  done
  HARNESSES=("${deduped[@]}")
  ok "Harnesses: ${HARNESSES[*]}"
}

# ── diretórios de destino por harness ─────────

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

# ── seleção de vendor (opcional, só OpenCode) ──

FLAG_VENDOR=""   # valor de --vendor (se fornecido)

select_vendor_optional() {
  info "Vendor do OpenCode (Enter para usar o default $OPENCODE_MAIN):"
  echo ""
  local i=1
  for vendor in "${VENDOR_NAMES[@]}"; do
    printf "        ${BOLD}%d)${RESET} %-18s  main: %s\n" \
      "$i" "$vendor" "${MODEL_MAIN[$vendor]}"
    i=$((i + 1))
  done
  echo ""
  local choice
  # lê de /dev/tty para funcionar em curl|bash (stdin = pipe).
  # se /dev/tty não disponível, mantém o default.
  if ! read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    Número do vendor [1-${#VENDOR_NAMES[@]}] (Enter = padrão): ")" choice </dev/tty; then
    choice=""
  fi
  if [[ "$choice" =~ ^[0-9]+$ ]] && (( choice >= 1 && choice <= ${#VENDOR_NAMES[@]} )); then
    VENDOR="${VENDOR_NAMES[$((choice - 1))]}"
    OPENCODE_MAIN="${MODEL_MAIN[$VENDOR]}"
  else
    if [[ -n "$choice" ]]; then
      warn "Opção inválida; mantendo o padrão."
    fi
    VENDOR="(padrão)"
  fi
  echo ""
  ok "Modelo principal OpenCode: ${BOLD}$OPENCODE_MAIN${RESET}"
  echo ""
}

# resolve --vendor flag value into OPENCODE_MAIN
resolve_vendor_flag() {
  local name="$1"
  if [[ -v MODEL_MAIN["$name"] ]]; then
    VENDOR="$name"
    OPENCODE_MAIN="${MODEL_MAIN[$name]}"
    ok "Vendor: $VENDOR → modelo principal: $OPENCODE_MAIN"
  else
    echo -e "${RED}${BOLD}[erro]${RESET} vendor inválido: '$name'. Valores aceitos: ${!MODEL_MAIN[*]}"
    exit 1
  fi
}

# ── argparse ─────────────────────────────────

FORCE=false
LOCAL=false
FETCH_CMD=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --force|-f) FORCE=true; shift ;;
    --local|-l) LOCAL=true; shift ;;
    --harness)
      [[ $# -ge 2 ]] || { echo -e "${RED}${BOLD}[erro]${RESET} --harness requer um valor."; exit 1; }
      FLAG_HARNESS="$2"; shift 2 ;;
    --harness=*) FLAG_HARNESS="${1#--harness=}"; shift ;;
    --vendor)
      [[ $# -ge 2 ]] || { echo -e "${RED}${BOLD}[erro]${RESET} --vendor requer um valor."; exit 1; }
      FLAG_VENDOR="$2"; shift 2 ;;
    --vendor=*) FLAG_VENDOR="${1#--vendor=}"; shift ;;
    --help|-h)
      echo ""
      echo -e "  ${BOLD}Uso:${RESET} install.sh [opções]"
      echo ""
      echo "  Instala agentes, skills e commands em um ou mais harnesses"
      echo "  (OpenCode, Claude Code, Codex), escolhidos antes da instalação."
      echo ""
      echo "  Opções:"
      echo "    --force, -f              Substituir todos os arquivos sem perguntar"
      echo "    --local, -l              Instalar a partir dos arquivos locais do repositório"
      echo "    --harness <lista>        Harness(es) a instalar sem menu interativo."
      echo "                             Valores: opencode, claude, codex, all (ou combinações"
      echo "                             separadas por vírgula/espaço, ex.: opencode,claude)"
      echo "    --vendor <nome>          Vendor para o OpenCode sem menu interativo (ignorado se OpenCode"
      echo "                             não estiver nos harnesses selecionados). Vendor inválido → exit 1."
      echo "                             Valores: anthropic openai google groq amazon-bedrock github-copilot"
      echo "    --help,  -h              Exibir esta ajuda"
      echo ""
      echo "  Overrides de diretório (env vars):"
      echo "    OPENCODE_DIR        base do OpenCode (default ~/.config/opencode)"
      echo "    CLAUDE_DIR          base do Claude Code (default ~/.claude)"
      echo "    CODEX_DIR           base do Codex (default ~/.codex)"
      echo "    CODEX_SKILLS_DIR    skills do Codex (default ~/.agents/skills)"
      echo ""
      exit 0
      ;;
    *) echo -e "${RED}${BOLD}[erro]${RESET} opção desconhecida: '$1'"; exit 1 ;;
  esac
done

# ── verificações ─────────────────────────────

if ! $LOCAL; then
  if command -v curl &>/dev/null; then
    FETCH_CMD="curl"
  elif command -v wget &>/dev/null; then
    FETCH_CMD="wget"
  else
    echo -e "${RED}${BOLD}[erro]${RESET} curl ou wget são necessários para a instalação remota."
    echo "       Use --local para instalar a partir de arquivos locais."
    exit 1
  fi
fi

# baixa um arquivo remoto para o destino indicado.
fetch_remote() {
  local rel="$1"   # caminho relativo no repositório (ex: agents/coder/body.md)
  local dst="$2"
  if [[ "$FETCH_CMD" == "curl" ]]; then
    curl -fsSL "$REPO_URL/$rel" -o "$dst" || { echo "[erro] Falha ao baixar: $REPO_URL/$rel"; return 1; }
  else
    wget -qO "$dst" "$REPO_URL/$rel" || { echo "[erro] Falha ao baixar: $REPO_URL/$rel"; return 1; }
  fi
}

# ── fonte por tipo (local ou remoto) ──────────
#
# Em modo local, a fonte é o próprio repositório.
# Em modo remoto, baixamos os arquivos por diretório para um temp e
# devolvemos esse temp como diretório fonte.

REMOTE_TMP=""
cleanup() {
  [[ -n "$REMOTE_TMP" && -d "$REMOTE_TMP" ]] && rm -rf "$REMOTE_TMP"
  return 0
}
trap cleanup EXIT

# prepara o diretório fonte de um agente/command (contém <harness>.yml + body.md).
# uso: src_dir=$(prepare_assembled_src agents coder opencode)
prepare_assembled_src() {
  local kind="$1"      # agents | commands
  local name="$2"
  local harness="$3"   # opencode | claude
  if $LOCAL; then
    echo "$SCRIPT_DIR/$kind/$name"
    return
  fi
  local tmp="$REMOTE_TMP/$kind/$name"
  mkdir -p "$tmp"
  fetch_remote "$kind/$name/${harness}.yml" "$tmp/${harness}.yml"
  fetch_remote "$kind/$name/body.md" "$tmp/body.md"
  echo "$tmp"
}

# prepara o body.md de um command (usado pelo Codex, body-only).
prepare_body_src() {
  local kind="$1"   # commands
  local name="$2"
  if $LOCAL; then
    echo "$SCRIPT_DIR/$kind/$name/body.md"
    return
  fi
  local tmp="$REMOTE_TMP/$kind/$name"
  mkdir -p "$tmp"
  fetch_remote "$kind/$name/body.md" "$tmp/body.md"
  echo "$tmp/body.md"
}

# ── instalação por harness ────────────────────

install_agents() {
  local harness="$1"
  [[ -n "$H_AGENTS" ]] || return 0   # codex: sem agentes
  mkdir -p "$H_AGENTS"
  info "Instalando agentes em $H_AGENTS"
  echo ""
  for name in "${AGENT_NAMES[@]}"; do
    echo -e "  ${BOLD}$name${RESET}"
    local dst="$H_AGENTS/$name.md"
    if ! check_overwrite "$dst" "$name"; then
      continue
    fi
    local src_dir
    src_dir="$(prepare_assembled_src agents "$name" "$harness")"
    assemble "$src_dir" "$harness" "$dst"
    apply_model "$dst" "$harness"
    installed
  done
  echo ""
}

install_skills() {
  mkdir -p "$H_SKILLS"
  info "Instalando skills em $H_SKILLS"
  echo ""
  for name in "${SKILL_NAMES[@]}"; do
    echo -e "  ${BOLD}$name${RESET}"
    local dst="$H_SKILLS/$name"
    if ! check_overwrite "$dst" "$name"; then
      continue
    fi
    rm -rf "$dst"
    mkdir -p "$dst"
    if $LOCAL; then
      cp -R "$SCRIPT_DIR/skills/$name/." "$dst/"
    else
      # modo remoto: baixamos apenas SKILL.md. Skills com references/ não são
      # suportadas no modo remoto desta fase.
      fetch_remote "skills/$name/SKILL.md" "$dst/SKILL.md"
    fi
    installed
  done
  echo ""
}

install_commands() {
  local harness="$1"
  [[ -n "$H_COMMANDS" ]] || return 0   # codex: tratado por install_codex_prompts
  mkdir -p "$H_COMMANDS"
  info "Instalando commands em $H_COMMANDS"
  echo ""
  for name in "${COMMAND_NAMES[@]}"; do
    echo -e "  ${BOLD}$name${RESET}"
    local dst="$H_COMMANDS/$name.md"
    if ! check_overwrite "$dst" "$name"; then
      continue
    fi
    local src_dir
    src_dir="$(prepare_assembled_src commands "$name" "$harness")"
    assemble "$src_dir" "$harness" "$dst"
    apply_model "$dst" "$harness"
    installed
  done
  echo ""
}

# prompts body-only para o Codex (sem frontmatter de `agent:`, que não se aplica).
install_codex_prompts() {
  [[ -n "$H_PROMPTS" ]] || return 0
  mkdir -p "$H_PROMPTS"
  info "Instalando prompts em $H_PROMPTS"
  echo ""
  for name in "${COMMAND_NAMES[@]}"; do
    echo -e "  ${BOLD}$name${RESET}"
    local dst="$H_PROMPTS/$name.md"
    if ! check_overwrite "$dst" "$name"; then
      continue
    fi
    local src
    src="$(prepare_body_src commands "$name")"
    cp "$src" "$dst"
    installed
  done
  echo ""
}

# AGENTS.md de orquestração para o Codex.
# NOTA: AGENTS.md existe no disco mas é gitignored (não está no GitHub).
# Em --local a cópia local funciona; em modo remoto o download retorna 404,
# por isso a falha remota é não-fatal (warn e segue).
install_codex_agentsmd() {
  [[ -n "$H_AGENTSMD" ]] || return 0
  info "Instalando AGENTS.md em $H_AGENTSMD"
  mkdir -p "$(dirname "$H_AGENTSMD")"
  if ! check_overwrite "$H_AGENTSMD" "AGENTS.md"; then
    return 0
  fi
  if $LOCAL; then
    cp "$SCRIPT_DIR/AGENTS.md" "$H_AGENTSMD"
    installed
  else
    if fetch_remote "AGENTS.md" "$H_AGENTSMD"; then
      installed
    else
      warn "AGENTS.md indisponível no modo remoto (gitignored); pulando."
      rm -f "$H_AGENTSMD"
    fi
  fi
  echo ""
}

# ── resumo final ──────────────────────────────

print_summary() {
  echo ""
  ok "Instalação concluída."
  echo ""
  for h in "${HARNESSES[@]}"; do
    harness_paths "$h"
    echo -e "  ${BOLD}• $h${RESET}"
    echo "      skills:   $H_SKILLS"
    [[ -n "$H_AGENTS"   ]] && echo "      agents:   $H_AGENTS"
    [[ -n "$H_COMMANDS" ]] && echo "      commands: $H_COMMANDS"
    [[ -n "$H_PROMPTS"  ]] && echo "      prompts:  $H_PROMPTS"
    [[ -n "$H_AGENTSMD" ]] && echo "      AGENTS.md: $H_AGENTSMD"
    case "$h" in
      opencode) echo "      modelo principal: $OPENCODE_MAIN" ;;
      claude)   echo "      modelo principal: sonnet" ;;
      codex)    echo "      modelo: herdado da sessão" ;;
    esac
  done
  echo ""
  echo -e "  Reinicie o harness para carregar os novos agentes e skills."
  echo ""
}

# ── fluxo principal ───────────────────────────

echo ""
echo -e "${BOLD}  coder — multi-harness agents & skills${RESET}"
echo "  ─────────────────────────────────────"
echo ""

if [[ -n "$FLAG_HARNESS" ]]; then
  resolve_harness_flag "$FLAG_HARNESS"
else
  select_harness
fi

case " ${HARNESSES[*]} " in
  *" opencode "*)
    if [[ -n "$FLAG_VENDOR" ]]; then
      resolve_vendor_flag "$FLAG_VENDOR"
    else
      select_vendor_optional
    fi
    ;;
esac

if ! $LOCAL; then
  REMOTE_TMP="$(mktemp -d)"
fi

for h in "${HARNESSES[@]}"; do
  harness_paths "$h"
  info "Instalando para: ${BOLD}$h${RESET}"
  echo ""
  install_agents "$h"
  install_skills
  install_commands "$h"
  if [[ "$h" == "codex" ]]; then
    install_codex_prompts
    install_codex_agentsmd
  fi
  ok "Concluído: $h"
  echo ""
done

print_summary
