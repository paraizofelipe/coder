Acione o agente `qa` com a skill `validate-implementation` para validar funcionalmente as modificações da branch atual. `$ARGUMENTS` é um foco/escopo opcional (regra de negócio ou fluxo a priorizar); se vazio, valide todas as modificações da branch.

## Passos

### 1. Enquadrar o alvo

- Identificar a branch atual (`git rev-parse --abbrev-ref HEAD`)
- Se for `main`/`master`, alertar que não há feature a validar e pedir orientação
- Anunciar em 1 linha o alvo da validação (branch, base e foco de `$ARGUMENTS`)

### 2. Acionar o `analyzer`

- Mapear o diff da branch contra a base e as regras de negócio/fluxos afetados
- Descobrir os serviços e dependências que o projeto acessa (env, configs, compose, clientes HTTP/SDK, conexões de banco, MCPs, manifests)

### 3. Montar o plano de testes

- Derivar casos por regra de negócio (smoke, black-box, e2e, regressão) com tipo, objetivo, pré-condições, passos, dados, resultado esperado, serviços exigidos, classe (`leitura`/`mutação`) e ambiente (`hml`)

### 4. Consolidar serviços e acessos

- Listar cada serviço com método de acesso (MCP/CLI/curl) e classe de uso

### 5. Salvar o plano e pedir aprovação

- Gravar `.coder/tests-<timestamp>.md` (mesmo padrão do `.coder/task-*.md`)
- Apresentar resumo (≤15 linhas) e aguardar "sim" explícito antes de prosseguir

### 6. Validar os acessos

- Probe read-only por serviço, mostrando o comando exato antes de executar (acionar o `infra` para ambiente/logs)
- Montar a tabela `OK`/`FALTA`/`FALHA` + remediação; casos sem acesso essencial ficam `BLOQUEADO`

### 7. Executar os testes

- Apenas acessos `OK`, em HML; `leitura` roda direto; `mutação` exige confirmação por item
- Registrar evidência real por caso (request/response, query/result, trecho de log)

### 8. Reportar

- Status por caso (`PASSOU`/`FALHOU`/`BLOQUEADO`/`PULADO`) com evidência
- Para falhas, indicar a regra de negócio violada e encaminhar pendências de correção ao `coder`
