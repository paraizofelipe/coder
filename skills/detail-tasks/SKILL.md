---
name: detail-tasks
description: Skill do subagente detailer. Enriquece cada task esqueleto do planner com motivação, arquivos, preview de código, estratégia de teste, critérios de aceite, contrato de interface, definition of done e esforço estimado.
---

Você está executando a skill `detail_tasks`. Recebe o `TaskGraph` esqueleto do `planner`, o relatório do `analyzer` e devolve cada task enriquecida no formato `<output_format>` do agente `detailer`.

<instructions>
### 1. Para cada task do TaskGraph, em ordem
- Releia título + descrição + dependências
- Identifique vizinhos relevantes (tasks que esta depende ou que dependem dela)
- Verifique se algum arquivo foi tocado por task anterior; se sim, não o repita em `arquivos_afetados` desta task

### 2. Drille apenas o necessário
- Use `lsp` / grep / glob para confirmar assinaturas reais antes de escrever `preview_de_codigo` ou `contrato_de_interface`
- Não relista diretórios desconhecidos — confie no `analyzer`

### 3. Preencha os campos obrigatórios da task
- `por_que`: 1 a 2 linhas ancoradas no analyzer ou na decisão registrada
- `objetivo`: 1 frase com resultado observável
- `arquivos_afetados`: paths exatos, mantendo `[parcial]` quando aplicável
- `estrategia_de_teste`: arquivos de teste a criar/atualizar, cenários positivo/negativo/borda, comando exato (use o `Comandos disponíveis` do analyzer)
- `criterios_de_aceite`: 2 a 5 observáveis concretos
- `done_when`: 2 a 4 itens de higiene (lint/types/format/test) — use os comandos do analyzer
- `esforco_estimado`: 1 a 5 dias úteis (>3 só multi-camada)

### 4. Preencha os campos contextuais quando aplicar
- `preview_de_codigo`: apenas trecho mínimo relevante, com `tipo` em `nova_funcao | novo_arquivo | modificacao | referencia`
- `contrato_de_interface`: sempre que esta task expõe algo consumido por outra; caso contrário, `n/a`
- `arquivos_proibidos`: apenas quando há risco real de mexer em arquivo vizinho que não pertence à task

### 5. Sizing check antes de devolver
- Se `arquivos_afetados` tem 5+ entradas não correlatas, sinalize no `por_que` ou no preview que o recorte do `planner` precisa de revisão
- Se `criterios_de_aceite` contém algo subjetivo ("código limpo", "boa performance"), substitua por observável concreto

### 6. Devolva no formato do agente
Siga o `<output_format>` definido em `detailer.md` — uma task por bloco, separadas por `---`.
</instructions>

<rules>
- Nunca invente assinaturas — sempre confirmar via LSP/grep/glob
- `preview_de_codigo` é trecho mínimo, nunca arquivo inteiro
- Toda task tem testes embutidos; nunca referencie uma "task de testes" separada
- `criterios_de_aceite` só observáveis; `done_when` só higiene
- Não duplique arquivos entre tasks
</rules>
