package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
)

// assemble monta o arquivo final do command:
//
//	---
//	<conteúdo de <harness>.yml>
//	---
//
//	<conteúdo de body.md>
//
// A quebra de linha após o YAML é garantida aqui: um .yml sem newline final
// colaria o delimitador na última chave e produziria frontmatter inválido.
func assemble(src fs.FS, dir, harnessName string) ([]byte, error) {
	front, err := fs.ReadFile(src, path.Join(dir, harnessName+".yml"))
	if err != nil {
		return nil, fmt.Errorf("frontmatter de %s: %w", harnessName, err)
	}
	body, err := fs.ReadFile(src, path.Join(dir, "body.md"))
	if err != nil {
		return nil, fmt.Errorf("corpo de %s: %w", dir, err)
	}

	var out bytes.Buffer
	out.WriteString("---\n")
	out.Write(front)
	if !bytes.HasSuffix(front, []byte("\n")) {
		out.WriteString("\n")
	}
	out.WriteString("---\n\n")
	out.Write(body)
	return out.Bytes(), nil
}
