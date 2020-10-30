package main

import (
	"os"
	"path/filepath"
)

func viteInit(path string) {
	paths := []string{"build", "pages", "static", "templates"}
	var dirPath string
	for _, p := range paths {
		dirPath = filepath.Join(path, p)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			printErr(err)
			return
		}
	}
	fp, _ := filepath.Abs(path)
	printMsg("created project:", fp)
}
