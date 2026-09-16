package main

import (
	"bytes"
	"strings"
	"testing"
)

func alvos(e escopo, nomes ...string) []harness {
	sel := make(map[string]bool, len(nomes))
	for _, n := range nomes {
		sel[n] = true
	}
	return comEscopo(canonical(sel), e)
}

// O alcance é o que a seleção realmente atinge, não o que ela pede. Instalar
// para claude num projeto entrega as skills a quatro harnesses.
func TestAlcanceNomeiaQuemMaisLeODestino(t *testing.T) {
	casos := []struct {
		nome      string
		targets   []harness
		contem    []string
		naoContem []string
	}{
		{
			nome:    "claude no projeto chega a três a mais",
			targets: alvos(escopo{raiz: "/p"}, "claude"),
			contem:  []string{"claude", "opencode", "omp", "copilot"},
		},
		{
			nome:      "claude no global só vaza para o opencode",
			targets:   alvos(escopo{}, "claude"),
			contem:    []string{"opencode"},
			naoContem: []string{"omp", "copilot"},
		},
		{
			nome:    "omp vaza nos dois escopos",
			targets: alvos(escopo{}, "omp"),
			contem:  []string{"opencode", "copilot"},
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			alcance, _ := sobreposicao(caso.targets)
			junto := strings.Join(alcance, "\n")
			if junto == "" {
				t.Fatal("esperava ao menos uma linha de alcance")
			}
			for _, termo := range caso.contem {
				if !strings.Contains(junto, termo) {
					t.Errorf("alcance não menciona %q:\n%s", termo, junto)
				}
			}
			for _, termo := range caso.naoContem {
				if strings.Contains(junto, termo) {
					t.Errorf("alcance menciona %q, que não lê esse destino:\n%s", termo, junto)
				}
			}
		})
	}
}

// Destino exclusivo não rende aviso: linha em toda instalação vira ruído e
// para de destacar as que importam.
func TestSemSobreposicaoNaoHaAviso(t *testing.T) {
	for _, e := range []escopo{{}, {raiz: "/p"}} {
		alcance, redundantes := sobreposicao(alvos(e, "copilot"))
		if len(alcance) != 0 || len(redundantes) != 0 {
			t.Errorf("escopo %s: o copilot não é lido por ninguém; avisos: %v %v",
				e.resumo(), alcance, redundantes)
		}
	}
}

// A redundância é a informação acionável: o alvo escolhido já é alcançado
// pelo destino de outro alvo escolhido, e instalar nele só duplica.
func TestRedundanciaApontaOAlvoJaAlcancado(t *testing.T) {
	_, redundantes := sobreposicao(alvos(escopo{raiz: "/p"}, "claude", "omp"))
	junto := strings.Join(redundantes, "\n")
	if !strings.Contains(junto, "omp") || !strings.Contains(junto, "claude") {
		t.Errorf("esperava que omp fosse apontado como alcançado por claude:\n%s", junto)
	}

	// No global a história é outra: o OMP não lê ~/.claude.
	_, redundantes = sobreposicao(alvos(escopo{}, "claude", "omp"))
	if len(redundantes) != 0 {
		t.Errorf("no escopo global o omp não é alcançado por claude: %v", redundantes)
	}
}

// Invariante, não lista fixa: harness novo sem entrada na tabela passaria a
// prometer exclusividade que ninguém mediu.
func TestTodoHarnessTemEntradaNaTabelaDeLeitores(t *testing.T) {
	for _, h := range harnesses {
		if _, ok := leitoresExtras[h.name]; !ok {
			t.Errorf("%s não tem entrada em leitoresExtras", h.name)
		}
	}
	for nome, l := range leitoresExtras {
		if _, ok := lookupHarness(nome); !ok {
			t.Errorf("leitoresExtras cita harness inexistente: %q", nome)
		}
		for _, leitor := range append(append([]string{}, l.global...), l.projeto...) {
			if _, ok := lookupHarness(leitor); !ok {
				t.Errorf("%s: leitor inexistente %q", nome, leitor)
			}
			if leitor == nome {
				t.Errorf("%s: um harness não é leitor extra de si mesmo", nome)
			}
		}
	}
}

// A tela só fala quando tem o que dizer: aviso em toda instalação viraria
// ruído, e é a ausência dele que passa a informar exclusividade.
func TestNoticeDeSobreposicaoCalaQuandoNaoHaNada(t *testing.T) {
	var buf bytes.Buffer
	u := &ui{out: &buf}

	u.sobreposicaoNotice(alvos(escopo{}, "copilot"))
	if buf.Len() != 0 {
		t.Errorf("destino exclusivo não deveria imprimir nada; imprimiu:\n%s", buf.String())
	}

	u.sobreposicaoNotice(alvos(escopo{raiz: "/p"}, "claude", "omp"))
	saida := buf.String()
	for _, termo := range []string{"opencode", "copilot", "duplica"} {
		if !strings.Contains(saida, termo) {
			t.Errorf("a saída não menciona %q:\n%s", termo, saida)
		}
	}
}
