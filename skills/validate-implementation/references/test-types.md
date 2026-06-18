# Tipos de teste

Referência consultada ao montar o plano (Passo 3 da skill). Para cada modificação/regra de negócio afetada, escolher um ou mais tipos abaixo. A lista não é exaustiva — adapte ao que faz sentido para o fluxo.

## Smoke

- **Definição:** verificação rápida e superficial de que o serviço/fluxo "sobe e responde" — sanidade básica, não profundidade.
- **Quando aplicar:** logo após confirmar acessos, para garantir que a aplicação e suas dependências estão de pé antes dos testes mais caros.
- **Exemplo:** `GET /health` retorna 200; o consumer conecta na fila; a aplicação está `Healthy` no ArgoCD.
- **Classe típica:** leitura.

## Black-box

- **Definição:** validação do comportamento observável de uma funcionalidade pela sua interface (entrada → saída), sem olhar a implementação interna.
- **Quando aplicar:** para cada regra de negócio alterada, exercitar entradas representativas (caminho feliz, entrada inválida, borda) e conferir a saída esperada.
- **Exemplo:** enviar um payload válido a um endpoint e validar o corpo/HTTP status; consultar um registro e conferir os campos calculados pela regra.
- **Classe típica:** leitura (consultas) ou mutação (criação/atualização via API).

## E2E (ponta a ponta)

- **Definição:** exercita o fluxo completo da regra de negócio atravessando múltiplos serviços (ex.: API → processamento → persistência → efeito observável).
- **Quando aplicar:** quando a modificação afeta um fluxo que cruza componentes e o valor está em validar a integração real, não as partes isoladas.
- **Exemplo:** disparar um evento/requisição, aguardar o processamento e confirmar o resultado no banco e/ou nos logs do worker.
- **Classe típica:** frequentemente mutação (cria dados ao longo do fluxo) — exige confirmação por item.

## Regressão

- **Definição:** confirma que comportamentos pré-existentes **não** quebraram com a modificação.
- **Quando aplicar:** sempre que a mudança toca código compartilhado ou um fluxo crítico que já funcionava; cobrir os cenários adjacentes ao que mudou.
- **Exemplo:** revalidar o caminho feliz antigo e casos vizinhos que dependem do mesmo módulo/endpoint alterado.
- **Classe típica:** leitura, quando possível; mutação se o cenário antigo envolvia escrita.

## Orientações gerais

- Preferir `leitura` sempre que a regra puder ser validada sem mutar estado.
- Marcar explicitamente cada caso como `leitura` ou `mutação` — isso determina se a execução exige confirmação por item (Passo 9).
- Em HML, mutações são aceitáveis sob confirmação; em PROD, evitar mutação e exigir confirmação reforçada para qualquer execução.
