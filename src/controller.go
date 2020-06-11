package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"tx-url-shortener/api"
	"tx-url-shortener/view"
)

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
