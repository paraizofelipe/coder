# Decomposição: prefactor, corte e arestas

Carregue nos **passos 2 e 3**, antes de montar a proposta para aprovação.

## As regras de corte vivem em outro lugar

Bala traçante, uma seam por fatia, raio de alcance e expandir–contrair estão em `vertical-slices.md` da skill `implement`, e é de lá que vêm — duplicá-las aqui garantiria divergência entre as duas skills em poucas semanas.

O que este arquivo acrescenta é o que só existe quando o trabalho vira grafo distribuído: **prefactor** e **arestas**.

## Prefactor

> Deixe a mudança fácil, depois faça a mudança fácil.

Antes de cortar, olhe o código que vai receber o trabalho e pergunte: existe uma reorganização que torne as fatias seguintes triviais?

| Sinal de que há prefactor | Exemplo |
|---|---|
| A mesma condicional se repete em vários pontos e a mudança adicionaria mais um ramo em cada | `if cartao / else boleto` em quatro lugares, e agora entra Pix |
| A feature precisa de um dado que hoje não atravessa a fronteira | O handler não recebe o contexto que a regra nova exige |
| Existe duplicação que a mudança obrigaria a alterar em dois lugares iguais | Duas cópias da validação, e a regra muda nas duas |

Características de um card de prefactor:

- **Não muda comportamento.** A suíte existente fica verde, sem teste novo e sem teste alterado
- **Não entrega nada ao usuário.** A entrega é a estrutura
- **É o card `01`**, e bloqueia todo card que depende da estrutura nova

Sem sinal, **não invente**. Prefactor especulativo é refatoração que ninguém pediu, e a `implement` reprovaria isso como escopo indevido.

Preparação feita depois da mudança não é prefactor: é conserto, e já custou o retrabalho que o prefactor existia para evitar.

## Arestas de bloqueio

Uma aresta declara: **este card não é observável antes daquele terminar**.

| É aresta | Não é aresta |
|---|---|
| B usa a interface que A cria | "Faz mais sentido fazer A primeiro" |
| B lê o campo que a migração de A adiciona | A e B tocam o mesmo arquivo |
| B só é verificável com o endpoint que A expõe | A é mais importante que B |
| Todo card de migração depende do card de expandir | Preferência de quem vai pegar |

Aresta inventada é o defeito mais caro desta skill: ela **serializa trabalho que poderia correr em paralelo**, e ninguém percebe — o grafo parece correto, só está mais lento do que precisava.

Em caso de dúvida, pergunte: *se as duas pessoas começassem agora, a segunda ficaria bloqueada de verdade?* Se a resposta é "não, só ficaria um pouco mais difícil", não é aresta.

### Fronteira

A **fronteira** é o conjunto de cards cujos bloqueios já fecharam. No início, são os cards sem bloqueio nenhum — e é o número deles que diz quantas pessoas podem começar no primeiro dia.

Fronteira inicial de um card só, num grafo de dez, é sinal de que o corte virou uma fila. Reveja as arestas antes de apresentar.

### Caminho mais longo

A maior sequência de cards encadeados é o piso do prazo: nenhuma quantidade de pessoas a encurta. Vale medir e mostrar na proposta — é a informação que torna a granularidade discutível de verdade.

## Mudança larga

Quando o raio de alcance abre muitos chamadores, o expandir–contrair de `vertical-slices.md` vira três grupos de cards:

| Card | Bloqueado por |
|---|---|
| Expandir — adiciona o novo ao lado do antigo | Nada |
| Migrar (um card por lote, dimensionado pelo raio) | O card de expandir |
| Contrair — remove o antigo | **Todos** os cards de migrar |

Cada lote mantém a suíte verde sozinho, porque o antigo ainda existe. Quando nem isso se sustenta, mantenha a sequência e faça os lotes convergirem num card final de integrar-e-verificar, onde o verde é prometido.
