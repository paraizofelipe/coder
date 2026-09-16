package main

import (
	"fmt"
	"strings"
)

// leitores são os outros harnesses que varrem o destino de um harness, por
// escopo. Os quatro convergiram para o mesmo formato de skill e passaram a
// ler os diretórios uns dos outros de propósito — então escolher o destino
// não escolhe quem enxerga.
type leitores struct {
	global  []string
	projeto []string
}

// leitoresExtras foi medido, não deduzido da documentação: `opencode debug
// skill` e `copilot skill list` num container limpo por combinação, e `omp
// config list` para os defaults do OMP.
//
// A tabela existe porque a documentação de cada harness descreve só o que ele
// próprio lê. Ninguém publica a matriz inversa — que é a que importa para
// quem instala.
var leitoresExtras = map[string]leitores{
	// ~/.config/opencode e .opencode/ só o OpenCode varre. O OMP lê
	// .opencode/commands no projeto, mas não as skills: não há provider de
	// skills do OpenCode no `omp config list`.
	"opencode": {},

	// O caso que mais surpreende. O OpenCode auto-carrega
	// ~/.claude/skills como "external skill"; no projeto, somam-se o OMP
	// (skills.enableClaudeProject, ligado) e o Copilot, que aceita
	// .claude/skills no repositório.
	"claude": {
		global:  []string{"opencode"},
		projeto: []string{"opencode", "omp", "copilot"},
	},

	// ~/.agents é o diretório canônico do padrão Agent Skills, e por isso o
	// mais lido: OpenCode por auto-carga, Copilot por documentação própria.
	"omp": {
		global:  []string{"opencode", "copilot"},
		projeto: []string{"opencode", "copilot"},
	},

	// ~/.copilot e .github/ são os únicos destinos exclusivos dos quatro.
	"copilot": {},
}

func (h harness) leitores() []string {
	l := leitoresExtras[h.name]
	if h.projeto() {
		return l.projeto
	}
	return l.global
}

// sobreposicao descreve o alcance real de uma seleção, em duas listas
// separadas porque são duas decisões diferentes:
//
//   - alcance diz para quem mais o conteúdo vai. É informação, não defeito:
//     pode ser exatamente o que se quer.
//   - redundantes aponta alvo escolhido que já é alcançado pelo destino de
//     outro alvo escolhido. Essa é a acionável — cópia a mais não é neutra:
//     quando as duas divergem, qual delas o harness carrega passa a variar.
//
// Nenhuma das duas muda o que é instalado. A tabela descreve os defaults, e
// quem mexeu nos settings do OMP ou nas env vars do OpenCode tem outra
// matriz — que o instalador não tem como enxergar. Decidir por quem roda,
// com base em configuração que não se vê, seria pior do que informar.
func sobreposicao(targets []harness) (alcance, redundantes []string) {
	selecionado := make(map[string]bool, len(targets))
	for _, h := range targets {
		selecionado[h.name] = true
	}

	for _, h := range targets {
		todos := h.leitores()
		if len(todos) == 0 {
			continue
		}
		alcance = append(alcance, fmt.Sprintf("%s também é lido por: %s",
			h.skillsDir(), strings.Join(todos, ", ")))

		var jaCobertos []string
		for _, leitor := range todos {
			if selecionado[leitor] {
				jaCobertos = append(jaCobertos, leitor)
			}
		}
		if len(jaCobertos) > 0 {
			redundantes = append(redundantes, fmt.Sprintf(
				"%s já %s alcançado%s pelo destino de %s — instalar neles duplica o conteúdo",
				strings.Join(jaCobertos, ", "), verbo(len(jaCobertos)), plural(len(jaCobertos)), h.name))
		}
	}
	return alcance, redundantes
}

func verbo(n int) string {
	if n == 1 {
		return "é"
	}
	return "são"
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
