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
    init PATH                   create vite project at PATH
    build [--drafts]            builds the current project
    new PATH                    create a new markdown post
    serve [HOST:PORT]           serves the 'build' directory
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
		if err := commands.Init(initPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: init: %+v\n", err)
		}

	case "build":
		var drafts bool
		if len(args) > 2 && args[2] == "--drafts" {
			drafts = true
		}
		if err := commands.Build(drafts); err != nil {
			fmt.Fprintf(os.Stderr, "error: build: %+v\n", err)
		}

	case "new":
		if len(args) <= 2 {
			fmt.Println(helpStr)
			return
		}
		newPath := args[2]
		if err := commands.New(newPath); err != nil {
			fmt.Fprintf(os.Stderr, "error: new: %+v\n", err)
		}
	case "serve":
		var addr string
		if len(args) == 3 {
			addr = args[2]
		} else {
			addr = ":9191"
		}
		if err := commands.Serve(addr); err != nil {
			fmt.Fprintf(os.Stderr, "error: serve: %+v\n", err)
		}
	default:
		fmt.Println(helpStr)
	}

}
