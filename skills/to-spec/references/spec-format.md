# Formato do `.coder/spec-AAAAMMDD-HHMMSS.md`

Carregue no **passo 4**. Sete seções, nesta ordem, com títulos literais. O `Histórico de iterações` só aparece a partir do primeiro ajuste.

Os títulos em português correspondem ao spec canônico assim: `Problema` = Problem Statement, `Solução` = Solution, `Histórias de usuário` = User Stories, `Decisões de implementação` = Implementation Decisions, `Decisões de teste` = Testing Decisions, `Fora do escopo` = Out of Scope, `Notas adicionais` = Further Notes.

## Template

````markdown
# Spec — <título curto do que será construído>

> Base: `.coder/plan-<branch>.md` | conversa da sessão · Branch: `<branch atual>` · <AAAA-MM-DD HH:MM>

## Problema

O problema, na perspectiva de quem o sofre. Quem é afetado, o que não consegue fazer hoje e qual o custo disso. Sem solução aqui.

## Solução

A solução, na perspectiva de quem usa. O que passa a ser possível depois da entrega. Sem plano de execução aqui.

## Histórias de usuário

Lista numerada, longa e exaustiva, cobrindo todos os aspectos do que foi decidido — inclusive fluxos alternativos, estados de erro e casos-limite.

1. Como <ator>, quero <capacidade>, para que <benefício>
2. Como <ator>, quero <capacidade>, para que <benefício>

## Decisões de implementação

Decisões já tomadas, agrupadas por tema. Cada uma com a escolha e o porquê:

- **<tema>** — <escolha>. <justificativa e trade-off aceito>

Cobre, quando aplicável: módulos criados ou modificados, interfaces afetadas, decisões arquiteturais, mudanças de schema, contratos de API e eventos, compatibilidade e migração, interações específicas.

## Decisões de teste

### Seams acordadas

| Seam | Existente / nova | O que observa |
|---|---|---|
| <nome> | existente | <comportamento verificável> |

### O que caracteriza um bom teste aqui

- Observa comportamento externo na seam acordada, nunca detalhe de implementação
- <critério específico deste spec>

### Prior art

- <testes equivalentes já existentes no projeto que servem de modelo>

## Fora do escopo

- <limite explícito> — <motivo pelo qual foi recusado>

## Notas adicionais

- Riscos e pressupostos, com mitigação
- Conflitos detectados entre decisão e ADR ou convenção vigente
- Pendências: `Nenhuma.` ou a decisão que falta e por que ela bloqueia

## Histórico de iterações

- <AAAA-MM-DD> — <motivo do ajuste e o que mudou>
````

## Regras de preenchimento

| Seção | Cuidado |
|---|---|
| `Problema` | Se descrever a solução, está errada. O problema existe antes de qualquer decisão técnica |
| `Solução` | Perspectiva de uso, não de implementação. "O usuário passa a ver X", não "o serviço Y chama Z" |
| `Histórias de usuário` | Exaustiva é o ponto. Fluxo alternativo, erro e caso-limite são histórias, não notas de rodapé |
| `Decisões de implementação` | Sem caminho de arquivo e sem snippet. Nomeie módulo e interface pelo vocabulário do domínio |
| `Decisões de teste` | As seams vêm confirmadas do passo 3. Nunca escreva esta seção com seam não acordada |
| `Fora do escopo` | Seção vazia é sinal de spec fraco — o que foi recusado costuma ser a informação mais útil |
| `Notas adicionais` | Pendência material aqui significa spec não pronto. Diga isso no resumo |

## Exceção de snippet

A regra é prosa. A única exceção é um trecho que codifica a decisão com mais precisão do que qualquer texto — máquina de estados, shape de tipo, schema, contrato. Nesses casos:

- Inclua apenas as linhas que carregam a decisão, nunca um exemplo funcional
- Coloque-o dentro da decisão a que pertence, não em seção própria
- Anote a origem em uma linha (`decidido em /planning`, `extraído do contrato atual`)

## Trabalho arquitetural

Refactor e mudança de fronteira de módulo cabem mal em histórias de usuário. Nesses casos, concentre o conteúdo em `Decisões de implementação` e `Decisões de teste`, e mantenha `Histórias de usuário` curta, com as poucas histórias reais que existirem — não invente histórias para preencher a seção.
