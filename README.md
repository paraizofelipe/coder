# coder — fluxo `planning` → `to-spec`

Branch isolada com o fluxo de especificação: uma entrevista que decide (`planning`) e uma síntese que registra (`to-spec`). Nada mais.

Markdown puro mais um shell script. Sem código executável, dependências ou build.

## Índice

- [O fluxo](#o-fluxo)
- [Estrutura do repositório](#estrutura-do-repositório)
- [Skills](#skills)
- [Commands](#commands)
- [Artefatos gerados](#artefatos-gerados)
- [Instalação](#instalação)
- [Diretórios de instalação](#diretórios-de-instalação)
- [Opções do instalador](#opções-do-instalador)
- [Convenções do projeto](#convenções-do-projeto)

## O fluxo

```text
/planning  ─── entrevista até a árvore de decisões esvaziar ──→  .coder/plan.md
    │
    └─→ /to-spec  ─── sintetiza, confirma as seams, não entrevista ──→  .coder/spec-<timestamp>.md
```

| Etapa | Papel | Pode perguntar? |
|---|---|---|
| `planning` | **Decide.** Entrevista por rodadas até nenhuma decisão material ficar implícita | Sim — é a razão de existir |
| `to-spec` | **Registra.** Converte as decisões em um spec que sobrevive ao fim da janela de contexto | Só a confirmação das seams |

A separação é o ponto: quem decide não escreve o documento final, e quem escreve não reabre decisão. Um spec que afirma algo que ninguém decidiu é defeito, não iniciativa.

**Quando usar o `to-spec`:** quando o trabalho atravessa várias sessões. Se cabe em uma janela de contexto, o spec não paga o próprio custo.

## Estrutura do repositório

```text
skills/
  planning/
    SKILL.md
  to-spec/
    SKILL.md
    references/
      plan-to-spec-map.md    mapeia .coder/plan.md → seções do spec
      seams.md               conceito, heurísticas e anti-padrões
      spec-format.md         template do .coder/spec-*.md
commands/
  to-spec/
    body.md
    opencode.yml
    claude.yml
    pi.yml
install.sh
```

Esta branch não contém agentes. As duas skills rodam no agente ativo do harness.

## Skills

Formato [Agent Skills](https://agentskills.io/specification): uma pasta por skill, `SKILL.md` obrigatório, `references/` opcional carregado sob demanda.

| Skill | O que faz | Saída |
|---|---|---|
| `planning` | Entrevista rigorosa: investiga os fatos no repositório, modela a árvore de decisões e pergunta por rodadas até a fronteira esvaziar | `.coder/plan.md` |
| `to-spec` | Sintetiza as decisões já tomadas em um spec de sete seções, confirmando antes as seams de teste | `.coder/spec-AAAAMMDD-HHMMSS.md` |

### Progressive disclosure no `to-spec`

O `SKILL.md` tem 93 linhas e nunca carrega o template nem a teoria de seams. Cada reference entra no passo que precisa dela:

| Passo | Carrega |
|---|---|
| 1 — reunir decisões | `references/plan-to-spec-map.md` |
| 3 — propor seams | `references/seams.md` |
| 4 — escrever | `references/spec-format.md` |

Validar antes de commitar:

```bash
npx -y skills-ref validate ./skills/to-spec
```

## Commands

| Command | Harnesses | Efeito |
|---|---|---|
| `/to-spec [foco]` | OpenCode, Claude Code, Codex, Pi | Executa a skill `to-spec` sobre o `.coder/plan.md` ou a conversa atual |

A skill `planning` não tem command: invoque-a pelo nome (`/planning`, quando o harness expõe skills como slash commands) ou peça o planejamento em linguagem natural.

## Artefatos gerados

Ambos ficam em `.coder/`, no repositório onde o fluxo roda.

| Arquivo | Origem | Ciclo de vida |
|---|---|---|
| `.coder/plan.md` | `planning` | Atualizado na mesma solicitação, com histórico de iterações |
| `.coder/spec-AAAAMMDD-HHMMSS.md` | `to-spec` | Um por solicitação; ajustes na mesma sessão atualizam o mesmo arquivo |

## Instalação

### Via curl

```bash
curl -fsSL https://raw.githubusercontent.com/paraizofelipe/coder/feat-new-flow/install.sh | bash
```

### A partir do repositório local

```bash
git clone https://github.com/paraizofelipe/coder.git
cd coder
git switch feat-new-flow
./install.sh --local
```

### Seleção de harness

```bash
./install.sh --harness claude            # só Claude Code
./install.sh --harness opencode,claude   # dois harnesses
./install.sh --harness all               # todos
```

Sem a flag, o instalador exibe o menu interativo.

## Diretórios de instalação

| Harness | Skills | Commands / prompts | AGENTS.md |
|---|---|---|---|
| OpenCode | `~/.config/opencode/skills/` | `~/.config/opencode/commands/` | — |
| Claude Code | `~/.claude/skills/` | `~/.claude/commands/` | — |
| Codex | `~/.agents/skills/` | `~/.codex/prompts/` (body-only) | `~/.codex/AGENTS.md` |
| Pi | `~/.agents/skills/` | `~/.pi/agent/prompts/` (montado com `pi.yml`) | `~/.pi/agent/AGENTS.md` |

Codex e Pi compartilham `~/.agents/skills` de propósito: evita colisão de nomes quando os dois coexistem.

Overrides por variável de ambiente: `OPENCODE_DIR`, `CLAUDE_DIR`, `CODEX_DIR`, `CODEX_SKILLS_DIR`, `PI_DIR`, `PI_SKILLS_DIR`.

## Opções do instalador

| Flag | Efeito |
|---|---|
| `--force`, `-f` | Substitui tudo sem perguntar |
| `--local`, `-l` | Instala dos arquivos locais em vez de baixar do GitHub |
| `--harness <lista>` | `opencode`, `claude`, `codex`, `pi`, `all` — vírgula ou espaço |
| `--help`, `-h` | Ajuda |

Em conflito, o prompt aceita `s` (só este), `n`/Enter (pular) e `todos` (este e todos os seguintes).

Não há seleção de vendor nesta branch: ela só existia para injetar modelo nos agentes primários, que não estão aqui.

## Convenções do projeto

| Item | Regra |
|---|---|
| Idioma | Português do Brasil em todo conteúdo; commits em inglês |
| Commits | Conventional Commits: `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`, `test:` |
| Skills | Pasta kebab-case igual ao campo `name`; `SKILL.md` abaixo de 500 linhas |
| References | Um nível só (`references/<arquivo>.md`); material extenso sai do `SKILL.md` |
| XML tags | `<role>`, `<context>`, `<workflow>`, `<rules>`, `<checklist>`, `<output_format>` |
| `<output_format>` | Obrigatório — é o contrato de resposta da skill |

O `install.sh` desta branch aponta o `REPO_URL` para `feat-new-flow`. Ao integrar em `main`, voltar para `/main`.
