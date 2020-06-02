package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
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
)

func initRouter() *mux.Router {
	memoryAdapter, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(config.Config.CacheSize),
	)
	if err != nil {
		panic(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memoryAdapter),
		cache.ClientWithTTL(10*time.Minute),
	)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(handlers.ProxyHeaders)
	router.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})
	router.Use(handlers.RecoveryHandler())
	router.Use(cacheClient.Middleware)
	router.Use(api.AuthHandler)

	router.HandleFunc("/urls", api.NewShortURLEndpoint).Methods("PUT")
	router.HandleFunc("/urls/{code}", api.ShortURLEndpoint).Methods("GET")
	router.NotFoundHandler = api.NotFoundHandler()

	return router
}

func generateAPIKey() {
	apiKey := &model.APIKey{
		Time: time.Now().Unix(),
	}

	err := apiKey.GenerateToken()
	if err != nil {
		panic(err)
	}

	err = database.DbMap.Insert(apiKey)
	if err != nil {
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
		Addr:    config.Config.ListenAddress,
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
