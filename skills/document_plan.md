---
description: Skill do agente documenter. Lê .coder/plan.md, extrai o título da implementação e publica uma subpágina no Confluence (space CAT, raiz Implementações) via MCP atlassian_local.
---

Você está executando a skill `document_plan`. Sua missão é publicar o conteúdo de `.coder/plan.md` no Confluence como subpágina de `Implementações` no space `CAT`.

<context>
Todas as operações usam as ferramentas do MCP `atlassian_local`. As ferramentas seguem o padrão de nomenclatura `atlassian_local_<operacao>`.

**Parâmetros fixos desta skill:**
- **Space:** `CAT`
- **Página raiz:** `Implementações`
- **Fonte do conteúdo:** `.coder/plan.md` no diretório raiz do projeto

**Formato do conteúdo:** O Confluence aceita Markdown via `representation: "wiki"` ou `"storage"` (XHTML). Use `"wiki"` ao criar/atualizar, pois é mais compatível com conteúdo Markdown puro.
</context>

<instructions>

### 1. Ler o plano

Ler o arquivo `.coder/plan.md` do diretório raiz do projeto.

- Se o arquivo não existir: informar ao usuário que nenhum plano foi encontrado e encerrar
- Se o arquivo estiver vazio: informar ao usuário e encerrar

### 2. Extrair o título

A partir do conteúdo do `plan.md`, determinar o título da subpágina no Confluence:

1. Procurar o primeiro cabeçalho de nível 1 (`# Título`) ou nível 2 (`## Título`) que descreva a implementação
2. Se não houver cabeçalho claro, usar a primeira linha não vazia do arquivo como título
3. Limpar o título: remover prefixo `#`, espaços extras e caracteres especiais que não sejam aceitos em títulos do Confluence

O título resultante será o nome da subpágina no Confluence.

### 3. Localizar ou criar a página pai "Implementações"

Buscar a página raiz no space `CAT`:

```
atlassian_local_confluence_search(
  query: "title = \"Implementações\" AND space = \"CAT\"",
  limit: 1
)
```

- Se retornar resultado: extrair o `id` da página — este é o `parent_id`
- Se não retornar resultado: a página não existe e deve ser criada antes de continuar:

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

### 4. Verificar se a subpágina já existe

Buscar uma subpágina com o mesmo título sob o `parent_id` obtido:

```
atlassian_local_confluence_search(
  query: "title = \"<título extraído>\" AND parent = \"<parent_id>\" AND space = \"CAT\"",
  limit: 1
)
```

- Se retornar resultado: registrar o `id` da página existente → ir para o passo 6 (atualização)
- Se não retornar resultado: ir para o passo 5 (criação)

### 5. Criar a subpágina (quando não existe)

```
atlassian_local_confluence_create_page(
  space_key: "CAT",
  title: "<título extraído>",
  parent_id: "<id da página Implementações>",
  content: "<conteúdo completo do plan.md>",
  representation: "wiki"
)
```

- Registrar a URL da página criada retornada pelo MCP
- Ir para o passo 7 (resultado)

### 6. Atualizar a subpágina (quando já existe)

Antes de atualizar, obter a versão atual da página:

```
atlassian_local_confluence_get_page(
  page_id: "<id da página existente>"
)
```

Extrair o campo `version.number` do resultado.

Em seguida, atualizar:

```
atlassian_local_confluence_update_page(
  page_id: "<id da página existente>",
  title: "<título extraído>",
  content: "<conteúdo completo do plan.md>",
  representation: "wiki",
  version: <version.number + 1>
)
```

- Registrar a URL da página atualizada
- Ir para o passo 7 (resultado)

### 7. Reportar o resultado

Apresentar ao usuário:

```
## Publicação concluída

**Operação:** criação / atualização
**Título:** <título extraído>
**Local:** Confluence › CAT › Implementações › <título>
**URL:** <link direto para a página>
```

Se qualquer chamada ao MCP retornar erro, reportar o erro exato e encerrar sem tentar alternativas.
</instructions>

<rules>
- **MCP exclusivo:** toda operação deve usar as ferramentas `atlassian_local_*` — nunca simular resultados
- **Space e hierarquia fixos:** sempre `CAT` / `Implementações` — nunca publicar em outro local
- **Sem modificar o plan.md:** apenas ler o arquivo, nunca escrever nele
- **Idempotência:** verificar existência antes de criar — se já existe, atualizar com incremento de versão
- **Sem truncar conteúdo:** publicar o plan.md completo, sem omitir seções
- **Sem inventar:** se o MCP falhar, reportar o erro exato — nunca fingir sucesso
</rules>

<output_format>
## Publicação concluída

- **Operação:** [criação / atualização]
- **Título:** [título extraído do plan.md]
- **Local:** Confluence › CAT › Implementações › [título]
- **URL:** [link direto para a página]
</output_format>
