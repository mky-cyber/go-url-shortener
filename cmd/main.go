package main

import (
	"database/sql"
	"flag"
	"go-url-shortener/internal/api"
	"go-url-shortener/internal/models"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := flag.String("addr", ":"+port, "HTTP network address")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbFile := "../db/update/database.db"
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	URLShortener := &models.ShortenerDBModel{DB: db}
	app := api.NewApp(URLShortener)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
