package server

import (
	"net/http"

	ap "github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/handler"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func Routes(app *ap.Application) http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", handler.Home(app))
	mux.HandleFunc("/snippet/view", handler.SnippetView(app))
	mux.HandleFunc("/snippet/create", handler.SnippetCreate(app))

	standard := alice.New(app.RecoverPanic, app.LogRequest, ap.SecureHeaders)
	return standard.Then(mux)
}
