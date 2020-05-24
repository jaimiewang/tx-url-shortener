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

func ShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	shortURL, err := model.FindShortURL(vars["code"])
	if err == sql.ErrNoRows {
		http.Error(w, NewAPIError("not found").Error(), http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	}

	err = util.WriteJson(w, shortURL)
	if err != nil {
		panic(err)
	}
}

type newURLRequest struct {
	URL string `json:"url"`
}

type newURLResponse struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

func NewShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestData := &newURLRequest{}
	err := json.NewDecoder(r.Body).Decode(requestData)
	if err != nil {
		http.Error(w, NewAPIErrorFromError(err).Error(), http.StatusBadRequest)
		return
	}

	originalURL, err := util.ValidateURL(requestData.URL)
	if err != nil {
		http.Error(w, NewAPIErrorFromError(err).Error(), http.StatusBadRequest)
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

	err = util.WriteJson(w, newURLResponse{
		Code: shortURL.Code,
		URL:  config.Config.ShortURLPrefix + "/" + shortURL.Code,
	})
	if err != nil {
		panic(err)
	}
}
