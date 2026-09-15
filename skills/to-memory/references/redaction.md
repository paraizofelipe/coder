# Varredura antes do envio

Carregue no **passo 3**, antes de montar a confirmação.

## Por que esta varredura é mais severa que a do commit

O commit pode ser desfeito antes do push. O envio à memória, não: a superfície exposta ao agente **não tem exclusão de documento** — reenviar com o mesmo identificador substitui, mas nada remove. E o bank é compartilhado entre harnesses e projetos, então o que entra ali reaparece em conversas que não têm relação com este trabalho.

## Perfil de risco por artefato

Os três não são igualmente perigosos, porque os formatos do fluxo já filtram coisas diferentes:

| Artefato | Risco | Por quê |
|---|---|---|
| **Plano** | **Alto** | `Fatos confirmados` exige a fonte junto do fato. É exatamente ali que aterrissa log colado, URL interna, nome de cliente, ID de ticket e trecho de payload real |
| Andamento | Médio | O formato proíbe diff e conteúdo de arquivo, mas `Verificação do conjunto` carrega saída real de comando, e `Decisões tomadas` pode citar caminho de infraestrutura |
| Spec | **Baixo** | O formato já proíbe caminho de arquivo e snippet, salvo a exceção de decisão codificada. É o mais seguro dos três por construção |

Comece pelo plano. Se houver um único achado no conjunto, é lá que ele estará.

## O que procurar

| Categoria | Exemplos |
|---|---|
| Credencial | Token, chave de API, senha, string de conexão, cabeçalho de autorização |
| Dado pessoal | Nome de pessoa física, e-mail, CPF, CNPJ, telefone, endereço — inclusive em exemplo ou fixture |
| Identificação de cliente | Nome de cliente, razão social, identificador de conta, número de contrato |
| Topologia interna | Hostname interno, IP privado, caminho de bucket, nome de cluster, endpoint não público |
| Conteúdo capturado | Trecho de log de produção, payload real de requisição, linha de banco de dados |
| Referência interna | ID de ticket, URL de board, link de wiki interna — quando o bank de destino é pessoal |

A última linha depende do destino: identificador interno de trabalho é esperado no bank da organização e fora de lugar no bank pessoal. É por isso que o bank é resolvido no passo 1, antes da varredura.

## O que fazer com um achado

**Retenha o artefato e volte ao usuário.** Nomeie o que encontrou e onde, sem reproduzir o valor sensível na resposta.

Não faça, em nenhuma hipótese:

| Não faça | Por quê |
|---|---|
| Editar o artefato para "limpar" e enviar | O que está em `.coder/` é registro do fluxo. Alterá-lo para caber na memória falsifica o registro |
| Enviar uma versão recortada | A ferramenta pede conteúdo completo; um recorte é um resumo com outro nome |
| Enviar assim mesmo por ser "provavelmente inofensivo" | A avaliação do que é sensível é do usuário, e o envio não tem volta |

Retenção parcial é desfecho normal: envie os artefatos limpos, retenha o que tem achado, e diga claramente no relatório qual ficou de fora e por quê.

## Quando o usuário decide enviar assim mesmo

Decisão dele, e vale registrá-la. Confirme uma vez, nomeando o achado e o bank de destino, e siga. Não repita o alerta depois da confirmação.
