# Sintaxe Querify — kanban-force

Querify é a sintaxe usada nos filtros `where` das ferramentas do MCP `kanban-force` (ex.: `get_cards`, `search_boards`).

## Operadores

- **Wildcards:** `name:*termo*` — busca parcial (substring)
- **OR entre valores:** `status:active|done` — qualquer um dos valores
- **Múltiplos filtros (AND):** `boardId:abc,active:true` — todos os filtros precisam casar
- **OR entre campos:** `$OR(name:*termo*||desc:*termo*)` — termo presente em qualquer um dos campos

## Exemplos práticos

### Buscar card por código amigável
```
get_cards(where="name:*STK-76F4*", limit=1)
```

### Buscar cards de um board específico, ativos
```
get_cards(where="boardId:62bc98d8a1b2c3d4e5f6a7b8,active:true", limit=50)
```

### Buscar cards por nome em vários campos
```
get_cards(where="$OR(name:*pagamento*||desc:*pagamento*)", limit=20)
```

### Buscar cards em múltiplos status
```
get_cards(where="status:active|done|in_progress", limit=100)
```

## Conceitos fundamentais (IDs)

- **board_id:** ObjectId do quadro (24 caracteres hex)
- **card_id:** ObjectId do cartão (24 caracteres hex) — nunca usar código amigável como ID
- **column_id:** ObjectId da coluna (24 caracteres hex) — obtido via `get_board`
- **card_type:** ObjectId do tipo de cartão — obtido via `get_card_types`
- **friendlyId:** código legível do card (ex.: `STK-76F4`) — presente no campo `name` dos resultados de busca

**Regra:** todos os parâmetros `*_id` em chamadas ao MCP esperam ObjectId. Para códigos amigáveis, buscar o ObjectId correspondente via `get_cards` antes.
