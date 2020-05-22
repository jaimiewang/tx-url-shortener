package view

import (
	"database/sql"
	"errors"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/url"
	"strings"
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
	var shortUrl model.ShortURL
	vars := mux.Vars(r)

	err := database.DbMap.SelectOne(&shortUrl, "SELECT * FROM urls WHERE code=?", vars["code"])
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		panic(err)
	}

	http.Redirect(w, r, shortUrl.Original, http.StatusPermanentRedirect)
}

func NewShortURLView(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalUrl, err := url.ParseRequestURI(r.FormValue("url"))
	if err != nil {
		util.RenderTemplate(w, "failed.html", map[string]interface{}{
			"err": err,
		})
		return
	}

	if originalUrl.Host == "" || originalUrl.Scheme == "" {
		util.RenderTemplate(w, "failed.html", map[string]interface{}{
			"err": errors.New("host and scheme cannot be empty"),
		})
		return
	}

	if !strings.HasSuffix(originalUrl.Path, "/") {
		originalUrl.Path += "/"

		remoteAddress, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remoteAddress = r.RemoteAddr
		}

		shortUrl := &model.ShortURL{
			Original:  originalUrl.String(),
			IPAddress: remoteAddress,
			Time:      time.Now(),
		}

		create, err := shortUrl.GenerateCode()
		if err != nil {
			panic(err)
		}

		if create {
			err = database.DbMap.Insert(shortUrl)
			if err != nil {
				panic(err)
			}
		}

		util.RenderTemplate(w, "success.html", map[string]interface{}{
			"shortUrlPrefix": config.Config.ShortURLPrefix,
			"shortUrl":       shortUrl,
			"request":        r,
		})
	}
}
