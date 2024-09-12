package main

import (
	"flag"
	"go-url-shortener/internal/api"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	addr := flag.String("addr", ":"+port, "HTTP network address")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  api.Routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
