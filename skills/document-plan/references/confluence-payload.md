# Payload e chamadas MCP — Confluence

Referência detalhada das chamadas ao MCP `atlassian_local` usadas pela skill `document-plan`. Carregue este arquivo apenas quando precisar montar uma chamada concreta.

## Localizar a página pai "Implementações"

```
atlassian_local_confluence_search(
  query: "title = \"Implementações\" AND space = \"CAT\"",
  limit: 1
)
```

- Se retornar resultado: extrair o `id` da página — este é o `parent_id`
- Se não retornar: criar a página antes de continuar (próximo bloco)

## Criar a página raiz "Implementações" (quando ausente)

```
atlassian_local_confluence_create_page(
  space_key: "CAT",
  title: "Implementações",
  content: "Página raiz para registro de planos de implementação.",
  representation: "wiki"
)
```

- Extrair o `id` da página criada — este é o `parent_id`
- Se a criação falhar: reportar o erro exato retornado pelo MCP e encerrar

## Verificar se a subpágina já existe

```
atlassian_local_confluence_search(
  query: "title = \"<título extraído>\" AND parent = \"<parent_id>\" AND space = \"CAT\"",
  limit: 1
)
```

- Se retornar resultado: registrar o `id` da página existente → ir para atualização
- Se não retornar: ir para criação

## Criar a subpágina (quando não existe)

```
atlassian_local_confluence_create_page(
  space_key: "CAT",
  title: "<título extraído>",
  parent_id: "<id da página Implementações>",
  content: "<conteúdo completo do plan.md>",
  representation: "wiki"
)
```

- Registrar a URL retornada pelo MCP

## Comparar e atualizar a subpágina (quando já existe)

Obter a versão atual:

```
atlassian_local_confluence_get_page(
  page_id: "<id da página existente>"
)
```

Extrair:
- `version.number` — número da versão atual
- `body.wiki.value` (ou equivalente na representação retornada) — conteúdo atual no Confluence

**Comparar conteúdos** antes de atualizar:

- Normalizar ambos: remover espaços em branco no início e fim de cada linha, ignorar linhas em branco consecutivas
- Se idêntico após normalização: encerrar sem atualizar e reportar "sem alterações"
- Se houver diferenças: prosseguir com `update_page`

```
atlassian_local_confluence_update_page(
  page_id: "<id da página existente>",
  title: "<título extraído>",
  content: "<conteúdo completo do plan.md>",
  representation: "wiki",
  version: <version.number + 1>
)
```

- Registrar a URL atualizada

## Formato do conteúdo

O Confluence aceita Markdown via `representation: "wiki"` (mais compatível com Markdown puro) ou `"storage"` (XHTML). Use sempre `"wiki"` ao criar/atualizar nesta skill.

## Tratamento de erros

Se qualquer chamada ao MCP falhar, reportar o erro exato retornado e encerrar sem tentar alternativas — não inventar resultados.
