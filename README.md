# coder — fluxo `planning` → `to-spec` → `implement`

Branch isolada com o fluxo completo: uma entrevista que decide (`planning`), uma síntese que registra (`to-spec`) e uma execução que constrói (`implement`), apoiada por `tdd` no ciclo de cada fatia e por `code-review` na revisão. No fim, `to-memory` leva o que foi decidido para a memória de longo prazo — quando ela existir. Nada mais.

O conteúdo é markdown puro. O instalador é um binário Go que embute `skills/` e `commands/` com `go:embed` — existe para distribuir o markdown, não para participar do fluxo.

## Índice

- [O fluxo](#o-fluxo)
- [Estrutura do repositório](#estrutura-do-repositório)
- [Skills](#skills)
- [Commands](#commands)
- [Artefatos gerados](#artefatos-gerados)
- [Cards: quando o trabalho é distribuído](#cards-quando-o-trabalho-é-distribuído)
- [Memória de longo prazo](#memória-de-longo-prazo)
- [Instalação](#instalação)
- [Diretórios de instalação](#diretórios-de-instalação)
- [Opções do instalador](#opções-do-instalador)
- [Convenções do projeto](#convenções-do-projeto)

## O fluxo

```text
/planning ─── entrevista até a árvore de decisões esvaziar ──→ .coder/plan-<branch>.md
   │
   └─→ /to-spec ─── sintetiza, confirma as seams, não entrevista ──→ .coder/spec-<timestamp>.md
          │
          ├─→ /to-cards ─── quebra em cards com arestas de bloqueio ──→ .coder/cards-<branch>/
          │                                                              (quando for distribuir)
          └─→ /implement ─── fatias verticais, uma seam por vez ──→ código + .coder/impl-<branch>.md
                 │      │
                 │      ├── skill tdd ──────── ciclo vermelho → verde de cada fatia
                 │      └── skill code-review ─ dois eixos ao final: spec e padrões
                 │
                 └─→ /to-memory ─── envia os artefatos crus ──→ memória de longo prazo (opcional)
```

O `recall` da memória fecha o ciclo na outra ponta: o `planning` consulta o histórico de longo prazo **antes** de abrir a entrevista, quando o harness o expõe. Sem isso, o desenho seria um cano de mão única, e você reabriria em uma branch a decisão já resolvida em outra.

| Etapa | Papel | Pode perguntar? |
|---|---|---|
| `planning` | **Decide.** Entrevista por rodadas até nenhuma decisão material ficar implícita | Sim — é a razão de existir |
| `to-spec` | **Registra.** Converte as decisões em um spec que sobrevive ao fim da janela de contexto | Só a confirmação das seams |
| `to-cards` | **Recorta.** Quebra o spec em cards com arestas de bloqueio, para distribuir | Sim — aprova a granularidade antes de gravar |
| `implement` | **Executa.** Constrói em fatias verticais, testando nas seams acordadas | Só diante de lacuna material — e então para |
| `tdd` | **Prova.** Conduz o ciclo vermelho → verde de cada fatia | Só para confirmar seam não acordada |
| `code-review` | **Aponta.** Revisa a diff em dois eixos independentes, sem corrigir | Só o ponto fixo, quando não informado |
| `to-memory` | **Persiste.** Envia os artefatos crus para a memória de longo prazo e mantém as perguntas permanentes | Só a confirmação antes do envio |

A separação é o ponto: quem decide não escreve o documento final, quem escreve não reabre decisão, e quem executa não decide o que ficou em aberto. Um spec que afirma algo que ninguém decidiu é defeito, não iniciativa — e o mesmo vale para código que o spec não pediu.

**Quando usar o `to-spec`:** quando o trabalho atravessa várias sessões. Se cabe em uma janela de contexto, o spec não paga o próprio custo — vá de `planning` direto para `implement`.

## Estrutura do repositório

```text
skills/
  planning/
    SKILL.md
    references/
      decision-tree.md       15 dimensões de decisão com perguntas-sonda
      plan-format.md         nome do arquivo + template do plano
  to-spec/
    SKILL.md
    references/
      plan-to-spec-map.md    qual plano ler + mapeia suas seções no spec
      seams.md               conceito, heurísticas e anti-padrões
      spec-format.md         template do .coder/spec-*.md
  to-cards/
    SKILL.md
    references/
      decomposition.md       prefactor, arestas de bloqueio, fronteira
      card-format.md         nome do diretório + template do card
  implement/
    SKILL.md
    references/
      impl-format.md         nome do arquivo + template do andamento
      vertical-slices.md     como cortar, raio de alcance, expandir–contrair
      commit-gate.md         stage explícito, revisão da diff e confirmação
  tdd/
    SKILL.md
    references/
      test-quality.md        bom teste, anti-padrões, teste tautológico
      mocking.md             fronteira de sistema, injeção, determinismo
  code-review/
    SKILL.md
    references/
      spec-axis.md           faltando, escopo indevido, implementado errado
      standards-axis.md      convenção documentada + linha de base de smells
  to-memory/
    SKILL.md
    references/
      artifact-map.md        qual artefato vira qual documento, e quando enviar
      redaction.md           varredura antes do envio, perfil de risco por artefato
      knowledge-pages.md     as perguntas permanentes, e o que nunca vira página
commands/
  to-spec/       body.md + opencode.yml + claude.yml + omp.yml
  to-cards/      idem
  implement/     idem
  code-review/   idem
  to-memory/     idem
manifest.toml    o que é instalado, na ordem do fluxo
main.go          CLI: flags, orquestração, ajuda
embed.go         go:embed do manifesto, das skills e dos commands
harness.go       destinos, overrides de diretório, --harness
assemble.go      frontmatter do harness + body.md → .md final
installer.go     cópia das skills, gate de sobrescrita
tree.go          caminhada da árvore embutida: cópia e contagem
banner.go        a arte ASCII, a largura calculada dela e o traçado
output.go        cores, detecção de terminal, banner e resumo
prompt.go        formulários huh: seleção de harness e gate de conflito
install.sh       instalador anterior, em bash (ver Instalação)
```

Esta branch não contém agentes. As sete skills rodam no agente ativo do harness.

## Skills

Formato [Agent Skills](https://agentskills.io/specification): uma pasta por skill, `SKILL.md` obrigatório, `references/` opcional carregado sob demanda.

| Skill | O que faz | Saída |
|---|---|---|
| `planning` | Entrevista rigorosa: investiga os fatos no repositório, modela a árvore de decisões e pergunta por rodadas até a fronteira esvaziar | `.coder/plan-<branch>.md` |
| `to-spec` | Sintetiza as decisões já tomadas em um spec de sete seções, confirmando antes as seams de teste | `.coder/spec-AAAAMMDD-HHMMSS.md` |
| `to-cards` | Quebra o spec em cards tracer-bullet com arestas de bloqueio, incluindo prefactor quando necessário | `.coder/cards-<branch>/NN-slug.md` |
| `implement` | Corta o spec em fatias verticais e executa uma por vez, ou consome um card, verificando o conjunto ao final e preparando o commit | código, testes e `.coder/impl-<branch>.md` |
| `tdd` | Conduz o ciclo vermelho → verde de uma fatia na seam acordada, e diz o que faz o teste valer a pena guardar | teste + implementação da fatia |
| `code-review` | Revisa a diff contra um ponto fixo em dois eixos independentes: aderência ao spec e aderência aos padrões | relatório; não altera arquivo |
| `to-memory` | Envia os artefatos de `.coder/` crus ao serviço de memória e mantém as páginas de conhecimento | documentos e páginas no bank; não altera arquivo |

### Por que `tdd` e `code-review` são skills, e não references da `implement`

Porque são úteis fora dela. Revisar uma branch que ninguém implementou com `/implement`, ou rodar um ciclo vermelho → verde em uma correção pontual, são usos legítimos que uma reference não atende — reference é material carregado por um passo de outra skill, não porta de entrada.

O custo é real: três skills novas em vez de uma, três entradas no manifesto do instalador. A contrapartida é que a `implement` fica sendo o que a do Matt é — um despachante curto — em vez de absorver duas responsabilidades que já têm nome próprio.

### Sobre a ausência de `disable-model-invocation`

A skill do Matt usa esse campo para fechar a porta do modelo: só o usuário invoca. Aqui ele não existe — `skills-ref validate` aceita apenas `allowed-tools`, `compatibility`, `description`, `license`, `metadata` e `name`, e rejeita o resto.

A proteção equivalente está dentro da skill: a `implement` para quando não encontra decisões tomadas e seams acordadas, em vez de improvisá-las.

### Progressive disclosure

Nenhum `SKILL.md` carrega o template que só será usado no último passo. Cada reference entra no passo que precisa dela:

| Skill | Passo | Carrega |
|---|---|---|
| `planning` | 2 — modelar a árvore | `references/decision-tree.md` |
| `planning` | 4 — consolidar | `references/plan-format.md` |
| `to-spec` | 1 — reunir decisões | `references/plan-to-spec-map.md` |
| `to-spec` | 3 — propor seams | `references/seams.md` |
| `to-spec` | 4 — escrever | `references/spec-format.md` |
| `to-cards` | 2 e 3 — prefactor e arestas | `references/decomposition.md` |
| `to-cards` | 5 — publicar | `references/card-format.md` |
| `implement` | 1 — ancorar e retomar | `references/impl-format.md` (regra de nome) |
| `implement` | 2 — cortar as fatias | `references/vertical-slices.md` |
| `implement` | 6 — fechar | `references/commit-gate.md` |
| `tdd` | 2 — vermelho | `references/test-quality.md`, `references/mocking.md` |
| `code-review` | 4 — rodar os eixos | `references/spec-axis.md`, `references/standards-axis.md` |
| `to-memory` | 2 — selecionar | `references/artifact-map.md` |
| `to-memory` | 3 — varrer | `references/redaction.md` |
| `to-memory` | 5 — manter as páginas | `references/knowledge-pages.md` |

No `code-review` a separação tem um efeito extra: quando o harness oferece subagentes, cada eixo roda com **apenas** a sua reference no contexto. É isso que os mantém cegos um para o outro — a condição para que um eixo não mascare o achado do outro.

No `planning` a extração não foi por tamanho — o corpo cabia folgado no limite de 500 linhas da spec. Foi para separar o contrato de saída da prosa do workflow e, principalmente, para o guia de decisões poder crescer: as 15 dimensões deixaram de ser uma lista comprimida em uma linha e passaram a ter perguntas-sonda concretas e as dependências entre ramos.

O fallback textual de perguntas continua inline de propósito: são 3% do arquivo e ele só é acionado quando a ferramenta nativa de questionário não existe — justamente quando não se quer depender de mais uma leitura de arquivo.

Validar antes de commitar:

```bash
for s in planning to-spec to-cards implement tdd code-review to-memory; do npx -y skills-ref validate ./skills/$s; done
```

## Commands

| Command | Harnesses | Efeito |
|---|---|---|
| `/to-spec [foco]` | OpenCode, Claude Code, Oh My Pi | Executa a skill `to-spec` sobre o plano da branch atual ou a conversa atual |
| `/to-cards [spec]` | idem | Quebra o spec em cards, depois de você aprovar a granularidade |
| `/implement [spec, card ou recorte]` | idem | Executa a skill `implement` sobre um spec ou **um** card |
| `/code-review <ponto-fixo>` | idem | Executa a skill `code-review` sobre `<ponto-fixo>...HEAD` |
| `/to-memory [artefato]` | idem | Envia os artefatos da branch atual à memória de longo prazo, sob confirmação |

Têm command as skills que o **usuário** dispara. `planning`, `to-spec`, `to-cards`, `implement`, `code-review` e `to-memory` são portas de entrada; `tdd` não é — ela é acionada pela `implement` dentro do ciclo de cada fatia.

`planning` ficou sem command por ser a única cuja invocação não precisa de argumento nem de recorte: invoque-a pelo nome (`/planning`, quando o harness expõe skills como slash commands) ou peça o planejamento em linguagem natural.

## Artefatos gerados

Todos ficam em `.coder/`, no repositório onde o fluxo roda.

| Arquivo | Origem | Ciclo de vida |
|---|---|---|
| `.coder/plan-<branch-safe>.md` | `planning` | **Um por branch.** Atualizado na mesma solicitação, com histórico de iterações |
| `.coder/spec-AAAAMMDD-HHMMSS.md` | `to-spec` | Um por solicitação; ajustes na mesma sessão atualizam o mesmo arquivo |
| `.coder/cards-<branch-safe>/NN-slug.md` | `to-cards` | **Um diretório por branch**, um arquivo por card. Definição imutável depois de aprovada |
| `.coder/impl-<branch-safe>.md` | `implement` | **Um por branch.** Atualizado ao fim de cada fatia; é ele que permite retomar em sessão nova |

`code-review` e `to-memory` não geram artefato: a primeira devolve o relatório na conversa, a segunda escreve fora do repositório.

### Por que plano e andamento levam a branch no nome, e o spec não

O plano é **lido de volta** pelo `/to-spec`. Um caminho fixo (`.coder/plan.md`) seria um arquivo único no disco: `.coder/` está fora do versionamento, então ele não viaja com a branch e um `/planning` em outra branch sobrescreveria o anterior. Pior que a perda: o `/to-spec` seguiria lendo o caminho e produziria um spec plausível derivado do plano errado. Timestamp resolveria a perda, não o pareamento — a branch no nome resolve os dois.

O andamento da `implement` segue a mesma regra pela mesma razão, com um agravante: ele é lido de volta pela **própria skill**, na sessão seguinte, para saber quais fatias já fecharam. Ler o andamento de outra branch faria a skill pular fatias que nunca foram implementadas aqui.

O spec não é lido de volta por nenhuma skill deste fluxo — a `implement` recebe o caminho dele, não o descobre por convenção de nome. Então o timestamp basta para nunca colidir, e a branch de origem fica registrada no cabeçalho do documento.

| Branch atual | Plano | Cards | Andamento |
|---|---|---|---|
| `feat/rate-limit` | `.coder/plan-feat-rate-limit.md` | `.coder/cards-feat-rate-limit/` | `.coder/impl-feat-rate-limit.md` |
| `main` | `.coder/plan-main.md` | `.coder/cards-main/` | `.coder/impl-main.md` |
| `HEAD` desanexado | `.coder/plan-AAAAMMDD-HHMMSS.md` | `.coder/cards-AAAAMMDD-HHMMSS/` | `.coder/impl-AAAAMMDD-HHMMSS.md` |

Os cards seguem a mesma regra pelo mesmo motivo: são lidos de volta pelo `/implement`.

Rodar `/planning` na mesma branch com uma solicitação diferente não sobrescreve nada: a skill para e pergunta. A `implement` faz o mesmo quando o arquivo da branch aponta para outro spec.

## Cards: quando o trabalho é distribuído

`to-cards` só se paga quando o trabalho **sai de uma pessoa**: vários devs, ou muitas sessões com `/clear` entre elas. Uma pessoa numa sessão não precisa de cards — o `implement` corta as fatias internamente e executa direto.

| | Sem cards | Com cards |
|---|---|---|
| Quem corta | `implement`, no passo 2 | `to-cards`, com sua aprovação |
| Onde o corte vive | Só na sessão | `.coder/cards-<branch>/`, por arquivo |
| Granularidade | Apresentada e seguida | **Aprovada antes de gravar** |
| Execução | Fatias em sequência, na mesma sessão | Um card por invocação do `/implement` |

### Grafo, não fila

Cada card declara `Bloqueado por`. Card sem bloqueio está na **fronteira** e pode começar já — e o tamanho da fronteira inicial diz quantas pessoas conseguem trabalhar no primeiro dia.

Aresta é bloqueio real ("B não é observável sem A"), nunca preferência de ordem. Aresta inventada é o defeito mais caro da skill: serializa trabalho que poderia correr em paralelo, e o grafo continua parecendo correto.

O **caminho mais longo** é o piso do prazo — nenhuma quantidade de gente o encurta. A skill o mostra na proposta, porque é o que torna a granularidade discutível de verdade.

### Prefactor é o card 01

"Deixe a mudança fácil, depois faça a mudança fácil." Quando existe uma reorganização que torna as fatias seguintes triviais, ela vira o primeiro card: não muda comportamento, não entrega nada ao usuário, deixa a suíte verde sem teste novo nem alterado, e bloqueia todos os que dependem da estrutura nova.

Sem sinal claro, não se inventa — prefactor especulativo é refatoração que ninguém pediu, e o `code-review` a apontaria como escopo indevido.

### Card é definição; o andamento é estado

O arquivo do card não muda depois de aprovado. Quem registra progresso é `.coder/impl-<branch>.md`, que ganhou uma coluna `Card`. Duas fontes de estado divergem em uma semana — por isso só existe uma.

É também por isso que cards **não vão para a memória de longo prazo**: são descartáveis por construção. O porquê está no spec; o que não sobreviveu ao código está em `Desvios do spec`.

## Memória de longo prazo

O fluxo fala com um serviço de memória compartilhado entre harnesses e projetos — o Hindsight, na configuração atual — em dois pontos: `recall` na entrada do `planning` e `/to-memory` na saída.

### Opcional por construção

Nenhum dos dois pontos é obrigatório, e nenhum deles falha quando o serviço não existe:

| Estado | `planning` | `/to-memory` |
|---|---|---|
| Plugin ausente no harness | Segue a entrevista, e registra `Memória: indisponível` no cabeçalho do plano | Encerra dizendo que não enviou nada |
| Serviço fora do ar | Idem — uma tentativa, sem repetir | Encerra com o que já foi enviado |
| Bank sem resultado | Idem. Não é erro | Envia normalmente |

Duas propriedades fazem isso funcionar. A primeira é **detecção, não tentativa**: a skill verifica a lista de ferramentas do agente antes de chamar, o mesmo que já faz para escolher entre `AskUserQuestion`, `question` e `ask_user`. A segunda é **uma tentativa, sem retentativa**: ausência, falha e resultado vazio são desfechos equivalentes para o fluxo.

E degradar não é silenciar: o cabeçalho do plano registra que a entrevista rodou sem memória, porque quem ler depois precisa saber que ela pode ter reaberto decisão já tomada em outra sessão.

### O fallback não é outro sistema de memória

Sem o serviço, a fonte equivalente são os **planos que já existem em `.coder/` na máquina**: mesmo repositório, todo harness, zero configuração. Lê-los como evidência histórica é legítimo e cabe na investigação do passo 1 — o que as skills proíbem é **derivar** plano ou spec do plano de outra branch, e sobrescrevê-lo.

O que o fluxo deliberadamente **não** faz é cair para a memória nativa do harness. Ela continua capturando por conta própria, fora das skills, e não precisa que elas a acionem. Escrever nela a partir de um `SKILL.md` repetiria um erro conhecido: embutir caminho de um harness específico em conteúdo que roda em três.

### Por que o artefato vai cru

A ferramenta de ingestão pede o conteúdo completo e **proíbe resumir antes**, porque a síntese acontece do lado do servidor. Resumo enviado no lugar do documento entrega a interpretação de hoje em vez dos fatos, e o que foi descartado não volta.

O que um resumo apagaria é justamente o que dá valor a esses artefatos como memória: a coluna `Origem` das decisões, que separa o que o usuário escolheu do que foi inferido do código; o `Fora do escopo`, que é a informação que ninguém registra em outro lugar; e o `Desvios do spec`, que diz quais decisões não sobreviveram ao contato com o código.

### Por que as páginas guardam perguntas, não respostas

Uma página de conhecimento não armazena conteúdo curado: armazena uma pergunta que o servidor re-responde a cada consolidação. Você escreve a pergunta certa e a resposta se reconstrói conforme a evidência chega.

| | Documento | Página |
|---|---|---|
| Você fornece | O artefato bruto | Uma pergunta |
| Quem escreve o conteúdo | Você | O servidor |
| Quantos existem | Muitos, um por artefato | **Poucos e estáveis** |

Por isso não se cria página por feature ou por branch — isso é documento. As quatro perguntas permanentes do fluxo estão em `skills/to-memory/references/knowledge-pages.md`.

### Os nomes dos artefatos já são identificadores de documento

A ingestão usa o nome do arquivo como ID, e reenviar substitui. A regra de nome que o fluxo já adota produz o comportamento certo sem esforço:

| Artefato | No bank | Casa com |
|---|---|---|
| `plan-<branch>.md` | Um doc por branch, substituído | "um plano por branch, atualizado" |
| `impl-<branch>.md` | Um doc por branch, substituído | "atualizado ao fim de cada fatia" |
| `spec-<timestamp>.md` | Doc novo a cada spec | "snapshot" |

O caminho fixo `.coder/plan.md`, que o fluxo abandonou, colidiria aqui também: todos os planos de todas as branches viveriam como um documento só, cada envio apagando o anterior.

### Lacuna conhecida

A `code-review` não escreve artefato, então os achados do eixo Padrões **não chegam à memória**. Fechar isso exige decidir antes se o relatório passa a ser persistido em `.coder/` — decisão de escopo do fluxo, ainda não tomada.

## Instalação

O instalador é um binário único com as skills e os commands embutidos por `go:embed`. A instalação não acessa a rede: o conteúdo já está dentro do executável, e binário e conteúdo versionam juntos.

### Via `go install`

```bash
go install github.com/paraizofelipe/coder@latest
coder
```

### A partir do repositório

```bash
git clone https://github.com/paraizofelipe/coder.git
cd coder
git switch feat-new-flow
go run .
```

### Seleção de harness

```bash
coder --harness claude            # só Claude Code
coder --harness opencode,omp      # dois harnesses
coder --harness copilot           # só as skills, no ~/.copilot
coder --harness all               # todos
```

Sem a flag, abre o menu interativo — multi-seleção, espaço marca e enter confirma. O harness já presente na máquina vem marcado com `[instalado]` em amarelo; a ausência de marca é a informação sobre o resto.

Depois do destino vêm duas telas de artefatos, uma para skills e outra para commands, com tudo pré-marcado: Enter instala o conjunto inteiro, e desmarcar é que é a ação. Cada linha mostra em quais dos destinos escolhidos o artefato **já existe**:

```text
Instalação de Skills
[harness] marca onde o artefato já está instalado
> ✓ planning - [opencode]
  ✓ to-spec - [opencode, omp]
  ✓ tdd - [omp]
  ✓ to-memory
```

A marca cobre só os harnesses selecionados. Saber que uma skill está instalada num destino que você não escolheu não muda decisão nenhuma.

A detecção é a existência do diretório base — o mesmo caminho onde a instalação vai gravar, e o mesmo que os overrides de ambiente mudam. Apontar `CLAUDE_DIR` para um sandbox muda destino **e** marca, para o menu não descrever um lugar enquanto a instalação escreve em outro.

### Simulação

```bash
coder --dry-run --harness all
```

Percorre tudo — menu, conflitos, montagem de cada command — e não escreve um byte. Cada item sai como `simulado (N arquivos)` em vez de `instalado`, e o aviso de simulação aparece no começo e no fim, porque uma instalação longa rola a tela e o aviso do topo some.

A montagem do frontmatter roda mesmo assim: é o que detecta um `.yml` faltando para algum harness, defeito que de outra forma só apareceria no dia da instalação real.

### Sem terminal

Em CI, sob pipe ou com a saída redirecionada, não há onde desenhar o formulário. O binário detecta isso e não trava esperando uma resposta que não virá:

| Situação | Comportamento |
|---|---|
| Sem terminal e sem `--harness` | Para e diz para informar `--harness` |
| Sem terminal, com `--harness`, arquivo já existe | Pula o conflito e segue, em vez de sobrescrever |
| Sem terminal, com `--harness` e `--force` | Instala tudo, sem pergunta nenhuma |

A detecção é `IsTerminal` no descritor, não a existência dele: `/dev/null` é um character device e seria aceito por uma checagem ingênua.

### Adicionar uma skill

Duas coisas, nesta ordem:

1. Criar `skills/<nome>/` com `SKILL.md` e, se houver, `references/`
2. Acrescentar `"<nome>"` à lista `skills` do `manifest.toml`

O `go:embed` traz o diretório para dentro do binário automaticamente, mas o instalador percorre o manifesto, não o filesystem — sem o passo 2 a skill é embutida e nunca instalada. `TestManifestoCobreExatamenteOQueFoiEmbutido` falha nomeando exatamente o que ficou de fora, nos dois sentidos: diretório sem entrada, e entrada sem diretório.

A posição na lista é a posição no menu e no log. Ela é a ordem do fluxo, não alfabética.

Command novo segue a mesma regra, na lista `commands`, e precisa de um `.yml` para cada harness que lê commands — `TestTodoCommandTemFrontmatterParaTodosOsHarnesses` cobre isso.

### O `install.sh` continua aqui

É o instalador anterior, em bash, que baixa o conteúdo do GitHub arquivo a arquivo. Produz exatamente os mesmos 111 arquivos que o binário — a saída dos dois foi comparada byte a byte. Enquanto não houver release publicado, ele é o caminho para instalar em máquina sem Go.

## Diretórios de instalação

| Harness | Skills | Commands |
|---|---|---|
| OpenCode | `~/.config/opencode/skills/` | `~/.config/opencode/commands/` |
| Claude Code | `~/.claude/skills/` | `~/.claude/commands/` |
| Oh My Pi (`omp`) | `~/.agents/skills/` | `~/.agents/commands/` |
| GitHub Copilot | `~/.copilot/skills/` | — |

O destino do OMP é `~/.agents` porque é o que o provider **Agent Dirs** dele varre: `.agent/` e `.agents/`, tanto no walk-up do projeto quanto no home do usuário. É também o diretório canônico do padrão Agent Skills.

### O Copilot recebe skills e mais nada

A CLI do Copilot lê skills pessoais em `~/.copilot/skills` (override pela variável dela própria, `COPILOT_HOME`) e invoca cada uma por `/<nome>` — mesma porta do Claude Code. Command em markdown ela não lê: prompt file (`*.prompt.md`) continua sendo coisa de IDE, e as requisições para trazê-lo à linha de comando seguem abertas no repositório da CLI.

Por isso o alvo `copilot` instala as sete skills e nenhum dos cinco commands. Instalá-los ali deixaria cinco arquivos que nenhuma ferramenta abre — e que ninguém iria remover depois. O instalador diz isso em uma linha (`copilot não lê commands: a skill já responde a /<nome>`) em vez de calar: quem pediu cinco commands e viu zero concluiria que a instalação falhou pela metade.

Há uma sobreposição que vale conhecer: a CLI do Copilot também varre `~/.agents/skills`, que é exatamente o destino do alvo `omp`. Quem já instala no `omp` recebe as skills no Copilot sem escolher o alvo `copilot`. O alvo existe para quem não usa o OMP, e para quem prefere o diretório nativo.

### O OMP lê os outros dois, mas só no nível de projeto

O `omp` é um agregador de compatibilidade: sabe ler os diretórios do Claude Code e do OpenCode além dos próprios. Mas os defaults distinguem **projeto** de **usuário**, e é aí que mora a pegadinha.

| Setting do OMP | Default | Lê de |
|---|---|---|
| `skills.enableAgentsUser` | **ligado** | `~/.agents/skills` |
| commands do provider Agent Dirs | **sem flag — sempre ligado** | `~/.agents/commands` |
| `skills.enableClaudeProject` | ligado | `.claude/skills` (do projeto) |
| `commands.enableClaudeProject` | ligado | `.claude/commands` (do projeto) |
| `commands.enableOpencodeProject` | ligado | `.opencode/commands` (do projeto) |
| `skills.enableClaudeUser` | **desligado** | `~/.claude/skills` |
| `commands.enableClaudeUser` | **desligado** | `~/.claude/commands` |
| `commands.enableOpencodeUser` | **desligado** | `~/.config/opencode/commands` |

Consequência prática: **instalar no alvo `omp` é necessário, não opcional.** Instalar só no Claude Code ou só no OpenCode não faz o OMP enxergar nada no nível de usuário — os três settings de nível de usuário desses harnesses vêm desligados. O único caminho ligado por padrão é `~/.agents`, que é exatamente o destino do alvo `omp`.

Instalar em mais de um alvo não duplica: a deduplicação é por nome, e o provider de maior prioridade vence (Claude Code 80, Agent Dirs 70).

### Como a skill é invocada em cada harness

| Harness | Skill | Command |
|---|---|---|
| Claude Code | `/<nome-da-skill>` | `/<nome>` (redundante aqui) |
| Oh My Pi | `/skill:<nome>` (setting `skills.enableSkillCommands`, ligado) | `/<nome>`, com `argument-hint` no autocomplete |
| OpenCode | não invocável por slash — carregada pelo tool `skill`, que o **modelo** escolhe chamar | `/<nome>` |
| GitHub Copilot | `/<nome-da-skill>` (`/skills reload` recarrega sem reiniciar a sessão) | não tem |

É por isso que `to-spec`, `implement` e `code-review` têm command: no OpenCode ele é a única porta para o usuário disparar o fluxo, e no OMP ele troca `/skill:implement` por um nome curto com dica de argumento.

A `tdd` não precisa disso: quem a aciona é a `implement`, e o caminho que o modelo usa para carregá-la — o tool `skill` no OpenCode, a skill pelo nome nos outros — funciona sem command em todos eles.

Overrides por variável de ambiente: `OPENCODE_DIR`, `CLAUDE_DIR`, `OMP_AGENTS_DIR` e `COPILOT_HOME` — esta última é a variável da própria CLI do Copilot, não uma inventada aqui.

## Opções do instalador

| Flag | Efeito |
|---|---|
| `--dry-run`, `-n` | Percorre o fluxo inteiro sem gravar nada |
| `--force`, `-f` | Substitui tudo sem perguntar |
| `--harness <lista>` | `opencode`, `claude`, `omp`, `copilot`, `all` — vírgula ou espaço. Pula **os dois** menus e instala o manifesto inteiro |
| `--version`, `-v` | Versão do binário |
| `--help`, `-h` | Ajuda |

`--local` deixou de existir: o conteúdo vem embutido, então não há de onde escolher. A flag é aceita com um aviso, para não quebrar o hábito de quem vinha do script.

Em conflito, o menu oferece três saídas: pular, substituir e substituir todos os próximos. "Pular" é a opção destacada ao abrir — Enter preserva o arquivo existente.

| Variável | Efeito |
|---|---|
| `OPENCODE_DIR`, `CLAUDE_DIR`, `OMP_AGENTS_DIR`, `COPILOT_HOME` | Override do diretório base de cada harness |
| `NO_COLOR` | Desliga as cores da saída |
| `ACCESSIBLE` | Troca o TUI por prompts de texto, para leitor de tela |

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

O `install.sh` desta branch aponta o `REPO_URL` para `feat-new-flow`. Ao integrar em `main`, voltar para `/main`. O binário não tem esse problema: ele não baixa nada.
