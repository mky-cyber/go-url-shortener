package api

import (
	"go-url-shortener/internal/api/handler"
	"go-url-shortener/internal/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type App struct {
	urls models.ShortenerDataInterface
}

func NewApp(dataInterface models.ShortenerDataInterface) *App {
	return &App{
		urls: dataInterface,
	}
}

// pong just writes pong to response to test if the server is working
func pong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Write([]byte("pong"))
}

// Routes creates the application's routing table
func (app *App) Routes() http.Handler {
	router := httprouter.New()
	router.GET("/ping", pong)
	router.GET("/s/:shortenedURLKey", handler.OpenShortenedURL(app.urls))
	router.POST("/shorten", handler.ShortenedURL(app.urls))
	standard := alice.New()

	return standard.Then(router)
}
