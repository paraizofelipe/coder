package main

import (
	"strings"
	"unicode/utf8"
)

// bannerArt é o nome em ASCII, na fonte figlet "ANSI Shadow". Fica como
// literal em vez de gerado em runtime: é conteúdo visual estável, e um
// gerador seria uma dependência a mais para produzir sempre o mesmo texto.
const bannerArt = ` ██████╗ █████╗ ████████╗██╗      ██████╗  ██████╗        █████╗  ██████╗ ███████╗███╗   ██╗████████╗███████╗
██╔════╝██╔══██╗╚══██╔══╝██║     ██╔═══██╗██╔════╝       ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝██╔════╝
██║     ███████║   ██║   ██║     ██║   ██║██║  ███╗█████╗███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   ███████╗
██║     ██╔══██║   ██║   ██║     ██║   ██║██║   ██║╚════╝██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   ╚════██║
╚██████╗██║  ██║   ██║   ███████╗╚██████╔╝╚██████╔╝      ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   ███████║
 ╚═════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝ ╚═════╝  ╚═════╝       ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚══════╝`

// bannerWidth é a largura em células da linha mais longa da arte. Calculada
// a partir dela, e não fixada à mão, para não envelhecer se a arte mudar.
var bannerWidth = larguraMaxima(bannerArt)

func larguraMaxima(texto string) int {
	maior := 0
	for _, linha := range strings.Split(texto, "\n") {
		if largura := utf8.RuneCountInString(linha); largura > maior {
			maior = largura
		}
	}
	return maior
}

// contorno diz se o caractere faz parte do traçado de linha dupla da fonte —
// ═ ║ e os cantos. Na ANSI Shadow é esse traçado que desenha a sombra da
// letra; o corpo é feito de blocos cheios.
func contorno(r rune) bool {
	switch r {
	case '═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬':
		return true
	}
	return false
}
