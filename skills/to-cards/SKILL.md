---
name: to-cards
description: Quebra um spec em cards — fatias verticais tracer-bullet, cada uma declarando quem a bloqueia — gravados em `.coder/cards-<branch>/`. Publica só depois da sua aprovação da granularidade. Use depois de `to-spec`, quando o trabalho será distribuído entre pessoas ou atravessará muitas sessões.
---

<role>
Você está executando a skill `to-cards`. Transforma um spec em um **grafo de cards**: fatias verticais que atravessam a seam de ponta a ponta, cada uma declarando as que precisam terminar antes dela.

Não é lista ordenada. É grafo: a qualquer momento existe uma **fronteira** de cards cujos bloqueios já fecharam e que podem começar em paralelo, por pessoas diferentes.

Esta skill não decide e não implementa: ela recorta o que o spec já decidiu.
</role>

<context>
**Entrada** — `.coder/spec-AAAAMMDD-HHMMSS.md`. Havendo mais de um, pergunte qual é o alvo; não presuma o mais recente. Sem spec, pare e instrua o usuário a executar `/planning` e `/to-spec` antes.

**Saída** — `.coder/cards-<branch-safe>/NN-slug.md`, um arquivo por card, numerado em ordem de dependência. O formato está em `references/card-format.md`.

**Card é definição, não estado.** Depois de aprovado, o arquivo do card não muda. Quem registra andamento é `.coder/impl-<branch>.md`, e ele é a **única** fonte de estado — duas fontes divergem em uma semana.

**Quando usar** — quando o trabalho será distribuído entre pessoas, ou atravessará muitas sessões com `/clear` entre elas. Trabalho que uma pessoa faz numa sessão não precisa de cards: o `implement` corta as fatias internamente e executa direto.
</context>

<workflow>

### 1. Ancore no spec e no código
- Leia o spec inteiro. `Fora do escopo` vincula tanto quanto `Histórias de usuário`.
- Confirme que as seams estão acordadas em `Decisões de teste`. Card não inventa seam.
- Explore o repositório o necessário para nomear módulos e interfaces pelo vocabulário do domínio. Respeite os ADRs da área.

### 2. Procure o prefactor
- Antes de cortar, pergunte: existe uma preparação que torne as fatias seguintes triviais? Os sinais estão em `references/decomposition.md`, junto das regras de aresta do passo 3.
- Prefactor **não muda comportamento** e não entrega nada ao usuário: reorganiza o código para que a feature caiba sem espalhar ramos condicionais.
- Se existir, ele é o **card 1**, e bloqueia todos os que dependem da estrutura nova. A suíte fica verde ao final dele, como em qualquer card.
- Se não existir, não invente. Prefactor especulativo é refatoração que ninguém pediu.

### 3. Corte as fatias e declare as arestas
- Aplique as regras de corte de `vertical-slices.md` da skill `implement`: bala traçante, uma seam por fatia, dimensionada para caber em **uma janela de contexto nova**, raio de alcance, e expandir–contrair para mudança larga. Elas vivem lá, e não são duplicadas aqui, para as duas skills não divergirem.
- Para cada card, declare **Bloqueado por**: os cards que precisam terminar antes. Card sem bloqueio começa imediatamente.
- Aresta é bloqueio real — "B não é observável sem A" —, nunca preferência de ordem. `references/decomposition.md` traz a tabela do que é e do que não é aresta, mais a fronteira e o caminho mais longo, que entram na proposta.
- Caminho feliz primeiro; erro, estado vazio e caso-limite são cards próprios.

### 4. Aprove a granularidade — e espere
- Apresente a quebra numerada no formato de `<output_format>`, com título, bloqueios e o que cada card entrega.
- Pergunte explicitamente: a granularidade está certa? As arestas são bloqueios reais? Algum card deve ser fundido ou dividido?
- **Itere até o usuário aprovar. Sem aprovação, nenhum arquivo é escrito.** Card publicado é trabalho distribuído para outras pessoas; errar a granularidade depois custa retrabalho de todo mundo.

### 5. Publique
- Grave os arquivos seguindo `references/card-format.md`, numerados em ordem de dependência: bloqueadores antes dos bloqueados.
- Um card por arquivo. Nunca um arquivo único com todos.
- Se o diretório já existir para outro spec, **pare e pergunte** antes de sobrescrever.

</workflow>

<rules>
- **Fatia vertical, nunca horizontal.** Cada card atravessa todas as camadas que precisa e é verificável sozinho. "Todos os modelos" não é card.
- **Um card cabe em uma janela de contexto nova.** Esse é o tamanho: quem pegar o card não precisa da conversa que o gerou.
- **Aresta é bloqueio real**, não ordem preferida.
- **Prefactor vem primeiro ou não vem.** Preparação depois da mudança não é prefactor, é conserto.
- **Sem caminho de arquivo e sem snippet.** O card pode esperar semanas na fila; caminho envelhece antes de alguém pegá-lo. Mesma exceção do spec: trecho que codifica uma decisão com mais precisão que a prosa.
- **Card é definição, `.coder/impl-<branch>.md` é estado.** Nunca registre andamento no card.
- **Nada é escrito antes da aprovação** do passo 4.
- Esta skill não escreve código de produção, não cria branch e não commita.
</rules>

<checklist>
- [ ] O spec alvo foi identificado explicitamente.
- [ ] As seams vieram do spec; nenhuma foi inventada.
- [ ] A oportunidade de prefactor foi avaliada — e emitida como card 1, ou descartada conscientemente.
- [ ] Cada card atravessa uma seam de ponta a ponta e é verificável sozinho.
- [ ] Cada card cabe em uma janela de contexto nova.
- [ ] Toda aresta declarada é bloqueio real, não preferência de ordem.
- [ ] Mudança larga foi sequenciada como expandir–contrair, não forçada em um card.
- [ ] O usuário aprovou a granularidade antes de qualquer arquivo ser escrito.
- [ ] Um arquivo por card, numerado em ordem de dependência.
- [ ] Nenhum card contém caminho de arquivo, snippet fora da exceção, ou estado de andamento.
</checklist>

<output_format>
No passo 4, para aprovação:

```text
Quebra proposta — <N> cards a partir de `.coder/spec-AAAAMMDD-HHMMSS.md`.

1. <título> — bloqueado por: nenhum
   Entrega: <comportamento observável de ponta a ponta>
2. <título> — bloqueado por: 1
   Entrega: <comportamento observável de ponta a ponta>

Fronteira inicial: <cards que podem começar já>
<em 1 linha: por que o card 1 vem primeiro>

A granularidade está certa? As arestas são bloqueios reais? Algum card deve ser fundido ou dividido?
```

Depois da aprovação e da gravação:

```text
<N> cards em `.coder/cards-<branch>/`.

- Fronteira inicial: <cards sem bloqueio>
- Prefactor: <card 1, ou "nenhum necessário">
- Seams cobertas: <lista>
- Caminho mais longo: <N cards em sequência>

Cada card cabe em uma janela nova. Para executar: `/implement <card>`.
```
</output_format>
