# Mock

Carregue no **passo 2**, quando a fatia tocar uma fronteira externa.

## A linha

Mock só em **fronteira de sistema** — onde o seu código para e começa o de outra pessoa.

| Mocke | Não mocke |
|---|---|
| API externa (pagamento, e-mail, terceiros) | Suas próprias classes e módulos |
| Relógio, data, aleatoriedade, UUID | Colaborador interno do módulo sob teste |
| Sistema de arquivos, às vezes | Qualquer coisa que você controla |
| Banco de dados, às vezes — prefira um banco de teste real | A própria interface que está sendo testada |

Mockar colaborador interno faz o mock virar o objeto sob teste: o teste passa a provar que o seu dublê se comporta como você escreveu, o que é sempre verdade e nunca útil.

Quando bater a vontade de mockar algo interno, normalmente a seam está baixa demais. Suba a seam em vez de mockar.

## Projetar para ser mockável

### Injete a dependência externa

```typescript
// Fácil de mockar — a fronteira entra por parâmetro
function processarPagamento(pedido, clienteDePagamento) {
  return clienteDePagamento.cobrar(pedido.total);
}

// Difícil — a fronteira é criada por dentro, e o teste precisa interceptar o módulo
function processarPagamento(pedido) {
  const cliente = new ClienteStripe(process.env.STRIPE_KEY);
  return cliente.cobrar(pedido.total);
}
```

### Prefira interface no estilo SDK a um buscador genérico

```typescript
// BOM — cada operação é mockável isoladamente
const api = {
  buscarUsuario: (id) => fetch(`/usuarios/${id}`),
  buscarPedidos: (usuarioId) => fetch(`/usuarios/${usuarioId}/pedidos`),
  criarPedido: (dados) => fetch("/pedidos", { method: "POST", body: dados }),
};

// RUIM — o mock precisa de lógica condicional por endpoint
const api = {
  fetch: (endpoint, opcoes) => fetch(endpoint, opcoes),
};
```

Com o estilo SDK, cada dublê devolve uma forma só, o teste não carrega `if` no setup e fica evidente quais chamadas externas aquele teste exercita.

## Determinismo

Tempo, aleatoriedade e identificadores gerados são fronteira de sistema mesmo quando parecem detalhe. Não fixá-los produz teste que falha sozinho uma vez por mês, e teste intermitente é desativado pela equipe em poucas semanas.

Fixe o valor pela mesma via por onde ele entra no código — injeção ou utilitário do próprio projeto — nunca reimplementando o cálculo dentro do teste.

## Sinais de mock demais

- O setup do teste é maior que o comportamento verificado
- O teste quebra ao renomear uma função que o usuário nunca vê
- Ler o teste não revela que capacidade existe, só que colaboradores foram chamados
- Trocar a implementação por uma equivalente exige reescrever o dublê
