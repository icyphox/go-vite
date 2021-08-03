package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

func Init(path string) error {
	paths := []string{"build", "pages", "static", "templates"}
	var dirPath string

	for _, p := range paths {
		dirPath = filepath.Join(path, p)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	fp, _ := filepath.Abs(path)
	fmt.Printf("vite: created project at %q\n", fp)
	return nil
}
