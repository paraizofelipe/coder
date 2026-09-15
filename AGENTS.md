# AGENTS.md

## Sobre o repositório

Markdown puro: as skills e os commands do fluxo `planning` → `to-spec` → `implement`, para OpenCode, Claude Code e Oh My Pi (`omp`). Não há código executável, dependências, build ou testes. O único script é o `install.sh`.

Esta branch (`feat-new-flow`) não contém agentes: as sete skills rodam no agente ativo do harness.

## Estrutura

```text
skills/    7 subdiretórios — SKILL.md + references/ (opcional)
           planning · to-spec · to-cards · implement · tdd · code-review · to-memory
commands/  5 subdiretórios — body.md + opencode.yml + claude.yml + omp.yml
           to-spec · to-cards · implement · code-review · to-memory
install.sh                 — monta e copia para o diretório nativo de cada harness
```

## O fluxo

```text
/planning ──→ .coder/plan-<branch>.md ──→ /to-spec ──→ .coder/spec-<timestamp>.md
                                                          │
                                  /to-cards ←─────────────┤  (só ao distribuir)
                                      │                   │
                                  /implement ←────────────┘
                                      ├─ skill tdd ── ciclo de cada fatia
                                      ├─ skill code-review ── revisão final
                                      └─ /to-memory ── memória de longo prazo (opcional)
```

| Skill | Papel | Pode perguntar? |
|---|---|---|
| `planning` | Entrevista por rodadas até a árvore de decisões esvaziar. Investiga os fatos no repositório em vez de devolvê-los ao usuário | Sim |
| `to-spec` | Sintetiza as decisões em um spec de sete seções, confirmando antes as seams de teste | Só a confirmação das seams |
| `to-cards` | Quebra o spec em cards com arestas de bloqueio; prefactor é o card 01 | Sim — aprova a granularidade antes de gravar |
| `implement` | Corta o spec em fatias verticais e executa uma por vez, ou consome **um** card | Só diante de lacuna material — e então para |
| `tdd` | Ciclo vermelho → verde de cada fatia. Acionada pela `implement`, mas invocável sozinha | Só para confirmar seam não acordada |
| `code-review` | Revisa a diff contra um ponto fixo em dois eixos independentes: spec e padrões | Só o ponto fixo, quando não informado |
| `to-memory` | Envia os artefatos de `.coder/` crus à memória de longo prazo e mantém as páginas de conhecimento | Só a confirmação antes do envio |

Regra invariável: `planning` decide, `to-spec` registra, `implement` executa. O spec não reabre decisão e não inventa nenhuma; a implementação não decide o que ficou em aberto nem entrega o que ninguém pediu — afirmação ou código sem decisão por trás é defeito.

Use `to-spec` quando o trabalho atravessa várias sessões. Quando cabe em uma janela de contexto, vá de `planning` direto para `implement`.

## Convenções obrigatórias

- **Idioma:** todo conteúdo em português do Brasil; commits em inglês
- **Skills** (`SKILL.md`): frontmatter com `name` (kebab-case, igual ao nome da pasta) e `description` (o que faz e quando usar)
- **Progressive disclosure:** `SKILL.md` abaixo de 500 linhas; material extenso em `references/<arquivo>.md`, um nível só, carregado no passo que precisa dele
- **Commands:** frontmatter separado por harness — `opencode.yml`, `claude.yml` e `omp.yml` (este com `description` + `argument-hint`). Sem campo `agent:` nesta branch
- **XML tags** para estruturar conteúdo: `<role>`, `<context>`, `<workflow>`, `<rules>`, `<checklist>`, `<output_format>`
- **`<output_format>`** no final de cada skill define o contrato de resposta — nunca remover

## Artefatos do fluxo

| Arquivo | Origem | Observação |
|---|---|---|
| `.coder/plan-<branch-safe>.md` | `planning` | Um por branch (`/` vira `-`). Atualizado na mesma solicitação, com histórico de iterações. Solicitação diferente na mesma branch: a skill para e pergunta |
| `.coder/spec-AAAAMMDD-HHMMSS.md` | `to-spec` | Um por solicitação; ajuste na mesma sessão atualiza o mesmo arquivo |
| `.coder/cards-<branch-safe>/NN-slug.md` | `to-cards` | Um diretório por branch, um arquivo por card, numerado em ordem de dependência. Definição imutável depois de aprovada — nunca registra andamento |
| `.coder/impl-<branch-safe>.md` | `implement` | Um por branch, mesma regra de nome do plano. Atualizado ao fim de cada fatia; é o que permite retomar em sessão nova. **Única fonte de estado** |

`planning`, `to-spec`, `to-cards`, `code-review` e `to-memory` não escrevem código de produção. A `implement` escreve — e é a única —, mas **nenhuma** delas executa commit, push, branch, merge, rebase ou reset sem confirmação explícita do usuário. A `to-memory` também não escreve arquivo local: ela envia para fora, e todo envio é confirmado.

## Memória de longo prazo

`recall` na entrada do `planning`, `/to-memory` na saída. Ambos **opcionais**: a skill verifica a lista de ferramentas do agente antes de chamar, faz uma tentativa só, e trata ferramenta ausente, serviço fora e bank vazio como desfechos equivalentes.

Ausência não é falha, mas também não é silêncio: o cabeçalho do plano registra `Memória: indisponível`, porque a entrevista pode ter reaberto decisão já tomada em outra sessão.

O fallback são os planos já existentes em `.coder/` na máquina — **não** a memória nativa do harness, que captura sozinha e não deve ser acionada por um `SKILL.md`.

Enviar é irreversível e nunca resumido: a ingestão pede conteúdo bruto, e a síntese acontece no servidor. Páginas de conhecimento guardam **perguntas permanentes**, não respostas curadas — poucas e estáveis, nunca uma por feature.

## Onde cada harness encontra o fluxo

| Harness | Skills | Commands |
|---|---|---|
| OpenCode | `~/.config/opencode/skills/` | `~/.config/opencode/commands/` |
| Claude Code | `~/.claude/skills/` | `~/.claude/commands/` |
| Oh My Pi | `~/.agents/skills/` | `~/.agents/commands/` |

O OMP sabe ler os diretórios do Claude Code e do OpenCode, mas **só no nível de projeto** — `skills.enableClaudeUser`, `commands.enableClaudeUser` e `commands.enableOpencodeUser` vêm desligados. No nível de usuário, o único caminho ligado por padrão é `~/.agents`. Instalar no alvo `omp` é necessário, não opcional.

Invocação por harness: no Claude Code a skill já responde a `/<nome>`; no OMP ela responde a `/skill:<nome>` e o command dá o nome curto com `argument-hint`; no OpenCode a skill só é carregada pelo tool `skill`, escolhido pelo modelo — ali o command é a única porta do usuário. Por isso têm command as skills que o usuário dispara (`to-spec`, `to-cards`, `implement`, `code-review`, `to-memory`) e não a `tdd`, acionada pela `implement`.

Este `AGENTS.md` é documentação do repositório: o instalador não o copia para lugar nenhum.

## Validação

```bash
for s in planning to-spec to-cards implement tdd code-review to-memory; do npx -y skills-ref validate ./skills/$s; done
bash -n install.sh
```

O validador aceita só `allowed-tools`, `compatibility`, `description`, `license`, `metadata` e `name` no frontmatter. Campo de harness específico, como `disable-model-invocation`, é rejeitado.
