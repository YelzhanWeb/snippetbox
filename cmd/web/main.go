package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/handler"
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

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", handler.Home(app))
	mux.HandleFunc("/snippet/view", handler.SnippetView(app))
	mux.HandleFunc("/snippet/create", handler.SnippetCreate(app))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
