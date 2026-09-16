package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// copyTree grava uma subárvore do conteúdo embutido no disco, preservando
// references/.
func (in *installer) copyTree(srcDir, dstDir string) error {
	return fs.WalkDir(in.content, srcDir, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(strings.TrimPrefix(p, srcDir), "/")
		target := filepath.Join(dstDir, filepath.FromSlash(rel))
		if entry.IsDir() {
			return os.MkdirAll(target, dirPerm)
		}
		data, err := fs.ReadFile(in.content, p)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, filePerm)
	})
}

// countTree conta o que copyTree gravaria. Em simulação é o número que dá
// utilidade à saída: sem ele, o modo diria apenas que não fez nada.
func (in *installer) countTree(srcDir string) (int, error) {
	total := 0
	err := fs.WalkDir(in.content, srcDir, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			total++
		}
		return nil
	})
	return total, err
}
