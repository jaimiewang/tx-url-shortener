package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tx-url-shortener/database"
)

func ParseAPIForm(r *http.Request, i interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("not supported content-type: %s", contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		return err
	}

	return nil
}

func WriteAPIResponse(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(i)
	if err != nil {
		return
	}

	_, _ = w.Write(b)
}

func Error(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	WriteAPIResponse(w, map[string]interface{}{
		"error": error,
	})
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	Error(w, "not found", http.StatusNotFound)
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(NotFound)
}

var ErrEmptyAuthToken = errors.New("empty authorization token")
var ErrInvalidAuthToken = errors.New("invalid authorization token")

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
