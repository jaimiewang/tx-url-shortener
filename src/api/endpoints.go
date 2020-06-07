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

type shortURLResponse struct {
	IPAddress string `json:"ip_address"`
	Counter   int64  `json:"counter"`
	Code      string `json:"code"`
	CreatedAt int64  `json:"created_at"`
	Original  string `json:"original"`
}

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

	WriteAPIResponse(w, shortURLResponse{
		IPAddress: shortURL.IPAddress,
		Counter:   shortURL.Counter,
		Code:      shortURL.Code,
		CreatedAt: shortURL.CreatedAt,
		Original:  shortURL.Original,
	})
}

type newShortURLRequest struct {
	URL string `json:"url"`
}

type newShortURLResponse struct {
	Code    string `json:"code"`
	URL     string `json:"url"`
	Created bool   `json:"created"`
}

func NewShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	requestData := &newShortURLRequest{}
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

	create := !doubled
	if create {
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
	}

	err = trans.Commit()
	if err != nil {
		panic(err)
	}

	WriteAPIResponse(w, newShortURLResponse{
		Code:    shortURL.Code,
		URL:     config.Config.ShortURLPrefix + "/" + shortURL.Code,
		Created: create,
	})
}
