package apiv1

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

func APIError(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	util.WriteJsonResponse(w, apiError{Text: error})
}

var ErrEmptyAuthToken = errors.New("empty authorization token")
var ErrInvalidAuthToken = errors.New("invalid authorization token")
var ErrNotFound = errors.New("not found")

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimLeft(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			APIError(w, ErrEmptyAuthToken.Error(), http.StatusUnauthorized)
			return
		}

		ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM api_keys WHERE token=?", token)
		if err != nil {
			panic(err)
		}

		if ret == 0 {
			APIError(w, ErrInvalidAuthToken.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
