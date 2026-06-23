# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## O que é este repositório

Coleção de definições de agentes e skills para [OpenCode](https://opencode.ai), [Claude Code](https://claude.ai/code), [Codex](https://github.com/openai/codex) e [Pi](https://github.com/earendil-works/pi-coding-agent) que implementa um fluxo disciplinado de desenvolvimento de software assistido por IA. Não há código executável, dependências ou comandos de build — todos os arquivos são Markdown mais um shell script (`install.sh`).

## Estrutura

```
agents/                  → Definições de agentes (um subdiretório por agente)
  <name>/
    body.md              → Corpo do agente (instruções, XML tags)
    opencode.yml         → Frontmatter para OpenCode (description, mode, model se primário)
    claude.yml           → Frontmatter para Claude Code (Phase 2)
skills/                  → Skills no formato Agent Skills (https://agentskills.io/specification)
  <kebab-name>/
    SKILL.md             → Metadados (name, description) + instruções
    references/          → (opcional) material consultado sob demanda
commands/                → Slash commands (um subdiretório por command)
  <kebab-name>/
    body.md              → Corpo do command
    opencode.yml         → Frontmatter para OpenCode
    claude.yml           → Frontmatter para Claude Code (Phase 2)
    pi.yml               → Frontmatter para Pi (description, argument-hint)
install.sh               → Monta e copia agents/, skills/ e commands/ para os diretórios nativos dos harnesses; aplica modelo do vendor
```

## Agentes e skills

| Agente | Mode | Skill | Papel |
|---|---|---|---|
| `coder` | `primary` | `write-code` | Orquestrador de implementação — aciona todos os subagentes de desenvolvimento e delega operações Kanban ao `kanban` |
| `lead` | `primary` | `plan-implementation` | Orquestrador de planejamento — recebe feature/bug, delega análise, resolve ambiguidades com o usuário, gera `.coder/task-AAAAMMDD-HHMMSS.md` e delega implementação ao `coder` após aprovação |
| `documenter` | `primary` | `document-plan`, `get-plan` | Publica e sincroniza planos de implementação com o Confluence via MCP `atlassian_local` |
| `kanban` | `primary` | `kanban-force` | Gerencia cards e boards via MCP `kanban-force` |
| `infra` | `primary` | `query-argocd` | Consulta aplicações no ArgoCD via MCPs `argocd-api-prod`, `argocd-worker-prod` e `argocd-hml` (também atua como subagente do `coder`) |
| `mr_reviewer` | `primary` | `review-mr` | Revisa Merge Requests do GitLab via CLI `glab`; aciona o `analyzer` para julgar comentários inline e posta respostas/aprovações sob confirmação (também atua como subagente do `coder`) |
| `qa` | `primary` | `validate-implementation` | Valida funcionalmente as modificações da branch atual (smoke, black-box, e2e, regressão) contra serviços reais; aciona `analyzer` (modificações + descoberta de serviços) e `infra` (ambiente/logs); grava `.coder/tests-AAAAMMDD-HHMMSS.md`, valida acessos e executa em HML sob confirmação. Nunca altera código |
| `analyzer` | `subagent` | `analyse-code` | Inspeciona a codebase antes de qualquer modificação (subagente comum a `lead`, `coder` e `qa`) |
| `clarifier` | `subagent` | `clarify-intent` | Transforma ambiguidades brutas do `analyzer` em perguntas com opções e recomendação justificada — não conversa com o usuário, quem apresenta é o `lead` |
| `planner` | `subagent` | `plan-tasks` | Produz o TaskGraph esqueleto (lista de tasks com dependências e riscos) a partir da intenção esclarecida e do relatório do `analyzer` |
| `detailer` | `subagent` | `detail-tasks` | Enriquece cada task do TaskGraph com motivação, arquivos, preview, estratégia de teste, critérios, contrato, done when e esforço |
| `tester` | `subagent` | `test-code` | Cria e executa testes com abordagem TDD |
| `code_reviewer` | `subagent` | `review-code` | Revisão técnica — Camada 1 |
| `business_reviewer` | `subagent` | `review-code` | Revisão de negócio/segurança — Camada 2 (mesma skill, papel diferente) |
| `versioner` | `subagent` | `version-code` | Executa operações Git com confirmação explícita; herda modelo do agente que o aciona |

## Fluxo calibrado por impacto

O `coder` triado a solicitação em um de quatro níveis **antes** de decidir o fluxo. O nível é anunciado em 1 linha ao usuário antes de qualquer ação. Em todos os níveis valem: confirmação antes de tocar em arquivo, confirmação antes de versionar, nenhuma alteração em `main`/`master`, e decisão sobre testes sempre é do `tester`.

| Nível | Critérios | Fluxo resumido |
|---|---|---|
| **Trivial** | Typo, doc, rename local, ≤ ~5 linhas em 1 arquivo, sem regra de negócio | resumo 1 linha → confirma → implementa → `tester` decide se cabe teste (normalmente não) → confirma → `versioner` |
| **Pequena** | 1–2 arquivos, função isolada, bug com causa visível, sem mudança de API/negócio | `analyzer` focado (opcional) → plano inline → confirma → `tester` decide testes → implementa → `tester` executa → `code_reviewer` (opcional) → confirma → `versioner` |
| **Média** | 3–5 arquivos, comportamento contido, refactor localizado | `analyzer` focado → plano inline → confirma → `tester` → implementa → `tester` executa + regressão → `code_reviewer` → `business_reviewer` (se houver risco) → confirma → `versioner` |
| **Grande** | Nova feature, múltiplos módulos, segurança/auth/dados, regra de negócio, API pública, ambiguidade séria, > 5 arquivos | `analyzer` completo → `.coder/plan.md` + loop de ambiguidades → confirma → `tester` (red) → implementa → `tester` (green + regressão) → `code_reviewer` → `business_reviewer` (obrigatório) → confirma → `versioner` |

**Escalonamento obrigatório** (vira **Grande** independente do tamanho): autenticação, autorização, criptografia, secrets, schema/migrations, contratos de API/eventos, pagamento/billing, PII/PCI, ambiguidade no pedido, falta de contexto para mapear arquivos.

**Em caso de dúvida entre dois níveis, escolher o mais alto.**

**Regra de branch:** o `versioner` verifica se a branch atual é `main` ou `master`. Se for, consulta o usuário sobre o nome da nova branch (ou gera um nome curto em kebab-case). Nenhum arquivo é alterado na branch principal. Se já estiver em uma branch de trabalho, reporta ao `coder` e aguarda instrução.

**Worktrees (`.wt/`):** em **Média e Grande**, o trabalho é isolado em uma git worktree dedicada em `.wt/<branch-safe>` dentro do repositório (ignorada pelo Git via `.gitignore`, uma por branch). Decisão de entrada: na `main`/`master`, nunca usar a principal — cria-se branch nova + worktree; em branch de trabalho, pergunta-se ao usuário se usa a atual (o trabalho fica no repositório principal, pois o Git não permite worktree de uma branch já ativa) ou cria uma nova (worktree em `.wt/<nova>`). Quando o trabalho roda na worktree, todos os subagentes operam com o diretório de trabalho lá. **Trivial e Pequena** seguem o fluxo simples, sem worktree. A revisão de MR (`mr_reviewer`) também usa `.wt/<branch>`. **Limpeza:** worktrees são mantidas e reaproveitadas; a remoção é sempre sob confirmação, por ciclo de vida (oferta ao integrar/mergear a branch) ou varredura sob demanda ("limpar worktrees"), com salvaguardas (working tree limpo, branch mergeada/gone, sem `--force`, `git worktree prune` ao final).

**`.coder/plan.md`:** criado **apenas no nível Grande**. Contém a solicitação original, resumo da análise, tabela de ambiguidades com decisões tomadas, plano de ação e riscos. Em Média o plano é inline na resposta; em Pequena são bullets; em Trivial é uma linha.

**Solicitação Kanban** (ID de card ou operação de board/card): `coder → kanban`

**Solicitação de Merge Request** (verificar/analisar/responder/aprovar/abrir): `coder → mr_reviewer → (analyzer por comentário) → confirmação do usuário → glab`

**Solicitação mista:** Kanban primeiro, fluxo de código depois

## Fluxo de planejamento (Lead)

O `lead` é um agente primário paralelo ao `coder`, escolhido quando o usuário quer **um plano técnico em tasks revisáveis antes de implementar** (feature nova, fix complexo, refactor que cruza módulos). O `coder` continua sendo o ponto de entrada padrão para mudanças triviais/pequenas/médias que não precisam de quebra prévia.

Pipeline:

```
lead
 ├─ analyzer (analyse-code)                    inspeciona codebase, identifica ambiguidades brutas
 ├─ clarifier (clarify-intent, se houver amb.) formata perguntas com opções + recomendação
 ├─ <loop de decisões com o usuário>           uma pergunta por vez, decisão registrada
 ├─ planner (plan-tasks)                       TaskGraph esqueleto (tasks + dependências + riscos)
 ├─ detailer (detail-tasks)                    enriquece cada task com preview/testes/critérios/contrato
 ├─ grava .coder/task-*.md                     documento canônico
 ├─ apresenta RESUMO (≤15 linhas, não o documento inteiro)
 └─ aprovação explícita → delega ao coder      (toda implementação, branch, testes, review, versionamento)
```

Quem implementa continua sendo o `coder` — após aprovação do `.coder/task-AAAAMMDD-HHMMSS.md`, o `lead` faz hand-off completo e o `coder` aplica seu fluxo normal (triagem, TDD, reviewers, versioner) para cada task selecionada.

**`.coder/task-AAAAMMDD-HHMMSS.md`:** artefato canônico do `lead`, gravado em `.coder/` com o timestamp da criação (ex.: `.coder/task-20260716-102717.md`) — diferente de `.coder/plan.md`, que é do `coder` no nível Grande. Contém solicitação original, contexto técnico, decisões tomadas, riscos e tasks detalhadas. Cada nova solicitação de planejamento gera um novo arquivo; ajustes na mesma sessão atualizam o arquivo já criado, anexando o bloco `## Histórico de iterações` e preservando decisões anteriores.

**Regras invariáveis:** o `lead` nunca escreve código de produção, nunca cria branch, nunca commita; tudo isso continua via `coder → versioner`.

## Fluxo de validação (QA)

O `qa` é um agente primário paralelo ao `coder`/`lead`, escolhido quando o usuário quer **validar funcionalmente as modificações da branch atual contra serviços reais** — smoke, black-box, e2e e regressão — focando regras de negócio e caçando falhas no fluxo. É distinto do `tester` (que faz TDD de código na codebase): o `qa` exercita o sistema **em execução**.

Pipeline (skill `validate-implementation`):

```
qa
 ├─ enquadra a branch atual (não main/master)
 ├─ analyzer (analyse-code)        modificações da branch + descoberta de serviços que o projeto acessa
 ├─ monta o plano de testes        casos por regra de negócio (tipo, classe leitura/mutação, serviços)
 ├─ grava .coder/tests-*.md        artefato canônico do qa (timestamp da criação)
 ├─ apresenta RESUMO + aprovação   ≤15 linhas; sem "sim" explícito não avança
 ├─ valida acessos                 probe read-only por serviço (MCP/CLI/curl); infra para ambiente/logs
 ├─ gate de acessos                OK/FALTA/FALHA; sem acesso essencial → BLOQUEADO
 ├─ executa testes (HML)           leitura direto; mutação confirmada por item; registra evidência real
 └─ relatório + pendências         status por caso; correções encaminhadas ao coder
```

**`.coder/tests-AAAAMMDD-HHMMSS.md`:** artefato canônico do `qa`, gravado em `.coder/` com o timestamp da criação (mesmo padrão do `.coder/task-*.md` do `lead`). Contém alvo, contexto técnico, serviços/acessos, plano de testes e histórico de iterações.

**Regras invariáveis:** o `qa` nunca escreve código de produção, nunca cria branch, nunca commita, nunca corrige código; HML é o ambiente padrão (PROD nunca default); mutações exigem confirmação por item; nada de inventar resultado (sem acesso = BLOQUEADO). Correções vão para o `coder`.

## Padrão Agent Skills

As skills seguem a especificação aberta **Agent Skills** (https://agentskills.io/specification), o mesmo padrão adotado por OpenCode, Claude Code, Cursor, Codex e outros clientes. Pontos essenciais:

- Cada skill é uma **pasta** dentro de `skills/` contendo um arquivo `SKILL.md`
- O nome da pasta é em kebab-case (`[a-z0-9-]`, ≤ 64 chars) e **deve ser igual ao campo `name`** do frontmatter
- `references/` é opcional e contém material extenso que o agente carrega só quando precisa (templates, payloads detalhados, checklists). Use referências relativas de **um nível só** (`references/<arquivo>.md`)
- O `SKILL.md` deve ficar **abaixo de 500 linhas**; mover material extenso para `references/`
- Carregamento progressivo: metadados (`name` + `description`) na descoberta → corpo do `SKILL.md` na ativação → arquivos em `references/` sob demanda

Para validar uma skill antes de commitar:

```bash
npx -y skills-ref validate ./skills/<nome>
```

## Convenções ao editar arquivos

### Frontmatter YAML

**Agentes** (`agents/<name>/{body.md,opencode.yml,claude.yml}`) — feature do OpenCode, fora da spec Agent Skills:
- `description`, `mode` (`primary` | `subagent`), `model` (apenas primários), `temperature`
- Nunca editar `model:` manualmente — o `install.sh` sobrescreve esse campo na instalação
- Subagentes **não têm** campo `model` (herdam do agente chamador)

**Skills** (`skills/<name>/SKILL.md`) — conforme a spec Agent Skills:
- `name` (obrigatório) — kebab-case, igual ao nome da pasta
- `description` (obrigatório) — descreve **o que faz e quando usar**, ≤ 1024 caracteres, com palavras-chave que ajudem o agente a ativar a skill
- `license`, `compatibility`, `metadata`, `allowed-tools` (opcionais)

### Estrutura XML

Todo conteúdo é estruturado com XML tags. Tags em uso:

```
<role>  <responsibilities>  <rules>  <workflow>  <instructions>
<output_format>  <checklist>  <principles>  <criteria>  <context>  <code_navigation>
<objetivo>  <triage>  <subagents>  <priorities>
```

Não alterar tags XML sem verificar todos os arquivos que as utilizam. Não remover `<output_format>` de nenhum arquivo — é o contrato de resposta do agente.

### Formato de diff nas skills de revisão

Toda sugestão nos reviewers usa o formato estruturado:
```
path > linha > atual > sugerido > motivo
```

### Idioma

Português do Brasil em todo conteúdo (agentes, skills, comentários). Commits em inglês.

## Commits

Conventional Commits em inglês: `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`, `style:`, `perf:`, `test:`

## O que não fazer

- Não adicionar código executável, dependências ou configuração de build — o repositório é apenas Markdown + 1 shell script
- Não alterar a estrutura de XML tags sem verificar todos os arquivos que a utilizam
- Não remover `<output_format>` de nenhum arquivo
- Não editar `model:` manualmente nos agentes — o `install.sh` gerencia esse campo

## install.sh

Copia `agents/`, `skills/` e `commands/` para o diretório nativo de cada harness selecionado (OpenCode, Claude Code, Codex, Pi) e monta cada agente/command juntando `<harness>.yml` + `body.md`. Substitui o campo `model:` dos agentes primários via `sed` conforme o vendor escolhido (no Codex e no Pi `apply_model` é no-op). Subagentes não recebem `model`. O Codex e o Pi não recebem agentes nativos: as skills são copiadas, os commands viram prompts (`~/.codex/prompts/`, body-only; `~/.pi/agent/prompts/`, montado com `pi.yml`) e a orquestração vai pelo `AGENTS.md`.

Skills são copiadas como **diretórios completos** (`skills/<name>/SKILL.md` + `references/`), preservando a estrutura do padrão Agent Skills.

No modo remoto, o instalador faz `git clone --depth 1` do repositório em diretório temporário e copia a partir daí — git é dependência obrigatória nesse modo.

Flags: `--force` (sobrescreve sem perguntar), `--local` (usa arquivos locais em vez de clonar do GitHub).
