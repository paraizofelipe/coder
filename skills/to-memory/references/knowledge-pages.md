# Páginas de conhecimento

Carregue no **passo 5**, depois que os documentos já foram enviados.

## O que é uma página

Uma página **não guarda conteúdo curado**. Ela guarda uma pergunta permanente — a `source_query` — que o servidor re-responde a cada consolidação, reconstruindo o conteúdo a partir das observações e dos documentos acumulados.

A consequência prática inverte o instinto: você não escreve a memória, você **escreve a pergunta certa** e deixa a resposta ser reconstruída conforme a evidência chega.

| | Documento | Página |
|---|---|---|
| O que você fornece | O artefato bruto | Uma pergunta |
| Quem escreve o conteúdo | Você | O servidor, a cada consolidação |
| Quantos existem | Um por artefato, muitos | **Poucos e estáveis** |
| Muda quando | Você reenvia | A evidência muda |

## As perguntas permanentes deste fluxo

Cada uma se apoia em seções que os templates do fluxo **já produzem**, o que é o que as torna respondíveis:

| `page_id` | `source_query` | Alimentada por |
|---|---|---|
| `decisoes-arquiteturais` | Quais decisões de arquitetura foram tomadas neste trabalho, com que justificativa e que trade-off aceito? | `Decisões` do plano, `Decisões de implementação` do spec |
| `escopo-recusado` | O que foi explicitamente recusado de escopo, e por qual motivo? | `Fora do escopo` do plano e do spec |
| `decisoes-revertidas` | Que decisões registradas em plano ou spec foram contrariadas durante a implementação, e o que se aprendeu com isso? | `Desvios do spec` do andamento |
| `seams-do-projeto` | Que fronteiras de teste foram acordadas, com que justificativa, e quais se mostraram acertadas depois? | `Decisões de teste` do spec, `Seams acordadas` do andamento |

A coluna `Origem` da tabela de decisões é o filtro que dá valor à primeira página: decisão tomada pelo **usuário** é durável; inferência do repositório envelhece junto com o código que a sustentava.

## Criar, atualizar ou deixar quieto

1. **Liste as páginas existentes antes de qualquer coisa.** O bank é compartilhado entre projetos; páginas de bootstrap com nome genérico podem já ocupar o espaço.
2. Página que já cobre a pergunta → **atualize** a `source_query` se a redação melhorou, ou deixe como está. Nunca crie uma segunda para a mesma pergunta.
3. Página ausente → crie, com o `page_id` e a `source_query` da tabela acima.
4. Em dúvida entre atualizar e criar, **atualize**. Duplicata divide a evidência entre duas páginas e piora as duas.

## O que nunca vira página

| Item | Onde ele pertence |
|---|---|
| Uma página por feature, branch ou trabalho | Isso é documento — o passo 4 |
| Conteúdo já sintetizado por você | Página guarda pergunta; a resposta é do servidor |
| Pergunta que só tem uma resposta possível e fixa | Se não vai mudar com a evidência, é documento |
| Pergunta ampla demais ("o que sei sobre este projeto?") | Não converge; a página fica genérica e inútil |

## Sobre escopo por projeto

As páginas existentes podem carregar tags de projeto, mas a ferramenta de criação exposta ao agente aceita apenas `page_id`, `name` e `source_query` — **tag não é definível por aqui**. Quando o escopo importar, ele precisa estar no `page_id` e no texto da própria pergunta.
