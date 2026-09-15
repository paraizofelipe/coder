# Mapeamento plano → spec

Carregue no **passo 1**. Define qual arquivo ler, de onde vem cada seção do spec e, principalmente, **o que não migra**.

## Qual arquivo ler

`.coder/plan-<branch-safe>.md`, onde `<branch-safe>` é a branch atual (`git rev-parse --abbrev-ref HEAD`) com `/` trocado por `-`. Em `HEAD` desanexado, o `planning` grava `.coder/plan-AAAAMMDD-HHMMSS.md`.

| Situação | O que fazer |
|---|---|
| Existe plano da branch atual | Use-o como fonte canônica |
| Não existe | Sintetize do contexto da conversa |
| Existem planos de outras branches | **Ignore-os.** Descrevem outro trabalho; um spec derivado deles sai plausível e errado |
| Existem vários planos com timestamp (`HEAD` desanexado) | Pergunte ao usuário qual é o alvo. Não presuma o mais recente |

Registre no cabeçalho do spec qual arquivo de plano foi usado.

## Tabela de origem

| Seção do plano | Seção do spec | Transformação |
|---|---|---|
| `Solicitação original` | `Problema` | Reescreva na perspectiva de quem sofre o problema. Não cole o texto cru do pedido |
| `Objetivo e resultado esperado` | `Problema` + `Solução` | O resultado esperado vira a `Solução`; o público afetado ancora o `Problema` |
| `Fatos confirmados` | — | **Não migra.** Serve para ancorar o texto e evitar afirmação inventada |
| `Escopo > Incluído` | `Histórias de usuário` | Cada entrega vira uma ou mais histórias. Uma entrega que não gera história é sinal de escopo mal definido |
| `Escopo > Fora do escopo` | `Fora do escopo` | Migração direta, **preservando o motivo** de cada exclusão |
| `Decisões` | `Decisões de implementação` | Escolha + justificativa. Descarte a coluna `Origem`; agrupe por tema |
| `Fluxos, regras e casos-limite` | `Histórias de usuário` + `Decisões de implementação` | Fluxo alternativo e caso-limite viram história; regra invariante vira decisão |
| `Contratos, dados e compatibilidade` | `Decisões de implementação` | Contratos, estados, migração e compatibilidade |
| `Plano de ação` | — | **Não migra.** É sequência de execução; pertence ao fatiamento, não ao spec |
| `Critérios de aceite e verificação` | `Decisões de teste` + `Histórias de usuário` | O critério observável vira decisão de teste; o comportamento que ele prova vira história |
| `Riscos, pressupostos e mitigação` | `Notas adicionais` | Mantenha a mitigação junto do risco |
| `Decisões pendentes` | `Notas adicionais` | Se a pendência for **material**, interrompa: o spec não fica pronto |
| `Histórico de iterações` | — | **Não migra.** O spec tem histórico próprio, a partir da sua criação |

## As três coisas que nunca migram

| Item | Por quê |
|---|---|
| `Fatos confirmados` | O spec afirma decisões, não inventário de repositório. O fato entra diluído na prosa da decisão que ele justifica |
| `Plano de ação` | Ordem de execução envelhece na primeira surpresa da implementação. O spec descreve o destino |
| `Histórico de iterações` do plano | Confundir os dois históricos apaga a rastreabilidade de qual documento mudou quando |

## Quando não existe plano para a branch

Sintetize do contexto da conversa, aplicando o mesmo filtro:

- Só entra no spec o que foi **efetivamente decidido** — uma opção que o usuário escolheu, um limite que ele impôs, um fato que o repositório comprova.
- Alternativa discutida e não escolhida vira linha em `Fora do escopo` ou em `Notas adicionais`, nunca uma decisão.
- Sugestão sua que o usuário não respondeu **não é decisão**. Vira pendência.
- Se a conversa não contém decisões — só exploração ou a solicitação crua — pare e peça `/planning`.

## Teste de rastreabilidade

Antes de gravar, para cada afirmação do spec pergunte: *de onde isso veio?* Toda resposta deve ser uma destas três:

1. Decisão registrada no plano da branch atual ou tomada explicitamente na conversa
2. Fato observável no repositório (código, config, ADR, glossário)
3. Inferência **declarada como inferência** no próprio texto

Sem uma dessas três, a linha sai do spec.
