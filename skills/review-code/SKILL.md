---
name: review-code
description: Skill compartilhada pelos subagentes code_reviewer e business_reviewer. Cobre duas camadas de revisão — (1) qualidade técnica, padrões e testes (code_reviewer, logo após implementação) e (2) integridade com regras de negócio, boas práticas e segurança OWASP (business_reviewer, portão final antes do versionamento). Use quando o coder pedir revisão técnica ou revisão de negócio/segurança. Nenhum código é versionado sem passar por ambas. Devolve parecer APROVADO, APROVADO COM RESSALVAS ou REPROVADO.
---

Você está executando a skill `review-code`. Esta skill é usada em dois momentos distintos do fluxo:

- Pelo `code_reviewer`: revisão de qualidade técnica, padrões do projeto e cobertura de testes — logo após a implementação
- Pelo `business_reviewer`: revisão de integridade com regras de negócio, boas práticas e segurança — portão final antes do versionamento

Identifique em qual papel você está atuando e execute a revisão correspondente. Se acionado pelos dois agentes em sequência, execute ambas as camadas.

<context>
Para realizar uma revisão completa, você deve ter:
- O relatório do `analyzer` com padrões, convenções e arquitetura do projeto
- A descrição original da solicitação do usuário (incluindo regras de negócio)
- Todo o código criado ou alterado na implementação
- Os testes criados ou ajustados
- O resultado da execução dos testes
</context>

---

## Camada 1 — Revisão técnica (code_reviewer)

Execute esta camada quando acionado pelo `code_reviewer`, logo após a implementação.

<checklist id="camada-1">
**Corretude funcional**
- [ ] A implementação resolve exatamente o que foi solicitado pelo usuário?
- [ ] Há casos de borda não tratados?
- [ ] A lógica está correta e sem erros óbvios?
- [ ] O comportamento em situações de erro é adequado?

**Qualidade do código**
- [ ] O código é legível e compreensível sem comentários excessivos?
- [ ] Os nomes de variáveis, funções e classes são descritivos e claros?
- [ ] Há duplicação desnecessária que poderia ser eliminada?
- [ ] Cada função/método tem uma única responsabilidade clara?
- [ ] A complexidade está adequada ao problema (sem over-engineering)?
- [ ] Há código morto, comentado ou não utilizado?

**Aderência ao projeto**
- [ ] Segue a arquitetura e padrões identificados pelo `analyzer`?
- [ ] Usa as mesmas convenções de nomenclatura do restante do projeto?
- [ ] Utiliza bibliotecas e utilitários já disponíveis no projeto?
- [ ] O estilo de código é consistente com o existente?
- [ ] Os imports seguem o padrão do projeto?

**Cobertura de testes**
- [ ] Os testes cobrem os comportamentos alterados de forma significativa?
- [ ] Há testes para o caminho feliz e para cenários de erro?
- [ ] Os testes são determinísticos e não dependem de estado externo não controlado?
- [ ] Os testes fazem asserções relevantes (não apenas verificam que não lançou erro)?

**Segurança e riscos**
- [ ] Há vulnerabilidades introduzidas (injeção, exposição de dados, etc.)?
- [ ] Operações destrutivas têm proteção adequada?
- [ ] Dados sensíveis não estão sendo expostos em logs ou respostas?
- [ ] Há dependências circulares ou acoplamentos problemáticos?

**Impacto no sistema**
- [ ] A mudança pode causar regressões em outras partes do sistema?
- [ ] Há efeitos colaterais não intencionais?
- [ ] A performance foi considerada quando relevante?
</checklist>

---

## Camada 2 — Revisão de negócio, boas práticas e segurança (business_reviewer)

Execute esta camada quando acionado pelo `business_reviewer`, como portão final antes do versionamento.

O checklist completo (regras de negócio, boas práticas, OWASP Top 10 e preparação para produção) está em `references/owasp-checklist.md`. Carregue esse arquivo ao iniciar a Camada 2 e marque cada item antes de emitir o parecer.

<criteria>
**APROVADO:** Nenhum problema crítico ou importante identificado. Pode prosseguir para a próxima etapa.

**APROVADO COM RESSALVAS:** Há sugestões ou pontos de atenção, mas nenhum bloqueia a continuidade. O `coder` decide se corrige agora ou registra como pendência.

**REPROVADO:** Há problemas críticos ou importantes que impedem a continuidade. O `coder` deve corrigir e submeter para nova revisão antes de prosseguir.
</criteria>

<output_format>
### Resultado geral
**[APROVADO / APROVADO COM RESSALVAS / REPROVADO]**

### Camada executada
- [ ] Camada 1 — Revisão técnica (code_reviewer)
- [ ] Camada 2 — Revisão de negócio, boas práticas e segurança (business_reviewer)

Cada problema ou sugestão deve seguir obrigatoriamente este bloco:

```
📍 <path/do/arquivo> — linha <N>

**Atual:**
```<linguagem>
[trecho de código atual exato]
```

**Sugerido:**
```<linguagem>
[trecho de código como deve ficar]
```

**Motivo:** [explicação objetiva do problema e da correção proposta]
```

### Problemas críticos (bloqueiam aprovação)
[bloco acima para cada problema crítico]

### Problemas importantes (devem ser corrigidos)
[bloco acima para cada problema importante]

### Sugestões (melhorias desejáveis, não bloqueantes)
[bloco acima para cada sugestão]

### Pontos positivos
- [o que foi bem implementado]

### Ações necessárias antes de prosseguir
- [ ] [ação específica]
</output_format>
