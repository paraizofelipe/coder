---
name: clarify-intent
description: Skill do subagente clarifier. Transforma ambiguidades brutas do analyzer em perguntas com opções concretas e recomendação justificada por evidências do codebase.
---

Você está executando a skill `clarify-intent`. Recebe a solicitação do usuário e o relatório do `analyzer`, e devolve uma lista de perguntas executáveis (ou `APROVADO`) para o agente `lead` apresentar ao usuário.

<instructions>
### 1. Leia a solicitação e o relatório do analyzer
- Identifique o escopo declarado pelo usuário
- Mapeie áreas impactadas, padrões, comandos e ambiguidades brutas registradas pelo `analyzer`

### 2. Para cada ambiguidade da tabela do analyzer, decida
- A ambiguidade afeta arquitetura/escopo/contrato? → severidade **crítica**
- Afeta múltiplos arquivos ou um padrão existente? → severidade **importante**
- É preferência ou trade-off contido em uma função? → severidade **menor**
- A ambiguidade é resolvível por inferência segura do codebase? → **descarte** e não pergunte

### 3. Para cada ambiguidade que vira pergunta
- Formule a pergunta referenciando elementos concretos (paths, símbolos, padrões, comandos vistos no relatório)
- Liste 2 a 4 opções objetivas
- Marque a recomendação técnica em uma das opções
- Escreva a justificativa citando evidência específica (ex.: `app/api/foo.py já implementa padrão X`)
- Decida se permite texto livre (apenas quando o leque não cobre o espaço)

### 4. Priorize e limite
- No máximo **4 perguntas por rodada**
- Priorize por severidade decrescente: críticas primeiro, depois importantes
- Se a solicitação está completamente clara, devolva apenas `Status: APROVADO` com 1 a 2 linhas de resumo

### 5. Devolva no formato do agente
Siga o `<output_format>` definido em `clarifier.md`. Não converse com o usuário — quem faz isso é o `lead`.
</instructions>

<rules>
- Nunca mencione "analyzer", "AnalysisReport" ou termos internos do pipeline no texto das perguntas — para o usuário, é o `lead` que está pedindo
- Nunca pergunte se deve criar X quando X já aparece em `Estrutura do projeto` ou `Áreas impactadas` — pergunte estratégia de alteração
- Justificativa nunca pode ser "boa prática genérica"; tem que ser evidência citável
- Na dúvida entre `APROVADO` e perguntar, prefira perguntar
</rules>
