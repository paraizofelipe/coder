Acione o agente `mr_reviewer` com a skill `review-mr` para revisar o Merge Request identificado por `$ARGUMENTS`.

## Passos

### 1. Identificar o MR a partir de `$ARGUMENTS`

- Aceitar IID puro (`123`), IID com prefixo (`!123`) ou URL completa (`https://.../-/merge_requests/N`)
- Se `$ARGUMENTS` estiver vazio: listar MRs atribuídos com `glab mr list --assignee=@me -F json` e perguntar qual revisar
- Validar que o repositório local corresponde ao MR (`git remote get-url origin` × `web_url`)

### 2. Carregar o MR com `glab`

- `glab auth status` (abortar se falhar)
- `glab mr view <iid> --comments -F json` para metadados, descrição, `source_branch`, `diff_refs.head_sha` e `discussions[]`
- `glab mr diff <iid>` para mapear `path:linha` dos hunks

### 3. Preparar a worktree isolada da branch do MR

Garantir que a revisão seja feita sobre o código mais atual do MR, em uma worktree dedicada que não toca na branch atual do repositório principal:
- `git fetch origin <source_branch>` para atualizar as refs remotas
- Definir o caminho determinístico da worktree em `.wt/` dentro do repositório (`$RAIZ/.wt/<branch-safe>`, ignorada pelo Git), branch com `/` saneada para `-`, e checar `git worktree list --porcelain`
- Se **não existe**: `git worktree add <WT> <source_branch>` (ou `git worktree add --track -b <source_branch> <WT> origin/<source_branch>`)
- Se **já existe**: reaproveitar e **sempre atualizar antes de revisar** — `git -C <WT> fetch origin <source_branch>` + `git -C <WT> reset --hard origin/<source_branch>`
- Confirmar que `git -C <WT> rev-parse HEAD` corresponde ao `diff_refs.head_sha`; se divergir, alertar e não prosseguir
- Acionar o `analyzer` com o diretório de trabalho em `<WT>`

### 4. Revisar proativamente as modificações acionando o `analyzer`

Enviar ao `analyzer` (skill `analyse-code`) a descrição do MR + o diff completo. Receber achados (severidade `Crítico`/`Importante`/`Sugestão` + `path:linha` + correção) e veredito de aderência à descrição (**Condiz / Condiz parcialmente / Diverge**) — bugs, riscos e o que foi prometido vs. entregue.

### 5. Avaliar cada thread não resolvida acionando o `analyzer`

Para cada `discussion` com `resolved == false`, montar pacote (solicitação original, `path:linha`, hunk do diff, corpo do comentário) e delegar julgamento ao `analyzer`. Receber parecer **PROCEDE / NÃO PROCEDE / PARCIAL** com justificativa e resposta sugerida.

### 6. Apresentar tabelas consolidadas e aguardar confirmação

Mostrar a análise das modificações (achados + aderência) e as threads avaliadas (`discussion_id`, `path:linha`, autor, parecer, ação proposta). Para cada item, oferecer: postar comentário sugerido / postar editado / marcar como resolvida / ignorar. Nada vai ao GitLab sem "sim" explícito.

### 7. Executar ações confirmadas

Postar achados proativos via `glab mr note` (ou nova discussion inline via `glab api ... discussions` com `position`), respostas inline via `glab api ... discussions/<id>/notes`, resolver threads via `glab api ... -X PUT`. Aprovar (`glab mr approve`) ou abrir novo MR (`glab mr create`) somente sob solicitação explícita.

### 8. Reportar resultado

Resumir comandos executados, `note_id` retornados, status de aprovação e listar pendências de implementação para o `coder`.
