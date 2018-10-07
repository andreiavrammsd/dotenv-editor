package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"strconv"

	"io"

	"path/filepath"
	"strings"

	"github.com/andreiavrammsd/dotenv-editor/env"
)

func currentEnvHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(environment.Current())
	if err != nil {
		log.Println(err)
	}
}

func saveEnvHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	vars := make(map[string]env.Variable)
	err := json.Unmarshal([]byte(r.Form.Get("data")), &vars)
	if err != nil {
		log.Fatal(err)
	}

	out := ""
	src := r.Form.Get("src")
	if src != "" {
		out = environment.Update(src, vars)
	} else {
		out = environment.ToFile(vars)
	}

	content := []byte(out)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=\"env\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))

	if _, err := w.Write(content); err != nil {
		log.Println(err)
	}
}

func fileEnvHandler(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, r.ContentLength)
	_, err := r.Body.Read(data)
	if err != nil && err != io.EOF {
		log.Println(err)
		return
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err = encoder.Encode(environment.FromInput(string(data)))
	if err != nil {
		log.Println(err)
		return
	}
}

func defaultHandler(w http.ResponseWriter, _ *http.Request) {
	html := "ui/index.html"

	content, err := Asset(html)
	if err != nil {
		log.Println(err)
		return
	}

	t, err := template.New(html).Parse(string(content))
	if err != nil {
		log.Println(err)
		return
	}

	if err := t.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func faviconHandler(w http.ResponseWriter, _ *http.Request) {
	content, err := Asset("ui/favicon.png")
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if _, err := w.Write(content); err != nil {
		log.Println(err)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	ext := strings.TrimPrefix(filepath.Ext(path), ".")

	content, err := Asset(path)
	if err != nil {
		log.Println(err)
		return
	}

	contentTypes := map[string]string{
		"css": "text/css",
		"js":  "text/javascript",
	}

	w.Header().Set("Content-Type", contentTypes[ext])
	if _, err := w.Write(content); err != nil {
		log.Println(err)
	}
}
