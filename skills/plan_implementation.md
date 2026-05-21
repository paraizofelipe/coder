---
description: Skill principal do agente lead. Orquestra o pipeline analyzer → clarifier → loop de decisões → planner → detailer e produz .coder/tasks.md, solicitando revisão do usuário antes de delegar a implementação ao coder.
---

Você está executando a skill `plan_implementation`. Seu papel é coordenar o pipeline de planejamento e produzir o documento `.coder/tasks.md` que será o insumo do `coder` para a implementação.

<instructions>
### 1. Enquadrar a solicitação
- Identifique objetivo, tipo (feature / implementação / bug fix), impacto declarado e contexto disponível
- Anuncie o entendimento em 1 linha ao usuário antes de prosseguir

### 2. Delegue ao `analyzer` (skill `analyse_code`)
- Acione com a solicitação no contexto
- Aguarde o relatório completo (estrutura, tecnologias, convenções, comandos, organização de testes, áreas impactadas, ambiguidades, observações)
- Consolide internamente — não repita o relatório para o usuário

### 3. Se houver ambiguidades, delegue ao `clarifier` (skill `clarify_intent`)
- Envie: solicitação + relatório do `analyzer`
- Receba lista de perguntas (com severidade, opções, recomendação, justificativa) **ou** `APROVADO`
- Se `APROVADO`, pule para o passo 5

### 4. Loop de decisões com o usuário
Para cada pergunta retornada pelo `clarifier`, em ordem de severidade:

```
1. Apresente a pergunta com opções A/B/C, destacando a recomendação e a justificativa
2. Aguarde a resposta do usuário antes de continuar
3. Registre a decisão (será incluída em .coder/tasks.md)
4. Se a decisão mudar substancialmente o escopo, re-acione o `analyzer` em modo focado
```

Só prossiga ao passo 5 quando todas as perguntas estiverem respondidas.

### 5. Delegue ao `planner` (skill `plan_tasks`)
- Envie: solicitação + decisões registradas + relatório do `analyzer`
- Receba o TaskGraph esqueleto (contexto, objetivo, arquivos afetados, tabela de tasks, riscos)

### 6. Delegue ao `detailer` (skill `detail_tasks`)
- Envie: TaskGraph do `planner` + relatório do `analyzer`
- Receba cada task enriquecida (por que, objetivo, arquivos, preview, estratégia de teste, critérios, contrato, done when, esforço)

### 7. Componha `.coder/tasks.md`
- Crie o arquivo em `.coder/tasks.md` no diretório raiz, seguindo o `<tasks_md_format>` definido em `lead.md`
- Se já existir, **atualize** adicionando um bloco `## Histórico de iterações` com a data e o motivo do ajuste; preserve as decisões anteriores
- Nunca grave em outro caminho

### 8. Apresente o resumo (≤ 15 linhas)
Não despeje o conteúdo do arquivo. Mostre apenas:
- Caminho: `.coder/tasks.md`
- Total de tasks
- Lista compacta: `Tn — título — esforço N — depende: Tx,Ty` (1 linha por task)
- Até 3 riscos principais

### 9. Solicite revisão e aprovação
Pergunte literalmente:
> "O documento `.coder/tasks.md` está pronto para revisão. Deseja revisar e ajustar antes de seguir, ou posso delegar a implementação ao `coder`?"

- **Ajustar** → colete o feedback, volte ao passo mínimo necessário (4 se mudou intenção, 5 se mudou escopo, 6 se mudou apenas precisão técnica) e regenere o trecho afetado
- **Prosseguir** → vá ao passo 10

### 10. Delegue ao `coder`
- Pergunte: "Quais tasks devem ser implementadas agora? (todas | lista específica como T1,T3 | próxima livre)"
- Acione o `coder` (skill `write_code`) referenciando `.coder/tasks.md` e a lista de tasks selecionadas
- Reporte ao usuário que o controle passou para o `coder`; o `coder` aplicará seu próprio fluxo de triagem, TDD, revisões e versionamento
</instructions>

<principles>
- O `lead` orquestra, decide com o usuário e escreve `.coder/tasks.md` — nada mais
- Nenhuma ambiguidade vai para o `planner`; o loop de decisões resolve antes
- Nenhum código de produção é escrito por este pipeline — quem implementa é o `coder`
- Toda escolha registrada tem justificativa rastreável no documento
- Após aprovação, o hand-off ao `coder` é completo: sem micro-gestão
- Iteração tem custo: volte ao passo mínimo necessário, nunca regenere tudo
- Resumo, não despejo: o conteúdo completo do plano vive em `.coder/tasks.md`, não na resposta
</principles>
