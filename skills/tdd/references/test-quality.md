# O que distingue um bom teste

Carregue no **passo 2**, ao decidir o que asserir.

## O critério

Um bom teste observa comportamento por uma interface pública. O código interno pode mudar inteiro; o teste não. Ele se lê como especificação — "usuário com carrinho válido conclui a compra" diz que capacidade existe — e sobrevive a refatoração porque não sabe nada sobre estrutura.

| | Bom teste | Teste ruim |
|---|---|---|
| Observa | Comportamento pela interface pública | Estrutura interna |
| Nome descreve | **O quê** | **Como** |
| Em refatoração sem mudança de comportamento | Continua verde | Quebra |
| Asserção | Uma ideia lógica | Várias, ou nenhuma de verdade |
| Valor esperado | Fonte independente | Recalculado como o código calcula |

## Bom

```typescript
test("usuário com carrinho válido conclui a compra", async () => {
  const carrinho = criarCarrinho();
  carrinho.adicionar(produto);

  const resultado = await finalizarCompra(carrinho, meioDePagamento);

  expect(resultado.status).toBe("confirmada");
});
```

Usa só a interface pública, descreve uma capacidade, e nada nele depende de como `finalizarCompra` foi escrita por dentro.

## Anti-padrões

### Acoplado à implementação

```typescript
// RUIM — testa como, não o quê
test("finalizarCompra chama pagamento.processar", async () => {
  const pagamento = jest.mock(servicoDePagamento);
  await finalizarCompra(carrinho, meio);
  expect(pagamento.processar).toHaveBeenCalledWith(carrinho.total);
});
```

Sinais: mock de colaborador interno, teste de método privado, asserção sobre contagem ou ordem de chamadas, nome que descreve o caminho do código.

### Verificação por canal lateral

```typescript
// RUIM — contorna a interface para conferir
test("criarUsuario grava no banco", async () => {
  await criarUsuario({ nome: "Alice" });
  const linha = await db.query("SELECT * FROM usuarios WHERE nome = ?", ["Alice"]);
  expect(linha).toBeDefined();
});

// BOM — confere pela mesma interface que o usuário usaria
test("usuário criado fica recuperável", async () => {
  const usuario = await criarUsuario({ nome: "Alice" });
  const recuperado = await buscarUsuario(usuario.id);
  expect(recuperado.nome).toBe("Alice");
});
```

Conferir pelo banco congela o schema dentro do teste: trocar a forma de persistir quebra o teste sem que nenhum comportamento tenha mudado.

### Tautológico

```typescript
// RUIM — o esperado é recalculado como o código calcula; passa por construção
test("calcularTotal soma os itens", () => {
  const itens = [{ preco: 10 }, { preco: 5 }];
  const esperado = itens.reduce((soma, i) => soma + i.preco, 0);
  expect(calcularTotal(itens)).toBe(esperado);
});

// BOM — o esperado é um literal independente
test("calcularTotal soma os itens", () => {
  expect(calcularTotal([{ preco: 10 }, { preco: 5 }])).toBe(15);
});
```

Vale para o snapshot derivado à mão do mesmo jeito e para a constante comparada consigo mesma: se o teste nunca pode discordar do código, ele não verifica nada.

### Fatiamento horizontal

Escrever todos os testes primeiro e só depois toda a implementação verifica comportamento **imaginado**. Você testa a forma das coisas em vez do que o usuário faz, os testes ficam insensíveis a mudança real, e a estrutura de teste é fixada antes de você entender a implementação.

O oposto é a fatia vertical: um teste → uma implementação → o que isso ensinou decide o próximo.

## Antes de dar o ciclo por bom

- Trocar a implementação por outra, com o mesmo comportamento, mantém o teste verde?
- O nome diz o que passa a ser possível, sem citar função interna?
- O valor esperado veio de fora do código sob teste?
- Alguém que não escreveu o código entende, pelo teste, que capacidade existe?
