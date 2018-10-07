package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/andreiavrammsd/dotenv-editor/env"
)

var (
	environment env.Env
	addr        = flag.String("addr", ":8811", "-addr 127.0.0.1:8811")
)

func main() {
	flag.Parse()
	environment = env.New()

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/env/current", currentEnvHandler)
	http.HandleFunc("/env/save", saveEnvHandler)
	http.HandleFunc("/env/file", fileEnvHandler)

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
