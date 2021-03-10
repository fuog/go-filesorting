package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// file for the worker
type file struct {
	Locked bool
}

// FilePathWalkDir TODO: writing this Comment
func FilePathWalkDir(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != "." {
			return filepath.SkipDir
		}
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			fmt.Println(path)
		}

		return nil
	})
	return err
}
