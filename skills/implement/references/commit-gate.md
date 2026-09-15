# Portão de commit

Carregue no **passo 6**. O commit é do usuário; esta skill apenas o prepara e apresenta.

## Regra

Nenhuma operação que altere o histórico ou o estado compartilhado acontece sem um "sim" explícito na conversa:

| Operação | Sem confirmação |
|---|---|
| `commit` | Proibida |
| `push`, `push --force` | Proibida |
| `branch`, `switch`, `checkout` | Proibida |
| `merge`, `rebase`, `reset`, `stash` | Proibida |
| `add` de arquivos identificados | Permitida — é o preparo que torna a revisão possível |

Autorização vale para a operação apresentada, não para as próximas. Um "pode commitar" não autoriza o push.

## Preparo

### 1. Stage explícito

Adicione os arquivos **um a um**, pelos caminhos que você mesmo tocou:

```bash
git add <caminho> <caminho>
```

Nunca `git add .`, `git add -A` nem `git add -u`: eles arrastam para o commit arquivo que você não conhece — resto de outra tarefa, artefato local, segredo.

### 2. Revise o que foi preparado

```bash
git diff --cached --name-only
git diff --cached
```

Antes de apresentar, procure nesta diff:

| Procure | Por quê |
|---|---|
| Credencial, token, `.env`, chave, dado pessoal | Commit publica; remoção posterior não desfaz |
| `console.log`, `print`, `debugger`, breakpoint | Sobra de investigação, não de implementação |
| Código comentado e `TODO` recém-criado | Ou resolve, ou vira pendência registrada — não fica no meio |
| Arquivo que você não tocou nesta implementação | Entrou por engano no stage |
| Teste pulado, `skip`, `only`, asserção comentada | Verificação desligada silenciosamente |
| Mudança fora do escopo do spec | Refatoração oportunista que ninguém pediu |

Achou algo? Tire do stage e resolva antes de apresentar. Não apresente um commit que você mesmo apontaria na revisão.

### 3. Mensagem

Conventional Commits, assunto em inglês, minúsculo, abaixo de 70 caracteres, sem ponto final:

```text
<tipo>: <o que muda, no imperativo>
```

| Tipo | Quando |
|---|---|
| `feat` | Comportamento novo visível para quem usa |
| `fix` | Correção de comportamento errado |
| `refactor` | Estrutura muda, comportamento não |
| `test` | Só testes |
| `docs` | Só documentação |
| `chore` | Build, dependência, configuração |

O corpo é opcional e serve para o **porquê**, não para listar arquivos — a diff já os lista. Quando o spec existir, cite-o no corpo em uma linha.

Um commit por ideia. Se a mensagem precisa de "e" para descrever o que fez, provavelmente são dois commits.

## Apresentação

```text
Commit preparado:

  <tipo>: <assunto>

Arquivos (<N>):
  <caminho>
  <caminho>

Confirma para eu executar?
```

Sem resposta afirmativa, o trabalho fica preparado no stage e a skill encerra assim mesmo. Trabalho staged não é trabalho perdido.
