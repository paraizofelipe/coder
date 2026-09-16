#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────
# coder — instalador do fluxo planning → to-spec → implement
# (instala a partir do layout híbrido:
#   commands/<name>/{opencode.yml,claude.yml,omp.yml,body.md}
#   skills/<name>/SKILL.md[, references/])
#
# Harnesses suportados: OpenCode, Claude Code, Oh My Pi (omp), GitHub Copilot.
# Cada um recebe os artefatos no diretório e formato nativos. O Copilot recebe
# só as skills: a CLI dele já invoca a skill por /<nome> e não lê prompt file.
# ─────────────────────────────────────────────

# ATENÇÃO: aponta para a branch do fluxo novo. Ao integrar em main, voltar para /main.
REPO_URL="https://raw.githubusercontent.com/paraizofelipe/coder/feat-new-flow"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ── listas de nomes (não caminhos) ────────────

# ordem do fluxo, não alfabética
SKILL_NAMES=(
  planning
  to-spec
  to-cards
  implement
  tdd
  code-review
  to-memory
)

# referências (references/) por skill, baixadas no modo remoto
# (no modo local o `cp -R` já traz o diretório inteiro)
declare -A SKILL_REFERENCES=(
  [planning]="decision-tree.md plan-format.md"
  [to-spec]="plan-to-spec-map.md seams.md spec-format.md"
  [to-cards]="card-format.md decomposition.md"
  [implement]="commit-gate.md impl-format.md vertical-slices.md"
  [tdd]="mocking.md test-quality.md"
  [code-review]="spec-axis.md standards-axis.md"
  [to-memory]="artifact-map.md knowledge-pages.md redaction.md"
)

COMMAND_NAMES=(
  to-spec
  to-cards
  implement
  code-review
  to-memory
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
installed() { echo -e "        ${GREEN}↳ instalado${RESET}"; }

confirm() {
  local msg="$1"
  local answer
  read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    $msg [s/N] ")" answer </dev/tty
  [[ "$answer" =~ ^[sSyY]$ ]]
}

# ── montagem (commands) ───────────────────────
#
# Monta o arquivo final `.md` no formato:
#   ---
#   <conteúdo de <harness>.yml>
#   ---
#   <linha em branco>
#   <conteúdo de body.md>
#
assemble() {
  local dir="$1"       # diretório fonte (ex: commands/to-spec)
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

# verifica conflito + confirmação antes de sobrescrever.
# retorna 0 = pode prosseguir, 1 = pular.
# Opções interativas:
#   s = substitui apenas este item
#   n/Enter = pula este item
#   t = substitui este e todos os próximos conflitos
check_overwrite() {
  local dst="$1"
  local label="$2"
  if [[ -e "$dst" ]]; then
    warn "Já existe: $dst"
    if ! $FORCE && ! $OVERWRITE_ALL; then
      local answer
      if ! read -r -p "$(echo -e "${YELLOW}${BOLD}[?]${RESET}    Substituir $label? [s/N/todos] ")" answer </dev/tty; then
        answer=""
      fi
      case "$answer" in
        s|S|y|Y)
          ;;
        t|T|todos|Todos|TODOS|all|All|ALL)
          OVERWRITE_ALL=true
          ok "Todos os próximos conflitos serão sobrescritos."
          ;;
        *)
          skip
          return 1
          ;;
      esac
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
  echo "        3) omp"
  echo "        4) copilot"
  echo "        5) todos"
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
      3) HARNESSES+=(omp) ;;
      4) HARNESSES+=(copilot) ;;
      5) HARNESSES=(opencode claude omp copilot) ;;
    esac
  done
  if [[ ${#HARNESSES[@]} -eq 0 ]]; then
    warn "Nenhum harness selecionado."
    exit 1
  fi
  # dedup preservando ordem canônica: opencode → claude → omp → copilot
  local -a canonical=(opencode claude omp copilot)
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
  local -a canonical=(opencode claude omp copilot)
  local h token found
  for token in "${raw[@]}"; do
    if [[ "$token" == "all" ]]; then
      HARNESSES=(opencode claude omp copilot)
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
      echo -e "${RED}${BOLD}[erro]${RESET} harness inválido: '$token'. Use: opencode, claude, omp, copilot, all."
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
H_SKILLS=""; H_COMMANDS=""

harness_paths() {
  local h="$1"
  case "$h" in
    opencode)
      local base="${OPENCODE_DIR:-$HOME/.config/opencode}"
      H_SKILLS="$base/skills"; H_COMMANDS="$base/commands" ;;
    claude)
      local base="${CLAUDE_DIR:-$HOME/.claude}"
      H_SKILLS="$base/skills"; H_COMMANDS="$base/commands" ;;
    omp)
      # Oh My Pi descobre skills e slash commands pelo provider "Agent Dirs":
      # .agent/ e .agents/ no walk-up do projeto E no home do usuário.
      # ~/.agents é o destino nativo e segue o padrão Agent Skills.
      local base="${OMP_AGENTS_DIR:-$HOME/.agents}"
      H_SKILLS="$base/skills"; H_COMMANDS="$base/commands" ;;
    copilot)
      # COPILOT_HOME é a variável da própria CLI do Copilot, não uma inventada
      # aqui: inventar outra criaria dois lugares para apontar o mesmo diretório.
      # H_COMMANDS vazio porque o Copilot não lê command em markdown.
      local base="${COPILOT_HOME:-$HOME/.copilot}"
      H_SKILLS="$base/skills"; H_COMMANDS="" ;;
  esac
}

# ── argparse ─────────────────────────────────

FORCE=false
OVERWRITE_ALL=false
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
    --help|-h)
      echo ""
      echo -e "  ${BOLD}Uso:${RESET} install.sh [opções]"
      echo ""
      echo "  Instala as skills e os commands do fluxo planning → to-spec → implement"
      echo "  (OpenCode, Claude Code, Oh My Pi, GitHub Copilot), escolhidos antes da"
      echo "  instalação. O Copilot recebe só as skills: a CLI dele já invoca a skill"
      echo "  por /<nome> e não lê prompt file."
      echo ""
      echo "  Opções:"
      echo "    --force, -f              Substituir todos os arquivos sem perguntar"
      echo "    --local, -l              Instalar a partir dos arquivos locais do repositório"
      echo "    --harness <lista>        Harness(es) a instalar sem menu interativo."
      echo "                             Valores: opencode, claude, omp, copilot, all (ou"
      echo "                             combinações separadas por vírgula/espaço)"
      echo "    --help,  -h              Exibir esta ajuda"
      echo ""
      echo "  Overrides de diretório (env vars):"
      echo "    OPENCODE_DIR        base do OpenCode (default ~/.config/opencode)"
      echo "    CLAUDE_DIR          base do Claude Code (default ~/.claude)"
      echo "    OMP_AGENTS_DIR      base de agent dirs do Oh My Pi (default ~/.agents)"
      echo "    COPILOT_HOME        base do GitHub Copilot (default ~/.copilot)"
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
  local rel="$1"   # caminho relativo no repositório (ex: commands/to-spec/body.md)
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

# prepara o diretório fonte de um command (contém <harness>.yml + body.md).
# uso: src_dir=$(prepare_assembled_src commands to-spec opencode)
prepare_assembled_src() {
  local kind="$1"      # commands
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

# ── instalação por harness ────────────────────

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
      # modo remoto: baixa o SKILL.md e os arquivos de references/ declarados em SKILL_REFERENCES
      fetch_remote "skills/$name/SKILL.md" "$dst/SKILL.md"
      local refs="${SKILL_REFERENCES[$name]:-}"
      if [[ -n "$refs" ]]; then
        mkdir -p "$dst/references"
        local ref
        for ref in $refs; do
          fetch_remote "skills/$name/references/$ref" "$dst/references/$ref"
        done
      fi
    fi
    installed
  done
  echo ""
}

install_commands() {
  local harness="$1"
  # destino sem commands merece uma linha, não silêncio
  if [[ -z "$H_COMMANDS" ]]; then
    info "$harness não lê commands: a skill já responde a /<nome>"
    echo ""
    return 0
  fi
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
    installed
  done
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
    [[ -n "$H_COMMANDS" ]] && echo "      commands: $H_COMMANDS"
  done
  echo ""
  echo -e "  Reinicie o harness para carregar as novas skills e commands."
  echo ""
}

# ── fluxo principal ───────────────────────────

echo ""
echo -e "${BOLD}  coder — fluxo planning → to-spec → implement${RESET}"
echo "  ─────────────────────────────────────"
echo ""

if [[ -n "$FLAG_HARNESS" ]]; then
  resolve_harness_flag "$FLAG_HARNESS"
else
  select_harness
fi

if ! $LOCAL; then
  REMOTE_TMP="$(mktemp -d)"
fi

for h in "${HARNESSES[@]}"; do
  harness_paths "$h"
  info "Instalando para: ${BOLD}$h${RESET}"
  echo ""
  install_skills
  install_commands "$h"
  ok "Concluído: $h"
  echo ""
done

print_summary
