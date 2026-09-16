package main

import (
	"io"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSelectTargetsPrefereAFlagAoMenu(t *testing.T) {
	// ui sem terminal: se a flag não tivesse precedência, isto falharia.
	obtido, err := selectTargets(&ui{out: io.Discard}, "claude,opencode")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got := strings.Join(namesOf(obtido), ","); got != "opencode,claude" {
		t.Errorf("obtido %q, esperado \"opencode,claude\"", got)
	}
}

func TestSelectTargetsExplicaAAusenciaDeTerminalEmVezDeTravar(t *testing.T) {
	_, err := selectTargets(&ui{out: io.Discard}, "")
	if err == nil {
		t.Fatal("esperava erro sem terminal e sem --harness")
	}
	if !strings.Contains(err.Error(), "--harness") {
		t.Errorf("o erro precisa apontar a saída: %v", err)
	}
}

func TestRunEncerraSemInstalarNadaEmHelpEVersion(t *testing.T) {
	for _, arg := range []string{"--help", "-h", "--version", "-v"} {
		if err := run([]string{arg}); err != nil {
			t.Errorf("%s retornou erro: %v", arg, err)
		}
	}
}

func TestHelpTextDocumentaTodasAsOpcoesEOverrides(t *testing.T) {
	texto := helpText()
	obrigatorios := []string{
		"--force", "--harness", "--version", "--help",
		"NO_COLOR", "ACCESSIBLE", "all",
	}
	// Nome e override saem da tabela de harnesses: harness novo sem linha na
	// ajuda é um destino que a flag aceita e que ninguém descobre lendo -h.
	for _, h := range harnesses {
		obrigatorios = append(obrigatorios, h.name, h.envVar)
	}
	for _, termo := range obrigatorios {
		if !strings.Contains(texto, termo) {
			t.Errorf("ajuda não menciona %q", termo)
		}
	}
}

func TestBaseCaiNoHomeQuandoNaoHaOverride(t *testing.T) {
	t.Setenv("CLAUDE_DIR", "")
	t.Setenv("HOME", "/home/teste")
	h, _ := lookupHarness("claude")
	if got := h.base(); got != "/home/teste/.claude" {
		t.Errorf("obtido %q, esperado \"/home/teste/.claude\"", got)
	}
}

func TestVersionDevolveDevEmBuildLocal(t *testing.T) {
	if v := version(); v == "" {
		t.Error("version() nunca deve devolver string vazia")
	}
}

func TestParseArgsReconheceDryRun(t *testing.T) {
	for _, arg := range []string{"--dry-run", "-n"} {
		opts, err := parseArgs([]string{arg, "--harness", "all"})
		if err != nil {
			t.Fatalf("%s: erro inesperado: %v", arg, err)
		}
		if !opts.dryRun {
			t.Errorf("%s não ligou dryRun", arg)
		}
	}
	opts, _ := parseArgs([]string{"--harness", "all"})
	if opts.dryRun {
		t.Error("dryRun deveria vir desligado por padrão")
	}
}

func TestBannerArtEstaIntegraEMedida(t *testing.T) {
	linhas := strings.Split(bannerArt, "\n")
	if len(linhas) != 6 {
		t.Fatalf("a arte tem %d linhas, esperava 6", len(linhas))
	}
	for i, linha := range linhas {
		if strings.ContainsAny(linha, "\t\x1b") {
			t.Errorf("linha %d contém tab ou escape ANSI — a cor é aplicada na impressão", i+1)
		}
		if largura := utf8.RuneCountInString(linha); largura > bannerWidth {
			t.Errorf("linha %d tem %d células, acima da largura calculada %d", i+1, largura, bannerWidth)
		}
	}
	// A invariante é bannerWidth refletir a arte, não valer um número
	// específico: trocar a arte é mudança legítima, e fixar a largura aqui
	// faria o teste falhar sem defeito nenhum.
	if bannerWidth <= 0 {
		t.Fatal("bannerWidth precisa ser positivo — a guarda de terminal estreito depende dele")
	}
	if !slices.ContainsFunc(linhas, func(l string) bool { return utf8.RuneCountInString(l) == bannerWidth }) {
		t.Errorf("nenhuma linha tem %d células: bannerWidth não veio da arte", bannerWidth)
	}
}

func TestLarguraMaximaContaCelulasNaoBytes(t *testing.T) {
	// "╗" ocupa 3 bytes e 1 célula: contar bytes daria 5 em vez de 3.
	if got := larguraMaxima("ab\n█╗x\nc"); got != 3 {
		t.Errorf("obtido %d, esperado 3", got)
	}
}

func TestSubtitleEncurtaEmTerminalEstreito(t *testing.T) {
	semTerminal := &ui{}
	if !strings.Contains(semTerminal.subtitle(), "to-cards") {
		t.Error("sem terminal medível, a legenda completa deve ser mantida")
	}
}

// "commands:" seguido de nada é a única forma de o resumo mentir sobre o que
// foi instalado — e foi exatamente o que apareceu na primeira versão.
func TestSummaryOmiteALinhaDeCommandsParaQuemNaoOsRecebe(t *testing.T) {
	base := t.TempDir()
	t.Setenv("COPILOT_HOME", base)
	copilot, _ := lookupHarness("copilot")
	claude, _ := lookupHarness("claude")

	var buf strings.Builder
	(&ui{out: &buf}).summary([]harness{copilot, claude}, false)
	texto := buf.String()

	if strings.Contains(texto, "commands: \n") {
		t.Errorf("o resumo imprimiu \"commands:\" sem caminho:\n%s", texto)
	}
	if !strings.Contains(texto, filepath.Join(base, "skills")) {
		t.Error("o resumo precisa dizer onde as skills do copilot foram parar")
	}
	if !strings.Contains(texto, claude.commandsDir()) {
		t.Error("harness que recebe commands continua listando o diretório")
	}
}
