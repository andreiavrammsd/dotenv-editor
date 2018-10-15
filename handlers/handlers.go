package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"strconv"

	"io"

	"encoding/base64"

	"errors"

	"github.com/andreiavrammsd/dotenv-editor/env"
)

// Handlers for HTTP calls
type Handlers struct {
	env env.Env
}

// GetCurrent loads the env list from
// the current machine and generates a list
func (h Handlers) GetCurrent(w http.ResponseWriter, _ *http.Request) {
	data, err := json.Marshal(h.env.Current())
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	n, err := w.Write(data)
	if n != len(data) {
		log.Println(errors.New("error writing response"))
		return
	}
	if err != nil {
		log.Println(err)
	}
}

// SaveAsFile generates a dotenv file based on submitted vars list
func (h Handlers) SaveAsFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	var vars []env.Variable
	err := json.Unmarshal([]byte(r.Form.Get("data")), &vars)
	if err != nil {
		log.Fatal(err)
	}

	out := ""
	src := r.Form.Get("src")
	if src != "" {
		out = h.env.Sync(src, vars)
	} else {
		out = h.env.ToString(vars)
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
	input := make([]byte, r.ContentLength)
	n, err := r.Body.Read(input)
	if err != nil && err != io.EOF {
		log.Println(err)
		return
	}
	if int64(n) != r.ContentLength {
		log.Println(errors.New("error reading request body"))
		return
	}

	defer func() {
		if e := r.Body.Close(); e != nil {
			log.Println(e)
		}
	}()

	data, err := json.Marshal(h.env.FromInput(string(input)))
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	n, err = w.Write(data)
	if n != len(data) {
		log.Println(errors.New("error writing response"))
		return
	}
	if err != nil {
		log.Println(err)
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
	}

	favicon, err := Asset("ui/favicon.png")
	if err != nil {
		log.Println(err)
		return
	}
	faviconEncoded := make([]byte, base64.StdEncoding.EncodedLen(len(favicon)))
	base64.StdEncoding.Encode(faviconEncoded, favicon)

	data := struct {
		Favicon []byte
	}{
		faviconEncoded,
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

// Static files
func (Handlers) Static(w http.ResponseWriter, r *http.Request) {
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

// New initializes the handlers functions
func New(e env.Env) Handlers {
	return Handlers{e}
}
