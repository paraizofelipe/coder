---
name: document-plan
description: Skill do agente documenter. Lê .coder/plan.md, extrai o título da implementação e publica uma subpágina no Confluence (space CAT, raiz Implementações) via MCP atlassian_local. Use quando o usuário pedir para publicar, sincronizar ou atualizar um plano no Confluence. Idempotente — não atualiza se conteúdo for idêntico após normalização. Cria automaticamente a página raiz "Implementações" se ausente.
---

Você está executando a skill `document-plan`. Sua missão é publicar o conteúdo de `.coder/plan.md` no Confluence como subpágina de `Implementações` no space `CAT`.

<context>
Todas as operações usam as ferramentas do MCP `atlassian_local`. As ferramentas seguem o padrão de nomenclatura `atlassian_local_<operacao>`.

**Parâmetros fixos desta skill:**
- **Space:** `CAT`
- **Página raiz:** `Implementações`
- **Fonte do conteúdo:** `.coder/plan.md` no diretório raiz do projeto

Para o payload exato e os parâmetros de cada chamada MCP (`confluence_search`, `confluence_create_page`, `confluence_get_page`, `confluence_update_page`), consulte `references/confluence-payload.md`.
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

- Buscar a página raiz no space `CAT` (ver `references/confluence-payload.md` → "Localizar a página pai")
- Se existir: extrair o `id` — este é o `parent_id`
- Se não existir: criar a página antes de continuar (ver "Criar a página raiz" no reference)

### 4. Verificar se a subpágina já existe

Buscar uma subpágina com o mesmo título sob o `parent_id` (ver "Verificar se a subpágina já existe" no reference).

- Se existir: ir para o passo 6 (atualização)
- Se não existir: ir para o passo 5 (criação)

### 5. Criar a subpágina (quando não existe)

Executar `confluence_create_page` com `space_key`, `title`, `parent_id` e `content` do `plan.md` em `representation: "wiki"`. Registrar a URL retornada e ir para o passo 7.

### 6. Comparar e atualizar a subpágina (quando já existe)

1. Obter a versão atual da página com `confluence_get_page`
2. Extrair `version.number` e `body.wiki.value`
3. **Comparar o conteúdo atual com o `plan.md`:**
   - Normalizar ambos: remover espaços em branco no início e fim de cada linha, ignorar linhas em branco consecutivas
   - Se idêntico após normalização: informar "sem alterações" e **encerrar sem chamar update_page**
   - Se houver diferenças: chamar `confluence_update_page` com `version: <version.number + 1>` e o conteúdo do `plan.md`
4. Registrar a URL atualizada

### 7. Reportar o resultado

```
## Publicação concluída

**Operação:** criação / atualização / sem alterações
**Título:** <título extraído>
**Local:** Confluence › CAT › Implementações › <título>
**URL:** <link direto para a página>          ← omitir se operação for "sem alterações"
```

Se qualquer chamada ao MCP retornar erro, reportar o erro exato e encerrar sem tentar alternativas.
</instructions>

<rules>
- **MCP exclusivo:** toda operação deve usar as ferramentas `atlassian_local_*` — nunca simular resultados
- **Space e hierarquia fixos:** sempre `CAT` / `Implementações` — nunca publicar em outro local
- **Sem modificar o plan.md:** apenas ler o arquivo, nunca escrever nele
- **Idempotência:** verificar existência antes de criar — se já existe, comparar conteúdo antes de atualizar
- **Sem publicação desnecessária:** se o conteúdo do Confluence for idêntico ao `plan.md` (após normalização), encerrar sem chamar `update_page`
- **Sem truncar conteúdo:** publicar o plan.md completo, sem omitir seções
- **Sem inventar:** se o MCP falhar, reportar o erro exato — nunca fingir sucesso
</rules>

<output_format>
## Publicação concluída

- **Operação:** [criação / atualização / sem alterações]
- **Título:** [título extraído do plan.md]
- **Local:** Confluence › CAT › Implementações › [título]
- **URL:** [link direto para a página]
</output_format>
