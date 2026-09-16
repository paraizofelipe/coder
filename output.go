package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/x/term"
)

type palette struct {
	red, yellow, green, cyan, bold, reset string

	// banner tem duas entradas: a cor do corpo da letra e a do contorno de
	// linha dupla. Vazio desliga a cor da arte.
	banner []string
}

var colored = palette{
	red: "\033[0;31m", yellow: "\033[1;33m", green: "\033[0;32m",
	cyan: "\033[0;36m", bold: "\033[1m", reset: "\033[0m",
	banner: []string{
		"\033[38;5;33m",  // corpo: azul
		"\033[38;5;231m", // contorno: branco
	},
}

// ui concentra a saída e as perguntas. ttyIn/ttyOut existem separados de out
// porque o formulário precisa do terminal mesmo quando a saída do binário
// está sendo redirecionada para um arquivo ou um pipe.
type ui struct {
	out         io.Writer
	ttyIn       *os.File
	ttyOut      *os.File
	interactive bool
	c           palette
}

// newUI abre o terminal. Preferir /dev/tty a os.Stdin mantém as perguntas
// funcionando quando a entrada padrão é um pipe — o mesmo motivo pelo qual o
// install.sh lia de /dev/tty.
//
// Cada candidato passa por term.IsTerminal: um descritor pode ser abrível e
// ainda assim não ser terminal. /dev/null é o caso que importa, porque é um
// character device, é o stdin de qualquer processo sem console, e aceitá-lo
// faria o formulário abrir para ninguém e falhar no meio da instalação.
func newUI(out io.Writer) *ui {
	u := &ui{out: out, c: colored}
	if os.Getenv("NO_COLOR") != "" {
		u.c = palette{}
	}
	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		if term.IsTerminal(tty.Fd()) {
			u.ttyIn, u.ttyOut, u.interactive = tty, tty, true
			return u
		}
		_ = tty.Close()
	}
	if term.IsTerminal(os.Stdin.Fd()) {
		u.ttyIn, u.ttyOut, u.interactive = os.Stdin, os.Stderr, true
	}
	return u
}

func (u *ui) close() {
	if u.ttyIn != nil && u.ttyIn != os.Stdin {
		_ = u.ttyIn.Close()
	}
}

func (u *ui) bold(s string) string { return u.c.bold + s + u.c.reset }

func (u *ui) printf(format string, args ...any) {
	fmt.Fprintf(u.out, format+"\n", args...)
}

func (u *ui) blank()                  { fmt.Fprintln(u.out) }
func (u *ui) info(f string, a ...any) { u.tagged(u.c.cyan, "[info]", "  ", f, a...) }
func (u *ui) ok(f string, a ...any)   { u.tagged(u.c.green, "[ok]", "    ", f, a...) }
func (u *ui) warn(f string, a ...any) { u.tagged(u.c.yellow, "[warn]", "  ", f, a...) }
func (u *ui) tagged(color, tag, pad, f string, a ...any) {
	u.printf("%s%s%s%s%s", color+u.c.bold, tag, u.c.reset, pad, fmt.Sprintf(f, a...))
}

func (u *ui) item(name string) { u.printf("  %s", u.bold(name)) }
func (u *ui) skipped()         { u.printf("        %s↳ pulado%s", u.c.yellow, u.c.reset) }
func (u *ui) installed()       { u.printf("        %s↳ instalado%s", u.c.green, u.c.reset) }

// simulated é o par de installed no modo --dry-run. O verbo diferente é
// proposital: quem lê o log depois precisa distinguir os dois sem contexto.
func (u *ui) simulated(f string, a ...any) {
	u.printf("        %s↳ simulado (%s)%s", u.c.cyan, fmt.Sprintf(f, a...), u.c.reset)
}

// dryRunNotice abre e fecha a execução simulada. Aparece duas vezes de
// propósito: uma instalação longa rola a tela, e o aviso do topo some.
func (u *ui) dryRunNotice() {
	u.printf("  %s▲ SIMULAÇÃO — nada será gravado em disco%s", u.c.yellow, u.c.reset)
	u.blank()
}

// terminalWidth devolve a largura do terminal, ou 0 quando não há como medir.
func (u *ui) terminalWidth() int {
	if u.ttyOut == nil {
		return 0
	}
	largura, _, err := term.GetSize(u.ttyOut.Fd())
	if err != nil {
		return 0
	}
	return largura
}

// colorirArte pinta a linha em dois tons: corpo azul, contorno branco. Emite
// um escape só quando a cor muda, e não um por caractere — a saída ficaria
// dez vezes maior sem nenhuma diferença visual.
func (u *ui) colorirArte(linha string) string {
	if len(u.c.banner) < 2 {
		return linha
	}
	corpo, traco := u.c.banner[0], u.c.banner[1]

	var out strings.Builder
	atual := ""
	for _, r := range linha {
		desejada := corpo
		if contorno(r) {
			desejada = traco
		}
		if desejada != atual {
			out.WriteString(desejada)
			atual = desejada
		}
		out.WriteRune(r)
	}
	if atual != "" {
		out.WriteString(u.c.reset)
	}
	return out.String()
}

// printArt só desenha o nome quando ele cabe inteiro. Arte que quebra linha
// em terminal estreito vira ruído ilegível — pior do que não ter banner.
// Largura desconhecida conta como suficiente: quem não tem terminal está
// lendo um log, onde a quebra não desalinha nada.
func (u *ui) printArt() {
	if largura := u.terminalWidth(); largura > 0 && largura < bannerWidth+2 {
		return
	}
	for _, linha := range strings.Split(bannerArt, "\n") {
		u.printf("  %s", u.colorirArte(linha))
	}
	u.blank()
}

// subtitle encurta quando a linha inteira não cabe, pelo mesmo motivo da
// arte: legenda quebrada em duas linhas desalinha o cabeçalho.
func (u *ui) subtitle() string {
	const completa = "instalador do fluxo planning → to-spec → to-cards → implement"
	if largura := u.terminalWidth(); largura > 0 && largura < utf8.RuneCountInString(completa)+2 {
		return "instalador de skills e commands"
	}
	return completa
}

func (u *ui) banner() {
	u.blank()
	u.printArt()
	u.printf("  %s", u.bold(u.subtitle()))
	u.blank()
}

func (u *ui) summary(targets []harness, dryRun bool) {
	u.blank()
	if dryRun {
		u.ok("Simulação concluída. Nenhum arquivo foi criado ou alterado.")
	} else {
		u.ok("Instalação concluída.")
	}
	u.blank()
	for _, h := range targets {
		u.printf("  %s %s", u.bold("•"), u.bold(h.name))
		u.printf("      skills:   %s", h.skillsDir())
		// Linha de commands só para quem os recebe: "commands:" seguido de
		// nada é a única forma de o resumo mentir sobre o que foi instalado.
		if h.commands {
			u.printf("      commands: %s", h.commandsDir())
		}
	}
	u.blank()
	if dryRun {
		u.dryRunNotice()
		u.printf("  Rode sem --dry-run para aplicar.")
	} else {
		u.printf("  Reinicie o harness para carregar o que foi instalado.")
	}
	u.blank()
}

// sobreposicaoNotice conta o alcance real da seleção antes de instalar. Fica
// em silêncio quando não há sobreposição — aviso em toda execução viraria
// ruído, e é a ausência dele que passa a dizer "este destino é exclusivo".
//
// O alcance sai como informação e a redundância como aviso, porque são
// coisas diferentes: alcançar mais harnesses pode ser o que se quer; gravar
// a mesma skill duas vezes onde um só harness vai lê-las, não.
func (u *ui) sobreposicaoNotice(targets []harness) {
	alcance, redundantes := sobreposicao(targets)
	if len(alcance) == 0 && len(redundantes) == 0 {
		return
	}
	for _, linha := range alcance {
		u.info("%s", linha)
	}
	for _, linha := range redundantes {
		u.warn("%s", linha)
	}
	u.blank()
}

// statusReport relata o que a varredura encontrou. Destino vazio não vira
// linha: o comando descreve o que existe, e listar os oito caminhos
// possíveis afogaria o que importa.
func (u *ui) statusReport(destinos []destino, avisos []string) {
	if len(destinos) == 0 {
		u.info("Nenhum artefato deste instalador foi encontrado.")
		u.blank()
		return
	}
	for _, d := range destinos {
		u.printf("  %s %s · %s", u.bold("•"), u.bold(d.harness), rotulo(d.escopo))
		u.printf("      %s", d.dir)
		u.printf("      skills: %d   commands: %d", len(d.skills), len(d.commands))
	}
	u.blank()
	for _, aviso := range avisos {
		u.warn("%s", aviso)
	}
	if len(avisos) > 0 {
		u.printf("        Reinstale nos dois destinos, ou remova um: cópias que divergem")
		u.printf("        transformam a precedência em sorteio.")
		u.blank()
	}
}
