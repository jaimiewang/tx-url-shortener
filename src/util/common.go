package util

import (
	"encoding/json"
	"io"
	"math/rand"
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

func WriteJson(w io.Writer, i interface{}) error {
	bytes, err := json.Marshal(i)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}
