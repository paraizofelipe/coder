# Fatias verticais

Carregue no **passo 2**, ao transformar o spec em ordem de execução.

## O que é uma fatia

Uma fatia atravessa a seam acordada de ponta a ponta e entrega **comportamento observável**. Ao fim dela o sistema funciona; não fica pela metade esperando a próxima.

É a diferença entre fatiar **na vertical** e **na horizontal**:

| Corte | Como se parece | Por que falha ou funciona |
|---|---|---|
| **Vertical** — use este | "Usuário com carrinho válido conclui a compra" | Cada fatia é uma bala traçante: prova a ligação inteira cedo e responde ao que a anterior ensinou |
| **Horizontal** — evite | "Primeiro todos os modelos, depois todos os serviços, depois todas as rotas" | Nada é verificável até o fim. Você constrói a *forma* do sistema, não o comportamento dele, e descobre os erros de ligação no último dia |

## Como cortar

| Regra | Razão |
|---|---|
| **Uma fatia, uma seam** | Se a fatia precisa de duas seams para ser observada, ela é grande demais ou a seam está baixa demais |
| **Comece pela que mais reduz incerteza** | A primeira fatia é reconhecimento: ela revela o que o spec não podia saber sobre o código real |
| **Ordene por dependência real, não por camada** | Dependência real é "B não é observável sem A". Camada é organização de arquivo, não de comportamento |
| **Fatia que não cabe em um ciclo é duas fatias** | Se o vermelho → verde não fecha sem trocar de assunto no meio, o corte está errado |
| **Caminho feliz primeiro, exceções como fatias próprias** | Erro, estado vazio e caso-limite são histórias do spec; cada um vira a sua fatia |

## O que a primeira fatia deve provar

A ligação mais arriscada do trabalho, em sua forma mais fina. Não a mais fácil, não a mais "fundamental" na arquitetura: a que, se estiver errada, invalida o resto.

Feita ela, as seguintes engrossam a mesma ligação em vez de abrir frentes novas.

## Raio de alcance

Antes de começar, estime quanto do código cada fatia toca:

| Alcance | Sinal | Consequência no corte |
|---|---|---|
| **Contido** | Um módulo, sem mudança de interface pública | Fatia normal |
| **Largo** | Muitos chamadores, contrato compartilhado, schema | Use expandir–contrair (abaixo); nunca em uma fatia só |
| **Irreversível** | Migração de dados, mudança de contrato publicado, exclusão | Fatia própria, com verificação própria, e confirmação do usuário antes |

## Expandir–contrair, para mudança larga

Quando muitos chamadores dependem do que vai mudar, três fatias em vez de uma:

1. **Expandir** — introduza o novo ao lado do antigo. Nada quebra, porque nada deixou de existir
2. **Migrar** — mova os chamadores, em fatias por grupo, cada uma verificável
3. **Contrair** — remova o antigo, agora sem chamadores

Tentar isso em uma fatia só transforma toda a suíte em vermelho ao mesmo tempo — e aí o vermelho não diz mais nada sobre o que está errado.

## Apresentação

Liste as fatias em até 10 linhas, na ordem de execução:

```text
Fatias desta implementação:

1. <comportamento observável> — seam: <nome>
2. <comportamento observável> — seam: <nome>

<em 1 linha: o que a fatia 1 prova e por que ela vem primeiro>
```

Siga sem aguardar resposta. Interrompa apenas se o corte revelar decisão que o spec não tomou.
