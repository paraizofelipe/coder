package main

import (
	"fmt"
	"strings"
)

// destino é um lugar do disco com artefatos deste instalador dentro.
//
// O status não guarda registro do que instalou: varre o disco e responde
// pelo que encontra. Um arquivo de estado seria uma segunda verdade sobre a
// mesma coisa, e duas verdades divergem — basta alguém apagar um diretório
// à mão para o registro passar a mentir sobre o que existe.
type destino struct {
	harness  string
	escopo   escopo
	dir      string
	skills   []string
	commands []string
}

// varrerDestinos percorre os quatro harnesses nos escopos possíveis e devolve
// só os destinos ocupados. Raiz vazia significa "não estamos num
// repositório", e aí o escopo de projeto nem é procurado.
//
// Conta apenas o que o manifesto declara: skill de terceiro no mesmo
// diretório é legítima e não é assunto deste instalador.
func varrerDestinos(man manifest, raiz string) []destino {
	escopos := []escopo{{}}
	if raiz != "" {
		escopos = append(escopos, escopo{raiz: raiz})
	}

	var out []destino
	for _, e := range escopos {
		for _, h := range comEscopo(harnesses, e) {
			d := destino{harness: h.name, escopo: e, dir: h.base()}
			for _, nome := range man.Skills {
				if h.hasSkill(nome) {
					d.skills = append(d.skills, nome)
				}
			}
			for _, nome := range man.Commands {
				if h.hasCommand(nome) {
					d.commands = append(d.commands, nome)
				}
			}
			if len(d.skills) > 0 || len(d.commands) > 0 {
				out = append(out, d)
			}
		}
	}
	return out
}

// le diz se um harness varre um destino — o dele próprio, ou o de outro que
// a tabela de leitores reporta.
func le(leitor string, d destino) bool {
	if d.harness == leitor {
		return true
	}
	dono, ok := lookupHarness(d.harness)
	if !ok {
		return false
	}
	for _, outro := range comEscopo([]harness{dono}, d.escopo)[0].leitores() {
		if outro == leitor {
			return true
		}
	}
	return false
}

// duplicatas aponta skill que um mesmo harness enxerga em mais de um
// destino. É o achado que dá razão ao comando: com duas cópias, qual delas o
// harness carrega não é escolha de ninguém — foi medido variando entre
// execuções do mesmo binário, no mesmo disco.
//
// Cópias idênticas são inofensivas. O problema nasce quando uma é atualizada
// e a outra não, e aí não há como saber qual está no ar.
func duplicatas(destinos []destino) []string {
	var avisos []string
	for _, h := range harnesses {
		// A comparação é dentro de um escopo, nunca entre os dois: os quatro
		// harnesses definem que o projeto tem precedência sobre o diretório
		// de usuário, e a mesma skill nos dois é sobreposição prevista.
		//
		// O indefinido é outro — duas cópias no mesmo escopo, como
		// ~/.claude e ~/.agents, que nenhuma regra ordena.
		for _, noProjeto := range []bool{false, true} {
			var lidos []destino
			for _, d := range destinos {
				if d.escopo.projeto() == noProjeto && le(h.name, d) {
					lidos = append(lidos, d)
				}
			}
			if len(lidos) < 2 {
				continue
			}
			// Uma linha por harness, não por skill: sete linhas iguais com
			// nomes diferentes contam a mesma história uma vez só.
			repetidas := repetidasEntre(lidos)
			if len(repetidas) == 0 {
				continue
			}
			caminhos := make([]string, 0, len(lidos))
			for _, d := range lidos {
				caminhos = append(caminhos, d.dir)
			}
			avisos = append(avisos, fmt.Sprintf(
				"%s (escopo %s) enxerga %d skill(s) em mais de um destino (%s) — qual cópia vence não é definido: %s",
				h.name, rotulo(escopo{raiz: raizSe(noProjeto)}), len(repetidas),
				strings.Join(caminhos, ", "), strings.Join(repetidas, ", ")))
		}
	}
	return avisos
}

// raizSe devolve uma raiz simbólica só para o rótulo do escopo sair certo na
// mensagem — o caminho de verdade já está em cada destino listado.
func raizSe(noProjeto bool) string {
	if noProjeto {
		return "."
	}
	return ""
}

func repetidasEntre(destinos []destino) []string {
	vezes := map[string]int{}
	for _, d := range destinos {
		for _, nome := range d.skills {
			vezes[nome]++
		}
	}
	// A ordem sai do manifesto, não do mapa: saída estável entre execuções.
	var out []string
	for _, d := range destinos {
		for _, nome := range d.skills {
			if vezes[nome] > 1 {
				out = append(out, nome)
				vezes[nome] = 0 // já contada
			}
		}
	}
	return out
}

// rotulo é o escopo em uma palavra. O resumo() repetiria a raiz que a linha
// seguinte já mostra por inteiro.
func rotulo(e escopo) string {
	if e.projeto() {
		return "projeto"
	}
	return "global"
}
