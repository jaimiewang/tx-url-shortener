package util

import (
	"errors"
	"math/rand"
	"net/url"
	"strings"
)

const (
	Digits         = "0123456789"
	SpecialChars   = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
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

func ValidateURL(u string) (string, error) {
	url_, err := url.ParseRequestURI(u)
	if err != nil {
		return "", err
	}

	if url_.Host == "" || url_.Scheme == "" {
		return "", errors.New("host and scheme cannot be empty")
	}

	if !strings.HasSuffix(url_.Path, "/") {
		url_.Path += "/"
	}

	return url_.String(), nil
}
