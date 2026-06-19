# Matriz de serviços e acessos

Referência para descobrir os serviços que o projeto acessa (Passo 2/4) e para provar acesso a cada um antes de executar testes (Passo 7). A lista de tipos é **extensível** ao que o projeto realmente usa — o exemplo abaixo é ponto de partida, não limite.

## Parte A — Descoberta de serviços

O `analyzer` procura indícios de dependências externas em:

- **Variáveis de ambiente** — `.env*`, `*.env`, charts/values, `ConfigMap`/`Secret` (hosts, URLs, DSNs, tópicos)
- **Compose / containers** — `docker-compose*.yml`, `Dockerfile` (serviços vinculados: db, cache, broker)
- **Configuração de código** — clientes HTTP/SDK instanciados, `baseURL`/endpoints, strings de conexão de banco, nomes de fila/tópico
- **MCPs configurados** — config de MCP do projeto/harness (ex.: catlog, mongodb, argocd) e aliases disponíveis
- **Manifests de ambiente** — k8s/ArgoCD (`Application`, `Deployment`, `Service`, `Ingress`) para nome da app e ambiente

Saída: lista de serviços, cada um com **tipo** (DB / API / ambiente-logs / fila-broker / outro) e o identificador (alias, host, nome da app).

## Parte B — Prova de acesso por tipo

Para cada serviço, escolher o método disponível e rodar o probe **read-only** antes de testar. Sempre mostrar o comando exato ao usuário.

| Tipo de serviço | Métodos possíveis | Probe read-only (exemplo) | Remediação se falhar |
|---|---|---|---|
| **API (MCP)** | MCP catlog | `catlog_list_apis` / `catlog_list_operations api=<alias>` | `catlog sync <alias>`; conferir alias com `catlog apis` |
| **API (HTTP)** | curl | `curl -sS -o /dev/null -w '%{http_code}\n' <base>/health` | verificar URL/rede/credencial; checar se o serviço está no ar |
| **Banco de dados** | MCP mongodb / CLI | `list-collections` (MCP) ou `<cli> --eval 'ping'` | conferir string de conexão/credencial; abrir túnel/VPN |
| **Ambiente / logs** | MCP argocd / `infra` / kubectl | acionar `infra` (`query-argocd`) para status; `kubectl auth can-i get pods -n <ns>` | autenticar no cluster; conferir contexto/namespace |
| **Fila / broker** | CLI específica / MCP | comando de "describe"/"list topics" read-only | conferir host/credencial do broker |
| **GitLab** | `glab` | `glab auth status` | `glab auth login` |

Notas:
- O probe deve ser **inócuo**: listar, descrever, health, `auth status`, `can-i` — nunca criar/alterar/apagar.
- Ambiente padrão dos probes e testes: **HML**. Para PROD, exigir confirmação reforçada.
- Resultado do probe ∈ `OK` (acesso confirmado) / `FALTA` (ferramenta/credencial ausente) / `FALHA` (erro ao acessar). Sem `OK`, o teste dependente fica `BLOQUEADO` (não simular).
