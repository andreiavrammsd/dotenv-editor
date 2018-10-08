package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"strconv"

	"io"

	"encoding/base64"

	"github.com/andreiavrammsd/dotenv-editor/env"
)

// Handlers for HTTP calls
type Handlers struct {
	env env.Env
}

// GetCurrent loads the env list from
// the current machine and generates a list
func (h Handlers) GetCurrent(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(h.env.Current())
	if err != nil {
		log.Println(err)
	}
}

// SaveAsFile generates a dotenv file based on submitted vars list
func (h Handlers) SaveAsFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	vars := make(map[string]env.Variable)
	err := json.Unmarshal([]byte(r.Form.Get("data")), &vars)
	if err != nil {
		log.Fatal(err)
	}

	out := ""
	src := r.Form.Get("src")
	if src != "" {
		out = h.env.Update(src, vars)
	} else {
		out = h.env.ToFile(vars)
	}

	content := []byte(out)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=\"env\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))

	if _, err := w.Write(content); err != nil {
		log.Println(err)
	}
}

// LoadFromFile generates the env list from an uploaded file
func (h Handlers) LoadFromFile(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, r.ContentLength)
	_, err := r.Body.Read(data)
	if err != nil && err != io.EOF {
		log.Println(err)
		return
	}

	defer func() {
		if e := r.Body.Close(); e != nil {
			log.Println(e)
		}
	}()

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err = encoder.Encode(h.env.FromInput(string(data)))
	if err != nil {
		log.Println(err)
		return
	}
}

// Default sets up the UI
func (Handlers) Default(w http.ResponseWriter, _ *http.Request) {
	html := "ui/index.html"

	content, err := Asset(html)
	if err != nil {
		log.Println(err)
		return
	}

	funcMap := template.FuncMap{
		"src": func(b []byte) template.Srcset {
			return template.Srcset(b)
		},
		"css": func(b []byte) template.CSS {
			return template.CSS(b)
		},
		"js": func(b []byte) template.JS {
			return template.JS(b)
		},
	}

	favicon, err := Asset("ui/favicon.png")
	if err != nil {
		log.Println(err)
		return
	}
	faviconEncoded := make([]byte, base64.StdEncoding.EncodedLen(len(favicon)))
	base64.StdEncoding.Encode(faviconEncoded, favicon)

	css, err := Asset("ui/style.css")
	if err != nil {
		log.Println(err)
		return
	}

	js, err := Asset("ui/script.js")
	if err != nil {
		log.Println(err)
		return
	}

	data := struct {
		Favicon []byte
		CSS     []byte
		JS      []byte
	}{
		faviconEncoded,
		css,
		js,
	}

	t, err := template.New(html).Funcs(funcMap).Parse(string(content))
	if err != nil {
		log.Println(err)
		return
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}

// New initializes the handlers functions
func New(env env.Env) Handlers {
	return Handlers{env}
}
