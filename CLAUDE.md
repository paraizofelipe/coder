# CLAUDE.md

Guia para o Claude Code (claude.ai/code) ao trabalhar neste repositório.

## O que é esta branch

`feat-new-flow` isola o fluxo `planning` → `to-spec` → `implement`. Contém apenas essas skills, mais `to-cards` para distribuir, `tdd` e `code-review` que a `implement` aciona, `to-memory` como epílogo opcional, os commands correspondentes e o instalador. Os agentes e as demais skills do `coder` vivem na `main` e **não** devem ser trazidos para cá.

O conteúdo é markdown puro. O instalador é um binário Go na raiz do repositório, que embute `skills/` e `commands/` com `go:embed` — é o único código executável, e existe para distribuir o markdown, não para participar do fluxo.

## Estrutura

```text
skills/
  planning/SKILL.md    + references/{decision-tree,plan-format}.md
  to-spec/SKILL.md     + references/{plan-to-spec-map,seams,spec-format}.md
  to-cards/SKILL.md    + references/{decomposition,card-format}.md
  implement/SKILL.md   + references/{impl-format,vertical-slices,commit-gate}.md
  tdd/SKILL.md         + references/{test-quality,mocking}.md
  code-review/SKILL.md + references/{spec-axis,standards-axis}.md
  to-memory/SKILL.md   + references/{artifact-map,redaction,knowledge-pages}.md
commands/
  to-spec/{body.md,opencode.yml,claude.yml,omp.yml}
  to-cards/{...}
  implement/{...}
  code-review/{...}
  to-memory/{...}
main.go · embed.go · harness.go · assemble.go · installer.go
tree.go · banner.go · output.go · prompt.go
*_test.go
install.sh          instalador anterior, em bash
```

## O fluxo

| Etapa | Papel | Artefato |
|---|---|---|
| `planning` | Entrevista por rodadas até a árvore de decisões esvaziar. Investiga os fatos no repositório em vez de devolvê-los ao usuário | `.coder/plan-<branch-safe>.md` |
| `to-spec` | Sintetiza as decisões em um spec de sete seções. Não entrevista; a única pergunta permitida é a confirmação das seams | `.coder/spec-AAAAMMDD-HHMMSS.md` |
| `to-cards` | Quebra o spec em cards com arestas de bloqueio, prefactor como card 01. Só grava depois de o usuário aprovar a granularidade | `.coder/cards-<branch-safe>/NN-slug.md` |
| `implement` | Corta o spec em fatias verticais e executa uma por vez, ou consome **um** card, na seam acordada. Não decide o que ficou em aberto: para | código, testes e `.coder/impl-<branch-safe>.md` |
| `tdd` | Ciclo vermelho → verde de cada fatia. Acionada pela `implement`, mas invocável sozinha | teste + implementação da fatia |
| `code-review` | Revisa a diff contra um ponto fixo em dois eixos independentes — spec e padrões. Aponta, não corrige | relatório na conversa |
| `to-memory` | Envia os artefatos de `.coder/` **crus** ao serviço de memória de longo prazo e mantém as páginas de conhecimento | documentos e páginas no bank |

A fronteira entre elas é o contrato do fluxo: `planning` decide, `to-spec` registra, `implement` executa. Qualquer afirmação no spec que ninguém decidiu é defeito — e o mesmo vale para código que o spec não pediu.

A `implement` nunca commita sozinha: prepara o stage arquivo a arquivo, revisa a própria diff, monta a mensagem e **aguarda confirmação explícita**. O mesmo vale para push, branch, merge, rebase e reset, e para todo envio da `to-memory`.

## Cards

`to-cards` só entra quando o trabalho sai de uma pessoa — vários devs, ou muitas sessões com `/clear`. Uma pessoa numa sessão não usa cards: o `implement` corta internamente.

| Regra | Razão |
|---|---|
| Aresta é bloqueio real | Preferência de ordem serializa trabalho paralelizável sem ninguém perceber |
| Um card cabe em uma janela nova | Quem pegar o card não precisa da conversa que o gerou |
| Prefactor é o card 01, ou não existe | Preparação depois da mudança é conserto, não prefactor |
| Sem path e sem snippet | O card pode esperar semanas na fila |
| Card é definição; `impl-<branch>.md` é estado | Duas fontes de estado divergem em uma semana |
| Nada é gravado antes da aprovação | Card publicado é fila para outras pessoas |

As regras de corte moram em `implement/references/vertical-slices.md` e **não são duplicadas** no `to-cards` — ele aponta para lá, do mesmo jeito que o `implement` aponta para a skill `tdd`.

Cards não vão para o `to-memory`: são descartáveis por construção.

## Memória de longo prazo

Dois pontos de contato, ambos **opcionais e tolerantes à ausência**: `recall` na entrada do `planning` e `/to-memory` na saída.

| Regra | Razão |
|---|---|
| **Detecção, não tentativa** | A skill verifica a lista de ferramentas do agente antes de chamar — mesmo padrão já usado para escolher entre `AskUserQuestion`, `question` e `ask_user` |
| **Uma tentativa, sem retentativa** | Ferramenta ausente, serviço fora e bank vazio são desfechos equivalentes: siga sem memória |
| **Degradar não é silenciar** | O cabeçalho do plano registra `Memória: indisponível`, porque a entrevista pode ter reaberto decisão já tomada |
| **Nunca resumir antes de enviar** | A ingestão pede conteúdo bruto; a síntese acontece no servidor |
| **Página guarda pergunta, não resposta** | Poucas e estáveis. Página por feature é erro de camada — isso é documento |
| **Envio é irreversível** | A superfície exposta não tem exclusão de documento. Por isso a confirmação nomeia o bank de destino |

O fallback, sem serviço, são os planos que já existem em `.coder/` na máquina — **não** a memória nativa do harness, que segue capturando sozinha e não deve ser acionada por um `SKILL.md`. Embutir caminho de um harness em conteúdo que roda em três é o erro que essa branch inteira evita.

O plano leva a branch no nome porque é lido de volta pelo `/to-spec`: `.coder/` está fora do versionamento, então um caminho fixo seria um arquivo único no disco, sobrescrito pelo próximo `/planning` de outra branch — e o `/to-spec` seguiria lendo, produzindo um spec plausível a partir do plano errado. O andamento da `implement` segue a mesma regra, e é lido de volta pela própria skill na sessão seguinte. O spec não é lido de volta por convenção de nome — o caminho dele é passado —, então o timestamp basta; a branch de origem vai no cabeçalho.

`to-spec` só se paga quando o trabalho atravessa várias sessões. Em mudança que cabe em uma janela de contexto, o spec é custo sem retorno: vá de `planning` direto para `implement`.

## Convenções ao editar

### Skills

Padrão [Agent Skills](https://agentskills.io/specification):

- Uma pasta por skill, nome em kebab-case **igual** ao campo `name` do frontmatter
- `SKILL.md` abaixo de 500 linhas; material extenso vai para `references/`
- `references/` com um nível só (`references/<arquivo>.md`), carregado sob demanda pelo passo que precisa dele
- Frontmatter: `name` e `description` obrigatórios. A `description` diz **o que faz e quando usar**

O frontmatter aceita apenas `allowed-tools`, `compatibility`, `description`, `license`, `metadata` e `name`. Campos de harness específico — `disable-model-invocation`, entre outros — são rejeitados pelo validador; a restrição de invocação, quando necessária, mora dentro da skill.

Validar antes de commitar:

```bash
npx -y skills-ref validate ./skills/<nome>
timeout 180s env CI=true go test ./...
```

### Commands

`commands/<name>/` com `body.md` mais um `.yml` por harness que lê commands (`opencode.yml`, `claude.yml`, `omp.yml`). O instalador monta frontmatter + corpo. O Copilot não entra nessa lista: não lê command em markdown. Commands desta branch não declaram `agent:` — não há agentes aqui, então rodam no agente ativo.

Tem command a skill que o **usuário** dispara: `to-spec`, `implement` e `code-review`. A `tdd` não tem, porque quem a aciona é a `implement`. A `planning` também não, por não precisar de argumento.

O `omp.yml` usa `description` e `argument-hint` — os dois campos que o Oh My Pi lê no frontmatter de slash command.

### Estrutura XML

Tags em uso: `<role>`, `<context>`, `<workflow>`, `<rules>`, `<checklist>`, `<output_format>`.

Não remover `<output_format>` de nenhuma skill — é o contrato de resposta.

### Idioma

Português do Brasil em todo conteúdo. Commits em inglês, Conventional Commits.

## O instalador

Binário Go na raiz. Copia `skills/` e `commands/` para o diretório nativo de cada harness selecionado (OpenCode, Claude Code, Oh My Pi, GitHub Copilot) e monta cada command juntando `<harness>.yml` + `body.md`. O Copilot é a exceção: recebe só as skills. Skills são copiadas como diretórios completos, preservando `references/`.

| Arquivo | Responsabilidade |
|---|---|
| `main.go` | Flags, orquestração, ajuda |
| `manifest.toml` | O que é instalado, e em que ordem — dado, não código |
| `embed.go` | `go:embed` do conteúdo e leitura do manifesto |
| `harness.go` | Destinos, overrides de diretório, parsing de `--harness` |
| `assemble.go` | Frontmatter do harness + `body.md` no `.md` final |
| `installer.go` | Cópia das skills e o gate de sobrescrita |
| `banner.go` | A arte ASCII e a largura calculada dela |
| `output.go` | Cores, detecção de terminal, banner e resumo |
| `prompt.go` | Formulários `huh`: seleção de harness e gate de conflito |

| Regra | Razão |
|---|---|
| Conteúdo embutido, nunca baixado | A instalação vira atômica: sem rede, sem `REPO_URL` apontando para branch, sem download pela metade |
| O manifesto é TOML, não slice em Go | Adicionar uma skill vira editar um arquivo de dados, não recompilar uma decisão |
| O manifesto tem teste, não confiança | `TestManifestoCobreExatamenteOQueFoiEmbutido` falha se uma skill entrar no repositório sem entrar no `manifest.toml`, ou vice-versa |
| O teste de ordem não fixa a lista | Fixá-la duplicaria o manifesto e faria toda skill nova quebrar o teste sem defeito nenhum; ele verifica a invariante — não alfabética, começa em `planning` |
| Manifesto quebrado devolve erro, não pânico | É defeito de build, e quem roda merece a mensagem, não uma stack trace |
| Terminal é verificado com `IsTerminal` | `/dev/null` é character device: uma checagem ingênua abriria o formulário para ninguém e estouraria no meio da instalação |
| Sem terminal, conflito é pulado | Travar esperando resposta em CI é pior do que não substituir |
| `Height` explícita nos formulários | Terminal que não negocia tamanho encolhe o viewport e esconde opções — o usuário escolheria entre o que vê |
| O `prompt` é injetado no `installer` | É a seam: com ela, a instalação inteira é testável sem terminal |
| `--dry-run` pergunta, mas não grava | Percorrer o fluxo é o ponto; pular as perguntas simularia outra coisa |
| `[instalado]` só no harness detectado | Rótulo em toda linha vira ruído e para de destacar; a ausência de marca é a informação |
| A marca do artefato cobre só os destinos escolhidos | Saber que a skill existe num harness não selecionado não muda decisão, e alonga a linha |
| Artefatos vêm pré-marcados | Enter instala tudo, como antes da tela existir; o contrário obrigaria a marcar sete itens no caminho mais comum |
| A escolha é reordenada pelo manifesto | A ordem do fluxo é contrato — não se confia na ordem que o formulário devolve |
| `--harness` pula os dois menus | É o caminho de script, e o `--help` já promete "sem menu interativo" |
| Seção vazia não cria diretório | Zero skills selecionadas não deve deixar um `skills/` vazio para trás |
| A marca é texto, a cor é decoração | Sem cor — `NO_COLOR`, modo acessível — a informação continua legível |
| Arte e legenda somem se não couberem | Cabeçalho quebrado em terminal estreito é ruído; `bannerWidth` é calculado da arte, não fixado à mão |
| A arte pede 111 colunas | `catlog-agents` em ANSI Shadow tem 109 células — em terminal de 80 o banner simplesmente não aparece |
| Arte em dois tons | Corpo (`█`) em azul, contorno de linha dupla (`═║╔╗╚╝`) em branco — é o contorno que faz a sombra da fonte |
| Um escape por troca de cor | Um por caractere multiplicaria a saída por dez sem diferença visual: 26–40 escapes por linha, não 109 |
| `--dry-run` ainda monta o frontmatter | É o que detecta `.yml` faltando para um harness, sem esperar a instalação real |
| O Copilot não recebe command | A CLI dele invoca a skill por `/<nome>` e não lê prompt file: os cinco arquivos ficariam no disco sem ninguém para abrir |
| Destino sem command diz por quê | Pedir cinco e receber zero, em silêncio, parece instalação pela metade |
| O menu de commands some quando nenhum destino os lê | Cobrar uma decisão que não muda nada é pior do que não perguntar |
| `COPILOT_HOME`, não `COPILOT_DIR` | É a variável da própria CLI. Inventar outra criaria dois lugares para apontar o mesmo diretório, livres para discordar |

O `install.sh` continua no repositório como instalador anterior. A saída dos dois foi comparada byte a byte: 111 arquivos idênticos.

Diferenças em relação ao instalador da `main`:

| Removido | Motivo |
|---|---|
| `AGENT_NAMES` e `install_agents()` | Não há agentes nesta branch |
| Seleção de vendor (`--vendor`, `MODEL_MAIN`, `apply_model`) | Só existia para injetar modelo nos agentes primários |
| Harnesses `codex` e `pi` | Fora do escopo deste fluxo |
| `install_agentsmd()` | Existia só para Codex e Pi. O OMP descobre skills e commands nativamente, e um `~/.agents/AGENTS.md` seria injetado em toda sessão de todo projeto |

Adicionados:

- harness `omp` (Oh My Pi) → `~/.agents/skills` e `~/.agents/commands`, override por `OMP_AGENTS_DIR`
- harness `copilot` (GitHub Copilot) → `~/.copilot/skills`, override por `COPILOT_HOME`. Sem commands

Flags do binário: `--dry-run`, `--force`, `--harness <lista>`, `--version`, `--help`. O `--local` do `install.sh` é aceito com aviso de depreciação — o conteúdo vem embutido, então não há de onde escolher.

O `AGENTS.md` da raiz é documentação do repositório e não é instalado em lugar nenhum.

No modo remoto, o `REPO_URL` aponta para a branch `feat-new-flow`. **Ao integrar em `main`, voltar para `/main`** — senão o instalador continuará baixando de uma branch que pode não existir mais.

## O que não fazer

- Não trazer agentes, skills ou commands da `main` para esta branch
- Não adicionar código executável ao **conteúdo**: o Go existe para instalar o markdown, e o fluxo não depende dele
- Não registrar skill nova só no `manifest.toml` nem só no disco — o teste de manifesto cobre os dois lados
- Não remover `<output_format>` de nenhuma skill
- Não deixar o `to-spec` fazer perguntas além da confirmação das seams
- Não deixar a `implement` decidir o que o spec não decidiu, nem implementar o que ele não pediu
- Não deixar nenhuma skill executar commit, push, branch, merge, rebase ou reset sem confirmação explícita
- Não deixar a `code-review` corrigir o que aponta, nem fundir os dois eixos em uma lista só
- Não duplicar as regras de fatiamento no `to-cards` — elas vivem em `implement/references/vertical-slices.md`
- Não deixar o `to-cards` gravar antes da aprovação, nem registrar estado de andamento dentro do card
- Não deixar a `to-memory` resumir artefato antes de enviar, nem criar página de conhecimento por feature
- Não fazer o fluxo depender de serviço de memória: ausência dele é desfecho válido, registrado e não bloqueante
