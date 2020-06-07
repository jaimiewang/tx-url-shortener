package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tx-url-shortener/api"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
	"tx-url-shortener/view"
)

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(handlers.ProxyHeaders)
	router.Use(handlers.RecoveryHandler())

	viewRouter := router.PathPrefix("").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.Use(api.AuthHandler)
	apiRouter.NotFoundHandler = api.NotFoundHandler()

	viewRouter.HandleFunc("/{code}", view.ShortURLRedirectView)
	apiRouter.HandleFunc("/urls", api.NewShortURLEndpoint).Methods("PUT")
	apiRouter.HandleFunc("/urls/{code}", api.ShortURLEndpoint).Methods("GET")

	return router
}

func generateAPIKey() {
	apiKey := &model.APIKey{
		CreatedAt: time.Now().Unix(),
	}

	trans, err := database.DbMap.Begin()
	if err != nil {
		panic(err)
	}

	err = apiKey.GenerateToken(trans)
	if err != nil {
		_ = trans.Rollback()
		panic(err)
	}

	err = trans.Insert(apiKey)
	if err != nil {
		_ = trans.Rollback()
		panic(err)
	}

	err = trans.Commit()
	if err != nil {
		_ = trans.Rollback()
		panic(err)
	}

	fmt.Printf("Your new API key is: %s\n", apiKey.Token)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	configFileFlag := flag.String("config", "config.yml", "Configuration file.")
	generateAPIKeyFlag := flag.Bool("generate-api-key", false, "")
	flag.Parse()

	err := config.LoadConfig(*configFileFlag)
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

	if *generateAPIKeyFlag {
		generateAPIKey()
		return
	}

	router := initRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	<-sc
	_ = server.Shutdown(nil)
}
