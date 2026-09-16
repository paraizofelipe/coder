# AGENTS.md

## Sobre o repositório

As skills e os commands do fluxo `planning` → `to-spec` → `implement`, para OpenCode, Claude Code, Oh My Pi (`omp`) e GitHub Copilot — este último recebe só as skills. O conteúdo é markdown puro; o instalador é um binário Go na raiz que o embute com `go:embed`. O Go existe para distribuir o markdown — o fluxo em si não depende dele.

Esta branch (`feat-new-flow`) não contém agentes: as sete skills rodam no agente ativo do harness.

## Estrutura

```text
skills/    7 subdiretórios — SKILL.md + references/ (opcional)
           planning · to-spec · to-cards · implement · tdd · code-review · to-memory
commands/  5 subdiretórios — body.md + opencode.yml + claude.yml + omp.yml
           to-spec · to-cards · implement · code-review · to-memory
manifest.toml  lista do que é instalado, na ordem do fluxo
*.go       instalador — main · embed · harness · project · assemble · installer
           tree · banner · output · prompt
*_test.go  manifesto vs conteúdo embutido, montagem, gate de sobrescrita
install.sh instalador anterior, em bash; saída idêntica à do binário
Dockerfile ambiente descartável de teste, com os quatro harnesses instalados
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

| Harness | Skills (global) | Commands (global) | Raiz no projeto |
|---|---|---|---|
| OpenCode | `~/.config/opencode/skills/` | `~/.config/opencode/commands/` | `.opencode/` |
| Claude Code | `~/.claude/skills/` | `~/.claude/commands/` | `.claude/` |
| Oh My Pi | `~/.agents/skills/` | `~/.agents/commands/` | `.agents/` |
| GitHub Copilot | `~/.copilot/skills/` | — não lê command em markdown | `.github/` |

Os harnesses leem os diretórios uns dos outros: `.claude/` no projeto é lido pelos quatro; `~/.claude` também pelo OpenCode; `~/.agents` pelo OpenCode e pelo Copilot. Só `~/.config/opencode/`, `~/.copilot/` e `.github/` são exclusivos. `coder --status` lista o que está instalado e aponta cópia ambígua.

A instalação é **global ou de projeto**, nunca as duas na mesma execução. Chamado de dentro de um repositório git, o instalador pergunta; fora dele, vai de global sem perguntar. O Copilot é o único em que o nome do diretório muda entre os escopos: no usuário lê `~/.copilot`, no repositório lê `.github`.

O OMP sabe ler os diretórios do Claude Code e do OpenCode, mas **só no nível de projeto** — `skills.enableClaudeUser`, `commands.enableClaudeUser` e `commands.enableOpencodeUser` vêm desligados. No nível de usuário, o único caminho ligado por padrão é `~/.agents`. Instalar no alvo `omp` é necessário, não opcional — **no escopo global**. Instalando no projeto, `.claude/` e `.opencode/` já são lidos pelo OMP por padrão, e é o sentido inverso que passa a valer.

Invocação por harness: no Claude Code e no Copilot a skill já responde a `/<nome>`; no OMP ela responde a `/skill:<nome>` e o command dá o nome curto com `argument-hint`; no OpenCode a skill só é carregada pelo tool `skill`, escolhido pelo modelo — ali o command é a única porta do usuário. Por isso têm command as skills que o usuário dispara (`to-spec`, `to-cards`, `implement`, `code-review`, `to-memory`) e não a `tdd`, acionada pela `implement`.

Este `AGENTS.md` é documentação do repositório: o instalador não o copia para lugar nenhum. O mesmo vale para o `Dockerfile`, que serve só ao container de teste — `docker build -t coder-test . && docker run --rm -it coder-test`, sem volume, com o binário compilado dentro da imagem.

## Validação

```bash
for s in planning to-spec to-cards implement tdd code-review to-memory; do npx -y skills-ref validate ./skills/$s; done
timeout 180s env CI=true go test ./...
bash -n install.sh
```

Skill nova exige uma linha em `manifest.toml`: sem ela, `TestManifestoCobreExatamenteOQueFoiEmbutido` falha, e a skill seria embutida no binário sem nunca ser instalada.

A CLI do Copilot também varre `~/.agents/skills` — o destino do alvo `omp`. Quem instala no `omp` já alimenta o Copilot; o alvo `copilot` existe para quem não usa o OMP.

O validador aceita só `allowed-tools`, `compatibility`, `description`, `license`, `metadata` e `name` no frontmatter. Campo de harness específico, como `disable-model-invocation`, é rejeitado.
