Execute a skill `implement` para construir o que o spec decidiu, em fatias verticais testadas nas seams acordadas. `$ARGUMENTS` é opcional: o caminho de um spec específico, ou um recorte do trabalho a implementar. Vazio, use o spec mais recente após confirmar com o usuário qual é o alvo.

Esta skill **não decide**: o que implementar já foi decidido no `/planning` e registrado pelo `/to-spec`. Lacuna material interrompe a implementação.

## Passos

### 1. Ancorar no spec e retomar o andamento

- Ler o spec inteiro. `Fora do escopo` vincula tanto quanto `Decisões de implementação`
- Ler `.coder/impl-<branch>.md` da branch atual quando existir: fatia concluída não se repete
- Confirmar que as seams estão acordadas em `Decisões de teste` — sem isso, parar e pedir confirmação
- Descobrir os comandos reais de teste, typecheck e lint do projeto; não inventar comando que não existe

### 2. Cortar as fatias verticais

- Derivar as fatias das histórias e das seams; cada uma atravessa uma seam de ponta a ponta
- Ordenar por dependência real e por incerteza: a fatia que mais ensina vem primeiro
- Apresentar a lista em até 10 linhas e seguir

### 3. Executar o ciclo, uma fatia por vez

- Rodar o ciclo vermelho → verde da skill `tdd` na seam acordada
- Depois de cada fatia: typecheck mais o arquivo de teste da fatia, nunca a suíte inteira
- Atualizar `.coder/impl-<branch>.md` ao fechar cada fatia, antes de abrir a próxima
- Decisão material não tomada interrompe: registrar a pendência e trazê-la ao usuário

### 4. Verificar o conjunto

- Suíte completa uma vez, mais typecheck e lint
- Falha fora das fatias tocadas é regressão; teste pulado conta como falha

### 5. Revisar

- Executar a skill `code-review` sobre o merge-base com a branch base
- Corrigir defeito e desvio do spec; achado que exige decisão vira pendência no relatório

### 6. Fechar sob confirmação

- Atualizar o arquivo de andamento com estado final, pendências e decisões de execução
- Preparar o commit: stage explícito arquivo a arquivo, revisão da diff preparada, mensagem em Conventional Commits
- **Aguardar confirmação explícita.** Sem ela, nada de commit, push, branch ou reset
- Commit executado e trabalho que atravessou sessões: sugerir `/to-memory` em uma linha. Sugerir, nunca executar
