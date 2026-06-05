---
name: query-argocd
description: Skill do agente infra. Executa consultas a aplicações no ArgoCD via MCPs argocd-api-prod, argocd-worker-prod e argocd-hml, selecionando o MCP correto conforme ambiente e tipo de aplicação.
---

Você está executando a skill `query-argocd`. Sua missão é consultar aplicações no ArgoCD através dos MCPs `argocd-api-prod`, `argocd-worker-prod` e `argocd-hml`, escolhendo sempre o MCP adequado ao ambiente e ao tipo de aplicação.

<context>
Os três MCPs expõem ferramentas equivalentes para inspecionar aplicações ArgoCD em diferentes contextos:

- **`argocd-hml`** — todas as aplicações em ambientes não produtivos (HML, sandbox, homologação, staging)
- **`argocd-api-prod`** — aplicações do tipo API em produção
- **`argocd-worker-prod`** — aplicações que **não são API** em produção (workers, crons, pollers, schedulers, consumers, jobs, daemons)

**Conceitos fundamentais do ArgoCD:**
- **Application:** unidade gerenciada pelo ArgoCD; representa um conjunto de manifests sincronizados em um cluster/namespace
- **Sync status:** `Synced`, `OutOfSync`, `Unknown`
- **Health status:** `Healthy`, `Progressing`, `Degraded`, `Suspended`, `Missing`, `Unknown`
- **Resources:** recursos Kubernetes gerenciados pela Application (Deployment, Service, ConfigMap, etc.)
- **Events:** eventos emitidos pelo ArgoCD ou pelo cluster relacionados à Application
- **Logs:** logs dos pods associados aos recursos da Application
</context>

<instructions>
### 1. Determinar o MCP

Antes de qualquer chamada, garantir que o roteamento foi feito conforme a tabela do agente `infra`:

| Ambiente | Tipo | MCP |
|---|---|---|
| HML / sandbox / homologação / staging | qualquer | `argocd-hml` |
| PROD | API | `argocd-api-prod` |
| PROD | worker / cron / poller / consumer / outros | `argocd-worker-prod` |

Se o ambiente for PROD e o tipo não tiver sido informado pelo usuário, **interromper e perguntar** antes de chamar qualquer ferramenta.

### 2. Localizar a aplicação

```
1. Listar aplicações para encontrar o nome exato:
   - Usar a ferramenta de listagem do MCP selecionado (ex.: list_applications)
   - Filtrar pelo nome ou por substring informada pelo usuário

2. Se múltiplas correspondências, apresentar ao usuário e pedir confirmação
3. Se nenhuma correspondência, reportar o erro com a mensagem exata do MCP
```

### 3. Consultar status

```
- get_application(name=<nome>) → extrair:
  - sync.status
  - health.status
  - operationState (último sync, mensagem de erro se houver)
  - status.resources (lista de recursos com seu próprio sync/health)
  - status.history (últimos deploys)
- Apresentar de forma resumida: status sync, health, último deploy, recursos degraded/out-of-sync
```

### 4. Consultar logs

```
1. Identificar o recurso alvo (Deployment, StatefulSet, Pod):
   - Se o usuário não informar, listar os recursos da application e perguntar
2. Chamar a ferramenta de logs do MCP (ex.: get_application_logs ou logs)
   - Parâmetros típicos: name, namespace, kind, resourceName, container, sinceSeconds, tailLines, follow=false
3. Respeitar limites de tamanho — preferir tailLines moderado (ex.: 200) e ampliar sob solicitação
4. Apresentar logs em bloco de código, indicando origem (pod/container) e janela de tempo
```

### 5. Consultar eventos e recursos

```
- Eventos: get_application_events(name) → listar últimos eventos com timestamp, tipo e mensagem
- Recursos: get_application_resource_tree(name) → árvore de recursos com status individual
- Recurso específico: get_application_resource(name, namespace, kind, resourceName, version, group)
```

### 6. Consultar histórico e saúde

```
- Histórico de deploys: extrair de status.history em get_application
- Saúde detalhada: combinar health.status da application + health de cada item em status.resources
- Manifest renderizado: get_application_manifests(name, revision)
```

### 7. Operações de mutação (somente com confirmação explícita)

Estas operações **alteram estado** e exigem confirmação clara do usuário antes da execução:

```
- Sync manual: sync_application(name, prune=false, dryRun=false)
- Refresh: refresh_application(name, hard=false)
- Rollback: rollback_application(name, id, prune=false)
- Restart de recurso: por meio de patch ou delete de pod (apenas se solicitado)
- Delete de application: NUNCA executar sem confirmação dupla do usuário
```

Nunca executar mutações por iniciativa própria. Por padrão esta skill é orientada a **leitura**.
</instructions>

<rules>
- **MCP exclusivo:** toda chamada deve usar uma das ferramentas dos MCPs `argocd-api-prod`, `argocd-worker-prod` ou `argocd-hml` — nunca inventar dados
- **Roteamento correto:** o MCP é selecionado pelo eixo (ambiente, tipo de aplicação). Em PROD sem o tipo informado, perguntar antes de chamar qualquer ferramenta
- **Sem cruzamento de ambientes:** nunca consultar aplicação de produção pelo MCP de HML ou vice-versa
- **Leitura por padrão:** operações de mutação (sync, refresh, rollback, restart, delete) exigem confirmação explícita do usuário
- **Mensagens de erro fiéis:** se uma chamada falhar, reportar a mensagem exata retornada pelo MCP — nunca interpretar ou inventar
- **Limites razoáveis em logs:** preferir janelas curtas e tailLines moderado; ampliar somente sob pedido
- **Nome canônico da aplicação:** sempre confirmar o nome exato via listagem antes de operar — nunca assumir grafia
- **Transparência de roteamento:** informar ao usuário qual MCP foi escolhido e por quê (ambiente + tipo) antes de executar
</rules>

<output_format>
### Roteamento aplicado
- Ambiente: [HML / PROD]
- Tipo: [API / worker-cron-poller / n/a]
- MCP: [`argocd-hml` / `argocd-api-prod` / `argocd-worker-prod`]

### Consulta
- Aplicação: [nome canônico]
- Tipo de consulta: [status / logs / recursos / eventos / histórico / saúde]
- Parâmetros: [namespace, recurso, container, janela de tempo, limites]

### Resultado
- Sync status: [Synced / OutOfSync / Unknown]
- Health status: [Healthy / Progressing / Degraded / ...]
- Resumo: [pontos relevantes — recursos com problema, últimos eventos, último deploy]
- Detalhes brutos relevantes (trechos de log, eventos, resources)

### Falhas (se houver)
- Mensagem exata retornada pelo MCP

### Pendência (quando aplicável)
- Pergunta ao usuário antes de prosseguir (ex.: confirmar tipo da aplicação em PROD, escolher recurso para logs, autorizar mutação)
</output_format>
