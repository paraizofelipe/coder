# Formato do `.coder/tests-AAAAMMDD-HHMMSS.md`

Template literal do artefato gravado no Passo 5 da skill. O nome usa a data e hora locais da criação (ex.: `.coder/tests-20260618-143012.md`). O arquivo fica em `.coder/` (gitignored, igual ao `.coder/task-*.md` do `lead`).

```markdown
# Plano de Testes de QA

## Alvo
- Branch: <branch> (base: <base>)
- Foco: <regras de negócio / fluxos>

## Contexto técnico
[1 parágrafo consolidado do analyzer: o que mudou no comportamento observável e quais regras de negócio são afetadas]

## Serviços e acessos necessários
| Serviço | Tipo | Método de acesso | Classe |
|---|---|---|---|
| <id/alias> | DB / API / ambiente-logs / fila | MCP <x> / CLI <y> / curl | leitura / mutação |

## Plano de testes

### T1 — <título>
- Tipo: smoke | black-box | e2e | regressão
- Objetivo: <regra de negócio que valida>
- Pré-condições: <estado/dados necessários>
- Passos:
  1. <ação>
  2. <ação>
- Dados: <inputs / fixtures>
- Resultado esperado: <saída observável>
- Serviços: <lista de serviços exigidos>
- Classe: leitura | mutação
- Ambiente: hml

### T2 — <título>
[mesma estrutura]

## Riscos
- <risco de cobertura, de ambiente, de dados>

## Histórico de iterações
- [AAAA-MM-DD] <ajuste solicitado nesta sessão>
```

Regras do template:
- Cada caso recebe um id sequencial `T1`, `T2`, … e marca explicitamente a **classe** (`leitura`/`mutação`).
- O bloco `## Histórico de iterações` só é adicionado quando o arquivo é atualizado na mesma sessão; preservar decisões anteriores.
- Não adicionar meta-texto explicando o processo — apenas o plano.
