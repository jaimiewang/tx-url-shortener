package util

import (
	"errors"
	"math/rand"
	"net/url"
	"strings"
)

var (
	AsciiLowercase = []rune("abcdefghijklmnopqrstuvwxyz")
	AsciiUppercase = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	AsciiLetters   = append(AsciiLowercase, AsciiUppercase...)
)

func RandomString(n int, runes []rune) string {
	builder := strings.Builder{}
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
