---
name: tdd
description: Ciclo vermelho → verde em fatias verticais, testando apenas nas seams acordadas. Define o que é um bom teste, onde ele vive e os anti-padrões que o tornam inútil. Use ao implementar comportamento novo ou corrigir bug test-first, dentro da skill `implement` ou isoladamente.
---

<role>
Você está executando a skill `tdd`. Conduz o ciclo vermelho → verde de uma fatia e garante que o teste resultante valha a pena guardar.

Um teste vale a pena quando descreve comportamento que alguém se importa em preservar. Um teste que quebra em toda refatoração sem que o comportamento tenha mudado é custo, não rede de segurança.
</role>

<context>
**Entrada** — uma fatia vertical e a seam onde ela será observada. A seam vem acordada do spec (`Decisões de teste`) ou confirmada com o usuário antes de qualquer teste.

**Saída** — teste e implementação da fatia, nesta ordem, com o vermelho observado antes do verde.

**Quando usar** — comportamento novo, correção de bug e qualquer mudança cujo efeito seja observável em uma fronteira pública. Não se aplica a mudança puramente estrutural sem efeito observável: isso é refatoração, protegida pelos testes que já existem.

**Antes de começar**, alinhe o vocabulário: leia `CONTEXT.md`, `README` ou o glossário do projeto, para que o nome do teste e o da interface usem as palavras do domínio, não sinônimos novos. Respeite os ADRs da área tocada.
</context>

<workflow>

### 1. Fixe a seam e o comportamento
- Nomeie a seam onde o teste vai viver e o comportamento observável que a fatia entrega. Se a seam não estiver acordada, **pare e confirme** — nenhum teste é escrito em seam não acordada.
- Procure **prior art**: testes equivalentes já existentes no projeto. Fixtures, helpers e convenção de nome vêm deles, não de você.
- Escreva o nome do teste como uma frase de especificação: o que passa a ser possível, não como o código faz.

### 2. Vermelho
- Escreva **um** teste, na seam, para **um** comportamento.
- Execute-o e observe a falha. Confira que ela falha pelo motivo esperado: um teste que passa antes da implementação, ou que falha por erro de importação, não é vermelho — é ruído.
- O valor esperado vem de fonte independente: um literal conhecido, um exemplo do spec, um caso trabalhado à mão. Nunca recalculado do jeito que o código calcula.
- Consulte `references/test-quality.md` ao decidir o que asserir, e `references/mocking.md` quando a fatia tocar uma fronteira externa.

### 3. Verde
- Escreva o **mínimo** que faz o teste passar. Não antecipe o próximo teste nem adicione caminho que nenhum teste exercita.
- Execute o mesmo teste e observe o verde. Depois execute o arquivo inteiro, para flagrar o que a fatia possa ter quebrado por perto.
- Verde obtido enfraquecendo o teste não é verde. Se a asserção foi afrouxada para passar, o ciclo falhou.

### 4. Feche o ciclo
- Uma fatia, um ciclo. Volte ao passo 1 para a próxima fatia, levando o que esta ensinou — inclusive quando o aprendizado for que o corte seguinte mudou.
- Não refatore aqui. Estrutura se ajusta depois do verde, na revisão, com os testes já verdes protegendo a mudança.

</workflow>

<rules>
- **Vermelho antes de verde.** Sem falha observada, não há prova de que o teste testa alguma coisa.
- **Só em seam acordada.** Não se testa tudo; acordar as seams antes é o que faz o esforço cair no caminho crítico em vez de em toda borda.
- **Uma fatia por ciclo.** Escrever todos os testes primeiro e depois toda a implementação verifica comportamento imaginado, não real.
- **Nada de detalhe de implementação.** Sem método privado, sem contagem de chamadas, sem verificação por canal lateral (consultar o banco em vez de usar a interface).
- **Nada de teste tautológico.** Se a asserção recalcula o resultado como o código recalcula, ela passa por construção e nunca discordará dele.
- **Mock apenas em fronteira de sistema.** Colaborador interno mockado vira o objeto sob teste.
- **Refatoração fica fora do ciclo.** Ela pertence à revisão.
- Teste pulado ou removido para obter verde é falha reportável, nunca um detalhe de execução.
</rules>

<checklist>
- [ ] A seam estava acordada antes do primeiro teste.
- [ ] O prior art do projeto foi consultado antes de criar convenção nova.
- [ ] O vermelho foi observado, e pelo motivo esperado.
- [ ] O valor esperado veio de fonte independente do código sob teste.
- [ ] O verde saiu do mínimo necessário, sem caminho especulativo.
- [ ] Nenhuma asserção foi afrouxada para obter verde.
- [ ] Nenhum colaborador interno foi mockado.
- [ ] O teste sobrevive a uma refatoração que não mude comportamento.
</checklist>

<output_format>
Uma linha por ciclo, durante a execução:

```text
<n>. <comportamento> @ <seam> — vermelho: <motivo da falha> → verde: <o que foi implementado>
```

Ao encerrar a sequência de ciclos:

```text
Ciclos concluídos: <N>.

- Seams exercitadas: <lista>
- Testes adicionados: <N> em <arquivos>
- Última execução: <comando> — <N passaram, N falharam, N puladas>
- Pendências: <o que ficou sem cobertura e por quê, ou "nenhuma">
```
</output_format>
