package api

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
	"tx-url-shortener/config"
	"tx-url-shortener/database"
	"tx-url-shortener/model"
	"tx-url-shortener/util"
)

func ShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := &model.ShortURL{}

	err := database.DbMap.SelectOne(shortURL, "SELECT * FROM urls WHERE code=?", vars["code"])
	if err == sql.ErrNoRows {
		NotFound(w, r)
		return
	} else if err != nil {
		panic(err)
	}

	WriteAPIResponse(w, ShortURL{
		IPAddress: shortURL.IPAddress,
		Counter:   shortURL.Counter,
		Code:      shortURL.Code,
		CreatedAt: shortURL.CreatedAt,
		Original:  shortURL.Original,
		URL:       config.Config.ShortURLPrefix + "/" + shortURL.Code,
	})
}

func NewShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	requestData := &ShortenURLForm{}
	if err := ParseAPIForm(r, requestData); err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalURL, err := util.ValidateURL(requestData.URL)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	remoteAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteAddress = r.RemoteAddr
	}

	shortURL := &model.ShortURL{
		Original:  originalURL,
		IPAddress: remoteAddress,
		CreatedAt: time.Now().UTC().Unix(),
	}

	trans, err := database.DbMap.Begin()
	if err != nil {
		panic(err)
	}

	doubled, originalShortURL, err := shortURL.IsDoubled(trans)
	if err != nil {
		_ = trans.Rollback()
		panic(err)
	}

	statusCode := http.StatusCreated
	created := !doubled
	if created {
		err = shortURL.GenerateCode(trans)
		if err != nil {
			_ = trans.Rollback()
			panic(err)
		}

		err = trans.Insert(shortURL)
		if err != nil {
			_ = trans.Rollback()
			panic(err)
		}
	} else {
		shortURL = originalShortURL
		statusCode = http.StatusNotModified
	}

	err = trans.Commit()
	if err != nil {
		_ = trans.Rollback()
		panic(err)
	}

	w.WriteHeader(statusCode)
	WriteAPIResponse(w, ShortURL{
		IPAddress: shortURL.IPAddress,
		Counter:   shortURL.Counter,
		Code:      shortURL.Code,
		CreatedAt: shortURL.CreatedAt,
		Original:  shortURL.Original,
		URL:       config.Config.ShortURLPrefix + "/" + shortURL.Code,
	})
}
