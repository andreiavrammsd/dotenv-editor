package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/andreiavrammsd/dotenv-editor/env"
	"github.com/andreiavrammsd/dotenv-editor/handlers"
)

var (
	addr = flag.String("addr", ":8811", "-addr 127.0.0.1:8811")
)

func main() {
	flag.Parse()
	environment := env.New()
	h := handlers.New(environment)

	http.HandleFunc("/", h.Default)
	http.HandleFunc("/env/current", h.GetCurrent)
	http.HandleFunc("/env/save", h.SaveAsFile)
	http.HandleFunc("/env/file", h.LoadFromFile)

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
