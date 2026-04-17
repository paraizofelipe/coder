---
description: Consulta e carrega no contexto as informações de um card do kanban-force pelo friendlyID. Uso: /kanban_card STK-76F4
---

Utilize as ferramentas do MCP `kanban-force` para carregar no contexto as informações do card com friendlyID `$ARGUMENTS`.

## Passos

### 1. Buscar o card pelo friendlyID

Chamar `get_cards` com o filtro abaixo para localizar o card e excluir arquivados:

```
where: "name:*$ARGUMENTS*,archived:false"
limit: 5
```

- O campo `name` dos resultados contém o friendlyId (ex: `STK-76F4`)
- Confirmar que o resultado corresponde exatamente ao friendlyID `$ARGUMENTS`
- Se nenhum resultado for retornado: informar que o card `$ARGUMENTS` não foi encontrado ou está arquivado e encerrar

### 2. Carregar os detalhes completos

Com o `_id` (ObjectId de 24 caracteres hex) do card encontrado:

- `get_card(card_id)` — detalhes completos com subitens e hierarquia
- `get_card_comments(card_id)` — comentários registrados
- `get_card_movements(card_id)` — histórico de movimentações entre colunas

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
