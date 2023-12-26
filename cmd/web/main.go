package main

import (
	"log"
	"net/http"

	"github.com/akyrey/snippetbox/internal"
)

func main() {
	config := internal.Config{}
	config.Parse()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(config.StaticDir))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("starting server on %s\n", config.Addr)

	err := http.ListenAndServe(config.Addr, mux)
	log.Fatal(err)
}
