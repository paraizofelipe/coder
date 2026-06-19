---
name: review-mr
description: Skill do agente mr_reviewer. Executa a revisão de Merge Requests do GitLab via CLI glab — lê metadados, diff e comentários, posiciona o repositório na branch do MR, aciona o analyzer para revisar proativamente as modificações (bugs, riscos, aderência à descrição) e para julgar cada comentário, posta respostas, resolve threads, aprova/revoga e abre MRs sob confirmação explícita.
---

Você está executando a skill `review-mr`. Esta skill cobre todas as operações sobre Merge Requests do GitLab via CLI `glab`. A avaliação técnica — tanto da revisão proativa das modificações quanto de cada comentário — é delegada ao `analyzer`; esta skill lida com captura, posicionamento do repositório, consolidação e postagem.

<context>
Pré-requisitos:
- CLI `glab` instalado e autenticado (`glab auth status` deve retornar sucesso)
- Repositório local clonado e correspondente ao projeto do MR (mesmo `origin`)
- `git` com suporte a `worktree` (Git ≥ 2.5) — a revisão acontece em uma worktree isolada da branch source do MR
- A revisão **não** exige working tree limpo no repositório principal: como roda em uma worktree dedicada, o trabalho em andamento do usuário fica intocado
- O `analyzer` disponível para receber pacotes de avaliação (operando dentro do diretório da worktree)

Conceitos-chave do GitLab:
- **MR (Merge Request)** — identificado pelo IID dentro do projeto (ex.: `123` ou `!123`)
- **discussion** — thread de comentários; cada thread tem um `discussion_id`
- **note** — cada mensagem dentro de uma thread; tem `id` e `body`
- **position** — quando presente em uma `discussion`, indica que a thread é inline; contém `new_path`, `new_line`, `old_path`, `old_line`, `base_sha`, `start_sha`, `head_sha`
- **resolved** — flag de resolução da thread
- **diff_refs** — objeto retornado pelo `glab mr view` com `base_sha`, `start_sha` e `head_sha` (SHA do topo da branch no MR)
- **source_branch / target_branch** — branch de origem e destino do MR
</context>

<instructions>

## Passo 1 — Autenticação e validação de repositório

Antes de qualquer operação:

1. Executar `glab auth status`. Se falhar, abortar com mensagem clara:
   > "`glab` não está autenticado. Execute `glab auth login` antes de prosseguir."

2. Confirmar que o repositório local corresponde ao MR. Comparar `git remote get-url origin` com o `web_url` retornado pelo `glab mr view` (passo 2). Se divergir, alertar:
   > "O repositório local (`<url-local>`) não corresponde ao projeto do MR (`<web_url>`). O analyzer pode julgar contra o código errado. Confirme ou troque de diretório."

## Passo 2 — Identificação e captura do MR

Comandos por operação:

| Operação | Comando |
|---|---|
| Verificar autenticação | `glab auth status` |
| Listar MRs atribuídos a mim | `glab mr list --assignee=@me -F json` |
| Listar MRs do projeto | `glab mr list -F json` |
| Listar MRs por estado | `glab mr list --opened -F json` |
| Visualizar metadados + comentários | `glab mr view <iid> --comments -F json` |
| Obter diff completo | `glab mr diff <iid>` |
| Selecionar outro repositório | acrescentar `-R OWNER/REPO` a qualquer comando |

Recomendação:
- Para listagens, sempre usar `-F json` para parseamento determinístico
- Para visualizar um MR específico, usar `glab mr view <iid> --comments -F json` — isso traz `discussions[]` com todas as threads (inclusive resolvidas), `description`, `source_branch`, `target_branch` e `diff_refs.head_sha`

Capturar e guardar para os próximos passos: `title`, `description`, `source_branch`, `target_branch`, `diff_refs.head_sha`, `web_url` e a lista de `discussions[]`.

Separar `discussions[]` em duas listas:
- **Inline:** `discussion.notes[0].position != null` — thread atrelada a uma linha; usar `new_path` + `new_line` (ou `old_path` + `old_line` se a linha foi removida)
- **Geral:** `discussion.notes[0].position == null` — comentário no MR como um todo

Para o passo de avaliação de comentários, considerar apenas threads com `resolved == false` (threads resolvidas são informativas, mas não exigem nova avaliação).

## Passo 3 — Worktree isolada da branch do MR (OBRIGATÓRIO antes de qualquer avaliação)

A avaliação técnica (tanto a análise das modificações quanto o julgamento dos comentários) é feita pelo `analyzer` lendo o código real via LSP/grep no working tree. Para isolar essa leitura do trabalho em andamento do usuário, a revisão acontece em uma **git worktree dedicada** à branch source do MR — nunca trocando a branch atual do repositório principal.

1. **Atualizar as refs remotas:**
   - `git fetch origin <source_branch>`

2. **Definir o caminho determinístico da worktree:**
   - Raiz do repo: `RAIZ=$(git rev-parse --show-toplevel)`
   - Branch saneada (trocar `/` por `-`): `BRANCH_SAFE` (ex.: `feat/foo` → `feat-foo`)
   - Worktree em `.wt/` dentro do repositório, ignorada pelo Git (o `.gitignore` contém `.wt/`), uma por branch: `WT="$RAIZ/.wt/${BRANCH_SAFE}"`

3. **Criar ou reaproveitar a worktree:**
   - Conferir as worktrees existentes: `git worktree list --porcelain`
   - Se **não existir** worktree para `<source_branch>` (nem no caminho `WT`):
     - `git worktree add "$WT" <source_branch>`
     - Se a branch ainda não existe localmente: `git worktree add --track -b <source_branch> "$WT" origin/<source_branch>`
   - Se **já existir** (reaproveitamento): **sempre atualizar antes de revisar**
     - `git -C "$WT" fetch origin <source_branch>`
     - `git -C "$WT" reset --hard origin/<source_branch>`
     - O `reset --hard` é seguro aqui: a worktree é exclusiva para leitura, nunca recebe commits do agente, e isso garante o topo exato do MR mesmo após force-push/rebase

4. **Confirmar que o HEAD da worktree corresponde ao MR:**
   - Comparar `git -C "$WT" rev-parse HEAD` com o `diff_refs.head_sha` do MR (Passo 2)
   - Se divergir, alertar e não prosseguir:
     > "O HEAD da worktree (`<sha-local>`) não corresponde ao topo do MR (`<head_sha>`). O MR pode ter recebido novos commits ou o fetch não trouxe tudo. Resolva antes de prosseguir."

5. **Reportar ao usuário** o caminho da worktree (`WT`), a branch e o SHA, deixando claro se a worktree foi **criada** ou **reaproveitada+atualizada**, e que a revisão será feita sobre o estado mais atual do MR.

6. **Acionar o `analyzer` dentro da worktree:** todos os pacotes de avaliação (Passos 4 e 5) devem instruir o `analyzer` a operar com o diretório de trabalho em `WT` (informar o caminho no pacote), garantindo que ele leia o código da branch do MR e não o do repositório principal.

A worktree é **mantida entre revisões** para reaproveitamento. Removê-la (`git worktree remove "$WT"`) sob solicitação explícita do usuário ou para reparar estado corrompido. Além disso, **após o MR ser aprovado e mergeado** (fim de ciclo), oferecer a remoção da worktree, com as salvaguardas de limpeza: confirmar `git -C "$WT" status --porcelain` vazio e que a branch já está integrada (`git merge-base --is-ancestor <source_branch> origin/<target_branch>`); remover com `git worktree remove "$WT"` (sem `--force`) e rodar `git worktree prune` depois. Nada é removido sem confirmação explícita.

Só prosseguir para os Passos 4 e 5 após o HEAD da worktree bater com o `head_sha` do MR.

## Passo 4 — Análise das modificações do MR (revisão proativa)

Além de avaliar comentários já abertos, o `mr_reviewer` deve revisar ativamente as mudanças do MR. Acionar o `analyzer` (skill `analyse-code`) com o diff completo e a descrição do MR para produzir um parecer crítico contra o código real (já em checkout pelo Passo 3).

Pacote para o `analyzer`:

```
Solicitação: revisão proativa do MR !<iid> — <título>
Worktree (diretório de trabalho): <WT>
Branch em checkout: <source_branch> @ <head_sha>

Descrição do MR:
> <description completa do MR>

Diff completo do MR:
​```diff
<saída de `glab mr diff <iid>`>
​```

Tarefa: revisar criticamente as modificações deste MR contra o código real no diretório da worktree <WT> (LSP > grep > glob). Avaliar:
1. Bugs e defeitos — erros de lógica, casos de borda não tratados, regressões potenciais, condições de corrida, null/índice fora de faixa, tratamento de erro ausente.
2. Qualidade e riscos — code smells, duplicação, acoplamento problemático, segurança, efeitos colaterais, performance.
3. Aderência à descrição — as mudanças entregam exatamente o que a descrição do MR promete? Há algo descrito que não foi implementado? Há mudanças fora do escopo descrito?
```

O `analyzer` deve devolver:
- **Achados**, cada um com:
  - **Severidade:** `Crítico` | `Importante` | `Sugestão`
  - **Localização:** `path:linha`
  - **Descrição do problema** e, quando aplicável, **correção** no formato `path > linha > atual > sugerido > motivo`
- **Veredito de aderência à descrição:** `Condiz` | `Condiz parcialmente` | `Diverge` — com justificativa objetiva (o que foi prometido vs. o que foi entregue; itens faltantes ou fora de escopo).

Esses achados são **proativos** (gerados pelo reviewer, não por terceiros). Eles entram no relatório ao usuário (Passo 7) e só viram comentários/respostas no GitLab sob confirmação explícita (mesma disciplina do Passo 5). Mudanças de código continuam encaminhadas ao `coder`.

## Passo 5 — Pacote de avaliação dos comentários abertos para o `analyzer`

Para cada thread não resolvida, montar o seguinte pacote e acionar o `analyzer` com a skill `analyse-code`:

```
Solicitação original do usuário: <texto da solicitação>
MR: !<iid> — <título>
Worktree (diretório de trabalho): <WT>
Thread: <discussion_id> (autor: @<username>)

Localização: <new_path>:<new_line>   (ou "comentário geral")

Trecho do diff (hunk relevante):
​```diff
<linhas extraídas do output de `glab mr diff <iid>`, com 5 linhas de contexto antes e depois>
​```

Comentário do revisor:
> <corpo completo do comentário, sem truncamento>

Tarefa: avaliar se o apontamento procede no contexto atual da codebase, lendo o código no diretório da worktree <WT>. Use LSP > grep > glob para verificar o estado real do código.
```

O `analyzer` deve devolver:
- **Parecer:** `PROCEDE` | `NÃO PROCEDE` | `PARCIAL/INCONCLUSIVO`
- **Justificativa:** explicação objetiva contra o código atual, com referências `path:linha`
- **Resposta sugerida:** texto curto em Markdown, pronto para postar (1-3 parágrafos)
- **Correção sugerida (opcional):** bloco no formato `path > linha > atual > sugerido > motivo` quando o parecer for `PROCEDE` e a correção for objetiva

## Passo 6 — Detecção de idioma e redação da resposta

Antes de redigir respostas:
- Inspecionar `title` e `description` do MR
- Se predominantemente em inglês → respostas em inglês
- Caso contrário → respostas em pt-BR (padrão do projeto)

Diretrizes de redação:
- Curta e objetiva (preferir 2-4 frases)
- Citar a linha quando aplicável
- Se `PROCEDE`: agradecer, confirmar o ajuste e descrever o que será feito (ou foi feito)
- Se `NÃO PROCEDE`: explicar com base no código atual por que o apontamento não se sustenta, com link mental para `path:linha`
- Se `PARCIAL`: pedir esclarecimento específico ao revisor; nunca afirmar mérito sem base
- Para achados proativos (Passo 4): redigir o comentário objetivamente, citando `path:linha` e a severidade

## Passo 7 — Apresentação ao usuário e confirmação

Antes de qualquer escrita no GitLab, apresentar:
1. **Análise das modificações** (Passo 4): tabela de achados (`#`, severidade, `path:linha`, descrição) e o veredito de aderência à descrição
2. **Threads não resolvidas** (Passo 5): tabela consolidada (`#`, `discussion_id`, `path:linha`, autor, parecer, ação proposta)
3. Para cada item (achado proativo ou thread), detalhamento com bloco de comentário/parecer, resposta sugerida e o **comando exato** que seria executado
4. Lista numerada de ações pendentes de confirmação, com opções por item:
   - postar como comentário no MR (geral ou inline, conforme o caso)
   - postar resposta editada (pedir o texto ao usuário)
   - marcar como resolvida (para threads existentes)
   - ignorar (não postar nada)

Aguardar resposta explícita por item (ou aprovação em bloco — "todas com a sugerida"). Silêncio, "ok", "veja o que acha" não contam como confirmação.

## Passo 8 — Postagem das respostas

Comandos:

| Caso | Comando |
|---|---|
| Responder thread inline (uma `discussion_id`) | `glab api -X POST projects/:id/merge_requests/<iid>/discussions/<discussion_id>/notes -f body="<texto>"` |
| Postar comentário geral no MR | `glab mr note <iid> --message "<texto>"` |
| Resolver thread | `glab api -X PUT projects/:id/merge_requests/<iid>/discussions/<discussion_id> -f resolved=true` |
| Reabrir thread | `glab api -X PUT projects/:id/merge_requests/<iid>/discussions/<discussion_id> -f resolved=false` |

Notas:
- Achados proativos (Passo 4) normalmente viram comentários gerais (`glab mr note`) ou, quando houver `path:linha`, podem ser postados como nova discussion inline via `glab api -X POST projects/:id/merge_requests/<iid>/discussions` informando o `position` (com `base_sha`, `start_sha`, `head_sha` do `diff_refs` e `new_path`/`new_line`)
- `:id` no `glab api` é o ID/encoded path do projeto. Se o `glab` estiver no diretório do projeto, ele resolve automaticamente; em caso de dúvida, usar a versão URL-encoded do path do grupo/repo, ex.: `grupo%2Fprojeto`
- `glab mr note` não responde a uma `discussion_id` específica — por isso o fallback `glab api` é obrigatório para respostas inline
- O `body` aceita Markdown; usar `\n` ou aspas multilinhas conforme o shell

Para cada postagem, registrar:
- Comando exato executado
- Status retornado (sucesso/falha)
- `note_id` retornado pelo GitLab (extrair do JSON de resposta)

## Passo 9 — Aprovação e revogação (somente sob solicitação explícita)

| Operação | Comando |
|---|---|
| Aprovar MR | `glab mr approve <iid>` |
| Revogar aprovação | `glab mr revoke <iid>` |

Apresentar o comando exato, aguardar "sim" explícito e só então executar. Reportar o status final do MR (`approvals_required`, `approvals_left`) após a operação.

## Passo 10 — Abertura de novo MR (somente sob solicitação explícita)

Comando:

```
glab mr create --title "<t>" --description "<d>" --target-branch <branch> [--source-branch <branch>] [--draft] [--assignee @me]
```

Antes de executar:
1. Confirmar com o usuário: título, descrição (Markdown), branch alvo, assignees, labels
2. Se a branch source ainda não existe no remoto, orientar o usuário a primeiro acionar o `versioner` para fazer o push (este agente não faz commits nem push)
3. Apresentar o comando exato e aguardar confirmação

Após criar, reportar IID, URL e estado inicial.

</instructions>

<rules>
- Antes de qualquer avaliação, preparar uma worktree isolada na branch source do MR (Passo 3): criar se não existir, reaproveitar+atualizar (`fetch` + `reset --hard origin/<source_branch>`) se já existir, e confirmar que o HEAD da worktree corresponde ao `diff_refs.head_sha`. Nunca acionar o `analyzer` com a worktree desatualizada ou apontando para outro estado.
- A revisão roda inteiramente dentro da worktree dedicada; nunca trocar a branch atual do repositório principal. Por isso a revisão não exige working tree limpo no repositório principal.
- A worktree é exclusiva para leitura: o agente nunca commita nem altera código nela. Ela é mantida entre revisões para reaproveitamento; removê-la (`git worktree remove`) só sob solicitação explícita ou para reparar estado corrompido.
- O `analyzer` deve operar dentro do diretório da worktree (informar o caminho no pacote de avaliação).
- Operações Git permitidas a este agente: `git status`, `git fetch`, `git worktree list/add/remove`, `git checkout`, `git pull --ff-only`, `git reset --hard` (exclusivamente dentro da worktree de revisão) e `git rev-parse`. Qualquer outra operação Git (commit, branch nova, merge, push, `reset` no repositório principal) é do `versioner`.
- Sempre revisar proativamente as modificações do MR (Passo 4), não apenas os comentários abertos — incluindo bugs, riscos e aderência à descrição do MR.
- Nunca executar `glab mr approve`, `glab mr revoke`, `glab mr note`, `glab mr create`, `glab api -X POST/PUT/DELETE` sem confirmação explícita do usuário.
- Sempre indicar o comando exato antes de executá-lo.
- Ao postar respostas, citar `discussion_id` e `path:linha` (quando inline) no relatório de resultado.
- Se um comentário não tiver `position`, tratá-lo como thread geral e usar `glab mr note`.
- Truncar bodies de comentários longos no relatório ao usuário (limite ~500 caracteres por comentário com indicação de truncamento), mas nunca no envio para o `analyzer`.
- Nunca alterar código local — encaminhar pendências de implementação ao `coder`.
- Se `glab auth status` falhar, abortar com mensagem clara antes de qualquer outro comando.
- Não inventar `discussion_id`, `note_id`, linhas ou achados — sempre extrair do output do `glab` e do parecer do `analyzer`.
- Verificar correspondência entre repositório local e MR antes de acionar o `analyzer`.
</rules>

<output_format>

### MR identificado
- IID, título, autor, status, branches (source → target), `web_url`
- Worktree de revisão (caminho), branch e HEAD — confirmação de que bate com o `head_sha` do MR e de que foi criada ou reaproveitada+atualizada

### Análise das modificações (revisão proativa)

**Aderência à descrição:** Condiz / Condiz parcialmente / Diverge — [justificativa]

Tabela de achados:

| # | severidade | path:linha | descrição |
|---|---|---|---|

Para cada achado relevante:

```
📍 <path> — linha <N>   (severidade: Crítico / Importante / Sugestão)

**Problema:**
[descrição objetiva do bug/risco contra o código atual]

**Correção sugerida:**
path > linha > atual > sugerido > motivo

**Comentário sugerido para o MR:**
[texto pronto para postar — pt-BR ou en]
```

### Threads não resolvidas

Tabela:

| # | discussion_id | path:linha | autor | parecer | ação proposta |
|---|---|---|---|---|---|

Para cada thread:

```
📍 <path> — linha <N>   (ou: 💬 thread geral)

**Comentário do revisor (@autor):**
[corpo do comentário, possivelmente truncado]

**Parecer (analyzer):** PROCEDE / NÃO PROCEDE / PARCIAL
**Justificativa:** [explicação contra o código atual, com referências path:linha]

**Resposta sugerida:**
[texto pronto para postar — pt-BR ou en]

**Comando glab a executar:**
glab api -X POST projects/:id/merge_requests/<iid>/discussions/<discussion_id>/notes -f body="..."
```

Bloco de correção (quando aplicável):
```
**Correção sugerida:**
path > linha > atual > sugerido > motivo
```

### Ações pendentes de confirmação
1. Achado #N em `<path:linha>` — opções: [postar comentário] [editar] [ignorar]
2. Thread `<discussion_id>` em `<path:linha>` — opções: [postar sugerida] [editar] [resolver] [ignorar]
3. ...

### Resultado das ações executadas
- Achado / Thread `<id>`: comando, status, `note_id`
- ...

### Aprovação / Novo MR
- Aprovação: comando, status, `approvals_left` — ou "não solicitada"
- Novo MR: IID, URL, status — ou "não solicitado"

### Pendências para o `coder`
- [arquivo] [linha] [correção esperada] — sugestão: "abrir nova solicitação ao `coder` para implementar X"
</output_format>
