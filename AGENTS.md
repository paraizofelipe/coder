# AGENTS.md

## Sobre o repositório

Markdown puro: as skills e o command do fluxo `planning` → `to-spec`, para OpenCode, Claude Code, Codex e Pi. Não há código executável, dependências, build ou testes. O único script é o `install.sh`.

Esta branch (`feat-new-flow`) não contém agentes: as duas skills rodam no agente ativo do harness.

## Estrutura

```text
skills/    2 subdiretórios — SKILL.md + references/ (opcional)
commands/  1 subdiretório  — body.md + opencode.yml + claude.yml + pi.yml
install.sh                 — monta e copia para o diretório nativo de cada harness
```

## O fluxo

```text
/planning  ──→  .coder/plan.md  ──→  /to-spec  ──→  .coder/spec-<timestamp>.md
```

| Skill | Papel | Pode perguntar? |
|---|---|---|
| `planning` | Entrevista por rodadas até a árvore de decisões esvaziar. Investiga os fatos no repositório em vez de devolvê-los ao usuário | Sim |
| `to-spec` | Sintetiza as decisões em um spec de sete seções, confirmando antes as seams de teste | Só a confirmação das seams |

Regra invariável: `planning` decide, `to-spec` registra. O spec não reabre decisão e não inventa nenhuma — afirmação que ninguém decidiu é defeito.

Use `to-spec` quando o trabalho atravessa várias sessões. Quando cabe em uma janela de contexto, implemente direto.

## Convenções obrigatórias

- **Idioma:** todo conteúdo em português do Brasil; commits em inglês
- **Skills** (`SKILL.md`): frontmatter com `name` (kebab-case, igual ao nome da pasta) e `description` (o que faz e quando usar)
- **Progressive disclosure:** `SKILL.md` abaixo de 500 linhas; material extenso em `references/<arquivo>.md`, um nível só, carregado no passo que precisa dele
- **Commands:** frontmatter separado por harness — `opencode.yml`, `claude.yml` e `pi.yml` (este com `argument-hint`). Sem campo `agent:` nesta branch
- **XML tags** para estruturar conteúdo: `<role>`, `<context>`, `<workflow>`, `<rules>`, `<checklist>`, `<output_format>`
- **`<output_format>`** no final de cada skill define o contrato de resposta — nunca remover

## Artefatos do fluxo

| Arquivo | Origem | Observação |
|---|---|---|
| `.coder/plan.md` | `planning` | Atualizado na mesma solicitação, com histórico de iterações |
| `.coder/spec-AAAAMMDD-HHMMSS.md` | `to-spec` | Um por solicitação; ajuste na mesma sessão atualiza o mesmo arquivo |

Nenhuma das duas skills escreve código de produção, cria branch ou commita.

## Validação

```bash
npx -y skills-ref validate ./skills/to-spec
bash -n install.sh
```
