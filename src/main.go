package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
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

	err = model.InitModels()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.Use(handlers.ProxyHeaders)
	router.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})
	router.Use(handlers.RecoveryHandler())
	router.Use(csrf.Protect([]byte(config.Config.Secret), csrf.Secure(false)))

	router.HandleFunc("/", view.IndexView).Methods("GET")
	router.HandleFunc("/", view.NewShortURLView).Methods("POST")
	router.HandleFunc("/{code}", view.ShortURLView).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	server := &http.Server{
		Addr:    config.Config.ListenAddress,
		Handler: router,
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-sc
	_ = server.Shutdown(nil)
	os.Exit(0)
}
