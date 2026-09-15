# Eixo Spec

Carregue no **passo 4**. Quando rodar em subagente, este arquivo vai inteiro no contexto dele — é o único acesso que ele tem a estas instruções.

## A pergunta

O código faz o que o spec pediu? **Nem menos, nem mais.**

Este eixo não julga qualidade, estilo ou arquitetura. Código feio que implementa exatamente o pedido passa aqui. Código elegante que implementa outra coisa falha.

## Insumos

- A diff (`git diff <ponto-fixo>...HEAD`) e a lista de commits
- O spec: `Histórias de usuário`, `Decisões de implementação`, `Decisões de teste` e, principalmente, `Fora do escopo`
- O card, quando o trabalho veio de um: ele recorta o que **desta** entrega estava em jogo. Com card, o escopo a cobrar é o do card, não o do spec inteiro — requisito do spec que pertence a outro card **não** é achado de "faltando" aqui

## As três categorias

| Categoria | O que é | Como reconhecer |
|---|---|---|
| **Faltando** | O spec pediu e a diff não entrega, ou entrega parcialmente | Percorra as histórias uma a uma e procure na diff o comportamento de cada uma. História sem rastro na diff é achado |
| **Escopo indevido** | A diff entrega o que o spec não pediu | Comportamento novo sem história correspondente; item listado em `Fora do escopo` que apareceu assim mesmo; refatoração oportunista em área não tocada pelo trabalho |
| **Implementado errado** | O comportamento existe, mas não é o decidido | Caso-limite tratado ao contrário do que a história descreve; decisão de implementação contrariada; contrato divergente do acordado |

Escopo indevido pesa igual a requisito faltando. É trabalho que ninguém pediu, que ninguém revisou como decisão e que a equipe vai manter para sempre.

## Como percorrer

1. Enumere as histórias do spec. Para cada uma, procure na diff o comportamento observável que ela descreve — não o arquivo que "parece" implementá-la.
2. Confira `Decisões de implementação` item a item contra a diff. Decisão contrariada é achado mesmo quando o resultado ficou melhor: o lugar de mudar a decisão é o spec.
3. Leia `Fora do escopo` e verifique se algum item recusado entrou.
4. Confira `Decisões de teste`: os testes da diff vivem nas seams acordadas? Teste em seam não acordada é achado deste eixo, não do outro.
5. Percorra a diff em sentido inverso: para cada comportamento novo, aponte a história que o justifica. Sem história, é escopo indevido.

## Regras

- **Cite a linha do spec** em todo achado. Sem citação, é opinião sobre o que o spec deveria ter dito.
- Divergência que **o spec não cobre** não é achado deste eixo — é lacuna do spec. Registre como tal, em uma linha, sem transformá-la em defeito de quem implementou.
- Não proponha implementação alternativa. Este eixo descreve a distância entre o pedido e o entregue.
- Sem spec disponível, não execute o eixo. Reporte `nenhum spec disponível; eixo não executado`.

## Formato do achado

```text
1. [faltando | escopo indevido | implementado errado] <o que se observa>
   Spec: "<linha citada do spec>"
   Diff: <arquivo> — <trecho ou descrição do hunk>
```

Abaixo de 400 palavras no total. Ordene por gravidade: comportamento decidido que não existe vem antes de detalhe divergente.
