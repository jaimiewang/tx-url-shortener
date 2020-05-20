package view

import (
	"database/sql"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, shortUrl.Original, http.StatusPermanentRedirect)
}

func NewShortURLView(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(r.FormValue("url"))
	if err != nil || u.Host == "" || u.Scheme == "" {
		util.RenderTemplate(w, "failed.html", map[string]interface{}{"err": err})
		return
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	urlString := u.String()

	var shortUrl model.ShortURL
	err = database.DbMap.SelectOne(&shortUrl, "SELECT * FROM urls WHERE original=?", urlString)
	if err == sql.ErrNoRows {
		var urlCode string
		urlCodeLength := config.Conf.BaseURLLength

		urlsCount, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM urls")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var counter int64

		for {
			counter = 0
			for {
				if urlsCount >= 4 && counter >= urlsCount/4 {
					break
				}

				urlCode = util.RandomString(urlCodeLength)
				ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM urls WHERE code=?", urlCode)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if ret == 0 {
					break
				}

				counter++
			}

			if urlsCount >= 1 && counter == urlsCount {
				urlCodeLength += 1
				continue
			} else {
				break
			}
		}

		shortUrl = model.ShortURL{
			Original:  urlString,
			Code:      urlCode,
			IPAddress: r.RemoteAddr,
			Time:      time.Now(),
		}

		err = database.DbMap.Insert(&shortUrl)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.RenderTemplate(w, "success.html", map[string]interface{}{
		"shortUrl": shortUrl,
		"request":  r,
	})
}
