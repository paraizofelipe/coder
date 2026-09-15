Execute a skill `to-memory` para enviar os artefatos de `.coder/` da branch atual ao serviço de memória de longo prazo e manter as páginas de conhecimento. `$ARGUMENTS` é opcional: um artefato específico a enviar. Vazio, use os artefatos da branch atual.

Esta skill **não sintetiza**: ela envia o conteúdo bruto. O que ela cura são as perguntas permanentes, não as respostas.

Sem ferramenta de memória disponível no harness, **encerre sem erro** — os artefatos em `.coder/` seguem sendo o registro do fluxo.

## Passos

### 1. Verificar a capacidade e o destino

- Conferir na lista de ferramentas do agente se há memória de longo prazo (`agent_knowledge_*` do Hindsight, ou equivalente). Não chamar às cegas
- Ausente: encerrar pelo formato de indisponibilidade. Não é falha e não pede retentativa
- Descobrir o bank corrente antes de montar a confirmação — ele é escolhido pelo diretório, então muda de repositório para repositório

### 2. Selecionar os artefatos

- Plano e andamento da branch atual **substituem** o documento anterior; cada spec **cria** um novo
- Artefato de outra branch não entra, salvo se o usuário nomeá-lo

### 3. Varrer antes de enviar

- Procurar credencial, dado pessoal, identificação de cliente, topologia interna, log ou payload capturado
- O plano é o mais arriscado dos três; o spec, o mais seguro
- Achado retém aquele artefato e volta ao usuário. **Nunca editar o artefato para limpá-lo** — isso falsifica o registro do fluxo

### 4. Confirmar e enviar

- Apresentar bank de destino, arquivos, o que substitui ou cria, e o resultado da varredura; aguardar confirmação explícita
- Enviar o conteúdo **bruto e completo**. Nunca resumir, recortar ou reformatar
- Uma tentativa por artefato; falha encerra com o que foi enviado, reportado como está

### 5. Manter as páginas de conhecimento

- Listar as páginas existentes antes de criar qualquer uma; página que já cobre a pergunta é atualizada, nunca duplicada
- Página guarda uma pergunta permanente, não conteúdo curado. Não criar página por feature, branch ou trabalho

### 6. Reportar

- Documentos enviados, páginas tocadas, resultado da varredura e o que ficou de fora
- Não despejar o conteúdo enviado na resposta
