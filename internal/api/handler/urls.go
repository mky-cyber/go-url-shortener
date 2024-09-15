package handler

import (
	"encoding/json"
	"fmt"
	h "go-url-shortener/internal/api/http"
	"go-url-shortener/internal/models"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Max length for URLs
const maxURLLength = 2048

// openShortenedURL retrives the original URL using the shortened URL provided,
// then redirect the user to the original URL
func OpenShortenedURL(sd models.ShortenerDataInterface) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Retrieve the shortened URL from the path parameter
		shortenedURLKey := ps.ByName("shortenedURLKey")
		if !isValidURLKey(shortenedURLKey) {
			SendErrorResponse(w, "Shortened URL is invalid", http.StatusNotFound)
			return
		}

		// Check if the shortened URL exists in the db
		data, err := sd.Get(shortenedURLKey)
		if err != nil {
			http.Error(w, "Shortened URL not found", http.StatusNotFound)
			return
		}

		// Increase the clicks for monitor purpose
		err = sd.IncreaseClicks(shortenedURLKey)
		if err != nil {
			http.Error(w, "Unable to update the clicks", http.StatusInternalServerError)
			return
		}

		// Redirect to the original URL
		http.Redirect(w, r, data.OriginalURL, http.StatusSeeOther)
	}
}

func ShortenedURL(sd models.ShortenerDataInterface) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req h.URLRequest
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

		// setup the url dynamically based on the request host

		// TODO add a rate limit
		// TODO add a blacklist for banned urls
		// TODO check speical characters

		// Generate a unique key and save it in db
		shortenedURLKey := GenerateShortURLKey()

		// Handle concurrent processes
		var mu sync.Mutex
		mu.Lock()
		rowsAffected, err := sd.Insert(req.URL, shortenedURLKey, 0)
		mu.Unlock()

		if rowsAffected != 1 || err != nil {
			SendErrorResponse(w, "Failed to create shortened URL", http.StatusInternalServerError)
			return
		}
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		shortenedURL := url.URL{
			Scheme: scheme,
			Host:   r.Host,
			Path:   fmt.Sprintf("/s/%s", shortenedURLKey),
		}

		response := h.URLResponse{
			Result:  shortenedURL.String(),
			Message: "URL successfully shortened",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
