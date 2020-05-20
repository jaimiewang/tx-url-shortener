package util

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates *template.Template

func init() {
	var templateFiles []string
	var err error

	err = filepath.Walk("templates/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	templates, err = template.ParseFiles(templateFiles...)
	if err != nil {
		panic(err)
	}
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}
