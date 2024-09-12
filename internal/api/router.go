package api

import (
	"go-url-shortener/internal/api/handler"
	"go-url-shortener/internal/model"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// pong just writes pong to response to test if the server is working
func pong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("pong"))
}

// Routes creates the application's routing table
func Routes() http.Handler {
	// Initialize the URLShortener with some URLs
	shortener := &model.URLShortener{
		URLs: map[string]string{
			"abc123": "https://github.com/",
			"def456": "https://google.com/",
		},
	}

	router := httprouter.New()
	// TODO add a post endpoint that handles url shortener
	router.GET("/ping", pong)
	router.GET("/s/:shortenedURL", handler.OpenShortenedURL(shortener))
	standard := alice.New()

	return standard.Then(router)
}
