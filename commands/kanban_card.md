---
description: Consulta e carrega no contexto as informações de um card do kanban-force pelo friendlyID. Uso: /kanban_card STK-76F4
agent: kanban
---

Utilize as ferramentas do MCP `kanban-force` para carregar no contexto as informações do card com friendlyID `$ARGUMENTS`.

## Passos

### 1. Buscar o card pelo friendlyID

Chamar a ferramenta `kanban-force_get_cards` com os parâmetros:

```
where: "name:*$ARGUMENTS*,active:true"
limit: 5
```

- O campo `name` dos resultados contém o friendlyId (ex: `STK-76F4`)
- Confirmar que o resultado corresponde exatamente ao friendlyID `$ARGUMENTS`
- Se nenhum resultado for retornado: informar que o card `$ARGUMENTS` não foi encontrado ou está arquivado e encerrar

### 2. Carregar os detalhes completos

Com o `_id` (ObjectId de 24 caracteres hex) do card encontrado, chamar em sequência:

- `kanban-force_get_card` com `card_id` — detalhes completos com subitens e hierarquia
- `kanban-force_get_card_comments` com `card_id` — comentários registrados
- `kanban-force_get_card_movements` com `card_id` — histórico de movimentações entre colunas

### 3. Apresentar as informações carregadas

Exibir de forma estruturada:

```
## Card: $ARGUMENTS

**Nome:** [nome completo]
**Coluna atual:** [coluna]
**Tipo:** [tipo do card]
**Owners:** [lista]
**Tags:** [lista]
**Início:** [data]  |  **Prazo:** [data]

### Descrição
[conteúdo do campo desc]

### Comentários ([N])
[lista de comentários com autor e data]

### Histórico de movimentações
[lista de movimentações com coluna origem → destino e data]
```

Após apresentar, informar que o card está carregado no contexto e disponível para operações.
