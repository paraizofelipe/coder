# Checklist completo — Camada 2 (business_reviewer)

Checklist detalhado para o `business_reviewer`. Carregar este arquivo apenas durante a Camada 2, como portão final antes do versionamento.

## Integridade com as regras de negócio

- [ ] O comportamento implementado corresponde exatamente ao que foi solicitado?
- [ ] Todas as regras de negócio descritas na solicitação foram respeitadas?
- [ ] Há lógica que contradiz ou ignora alguma regra de negócio?
- [ ] Casos excepcionais previstos nas regras foram tratados?
- [ ] O fluxo de dados respeita as restrições de negócio (validações, limites, permissões)?
- [ ] Há efeitos colaterais que impactam outros processos de negócio?

## Boas práticas de desenvolvimento

- [ ] O código segue o princípio de menor privilégio?
- [ ] Há separação adequada entre lógica de negócio e infraestrutura?
- [ ] Erros e exceções são tratados e comunicados corretamente?
- [ ] Não há lógica crítica de negócio duplicada em múltiplos lugares?
- [ ] Operações com efeitos colaterais são adequadamente isoladas e controladas?
- [ ] Configurações sensíveis não estão hardcoded?
- [ ] Logs não expõem informações sensíveis (dados pessoais, credenciais, tokens)?

## Segurança — OWASP Top 10 e práticas essenciais

- [ ] **Injeção:** risco de SQL injection, command injection ou similar?
- [ ] **Autenticação:** mecanismos de autenticação e sessão sem bypass?
- [ ] **Exposição de dados:** dados pessoais, financeiros ou confidenciais protegidos?
- [ ] **Controle de acesso:** verificação de autorização antes de operações protegidas?
- [ ] **Configuração insegura:** permissões excessivas, modo debug ou valores padrão inseguros?
- [ ] **XSS:** entradas do usuário sanitizadas antes de serem renderizadas?
- [ ] **Deserialização insegura:** dados externos validados antes de deserialização?
- [ ] **Componentes vulneráveis:** uso de bibliotecas desatualizadas ou com CVEs conhecidos?
- [ ] **Logging insuficiente:** eventos críticos de segurança são registrados?
- [ ] **SSRF:** chamadas a URLs externas construídas com input do usuário?
- [ ] Tokens, chaves e segredos ausentes no código e nos logs?

## Preparação para produção

- [ ] O código funcionaria corretamente em produção (não apenas em desenvolvimento)?
- [ ] Há dependências de variáveis de ambiente documentadas?
- [ ] Operações de migração têm plano de rollback?
- [ ] O impacto em performance foi considerado para volumes reais?
