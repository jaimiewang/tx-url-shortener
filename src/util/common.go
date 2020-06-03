package util

import (
	"errors"
	"math/rand"
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
