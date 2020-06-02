package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
)

func RandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
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
