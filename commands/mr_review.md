---
description: Revisa um Merge Request do GitLab — lê metadados/diff/comentários via glab, aciona o analyzer para julgar cada apontamento na linha indicada e prepara respostas/aprovação sob confirmação. Uso: /mr_review 123  ou  /mr_review https://gitlab.com/grupo/projeto/-/merge_requests/123
agent: mr_reviewer
---

Acione o agente `mr_reviewer` com a skill `review_mr` para revisar o Merge Request identificado por `$ARGUMENTS`.

## Passos

### 1. Identificar o MR a partir de `$ARGUMENTS`

- Aceitar IID puro (`123`), IID com prefixo (`!123`) ou URL completa (`https://.../-/merge_requests/N`)
- Se `$ARGUMENTS` estiver vazio: listar MRs atribuídos com `glab mr list --assignee=@me -F json` e perguntar qual revisar
- Validar que o repositório local corresponde ao MR (`git remote get-url origin` × `web_url`)

### 2. Carregar o MR com `glab`

- `glab auth status` (abortar se falhar)
- `glab mr view <iid> --comments -F json` para metadados, descrição e `discussions[]`
- `glab mr diff <iid>` para mapear `path:linha` dos hunks

### 3. Avaliar cada thread não resolvida acionando o `analyzer`

Para cada `discussion` com `resolved == false`, montar pacote (solicitação original, `path:linha`, hunk do diff, corpo do comentário) e delegar julgamento ao `analyzer` (skill `analyse_code`). Receber parecer **PROCEDE / NÃO PROCEDE / PARCIAL** com justificativa e resposta sugerida.

### 4. Apresentar tabela consolidada e aguardar confirmação

Mostrar todas as threads avaliadas com `discussion_id`, `path:linha`, autor, parecer e ação proposta. Para cada thread, oferecer: postar resposta sugerida / postar resposta editada / marcar como resolvida / ignorar. Nada vai ao GitLab sem "sim" explícito.

### 5. Executar ações confirmadas

Postar respostas inline via `glab api ... discussions/<id>/notes`, comentários gerais via `glab mr note`, resolver threads via `glab api ... -X PUT`. Aprovar (`glab mr approve`) ou abrir novo MR (`glab mr create`) somente sob solicitação explícita.

### 6. Reportar resultado

Resumir comandos executados, `note_id` retornados, status de aprovação e listar pendências de implementação para o `coder`.