package main

import (
	"database/sql"
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

func initCacheClient() *cache.Client {
	memoryAdapter, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(config.Config.CacheSize),
	)
	if err != nil {
		panic(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memoryAdapter),
		cache.ClientWithTTL(3*time.Minute),
	)
	if err != nil {
		panic(err)
	}

	return cacheClient
}

func shortURLRedirectView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := model.ShortURL{}

	err := database.DbMap.SelectOne(&shortURL, "SELECT * FROM urls WHERE code=?", vars["code"])
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		panic(err)
	}

	shortURL.Counter++
	_, err = database.DbMap.Update(&shortURL)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, shortURL.Original, http.StatusPermanentRedirect)
}

func initRouter(cacheClient *cache.Client) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.Use(handlers.ProxyHeaders)
	router.Use(handlers.RecoveryHandler())
	router.Use(cacheClient.Middleware)

	basicRouter := router.PathPrefix("").Subrouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(api.AuthHandler)
	apiRouter.NotFoundHandler = api.NotFoundHandler()

	basicRouter.HandleFunc("/{code}", shortURLRedirectView)

	apiRouter.HandleFunc("/urls", api.NewShortURLEndpoint).Methods("PUT")
	apiRouter.HandleFunc("/urls/{code}", api.ShortURLEndpoint).Methods("GET")

	return router
}

func generateAPIKey() {
	apiKey := &model.APIKey{
		CreatedAt: time.Now().Unix(),
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

	cacheClient := initCacheClient()
	router := initRouter(cacheClient)
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
