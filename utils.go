package main

import (
	"fmt"
	"os"
	"strings"
)

func printMsg(s ...string) {
	fmt.Println("vite:", strings.Join(s, " "))
}

func printErr(e error) {
	fmt.Fprintln(os.Stderr, "error:", e)
}
