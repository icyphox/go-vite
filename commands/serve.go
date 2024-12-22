package commands

import (
	"fmt"
	"net/http"
)

func Serve(addr string) error {
	fs := http.FileServer(http.Dir("./build"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	fmt.Printf("vite: serving on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return err
	}
	return nil
}
