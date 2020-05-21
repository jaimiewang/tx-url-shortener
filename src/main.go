package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
	"tx-url-shortener/view"
)

func InitDatabaseTables() error {
	database.DbMap.AddTableWithName(model.ShortURL{}, "urls")

	err := database.DbMap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}

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

	err = InitDatabaseTables()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.Use(handlers.ProxyHeaders)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realPath := r.Header.Get("X-Real-Path")
			if realPath != "" {
				r.URL.RawPath = realPath
				r.URL.Path = url.PathEscape(realPath)
			}
			next.ServeHTTP(w, r)
		})
	})
	router.Use(csrf.Protect([]byte(config.Config.Secret), csrf.Secure(false)))

	router.HandleFunc("/", view.IndexView).Methods("GET")
	router.HandleFunc("/", view.NewShortURLView).Methods("POST")
	router.HandleFunc("/{code}", view.ShortURLView).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(config.Config.ListenAddress, router))
}
