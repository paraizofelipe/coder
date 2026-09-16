package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"charm.land/huh/v2"
)

// errSemTerminal separa "o usuário cancelou" de "não havia como perguntar".
var errSemTerminal = errors.New(
	"sem terminal interativo: informe os destinos com --harness <lista>")

// form prepara o formulário com o terminal correto. O modo acessível desliga
// o TUI e imprime prompts simples, para leitor de tela.
func (u *ui) form(fields ...huh.Field) *huh.Form {
	return huh.NewForm(huh.NewGroup(fields...)).
		WithAccessible(os.Getenv("ACCESSIBLE") != "").
		WithInput(u.ttyIn).
		WithOutput(u.ttyOut)
}

// harnessLabel marca no menu o harness já presente na máquina. Ausência de
// marca é a informação sobre o resto: um rótulo em cada linha viraria ruído
// e pararia de destacar qualquer coisa.
func (u *ui) harnessLabel(h harness) string {
	label := h.name + " — " + h.label
	if h.detected() {
		label += u.c.yellow + " [instalado]" + u.c.reset
	}
	return label
}

// rotulosDeEscopo descreve as duas saídas do menu. Separado do formulário
// para que o texto tenha teste sem abrir terminal — e o rótulo de projeto
// nomeia a raiz pela mesma razão da marca [instalado]: a tela não pode
// descrever um lugar e a instalação escrever em outro.
func rotulosDeEscopo(raiz string) (projeto, global string) {
	return "projeto — " + raiz,
		"global — diretórios de usuário de cada harness"
}

// selectEscopo pergunta antes de tudo, porque é a resposta daqui que define
// a base de cada harness — e, com ela, o [instalado] da tela seguinte.
func (u *ui) selectEscopo(raiz string) (escopo, error) {
	if !u.interactive {
		return escopo{}, errSemTerminal
	}
	projeto, global := rotulosDeEscopo(raiz)

	// Global vem primeiro e é a opção destacada ao abrir: Enter mantém o
	// comportamento de sempre, e instalar no repositório é escolha ativa.
	var escolhido string
	field := huh.NewSelect[string]().
		Title("Escopo da instalação").
		Description("Um dos dois, nunca os dois na mesma execução.").
		Options(
			huh.NewOption(global, escopoGlobal),
			huh.NewOption(projeto, escopoProjeto),
		).
		Height(6).
		Value(&escolhido)

	if err := u.form(field).Run(); err != nil {
		return escopo{}, err
	}
	return resolveEscopoFlag(escolhido, raiz)
}

func (u *ui) selectHarnesses(e escopo) ([]harness, error) {
	if !u.interactive {
		return nil, errSemTerminal
	}
	options := make([]huh.Option[string], 0, len(harnesses))
	for _, h := range comEscopo(harnesses, e) {
		options = append(options, huh.NewOption(u.harnessLabel(h), h.name))
	}

	var chosen []string
	field := huh.NewMultiSelect[string]().
		Title("Harness(es) de destino").
		Description("Cada um recebe os artefatos no diretório que ele varre.").
		Options(options...).
		// Altura fixa: sem ela, um terminal que não negocia tamanho encolhe
		// o viewport e esconde opções — o usuário escolhe entre o que vê.
		Height(len(options) + 3).
		Value(&chosen)

	if err := u.form(field).Run(); err != nil {
		return nil, err
	}
	if len(chosen) == 0 {
		return nil, errors.New("nenhum harness selecionado")
	}
	selected := make(map[string]bool, len(chosen))
	for _, name := range chosen {
		selected[name] = true
	}
	return canonical(selected), nil
}

// overwritePrompt espelha as três saídas do install.sh — pular, substituir e
// substituir tudo — com "pular" em primeiro, que é a opção destacada ao abrir.
func (u *ui) overwritePrompt(label, dst string) (overwriteDecision, error) {
	if !u.interactive {
		return overwriteSkip, nil
	}
	var decision overwriteDecision
	field := huh.NewSelect[overwriteDecision]().
		Title(fmt.Sprintf("Substituir %s?", label)).
		Description(dst).
		Options(
			huh.NewOption("Não, pular", overwriteSkip),
			huh.NewOption("Sim, substituir", overwriteReplace),
			huh.NewOption("Sim, e todos os próximos conflitos", overwriteAll),
		).
		Height(6).
		Value(&decision)

	if err := u.form(field).Run(); err != nil {
		return overwriteSkip, err
	}
	return decision, nil
}

// artifactLabel marca em quais dos harnesses escolhidos o artefato já existe.
//
// A lista fica restrita aos destinos selecionados de propósito: saber que uma
// skill está instalada num harness que o usuário não escolheu não muda
// decisão nenhuma, e alongaria a linha sem informar.
func (u *ui) artifactLabel(nome string, targets []harness, instalado func(harness, string) bool) string {
	var onde []string
	for _, h := range targets {
		if instalado(h, nome) {
			onde = append(onde, h.name)
		}
	}
	if len(onde) == 0 {
		return nome
	}
	return nome + " - " + u.c.yellow + "[" + strings.Join(onde, ", ") + "]" + u.c.reset
}

// filtrarNaOrdem devolve os escolhidos na ordem do manifesto. Não confia na
// ordem que o formulário devolve: a ordem do fluxo é contrato, e é ela que
// aparece no log da instalação.
func filtrarNaOrdem(todos, escolhidos []string) []string {
	escolha := make(map[string]bool, len(escolhidos))
	for _, nome := range escolhidos {
		escolha[nome] = true
	}
	out := make([]string, 0, len(escolhidos))
	for _, nome := range todos {
		if escolha[nome] {
			out = append(out, nome)
		}
	}
	return out
}

// pickArtifacts abre um menu com tudo pré-marcado: Enter instala o conjunto
// inteiro, como antes desta tela existir, e desmarcar é que passa a ser a
// ação. O contrário obrigaria a marcar sete itens no caminho mais comum.
func (u *ui) pickArtifacts(titulo string, nomes []string, targets []harness,
	instalado func(harness, string) bool) ([]string, error) {
	if len(nomes) == 0 {
		return nil, nil
	}
	options := make([]huh.Option[string], 0, len(nomes))
	for _, nome := range nomes {
		options = append(options, huh.NewOption(u.artifactLabel(nome, targets, instalado), nome).Selected(true))
	}

	var escolhidos []string
	field := huh.NewMultiSelect[string]().
		Title(titulo).
		Description("[harness] marca onde o artefato já está instalado").
		Options(options...).
		Height(len(options) + 3).
		Value(&escolhidos)

	if err := u.form(field).Run(); err != nil {
		return nil, err
	}
	return filtrarNaOrdem(nomes, escolhidos), nil
}

// selectArtifacts filtra o manifesto pelo que o usuário quer instalar. Sem
// terminal devolve o manifesto inteiro: é o caminho de script, e lá a
// ausência de escolha significa "tudo", não "nada".
func (u *ui) selectArtifacts(man manifest, targets []harness) (manifest, error) {
	if !u.interactive {
		return man, nil
	}
	skills, err := u.pickArtifacts("Instalação de Skills", man.Skills, targets, harness.hasSkill)
	if err != nil {
		return manifest{}, err
	}
	// O menu de commands só abre quando algum destino os lê. Pedir a escolha
	// e não instalar nada seria cobrar uma decisão sem efeito.
	disponiveis := man.Commands
	if !aceitaCommands(targets) {
		disponiveis = nil
	}
	commands, err := u.pickArtifacts("Instalação de Commands", disponiveis, targets, harness.hasCommand)
	if err != nil {
		return manifest{}, err
	}
	if len(skills) == 0 && len(commands) == 0 {
		return manifest{}, errors.New("nenhum artefato selecionado")
	}
	return manifest{Skills: skills, Commands: commands}, nil
}
