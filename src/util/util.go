package util

import (
	rand2 "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const (
	AsciiLowercase = "abcdefghijklmnopqrstuvwxyz"
	AsciiUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AsciiLetters   = AsciiLowercase + AsciiUppercase
)

func RandomString(n int, chars string) string {
	builder := strings.Builder{}
	runes := []rune(chars)
	for i := 0; i < n; i++ {
		builder.WriteRune(runes[rand.Intn(len(runes))])
	}

	return builder.String()
}

func ValidateURL(rawurl string) (string, error) {
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return "", err
	}

	if u.Host == "" || u.Scheme == "" {
		return "", errors.New("host and scheme cannot be empty")
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return u.String(), nil
}

func RandomToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand2.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func ParseAPIForm(r *http.Request, i interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("not supported content type: %s", contentType)
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
