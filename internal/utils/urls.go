package utils

import (
	"encoding/json"
	"fmt"
	h "go-url-shortener/internal/api/http"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"math/rand"

	urlverifier "github.com/davidmytton/url-verifier"
)

const (
	Charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	URLKeyLength = 16
)

// Verify if the URL supplied is a genuine and workable URL
func CheckGenuineURL(originalURL string) bool {
	verifier := urlverifier.NewVerifier()
	verifier.EnableHTTPCheck()
	result, err := verifier.Verify(originalURL)

	if err != nil || !result.HTTP.IsSuccess {
		return false
	}
	return true
}

func IsValidURL(URL string) bool {
	u, err := url.ParseRequestURI(URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}
	return true
}

func SendErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(h.URLResponse{
		Message: errorMessage,
	})
}

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
