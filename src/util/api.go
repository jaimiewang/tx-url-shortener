package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func RandomToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func ParseJsonForm(r *http.Request, i interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("invalid content type: %s", contentType)
	}

	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		return err
	}

	return nil
}

func WriteJsonResponse(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(i)
	if err != nil {
		return
	}

	_, _ = w.Write(b)
	return
}
