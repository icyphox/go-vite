package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	helpStr := `usage: vite [options]

A simple and minimal static site generator.

options:
    init PATH       create vite project at PATH
    build           builds the current project
    new PATH        create a new markdown post
        `

    // TODO: make arg parsing less shit

    if len(args) <= 1 {
        fmt.Println(helpStr)
    }

	switch args[1] {
	case "init":
		initPath := args[2]
		viteInit(initPath)
	case "build":
		viteBuild()
	case "new":
		// newPath := args[2]
		// viteNew(newPath)
	}

}
