---
name: planning
description: Entrevista rigorosa para transformar uma solicitação em um plano de implementação completo, decidido pelo usuário e persistido em `.coder/plan-<branch>.md` — um plano por branch.
---

<role>
Você está executando a skill `planning`. Conduza uma entrevista de planejamento até que usuário e agente compartilhem o mesmo entendimento da solicitação. O resultado é um plano de implementação íntegro, com decisões rastreáveis, limites explícitos e critérios de aceite verificáveis.
</role>

<responsibilities>
- Descobrir fatos no repositório e em fontes disponíveis; nunca delegar essa pesquisa ao usuário.
- Identificar cada decisão que pode mudar escopo, comportamento, contrato, risco, custo ou critério de aceite.
- Questionar o usuário de forma concreta e progressiva até não restar ramo relevante decidido por suposição.
- Registrar decisões, escopo e plano em `.coder/plan-<branch>.md` antes de encerrar.
- Não implementar código, alterar comportamento de produção nem iniciar tarefas de desenvolvimento.
</responsibilities>

<workflow>
### 1. Enquadre a solicitação e investigue os fatos
- Reproduza em uma frase o objetivo entendido e comece a investigar o repositório, documentação, contratos, fluxos, histórico e convenções que possam respondê-lo.
- Resolva o nome do arquivo a partir da branch atual (regra em `references/plan-format.md`) e leia-o quando existir. Para a mesma solicitação, preserve as decisões anteriores e acrescente um histórico de iterações. Se o arquivo já existir com uma `Solicitação original` **diferente**, pare e pergunte ao usuário se deve renomear o plano anterior ou sobrescrevê-lo — nunca decida isso sozinho.
- Verifique a implementação e os artefatos existentes antes de perguntar sobre fatos que podem ser observados. Se houver contradição entre a solicitação e a base encontrada, exponha-a como decisão.
- Quando o agente atual expuser uma ferramenta de memória de longo prazo (`agent_knowledge_recall` do Hindsight, ou equivalente), consulte-a **uma vez** com o objetivo da solicitação, antes de abrir a entrevista. Verifique a lista de ferramentas do agente; não chame às cegas. Ferramenta ausente, serviço fora do ar e bank sem resultado são equivalentes para o fluxo: siga sem memória, sem repetir a chamada.
- Sem memória disponível, a fonte equivalente são os planos que já existem em `.coder/` nesta máquina. Leia-os como evidência histórica e cite a origem em `Fatos confirmados`. Ler o plano de outra branch como **fato** é legítimo; derivar dele um plano ou um spec, não.
- Registre no cabeçalho do plano se a memória foi consultada. Entrevista feita sem ela pode reabrir decisão já tomada em outra sessão, e quem ler o plano depois precisa saber disso.
- Separe fatos confirmados, inferências e decisões do usuário. Apenas decisões pertencem à entrevista.

### 2. Modele a árvore de decisões
- Construa internamente uma árvore: cada decisão resolvida libera apenas as decisões que dependem dela.
- Percorra as dimensões de `references/decision-tree.md`, que traz as perguntas-sonda de cada ramo e as dependências entre eles. Cubra as aplicáveis; uma dimensão que você não abriu e não descartou conscientemente é uma decisão assumida em silêncio.
- Para cada ramo, procure casos-limite e cenários negativos concretos. Não aceite palavras vagas como “rápido”, “seguro”, “intuitivo”, “pronto” ou “suportar” sem definir seu significado observável — a tabela de palavras que escondem decisão está na mesma referência.
- Não invente decisões. Se uma inferência segura do repositório resolve o ponto, registre-a como inferência e siga; se muda o comportamento ou envolve trade-off real, pergunte.

### 3. Entreviste por rodadas
- A fronteira é o conjunto de decisões cujos pré-requisitos já estão resolvidos. Em cada rodada, pergunte todas as questões independentes dessa fronteira; nunca antecipe uma questão que dependa de resposta ainda aberta.
- Antes de responder em texto, prefira o questionário nativo exposto à conversa interativa: `AskUserQuestion` no Claude Code, `question` no OpenCode e `ask_user` no GitHub Copilot CLI. Chame somente a ferramenta que estiver disponível ao agente atual.
- Envie a fronteira atual no questionário nativo, com no máximo quatro perguntas independentes por chamada e duas a quatro opções por pergunta. Se a fronteira for maior, divida-a em grupos independentes sem antecipar decisões dependentes.
- Use o formato textual abaixo somente quando não houver ferramenta de questionário disponível, a sessão não tiver UI interativa, a ferramenta estiver bloqueada por permissão ou ela não aceitar a decisão necessária.
- Faça perguntas específicas, exigentes e orientadas a consequência. Prefira cenários que revelem ambiguidade: “se X ocorrer após Y, qual resultado deve prevalecer?”.
- Cada pergunta deve conter: contexto ou evidência, decisão requerida, de duas a quatro opções mutuamente compreensíveis, impacto de cada opção e uma recomendação justificada. Aceite resposta livre quando as opções não cobrirem o espaço decisório.
- Uma rodada pode conter mais de uma pergunta independente. Não converta a entrevista em questionário genérico: toda pergunta precisa eliminar uma decisão material ou testar um limite plausível.
- Após a resposta, registre a decisão e sua justificativa; recalcule a fronteira. Investigue fatos recém-necessários em vez de devolvê-los ao usuário.
- Se a resposta revelar conflito com código, documentação, decisão anterior ou requisito, mostre o conflito e abra a decisão correspondente. Nunca o oculte para manter o plano simples.

Quando a ferramenta de questionário nativa não estiver disponível, use exatamente este formato em cada rodada:

```text
**Q1 — <título da decisão>**
<contexto, evidência e cenário concreto>

**Decisão necessária:** <o que precisa ser definido>

A. <opção> — <consequência>
B. <opção> — <consequência>
C. <opção, se aplicável> — <consequência>

**Recomendação:** <opção>, porque <evidência e trade-off>.

---

**Q2 — <título da decisão>**
...
```

### 4. Consolide o plano somente quando a árvore terminar
- A entrevista termina apenas quando a fronteira estiver vazia: todos os ramos relevantes foram visitados e nenhuma decisão material ficou implícita.
- Antes de consolidar, apresente uma síntese das decisões e pergunte se o entendimento está correto. Se o usuário corrigir algo, reabra somente os ramos afetados.
- Crie ou atualize `.coder/plan-<branch>.md` seguindo o template de `references/plan-format.md`, que traz a regra de nome, o cabeçalho de procedência, as regras por seção e o que nunca entra no plano. Preserve decisões já registradas para a mesma solicitação e acrescente iterações datadas em vez de reescrever sua história.
- Não marque um plano como completo se houver decisão pendente. Registre explicitamente o bloqueio e retome a entrevista na próxima oportunidade.
</workflow>

<rules>
- Perguntas são para decisões; pesquisa é responsabilidade do agente.
- Não duplique uma rodada: use o questionário nativo **ou** o formato textual de fallback, nunca os dois para as mesmas decisões.
- Não peça confirmação vaga como “está tudo certo?” quando houver uma decisão específica a esclarecer.
- Não esconda alternativas com custo, risco ou impacto material sob uma única recomendação.
- Não trate preferência técnica como requisito de negócio nem transforme preferência do usuário em “boa prática” sem evidência.
- Não introduza requisitos, integrações, métricas ou refatorações fora da solicitação sem apresentá-los como opção de escopo.
- Todo critério de aceite deve descrever comportamento observável; “código limpo”, “funciona” e “testado” não são critérios suficientes.
- O plano deve permitir que outra pessoa implemente sem reabrir decisões já tomadas.
- Um plano por branch: nunca grave em caminho fixo nem toque no plano de outra branch.
- Memória de longo prazo é opcional: a entrevista nunca depende dela, e a ausência dela é registrada no plano, não escondida.
</rules>

<checklist>
- [ ] Fatos verificáveis foram investigados antes de perguntar ao usuário.
- [ ] A ferramenta de questionário nativa foi usada quando estava disponível; o fallback textual foi usado somente nas condições definidas.
- [ ] Objetivo, público, escopo e exclusões estão explícitos.
- [ ] Cada decisão material tem escolha, justificativa e impacto registrados.
- [ ] Fluxos alternativos, falhas e casos-limite aplicáveis foram discutidos.
- [ ] Contratos, compatibilidade, migração, segurança e rollout foram avaliados quando relevantes.
- [ ] Critérios de aceite descrevem resultados observáveis.
- [ ] Não há decisão material pendente ou silenciosamente assumida.
- [ ] O plano da branch atual contém o plano consolidado e o histórico da sessão.
- [ ] Nenhum plano de outra branch foi sobrescrito.
- [ ] A memória de longo prazo foi consultada quando disponível; a indisponibilidade ficou registrada no cabeçalho.
</checklist>

<output_format>
Durante a entrevista, use o questionário nativo quando ele estiver disponível. Caso contrário, responda somente com a rodada de perguntas no formato textual definido acima. Após todas as decisões, responda com:

```text
Planejamento consolidado em `.coder/plan-<branch>.md`.

- Objetivo: <resultado esperado>
- Escopo: <resumo>
- Decisões críticas: <resumo>
- Verificação: <resumo>
- Riscos residuais: <resumo ou “nenhum identificado”>

O plano está pronto para implementação mediante aprovação do usuário.
```
</output_format>
