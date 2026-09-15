---
name: to-memory
description: Envia os artefatos de `.coder/` da branch atual, crus e sob confirmação, para o serviço de memória de longo prazo (Hindsight), e mantém as páginas de conhecimento — perguntas permanentes que o servidor re-responde a cada consolidação. Use depois de `implement`, quando o trabalho atravessou sessões. Sem serviço de memória disponível, encerra sem erro.
---

<role>
Você está executando a skill `to-memory`. Leva o que o fluxo já registrou em `.coder/` para a memória de longo prazo, onde ele passa a atravessar projetos e sessões.

Esta skill não sintetiza e não cura: **envia o conteúdo bruto**. O serviço quer o artefato inteiro para sintetizar do lado dele; um resumo enviado no lugar do documento destrói justamente a matéria-prima da síntese.

O que ela cura são as **perguntas**, não as respostas. Página de conhecimento guarda uma pergunta permanente; o conteúdo é reconstruído pelo servidor conforme a evidência se acumula.
</role>

<context>
**Entrada** — os artefatos da branch atual em `.coder/`: o plano, o spec e o andamento. O mapa de qual vai como quê está em `references/artifact-map.md`.

**Saída** — documentos no bank de memória, mais as páginas de conhecimento criadas ou atualizadas. Nenhum arquivo local é escrito ou alterado.

**Capacidade necessária** — uma ferramenta de memória de longo prazo exposta ao agente atual: no Claude Code, o conjunto `agent_knowledge_*` do plugin Hindsight (`ingest_file`, `ingest`, `recall`, `get_current_bank`, `list_pages`, `create_page`, `update_page`); em outros harnesses, o equivalente. **Verifique a lista de ferramentas do agente; não chame às cegas.**

Sem essa capacidade, **encerre pelo `<output_format>` de indisponibilidade**. Não é erro, não é falha e não pede retentativa: os artefatos em `.coder/` seguem sendo o registro do fluxo, legíveis por qualquer harness sem configuração nenhuma.

**Quando usar** — quando o trabalho atravessou sessões. Se coube em uma janela de contexto, a captura automática do harness já o reteve, e o envio explícito é custo sem retorno.
</context>

<workflow>

### 1. Verifique a capacidade e o destino
- Confirme que a ferramenta de memória está disponível ao agente atual. Ausente, encerre por `<output_format>`.
- Descubra o bank corrente (`get_current_bank`) **antes** de montar a confirmação. O bank é escolhido pelo diretório de trabalho, então o mesmo fluxo em repositórios diferentes escreve em banks diferentes.
- Bank errado é erro de fronteira de dados, não detalhe de configuração: artefato de trabalho de cliente em bank pessoal não tem como ser desfeito. O nome do bank vai **na** mensagem de confirmação, não antes dela.

### 2. Selecione os artefatos
- Resolva os arquivos da branch atual pelo mapa de `references/artifact-map.md`, que traz o que cada um vira no bank e por que os nomes já resolvem a identidade dos documentos.
- Artefato de outra branch não entra. Se o usuário quiser enviar um específico, ele nomeia; você não varre `.coder/` inteiro por conta própria.
- Verifique o que já existe no bank (`recall` pelo nome do artefato) para saber se o envio substitui ou cria. Substituição é o caso normal do plano e do andamento.

### 3. Varra antes de enviar
- Siga `references/redaction.md`. Os três artefatos têm perfis de risco diferentes, e o plano é o mais arriscado dos três.
- Achado sensível **para o envio daquele artefato** e volta ao usuário. Não edite o arquivo para "limpar" e enviar assim mesmo: o artefato é registro do fluxo, e alterá-lo para caber na memória corrompe o registro.

### 4. Confirme e envie
- Apresente: bank de destino, arquivos, o que cada um substitui ou cria, e o resultado da varredura. Aguarde confirmação explícita.
- Envie o **conteúdo bruto e completo**, via a ferramenta de ingestão de arquivo. Nunca resuma, recorte ou reformate antes.
- Uma tentativa por artefato. Falha de rede ou de serviço encerra com o que foi enviado até ali, reportado como está — sem repetir.

### 5. Mantenha as páginas de conhecimento
- Consulte `references/knowledge-pages.md` para o conjunto de perguntas permanentes e a redação de cada `source_query`.
- Liste as páginas existentes antes de criar qualquer uma. Página que já cobre a pergunta é **atualizada**, nunca duplicada.
- Não crie página por trabalho, por feature ou por branch. Isso é documento, e documento é o passo 4. Página é pergunta permanente, e elas são poucas e estáveis.

### 6. Reporte e pare
- Resumo no formato de `<output_format>`. Não despeje o conteúdo enviado.

</workflow>

<rules>
- **Nunca resuma antes de enviar.** A ferramenta de ingestão pede conteúdo bruto; resumo no lugar do documento destrói o que o servidor usaria para sintetizar.
- **Envio é irreversível.** A superfície exposta ao agente não tem exclusão de documento — só de página. Reenviar com o mesmo identificador substitui; nada remove. Por isso a confirmação é explícita e nomeia o bank.
- **Uma tentativa, sem repetição.** Ferramenta ausente, serviço fora do ar e bank sem resultado são equivalentes para o fluxo: siga, ou encerre, sem retentativa.
- **Ausência de memória não é falha.** Encerrar sem enviar é um desfecho válido e deve ser dito com todas as letras, não escondido nem tratado como erro.
- **Não edite artefato para enviá-lo.** O que está em `.coder/` é registro do fluxo; ajustá-lo para caber na memória falsifica o registro.
- **Página é pergunta, não resposta.** Não escreva conteúdo curado em página, e não crie uma por feature.
- Nenhum arquivo local é criado ou alterado por esta skill, nem operação de Git é executada.
</rules>

<checklist>
- [ ] A disponibilidade da ferramenta foi verificada na lista do agente, não por tentativa.
- [ ] O bank de destino foi consultado e apareceu na mensagem de confirmação.
- [ ] Só artefatos da branch atual, ou os que o usuário nomeou, entraram na seleção.
- [ ] A varredura de conteúdo sensível rodou antes do envio, com o plano recebendo o escrutínio maior.
- [ ] O conteúdo foi enviado bruto e completo, sem resumo, recorte ou reformatação.
- [ ] O usuário confirmou explicitamente antes do primeiro envio.
- [ ] As páginas existentes foram listadas antes de criar qualquer uma.
- [ ] Nenhuma página foi criada por feature, branch ou trabalho específico.
- [ ] Nenhum arquivo local foi alterado.
</checklist>

<output_format>
Com serviço disponível, após o envio:

```text
Memória atualizada — bank `<bank>`.

- Documentos enviados: <N>
  - `<arquivo>` — <substituiu o anterior | novo>
- Páginas: <N criadas, N atualizadas, N inalteradas>
- Varredura: <sem achados | N artefatos retidos e por quê>
- Não enviado: <o que ficou de fora e por quê, ou "nada">

<próximo passo recomendado, em 1 linha>
```

Sem serviço de memória disponível, responda **apenas** com:

```text
Memória de longo prazo indisponível neste harness — nada foi enviado.

Os artefatos do fluxo seguem em `.coder/`, que continua sendo o registro:
<lista dos arquivos da branch atual>

Isso não é falha do fluxo. Para habilitar, instale o plugin de memória no harness em uso.
```
</output_format>
