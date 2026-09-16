package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"charm.land/huh/v2"
)

type options struct {
	force       bool
	dryRun      bool
	harnessList string
	scope       string
	showHelp    bool
	showStatus  bool
	showVersion bool
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		switch {
		case errors.Is(err, huh.ErrUserAborted):
			fmt.Fprintln(os.Stderr, "\nCancelado.")
			os.Exit(130)
		default:
			fmt.Fprintf(os.Stderr, "\033[0;31m\033[1m[erro]\033[0m %v\n", err)
			os.Exit(1)
		}
	}
}

func run(args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}
	if opts.showHelp {
		fmt.Print(helpText())
		return nil
	}
	if opts.showVersion {
		fmt.Println(version())
		return nil
	}

	ui := newUI(os.Stdout)
	defer ui.close()
	ui.banner()
	if opts.dryRun {
		ui.dryRunNotice()
	}

	// O status responde e encerra: é leitura do disco, não instalação, e
	// misturá-lo ao fluxo faria o binário perguntar coisas para nada.
	if opts.showStatus {
		return runStatus(ui)
	}

	esc, err := selectScope(ui, opts)
	if err != nil {
		return err
	}
	targets, err := selectTargets(ui, opts.harnessList, esc)
	if err != nil {
		return err
	}
	ui.ok("Escopo: %s", esc.resumo())
	ui.ok("Harnesses: %s", strings.Join(namesOf(targets), " "))
	ui.blank()
	// Antes de gravar, não depois: o alcance real da seleção é a informação
	// que muda a escolha, e depois da instalação ela vira só constatação.
	ui.sobreposicaoNotice(targets)

	man, err := loadManifest(content)
	if err != nil {
		return err
	}
	// --harness é o caminho de script, e o --help promete "sem menu
	// interativo": ele pula os três menus, não só o de destino.
	if opts.harnessList == "" {
		if man, err = ui.selectArtifacts(man, targets); err != nil {
			return err
		}
		ui.ok("Artefatos: %d skills, %d commands", len(man.Skills), len(man.Commands))
		ui.blank()
	}

	in := &installer{
		content: content,
		man:     man,
		ui:      ui,
		force:   opts.force,
		dryRun:  opts.dryRun,
		prompt:  ui.overwritePrompt,
	}
	if err := in.run(targets); err != nil {
		return err
	}
	ui.summary(targets, opts.dryRun)
	return nil
}

// runStatus varre os destinos e relata o que está lá. Não usa registro de
// instalação de propósito: o disco é a verdade, e um arquivo à parte seria
// uma segunda — livre para discordar depois de qualquer remoção manual.
func runStatus(ui *ui) error {
	man, err := loadManifest(content)
	if err != nil {
		return err
	}
	raiz, _ := raizDoProjetoAtual()
	destinos := varrerDestinos(man, raiz)
	ui.statusReport(destinos, duplicatas(destinos))
	return nil
}

// selectTargets resolve --harness quando ele veio, e só abre o menu quando
// não veio. Manter a flag com precedência é o que permite rodar o instalador
// em script, sem terminal.
//
// O escopo é carimbado aqui, num ponto só, nos dois caminhos: o menu usa-o
// para rotular o [instalado], mas quem decide o destino é sempre esta função.
func selectTargets(ui *ui, list string, e escopo) ([]harness, error) {
	targets, err := escolherTargets(ui, list, e)
	if err != nil {
		return nil, err
	}
	return comEscopo(targets, e), nil
}

func escolherTargets(ui *ui, list string, e escopo) ([]harness, error) {
	if list != "" {
		return resolveHarnessFlag(list)
	}
	return ui.selectHarnesses(e)
}

// selectScope decide o escopo antes de qualquer outro menu — e antes de
// propósito: é o escopo que define a base de cada harness, e é da base que
// sai a marca [instalado] da tela seguinte.
func selectScope(u *ui, opts options) (escopo, error) {
	raiz, _ := raizDoProjetoAtual()
	if opts.scope != "" {
		return resolveEscopoFlag(opts.scope, raiz)
	}
	if !devePerguntarEscopo(opts, raiz != "", u.interactive) {
		return escopo{}, nil
	}
	return u.selectEscopo(raiz)
}

// devePerguntarEscopo é a tabela inteira de quando o menu aparece. Fica
// separada do formulário para que a decisão tenha teste sem terminal.
//
// Fora de repositório não há escolha a cobrar, e o silêncio é a resposta
// certa: o escopo global é o comportamento que o instalador sempre teve.
func devePerguntarEscopo(opts options, emRepo, interativo bool) bool {
	return opts.scope == "" && opts.harnessList == "" && emRepo && interativo
}

func namesOf(targets []harness) []string {
	names := make([]string, 0, len(targets))
	for _, h := range targets {
		names = append(names, h.name)
	}
	return names
}

func parseArgs(args []string) (options, error) {
	var opts options
	for i := 0; i < len(args); i++ {
		switch arg := args[i]; arg {
		case "--force", "-f":
			opts.force = true
		case "--dry-run", "-n":
			opts.dryRun = true
		case "--help", "-h":
			opts.showHelp = true
		case "--version", "-v":
			opts.showVersion = true
		case "--status", "-s":
			opts.showStatus = true
		case "--harness":
			if i+1 >= len(args) {
				return opts, errors.New("--harness exige um valor. Ex.: --harness opencode,claude")
			}
			i++
			opts.harnessList = args[i]
		case "--scope":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--scope exige um valor: %s ou %s",
					escopoProjeto, escopoGlobal)
			}
			i++
			opts.scope = args[i]
		case "--local", "-l":
			// Existia para instalar a partir do repositório em vez da rede.
			// O conteúdo agora vem embutido no binário, então não há o que
			// escolher — avisar é mais honesto do que aceitar em silêncio.
			fmt.Fprintln(os.Stderr,
				"[aviso] --local não é mais necessário: o conteúdo vem embutido no binário.")
		default:
			if value, ok := strings.CutPrefix(arg, "--harness="); ok {
				opts.harnessList = value
				continue
			}
			if value, ok := strings.CutPrefix(arg, "--scope="); ok {
				opts.scope = value
				continue
			}
			return opts, fmt.Errorf("opção desconhecida: %q (use --help)", arg)
		}
	}
	return opts, nil
}

// version lê a versão do módulo gravada pelo `go install`. Em binário
// construído localmente não há tag, e o retorno é "dev".
func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "dev"
	}
	return info.Main.Version
}

func helpText() string {
	return `
  coder — instalador das skills e commands do fluxo

  Uso: coder [opções]

  Instala skills e commands nos diretórios nativos do OpenCode, do Claude
  Code, do Oh My Pi e do GitHub Copilot. O conteúdo vem embutido no binário:
  a instalação não acessa a rede.

  O Copilot recebe só as skills: na CLI dele a skill já responde a /<nome>,
  e prompt file segue sendo coisa de IDE, não da linha de comando.

  Opções:
    --dry-run, -n        Percorrer o fluxo sem gravar nada em disco
    --force, -f          Substituir todos os arquivos sem perguntar
    --harness <lista>    Destinos, sem nenhum menu interativo: instala o
                         manifesto inteiro. Valores: opencode, claude, omp,
                         copilot, all (ou combinações separadas por vírgula
                         ou espaço, ex.: opencode,claude)
    --scope <valor>      Onde instalar, sem menu: project ou global. Sem a
                         flag, o instalador pergunta quando é chamado de
                         dentro de um repositório git, e cai em global fora
                         dele. É um ou outro — nunca os dois na mesma
                         execução
    --status, -s         Listar o que já está instalado, nos dois escopos, e
                         apontar skill que um mesmo harness enxerga em mais
                         de um destino. Só lê o disco: não instala nada
    --version, -v        Exibir a versão
    --help, -h           Exibir esta ajuda

  Escopo de projeto (--scope project):

    Os artefatos vão para a raiz do repositório — a partir do .git, subindo
    do diretório atual —, no diretório que cada harness varre lá dentro:

      opencode  .opencode/     claude   .claude/
      omp       .agents/       copilot  .github/

    O Copilot é o caso em que o nome muda: no usuário ele lê ~/.copilot, no
    repositório lê .github. E os overrides abaixo não valem neste escopo:
    eles nomeiam a base de usuário.

  Overrides de diretório (variáveis de ambiente, só no escopo global):
    OPENCODE_DIR         base do OpenCode (default ~/.config/opencode)
    CLAUDE_DIR           base do Claude Code (default ~/.claude)
    OMP_AGENTS_DIR       base de agent dirs do Oh My Pi (default ~/.agents)
    COPILOT_HOME         base do GitHub Copilot (default ~/.copilot); é a
                         variável da própria CLI, não uma inventada aqui
    NO_COLOR             desliga as cores da saída
    ACCESSIBLE           usa prompts simples, sem TUI, para leitor de tela

`
}
