#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────
# coder — OpenCode agents & skills installer
# (instala a partir do layout híbrido:
#   agents/<name>/{opencode.yml,body.md}
#   commands/<name>/{opencode.yml,body.md}
#   skills/<name>/SKILL.md[, references/])
# ─────────────────────────────────────────────

REPO_URL="https://raw.githubusercontent.com/paraizofelipe/coder/main"

OPENCODE_DIR="${OPENCODE_DIR:-$HOME/.config/opencode}"
AGENTS_DST="$OPENCODE_DIR/agents"
SKILLS_DST="$OPENCODE_DIR/skills"
COMMANDS_DST="$OPENCODE_DIR/commands"

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
  version-code
  write-code
)

COMMAND_NAMES=(
  doc-plan
  get-plan
  kanban-card
  mr-review
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
install() { echo -e "        ${GREEN}↳ instalado${RESET}"; }

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
  local harness="$2"   # opencode
  local dst="$3"       # caminho final do .md
  {
    echo "---"
    cat "$dir/${harness}.yml"
    echo "---"
    echo ""
    cat "$dir/body.md"
  } > "$dst"
}

# substitui o token do modelo principal do OpenCode.
# agentes subagentes não carregam o token → no-op.
resolve_opencode_model() {
  local file="$1"
  local model="$2"
  sed -i "s|__OPENCODE_MAIN__|$model|g" "$file"
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

# ── seleção de vendor ─────────────────────────

VENDOR=""

select_vendor() {
  info "Selecione o vendor de modelos (Enter para manter o padrão: $OPENCODE_MAIN):"
  echo ""

  local i=1
  for vendor in "${VENDOR_NAMES[@]}"; do
    printf "        ${BOLD}%d)${RESET} %-18s  main: %s\n" \
      "$i" "$vendor" "${MODEL_MAIN[$vendor]}"
    i=$((i + 1))
  done

  echo ""

  local choice
  # lê uma única vez; stdin fechado/EOF/entrada vazia → mantém o default.
  if ! read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    Número do vendor [1-${#VENDOR_NAMES[@]}] (Enter = padrão): ")" choice; then
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
  ok "Vendor selecionado: ${BOLD}$VENDOR${RESET}"
  ok "  modelo principal : $OPENCODE_MAIN"
  echo ""
}

# ── argparse ─────────────────────────────────

FORCE=false
LOCAL=false
FETCH_CMD=""

for arg in "$@"; do
  case "$arg" in
    --force|-f) FORCE=true ;;
    --local|-l) LOCAL=true ;;
    --help|-h)
      echo ""
      echo -e "  ${BOLD}Uso:${RESET} install.sh [opções]"
      echo ""
      echo "  Opções:"
      echo "    --force, -f    Substituir todos os arquivos sem perguntar"
      echo "    --local, -l    Instalar a partir dos arquivos locais do repositório"
      echo "    --help,  -h    Exibir esta ajuda"
      echo ""
      exit 0
      ;;
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

# ── resolução de diretório fonte por tipo ─────
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

# prepara o diretório fonte de um agente/command (contém opencode.yml + body.md)
# uso: src_dir=$(prepare_assembled_src agents coder)
prepare_assembled_src() {
  local kind="$1"   # agents | commands
  local name="$2"
  if $LOCAL; then
    echo "$SCRIPT_DIR/$kind/$name"
    return
  fi
  local tmp="$REMOTE_TMP/$kind/$name"
  mkdir -p "$tmp"
  fetch_remote "$kind/$name/opencode.yml" "$tmp/opencode.yml"
  fetch_remote "$kind/$name/body.md" "$tmp/body.md"
  echo "$tmp"
}

# ── instalação ───────────────────────────────

echo ""
echo -e "${BOLD}  coder — OpenCode agents & skills${RESET}"
echo "  ─────────────────────────────────"
echo ""

select_vendor

mkdir -p "$AGENTS_DST" "$SKILLS_DST" "$COMMANDS_DST"

if ! $LOCAL; then
  REMOTE_TMP="$(mktemp -d)"
fi

# agentes
info "Instalando agentes em $AGENTS_DST"
echo ""
for name in "${AGENT_NAMES[@]}"; do
  echo -e "  ${BOLD}$name${RESET}"
  dst="$AGENTS_DST/$name.md"
  if ! check_overwrite "$dst" "$name"; then
    continue
  fi
  src_dir="$(prepare_assembled_src agents "$name")"
  assemble "$src_dir" opencode "$dst"
  resolve_opencode_model "$dst" "$OPENCODE_MAIN"
  install
done

echo ""

# skills (copia a pasta inteira)
info "Instalando skills em $SKILLS_DST"
echo ""
for name in "${SKILL_NAMES[@]}"; do
  echo -e "  ${BOLD}$name${RESET}"
  dst="$SKILLS_DST/$name"
  if ! check_overwrite "$dst" "$name"; then
    continue
  fi
  rm -rf "$dst"
  mkdir -p "$dst"
  if $LOCAL; then
    cp -R "$SCRIPT_DIR/skills/$name/." "$dst/"
  else
    # Fase 1: baixamos apenas SKILL.md. Skills com references/ não são
    # suportadas no modo remoto desta fase.
    fetch_remote "skills/$name/SKILL.md" "$dst/SKILL.md"
  fi
  install
done

echo ""

# commands
info "Instalando commands em $COMMANDS_DST"
echo ""
for name in "${COMMAND_NAMES[@]}"; do
  echo -e "  ${BOLD}$name${RESET}"
  dst="$COMMANDS_DST/$name.md"
  if ! check_overwrite "$dst" "$name"; then
    continue
  fi
  src_dir="$(prepare_assembled_src commands "$name")"
  assemble "$src_dir" opencode "$dst"
  resolve_opencode_model "$dst" "$OPENCODE_MAIN"
  install
done

echo ""
ok "Instalação concluída."
echo ""
echo -e "  Vendor : ${BOLD}$VENDOR${RESET}"
echo "  Modelo : $OPENCODE_MAIN"
echo ""
echo "  Agentes disponíveis:"
for name in "${AGENT_NAMES[@]}"; do
  echo "    • $name"
done
echo ""
echo "  Skills disponíveis:"
for name in "${SKILL_NAMES[@]}"; do
  echo "    • $name"
done
echo ""
echo "  Commands disponíveis:"
for name in "${COMMAND_NAMES[@]}"; do
  echo "    • /$name"
done
echo ""
echo -e "  Reinicie o OpenCode para carregar os novos agentes e skills."
echo ""
