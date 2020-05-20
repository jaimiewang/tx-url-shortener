package util

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bytes)
}
