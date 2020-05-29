package apiv1

import (
	"database/sql"
	"encoding/json"
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
	Time      int64  `json:"time"`
	Original  string `json:"original"`
}

func ShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := &model.ShortURL{}

	err := database.DbMap.SelectOne(shortURL, "SELECT * FROM urls WHERE code=?", vars["code"])
	if err == sql.ErrNoRows {
		APIError(w, ErrNotFound.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	util.WriteJsonResponse(w, shortURLResponse{
		IPAddress: shortURL.IPAddress,
		Counter:   shortURL.Counter,
		Code:      shortURL.Code,
		Time:      shortURL.Time,
		Original:  shortURL.Original,
	})
}

type newShortURLRequest struct {
	URL string `json:"url"`
}

type newShortURLResponse struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

func NewShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	requestData := &newShortURLRequest{}
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		APIError(w, err.Error(), http.StatusBadRequest)
		return
	}

	originalURL, err := util.ValidateURL(requestData.URL)
	if err != nil {
		APIError(w, err.Error(), http.StatusBadRequest)
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

	util.WriteJsonResponse(w, newShortURLResponse{
		Code: shortURL.Code,
		URL:  config.Config.ShortURLPrefix + "/" + shortURL.Code,
	})
}
