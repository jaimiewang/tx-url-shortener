package apiv1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"tx-url-shortener/model"
)

type apiError struct {
	Text string `json:"error"`
}

func APIError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(apiError{Text: err.Error()})
	http.Error(w, string(b), status)
}

var ErrEmptyAuthToken = errors.New("empty authorization token")
var ErrInvalidAuthToken = errors.New("invalid authorization token")
var ErrNotFound = errors.New("not found")

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimLeft(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			APIError(w, ErrEmptyAuthToken, http.StatusUnauthorized)
			return
		}

		apiKey, err := model.GetAPIKey(token)
		if err != nil {
			panic(err)
		}

		if apiKey == nil {
			APIError(w, ErrInvalidAuthToken, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
