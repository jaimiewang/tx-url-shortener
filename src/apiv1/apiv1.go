package apiv1

import (
	"encoding/json"
	"net/http"
	"strings"
	"tx-url-shortener/database"
)

type apiError struct {
	Text string `json:"error"`
}

func (apiErr apiError) Error() string {
	bytes, _ := json.Marshal(apiErr)
	return string(bytes)
}

func NewAPIError(text string) error {
	return apiError{Text: text}
}

func NewAPIErrorFromError(err error) error {
	return NewAPIError(err.Error())
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimLeft(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, NewAPIError("empty authorization token").Error(), http.StatusUnauthorized)
			return
		}

		ret, err := database.DbMap.SelectInt("SELECT COUNT(*) FROM api_keys WHERE token=?", token)
		if err != nil {
			panic(err)
		}

		if ret == 0 {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, NewAPIError("invalid authorization token").Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
