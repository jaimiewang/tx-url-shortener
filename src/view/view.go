package view

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.Config.NotFoundRedirect, http.StatusPermanentRedirect)
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(NotFound)
}

func ShortURLRedirectView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := model.ShortURL{}

	err := database.DbMap.SelectOne(&shortURL, "SELECT * FROM urls WHERE code=?", vars["code"])
	if err == sql.ErrNoRows {
		NotFound(w, r)
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
