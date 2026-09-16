package main

import (
	"strings"
	"testing"
)

func TestContornoReconheceSoALinhaDupla(t *testing.T) {
	for _, r := range "═║╔╗╚╝╠╣╦╩╬" {
		if !contorno(r) {
			t.Errorf("%q faz parte do traçado e deveria ser contorno", r)
		}
	}
	for _, r := range "█ ABC─│┌┐" {
		if contorno(r) {
			t.Errorf("%q não é linha dupla e não deveria ser contorno", r)
		}
	}
}

func TestColorirArteSeparaCorpoDeContorno(t *testing.T) {
	azul, branco := colored.banner[0], colored.banner[1]

	got := (&ui{c: colored}).colorirArte("██╗═█")
	esperado := azul + "██" + branco + "╗═" + azul + "█" + colored.reset
	if got != esperado {
		t.Errorf("obtido  %q\nesperado %q", got, esperado)
	}
}

func TestColorirArteEmiteUmEscapePorTrocaDeCorNaoPorCaractere(t *testing.T) {
	got := (&ui{c: colored}).colorirArte("████████")
	if n := strings.Count(got, "\033[38;5;"); n != 1 {
		t.Errorf("oito caracteres da mesma cor geraram %d escapes, esperava 1", n)
	}
}

func TestColorirArteDevolveALinhaIntactaSemCor(t *testing.T) {
	const linha = " ██████╗ █████╗"
	if got := (&ui{}).colorirArte(linha); got != linha {
		t.Errorf("com cores desligadas a linha precisa sair intacta; obtido %q", got)
	}
}

func TestColorirArtePreservaTodosOsCaracteresDaArte(t *testing.T) {
	// A pintura não pode perder nem duplicar caractere: a arte é alinhada
	// por coluna, e um caractere a mais desmancha o desenho.
	for _, linha := range strings.Split(bannerArt, "\n") {
		colorida := (&ui{c: colored}).colorirArte(linha)
		limpa := colorida
		for _, escape := range []string{colored.banner[0], colored.banner[1], colored.reset} {
			limpa = strings.ReplaceAll(limpa, escape, "")
		}
		if limpa != linha {
			t.Errorf("a linha mudou ao ser colorida:\n antes:  %q\n depois: %q", linha, limpa)
		}
	}
}
