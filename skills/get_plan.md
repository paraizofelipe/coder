---
description: Skill do agente documenter. Busca o plano de implementação publicado no Confluence (space CAT, raiz Implementações) e salva em .coder/plan.md no projeto local.
---

Você está executando a skill `get_plan`. Sua missão é buscar um plano de implementação publicado no Confluence e salvá-lo no arquivo `.coder/plan.md` do projeto local.

<context>
Todas as operações usam as ferramentas do MCP `atlassian_local`. As ferramentas seguem o padrão de nomenclatura `atlassian_local_<operacao>`.

**Parâmetros fixos desta skill:**
- **Space:** `CAT`
- **Página raiz:** `Implementações`
- **Destino local:** `.coder/plan.md` no diretório raiz do projeto
</context>

<instructions>

### 1. Localizar a página pai "Implementações"

Buscar a página raiz no space `CAT`:

```
atlassian_local_confluence_search(
  query: "title = \"Implementações\" AND space = \"CAT\"",
  limit: 1
)
```

- Se retornar resultado: extrair o `id` da página — este é o `parent_id`
- Se não retornar resultado: informar ao usuário que a página `Implementações` não foi encontrada no space `CAT` e encerrar

### 2. Listar as subpáginas de "Implementações"

Buscar os filhos diretos da página pai:

```
atlassian_local_confluence_get_page_children(
  page_id: "<parent_id>",
  limit: 50
)
```

- Se retornar **0 resultados**: informar ao usuário que não há planos publicados sob `Implementações` e encerrar
- Se retornar **1 resultado**: selecionar automaticamente essa página e ir para o passo 3
- Se retornar **2 ou mais resultados**: apresentar a lista numerada ao usuário e aguardar a escolha:

```
Foram encontrados os seguintes planos em Confluence › CAT › Implementações:

  1) <título da página 1>
  2) <título da página 2>
  ...

Qual plano deseja baixar? Digite o número:
```

Após a escolha, registrar o `id` e o `title` da página selecionada.

### 3. Obter o conteúdo da página

```
atlassian_local_confluence_get_page(
  page_id: "<id da página selecionada>"
)
```

- Extrair o conteúdo da página no campo `body.wiki.value` (ou equivalente na representação retornada)
- Se a chamada falhar: reportar o erro exato e encerrar

### 4. Salvar em .coder/plan.md

1. Verificar se o diretório `.coder/` existe no diretório raiz do projeto
   - Se não existir: criar o diretório `.coder/`

2. Verificar se `.coder/plan.md` já existe localmente
   - Se existir: **sobrescrever** com o conteúdo obtido do Confluence

3. Escrever o conteúdo no arquivo `.coder/plan.md`

### 5. Reportar o resultado

```
## Plano sincronizado

**Origem:** Confluence › CAT › Implementações › <título da página>
**Destino:** .coder/plan.md
**Operação:** [criação / sobrescrita]
```

Se qualquer chamada ao MCP retornar erro, reportar o erro exato e encerrar.
</instructions>

<rules>
- **MCP exclusivo:** toda leitura do Confluence deve usar as ferramentas `atlassian_local_*` — nunca inventar conteúdo
- **Sobrescrever sem perguntar:** se `.coder/plan.md` já existir localmente, sobrescrever diretamente — não perguntar confirmação
- **Criar diretório se necessário:** garantir que `.coder/` exista antes de escrever o arquivo
- **Sem modificar o Confluence:** esta skill é somente leitura no Confluence — nunca criar ou alterar páginas
- **Sem inventar:** se o MCP falhar, reportar o erro exato
</rules>

<output_format>
## Plano sincronizado

- **Origem:** Confluence › CAT › Implementações › [título]
- **Destino:** `.coder/plan.md`
- **Operação:** [criação / sobrescrita]
</output_format>
