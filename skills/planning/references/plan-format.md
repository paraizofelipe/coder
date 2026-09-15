# Formato do `.coder/plan-<branch>.md`

Carregue no **passo 1**, para resolver o nome do arquivo, e no **passo 4**, ao consolidar. O documento segue o template literal abaixo: mesma ordem, mesmos títulos.

## Nome do arquivo

`.coder/plan-<branch-safe>.md`, onde `<branch-safe>` é a branch atual (`git rev-parse --abbrev-ref HEAD`) com `/` trocado por `-`.

| Branch atual | Arquivo |
|---|---|
| `feat/rate-limit` | `.coder/plan-feat-rate-limit.md` |
| `fix/timeout-cotacao` | `.coder/plan-fix-timeout-cotacao.md` |
| `release/2026.09` | `.coder/plan-release-2026.09.md` |
| `main` | `.coder/plan-main.md` |
| `HEAD` desanexado, ou fora de repositório Git | `.coder/plan-AAAAMMDD-HHMMSS.md` |

Um plano por branch. Nunca grave no caminho fixo `.coder/plan.md` nem toque no arquivo de outra branch — é ele que o `/to-spec` vai ler depois, e um plano trocado produz um spec plausível e errado.

### Quando o arquivo da branch já existe

| Situação | O que fazer |
|---|---|
| Mesma solicitação | Atualize o mesmo arquivo e acrescente uma linha em `## Histórico de iterações` |
| `Solicitação original` diferente | **Pare e pergunte** ao usuário: renomear o plano anterior ou sobrescrever. Não decida sozinho, e não funda dois planos sem relação no mesmo arquivo |

## Template

````markdown
# Plano de implementação

> Branch: `<branch atual>` · Memória: <consultada (N resultados) | indisponível> · Criado <AAAA-MM-DD HH:MM> · Atualizado <AAAA-MM-DD HH:MM>

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
````

## Regras por seção

| Seção | Cuidado |
|---|---|
| Cabeçalho `> Branch: …` | Registra a procedência do documento. `Criado` nunca muda; `Atualizado` acompanha a última iteração. `Memória` diz se a entrevista pôde consultar o histórico de longo prazo — `indisponível` avisa quem ler depois que uma decisão já tomada pode ter sido reaberta aqui |
| `Solicitação original` | Texto exato do usuário. Não reescreva nem "melhore" o pedido — é a âncora contra deriva de escopo |
| `Objetivo e resultado esperado` | Resultado observável, quem é afetado e qual valor entrega. Se não dá para observar, não é objetivo |
| `Fatos confirmados` | Só o que foi verificado no repositório ou em fonte disponível, **com a fonte**. Inferência não entra aqui |
| `Escopo > Incluído` | Entregas, não tarefas. Cada linha deve ser verificável no fim |
| `Escopo > Fora do escopo` | O limite explícito é o que impede o escopo de crescer sozinho depois. Seção vazia é sinal de entrevista rasa |
| `Decisões` | Uma linha por decisão material. `Origem` distingue o que o usuário escolheu do que foi inferido do repositório — essa coluna é o que torna o plano auditável |
| `Fluxos, regras e casos-limite` | Comportamento esperado **e** as exceções. Fluxo alternativo e caminho de erro pertencem aqui |
| `Contratos, dados e compatibilidade` | Interfaces, estados, migração. `não aplicável` é resposta válida, silêncio não |
| `Plano de ação` | Mudanças verificáveis com dependências reais. Não é cronograma nem lista de arquivos |
| `Critérios de aceite e verificação` | Cenário observável que prova a entrega. "Funciona", "testado" e "código limpo" não são critérios |
| `Riscos, pressupostos e mitigação` | Todo pressuposto não confirmado é risco. Registre a mitigação junto |
| `Decisões pendentes` | `Nenhuma.` quando a fronteira esvaziou. Qualquer outra coisa significa **plano incompleto** — diga isso ao usuário |
| `Histórico de iterações` | Uma linha por rodada de ajuste, com data e efeito. Nunca sobrescreva uma entrada anterior |

## O que nunca vai para o plano

| Item | Por quê |
|---|---|
| Decisão que o usuário não tomou | O plano é registro de acordo, não de sugestão. Se ninguém decidiu, é pendência |
| Preferência técnica apresentada como requisito de negócio | Confunde quem implementa sobre o que é negociável |
| Requisito, integração ou refatoração fora da solicitação | Se merece existir, é uma opção de escopo a ser decidida, não um item a ser assumido |
| Critério subjetivo | Não é verificável, logo não encerra a entrevista |
