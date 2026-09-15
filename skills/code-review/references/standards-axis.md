# Eixo Padrões

Carregue no **passo 4**. Quando rodar em subagente, este arquivo vai inteiro no contexto dele — é o único acesso que ele tem a estas instruções e à linha de base abaixo.

## A pergunta

O código segue as convenções **documentadas neste repositório**?

Este eixo não julga se a mudança era a certa. Isso é do eixo Spec.

## Insumos

- A diff (`git diff <ponto-fixo>...HEAD`) e a lista de commits
- Os arquivos de padrão encontrados no repositório: `AGENTS.md`, `CLAUDE.md`, `CONTRIBUTING.md`, `CODING_STANDARDS.md`, `.agents/rules/`, ADRs

## Duas fontes, uma hierarquia

1. **O que o repositório documenta** — a autoridade. Descumprimento é **violação**, e a citação é o arquivo mais a regra.
2. **A linha de base de smells** abaixo — vale mesmo quando o repositório não documenta nada. Todo item dela é **julgamento**, nunca violação.

Duas regras amarram isso:

- **O repositório vence.** Onde a convenção documentada endossa algo que a linha de base apontaria, suprima o smell.
- **Smell é heurística rotulada** — "possível inveja de dados" —, não veredito. Rotule como julgamento e siga.

E, antes de qualquer das duas: **pule o que a ferramenta já garante.** Formatador, linter e typecheck configurados no projeto não precisam de revisor.

## Linha de base de smells

Fowler, *Refactoring*, cap. 3. Cada item se lê *o que é* → *como resolver*. Aplique contra a diff, não contra o repositório inteiro.

| Smell | O que é | Como resolver |
|---|---|---|
| **Nome misterioso** | Função, variável ou tipo cujo nome não revela o que faz ou guarda | Renomeie; se nenhum nome honesto aparece, o desenho está confuso |
| **Código duplicado** | A mesma forma lógica em mais de um hunk ou arquivo da mudança | Extraia a forma comum e chame dos dois lados |
| **Inveja de dados** | Método que acessa mais os dados de outro objeto do que os próprios | Mova o método para junto dos dados que ele inveja |
| **Aglomerado de dados** | Os mesmos poucos campos ou parâmetros andando sempre juntos — um tipo querendo nascer | Junte-os em um tipo e passe esse tipo |
| **Obsessão por primitivo** | Primitivo ou string no lugar de um conceito de domínio que merece tipo próprio | Dê ao conceito um tipo pequeno e próprio |
| **Switches repetidos** | O mesmo `switch`/cascata de `if` sobre o mesmo tipo reaparecendo na mudança | Troque por polimorfismo, ou um mapa compartilhado pelos dois pontos |
| **Cirurgia com espingarda** | Uma mudança lógica obriga a editar muitos arquivos espalhados | Reúna no mesmo módulo o que muda junto |
| **Mudança divergente** | Um arquivo ou módulo editado por vários motivos sem relação | Separe, para que cada módulo mude por um motivo só |
| **Generalidade especulativa** | Abstração, parâmetro ou gancho para necessidade que o spec não tem | Remova; volte a inline até uma necessidade real aparecer |
| **Cadeia de mensagens** | Navegação longa `a.b().c().d()` que o chamador não deveria conhecer | Esconda o percurso atrás de um método no primeiro objeto |
| **Intermediário** | Classe ou função que basicamente só delega adiante | Corte e chame o destino real |
| **Herança recusada** | Subclasse que ignora ou sobrescreve quase tudo o que herda | Troque herança por composição |

## Segurança e dados sensíveis

Independem de o repositório documentar. Aponte sempre, como **violação**:

- Credencial, token, chave ou string de conexão na diff
- Dado pessoal em log, mensagem de erro ou fixture
- Entrada externa concatenada em query, comando ou caminho de arquivo
- Erro engolido em silêncio: `catch` vazio, exceção descartada sem tratamento nem registro

## Formato do achado

```text
1. [violação | julgamento] <o que se observa>
   Regra: <arquivo + regra documentada> — ou <nome do smell>
   Diff: <arquivo> — <trecho ou descrição do hunk>
```

Abaixo de 400 palavras no total. Violações antes de julgamentos; dentro de cada grupo, o de maior consequência primeiro.
