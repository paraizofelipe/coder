package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveHarnessFlagDevolveOrdemCanonicaSemRepeticao(t *testing.T) {
	casos := []struct {
		nome     string
		entrada  string
		esperado []string
	}{
		{"valor único", "claude", []string{"claude"}},
		{"separado por vírgula", "opencode,claude", []string{"opencode", "claude"}},
		{"separado por espaço", "omp opencode", []string{"opencode", "omp"}},
		{"fora da ordem canônica", "omp,claude,opencode", []string{"opencode", "claude", "omp"}},
		{"repetido colapsa", "claude,claude,claude", []string{"claude"}},
		// "all" é derivado da tabela: fixar a lista aqui faria todo harness
		// novo quebrar o teste sem nenhum defeito por trás.
		{"all expande para todos", "all", harnessNames()},
		{"maiúsculas são aceitas", "OpenCode,CLAUDE", []string{"opencode", "claude"}},
		{"all vence o resto da lista", "claude,all", harnessNames()},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			obtido, err := resolveHarnessFlag(caso.entrada)
			if err != nil {
				t.Fatalf("erro inesperado para %q: %v", caso.entrada, err)
			}
			if got := namesOf(obtido); strings.Join(got, ",") != strings.Join(caso.esperado, ",") {
				t.Errorf("para %q: obtido %v, esperado %v", caso.entrada, got, caso.esperado)
			}
		})
	}
}

func TestResolveHarnessFlagRejeitaEntradaInvalida(t *testing.T) {
	casos := []struct {
		nome    string
		entrada string
	}{
		{"harness inexistente", "codex"},
		{"harness removido", "pi"},
		{"string vazia", ""},
		{"só separadores", " , , "},
		{"válido misturado com inválido", "claude,codex"},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			if _, err := resolveHarnessFlag(caso.entrada); err == nil {
				t.Errorf("esperava erro para %q, não veio nenhum", caso.entrada)
			}
		})
	}
}

func TestHarnessUsaOverrideDeAmbienteQuandoPresente(t *testing.T) {
	t.Setenv("CLAUDE_DIR", "/tmp/destino-de-teste")
	h, ok := lookupHarness("claude")
	if !ok {
		t.Fatal("harness claude não encontrado")
	}
	if got := h.skillsDir(); got != "/tmp/destino-de-teste/skills" {
		t.Errorf("skillsDir: obtido %q", got)
	}
	if got := h.commandsDir(); got != "/tmp/destino-de-teste/commands" {
		t.Errorf("commandsDir: obtido %q", got)
	}
}

func TestParseArgsLeAsOpcoes(t *testing.T) {
	opts, err := parseArgs([]string{"--force", "--harness", "claude,omp"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !opts.force || opts.harnessList != "claude,omp" {
		t.Errorf("opções obtidas: %+v", opts)
	}

	opts, err = parseArgs([]string{"--harness=opencode", "-f"})
	if err != nil {
		t.Fatalf("erro inesperado na forma --harness=: %v", err)
	}
	if !opts.force || opts.harnessList != "opencode" {
		t.Errorf("opções obtidas: %+v", opts)
	}

	if _, err := parseArgs([]string{"--harness"}); err == nil {
		t.Error("--harness sem valor deveria falhar")
	}
	if _, err := parseArgs([]string{"--nao-existe"}); err == nil {
		t.Error("opção desconhecida deveria falhar")
	}
}

func TestDetectedSegueODiretorioBaseEOOverride(t *testing.T) {
	existente := t.TempDir()
	h, _ := lookupHarness("claude")

	t.Setenv("CLAUDE_DIR", existente)
	if !h.detected() {
		t.Error("diretório base existente deveria contar como detectado")
	}

	t.Setenv("CLAUDE_DIR", filepath.Join(existente, "nao-existe"))
	if h.detected() {
		t.Error("diretório ausente não deveria contar como detectado")
	}

	arquivo := filepath.Join(existente, "arquivo")
	if err := os.WriteFile(arquivo, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_DIR", arquivo)
	if h.detected() {
		t.Error("arquivo comum não é diretório de harness e não deveria ser detectado")
	}
}

func TestHarnessLabelMarcaSoOQueFoiDetectado(t *testing.T) {
	base := t.TempDir()
	presente := filepath.Join(base, "presente")
	if err := os.Mkdir(presente, 0o755); err != nil {
		t.Fatal(err)
	}

	semCor := &ui{}
	claude, _ := lookupHarness("claude")

	t.Setenv("CLAUDE_DIR", presente)
	if got := semCor.harnessLabel(claude); got != "claude — Claude Code [instalado]" {
		t.Errorf("harness presente: obtido %q", got)
	}

	t.Setenv("CLAUDE_DIR", filepath.Join(base, "ausente"))
	if got := semCor.harnessLabel(claude); got != "claude — Claude Code" {
		t.Errorf("harness ausente não deve receber marcação: obtido %q", got)
	}
}

func TestHarnessLabelColoreApenasAMarcacao(t *testing.T) {
	presente := t.TempDir()
	t.Setenv("OPENCODE_DIR", presente)
	opencode, _ := lookupHarness("opencode")

	got := (&ui{c: colored}).harnessLabel(opencode)
	if !strings.HasPrefix(got, "opencode — OpenCode") {
		t.Errorf("o nome do harness não pode vir colorido: %q", got)
	}
	if !strings.Contains(got, colored.yellow+" [instalado]"+colored.reset) {
		t.Errorf("a marcação deveria estar em amarelo e fechar o escape: %q", got)
	}
}

// A marca do menu e o destino da escrita precisam contar a mesma história:
// o Copilot não recebe command, então também não pode reportar um instalado.
func TestHarnessSemCommandsNaoApontaDiretorioNemReportaInstalado(t *testing.T) {
	base := t.TempDir()
	t.Setenv("COPILOT_HOME", base)
	h, ok := lookupHarness("copilot")
	if !ok {
		t.Fatal("harness copilot não encontrado")
	}
	if h.commands {
		t.Fatal("o Copilot não lê command em markdown")
	}
	if got := h.commandsDir(); got != "" {
		t.Errorf("harness sem commands não deve apontar diretório; apontou %q", got)
	}
	if got, esperado := h.skillsDir(), filepath.Join(base, "skills"); got != esperado {
		t.Errorf("obtido %q, esperado %q", got, esperado)
	}

	// Mesmo com o arquivo no disco, a resposta é não: o instalador nunca
	// escreveu ali, e marcar a linha prometeria algo que ele não fez.
	if err := os.MkdirAll(filepath.Join(base, "commands"), dirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "commands", "implement.md"), []byte("x"), filePerm); err != nil {
		t.Fatal(err)
	}
	if h.hasCommand("implement") {
		t.Error("harness sem commands não pode reportar command instalado")
	}
}

func TestAceitaCommandsBastaUmDestinoQueLeia(t *testing.T) {
	copilot, _ := lookupHarness("copilot")
	claude, _ := lookupHarness("claude")
	casos := []struct {
		nome     string
		alvos    []harness
		esperado bool
	}{
		{"nenhum destino", nil, false},
		{"só o copilot", []harness{copilot}, false},
		{"copilot acompanhado", []harness{copilot, claude}, true},
		{"todos", harnesses, true},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			if got := aceitaCommands(caso.alvos); got != caso.esperado {
				t.Errorf("obtido %v, esperado %v", got, caso.esperado)
			}
		})
	}
}
