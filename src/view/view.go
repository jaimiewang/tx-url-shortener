package view

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
)

func ShortURLRedirectView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := &model.ShortURL{}

	trans, err := database.DbMap.Begin()
	if err != nil {
		panic(err)
	}

	err = trans.SelectOne(shortURL, "SELECT * FROM urls WHERE code=?", vars["code"])
	if err == sql.ErrNoRows {
		_ = trans.Rollback()
		http.NotFound(w, r)
		return
	} else if err != nil {
		_ = trans.Rollback()
		panic(err)
	}

	shortURL.Counter++
	_, err = trans.Update(shortURL)
	if err != nil {
		_ = trans.Rollback()
		panic(err)
	}

	err = trans.Commit()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, shortURL.Original, http.StatusPermanentRedirect)
}
