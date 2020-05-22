package util

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		panic(err)
	}
}
