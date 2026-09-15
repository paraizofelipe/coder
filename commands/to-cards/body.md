Execute a skill `to-cards` para quebrar um spec em cards — fatias verticais tracer-bullet, cada uma declarando quem a bloqueia — gravados em `.coder/cards-<branch>/`. `$ARGUMENTS` é opcional: o caminho de um spec específico. Vazio, confirme com o usuário qual spec é o alvo.

Esta skill **não decide e não implementa**: recorta o que o spec já decidiu. E **não escreve nada** antes de o usuário aprovar a granularidade.

## Passos

### 1. Ancorar no spec e no código

- Ler o spec inteiro; `Fora do escopo` vincula tanto quanto as histórias
- Confirmar que as seams estão acordadas em `Decisões de teste`; card não inventa seam
- Explorar o repositório o necessário para nomear módulos e interfaces pelo vocabulário do domínio

### 2. Procurar o prefactor

- "Deixe a mudança fácil, depois faça a mudança fácil": existe preparação que torne as fatias seguintes triviais?
- Prefactor não muda comportamento e não entrega nada ao usuário — a suíte fica verde sem teste novo nem alterado
- Existindo, é o **card 01** e bloqueia os que dependem da estrutura nova. Não existindo, não inventar

### 3. Cortar as fatias e declarar as arestas

- Aplicar as regras de corte de `vertical-slices.md` da skill `implement`: bala traçante, uma seam por fatia, tamanho de uma janela de contexto nova, expandir–contrair para mudança larga
- Declarar **Bloqueado por** em cada card. Aresta é bloqueio real, nunca preferência de ordem — aresta inventada serializa trabalho que poderia correr em paralelo
- Caminho feliz primeiro; erro, estado vazio e caso-limite são cards próprios

### 4. Aprovar a granularidade

- Apresentar a quebra numerada com título, bloqueios, entrega, fronteira inicial e caminho mais longo
- Perguntar: a granularidade está certa? As arestas são bloqueios reais? Fundir ou dividir algum?
- **Iterar até aprovar. Sem aprovação, nenhum arquivo é escrito**

### 5. Publicar

- Gravar `.coder/cards-<branch>/NN-slug.md`, um card por arquivo, numerado em ordem de dependência
- Sem caminho de arquivo e sem snippet — o card pode esperar semanas na fila
- Card é definição; estado vive em `.coder/impl-<branch>.md` e só lá
- Diretório já existente para outro spec: parar e perguntar
