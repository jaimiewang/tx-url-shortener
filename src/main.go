package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
	"tx-url-shortener/api"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
	"tx-url-shortener/view"
)

func Setup(configFilename string) error {
	rand.Seed(time.Now().UTC().UnixNano())

	err := config.LoadConfig(configFilename)
	if err != nil {
		return err
	}

	err = database.InitDatabase()
	if err != nil {
		return err
	}

	err = model.InitModels()
	if err != nil {
		return err
	}

	return nil
}

func InitRouter() *mux.Router {
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

func GenerateAPIKey() (string, error) {
	apiKey := &model.APIKey{
		CreatedAt: time.Now().Unix(),
	}

	trans, err := database.DbMap.Begin()
	if err != nil {
		return "", err
	}

	err = apiKey.GenerateToken(trans)
	if err != nil {
		_ = trans.Rollback()
		return "", err
	}

	err = trans.Insert(apiKey)
	if err != nil {
		_ = trans.Rollback()
		return "", err
	}

	err = trans.Commit()
	if err != nil {
		_ = trans.Rollback()
		return "", err
	}

	return apiKey.Token, nil
}

func main() {
	configFilename := flag.String("config", "config.yml", "Configuration file.")
	generateAPIKey := flag.Bool("generate-api-key", false, "")
	flag.Parse()

	err := Setup(*configFilename)
	if err != nil {
		panic(err)
	}

	if *generateAPIKey {
		apiKey, err := GenerateAPIKey()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Your new API key is: %s\n", apiKey)
		return
	}

	router := InitRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
