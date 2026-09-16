package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// escopo diz onde os artefatos caem. A raiz vazia é o escopo global — o
// comportamento que sempre existiu —, e é ela o zero value de propósito:
// harness que ninguém marcou continua apontando para o home.
//
// Um estado só, e não um par "modo + raiz": dois campos obrigados a
// concordar acabam discordando, e o dia em que discordassem o menu diria
// "projeto" enquanto a escrita cairia no home.
type escopo struct {
	raiz string // raiz do repositório; vazia no escopo global
}

func (e escopo) projeto() bool { return e.raiz != "" }

// resumo é o escopo em uma linha da saída. Só o de projeto nomeia um
// caminho: no global não há um único a citar, porque cada harness tem o seu
// — e o resumo do fim da instalação já lista os quatro.
func (e escopo) resumo() string {
	if e.projeto() {
		return "projeto (" + e.raiz + ")"
	}
	return "global"
}

// harness é um agente de destino. Cada um recebe os artefatos no diretório
// nativo que ele varre; o instalador não tenta unificar caminhos.
type harness struct {
	name       string // identificador aceito em --harness
	label      string // nome exibido no menu
	envVar     string // override do diretório base
	defaultDir string // relativo ao HOME, quando não há override

	// projectDir é o diretório que o harness varre dentro do repositório.
	// Não é derivável do defaultDir: em dois dos quatro o nome é outro.
	projectDir string

	// commands diz se o harness lê command em markdown. O Copilot não lê:
	// na CLI a skill já responde a /<nome>, e prompt file continua sendo
	// coisa de IDE. Instalar commands ali criaria arquivo que ninguém abre.
	commands bool

	// escopo é preenchido por comEscopo, depois da escolha. Guardá-lo dentro
	// do próprio harness é o que mantém skillsDir, detected e hasSkill
	// intactos: todos continuam lendo base(), e a marca do menu passa a
	// seguir o escopo sem que nenhum deles saiba que ele existe.
	escopo
}

// harnesses está em ordem canônica: ela define a sequência de instalação e a
// ordem do menu, independente de como o usuário digitou --harness.
var harnesses = []harness{
	{name: "opencode", label: "OpenCode", envVar: "OPENCODE_DIR",
		defaultDir: ".config/opencode", projectDir: ".opencode", commands: true},
	{name: "claude", label: "Claude Code", envVar: "CLAUDE_DIR",
		defaultDir: ".claude", projectDir: ".claude", commands: true},
	{name: "omp", label: "Oh My Pi", envVar: "OMP_AGENTS_DIR",
		defaultDir: ".agents", projectDir: ".agents", commands: true},
	// COPILOT_HOME quebra o padrão <NOME>_DIR de propósito: é a variável da
	// própria CLI do Copilot. Inventar uma segunda criaria dois lugares para
	// apontar o mesmo diretório, livres para discordar.
	//
	// E o diretório de projeto dele é .github, não .copilot: é lá que a CLI
	// procura skill versionada junto com o repositório.
	{name: "copilot", label: "GitHub Copilot", envVar: "COPILOT_HOME",
		defaultDir: ".copilot", projectDir: ".github", commands: false},
}

// comEscopo carimba o escopo escolhido nos harnesses. É o ponto único onde
// isso acontece: o resto do instalador recebe harness já apontando para o
// lugar certo e não precisa saber que houve escolha.
func comEscopo(hs []harness, e escopo) []harness {
	out := make([]harness, 0, len(hs))
	for _, h := range hs {
		h.escopo = e
		out = append(out, h)
	}
	return out
}

func (h harness) base() string {
	// No escopo de projeto o override de ambiente não entra: essas variáveis
	// nomeiam a base de usuário, e deixá-las vencer aqui faria o menu
	// prometer o repositório enquanto a instalação escreve no home.
	if h.projeto() {
		return filepath.Join(h.raiz, h.projectDir)
	}
	if dir := os.Getenv(h.envVar); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return h.defaultDir
	}
	return filepath.Join(home, h.defaultDir)
}

// detected diz se o harness parece presente na máquina. O sinal é o próprio
// diretório base existir — o mesmo caminho onde a instalação vai gravar, o
// que faz a marca responder exatamente à pergunta que importa no menu.
//
// Segue o override de ambiente de propósito: apontar CLAUDE_DIR para um
// sandbox deve mudar tanto o destino quanto a marca, senão o menu descreve
// um lugar e a instalação escreve em outro.
func (h harness) detected() bool {
	info, err := os.Stat(h.base())
	return err == nil && info.IsDir()
}

func (h harness) skillsDir() string { return filepath.Join(h.base(), "skills") }

// commandsDir devolve caminho vazio para harness que não tem commands, em vez
// de um diretório plausível que ninguém lê.
func (h harness) commandsDir() string {
	if !h.commands {
		return ""
	}
	return filepath.Join(h.base(), "commands")
}

// hasSkill e hasCommand dizem se o artefato já está instalado neste harness.
// São o que alimenta a marca do menu de seleção — e usam exatamente o mesmo
// caminho que a instalação vai escrever, para a marca não mentir.
func (h harness) hasSkill(name string) bool {
	return caminhoExiste(filepath.Join(h.skillsDir(), name))
}

func (h harness) hasCommand(name string) bool {
	if !h.commands {
		return false
	}
	return caminhoExiste(filepath.Join(h.commandsDir(), name+".md"))
}

// aceitaCommands diz se ao menos um dos destinos lê command em markdown.
func aceitaCommands(targets []harness) bool {
	for _, h := range targets {
		if h.commands {
			return true
		}
	}
	return false
}

func caminhoExiste(caminho string) bool {
	_, err := os.Lstat(caminho)
	return err == nil
}

func lookupHarness(name string) (harness, bool) {
	for _, h := range harnesses {
		if h.name == name {
			return h, true
		}
	}
	return harness{}, false
}

func harnessNames() []string {
	names := make([]string, 0, len(harnesses))
	for _, h := range harnesses {
		names = append(names, h.name)
	}
	return names
}

// canonical devolve os harnesses selecionados na ordem canônica, sem
// repetição — o que dispensa qualquer deduplicação no ponto de chamada.
func canonical(selected map[string]bool) []harness {
	out := make([]harness, 0, len(selected))
	for _, h := range harnesses {
		if selected[h.name] {
			out = append(out, h)
		}
	}
	return out
}

// resolveHarnessFlag traduz o valor de --harness em harnesses. Aceita "all" e
// listas separadas por vírgula ou espaço. É pura de propósito: não lê
// ambiente, não escreve saída e não toca no disco, então o teste cobre a
// tabela inteira de entradas sem preparar nada.
func resolveHarnessFlag(input string) ([]harness, error) {
	fields := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	})
	if len(fields) == 0 {
		return nil, fmt.Errorf("--harness exige ao menos um valor. Use: %s, all",
			strings.Join(harnessNames(), ", "))
	}

	selected := make(map[string]bool, len(fields))
	for _, field := range fields {
		token := strings.ToLower(field)
		if token == "all" {
			return append([]harness(nil), harnesses...), nil
		}
		if _, ok := lookupHarness(token); !ok {
			return nil, fmt.Errorf("harness inválido: %q. Use: %s, all",
				field, strings.Join(harnessNames(), ", "))
		}
		selected[token] = true
	}
	return canonical(selected), nil
}
