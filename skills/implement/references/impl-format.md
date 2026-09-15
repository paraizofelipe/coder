# Formato do `.coder/impl-<branch>.md`

Carregue no **passo 1**, para resolver o nome do arquivo e retomar o andamento, e ao fim de cada fatia no **passo 3**, ao registrar.

## Nome do arquivo

`.coder/impl-<branch-safe>.md`, onde `<branch-safe>` é a branch atual (`git rev-parse --abbrev-ref HEAD`) com `/` trocado por `-` — a mesma regra do plano.

| Branch atual | Arquivo |
|---|---|
| `feat/rate-limit` | `.coder/impl-feat-rate-limit.md` |
| `fix/timeout-cotacao` | `.coder/impl-fix-timeout-cotacao.md` |
| `main` | `.coder/impl-main.md` |
| `HEAD` desanexado, ou fora de repositório Git | `.coder/impl-AAAAMMDD-HHMMSS.md` |

A implementação acontece em uma branch; o andamento dela também. Nunca grave no caminho fixo `.coder/impl.md` nem toque no arquivo de outra branch.

### Quando o arquivo da branch já existe

| Situação | O que fazer |
|---|---|
| Mesmo spec | **Retome.** Fatias com `concluída` não se repetem; decisões registradas não se reabrem. Continue da primeira fatia pendente |
| Spec diferente | **Pare e pergunte** ao usuário: retomar mesmo assim, arquivar o anterior ou sobrescrever. Dois trabalhos sem relação não convivem no mesmo arquivo |
| Fatias pendentes com o working tree sujo | Exponha o estado antes de continuar. Pode ser trabalho interrompido no meio de uma fatia |

## Template

````markdown
# Andamento da implementação

> Branch: `<branch atual>` · Spec: `.coder/spec-AAAAMMDD-HHMMSS.md` · Iniciado <AAAA-MM-DD HH:MM> · Atualizado <AAAA-MM-DD HH:MM>

## Seams acordadas
| Seam | Origem |
|---|---|
| <nome> | `Decisões de teste` do spec |

## Fatias
| # | Fatia | Card | Seam | Estado | Verificação |
|---|---|---|---|---|---|
| 1 | <comportamento observável que a fatia entrega> | `01-slug` | <seam> | concluída | typecheck ok · `<arquivo de teste>` ok |
| 2 | <comportamento observável que a fatia entrega> | — | <seam> | pendente | — |

## Decisões tomadas durante a implementação
| Decisão | Escolha | Por quê | Origem |
|---|---|---|---|
| <tema> | <o que foi feito> | <razão> | spec / convenção do repositório / usuário |

## Desvios do spec
- <o que o spec dizia, o que foi feito e por quê> — ou `Nenhum.`

## Pendências e bloqueios
- <decisão que falta, ou falha não resolvida> — ou `Nenhuma.`

## Verificação do conjunto
- Suíte completa: <resultado e data> — ou `ainda não executada`
- Typecheck: <resultado> · Lint: <resultado ou `inexistente no projeto`>

## Histórico de sessões
- <AAAA-MM-DD HH:MM> — <o que avançou nesta sessão>
````

## Regras por seção

| Seção | Cuidado |
|---|---|
| Cabeçalho | Nomeia o spec de origem. É ele que distingue "retomar" de "outro trabalho na mesma branch" |
| `Seams acordadas` | Copiadas do spec, não escolhidas aqui. Seam que aparecer nesta tabela sem estar no spec é defeito |
| `Fatias` | Estados possíveis: `pendente`, `em andamento`, `concluída`, `bloqueada`. Uma fatia só vira `concluída` depois de verde e verificada |
| `Fatias`, coluna `Card` | O card de origem, quando houver; `—` quando a fatia foi cortada aqui. **Este arquivo é a única fonte de estado** — o card é definição e nunca registra andamento |
| `Decisões tomadas durante a implementação` | Só decisões **de execução** — nome de variável, ordem de arquivo, uso de utilitário existente. Decisão de produto ou arquitetura aqui significa que a implementação passou por cima do spec |
| `Desvios do spec` | Todo desvio entra, inclusive o justificado. `Nenhum.` é resposta válida; silêncio não |
| `Pendências e bloqueios` | Enquanto houver conteúdo aqui, a implementação não está pronta — e o relatório precisa dizer isso |
| `Verificação do conjunto` | Resultado real dos comandos, com números. "Passou" sem contagem não é registro |
| `Histórico de sessões` | Uma linha por sessão. É o que permite a sessão seguinte saber onde o trabalho parou |

## O que nunca vai para o arquivo

| Item | Por quê |
|---|---|
| Fatia marcada como concluída antes de verde | Transforma o registro em intenção; a sessão seguinte confia e segue errado |
| Decisão de produto ou de escopo | O lugar dela é o plano e o spec. Se apareceu aqui, a implementação decidiu o que não devia |
| Diff, snippet ou conteúdo de arquivo | O código está no repositório; duplicá-lo aqui só cria uma segunda versão para envelhecer |
| Resultado de teste que não foi executado | Registro falso de verificação é pior do que registro nenhum |
