package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func uiSilencioso() *ui { return &ui{out: io.Discard} }

// manifestoReal carrega o manifest.toml embutido. Todo teste que percorre o
// conteúdo passa por aqui, então um manifesto quebrado derruba a suíte
// inteira com a mensagem certa, em vez de um nil silencioso.
func manifestoReal(t *testing.T) manifest {
	t.Helper()
	m, err := loadManifest(content)
	if err != nil {
		t.Fatalf("manifesto embutido não carregou: %v", err)
	}
	return m
}

func TestAssembleMontaFrontmatterMaisCorpo(t *testing.T) {
	src := fstest.MapFS{
		"commands/exemplo/claude.yml": {Data: []byte("description: \"teste\"\n")},
		"commands/exemplo/body.md":    {Data: []byte("# Corpo\n")},
	}
	obtido, err := assemble(src, "commands/exemplo", "claude")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	esperado := "---\ndescription: \"teste\"\n---\n\n# Corpo\n"
	if string(obtido) != esperado {
		t.Errorf("obtido:\n%q\nesperado:\n%q", obtido, esperado)
	}
}

func TestAssembleFechaOFrontmatterQuandoOYmlNaoTerminaEmNovaLinha(t *testing.T) {
	src := fstest.MapFS{
		"commands/exemplo/omp.yml": {Data: []byte("description: \"sem newline\"")},
		"commands/exemplo/body.md": {Data: []byte("corpo")},
	}
	obtido, err := assemble(src, "commands/exemplo", "omp")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !strings.Contains(string(obtido), "\"sem newline\"\n---\n") {
		t.Errorf("delimitador colou na última chave:\n%q", obtido)
	}
}

func TestAssembleFalhaQuandoOHarnessNaoTemFrontmatter(t *testing.T) {
	src := fstest.MapFS{
		"commands/exemplo/claude.yml": {Data: []byte("description: \"x\"\n")},
		"commands/exemplo/body.md":    {Data: []byte("corpo")},
	}
	if _, err := assemble(src, "commands/exemplo", "opencode"); err == nil {
		t.Error("esperava erro para harness sem .yml correspondente")
	}
}

func TestCheckOverwriteRespeitaCadaResposta(t *testing.T) {
	casos := []struct {
		nome     string
		resposta overwriteDecision
		force    bool
		escreve  bool
		ficaTudo bool
	}{
		{"pular mantém o arquivo", overwriteSkip, false, false, false},
		{"substituir escreve uma vez", overwriteReplace, false, true, false},
		{"todos escreve e fixa a decisão", overwriteAll, false, true, true},
		{"force não pergunta", overwriteSkip, true, true, false},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			dst := filepath.Join(t.TempDir(), "existente.md")
			if err := os.WriteFile(dst, []byte("antigo"), filePerm); err != nil {
				t.Fatal(err)
			}
			in := &installer{
				ui:    uiSilencioso(),
				force: caso.force,
				prompt: func(string, string) (overwriteDecision, error) {
					return caso.resposta, nil
				},
			}
			escreve, err := in.checkOverwrite(dst, "exemplo")
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if escreve != caso.escreve {
				t.Errorf("escreve: obtido %v, esperado %v", escreve, caso.escreve)
			}
			if in.overwriteAll != caso.ficaTudo {
				t.Errorf("overwriteAll: obtido %v, esperado %v", in.overwriteAll, caso.ficaTudo)
			}
		})
	}
}

func TestCheckOverwriteNaoPerguntaQuandoODestinoNaoExiste(t *testing.T) {
	in := &installer{
		ui: uiSilencioso(),
		prompt: func(string, string) (overwriteDecision, error) {
			t.Fatal("não deveria perguntar sobre arquivo inexistente")
			return overwriteSkip, nil
		},
	}
	escreve, err := in.checkOverwrite(filepath.Join(t.TempDir(), "novo.md"), "novo")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !escreve {
		t.Error("destino inexistente deveria liberar a escrita")
	}
}

func TestCheckOverwritePulaConflitoSemTerminal(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "existente.md")
	if err := os.WriteFile(dst, []byte("antigo"), filePerm); err != nil {
		t.Fatal(err)
	}
	in := &installer{ui: uiSilencioso(), prompt: nil}
	escreve, err := in.checkOverwrite(dst, "exemplo")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if escreve {
		t.Error("sem terminal, o conflito deve ser pulado em vez de sobrescrito")
	}
}

func TestInstalacaoCompletaGravaSkillsComReferencesECommandsMontados(t *testing.T) {
	base := t.TempDir()
	t.Setenv("CLAUDE_DIR", base)
	h, _ := lookupHarness("claude")

	man := manifestoReal(t)
	in := &installer{content: content, man: man, ui: uiSilencioso(), force: true}
	if err := in.run([]harness{h}); err != nil {
		t.Fatalf("instalação falhou: %v", err)
	}

	for _, name := range man.Skills {
		if _, err := os.Stat(filepath.Join(base, "skills", name, "SKILL.md")); err != nil {
			t.Errorf("SKILL.md ausente para %s: %v", name, err)
		}
	}
	// planning tem references/: a cópia precisa trazer a subárvore inteira.
	if _, err := os.Stat(filepath.Join(base, "skills", "planning", "references", "plan-format.md")); err != nil {
		t.Errorf("references/ não foi copiado: %v", err)
	}

	for _, name := range man.Commands {
		data, err := os.ReadFile(filepath.Join(base, "commands", name+".md"))
		if err != nil {
			t.Fatalf("command %s ausente: %v", name, err)
		}
		if !strings.HasPrefix(string(data), "---\n") || strings.Count(string(data), "\n---\n") == 0 {
			t.Errorf("command %s não tem frontmatter delimitado:\n%.80q", name, data)
		}
	}
}

func TestInstalacaoPreservaArquivoExistenteQuandoARespostaEhPular(t *testing.T) {
	base := t.TempDir()
	t.Setenv("OMP_AGENTS_DIR", base)
	h, _ := lookupHarness("omp")

	intocado := filepath.Join(base, "commands", "implement.md")
	if err := os.MkdirAll(filepath.Dir(intocado), dirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(intocado, []byte("conteúdo do usuário"), filePerm); err != nil {
		t.Fatal(err)
	}

	in := &installer{
		content: content,
		man:     manifestoReal(t),
		ui:      uiSilencioso(),
		prompt: func(_, dst string) (overwriteDecision, error) {
			if dst == intocado {
				return overwriteSkip, nil
			}
			return overwriteReplace, nil
		},
	}
	if err := in.run([]harness{h}); err != nil {
		t.Fatalf("instalação falhou: %v", err)
	}

	data, err := os.ReadFile(intocado)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "conteúdo do usuário" {
		t.Errorf("arquivo do usuário foi sobrescrito apesar do 'pular': %q", data)
	}
}

// O manifesto existe só para fixar a ordem do fluxo. Este teste é o que
// impede que ele vire uma segunda fonte de verdade: uma skill nova no
// repositório sem entrada aqui — ou uma entrada sem skill — quebra a build.
func TestManifestoCobreExatamenteOQueFoiEmbutido(t *testing.T) {
	man := manifestoReal(t)
	casos := []struct {
		raiz      string
		manifesto []string
	}{
		{"skills", man.Skills},
		{"commands", man.Commands},
	}

	for _, caso := range casos {
		t.Run(caso.raiz, func(t *testing.T) {
			entradas, err := fs.ReadDir(content, caso.raiz)
			if err != nil {
				t.Fatalf("leitura de %s: %v", caso.raiz, err)
			}
			var noDisco []string
			for _, entrada := range entradas {
				if entrada.IsDir() {
					noDisco = append(noDisco, entrada.Name())
				}
			}
			for _, nome := range noDisco {
				if !slices.Contains(caso.manifesto, nome) {
					t.Errorf("%s/%s está embutido mas não foi registrado no manifesto", caso.raiz, nome)
				}
			}
			for _, nome := range caso.manifesto {
				if !slices.Contains(noDisco, nome) {
					t.Errorf("%s/%s está no manifesto mas não foi embutido", caso.raiz, nome)
				}
			}
		})
	}
}

// Cada command precisa de um .yml por harness: a falta de um só aparece no
// dia em que alguém instala naquele harness.
func TestTodoCommandTemFrontmatterParaTodosOsHarnesses(t *testing.T) {
	for _, name := range manifestoReal(t).Commands {
		for _, h := range harnesses {
			if !h.commands {
				continue
			}
			if _, err := assemble(content, "commands/"+name, h.name); err != nil {
				t.Errorf("command %s não monta para %s: %v", name, h.name, err)
			}
		}
	}
}

func TestDryRunNaoCriaNadaNoDisco(t *testing.T) {
	base := filepath.Join(t.TempDir(), "destino-que-nao-existe")
	t.Setenv("CLAUDE_DIR", base)
	h, _ := lookupHarness("claude")

	in := &installer{content: content, man: manifestoReal(t), ui: uiSilencioso(), dryRun: true, force: true}
	if err := in.run([]harness{h}); err != nil {
		t.Fatalf("simulação falhou: %v", err)
	}

	if _, err := os.Stat(base); !os.IsNotExist(err) {
		t.Errorf("--dry-run criou %s; esperava que nem o diretório base existisse", base)
	}
}

func TestDryRunNaoAlteraArquivoExistenteNemComForce(t *testing.T) {
	base := t.TempDir()
	t.Setenv("OMP_AGENTS_DIR", base)
	h, _ := lookupHarness("omp")

	comando := filepath.Join(base, "commands", "implement.md")
	skill := filepath.Join(base, "skills", "tdd", "SKILL.md")
	for _, caminho := range []string{comando, skill} {
		if err := os.MkdirAll(filepath.Dir(caminho), dirPerm); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(caminho, []byte("original"), filePerm); err != nil {
			t.Fatal(err)
		}
	}

	in := &installer{content: content, man: manifestoReal(t), ui: uiSilencioso(), dryRun: true, force: true}
	if err := in.run([]harness{h}); err != nil {
		t.Fatalf("simulação falhou: %v", err)
	}

	for _, caminho := range []string{comando, skill} {
		data, err := os.ReadFile(caminho)
		if err != nil {
			t.Fatalf("%s sumiu durante a simulação: %v", caminho, err)
		}
		if string(data) != "original" {
			t.Errorf("%s foi alterado em --dry-run: %q", caminho, data)
		}
	}
}

func TestDryRunAindaPerguntaSobreConflito(t *testing.T) {
	base := t.TempDir()
	t.Setenv("CLAUDE_DIR", base)
	h, _ := lookupHarness("claude")

	existente := filepath.Join(base, "skills", "planning")
	if err := os.MkdirAll(existente, dirPerm); err != nil {
		t.Fatal(err)
	}

	var perguntou int
	in := &installer{
		content: content,
		man:     manifestoReal(t),
		ui:      uiSilencioso(),
		dryRun:  true,
		prompt: func(string, string) (overwriteDecision, error) {
			perguntou++
			return overwriteReplace, nil
		},
	}
	if err := in.run([]harness{h}); err != nil {
		t.Fatalf("simulação falhou: %v", err)
	}
	if perguntou != 1 {
		t.Errorf("conflito perguntado %d vezes, esperava 1 — a simulação precisa exercitar o fluxo inteiro", perguntou)
	}
}

func TestCountTreeContaOsArquivosQueSeriamGravados(t *testing.T) {
	in := &installer{content: content}
	total, err := in.countTree("skills/planning")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	// SKILL.md + references/decision-tree.md + references/plan-format.md
	if total != 3 {
		t.Errorf("obtido %d arquivos para skills/planning, esperado 3", total)
	}
}

// O Copilot recebe as skills e mais nada. Instalar os commands ali deixaria
// cinco arquivos que nenhuma ferramenta abre — e que ninguém iria remover.
func TestCopilotRecebeSkillsENenhumCommand(t *testing.T) {
	base := t.TempDir()
	t.Setenv("COPILOT_HOME", base)
	h, _ := lookupHarness("copilot")

	man := manifestoReal(t)
	if len(man.Commands) == 0 {
		t.Fatal("o manifesto precisa ter commands para este teste significar algo")
	}
	in := &installer{content: content, man: man, ui: uiSilencioso(), force: true}
	if err := in.run([]harness{h}); err != nil {
		t.Fatalf("instalação falhou: %v", err)
	}

	for _, name := range man.Skills {
		if _, err := os.Stat(filepath.Join(base, "skills", name, "SKILL.md")); err != nil {
			t.Errorf("SKILL.md ausente para %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(base, "skills", "planning", "references", "plan-format.md")); err != nil {
		t.Errorf("references/ não foi copiado: %v", err)
	}
	if _, err := os.Stat(filepath.Join(base, "commands")); !os.IsNotExist(err) {
		t.Errorf("o diretório de commands não deveria ter sido criado (err=%v)", err)
	}
}
