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
  "agents/tech_reviewer.md"
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
  read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    $msg [s/N] ")" answer
  [[ "$answer" =~ ^[sSyY]$ ]]
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

mkdir -p "$AGENTS_DST" "$SKILLS_DST"

install_file() {
  local src="$1"       # caminho relativo (ex: agents/analyzer.md)
  local dst_dir="$2"   # diretório de destino
  local filename
  filename="$(basename "$src")"
  local dst="$dst_dir/$filename"

  echo -e "  ${BOLD}$filename${RESET}"

  # verifica se já existe
  if [[ -f "$dst" ]]; then
    warn "Já existe: $dst"
    if ! $FORCE; then
      if ! confirm "Substituir $filename?"; then
        skip
        return
      fi
    fi
  fi

  # copia ou baixa
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

  install
}

# agentes
info "Instalando agentes em $AGENTS_DST"
echo ""
for agent in "${AGENTS[@]}"; do
  install_file "$agent" "$AGENTS_DST"
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
