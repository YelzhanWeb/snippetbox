package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/server"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &app.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  server.Routes(app),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
