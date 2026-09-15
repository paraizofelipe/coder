# Árvore de decisões

Carregue no **passo 2**, ao modelar a árvore. Cada dimensão abaixo é um ramo possível: cubra as aplicáveis, ignore as que a solicitação não toca. Uma dimensão que você não abriu e não descartou conscientemente é uma decisão assumida em silêncio.

## Ordem: o que libera o quê

A fronteira é o conjunto de decisões cujos pré-requisitos já estão resolvidos. Estas dependências são as que mais aparecem:

| Esta decisão | Só faz sentido depois de |
|---|---|
| Regras e limites | Resultado de negócio |
| Dados e estados | Fluxos |
| Contratos e compatibilidade | Dados e estados |
| Migração | Dados e estados + contratos |
| Rollout e reversão | Migração |
| Critérios de aceite | Fluxos + regras |

Perguntar sobre migração antes de fechar o modelo de dados produz uma resposta que você vai ter que jogar fora.

## As dimensões

### 1. Resultado de negócio
**Decide:** o que muda para o negócio quando isto existir, e como se sabe que mudou.
- Qual número, comportamento ou capacidade é diferente depois da entrega?
- Se nada for entregue, qual o custo concreto de não fazer?

*Ficou implícito quando:* o plano descreve o que será construído, mas não o que passa a ser possível.

### 2. Usuários e permissões
**Decide:** quem executa a ação, quem vê o resultado e quem não pode nem uma coisa nem outra.
- Existe mais de um perfil? Eles veem a mesma coisa?
- O que acontece quando alguém sem permissão tenta: erro, silêncio ou ausência da opção?

*Ficou implícito quando:* a solicitação diz "o usuário" no singular sem nunca nomear qual.

### 3. Fluxos principais e alternativos
**Decide:** o caminho feliz e os caminhos que saem dele.
- Qual a sequência exata do caminho que funciona?
- Onde o fluxo pode ser abandonado no meio, e o que fica no sistema quando isso acontece?
- Existe caminho de volta (cancelar, editar, desfazer)?

*Ficou implícito quando:* só o caminho feliz foi descrito.

### 4. Regras e limites
**Decide:** o que é permitido, o que é proibido e quais são os valores de fronteira.
- Qual o mínimo e o máximo de cada entrada? O que acontece exatamente no limite?
- Existe regra que depende de tempo, de estado anterior ou de quem está agindo?
- Duas regras podem conflitar? Qual vence?

*Ficou implícito quando:* há adjetivos ("válido", "razoável", "recente") sem número ou condição.

### 5. Dados e estados
**Decide:** o que é persistido, em quais estados uma entidade pode estar e quais transições existem.
- Quais estados existem, e quais transições são proibidas?
- O dado é criado, atualizado ou versionado? Pode ser apagado?
- Alguma informação é derivada em vez de armazenada?

*Ficou implícito quando:* o plano fala em "salvar" sem dizer o que é imutável depois.

### 6. Contratos e compatibilidade
**Decide:** o que outros sistemas consomem e o que não pode quebrar.
- Quem consome esta interface hoje? Existe consumidor fora do seu controle?
- A mudança é compatível para trás? Se não, quem precisa mudar junto?
- O contrato é versionado? Qual a política de depreciação?

*Ficou implícito quando:* a mudança de payload é tratada como detalhe interno.

### 7. Falhas
**Decide:** o que o sistema faz quando algo dá errado.
- Qual dependência externa pode falhar, e o que o usuário vê quando ela falha?
- A operação é idempotente? O que acontece se for repetida?
- Falha parcial deixa o sistema em estado consistente?

*Ficou implícito quando:* só o sucesso e o "erro genérico" foram considerados.

### 8. Segurança e privacidade
**Decide:** o que é protegido, de quem, e o que não pode vazar em log nem em resposta.
- Há dado pessoal ou sensível envolvido? Ele aparece em log, erro ou telemetria?
- A autorização é verificada no servidor ou só escondida na interface?
- Entrada de terceiro é validada antes de virar consulta, comando ou HTML?

*Ficou implícito quando:* a discussão tratou apenas de autenticação, nunca de autorização.

### 9. Migração
**Decide:** o que fazer com o que já existe.
- Existe dado no formato antigo? Ele é convertido, mantido em paralelo ou descartado?
- A migração pode rodar com o sistema no ar?
- O que acontece com registros que não se encaixam no novo formato?

*Ficou implícito quando:* o plano descreve o estado final sem citar o estado atual.

### 10. Observabilidade
**Decide:** como se descobre que quebrou, antes do usuário avisar.
- Que evento precisa ser registrado para investigar um problema depois?
- Existe métrica que prova que a entrega funcionou em produção?
- O que deveria disparar alerta?

*Ficou implícito quando:* não há como distinguir "ninguém usou" de "todo mundo falhou".

### 11. Desempenho
**Decide:** quanto é rápido o bastante, com qual volume.
- Qual o volume esperado hoje e em um ano?
- Qual latência é aceitável, e medida em qual percentil?
- Existe operação que cresce com o número de registros?

*Ficou implícito quando:* apareceu a palavra "rápido" sem número.

### 12. Acessibilidade
**Decide:** quem consegue usar além do caminho com mouse e visão.
- A ação é alcançável e operável só pelo teclado?
- O estado (erro, carregando, sucesso) é comunicado além da cor?
- Conteúdo que muda sem recarregar é anunciado?

*Ficou implícito quando:* a interface foi descrita só visualmente.

### 13. Rollout
**Decide:** como isto chega em produção.
- Vai para todos de uma vez ou é gradual?
- Depende de flag, de configuração ou de ordem entre serviços?
- Precisa de coordenação com outro time ou com uma janela específica?

*Ficou implícito quando:* o plano termina no merge.

### 14. Reversão
**Decide:** o que fazer quando der errado depois de publicado.
- Dá para desligar sem deploy?
- O que acontece com dado criado pela versão nova se a antiga voltar?
- Existe ponto sem volta? Onde exatamente?

*Ficou implícito quando:* a resposta para "e se der errado" é "a gente corrige".

### 15. Testes e critérios de aceite
**Decide:** o que prova que está pronto.
- Qual cenário observável, executável por outra pessoa, prova cada entrega?
- Quais casos negativos e de borda precisam estar cobertos?
- O que **não** será testado, e por quê?

*Ficou implícito quando:* o critério é "funciona".

## Palavras que escondem decisão

Nenhuma destas encerra um ramo. Cada uma exige a pergunta da direita antes de virar registro:

| Palavra | Pergunte |
|---|---|
| rápido | Qual latência, em qual percentil, com qual volume? |
| seguro | Protegido contra quem, fazendo o quê? |
| intuitivo | Qual ação o usuário consegue completar sem instrução? |
| pronto | Qual cenário observável prova isso? |
| suportar | O que exatamente o sistema faz quando recebe isso? |
| escalável | Até qual volume, e o que quebra primeiro depois dele? |
| simples | Simples para quem: quem usa, quem mantém ou quem integra? |
| robusto | Resiste a qual falha específica, e como se comporta nela? |
