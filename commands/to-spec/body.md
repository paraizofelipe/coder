Execute a skill `to-spec` para transformar as decisões já tomadas em um spec gravado em `.coder/spec-AAAAMMDD-HHMMSS.md`. `$ARGUMENTS` é um foco/escopo opcional (recorte do plano a especificar); se vazio, cubra todas as decisões disponíveis.

Esta skill **não entrevista**: ela registra decisões que já existem. A única pergunta permitida é a confirmação das seams.

## Passos

### 1. Reunir a base de decisões

- Ler `.coder/plan-<branch>.md` da branch atual quando existir — é a fonte canônica, produzida por `/planning`
- Sem plano para a branch atual, sintetizar do contexto da conversa. Nunca usar o plano de outra branch
- Sem nenhum dos dois com decisões tomadas, parar e instruir o usuário a executar `/planning` antes
- Classificar cada item em decisão registrada, inferência do repositório ou lacuna; lacuna nunca vira decisão

### 2. Ancorar no estado real do código

- Explorar apenas o necessário para nomear módulos, interfaces e integrações com precisão
- Usar o vocabulário de domínio do projeto e respeitar ADRs e convenções da área tocada
- Não revalidar fatos que o plano já confirmou

### 3. Propor as seams e confirmar

- Propor as fronteiras de teste: preferir existentes, a mais alta possível, o ideal é uma
- Apresentar em até 10 linhas e aguardar confirmação explícita — sem ela, não escrever o spec

### 4. Escrever o spec

- Gravar `.coder/spec-AAAAMMDD-HHMMSS.md` com as seções Problema, Solução, Histórias de usuário, Decisões de implementação, Decisões de teste, Fora do escopo e Notas adicionais
- Sem caminhos de arquivo e sem snippets, exceto trecho que codifica uma decisão
- Ajuste na mesma sessão atualiza o mesmo arquivo e acrescenta `## Histórico de iterações`

### 5. Reportar

- Resumo de até 15 linhas: caminho, objetivo, nº de histórias, seams acordadas, fora do escopo e pendências
- Não despejar o documento na resposta
- Pendência material → dizer explicitamente que o spec não está pronto e recomendar `/planning`
