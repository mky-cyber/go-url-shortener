package handler

import (
	"fmt"
	"go-url-shortener/internal/model"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
