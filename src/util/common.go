package util

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
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

func Serialize(value interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Deserialize(b []byte, value interface{}) error {
	buf := bytes.NewBuffer(b)

	if err := gob.NewDecoder(buf).Decode(value); err != nil {
		return err
	}

	return nil
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

func WriteJson(w io.Writer, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
