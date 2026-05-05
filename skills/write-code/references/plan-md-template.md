# Template `.coder/plan.md`

Estrutura obrigatória do arquivo `.coder/plan.md` criado pelo `coder` na fase de planejamento.

```markdown
# Plano de Implementação

## Solicitação original
[texto exato da solicitação do usuário]

## Resumo da análise
[estrutura do projeto, padrões relevantes, áreas impactadas]

## Ambiguidades identificadas
| # | Questão | Status | Decisão |
|---|---------|--------|---------|
| 1 | [descrição] | ⏳ Pendente | — |

## Plano de ação
- [ ] [o que será feito e por quê]

## Riscos e pontos de atenção
- [lista]
```

## Regras de manutenção do arquivo

- Se `.coder/plan.md` já existir, **atualizá-lo** em vez de substituir
- Cada decisão de ambiguidade deve mover o status de `⏳ Pendente` para `✅ Resolvida` e preencher a coluna **Decisão**
- Atualizar a seção **Plano de ação** sempre que uma decisão alterar o escopo ou abordagem
- Se nenhuma ambiguidade for identificada pelo `analyzer`, registrar explicitamente: `Nenhuma ambiguidade identificada — solicitação clara.`
- Somente prosseguir para a próxima etapa quando todas as ambiguidades estiverem `✅ Resolvida`
