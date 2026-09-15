# Mapa dos artefatos

Carregue no **passo 2**, ao selecionar o que enviar.

## A identidade do documento já está no nome do arquivo

A ferramenta de ingestão usa o **nome do arquivo como identificador** do documento, e reenviar com o mesmo identificador **substitui** o anterior. Isso não precisa ser construído: a regra de nome do fluxo já produz o comportamento certo.

| Artefato | Identidade no bank | Reenvio | Casa com a regra do fluxo |
|---|---|---|---|
| `.coder/plan-<branch-safe>.md` | Um documento por branch | **Substitui** | "Um plano por branch, atualizado na mesma solicitação" |
| `.coder/impl-<branch-safe>.md` | Um documento por branch | **Substitui** | "Atualizado ao fim de cada fatia" |
| `.coder/spec-AAAAMMDD-HHMMSS.md` | Um documento por spec | **Cria** | "Um por solicitação; snapshot" |

O caminho fixo `.coder/plan.md`, que o fluxo abandonou, colidiria aqui também: todos os planos de todas as branches viveriam como **um único documento** no bank, cada envio apagando o anterior. O mesmo defeito, uma camada acima.

## Quando cada um vale o envio

| Artefato | Envie quando | Não envie quando |
|---|---|---|
| Plano | A entrevista fechou e as decisões sobreviveram à implementação | A entrevista parou no meio, com decisões pendentes |
| Spec | Existe e guiou uma implementação | Foi escrito e abandonado sem virar código |
| Andamento | A implementação fechou, com `Desvios do spec` preenchido | Ainda há fatias pendentes |

O gatilho natural é o mesmo do commit: o momento em que o trabalho vira permanente no repositório é o momento em que vale torná-lo permanente na memória. Antes disso, decisão ainda pode ser revertida — e é o `Desvios do spec` do andamento que registra quais foram.

## O que não é enviado

| Artefato | Por quê |
|---|---|
| `.coder/cards-<branch>/` | Cards são **descartáveis por construção**: recortam a execução, não registram decisão. O porquê está no spec, e o que não sobreviveu ao código está em `Desvios do spec`. Enviá-los encheria o bank de unidades de trabalho sem valor durável |

## Por que o conteúdo vai bruto

A ferramenta pede o conteúdo completo e **proíbe resumir antes**. A razão é estrutural: a síntese acontece no servidor, a partir da evidência acumulada. Um resumo enviado no lugar do documento entrega ao servidor a sua interpretação de hoje em vez dos fatos, e ele não tem como recuperar o que você descartou.

O que torna esses artefatos valiosos como memória é exatamente o que um resumo apagaria:

| Seção | Por que ela importa depois |
|---|---|
| `Decisões`, coluna `Origem` | Distingue o que o usuário escolheu do que foi inferido do repositório. Decisão do usuário é durável; inferência envelhece com o código |
| `Fora do escopo` | O que foi recusado, e por quê. É a informação que ninguém registra em outro lugar e que volta a ser discutida a cada seis meses |
| `Desvios do spec` | Onde a decisão registrada não sobreviveu ao contato com o código. É o sinal de correção da memória |
| `Decisões de teste` | Que seams foram acordadas, e com que justificativa |

## O que não é coberto

A skill `code-review` não escreve artefato: ela devolve o relatório na conversa. Então os achados do eixo Padrões **não chegam à memória** por este caminho.

Isso é uma lacuna conhecida, não um descuido. Fechá-la exige decidir antes se o relatório passa a ser persistido em `.coder/` — decisão de escopo do fluxo, que não cabe a esta skill tomar.
