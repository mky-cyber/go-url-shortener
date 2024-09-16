package handler

import (
	"encoding/json"
	h "go-url-shortener/internal/api/http"
	"net/http"
	"net/url"

	urlverifier "github.com/davidmytton/url-verifier"
)

func SendErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(h.URLResponse{
		Message: errorMessage,
	})
}

func IsValidURL(URL string) bool {
	u, err := url.ParseRequestURI(URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}
	return true
}

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
