package commands

import (
	"fmt"
	"log"
	"net/http"
)

func requestLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		log.Printf("%s\t%s", r.Method, r.URL.Path)
	})
}

func Serve(addr string) error {
	fs := http.FileServer(http.Dir("./build"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	fmt.Printf("vite: serving on %s\n", addr)
	if err := http.ListenAndServe(addr, requestLog(mux)); err != nil {
		return err
	}
	return nil
}
