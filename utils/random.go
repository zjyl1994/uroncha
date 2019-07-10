package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomString(length int) string {
	b := make([]rune, length)
	lenLetter := len(letterRunes)
	for i := range b {
		b[i] = letterRunes[rand.Intn(lenLetter)]
	}
	return string(b)
}
