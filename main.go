package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/csrf"
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
	"tx-url-shortener/apiv1"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
	"tx-url-shortener/view"
)

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
		os.Exit(0)
	}

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

	viewRouter := router.PathPrefix("/").Subrouter()
	apiv1Router := router.PathPrefix("/api/v1").Subrouter()
	staticRouter := router.PathPrefix("/static").Subrouter()

	viewRouter.Use(csrf.Protect([]byte(config.Config.Secret), csrf.Secure(false)))
	apiv1Router.Use(apiv1.AuthHandler)

	viewRouter.HandleFunc("/", view.IndexView).Methods("GET")
	viewRouter.HandleFunc("/", view.NewShortURLView).Methods("POST")
	viewRouter.HandleFunc("/{code}", view.ShortURLView).Methods("GET")

	apiv1Router.HandleFunc("/urls", apiv1.NewShortURLEndpoint).Methods("PUT")
	apiv1Router.HandleFunc("/urls/{code}", apiv1.ShortURLEndpoint).Methods("GET")

	staticRouter.PathPrefix("").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

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
	os.Exit(0)
}
