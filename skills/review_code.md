---
description: Skill compartilhada pelos subagentes viewer e reviewer. Cobre duas camadas de revisão: (1) qualidade técnica, padrões e testes (viewer) e (2) integridade com regras de negócio, boas práticas e segurança (reviewer). Nenhum código é versionado sem passar por ambas.
---

Você está executando a skill `review_code`. Esta skill é usada em dois momentos distintos do fluxo:

- Pelo `viewer`: revisão de qualidade técnica, padrões do projeto e cobertura de testes — logo após a implementação
- Pelo `reviewer`: revisão de integridade com regras de negócio, boas práticas e segurança — portão final antes do versionamento

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

## Camada 1 — Revisão técnica (viewer)

Execute esta camada quando acionado pelo `viewer`, logo após a implementação.

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

## Camada 2 — Revisão de negócio, boas práticas e segurança (reviewer)

Execute esta camada quando acionado pelo `reviewer`, como portão final antes do versionamento.

<checklist id="camada-2">
**Integridade com as regras de negócio**
- [ ] O comportamento implementado corresponde exatamente ao que foi solicitado?
- [ ] Todas as regras de negócio descritas na solicitação foram respeitadas?
- [ ] Há lógica que contradiz ou ignora alguma regra de negócio?
- [ ] Casos excepcionais previstos nas regras foram tratados?
- [ ] O fluxo de dados respeita as restrições de negócio (validações, limites, permissões)?
- [ ] Há efeitos colaterais que impactam outros processos de negócio?

**Boas práticas de desenvolvimento**
- [ ] O código segue o princípio de menor privilégio?
- [ ] Há separação adequada entre lógica de negócio e infraestrutura?
- [ ] Erros e exceções são tratados e comunicados corretamente?
- [ ] Não há lógica crítica de negócio duplicada em múltiplos lugares?
- [ ] Operações com efeitos colaterais são adequadamente isoladas e controladas?
- [ ] Configurações sensíveis não estão hardcoded?
- [ ] Logs não expõem informações sensíveis (dados pessoais, credenciais, tokens)?

**Segurança — OWASP Top 10 e práticas essenciais**
- [ ] **Injeção:** Risco de SQL injection, command injection ou similar?
- [ ] **Autenticação:** Mecanismos de autenticação e sessão sem bypass?
- [ ] **Exposição de dados:** Dados pessoais, financeiros ou confidenciais protegidos?
- [ ] **Controle de acesso:** Verificação de autorização antes de operações protegidas?
- [ ] **Configuração insegura:** Permissões excessivas, modo debug ou valores padrão inseguros?
- [ ] **XSS:** Entradas do usuário sanitizadas antes de serem renderizadas?
- [ ] **Deserialização insegura:** Dados externos validados antes de deserialização?
- [ ] **Componentes vulneráveis:** Uso de bibliotecas desatualizadas ou com CVEs conhecidos?
- [ ] **Logging insuficiente:** Eventos críticos de segurança são registrados?
- [ ] **SSRF:** Chamadas a URLs externas construídas com input do usuário?
- [ ] Tokens, chaves e segredos ausentes no código e nos logs?

**Preparação para produção**
- [ ] O código funcionaria corretamente em produção (não apenas em desenvolvimento)?
- [ ] Há dependências de variáveis de ambiente documentadas?
- [ ] Operações de migração têm plano de rollback?
- [ ] O impacto em performance foi considerado para volumes reais?
</checklist>

<criteria>
**APROVADO:** Nenhum problema crítico ou importante identificado. Pode prosseguir para a próxima etapa.

**APROVADO COM RESSALVAS:** Há sugestões ou pontos de atenção, mas nenhum bloqueia a continuidade. O `coder` decide se corrige agora ou registra como pendência.

**REPROVADO:** Há problemas críticos ou importantes que impedem a continuidade. O `coder` deve corrigir e submeter para nova revisão antes de prosseguir.
</criteria>

<output_format>
### Resultado geral
**[APROVADO / APROVADO COM RESSALVAS / REPROVADO]**

### Camada executada
- [ ] Camada 1 — Revisão técnica (viewer)
- [ ] Camada 2 — Revisão de negócio, boas práticas e segurança (reviewer)

### Problemas críticos (bloqueiam aprovação)
- [descrição, arquivo e linha, sugestão de correção]

### Problemas importantes (devem ser corrigidos)
- [descrição, arquivo e linha, sugestão de correção]

### Sugestões (melhorias desejáveis, não bloqueantes)
- [descrição e justificativa]

### Pontos positivos
- [o que foi bem implementado]

### Ações necessárias antes de prosseguir
- [ ] [ação específica]
</output_format>
