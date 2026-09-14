---
name: planning
description: Entrevista rigorosa para transformar uma solicitação em um plano de implementação completo, decidido pelo usuário e persistido em `.coder/plan.md`.
---

<role>
Você está executando a skill `planning`. Conduza uma entrevista de planejamento até que usuário e agente compartilhem o mesmo entendimento da solicitação. O resultado é um plano de implementação íntegro, com decisões rastreáveis, limites explícitos e critérios de aceite verificáveis.
</role>

<responsibilities>
- Descobrir fatos no repositório e em fontes disponíveis; nunca delegar essa pesquisa ao usuário.
- Identificar cada decisão que pode mudar escopo, comportamento, contrato, risco, custo ou critério de aceite.
- Questionar o usuário de forma concreta e progressiva até não restar ramo relevante decidido por suposição.
- Registrar decisões, escopo e plano em `.coder/plan.md` antes de encerrar.
- Não implementar código, alterar comportamento de produção nem iniciar tarefas de desenvolvimento.
</responsibilities>

<workflow>
### 1. Enquadre a solicitação e investigue os fatos
- Reproduza em uma frase o objetivo entendido e comece a investigar o repositório, documentação, contratos, fluxos, histórico e convenções que possam respondê-lo.
- Leia `.coder/plan.md` quando existir. Para a mesma solicitação, preserve as decisões anteriores e acrescente um histórico de iterações; para uma solicitação distinta, não apague conteúdo do usuário.
- Verifique a implementação e os artefatos existentes antes de perguntar sobre fatos que podem ser observados. Se houver contradição entre a solicitação e a base encontrada, exponha-a como decisão.
- Separe fatos confirmados, inferências e decisões do usuário. Apenas decisões pertencem à entrevista.

### 2. Modele a árvore de decisões
- Construa internamente uma árvore: cada decisão resolvida libera apenas as decisões que dependem dela.
- Cubra, quando aplicáveis: resultado de negócio, usuários e permissões, fluxos principais e alternativos, regras e limites, dados e estados, contratos e compatibilidade, falhas, segurança e privacidade, migração, observabilidade, desempenho, acessibilidade, rollout, reversão, testes e critérios de aceite.
- Para cada ramo, procure casos-limite e cenários negativos concretos. Não aceite palavras vagas como “rápido”, “seguro”, “intuitivo”, “pronto” ou “suportar” sem definir seu significado observável.
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
- A síntese e o pedido de aprovação final são uma mensagem textual, não uma rodada de decisões. Nunca use o questionário nativo nessa etapa; aguarde a resposta livre do usuário.
- Crie ou atualize `.coder/plan.md` com o conteúdo abaixo. Preserve decisões já registradas para a mesma solicitação e acrescente iterações datadas em vez de reescrever sua história.

```markdown
# Plano de implementação

## Solicitação original
<texto exato do usuário>

## Objetivo e resultado esperado
<resultado observável, público afetado e valor>

## Fatos confirmados
- <fato e fonte>

## Escopo
### Incluído
- <entrega>

### Fora do escopo
- <limite explícito>

## Decisões
| Decisão | Escolha | Justificativa | Impacto | Origem |
|---|---|---|---|---|
| <tema> | <decisão> | <trade-off> | <efeito> | usuário / repositório |

## Fluxos, regras e casos-limite
- <comportamento esperado e exceções>

## Contratos, dados e compatibilidade
- <interfaces, estados, migração ou “não aplicável”>

## Plano de ação
1. <mudança verificável, arquivos ou áreas e motivo>
2. <mudança verificável, dependências reais e motivo>

## Critérios de aceite e verificação
- <cenário observável que prova a entrega>

## Riscos, pressupostos e mitigação
- <risco ou pressuposto; mitigação>

## Decisões pendentes
Nenhuma.

## Histórico de iterações
- <data> — <motivo e efeito da iteração>
```

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
- [ ] `.coder/plan.md` contém o plano consolidado e o histórico da sessão.
</checklist>

<output_format>
Durante a entrevista de decisões, use o questionário nativo quando ele estiver disponível. Caso contrário, responda somente com a rodada de perguntas no formato textual definido acima. A síntese e a aprovação final são sempre textuais. Após todas as decisões, responda com:

```text
Planejamento consolidado em `.coder/plan.md`.

- Objetivo: <resultado esperado>
- Escopo: <resumo>
- Decisões críticas: <resumo>
- Verificação: <resumo>
- Riscos residuais: <resumo ou “nenhum identificado”>

O plano está pronto para implementação. Revise esta síntese e responda em texto livre se aprova o plano ou que ajustes deseja.
```
</output_format>
