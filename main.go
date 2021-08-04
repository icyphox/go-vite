package main

import (
	"fmt"
	"os"

	"git.icyphox.sh/vite/commands"
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

	if len(args) <= 1 {
		fmt.Println(helpStr)
		return
	}

	switch args[1] {
	case "init":
		if len(args) <= 2 {
			fmt.Println(helpStr)
			return
		}
		initPath := args[2]
		err := commands.Init(initPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: init: %+v\n", err)
		}

	case "build":
		err := commands.Build()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: build: %+v\n", err)
		}

	case "new":
		if len(args) <= 2 {
			fmt.Println(helpStr)
			return
		}
		newPath := args[2]
		err := commands.New(newPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: new: %+v\n", err)
		}
	}

}
