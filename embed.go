package main

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/BurntSushi/toml"
)

// content carrega o manifesto, as skills e os commands para dentro do
// binário. É o que torna a instalação atômica: sem rede, sem REPO_URL e sem
// download arquivo a arquivo. Conteúdo e instalador versionam juntos.
//
// O go:embed ignora arquivos iniciados por "." ou "_", o que já exclui
// qualquer sobra de editor sem precisar de filtro explícito.
//
//go:embed manifest.toml skills commands
var content embed.FS

const manifestPath = "manifest.toml"

// manifest é a lista do que instalar, na ordem do fluxo. Vive em TOML, e não
// em uma slice no código, para que adicionar uma skill seja editar um arquivo
// de dados — não recompilar uma decisão.
type manifest struct {
	Skills   []string `toml:"skills"`
	Commands []string `toml:"commands"`
}

// loadManifest lê e valida o manifesto embutido. Devolve erro em vez de
// entrar em pânico no init: manifesto quebrado é defeito de build, e quem
// roda o binário merece a mensagem do que está errado, não uma stack trace.
func loadManifest(fsys fs.FS) (manifest, error) {
	data, err := fs.ReadFile(fsys, manifestPath)
	if err != nil {
		return manifest{}, fmt.Errorf("manifesto embutido ausente: %w", err)
	}
	var m manifest
	if err := toml.Unmarshal(data, &m); err != nil {
		return manifest{}, fmt.Errorf("%s inválido: %w", manifestPath, err)
	}
	if len(m.Skills) == 0 {
		return manifest{}, fmt.Errorf("%s não lista nenhuma skill: não há o que instalar", manifestPath)
	}
	return m, nil
}
