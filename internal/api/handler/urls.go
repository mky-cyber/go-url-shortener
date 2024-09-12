package handler

import (
	"encoding/json"
	"fmt"
	h "go-url-shortener/internal/api/http"
	"go-url-shortener/internal/model"
	"net/http"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Max length for URLs
const maxURLLength = 2048

// openShortenedURL retrives the original URL using the shortened URL provided,
// then redirect the user to the original URL
func OpenShortenedURL(shortener *model.URLShortener) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Retrieve the shortened URL from the path parameter
		shortenedURL := ps.ByName("shortenedURL")
		if shortenedURL == "" {
			http.Error(w, "Shortened URL is empty", http.StatusBadRequest)
			return
		}

		fmt.Printf("Attempting to retrieve %s.\n", shortenedURL)

		// Check if the shortened URL exists in the map
		originalURL, found := shortener.URLs[shortenedURL]
		if !found {
			http.Error(w, "Shortened URL not found", http.StatusNotFound)
			return
		}

		fmt.Printf("Redirecting to %s.\n", originalURL)

		// Redirect to the original URL
		http.Redirect(w, r, originalURL, http.StatusSeeOther)
	}
}

func ShortenedURL(shortener *model.URLShortener) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req h.URLRequest
		fmt.Printf("Attempting to shorten URL\n")
		// Decode the JSON from request body
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			SendErrorResponse(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		// Check if the URL is empty or missing
		if strings.TrimSpace(req.URL) == "" {
			SendErrorResponse(w, "Missing url in the request payload", http.StatusBadRequest)
			return
		}

		// Check if the URL is valid
		if !IsValidURL(req.URL) {
			SendErrorResponse(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		// Check if the URL is too long
		if len(req.URL) > maxURLLength {
			SendErrorResponse(w, fmt.Sprintf("URL exceeds the maximum length of %d characters", maxURLLength), http.StatusBadRequest)
			return
		}

		// Check if the URL is genuine
		if !CheckGenuineURL(req.URL) {
			SendErrorResponse(w, "The URL was not reachable", http.StatusBadRequest)
			return
		}

		// Check if the URL already exists
		for shortKey, original := range shortener.URLs {
			fmt.Printf("Checking if URL exists \n")
			if original == req.URL {
				response := h.URLResponse{
					Result:  fmt.Sprintf("http://localhost:8080/s/%s", shortKey),
					Message: "URL is already shortened",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		// TODO add a rate limit
		// TODO add a blacklist for banned urls
		// TODO check speical characters

		fmt.Printf("Generating shortened URL for %s \n", req.URL)
		// Generate a unique key and save it in the map
		shortenedURLKey := GenerateShortURLKey(6)

		// Handle corrency
		var mu sync.Mutex
		mu.Lock()
		shortener.URLs[shortenedURLKey] = req.URL
		mu.Unlock()

		response := h.URLResponse{
			Result:  fmt.Sprintf("http://localhost:8080/s/%s", shortenedURLKey),
			Message: "URL successfully shortened",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
