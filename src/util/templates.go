package util

import (
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template = template.Must(template.ParseGlob("templates/*"))

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}
