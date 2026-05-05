---
name: kanban-force
description: Skill do agente kanban. Cria, move, atualiza, consulta, bloqueia, comenta e arquiva cards e boards via MCP kanban-force. Use quando o usuário mencionar cards, boards, IDs amigáveis (ex.: STK-76F4, UST-FF51), movimentação entre colunas, comentários, bloqueios, riscos, métricas Kanban ou operações destrutivas (arquivar/descartar/deletar). Operações destrutivas e mutações exigem confirmação explícita do usuário. Toda chamada ao MCP usa ObjectId (24 hex), nunca friendlyId.
---

Você está executando a skill `kanban-force`. Sua missão é executar operações no sistema Kanban através das ferramentas do MCP `kanban-force`.

<context>
Todas as ferramentas do MCP `kanban-force` operam sobre boards, cards, comentários, riscos e usuários. Seguem o padrão de nomenclatura `kanban-force_<operacao>`.

**Conceitos fundamentais e sintaxe Querify dos filtros `where`:** ver `references/querify-syntax.md`.

**Template obrigatório do campo `desc`:** ver `references/description-template.md`. Markdown apenas para exibição ao usuário; HTML obrigatório no envio ao MCP.
</context>

<instructions>
### 1. Inicializar sessão — carregar board

Ao receber o nome ou ID do board:

```
1. Buscar o board:
   - Se ID fornecido: get_board(board_id)
   - Se nome fornecido: search_boards(search_term="nome") → pegar o _id do resultado
   - Se ambíguo: listar opções com get_boards e pedir confirmação

2. Carregar estrutura do board:
   - get_board(board_id) → extrair columns (array com _id e name de cada coluna)
   - get_card_types() → extrair tipos disponíveis (_id, name, color)

3. Apresentar ao usuário:
   - Nome e ID do board
   - Lista de colunas: nome → ID
   - Lista de tipos de card: nome → ID
```

Guardar esses dados como contexto para todas as operações subsequentes.

### 2. Criar card

```
1. Coletar do usuário: ID da task, título, motivação, contexto, impacto, evidência (se bug),
   plano de execução, especificações, critérios de conclusão, tipo e coluna inicial
   - Se o usuário não informar tipo ou coluna, perguntar
   - Sugerir valores padrão quando fizer sentido

2. Montar `desc` seguindo o template em `references/description-template.md`:
   - Apresentar ao usuário em Markdown para revisão
   - Converter para HTML antes de enviar ao MCP

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
      tags=[...],       # opcional, minúsculas, sem espaços
      dt_start="...",   # opcional, formato ISO
      dt_end="...",     # opcional, formato ISO
      size="..."        # opcional
    )

5. Confirmar com o usuário antes de executar
6. Reportar resultado com nome e código amigável do card criado
```

### 3. Mover card

```
1. Localizar o card:
   - Por código amigável: get_cards(where="name:*CODIGO*", limit=1) → extrair _id
   - Por nome: get_cards(where="name:*nome*", limit=5) → confirmar qual

2. Identificar coluna de destino:
   - Mapear o nome informado pelo usuário para o column_id do board

3. Executar:
   - Mesmo board: move_card(card_id, column_id)
   - Outro board: transfer_card(card_id, column_id_do_outro_board)

4. Confirmar antes de executar
5. Reportar resultado
```

### 4. Atualizar card

**CRÍTICO:** A API exige o objeto completo do card. Nunca enviar campos parciais.

```
1. Buscar card atual:
   get_cards(where="name:*CODIGO*", limit=1, include="type,owners")

2. Extrair TODOS os campos do resultado

3. Modificar apenas os campos solicitados pelo usuário
   - Se o campo `desc` for alterado: o novo valor deve estar em HTML
   - Apresentar o conteúdo ao usuário em Markdown para leitura/revisão antes de converter
   - Enviar ao MCP somente em HTML

4. Chamar update_card passando TODOS os campos:
   - card_type deve ser enviado como objeto completo ou JSON string
     Ex.: {"_id":"62bc98d8...","name":"TST","color":"#00BB00","active":true}
   - friendly_id deve manter o valor original
   - desc deve estar em HTML — nunca enviar Markdown no campo desc

5. Confirmar com o usuário antes de executar
6. Reportar resultado
```

### 5. Consultar cards e board

Ver `references/querify-syntax.md` para a sintaxe completa dos filtros.

```
- Listar cards do board: get_cards(where="boardId:<board_id>", limit=50)
- Buscar por código: get_cards(where="name:*STK-76F4*", limit=1)
- Buscar por nome: get_cards(where="name:*termo*")
- Detalhes hierárquicos: get_card(card_id)
- Métricas: get_card_metrics(card_id)
- Histórico: get_card_movements(card_id)
- Comentários: get_card_comments(card_id)
- Anexos: get_card_attachments(card_id)
- Riscos: get_card_risks(card_id)
```

### 6. Bloquear / Desbloquear

```
- Bloquear: block_card(card_id, reason="motivo obrigatório")
- Desbloquear: unblock_card(card_id, reason="motivo opcional")
- Sempre informar o usuário sobre o estado atual antes de alterar
```

### 7. Comentários

```
- Listar: get_card_comments(card_id)
- Criar: create_card_comment(card_id, description="texto")
- Atualizar: update_card_comment(card_id, comment_id, description="novo texto")
- Deletar: delete_card_comment(card_id, comment_id) — apenas próprios comentários
```

### 8. Riscos

```
- Listar: get_card_risks(card_id)
- Criar: create_card_risk(card_id, name, probability=0-100, impact=0.0-1.0, ...)
- Atualizar: update_card_risk(card_id, risk_id, name, probability, impact, ...)
- Deletar: delete_card_risk(card_id, risk_id)
```

### 9. Operações destrutivas

**Exigem confirmação explícita do usuário antes de executar:**

```
- Arquivar: archive_cards(card_ids=[...])
- Descartar: discard_card(card_id, reason="motivo")
- Deletar: delete_card(card_id) — permanente, requer permissão de admin
```
</instructions>

<rules>
- **MCP exclusivo:** toda operação deve usar as ferramentas `kanban-force_*` — nunca simular resultados
- **ObjectId obrigatório:** todos os parâmetros `*_id` esperam ObjectId de 24 caracteres hex, nunca códigos amigáveis
- **Buscar antes de operar:** sempre localizar o card via `get_cards` antes de mover, atualizar ou executar qualquer ação
- **Update completo:** ao atualizar, buscar o card inteiro e reenviar todos os campos
- **card_type como objeto:** no update, enviar o type como objeto completo com `_id`, `name`, `color`, `active`
- **Confirmação:** criar, mover, atualizar, bloquear e operações destrutivas exigem confirmação do usuário
- **Sem invenção:** se uma operação falhar, reportar o erro exato retornado pelo MCP
- **Template de descrição obrigatório:** em criação de card, o campo `desc` deve seguir exatamente o template em `references/description-template.md`
- **desc sempre em HTML:** o campo `desc` enviado ao MCP em criação ou atualização deve estar obrigatoriamente em HTML — Markdown é aceito apenas para exibição ao usuário, nunca para envio ao MCP
</rules>

<output_format>
### Contexto do board
- Board: [nome] (`board_id`)
- Colunas: [lista nome → ID]
- Tipos: [lista nome → ID]

### Operação proposta (antes da confirmação)
- Ação: [criar / mover / atualizar / bloquear / arquivar / deletar / comentar]
- Card: [nome e código amigável, se existente]
- Detalhes: [campos que serão definidos ou alterados]
- Pergunta: "Posso prosseguir com essa operação?"

### Resultado (após execução)
- Status: [sucesso / falha]
- Card: [nome, código amigável, coluna atual]
- Detalhes: [o que foi feito ou mensagem de erro]
</output_format>
