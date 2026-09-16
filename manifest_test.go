package main

import (
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoadManifestLeListasNaOrdemDeclarada(t *testing.T) {
	fsys := fstest.MapFS{manifestPath: {Data: []byte(`
skills = ["planning", "implement", "tdd"]
commands = ["implement"]
`)}}

	m, err := loadManifest(fsys)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if got := strings.Join(m.Skills, ","); got != "planning,implement,tdd" {
		t.Errorf("a ordem declarada precisa ser preservada; obtido %q", got)
	}
	if len(m.Commands) != 1 || m.Commands[0] != "implement" {
		t.Errorf("commands: obtido %v", m.Commands)
	}
}

func TestLoadManifestAceitaListaDeCommandsVazia(t *testing.T) {
	fsys := fstest.MapFS{manifestPath: {Data: []byte(`skills = ["planning"]`)}}
	m, err := loadManifest(fsys)
	if err != nil {
		t.Fatalf("skill sem command é configuração válida: %v", err)
	}
	if len(m.Commands) != 0 {
		t.Errorf("commands deveria vir vazio, veio %v", m.Commands)
	}
}

func TestLoadManifestFalhaComMensagemUtil(t *testing.T) {
	casos := []struct {
		nome    string
		arquivo fstest.MapFS
		trecho  string
	}{
		{
			"manifesto ausente",
			fstest.MapFS{},
			"ausente",
		},
		{
			"TOML malformado",
			fstest.MapFS{manifestPath: {Data: []byte("skills = [\"planning\"")}},
			"inválido",
		},
		{
			"tipo errado no campo",
			fstest.MapFS{manifestPath: {Data: []byte("skills = \"planning\"")}},
			"inválido",
		},
		{
			"sem nenhuma skill",
			fstest.MapFS{manifestPath: {Data: []byte("commands = [\"implement\"]")}},
			"não há o que instalar",
		},
		{
			"lista de skills vazia",
			fstest.MapFS{manifestPath: {Data: []byte("skills = []")}},
			"não há o que instalar",
		},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			_, err := loadManifest(caso.arquivo)
			if err == nil {
				t.Fatal("esperava erro, não veio nenhum")
			}
			if !strings.Contains(err.Error(), caso.trecho) {
				t.Errorf("a mensagem precisa explicar o defeito; obtida: %v", err)
			}
		})
	}
}

// A invariante é "a ordem é a do fluxo", não uma lista específica. Fixar a
// lista inteira aqui a duplicaria, e toda skill nova quebraria o teste sem
// nenhum defeito real — que é exatamente o contrário do que ele deve fazer.
func TestManifestoEmbutidoNaoEstaEmOrdemAlfabetica(t *testing.T) {
	m, err := loadManifest(content)
	if err != nil {
		t.Fatalf("o manifesto que vai no binário não carregou: %v", err)
	}
	if m.Skills[0] != "planning" {
		t.Errorf("o fluxo começa em planning; o manifesto começa em %q", m.Skills[0])
	}
	alfabetica := slices.Clone(m.Skills)
	slices.Sort(alfabetica)
	if slices.Equal(m.Skills, alfabetica) {
		t.Error("o manifesto está em ordem alfabética — a ordem precisa ser a do fluxo")
	}
}
