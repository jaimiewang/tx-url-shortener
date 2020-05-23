package apiv1

import (
	"encoding/json"
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

type newUrlRequest struct {
	URL string `json:"url"`
}

type newUrlResponse struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

func NewShortURLEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := &newUrlRequest{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, NewAPIErrorFromError(err).Error(), http.StatusBadRequest)
		return
	}

	originalUrl, err := url.ParseRequestURI(data.URL)
	if err != nil {
		http.Error(w, NewAPIErrorFromError(err).Error(), http.StatusBadRequest)
		return
	}

	if originalUrl.Host == "" || originalUrl.Scheme == "" {
		http.Error(w, NewAPIError("").Error(), http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(originalUrl.Path, "/") {
		originalUrl.Path += "/"
	}

	remoteAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteAddress = r.RemoteAddr
	}

	shortUrl := &model.ShortURL{
		Original:  originalUrl.String(),
		IPAddress: remoteAddress,
		Time:      time.Now().UTC().Unix(),
	}

	err = shortUrl.GenerateCode()
	if err != nil {
		panic(err)
	}

	err = database.DbMap.Insert(shortUrl)
	if err != nil {
		panic(err)
	}

	err = util.WriteJson(w, newUrlResponse{
		Code: shortUrl.Code,
		URL:  config.Config.ShortURLPrefix + "/" + shortUrl.Code,
	})
	if err != nil {
		panic(err)
	}
}
