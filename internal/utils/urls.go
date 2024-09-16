package utils

import (
	"fmt"
	"regexp"
	"time"

	"math/rand"
)

const (
	Charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	URLKeyLength = 16
)

func IsValidURLKey(key string) bool {
	if len(key) != URLKeyLength {
		return false
	}
	// use the same charset that generates the key to check if the received key is valid
	pattern := fmt.Sprintf("^[%s]+$", regexp.QuoteMeta(Charset))
	match, err := regexp.MatchString(pattern, key)
	if err != nil {
		return false
	}

	return match
}

func randomGenerator() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func generateRandomKey(length int) string {
	r := randomGenerator()
	key := make([]byte, length)
	for i := range key {
		key[i] = Charset[r.Intn(len(Charset))]
	}
	return string(key)
}

func GenerateShortURLKey() string {
	key := generateRandomKey(URLKeyLength)
	return key
}
