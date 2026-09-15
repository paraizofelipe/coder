---
name: to-spec
description: Sintetiza decisões já tomadas (de `.coder/plan-<branch>.md` ou da conversa atual) em um spec executável gravado em `.coder/spec-AAAAMMDD-HHMMSS.md`. Não entrevista — registra. Use depois da skill `planning`, quando o trabalho não cabe em uma única sessão e precisa sobreviver ao fim da janela de contexto.
---

<role>
Você está executando a skill `to-spec`. Transforma decisões **já tomadas** em um spec: o documento que sobrevive quando a janela de contexto acaba. Uma sessão nova deve conseguir retomar o trabalho a partir dele sem que o usuário reexplique nada.

Esta skill não decide e não valida — ela registra. Toda afirmação do spec precisa ser rastreável a uma decisão do usuário, a um fato do repositório ou a uma inferência declarada como tal.
</role>

<context>
**Entrada**, nesta ordem de precedência:

1. `.coder/plan-<branch>.md` — o plano da **branch atual**, produzido pela skill `planning`; fonte canônica quando existe
2. O contexto da conversa atual, quando não há plano para a branch atual

Se não houver nenhum dos dois com decisões efetivamente tomadas, **pare** e instrua o usuário a executar `/planning` antes.

Nunca use o plano de outra branch. Ele descreve outro trabalho, e um spec derivado dele sai plausível e errado — falha pior do que não produzir spec nenhum.

**Saída** — `.coder/spec-AAAAMMDD-HHMMSS.md`, onde `AAAAMMDD-HHMMSS` é a data e a hora locais da criação (ex.: `.coder/spec-20260915-143012.md`).

**Quando usar** — quando o trabalho atravessa várias sessões. Se ele cabe em uma única janela de contexto, o spec não paga o próprio custo: implemente direto.
</context>

<workflow>

### 1. Reúna a base de decisões
- Leia o plano da branch atual quando existir e mapeie suas seções com `references/plan-to-spec-map.md`, que traz a regra de nome do arquivo. Quando não existir, sintetize do contexto da conversa.
- Classifique cada item em **decisão registrada**, **inferência do repositório** ou **lacuna**.
- Lacuna nunca vira decisão. Ou entra em `Fora do escopo` com o motivo, ou vira pendência explícita — e, se for material, interrompe o spec (passo 5).
- Não reabra decisões fechadas nem faça novas perguntas de escopo aqui.

### 2. Ancore no estado real do código
- Explore o repositório apenas o necessário para nomear módulos, interfaces e pontos de integração com precisão.
- Use o vocabulário de domínio do projeto (glossário, `CONTEXT.md`, README, nomes já presentes no código). Não introduza sinônimos.
- Respeite ADRs e convenções vigentes na área tocada. Se uma decisão registrada conflitar com alguma delas, registre o conflito em `Notas adicionais` — não o resolva sozinho.
- Não revalide fatos que o plano já confirmou.

### 3. Proponha as seams e confirme — única interação permitida
- Consulte `references/seams.md` para o conceito, as heurísticas de escolha e o formato de apresentação.
- Prefira seams existentes a novas; escolha a mais alta possível; menos é melhor — o ideal é uma.
- Apresente a proposta em até 10 linhas e aguarde confirmação. **Sem confirmação, não escreva o spec.**

### 4. Escreva o spec
- Siga o template de `references/spec-format.md`, sem alterar a ordem nem os títulos das seções.
- Grave em `.coder/spec-AAAAMMDD-HHMMSS.md`. Nunca fora de `.coder/`, nunca com outro padrão de nome.
- Ajuste pedido na **mesma sessão** → atualize o **mesmo arquivo** e acrescente um bloco `## Histórico de iterações` com data e motivo, preservando o conteúdo anterior. Nova solicitação → novo arquivo com timestamp próprio.

### 5. Apresente o resumo e pare
- Resumo de até 15 linhas, no formato de `<output_format>`. Não despeje o documento.
- Se restou decisão material pendente, diga isso explicitamente, aponte o que falta e recomende voltar ao `/planning`. Não marque o spec como pronto.

</workflow>

<rules>
- **Não entreviste.** A confirmação das seams (passo 3) é a única pergunta permitida. Decisões de produto, escopo e arquitetura já foram tomadas antes desta skill.
- **Decisão inventada é defeito.** Se o spec afirma algo que ninguém decidiu e o repositório não comprova, remova ou marque como pendência.
- **Sem caminhos de arquivo e sem snippets de código.** Eles envelhecem antes de o spec ser consumido. Exceção única: um trecho que codifica uma decisão com mais precisão do que a prosa (máquina de estados, schema, shape de tipo, contrato) — mínimo, aparado ao que decide, com nota de origem.
- **Nenhum código de produção, nenhuma branch, nenhum commit.** Esta skill escreve exclusivamente `.coder/spec-AAAAMMDD-HHMMSS.md`.
- **Não fatie em tasks aqui.** O spec descreve o destino e o porquê; a sequência de execução pertence à etapa seguinte do fluxo.
- **Sem meta-texto.** O arquivo segue o template literal, sem comentários explicando o processo desta skill.
- Histórias de usuário e critérios descrevem comportamento observável. "Funciona", "rápido" e "testado" não são critérios.
</rules>

<checklist>
- [ ] A base veio do plano da branch atual ou de decisões efetivamente tomadas na conversa.
- [ ] Nenhum plano de outra branch foi usado como base.
- [ ] Nenhuma seção foi preenchida com decisão que o usuário não tomou.
- [ ] As seams foram propostas e confirmadas antes da escrita.
- [ ] O spec usa o vocabulário de domínio do projeto.
- [ ] `Fora do escopo` tem conteúdo real — o que foi recusado é a parte mais útil do documento.
- [ ] Não há caminho de arquivo nem snippet fora da exceção de decisão codificada.
- [ ] `.coder/spec-AAAAMMDD-HHMMSS.md` existe e segue o template.
- [ ] O resumo apresentado tem no máximo 15 linhas.
</checklist>

<output_format>
Após gravar o arquivo, responda **apenas** com:

```text
Spec consolidado em `.coder/spec-AAAAMMDD-HHMMSS.md`.

- Objetivo: <resultado observável em 1 linha>
- Histórias de usuário: <N>
- Seams acordadas: <lista curta>
- Decisões registradas: <N> — <as 2 ou 3 de maior impacto>
- Fora do escopo: <resumo em 1 linha>
- Pendências: <o que falta decidir, ou "nenhuma">

<próximo passo recomendado, em 1 linha>
```

Quando houver pendência material, substitua o próximo passo por: `Decisão pendente — execute /planning antes de consumir este spec.`
</output_format>
