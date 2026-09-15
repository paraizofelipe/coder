# Seams

Carregue no **passo 3**, antes de propor a confirmação ao usuário.

## O que é

Uma **seam** é a fronteira pública onde o comportamento é observado sem alcançar o interior da implementação. É onde o teste vive.

Acordar as seams antes de escrever o spec é o que faz o documento valer: a seam acordada vira o contrato que quem implementa, quem testa e quem revisa vão usar depois. Uma seam que ninguém acordou vira retrabalho na revisão.

## Heurísticas de escolha

| Regra | Razão |
|---|---|
| **Prefira a seam existente** | Uma fronteira já testada no projeto tem infra, fixtures e convenção prontas |
| **Escolha a mais alta possível** | Quanto mais alta, mais comportamento real ela observa e menos ela trava refatoração interna |
| **Menos é melhor — o ideal é uma** | Cada seam adicional é mais superfície para manter e mais acoplamento entre teste e estrutura |
| **Seam nova só no ponto mais alto viável** | Se precisar criar uma, proponha-a onde ela ainda observa comportamento, não onde é conveniente |

## Onde procurar candidatas

- Rota ou handler HTTP; resolver de GraphQL
- Comando de CLI, com entrada e saída observáveis
- Função ou classe exportada de um módulo — a interface pública, não o interno
- Publisher/consumer de evento ou mensagem de fila
- Adapter de persistência ou de serviço externo, quando o contrato dele é o que está em jogo
- Componente de UI pela interação do usuário, não pelo estado interno

## Anti-padrões

| Anti-padrão | Por que falha |
|---|---|
| Testar detalhe de implementação | O teste quebra em toda refatoração e não prova comportamento |
| Uma seam por camada | Multiplica testes que provam a mesma coisa em alturas diferentes |
| Uma seam por arquivo | Confunde unidade de código com unidade de comportamento |
| Mockar fronteira interna do próprio módulo | O mock passa a ser o que está sendo testado |
| Adiar a escolha para a implementação | A seam deixa de ser decisão acordada e vira acidente |

## Formato de apresentação

Apresente ao usuário em até 10 linhas e aguarde confirmação:

```text
Seams propostas para este spec:

1. <nome da seam> — [existente | nova]
   Observa: <comportamento que se torna verificável nessa fronteira>
   Por quê: <evidência: prior art no projeto, altura, ou ausência de alternativa>

<justificativa em 1 linha para a quantidade proposta>

Confirma essas seams ou quer ajustar antes de eu escrever o spec?
```

Se o usuário ajustar, registre o motivo — ele entra no spec junto da seam, em `Decisões de teste`.
