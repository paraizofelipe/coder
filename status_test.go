package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// instalar grava o manifesto inteiro num destino, que é o estado de partida
// de todo teste daqui: o status lê o disco, então o disco precisa existir.
func instalar(t *testing.T, targets ...harness) manifest {
	t.Helper()
	man := manifestoReal(t)
	in := &installer{content: content, man: man, ui: uiSilencioso(), force: true}
	if err := in.run(targets); err != nil {
		t.Fatalf("instalação falhou: %v", err)
	}
	return man
}

func TestStatusEncontraOQueEstaNoDisco(t *testing.T) {
	base := t.TempDir()
	t.Setenv("CLAUDE_DIR", base)
	t.Setenv("OPENCODE_DIR", filepath.Join(t.TempDir(), "vazio"))
	t.Setenv("OMP_AGENTS_DIR", filepath.Join(t.TempDir(), "vazio"))
	t.Setenv("COPILOT_HOME", filepath.Join(t.TempDir(), "vazio"))

	claude, _ := lookupHarness("claude")
	man := instalar(t, claude)

	encontrados := varrerDestinos(man, "")
	if len(encontrados) != 1 {
		t.Fatalf("esperava um destino ocupado, obtive %d: %+v", len(encontrados), encontrados)
	}
	achado := encontrados[0]
	if achado.harness != "claude" || achado.escopo.projeto() {
		t.Errorf("destino errado: %+v", achado)
	}
	if len(achado.skills) != len(man.Skills) || len(achado.commands) != len(man.Commands) {
		t.Errorf("contagem errada: %d skills, %d commands", len(achado.skills), len(achado.commands))
	}
}

// Destino vazio não vira linha: o status descreve o que existe, e listar os
// oito caminhos possíveis afogaria a informação.
func TestStatusIgnoraDestinoVazio(t *testing.T) {
	for _, env := range []string{"CLAUDE_DIR", "OPENCODE_DIR", "OMP_AGENTS_DIR", "COPILOT_HOME"} {
		t.Setenv(env, filepath.Join(t.TempDir(), "nao-existe"))
	}
	if encontrados := varrerDestinos(manifestoReal(t), ""); len(encontrados) != 0 {
		t.Errorf("nada instalado deveria dar lista vazia; obtive %+v", encontrados)
	}
}

func TestStatusEnxergaOsDoisEscoposQuandoHaProjeto(t *testing.T) {
	base := t.TempDir()
	t.Setenv("CLAUDE_DIR", base)
	for _, env := range []string{"OPENCODE_DIR", "OMP_AGENTS_DIR", "COPILOT_HOME"} {
		t.Setenv(env, filepath.Join(t.TempDir(), "vazio"))
	}
	raiz := tempRaiz(t)

	claude, _ := lookupHarness("claude")
	man := instalar(t, claude)
	instalar(t, comEscopo([]harness{claude}, escopo{raiz: raiz})...)

	encontrados := varrerDestinos(man, raiz)
	if len(encontrados) != 2 {
		t.Fatalf("esperava os dois escopos, obtive %d", len(encontrados))
	}
	var temGlobal, temProjeto bool
	for _, e := range encontrados {
		if e.escopo.projeto() {
			temProjeto = true
		} else {
			temGlobal = true
		}
	}
	if !temGlobal || !temProjeto {
		t.Errorf("faltou um escopo: global=%v projeto=%v", temGlobal, temProjeto)
	}
}

// O achado que justifica o comando: a mesma skill em dois destinos que um só
// harness lê. Medido no container, qual cópia vence varia entre execuções —
// então duas cópias divergentes são um sorteio, não uma preferência.
func TestStatusApontaCopiaVisivelDuasVezesAoMesmoHarness(t *testing.T) {
	claudeDir, ompDir := t.TempDir(), t.TempDir()
	t.Setenv("CLAUDE_DIR", claudeDir)
	t.Setenv("OMP_AGENTS_DIR", ompDir)
	t.Setenv("OPENCODE_DIR", filepath.Join(t.TempDir(), "vazio"))
	t.Setenv("COPILOT_HOME", filepath.Join(t.TempDir(), "vazio"))

	claude, _ := lookupHarness("claude")
	omp, _ := lookupHarness("omp")
	man := instalar(t, claude, omp)

	avisos := duplicatas(varrerDestinos(man, ""))
	junto := strings.Join(avisos, "\n")
	// O OpenCode lê ~/.claude e ~/.agents; nenhum dos dois é destino dele.
	if !strings.Contains(junto, "opencode") {
		t.Errorf("o opencode lê os dois destinos e deveria ser apontado:\n%s", junto)
	}
	if !strings.Contains(junto, claudeDir) || !strings.Contains(junto, ompDir) {
		t.Errorf("o aviso precisa nomear os dois caminhos:\n%s", junto)
	}
}

// Uma cópia só nunca é ambígua, por mais harnesses que a leiam.
func TestStatusNaoAcusaDuplicataComUmDestinoSo(t *testing.T) {
	t.Setenv("CLAUDE_DIR", t.TempDir())
	for _, env := range []string{"OPENCODE_DIR", "OMP_AGENTS_DIR", "COPILOT_HOME"} {
		t.Setenv(env, filepath.Join(t.TempDir(), "vazio"))
	}
	claude, _ := lookupHarness("claude")
	man := instalar(t, claude)

	if avisos := duplicatas(varrerDestinos(man, "")); len(avisos) != 0 {
		t.Errorf("um destino só não é ambíguo; avisos: %v", avisos)
	}
}

// Skill solta no destino, que o manifesto não conhece, não entra na conta: o
// status responde pelo que este instalador colocou lá.
func TestStatusContaSoOQueOManifestoDeclara(t *testing.T) {
	base := t.TempDir()
	t.Setenv("CLAUDE_DIR", base)
	for _, env := range []string{"OPENCODE_DIR", "OMP_AGENTS_DIR", "COPILOT_HOME"} {
		t.Setenv(env, filepath.Join(t.TempDir(), "vazio"))
	}
	claude, _ := lookupHarness("claude")
	man := instalar(t, claude)

	if err := os.MkdirAll(filepath.Join(base, "skills", "skill-de-terceiro"), dirPerm); err != nil {
		t.Fatal(err)
	}
	encontrados := varrerDestinos(man, "")
	if len(encontrados[0].skills) != len(man.Skills) {
		t.Errorf("skill de fora do manifesto entrou na conta: %v", encontrados[0].skills)
	}
}

// Global e projeto não são ambiguidade: os quatro harnesses definem que o
// projeto tem precedência sobre o diretório de usuário. O indefinido medido
// foi outro — duas cópias no MESMO escopo, sem regra que as ordene.
func TestStatusNaoAcusaDuplicataEntreEscoposDiferentes(t *testing.T) {
	t.Setenv("CLAUDE_DIR", t.TempDir())
	for _, env := range []string{"OPENCODE_DIR", "OMP_AGENTS_DIR", "COPILOT_HOME"} {
		t.Setenv(env, filepath.Join(t.TempDir(), "vazio"))
	}
	raiz := tempRaiz(t)
	claude, _ := lookupHarness("claude")
	man := instalar(t, claude)
	instalar(t, comEscopo([]harness{claude}, escopo{raiz: raiz})...)

	if avisos := duplicatas(varrerDestinos(man, raiz)); len(avisos) != 0 {
		t.Errorf("a mesma skill em escopos diferentes é sobreposição prevista, não sorteio:\n%v", avisos)
	}
}

func TestParseArgsLeOStatus(t *testing.T) {
	for _, arg := range []string{"--status", "-s"} {
		opts, err := parseArgs([]string{arg})
		if err != nil {
			t.Fatalf("%s: erro inesperado: %v", arg, err)
		}
		if !opts.showStatus {
			t.Errorf("%s não ligou showStatus", arg)
		}
	}
	if !strings.Contains(helpText(), "--status") {
		t.Error("a ajuda não menciona --status")
	}
}
