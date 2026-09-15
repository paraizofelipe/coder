---
name: code-review
description: Revisa a diff entre `HEAD` e um ponto fixo em dois eixos independentes — aderência ao spec e aderência aos padrões do repositório — e reporta os dois lado a lado, sem fundi-los. Use ao revisar uma branch, um MR ou o trabalho em andamento, e no fim da skill `implement`.
---

<role>
Você está executando a skill `code-review`. Revisa mudanças em dois eixos que não se substituem:

- **Spec** — o código faz o que o spec pediu? Nem menos, nem mais.
- **Padrões** — o código segue as convenções documentadas deste repositório?

Uma mudança pode passar em um eixo e falhar no outro: código impecável que implementa a coisa errada passa em Padrões e falha em Spec; código que faz exatamente o pedido violando toda convenção falha em Padrões e passa em Spec. Reportar os dois separados é o que impede um de mascarar o outro.

Esta skill **não corrige**. Ela aponta, com evidência. A correção é decisão de quem implementa.
</role>

<context>
**Entrada** — um ponto fixo (commit, branch, tag ou merge-base) e, quando existir, o spec de origem.

Quando acionada pela skill `implement`, o ponto fixo é o merge-base com a branch base e o spec é o `.coder/spec-AAAAMMDD-HHMMSS.md` que guiou a implementação. Quando acionada direto pelo usuário sem ponto fixo, **pergunte** — é a única pergunta obrigatória desta skill.

**Saída** — um relatório com os dois eixos em seções separadas. Nenhum arquivo é escrito e nenhuma linha de código é alterada.

**Paralelismo** — quando o harness oferecer subagentes, rode os dois eixos em paralelo, cada um recebendo **apenas** a sua reference e nada do outro eixo: a separação de contexto é o que os mantém independentes. Sem subagentes, rode em sequência, um eixo inteiro por vez, sem consultar os achados do outro no meio.
</context>

<workflow>

### 1. Fixe o ponto de comparação
- Resolva o que o usuário indicou: `git rev-parse <ponto-fixo>`. Ref inválida falha aqui, não dentro dos eixos.
- Capture a diff uma vez, com três pontos, para comparar contra o merge-base: `git diff <ponto-fixo>...HEAD`. Liste também os commits: `git log <ponto-fixo>..HEAD --oneline`.
- Diff vazia encerra a revisão com essa constatação. Não invente achado sobre código que não mudou.

### 2. Localize o spec
Nesta ordem, parando no primeiro que existir:

1. O card de `.coder/cards-<branch>/` que guiou o trabalho, mais o spec que o cabeçalho dele nomeia
2. `.coder/spec-AAAAMMDD-HHMMSS.md` citado pela skill `implement` ou pelo arquivo de andamento da branch
3. Um caminho que o usuário tenha passado como argumento
4. `.coder/plan-<branch>.md` da branch atual
5. Referência a issue ou ticket nas mensagens de commit

Não encontrando nada, pergunte ao usuário onde está o spec. Se ele disser que não há, **pule o eixo Spec** e registre isso no relatório — não o substitua por suposição sobre a intenção.

### 3. Localize os padrões
- Procure o que o repositório documenta sobre como escrever código: `AGENTS.md`, `CLAUDE.md`, `CONTRIBUTING.md`, `CODING_STANDARDS.md`, arquivos em `.agents/rules/`, ADRs.
- Ignore o que a ferramenta já garante: formatador, linter e typecheck configurados no projeto. Apontar o que o CI reprova sozinho é ruído.

### 4. Rode os dois eixos
- **Spec** — siga `references/spec-axis.md`.
- **Padrões** — siga `references/standards-axis.md`, que traz também a linha de base de smells aplicada quando o repositório documenta pouco.
- Em subagente, passe a reference do eixo **inteira**: ele não tem outro acesso a ela.

### 5. Agregue sem reordenar
- Apresente os dois eixos em seções próprias, na ordem `Spec` e depois `Padrões`.
- **Não funda, não reordene e não escolha um vencedor entre eixos.** Reordenar é exatamente o que a separação existe para impedir.
- Encerre com o total por eixo e o pior achado **dentro de cada eixo**.

</workflow>

<rules>
- **Não corrija nada.** Esta skill lê a diff e escreve um relatório. Nenhum arquivo é alterado, nenhum comando de Git que mude estado é executado.
- **Todo achado carrega evidência**: o trecho da diff e a linha do spec ou a regra documentada que ele contraria. Achado sem citação é opinião.
- **Distinga violação de julgamento.** Padrão documentado descumprido é violação; smell é heurística rotulada ("possível inveja de dados"), nunca veredito.
- **O repositório manda.** Convenção documentada no projeto vence a linha de base de smells, inclusive quando endossa o que a linha de base apontaria.
- **Não aponte o que a ferramenta já pega.** Formatação, import não usado e erro de tipo são do linter e do typecheck.
- **Escopo indevido é achado do eixo Spec**, com o mesmo peso de requisito faltando. Código a mais também é divergência.
- Não revise arquivo fora da diff, por mais tentador que seja o que houver ali.
</rules>

<checklist>
- [ ] O ponto fixo foi resolvido e a diff é não vazia.
- [ ] A origem do spec foi identificada, ou a ausência dele foi registrada.
- [ ] As fontes de padrão do repositório foram localizadas antes de aplicar a linha de base.
- [ ] Os dois eixos rodaram sem ver os achados um do outro.
- [ ] Cada achado cita a diff e a regra ou a linha do spec que o sustenta.
- [ ] Violações e julgamentos estão rotulados de forma distinta.
- [ ] Nada que o linter, o formatador ou o typecheck já pegam foi reportado.
- [ ] Os eixos foram apresentados separados, sem fusão nem reordenação.
- [ ] Nenhum arquivo foi alterado durante a revisão.
</checklist>

<output_format>
```text
Revisão de <ponto-fixo>...HEAD — <N arquivos>, <N commits>.

## Spec
Origem: `<caminho do spec>` — ou `nenhum spec disponível; eixo não executado`

1. [faltando | escopo indevido | implementado errado] <achado>
   Spec: "<linha citada>"
   Diff: <arquivo> — <trecho ou descrição do hunk>

## Padrões

1. [violação | julgamento] <achado>
   Regra: <arquivo + regra> — ou <nome do smell>
   Diff: <arquivo> — <trecho ou descrição do hunk>

---
Spec: <N> achados — pior: <resumo em 1 linha, ou "nenhum">
Padrões: <N> achados — pior: <resumo em 1 linha, ou "nenhum">
```

Nenhum achado em um eixo imprime `Nenhum achado.` sob o título dele. O eixo não desaparece do relatório.
</output_format>
