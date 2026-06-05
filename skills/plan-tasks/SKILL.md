---
name: plan-tasks
description: Skill do subagente planner. Transforma intenção esclarecida + relatório do analyzer em um TaskGraph executável com tasks dimensionadas, dependências reais e riscos.
---

Você está executando a skill `plan_tasks`. Recebe a solicitação esclarecida (com decisões já tomadas pelo usuário via `lead`), o relatório do `analyzer` e devolve o esqueleto do plano técnico no formato `<output_format>` do agente `planner`.

<instructions>
### 1. Consolide contexto e objetivo
- Releia a solicitação + decisões registradas
- Use `Estrutura do projeto`, `Áreas impactadas` e `Comandos disponíveis` do `analyzer` como **base canônica** — não revalide nem re-inspecione esses itens

### 2. Liste arquivos afetados
- Comece pelo conjunto da seção `Áreas impactadas` do `analyzer` — esses paths são aceitos sem revalidação
- Só adicione paths **novos**, fora do relatório do `analyzer`, quando houver lacuna concreta identificada durante o planejamento; nesses casos sim, confirme via LSP/grep/glob antes de incluir
- Mantenha a marcação `[parcial]` quando o analyzer marcou assim

### 3. Quebre em tasks
- Uma mudança lógica por task (single responsibility)
- Cada task deve ser revisável em poucas horas e testável de forma isolada
- Identifique dependências reais (assinatura compartilhada, schema, contrato) — preencha `Depende de` apenas nesses casos
- Tasks que tocam arquivos diferentes sem acoplamento ficam paralelas (sem dependência)
- **Não crie tasks separadas para testes** — testes ficam embutidos na task que entrega o código que valida

### 4. Sizing check
- Se uma task parece > 3 dias úteis OU toca > 5 arquivos não-correlatos, **quebre antes de devolver**
- Se uma task ficou vaga, registre em `Riscos` para o `lead` decidir se precisa nova rodada de clarificação

### 5. Levante riscos
- Pontos que podem escalar escopo durante implementação
- Convenções do codebase que podem entrar em conflito com a abordagem proposta
- Áreas onde o `analyzer` marcou `[parcial]` e podem revelar surpresas

### 6. Devolva no formato do agente
Siga o `<output_format>` definido em `planner.md`. Não detalhe code preview, estratégia de teste ou contrato de interface — isso é responsabilidade do `detailer`.
</instructions>

<rules>
- Não re-inspecione áreas já cobertas pelo `analyzer`
- `Depende de` reflete bloqueio real, não preferência de ordem
- Nunca crie task "implementar testes" — testes pertencem à task que entrega o código
- Critérios de aceite na coluna `Descrição` devem ser observáveis (curl/CLI/pytest), não subjetivos
</rules>
