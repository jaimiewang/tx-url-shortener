package api

import (
	"errors"
	"net/http"
	"strings"
	"tx-url-shortener/database"
	"tx-url-shortener/util"
)

type apiError struct {
	Text string `json:"error"`
}

func Error(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	util.WriteJsonResponse(w, apiError{Text: error})
}

var ErrEmptyAuthToken = errors.New("empty authorization token")
var ErrInvalidAuthToken = errors.New("invalid authorization token")

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, "not found", http.StatusNotFound)
	})
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimLeft(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			Error(w, ErrEmptyAuthToken.Error(), http.StatusUnauthorized)
			return
		}

		ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM api_keys WHERE token=?", token)
		if err != nil {
			panic(err)
		}

		if ret == 0 {
			Error(w, ErrInvalidAuthToken.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
