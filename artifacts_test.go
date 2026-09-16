package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// prepara um harness com alguns artefatos já instalados no disco.
func comArtefatos(t *testing.T, envVar string, skills, commands []string) harness {
	t.Helper()
	base := t.TempDir()
	t.Setenv(envVar, base)
	for _, nome := range skills {
		if err := os.MkdirAll(filepath.Join(base, "skills", nome), dirPerm); err != nil {
			t.Fatal(err)
		}
	}
	if len(commands) > 0 {
		if err := os.MkdirAll(filepath.Join(base, "commands"), dirPerm); err != nil {
			t.Fatal(err)
		}
	}
	for _, nome := range commands {
		if err := os.WriteFile(filepath.Join(base, "commands", nome+".md"), []byte("x"), filePerm); err != nil {
			t.Fatal(err)
		}
	}
	for _, h := range harnesses {
		if h.envVar == envVar {
			return h
		}
	}
	t.Fatalf("harness com envVar %q não existe", envVar)
	return harness{}
}

func TestHasSkillEHasCommandOlhamOCaminhoDaInstalacao(t *testing.T) {
	h := comArtefatos(t, "CLAUDE_DIR", []string{"planning"}, []string{"implement"})

	if !h.hasSkill("planning") {
		t.Error("skill presente no disco deveria ser detectada")
	}
	if h.hasSkill("tdd") {
		t.Error("skill ausente não deveria ser detectada")
	}
	if !h.hasCommand("implement") {
		t.Error("command presente deveria ser detectado")
	}
	if h.hasCommand("planning") {
		t.Error("command ausente não deveria ser detectado; skill de mesmo nome não conta")
	}
}

func TestArtifactLabelListaOsHarnessesOndeJaExiste(t *testing.T) {
	claude := comArtefatos(t, "CLAUDE_DIR", []string{"planning", "tdd"}, nil)
	omp := comArtefatos(t, "OMP_AGENTS_DIR", []string{"planning"}, nil)
	opencode := comArtefatos(t, "OPENCODE_DIR", nil, nil)

	alvos := []harness{opencode, claude, omp}
	semCor := &ui{}

	casos := []struct {
		nome     string
		esperado string
	}{
		{"planning", "planning - [claude, omp]"},
		{"tdd", "tdd - [claude]"},
		{"to-spec", "to-spec"},
	}
	for _, caso := range casos {
		if got := semCor.artifactLabel(caso.nome, alvos, harness.hasSkill); got != caso.esperado {
			t.Errorf("obtido %q, esperado %q", got, caso.esperado)
		}
	}
}

func TestArtifactLabelIgnoraHarnessNaoSelecionado(t *testing.T) {
	claude := comArtefatos(t, "CLAUDE_DIR", []string{"planning"}, nil)
	opencode := comArtefatos(t, "OPENCODE_DIR", nil, nil)

	// claude tem a skill, mas não está entre os alvos: a marca some.
	_ = claude
	if got := (&ui{}).artifactLabel("planning", []harness{opencode}, harness.hasSkill); got != "planning" {
		t.Errorf("harness fora da seleção não deve aparecer na marca; obtido %q", got)
	}
}

func TestArtifactLabelColoreApenasAMarcacao(t *testing.T) {
	claude := comArtefatos(t, "CLAUDE_DIR", []string{"planning"}, nil)

	got := (&ui{c: colored}).artifactLabel("planning", []harness{claude}, harness.hasSkill)
	if !strings.HasPrefix(got, "planning") {
		t.Errorf("o nome do artefato não pode vir colorido: %q", got)
	}
	if !strings.Contains(got, " - "+colored.yellow+"[claude]"+colored.reset) {
		t.Errorf("o separador fica fora da cor, e só a marcação em amarelo: %q", got)
	}
}

func TestFiltrarNaOrdemUsaAOrdemDoManifestoNaoADaEscolha(t *testing.T) {
	manifesto := []string{"planning", "to-spec", "to-cards", "implement", "tdd"}

	casos := []struct {
		nome       string
		escolhidos []string
		esperado   string
	}{
		{"escolha embaralhada", []string{"tdd", "planning", "to-cards"}, "planning,to-cards,tdd"},
		{"escolha completa", manifesto, "planning,to-spec,to-cards,implement,tdd"},
		{"escolha vazia", nil, ""},
		{"nome desconhecido é ignorado", []string{"planning", "inexistente"}, "planning"},
	}
	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			got := strings.Join(filtrarNaOrdem(manifesto, caso.escolhidos), ",")
			if got != caso.esperado {
				t.Errorf("obtido %q, esperado %q", got, caso.esperado)
			}
		})
	}
}

func TestSelectArtifactsDevolveOManifestoInteiroSemTerminal(t *testing.T) {
	completo := manifest{Skills: []string{"planning", "tdd"}, Commands: []string{"implement"}}

	got, err := (&ui{}).selectArtifacts(completo, harnesses)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strings.Join(got.Skills, ",") != "planning,tdd" || strings.Join(got.Commands, ",") != "implement" {
		t.Errorf("sem terminal, ausência de escolha significa tudo; obtido %+v", got)
	}
}

func TestSecaoVaziaNaoCriaDiretorio(t *testing.T) {
	base := t.TempDir()
	t.Setenv("CLAUDE_DIR", base)
	h, _ := lookupHarness("claude")

	// só commands selecionados: skills/ não deve nascer vazio.
	in := &installer{
		content: content,
		man:     manifest{Commands: []string{"implement"}},
		ui:      uiSilencioso(),
		force:   true,
	}
	if err := in.run([]harness{h}); err != nil {
		t.Fatalf("instalação falhou: %v", err)
	}
	if _, err := os.Stat(filepath.Join(base, "skills")); !os.IsNotExist(err) {
		t.Error("nenhuma skill selecionada não deveria criar o diretório skills/")
	}
	if _, err := os.Stat(filepath.Join(base, "commands", "implement.md")); err != nil {
		t.Errorf("o command selecionado deveria ter sido instalado: %v", err)
	}
}
