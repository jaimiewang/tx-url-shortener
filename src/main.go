package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/view"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	err := config.LoadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	err = database.InitDatabase()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.Use(csrf.Protect([]byte(config.Conf.Secret), csrf.Secure(false)))
	router.Use(handlers.ProxyHeaders)

	router.HandleFunc("/", view.IndexView).Methods("GET")
	router.HandleFunc("/", view.NewShortURLView).Methods("POST")
	router.HandleFunc("/{code}", view.ShortURLView).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(config.Conf.ListenAddress, router))
}
