package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tempRaiz devolve um diretório temporário com os symlinks já resolvidos. No
// macOS o t.TempDir() vive sob /var, que é link para /private/var: sem isso a
// comparação de caminhos falharia por um prefixo que ninguém escreveu.
func tempRaiz(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestRaizDoProjetoAchaODiretorioQueContemGit(t *testing.T) {
	raiz := tempRaiz(t)
	if err := os.Mkdir(filepath.Join(raiz, ".git"), dirPerm); err != nil {
		t.Fatal(err)
	}

	got, ok := raizDoProjeto(raiz)
	if !ok {
		t.Fatal("diretório com .git deveria ser reconhecido como raiz")
	}
	if got != raiz {
		t.Errorf("obtido %q, esperado %q", got, raiz)
	}
}

// Worktree e submódulo gravam .git como arquivo apontando para o diretório
// real. Uma checagem que exigisse diretório não enxergaria nenhum dos dois —
// e este repositório trabalha com worktrees em .wt/.
func TestRaizDoProjetoAceitaGitComoArquivo(t *testing.T) {
	raiz := tempRaiz(t)
	conteudo := []byte("gitdir: /outro/lugar/.git/worktrees/feat\n")
	if err := os.WriteFile(filepath.Join(raiz, ".git"), conteudo, filePerm); err != nil {
		t.Fatal(err)
	}

	got, ok := raizDoProjeto(raiz)
	if !ok {
		t.Fatal(".git como arquivo deveria contar como repositório")
	}
	if got != raiz {
		t.Errorf("obtido %q, esperado %q", got, raiz)
	}
}

// A instalação vai para a raiz, não para o cwd: é lá que os quatro harnesses
// procuram no walk-up deles. Rodar o instalador de dentro de cmd/ não pode
// criar um .claude/ perdido no meio da árvore.
func TestRaizDoProjetoSobeAPartirDeUmaSubpasta(t *testing.T) {
	raiz := tempRaiz(t)
	if err := os.Mkdir(filepath.Join(raiz, ".git"), dirPerm); err != nil {
		t.Fatal(err)
	}
	fundo := filepath.Join(raiz, "cmd", "coder", "interno")
	if err := os.MkdirAll(fundo, dirPerm); err != nil {
		t.Fatal(err)
	}

	got, ok := raizDoProjeto(fundo)
	if !ok {
		t.Fatal("a busca deveria subir até achar o .git")
	}
	if got != raiz {
		t.Errorf("obtido %q, esperado a raiz %q", got, raiz)
	}
}

// Fora de repositório a resposta é "não", e a busca para na raiz do
// filesystem em vez de girar para sempre.
func TestRaizDoProjetoDevolveFalsoForaDeRepositorio(t *testing.T) {
	if raiz, ok := raizDoProjeto(tempRaiz(t)); ok {
		t.Errorf("diretório sem .git não é repositório; devolveu %q", raiz)
	}
}

func TestParseArgsLeOEscopo(t *testing.T) {
	opts, err := parseArgs([]string{"--scope", "project"})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if opts.scope != "project" {
		t.Errorf("obtido %q, esperado \"project\"", opts.scope)
	}

	opts, err = parseArgs([]string{"--scope=global", "-n"})
	if err != nil {
		t.Fatalf("erro inesperado na forma --scope=: %v", err)
	}
	if opts.scope != "global" || !opts.dryRun {
		t.Errorf("opções obtidas: %+v", opts)
	}

	if _, err := parseArgs([]string{"--scope"}); err == nil {
		t.Error("--scope sem valor deveria falhar")
	}
}

func TestResolveEscopoFlagTraduzOsDoisValores(t *testing.T) {
	raiz := filepath.Join("/tmp", "projeto-de-teste")

	e, err := resolveEscopoFlag("project", raiz)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if e.raiz != raiz {
		t.Errorf("escopo de projeto deveria apontar a raiz; obtido %q", e.raiz)
	}

	// A flag vence a detecção: estar dentro de um repositório não pode
	// arrastar para o projeto quem pediu global explicitamente.
	e, err = resolveEscopoFlag("global", raiz)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if e.projeto() {
		t.Errorf("--scope global deveria ignorar a raiz detectada; obtido %q", e.raiz)
	}
}

// Instalar "no projeto" sem projeto nenhum espalharia .claude/ por qualquer
// diretório de onde o binário fosse chamado. Melhor parar e dizer.
func TestResolveEscopoFlagRecusaProjetoForaDeRepositorio(t *testing.T) {
	_, err := resolveEscopoFlag("project", "")
	if err == nil {
		t.Fatal("esperava erro ao pedir escopo de projeto fora de repositório")
	}
	if !strings.Contains(err.Error(), "git") {
		t.Errorf("o erro precisa dizer que não há repositório: %v", err)
	}
}

func TestResolveEscopoFlagRejeitaValorInvalido(t *testing.T) {
	for _, valor := range []string{"projeto", "local", "usuario", ""} {
		if _, err := resolveEscopoFlag(valor, "/tmp/x"); err == nil {
			t.Errorf("esperava erro para %q", valor)
		}
	}
}

// A tabela inteira de quando o menu de escopo aparece. Separada do
// formulário de propósito: é a decisão que precisa de teste, não o desenho.
func TestDevePerguntarEscopoCobreATabelaInteira(t *testing.T) {
	casos := []struct {
		nome       string
		opts       options
		emRepo     bool
		interativo bool
		esperado   bool
	}{
		{"repo, terminal, sem flag", options{}, true, true, true},
		{"a flag já decidiu", options{scope: escopoGlobal}, true, true, false},
		{"fora de repositório não há o que perguntar", options{}, false, true, false},
		{"sem terminal não há onde perguntar", options{}, true, false, false},
		// --harness é o caminho de script, e o --help promete "sem menu
		// interativo": ele pula os três menus, não só os dois antigos.
		{"--harness pula o menu", options{harnessList: "claude"}, true, true, false},
	}

	for _, caso := range casos {
		t.Run(caso.nome, func(t *testing.T) {
			got := devePerguntarEscopo(caso.opts, caso.emRepo, caso.interativo)
			if got != caso.esperado {
				t.Errorf("obtido %v, esperado %v", got, caso.esperado)
			}
		})
	}
}

func TestSelectScopePrefereAFlagAoMenu(t *testing.T) {
	raiz := tempRaiz(t)
	if err := os.Mkdir(filepath.Join(raiz, ".git"), dirPerm); err != nil {
		t.Fatal(err)
	}
	t.Chdir(raiz)

	e, err := selectScope(&ui{out: io.Discard}, options{scope: escopoProjeto})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if e.raiz != raiz {
		t.Errorf("obtido %q, esperada a raiz detectada %q", e.raiz, raiz)
	}
}

// Fora de repositório o instalador segue como sempre foi, sem perguntar e
// sem erro: o escopo global é o comportamento de origem.
func TestSelectScopeCaiNoGlobalForaDeRepositorio(t *testing.T) {
	t.Chdir(tempRaiz(t))

	e, err := selectScope(&ui{out: io.Discard}, options{})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if e.projeto() {
		t.Errorf("esperava escopo global; obtido %q", e.raiz)
	}
}

// O rótulo do menu mostra onde a instalação vai cair — a mesma regra da
// marca [instalado]: a tela não pode descrever um lugar e escrever em outro.
func TestRotulosDeEscopoMostramAOndeVaiCair(t *testing.T) {
	raiz := filepath.Join("/tmp", "projeto-de-teste")
	projeto, global := rotulosDeEscopo(raiz)
	if !strings.Contains(projeto, raiz) {
		t.Errorf("o rótulo de projeto precisa nomear a raiz: %q", projeto)
	}
	if strings.Contains(global, raiz) {
		t.Errorf("o rótulo global não pode citar a raiz do projeto: %q", global)
	}
}

// O resumo do escopo nomeia a raiz só quando ela existe: "global (…)" não
// tem um caminho único para citar, porque cada harness tem o seu.
func TestResumoDoEscopoNomeiaARaizSoNoProjeto(t *testing.T) {
	raiz := filepath.Join("/tmp", "projeto-de-teste")
	if got, esperado := (escopo{raiz: raiz}).resumo(), "projeto ("+raiz+")"; got != esperado {
		t.Errorf("obtido %q, esperado %q", got, esperado)
	}
	if got := (escopo{}).resumo(); got != "global" {
		t.Errorf("obtido %q, esperado \"global\"", got)
	}
}

// A prova de ponta a ponta: com o escopo de projeto, os arquivos caem na
// raiz do repositório e a base de usuário não é tocada — nem a apontada por
// override de ambiente, que aqui não tem vez.
func TestInstalacaoNoEscopoDeProjetoGravaNaRaizENaoNoHome(t *testing.T) {
	raiz := tempRaiz(t)
	if err := os.Mkdir(filepath.Join(raiz, ".git"), dirPerm); err != nil {
		t.Fatal(err)
	}
	home := tempRaiz(t)
	t.Setenv("CLAUDE_DIR", home)

	claude, _ := lookupHarness("claude")
	targets := comEscopo([]harness{claude}, escopo{raiz: raiz})

	man := manifestoReal(t)
	in := &installer{content: content, man: man, ui: uiSilencioso(), force: true}
	if err := in.run(targets); err != nil {
		t.Fatalf("instalação falhou: %v", err)
	}

	for _, name := range man.Skills {
		if _, err := os.Stat(filepath.Join(raiz, ".claude", "skills", name, "SKILL.md")); err != nil {
			t.Errorf("SKILL.md ausente para %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(raiz, ".claude", "skills", "planning", "references", "plan-format.md")); err != nil {
		t.Errorf("references/ não foi copiado: %v", err)
	}
	for _, name := range man.Commands {
		if _, err := os.Stat(filepath.Join(raiz, ".claude", "commands", name+".md")); err != nil {
			t.Errorf("command %s ausente: %v", name, err)
		}
	}

	entradas, err := os.ReadDir(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(entradas) != 0 {
		t.Errorf("o escopo de projeto não pode escrever na base de usuário; encontrei %d entradas", len(entradas))
	}
}

// O Copilot é o caso em que o diretório de projeto tem outro nome: as skills
// vão para .github/skills, e nenhum command aparece.
func TestInstalacaoNoProjetoLevaOCopilotParaGithub(t *testing.T) {
	raiz := tempRaiz(t)
	copilot, _ := lookupHarness("copilot")
	targets := comEscopo([]harness{copilot}, escopo{raiz: raiz})

	in := &installer{content: content, man: manifestoReal(t), ui: uiSilencioso(), force: true}
	if err := in.run(targets); err != nil {
		t.Fatalf("instalação falhou: %v", err)
	}

	if _, err := os.Stat(filepath.Join(raiz, ".github", "skills", "planning", "SKILL.md")); err != nil {
		t.Errorf("skill ausente em .github/skills: %v", err)
	}
	if caminhoExiste(filepath.Join(raiz, ".copilot")) {
		t.Error("o escopo de projeto do Copilot não é .copilot")
	}
	if caminhoExiste(filepath.Join(raiz, ".github", "commands")) {
		t.Error("o Copilot não recebe commands em nenhum escopo")
	}
}
