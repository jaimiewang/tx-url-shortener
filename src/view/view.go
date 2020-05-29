package view

import (
	"database/sql"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
	"tx-url-shortener/util"
)

func IndexView(w http.ResponseWriter, r *http.Request) {
	util.RenderTemplate(w, "index.html", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func ShortURLView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := &model.ShortURL{}

	err := database.DbMap.SelectOne(shortURL, "SELECT * FROM urls WHERE code=?", vars["code"])
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

func NewShortURLView(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalURL, err := util.ValidateURL(r.FormValue("url"))
	if err != nil {
		util.RenderTemplate(w, "failed.html", map[string]interface{}{
			"err": err,
		})
		return
	}

	remoteAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteAddress = r.RemoteAddr
	}

	shortURL := &model.ShortURL{
		Original:  originalURL,
		IPAddress: remoteAddress,
		Time:      time.Now().UTC().Unix(),
	}

	err = shortURL.GenerateCode()
	if err != nil {
		panic(err)
	}

	err = database.DbMap.Insert(shortURL)
	if err != nil {
		panic(err)
	}

	util.RenderTemplate(w, "success.html", map[string]interface{}{
		"shortURLPrefix": config.Config.ShortURLPrefix,
		"shortURL":       shortURL,
		"request":        r,
	})
}
