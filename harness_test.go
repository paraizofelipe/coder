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

// O escopo de projeto sai da raiz do repositório e não consulta o override de
// ambiente: OPENCODE_DIR e companhia nomeiam a base de usuário, e deixá-los
// vencer faria o menu dizer "projeto" enquanto a escrita cai em outro lugar.
func TestBaseNoEscopoDeProjetoSaiDaRaizEIgnoraOOverride(t *testing.T) {
	t.Setenv("CLAUDE_DIR", filepath.Join("/tmp", "base-global-de-teste"))
	raiz := filepath.Join("/tmp", "projeto-de-teste")
	claude, _ := lookupHarness("claude")

	h := comEscopo([]harness{claude}, escopo{raiz: raiz})[0]
	if got, esperado := h.skillsDir(), filepath.Join(raiz, ".claude", "skills"); got != esperado {
		t.Errorf("skillsDir: obtido %q, esperado %q", got, esperado)
	}
	if got, esperado := h.commandsDir(), filepath.Join(raiz, ".claude", "commands"); got != esperado {
		t.Errorf("commandsDir: obtido %q, esperado %q", got, esperado)
	}
}

// Raiz vazia é o escopo global, e é o zero value: harness que ninguém marcou
// continua apontando para o home, como antes de existir escopo.
func TestEscopoVazioMantemOComportamentoGlobal(t *testing.T) {
	t.Setenv("CLAUDE_DIR", "/tmp/base-global-de-teste")
	claude, _ := lookupHarness("claude")

	h := comEscopo([]harness{claude}, escopo{})[0]
	if h.projeto() {
		t.Error("raiz vazia não pode contar como escopo de projeto")
	}
	if got := h.skillsDir(); got != "/tmp/base-global-de-teste/skills" {
		t.Errorf("obtido %q", got)
	}
}

// Invariante, não lista fixa: harness novo sem diretório de projeto cairia na
// raiz do repositório, criando um "<raiz>/skills" que ninguém varre.
func TestTodoHarnessDeclaraDiretorioDeProjeto(t *testing.T) {
	for _, h := range harnesses {
		if h.projectDir == "" {
			t.Errorf("%s não declara projectDir", h.name)
		}
		if strings.HasPrefix(h.projectDir, "/") || strings.Contains(h.projectDir, "..") {
			t.Errorf("%s: projectDir precisa ser relativo à raiz; obtido %q", h.name, h.projectDir)
		}
	}
}

// O Copilot é o caso em que o nome do diretório de projeto não é o do global:
// a CLI lê .github/skills no repositório, e ~/.copilot/skills no usuário.
func TestCopilotNoProjetoVaiParaGithubEContinuaSemCommands(t *testing.T) {
	raiz := filepath.Join("/tmp", "projeto-de-teste")
	copilot, _ := lookupHarness("copilot")

	h := comEscopo([]harness{copilot}, escopo{raiz: raiz})[0]
	if got, esperado := h.skillsDir(), filepath.Join(raiz, ".github", "skills"); got != esperado {
		t.Errorf("obtido %q, esperado %q", got, esperado)
	}
	if got := h.commandsDir(); got != "" {
		t.Errorf("o escopo não muda quem lê command; apontou %q", got)
	}
}

// A marca do menu usa o mesmo caminho da escrita — inclusive no escopo de
// projeto, onde a base de usuário pode contar a história oposta.
func TestMarcaDeInstaladoSegueOEscopoEscolhido(t *testing.T) {
	raiz := t.TempDir()
	if err := os.MkdirAll(filepath.Join(raiz, ".claude", "skills", "planning"), dirPerm); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_DIR", filepath.Join(t.TempDir(), "sem-nada"))
	claude, _ := lookupHarness("claude")

	noProjeto := comEscopo([]harness{claude}, escopo{raiz: raiz})[0]
	if !noProjeto.detected() || !noProjeto.hasSkill("planning") {
		t.Error("o escopo de projeto deveria enxergar o que está em <raiz>/.claude")
	}

	global := comEscopo([]harness{claude}, escopo{})[0]
	if global.detected() || global.hasSkill("planning") {
		t.Error("o escopo global não pode reportar o que só existe no projeto")
	}
}
