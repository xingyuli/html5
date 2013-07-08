package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	index     *template.Template
	templates = make(map[string]*template.Template)
)

const (
	TEMPLATE_DIR = "./pages"
)

func init() {
	index = template.Must(template.ParseFiles("./index.html"))
	filepath.Walk(TEMPLATE_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			log.Println("Loading template:", path)
			templateName := path[strings.Index(path, string(os.PathSeparator))+1:]
			templates[templateName] = template.Must(template.ParseFiles(path))
		}
		return nil
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pages := []string{}
		for name, _ := range templates {
			pages = append(pages, name)
		}
		index.Execute(w, map[string]interface{}{"pages": pages})
	})

	http.HandleFunc("/html5", func(w http.ResponseWriter, r *http.Request) {
		pageParam := r.FormValue("page")
		if page, ok := templates[pageParam]; ok {
			page.Execute(w, nil)
		} else {
			http.NotFound(w, r)
		}
	})

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal("Failed to startup the server", err)
	}
}
