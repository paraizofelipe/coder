# Formato dos cards

Carregue no **passo 5**, ao gravar. O passo 4 não escreve nada.

## Onde os cards moram

`.coder/cards-<branch-safe>/NN-slug.md` — um diretório por branch, um arquivo por card.

`<branch-safe>` é a branch atual (`git rev-parse --abbrev-ref HEAD`) com `/` trocado por `-`, a mesma regra do plano e do andamento. Cards são **lidos de volta** pelo `/implement`, e por isso levam a branch no nome: `.coder/` está fora do versionamento, então um diretório fixo faria os cards de uma branch aparecerem para outra.

| Branch atual | Diretório |
|---|---|
| `feat/rate-limit` | `.coder/cards-feat-rate-limit/` |
| `fix/timeout-cotacao` | `.coder/cards-fix-timeout-cotacao/` |
| `main` | `.coder/cards-main/` |
| `HEAD` desanexado, ou fora de repositório Git | `.coder/cards-AAAAMMDD-HHMMSS/` |

`NN` começa em `01` e segue a **ordem de dependência**: bloqueador sempre com número menor que o bloqueado. `slug` é o título em kebab-case.

```text
.coder/cards-feat-rate-limit/
  01-extrair-limitador-para-interface.md
  02-limitar-por-chave-de-api.md
  03-responder-429-com-retry-after.md
```

### Quando o diretório já existe

| Situação | O que fazer |
|---|---|
| Mesmo spec | Pare e pergunte: os cards anteriores já podem estar em execução por outra pessoa |
| Spec diferente | **Pare e pergunte** antes de qualquer coisa. Sobrescrever apaga trabalho distribuído |

## Template

````markdown
# <NN>: <título curto, no imperativo>

> Spec: `.coder/spec-AAAAMMDD-HHMMSS.md` · Seam: `<seam acordada>` · <AAAA-MM-DD>

## O que entrega

O comportamento de ponta a ponta que este card torna possível, na perspectiva de quem usa.
Não é lista de camadas nem de arquivos.

## Bloqueado por

- `<NN>: <título>` — ou `Nenhum (pode começar imediatamente)`

## Critérios de aceite

- [ ] <cenário observável que prova a entrega>
- [ ] <caminho alternativo ou de erro, quando fizer parte deste card>

## Notas

<restrição herdada do spec, decisão codificada, ou `Nenhuma.`>
````

## Regras por seção

| Seção | Cuidado |
|---|---|
| Cabeçalho | Nomeia o spec e a seam. É o que permite alguém pegar o card semanas depois sem a conversa que o gerou |
| `O que entrega` | Comportamento observável de ponta a ponta. "Criar o model de X" não é entrega — é camada |
| `Bloqueado por` | Só bloqueio real. `Nenhum` é a resposta mais valiosa: significa que o card está na fronteira |
| `Critérios de aceite` | Verificáveis. "Funciona" e "testado" não são critérios |
| `Notas` | Onde entra a exceção de snippet, com a origem anotada |

## O que nunca vai para o card

| Item | Por quê |
|---|---|
| Caminho de arquivo | O card pode esperar semanas na fila; o caminho envelhece antes de alguém pegá-lo |
| Snippet de código | Mesma razão. Exceção única: trecho que codifica uma decisão com mais precisão que a prosa — máquina de estados, schema, shape de tipo |
| Estado de andamento | Card é definição. Estado vive em `.coder/impl-<branch>.md`, e só lá |
| Decisão que o spec não tomou | O card recorta o que foi decidido; não decide |
| Lista de camadas a construir | Isso é fatia horizontal disfarçada de card |

## Card de prefactor

Segue o mesmo template, com duas diferenças que precisam estar visíveis:

- `O que entrega` diz explicitamente que **não há mudança de comportamento** — a entrega é a estrutura que torna os cards seguintes triviais
- `Critérios de aceite` inclui: a suíte existente permanece verde, sem teste novo e sem teste alterado

É sempre o card `01`, e todo card que depende da estrutura nova o declara em `Bloqueado por`.
