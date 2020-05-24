package apiv1

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"time"
	"tx-url-shortener/config"
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
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	shortURL, err := model.GetShortURL(vars["code"])
	if err == sql.ErrNoRows {
		APIError(w, ErrNotFound, http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	err = util.WriteJson(w, shortURLResponse{
		IPAddress: shortURL.IPAddress,
		Counter:   shortURL.Counter,
		Code:      shortURL.Code,
		Time:      shortURL.Time,
		Original:  shortURL.Original,
	})
	if err != nil {
		panic(err)
	}
}

type newShortURLRequest struct {
	URL string `json:"url"`
}

type newShortURLResponse struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

func NewShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestData := &newShortURLRequest{}
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		APIError(w, err, http.StatusBadRequest)
		return
	}

	originalURL, err := util.ValidateURL(requestData.URL)
	if err != nil {
		APIError(w, err, http.StatusBadRequest)
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

	err = model.SaveShortURL(shortURL)
	if err != nil {
		panic(err)
	}

	err = util.WriteJson(w, newShortURLResponse{
		Code: shortURL.Code,
		URL:  config.Config.ShortURLPrefix + "/" + shortURL.Code,
	})
	if err != nil {
		panic(err)
	}
}
