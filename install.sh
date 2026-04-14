#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────
# coder — OpenCode agents & skills installer
# ─────────────────────────────────────────────

REPO_URL="https://raw.githubusercontent.com/paraizofelipe/coder/main"
OPENCODE_DIR="${OPENCODE_DIR:-$HOME/.opencode}"
AGENTS_DST="$OPENCODE_DIR/agents"
SKILLS_DST="$OPENCODE_DIR/skills"

AGENTS=(
  "agents/analyzer.md"
  "agents/business_reviewer.md"
  "agents/coder.md"
  "agents/code_reviewer.md"
  "agents/tester.md"
  "agents/versioner.md"
)

SKILLS=(
  "skills/analyse_code.md"
  "skills/review_code.md"
  "skills/test_code.md"
  "skills/version_code.md"
  "skills/write_code.md"
)

# agentes que usam modelo light (tarefas simples)
LIGHT_AGENTS=("agents/versioner.md")

# ── mapa de vendors e modelos ─────────────────

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
  [openai]="openai/gpt-5.3-codex"
  [google]="google/gemini-3.1-pro"
  [groq]="groq/llama-3.3-70b-versatile"
  [amazon-bedrock]="amazon-bedrock/amazon.nova-pro-v1"
  [github-copilot]="github-copilot/claude-sonnet-4-5"
)

declare -A MODEL_LIGHT=(
  [anthropic]="anthropic/claude-haiku-4-5"
  [openai]="openai/codex-mini-latest"
  [google]="google/gemini-3-flash"
  [groq]="groq/llama-3.1-8b-instant"
  [amazon-bedrock]="amazon-bedrock/amazon.nova-lite-v1"
  [github-copilot]="github-copilot/gpt-4o-mini"
)

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

# ── seleção de vendor ─────────────────────────

VENDOR=""

select_vendor() {
  info "Selecione o vendor de modelos:"
  echo ""

  local i=1
  for vendor in "${VENDOR_NAMES[@]}"; do
    printf "        ${BOLD}%d)${RESET} %-18s  main: %-42s  light: %s\n" \
      "$i" "$vendor" "${MODEL_MAIN[$vendor]}" "${MODEL_LIGHT[$vendor]}"
    i=$((i + 1))
  done

  echo ""

  local choice
  while true; do
    read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    Número do vendor [1-${#VENDOR_NAMES[@]}]: ")" choice </dev/tty
    if [[ "$choice" =~ ^[0-9]+$ ]] && (( choice >= 1 && choice <= ${#VENDOR_NAMES[@]} )); then
      VENDOR="${VENDOR_NAMES[$((choice - 1))]}"
      break
    fi
    echo -e "        ${RED}Opção inválida. Digite um número entre 1 e ${#VENDOR_NAMES[@]}.${RESET}"
  done

  echo ""
  ok "Vendor selecionado: ${BOLD}$VENDOR${RESET}"
  ok "  modelo principal : ${MODEL_MAIN[$VENDOR]}"
  ok "  modelo versioner : ${MODEL_LIGHT[$VENDOR]}"
  echo ""
}

# ── argparse ─────────────────────────────────

FORCE=false
LOCAL=false

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

# ── instalação ───────────────────────────────

echo ""
echo -e "${BOLD}  coder — OpenCode agents & skills${RESET}"
echo "  ─────────────────────────────────"
echo ""

select_vendor

mkdir -p "$AGENTS_DST" "$SKILLS_DST"

is_light_agent() {
  local src="$1"
  for light in "${LIGHT_AGENTS[@]}"; do
    [[ "$src" == "$light" ]] && return 0
  done
  return 1
}

patch_model() {
  local file="$1"
  local model="$2"
  local tmp
  tmp=$(mktemp)
  sed "s|^model:.*|model: $model|" "$file" > "$tmp" && mv "$tmp" "$file"
}

install_file() {
  local src="$1"       # caminho relativo (ex: agents/analyzer.md)
  local dst_dir="$2"   # diretório de destino
  local model="${3:-}" # modelo a substituir (vazio = sem substituição)
  local filename
  filename="$(basename "$src")"
  local dst="$dst_dir/$filename"

  echo -e "  ${BOLD}$filename${RESET}"

  if [[ -f "$dst" ]]; then
    warn "Já existe: $dst"
    if ! $FORCE; then
      if ! confirm "Substituir $filename?"; then
        skip
        return
      fi
    fi
  fi

  if $LOCAL; then
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    cp "$script_dir/$src" "$dst"
  else
    if [[ "$FETCH_CMD" == "curl" ]]; then
      curl -fsSL "$REPO_URL/$src" -o "$dst"
    else
      wget -qO "$dst" "$REPO_URL/$src"
    fi
  fi

  if [[ -n "$model" ]]; then
    patch_model "$dst" "$model"
  fi

  install
}

# agentes
info "Instalando agentes em $AGENTS_DST"
echo ""
for agent in "${AGENTS[@]}"; do
  if is_light_agent "$agent"; then
    install_file "$agent" "$AGENTS_DST" "${MODEL_LIGHT[$VENDOR]}"
  else
    install_file "$agent" "$AGENTS_DST" "${MODEL_MAIN[$VENDOR]}"
  fi
done

echo ""

# skills
info "Instalando skills em $SKILLS_DST"
echo ""
for skill in "${SKILLS[@]}"; do
  install_file "$skill" "$SKILLS_DST"
done

echo ""
ok "Instalação concluída."
echo ""
echo "  Vendor : ${BOLD}$VENDOR${RESET}"
echo "  Modelos: ${MODEL_MAIN[$VENDOR]} / ${MODEL_LIGHT[$VENDOR]}"
echo ""
echo "  Agentes disponíveis:"
for agent in "${AGENTS[@]}"; do
  echo "    • $(basename "$agent" .md)"
done
echo ""
echo "  Skills disponíveis:"
for skill in "${SKILLS[@]}"; do
  echo "    • $(basename "$skill" .md)"
done
echo ""
echo -e "  Reinicie o OpenCode para carregar os novos agentes e skills."
echo ""
