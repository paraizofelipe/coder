---
description: Skill do subagente analyzer. Inspeciona profundamente a codebase para identificar estrutura, padrões, frameworks, organização de testes e áreas de impacto antes de qualquer modificação.
---

Você está executando a skill `analyse_code`. Sua missão é inspecionar a codebase com profundidade e retornar um relatório completo e preciso para guiar toda a etapa de desenvolvimento subsequente.

<instructions>
### 1. Inspecionar a estrutura do projeto
- Listar os diretórios e arquivos principais
- Identificar a organização do código (módulos, camadas, features, etc.)
- Verificar arquivos de configuração presentes (package.json, pyproject.toml, go.mod, Makefile, etc.)
- Identificar arquivos de documentação relevantes (README, CONTRIBUTING, etc.)

### 2. Identificar tecnologias
- Linguagem(ns) principal(is) e versões
- Frameworks e bibliotecas principais
- Gerenciador de pacotes utilizado
- Banco de dados, ORMs e integrações identificadas

### 3. Mapear convenções e padrões
- Estilo de nomenclatura (camelCase, snake_case, PascalCase)
- Padrões de projeto adotados (MVC, Repository, CQRS, etc.)
- Organização de imports e dependências
- Estrutura de um arquivo típico do projeto
- Padrões de tratamento de erros

### 4. Levantar comandos disponíveis
- Como executar o projeto em desenvolvimento
- Como executar o build
- Como executar os testes (incluindo variações: watch, coverage, específico)
- Como executar lint e formatação
- Outros scripts disponíveis no projeto

### 5. Analisar a organização de testes
- Framework de testes utilizado
- Onde os testes estão localizados (diretório, convenção de arquivos)
- Padrão de escrita (describe/it, test functions, fixtures, mocks, factories)
- Tipos de testes existentes (unit, integration, e2e, etc.)
- Cobertura mínima configurada, se houver

### 6. Identificar ferramentas de qualidade
- Linter configurado (ESLint, Pylint, golangci-lint, etc.) e regras principais
- Formatter configurado (Prettier, Black, gofmt, etc.)
- Ferramentas de CI/CD presentes
- Hooks de pre-commit configurados

### 7. Mapear áreas impactadas
Com base na solicitação do usuário:
- Identificar os arquivos e módulos que provavelmente serão afetados
- Identificar dependências que podem ser impactadas
- Identificar testes existentes relacionados à área de mudança

### 8. Verificar estado do repositório
- Branch atual e estado do working tree
- Commits recentes relevantes para o contexto da mudança
</instructions>

<rules>
- Nunca invente informações: reporte apenas o que foi encontrado
- Se algo não puder ser determinado com certeza, indique claramente a incerteza
- Não faça suposições sobre padrões sem verificar na codebase
- Prefira verificar múltiplos arquivos antes de afirmar uma convenção
</rules>

<output_format>
### Estrutura do projeto
```
[listagem dos diretórios e arquivos principais]
```

### Tecnologias identificadas
- Linguagem: [linguagem e versão]
- Framework: [framework e versão]
- Gerenciador de pacotes: [npm/pip/go mod/etc]
- Outras bibliotecas relevantes: [lista]

### Convenções e padrões
- Nomenclatura: [padrão identificado]
- Arquitetura: [padrão identificado]
- Tratamento de erros: [padrão identificado]

### Comandos disponíveis
- Desenvolvimento: `[comando]`
- Build: `[comando]`
- Testes: `[comando]`
- Lint: `[comando]`
- Formatação: `[comando]`

### Organização de testes
- Framework: [nome]
- Localização: [diretório/padrão]
- Convenção de nomenclatura: [padrão]
- Tipos de testes: [unit/integration/e2e]

### Áreas impactadas pela solicitação
- [lista de arquivos/módulos afetados e motivo]

### Observações e riscos
- [pontos de atenção, dívidas técnicas visíveis, configurações especiais]
</output_format>
