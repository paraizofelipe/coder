# CLAUDE.md

Guia para o Claude Code (claude.ai/code) ao trabalhar neste repositório.

## O que é esta branch

`feat-new-flow` isola o fluxo de especificação `planning` → `to-spec`. Contém apenas essas duas skills, o command `/to-spec` e o instalador correspondente. Os agentes e as demais skills do `coder` vivem na `main` e **não** devem ser trazidos para cá.

Markdown puro mais um shell script. Sem código executável, dependências, build ou testes automatizados.

## Estrutura

```text
skills/
  planning/SKILL.md
  to-spec/SKILL.md + references/{plan-to-spec-map,seams,spec-format}.md
commands/
  to-spec/{body.md,opencode.yml,claude.yml,pi.yml}
install.sh
```

## O fluxo

| Etapa | Papel | Artefato |
|---|---|---|
| `planning` | Entrevista por rodadas até a árvore de decisões esvaziar. Investiga os fatos no repositório em vez de devolvê-los ao usuário | `.coder/plan.md` |
| `to-spec` | Sintetiza as decisões em um spec de sete seções. Não entrevista; a única pergunta permitida é a confirmação das seams | `.coder/spec-AAAAMMDD-HHMMSS.md` |

A fronteira entre as duas é o contrato do fluxo: `planning` decide, `to-spec` registra. Qualquer afirmação no spec que ninguém decidiu é defeito.

`to-spec` só se paga quando o trabalho atravessa várias sessões. Em mudança que cabe em uma janela de contexto, o spec é custo sem retorno.

## Convenções ao editar

### Skills

Padrão [Agent Skills](https://agentskills.io/specification):

- Uma pasta por skill, nome em kebab-case **igual** ao campo `name` do frontmatter
- `SKILL.md` abaixo de 500 linhas; material extenso vai para `references/`
- `references/` com um nível só (`references/<arquivo>.md`), carregado sob demanda pelo passo que precisa dele
- Frontmatter: `name` e `description` obrigatórios. A `description` diz **o que faz e quando usar**

Validar antes de commitar:

```bash
npx -y skills-ref validate ./skills/<nome>
```

### Commands

`commands/<name>/` com `body.md` mais um `.yml` por harness (`opencode.yml`, `claude.yml`, `pi.yml`). O `install.sh` monta frontmatter + corpo. Commands desta branch não declaram `agent:` — não há agentes aqui, então rodam no agente ativo.

### Estrutura XML

Tags em uso: `<role>`, `<context>`, `<workflow>`, `<rules>`, `<checklist>`, `<output_format>`.

Não remover `<output_format>` de nenhuma skill — é o contrato de resposta.

### Idioma

Português do Brasil em todo conteúdo. Commits em inglês, Conventional Commits.

## install.sh

Copia `skills/` e `commands/` para o diretório nativo de cada harness selecionado (OpenCode, Claude Code, Codex, Pi) e monta cada command juntando `<harness>.yml` + `body.md`. Skills são copiadas como diretórios completos, preservando `references/`.

Diferenças em relação ao instalador da `main`:

| Removido | Motivo |
|---|---|
| `AGENT_NAMES` e `install_agents()` | Não há agentes nesta branch |
| Seleção de vendor (`--vendor`, `MODEL_MAIN`, `apply_model`) | Só existia para injetar modelo nos agentes primários |

Flags: `--force`, `--local`, `--harness <lista>`, `--help`.

No modo remoto, o `REPO_URL` aponta para a branch `feat-new-flow`. **Ao integrar em `main`, voltar para `/main`** — senão o instalador continuará baixando de uma branch que pode não existir mais.

## O que não fazer

- Não trazer agentes, skills ou commands da `main` para esta branch
- Não adicionar código executável, dependências ou configuração de build
- Não remover `<output_format>` de nenhuma skill
- Não deixar o `to-spec` fazer perguntas além da confirmação das seams
