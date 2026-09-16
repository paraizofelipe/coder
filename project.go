package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// raizDoProjeto sobe a partir de dir até achar um .git e devolve o diretório
// que o contém. Devolve a raiz, e não o cwd, porque é dela que os quatro
// harnesses partem no walk-up: instalar no cwd deixaria um .claude/ perdido
// no meio da árvore, que nenhum deles varre.
//
// Aceita .git como arquivo além de diretório — worktree e submódulo gravam um
// arquivo apontando para o diretório real, e uma checagem que exigisse
// diretório não enxergaria nenhum dos dois.
//
// É separada do cwd de propósito: recebendo o diretório por parâmetro, o
// teste cobre a subida inteira sem trocar o diretório do processo.
func raizDoProjeto(dir string) (string, bool) {
	for {
		if caminhoExiste(filepath.Join(dir, ".git")) {
			return dir, true
		}
		pai := filepath.Dir(dir)
		// filepath.Dir devolve o próprio caminho na raiz do filesystem: é o
		// que encerra a busca em vez de girar para sempre.
		if pai == dir {
			return "", false
		}
		dir = pai
	}
}

// raizDoProjetoAtual é raizDoProjeto a partir de onde o instalador foi
// chamado. Diretório de trabalho ilegível conta como "não há projeto", que é
// o mesmo desfecho de rodar fora de um repositório.
func raizDoProjetoAtual() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}
	return raizDoProjeto(dir)
}

// Valores aceitos em --scope. Ficam em inglês como o resto das flags: a
// ajuda explica em português, mas a interface de linha de comando não troca
// de idioma no meio.
const (
	escopoProjeto = "project"
	escopoGlobal  = "global"
)

// resolveEscopoFlag traduz o valor de --scope em escopo. A raiz vazia é a
// resposta da detecção para "não estamos em um repositório".
//
// A flag vence a detecção nos dois sentidos: --scope global não é arrastado
// para o projeto por estar dentro de um repositório, e --scope project não é
// silenciosamente rebaixado para global por não estar em um.
func resolveEscopoFlag(valor, raiz string) (escopo, error) {
	switch strings.ToLower(strings.TrimSpace(valor)) {
	case escopoGlobal:
		return escopo{}, nil
	case escopoProjeto:
		if raiz == "" {
			return escopo{}, fmt.Errorf(
				"--scope %s exige um repositório git: nenhum .git foi encontrado a partir daqui",
				escopoProjeto)
		}
		return escopo{raiz: raiz}, nil
	default:
		return escopo{}, fmt.Errorf("escopo inválido: %q. Use: %s, %s",
			valor, escopoProjeto, escopoGlobal)
	}
}
