Execute a skill `code-review` para revisar a diff entre `HEAD` e um ponto fixo, em dois eixos independentes. `$ARGUMENTS` é o ponto fixo: commit, branch, tag ou `main`. Vazio, pergunte ao usuário — é a única pergunta obrigatória desta skill.

Os dois eixos não se substituem. Código impecável que implementa a coisa errada passa em Padrões e falha em Spec; código que faz o pedido violando toda convenção falha em Padrões e passa em Spec.

Esta skill **não corrige**: ela aponta, com evidência.

## Passos

### 1. Fixar o ponto de comparação

- Resolver a ref (`git rev-parse`) antes de qualquer análise; ref inválida falha aqui
- Capturar a diff uma vez, com três pontos: `git diff <ponto-fixo>...HEAD`, mais `git log <ponto-fixo>..HEAD --oneline`
- Diff vazia encerra a revisão com essa constatação

### 2. Localizar o spec

- Nesta ordem: spec citado pela `implement` ou pelo `.coder/impl-<branch>.md`, caminho passado como argumento, `.coder/plan-<branch>.md` da branch, referência a issue nos commits
- Sem nenhum, perguntar. Se não houver spec, **pular o eixo Spec** e registrar isso — não substituir por suposição

### 3. Localizar os padrões

- `AGENTS.md`, `CLAUDE.md`, `CONTRIBUTING.md`, `CODING_STANDARDS.md`, `.agents/rules/`, ADRs
- Ignorar o que formatador, linter e typecheck já garantem

### 4. Rodar os dois eixos

- **Spec**: requisito faltando, escopo indevido e implementado errado — cada achado citando a linha do spec
- **Padrões**: violação de convenção documentada, mais a linha de base de smells como julgamento. O repositório vence a linha de base
- Com subagentes disponíveis, rodar em paralelo, cada um cego para o outro eixo; sem eles, em sequência, um eixo inteiro por vez

### 5. Agregar sem reordenar

- Duas seções separadas, `Spec` e depois `Padrões`, sem fundir nem escolher vencedor entre eixos
- Encerrar com o total por eixo e o pior achado dentro de cada um
