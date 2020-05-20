package util

import (
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func init() {
	var err error

	templates, err = template.ParseGlob("templates/*")
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
