package main

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const (
	dirPerm  os.FileMode = 0o755
	filePerm os.FileMode = 0o644
)

// overwriteDecision é a resposta de quem está no terminal diante de um
// arquivo que já existe.
type overwriteDecision int

const (
	overwriteSkip overwriteDecision = iota
	overwriteReplace
	overwriteAll
)

// promptFunc é a seam do instalador: todo o resto é filesystem puro, então
// trocar esta função por uma resposta fixa torna a instalação inteira
// testável sem terminal.
type promptFunc func(label, dst string) (overwriteDecision, error)

type installer struct {
	content fs.FS
	man     manifest
	ui      *ui
	force   bool
	dryRun  bool
	prompt  promptFunc

	overwriteAll bool // fica ligado depois da resposta "todos"
}

// checkOverwrite decide se o destino pode ser escrito. --force e a resposta
// "todos" desligam a pergunta daí em diante; sem terminal, o conflito é
// pulado em vez de travar esperando uma resposta que não virá.
func (in *installer) checkOverwrite(dst, label string) (bool, error) {
	if _, err := os.Lstat(dst); err != nil {
		return true, nil
	}
	in.ui.warn("Já existe: %s", dst)
	if in.force || in.overwriteAll {
		return true, nil
	}
	if in.prompt == nil {
		in.ui.skipped()
		return false, nil
	}

	decision, err := in.prompt(label, dst)
	if err != nil {
		return false, err
	}
	switch decision {
	case overwriteAll:
		in.overwriteAll = true
		in.ui.ok("Todos os próximos conflitos serão sobrescritos.")
		return true, nil
	case overwriteReplace:
		return true, nil
	default:
		in.ui.skipped()
		return false, nil
	}
}

// ensureDir concentra a criação de diretório para que a simulação tenha um
// ponto único de saída, em vez de um `if` espalhado por cada chamada.
func (in *installer) ensureDir(dir string) error {
	if in.dryRun {
		return nil
	}
	return os.MkdirAll(dir, dirPerm)
}

func (in *installer) installSkills(h harness) error {
	// Seção vazia não cria diretório: selecionar zero skills não deve
	// deixar um ~/.claude/skills vazio para trás.
	if len(in.man.Skills) == 0 {
		return nil
	}
	dir := h.skillsDir()
	if err := in.ensureDir(dir); err != nil {
		return err
	}
	in.ui.info("Instalando skills em %s", dir)
	in.ui.blank()

	for _, name := range in.man.Skills {
		in.ui.item(name)
		dst := filepath.Join(dir, name)
		write, err := in.checkOverwrite(dst, name)
		if err != nil {
			return err
		}
		if !write {
			continue
		}
		origem := path.Join("skills", name)
		if in.dryRun {
			total, err := in.countTree(origem)
			if err != nil {
				return err
			}
			in.ui.simulated("%d arquivos", total)
			continue
		}
		// dst é sempre skillsDir + um nome do manifesto embutido, nunca
		// entrada do usuário — o remove não alcança caminho arbitrário.
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
		if err := in.copyTree(origem, dst); err != nil {
			return err
		}
		in.ui.installed()
	}
	in.ui.blank()
	return nil
}

func (in *installer) installCommands(h harness) error {
	if len(in.man.Commands) == 0 {
		return nil
	}
	// Destino sem commands merece uma linha, não silêncio: sem ela, quem
	// pediu cinco commands e viu zero conclui que a instalação falhou.
	if !h.commands {
		in.ui.info("%s não lê commands: a skill já responde a /<nome>", h.name)
		in.ui.blank()
		return nil
	}
	dir := h.commandsDir()
	if err := in.ensureDir(dir); err != nil {
		return err
	}
	in.ui.info("Instalando commands em %s", dir)
	in.ui.blank()

	for _, name := range in.man.Commands {
		in.ui.item(name)
		dst := filepath.Join(dir, name+".md")
		write, err := in.checkOverwrite(dst, name)
		if err != nil {
			return err
		}
		if !write {
			continue
		}
		data, err := assemble(in.content, path.Join("commands", name), h.name)
		if err != nil {
			return err
		}
		// A montagem roda mesmo em simulação: é ela que detecta um .yml
		// faltando para o harness, que é o defeito que só apareceria no
		// dia da instalação real.
		if in.dryRun {
			in.ui.simulated("%d bytes", len(data))
			continue
		}
		if err := os.WriteFile(dst, data, filePerm); err != nil {
			return err
		}
		in.ui.installed()
	}
	in.ui.blank()
	return nil
}

func (in *installer) run(targets []harness) error {
	for _, h := range targets {
		in.ui.info("Instalando para: %s", in.ui.bold(h.name))
		in.ui.blank()
		if err := in.installSkills(h); err != nil {
			return err
		}
		if err := in.installCommands(h); err != nil {
			return err
		}
		in.ui.ok("Concluído: %s", h.name)
		in.ui.blank()
	}
	return nil
}
