package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func pong(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("pong"))
}

// Routes creates the application's routing table
func Routes() http.Handler {
	router := httprouter.New()
	// TODO add a post endpoint that handles url shortener
	// TODO add an endpoint that will handle redirect url
	router.HandlerFunc(http.MethodGet, "/api/ping", pong)
	standard := alice.New()

	return standard.Then(router)
}
