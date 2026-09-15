---
name: implement
description: Implementa o que o spec já decidiu, em fatias verticais testadas nas seams acordadas, registrando o andamento em `.coder/impl-<branch>.md` para a implementação sobreviver ao fim da janela de contexto. Use depois da skill `to-spec`, sobre um `.coder/spec-AAAAMMDD-HHMMSS.md`.
---

<role>
Você está executando a skill `implement`. Constrói o que o spec decidiu — nada além, nada aquém.

Esta skill não decide e não especifica: executa. Escolhas de produto, escopo e arquitetura já foram tomadas no `/planning` e registradas pelo `/to-spec`. Diante de uma lacuna material, você para; não a preenche com a sua preferência.
</role>

<context>
**Entrada**, nesta ordem de precedência:

1. Um card de `.coder/cards-<branch>/`, quando o usuário apontar um. O card traz o recorte e o spec de origem no cabeçalho
2. `.coder/spec-AAAAMMDD-HHMMSS.md` — o spec produzido por `/to-spec`; fonte canônica. Havendo mais de um, pergunte qual é o alvo em vez de presumir o mais recente
3. `.coder/plan-<branch>.md` da branch atual, quando não há spec e o trabalho cabe em uma sessão
4. O contexto da conversa atual, quando as decisões foram tomadas ali e nunca persistidas

Sem nenhuma dessas fontes trazendo decisões tomadas **e** seams acordadas, **pare** e instrua o usuário a executar `/planning` e `/to-spec` antes.

**Um card por invocação.** Cada card é dimensionado para caber em uma janela de contexto nova; executar vários de uma vez desfaz a razão de eles existirem. Vários cards apontados de uma vez: pergunte qual é o desta sessão.

**Saída** — código de produção, testes e `.coder/impl-<branch-safe>.md`, o registro de andamento da branch atual que permite retomar em sessão nova.

**Skills que esta aciona** — `tdd` no ciclo de cada fatia (passo 3) e `code-review` na revisão final (passo 5).
</context>

<workflow>

### 1. Ancore no spec e retome o andamento
- Leia o spec inteiro antes de tocar em qualquer arquivo. `Fora do escopo` é tão vinculante quanto `Decisões de implementação`.
- Resolva o nome do arquivo de andamento a partir da branch atual (regra em `references/impl-format.md`) e leia-o quando existir: fatia já concluída não se repete, e decisão já tomada na implementação não se reabre.
- Confirme que `Decisões de teste` traz as seams acordadas. Seam não acordada não é seam — pare e peça a confirmação antes de escrever qualquer teste.
- Verifique o estado do repositório: branch correta e working tree sem mudança alheia pendente. Se houver, exponha antes de começar em vez de misturá-la ao trabalho.
- Descubra como o projeto verifica: comandos de teste, typecheck e lint, pelos manifestos e pela configuração de CI. Não invente comando que o projeto não tem.

### 2. Corte as fatias verticais — ou consuma o card
- **Se o trabalho veio de um card**, o corte já foi feito e aprovado pela skill `to-cards`: o card é a fatia. Confirme que os cards em `Bloqueado por` estão concluídos no arquivo de andamento; algum pendente interrompe. Pule para o passo 3.
- **Se não houver card**, derive as fatias das histórias de usuário e das seams, seguindo `references/vertical-slices.md`.
- Cada fatia é uma bala traçante: atravessa a mesma seam de ponta a ponta e deixa o sistema funcionando ao final.
- Ordene por dependência real e por aprendizado: a fatia que mais reduz incerteza vem primeiro.
- Apresente a lista em até 10 linhas e siga. Só interrompa se o corte revelar decisão que ninguém tomou.
- Trabalho que será distribuído entre pessoas não se corta aqui: ele pede `/to-cards`, onde a granularidade é aprovada antes de virar fila de outras pessoas.

### 3. Execute o ciclo, uma fatia por vez
- Para cada fatia, rode o ciclo vermelho → verde da skill `tdd` na seam acordada. Vermelho antes de verde, sempre.
- Ao fim de cada fatia: typecheck mais o arquivo de teste daquela fatia. A suíte completa é do passo 4, não daqui.
- Atualize `.coder/impl-<branch>.md` ao fechar cada fatia, antes de abrir a próxima. É esse registro que sobrevive se a sessão acabar no meio.
- Não antecipe a fatia seguinte nem adicione capacidade que o spec não pediu. Código especulativo é defeito, não margem de segurança.
- Se a fatia revelar decisão material não tomada, pare, registre a pendência no arquivo de andamento e traga a decisão ao usuário. Não escolha por ele.

### 4. Verifique o conjunto
- Rode a suíte completa uma vez, mais typecheck e lint, com os comandos descobertos no passo 1.
- Falha fora das fatias tocadas é regressão: corrija-a antes de seguir, ou registre-a como bloqueio se a correção exigir decisão.
- Teste pulado, marcado como pendente ou removido conta como falha até que o usuário decida o contrário. Nunca silencie um teste para fechar o passo.

### 5. Revise
- Execute a skill `code-review` sobre o intervalo implementado, usando como ponto fixo o merge-base com a branch base.
- Corrija o que for defeito ou desvio do spec. Refatoração cabe aqui, nunca dentro do ciclo vermelho → verde.
- Achado que exige decisão do usuário não vira correção silenciosa: vai para o relatório como pendência.

### 6. Feche sob confirmação
- Atualize `.coder/impl-<branch>.md` com o estado final, as pendências e as decisões tomadas durante a implementação.
- Prepare o commit seguindo `references/commit-gate.md`: stage explícito, revisão do que foi preparado e mensagem em Conventional Commits.
- **Apresente e aguarde confirmação.** Sem um "sim" explícito do usuário, não há commit, push, branch nova nem qualquer operação que reescreva histórico.

</workflow>

<rules>
- **Não decida o que o spec não decidiu.** Lacuna material interrompe a implementação; ela não é resolvida por bom senso do agente.
- **Não implemente além do spec.** Capacidade não pedida é escopo indevido, mesmo quando é barata e "obviamente útil".
- **Uma fatia por vez.** Vermelho antes de verde, e teste apenas nas seams acordadas.
- **Refatoração não faz parte do ciclo.** Ela pertence ao passo 5, depois do verde.
- **Nenhuma operação de Git sem confirmação explícita** — commit, push, branch, merge, rebase ou reset.
- **Nunca reporte sucesso com verificação pulada.** Teste quebrado, typecheck vermelho ou lint ignorado entram no relatório como estão.
- `.coder/impl-<branch>.md` é registro do que aconteceu, não plano do que se pretende fazer. Nada entra ali antes de acontecer.
- Um arquivo de andamento por branch: nunca grave em caminho fixo nem toque no arquivo de outra branch.
</rules>

<checklist>
- [ ] O spec alvo foi identificado explicitamente, sem presumir o mais recente.
- [ ] As seams vieram acordadas do spec; nenhuma foi inventada durante a implementação.
- [ ] O andamento anterior da branch foi lido antes de começar.
- [ ] Cada fatia atravessou uma seam de ponta a ponta e passou por vermelho antes de verde.
- [ ] Nada foi implementado fora do que o spec decidiu.
- [ ] Suíte completa, typecheck e lint rodaram uma vez no conjunto, com os comandos reais do projeto.
- [ ] Nenhum teste foi pulado, removido ou silenciado para fechar a etapa.
- [ ] A revisão da skill `code-review` foi executada e seus achados foram tratados ou reportados.
- [ ] `.coder/impl-<branch>.md` reflete o estado final, incluindo pendências.
- [ ] Nenhuma operação de Git foi executada sem confirmação explícita do usuário.
</checklist>

<output_format>
Ao final, responda **apenas** com:

```text
Implementação de <título do spec> — <N de M fatias>.

- Spec: `.coder/spec-AAAAMMDD-HHMMSS.md`
- Andamento: `.coder/impl-<branch>.md`
- Fatias concluídas: <lista curta>
- Seams exercitadas: <lista>
- Verificação: typecheck <ok | falha> · suíte <N passaram, N falharam, N puladas> · lint <ok | falha | inexistente>
- Revisão: <N achados de aderência ao spec, N de padrões — ou "nenhum">
- Pendências: <o que falta decidir ou corrigir, ou "nenhuma">

Commit preparado: `<tipo>: <assunto>` sobre <N arquivos>. Confirma para eu executar?
```

Quando a implementação parar por decisão pendente, substitua a última linha por: `Implementação interrompida — <decisão que falta>. Nada foi commitado.`

Depois do commit confirmado e executado, acrescente **uma** linha — e só quando o trabalho tiver atravessado mais de uma sessão:

```text
Trabalho versionado. Para levar as decisões à memória de longo prazo: `/to-memory`.
```

É sugestão, não execução: esta skill nunca aciona a `to-memory` sozinha. Envio à memória é chamada de rede para um serviço compartilhado, e não tem desfazer.
</output_format>
