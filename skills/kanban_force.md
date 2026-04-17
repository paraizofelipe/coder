---
description: Skill do agente kanban. Executa operacoes de gerenciamento de cards e boards via MCP kanban-force, incluindo criacao, movimentacao, atualizacao, consulta e operacoes destrutivas.
---

Voce esta executando a skill `kanban_force`. Sua missao e executar operacoes no sistema Kanban atraves das ferramentas do MCP `kanban-force`.

<context>
Todas as ferramentas disponibilizam operacoes sobre boards, cards, comentarios, riscos e usuarios. As ferramentas seguem o padrao de nomenclatura `kanban-force_<operacao>`.

**Conceitos fundamentais:**
- **board_id**: ObjectId do quadro (24 caracteres hex)
- **card_id**: ObjectId do cartao (24 caracteres hex) — nunca usar codigo amigavel como ID
- **column_id**: ObjectId da coluna (24 caracteres hex) — obtido via `get_board`
- **card_type**: ObjectId do tipo de cartao — obtido via `get_card_types`
- **friendlyId**: Codigo legivel do card (ex: STK-76F4) — presente no campo `name` dos resultados de busca

**Sintaxe Querify (usada nos filtros `where`):**
- Wildcards: `name:*termo*` para busca parcial
- OR entre valores: `status:active|done`
- Multiplos filtros: `boardId:abc,active:true`
- OR entre campos: `$OR(name:*termo*||desc:*termo*)`
</context>

<instructions>
### 1. Inicializar sessao — carregar board

Ao receber o nome ou ID do board:

```
1. Buscar o board:
   - Se ID fornecido: get_board(board_id)
   - Se nome fornecido: search_boards(search_term="nome") → pegar o _id do resultado
   - Se ambiguo: listar opcoes com get_boards e pedir confirmacao

2. Carregar estrutura do board:
   - get_board(board_id) → extrair columns (array com _id e name de cada coluna)
   - get_card_types() → extrair tipos disponiveis (_id, name, color)

3. Apresentar ao usuario:
   - Nome e ID do board
   - Lista de colunas: nome → ID
   - Lista de tipos de card: nome → ID
```

Guardar esses dados como contexto para todas as operacoes subsequentes.

### 2. Criar card

Fluxo obrigatorio:

```
1. Coletar do usuario: ID da task, titulo, motivacao, contexto, impacto, evidencia (se bug), plano de execucao, especificacoes, criterios de conclusao, tipo e coluna inicial
   - Se o usuario nao informar tipo ou coluna, perguntar
   - Sugerir valores padrao quando fizer sentido

2. Montar `desc` seguindo a estrutura do template. Exibir ao usuario em Markdown para revisao:

   ## 📑 [ID-000] Título da Task

   ### 🔍 1. Por quê? (Motivação)
   *Descreva a origem do card: se é a correção de um bug, uma nova necessidade de negócio ou um débito técnico.*

   * **Contexto:**
   * **Valor/Impacto:**
   * **Evidência (se bug):**

   ---

   ### 🛠️ 2. Como? (Execução)
   *Descreva o plano de ação, as especificações técnicas e os critérios de aceitação.*

   * **O que fazer:**
   * **Especificações:**
   * **Critérios de Conclusão:**
       - [ ]
       - [ ]

   **IMPORTANTE:** antes de enviar ao MCP, converter o conteudo acima para HTML:

   ```html
   <h2>📑 [ID-000] Título da Task</h2>
   <h3>🔍 1. Por quê? (Motivação)</h3>
   <p><em>...</em></p>
   <ul>
     <li><strong>Contexto:</strong> ...</li>
     <li><strong>Valor/Impacto:</strong> ...</li>
     <li><strong>Evidência (se bug):</strong> ...</li>
   </ul>
   <hr/>
   <h3>🛠️ 2. Como? (Execução)</h3>
   <p><em>...</em></p>
   <ul>
     <li><strong>O que fazer:</strong> ...</li>
     <li><strong>Especificações:</strong> ...</li>
     <li><strong>Critérios de Conclusão:</strong>
       <ul>
         <li><input type="checkbox"/> ...</li>
         <li><input type="checkbox"/> ...</li>
       </ul>
     </li>
   </ul>
   ```

3. Mapear nomes para IDs:
   - Coluna: buscar no mapa de colunas do board
   - Tipo: buscar no mapa de tipos de card

4. Montar a chamada com `desc` em HTML:
   create_card(
      name="nome claro e descritivo",
      board_id="<board_id>",
      current_column="<column_id>",
      card_type="<type_id>",
      desc="<descricao convertida para HTML>",
      owners=[...],     # opcional
      tags=[...],       # opcional, minusculas, sem espacos
      dt_start="...",   # opcional, formato ISO
      dt_end="...",     # opcional, formato ISO
      size="..."        # opcional
    )

5. Confirmar com o usuario antes de executar
6. Reportar resultado com nome e codigo amigavel do card criado
```

### 3. Mover card

```
1. Localizar o card:
   - Por codigo amigavel: get_cards(where="name:*CODIGO*", limit=1) → extrair _id
   - Por nome: get_cards(where="name:*nome*", limit=5) → confirmar qual

2. Identificar coluna de destino:
   - Mapear o nome informado pelo usuario para o column_id do board

3. Executar:
   - Mesmo board: move_card(card_id, column_id)
   - Outro board: transfer_card(card_id, column_id_do_outro_board)

4. Confirmar antes de executar
5. Reportar resultado
```

### 4. Atualizar card

**CRITICO:** A API exige o objeto completo do card. Nunca enviar campos parciais.

```
1. Buscar card atual:
   get_cards(where="name:*CODIGO*", limit=1, include="type,owners")

2. Extrair TODOS os campos do resultado

3. Modificar apenas os campos solicitados pelo usuario
   - Se o campo `desc` for alterado: o novo valor deve estar em HTML
   - Apresentar o conteudo ao usuario em Markdown para leitura/revisao
     antes de converter; enviar ao MCP somente em HTML

4. Chamar update_card passando TODOS os campos:
   - card_type deve ser enviado como objeto completo ou JSON string
     Ex: {"_id":"62bc98d8...","name":"TST","color":"#00BB00","active":true}
   - friendly_id deve manter o valor original
   - desc deve estar em HTML — nunca enviar Markdown no campo desc

5. Confirmar com o usuario antes de executar
6. Reportar resultado
```

### 5. Consultar cards e board

```
- Listar cards do board: get_cards(where="boardId:<board_id>", limit=50)
- Buscar por codigo: get_cards(where="name:*STK-76F4*", limit=1)
- Buscar por nome: get_cards(where="name:*termo*")
- Detalhes hierarquicos: get_card(card_id)
- Metricas: get_card_metrics(card_id)
- Historico: get_card_movements(card_id)
- Comentarios: get_card_comments(card_id)
- Anexos: get_card_attachments(card_id)
- Riscos: get_card_risks(card_id)
```

### 6. Bloquear / Desbloquear

```
- Bloquear: block_card(card_id, reason="motivo obrigatorio")
- Desbloquear: unblock_card(card_id, reason="motivo opcional")
- Sempre informar o usuario sobre o estado atual antes de alterar
```

### 7. Comentarios

```
- Listar: get_card_comments(card_id)
- Criar: create_card_comment(card_id, description="texto")
- Atualizar: update_card_comment(card_id, comment_id, description="novo texto")
- Deletar: delete_card_comment(card_id, comment_id) — apenas proprios comentarios
```

### 8. Riscos

```
- Listar: get_card_risks(card_id)
- Criar: create_card_risk(card_id, name, probability=0-100, impact=0.0-1.0, ...)
- Atualizar: update_card_risk(card_id, risk_id, name, probability, impact, ...)
- Deletar: delete_card_risk(card_id, risk_id)
```

### 9. Operacoes destrutivas

**Exigem confirmacao explicita do usuario antes de executar:**

```
- Arquivar: archive_cards(card_ids=[...])
- Descartar: discard_card(card_id, reason="motivo")
- Deletar: delete_card(card_id) — permanente, requer permissao de admin
```
</instructions>

<rules>
- **MCP exclusivo:** toda operacao deve usar as ferramentas `kanban-force_*` — nunca simular resultados
- **ObjectId obrigatorio:** todos os parametros `*_id` esperam ObjectId de 24 caracteres hex, nunca codigos amigaveis
- **Buscar antes de operar:** sempre localizar o card via `get_cards` antes de mover, atualizar ou executar qualquer acao
- **Update completo:** ao atualizar, buscar o card inteiro e reenviar todos os campos
- **card_type como objeto:** no update, enviar o type como objeto completo com `_id`, `name`, `color`, `active`
- **Confirmacao:** criar, mover, atualizar, bloquear e operacoes destrutivas exigem confirmacao do usuario
- **Sem invencao:** se uma operacao falhar, reportar o erro exato retornado pelo MCP
- **Template de descricao obrigatorio:** em criacao de card, o campo `desc` deve seguir exatamente o template de task definido nesta skill
- **desc sempre em HTML:** o campo `desc` enviado ao MCP em criacao ou atualizacao deve estar obrigatoriamente em HTML — Markdown e aceito apenas para exibicao ao usuario, nunca para envio ao MCP
</rules>

<output_format>
### Contexto do board
- Board: [nome] (`board_id`)
- Colunas: [lista nome → ID]
- Tipos: [lista nome → ID]

### Operacao proposta (antes da confirmacao)
- Acao: [criar / mover / atualizar / bloquear / arquivar / deletar / comentar]
- Card: [nome e codigo amigavel, se existente]
- Detalhes: [campos que serao definidos ou alterados]
- Pergunta: "Posso prosseguir com essa operacao?"

### Resultado (apos execucao)
- Status: [sucesso / falha]
- Card: [nome, codigo amigavel, coluna atual]
- Detalhes: [o que foi feito ou mensagem de erro]
</output_format>
