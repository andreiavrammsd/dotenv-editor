package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/andreiavrammsd/dotenv-editor/env"
	"github.com/andreiavrammsd/dotenv-editor/handlers"
)

func main() {
	e := env.New()
	h := handlers.New(e)
	http.HandleFunc("/", h.Default)
	http.HandleFunc("/ui/", h.Static)
	http.HandleFunc("/env/current", h.GetCurrent)
	http.HandleFunc("/env/save", h.SaveAsFile)
	http.HandleFunc("/env/file", h.LoadFromFile)

	addr := flag.String("addr", ":8811", "-addr 127.0.0.1:8811")
	flag.Parse()

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
