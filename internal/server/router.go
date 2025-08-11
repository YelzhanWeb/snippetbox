package server

import (
	"net/http"

	"github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/handler"
)

// The routes() method returns a servemux containing our application routes.
func Routes(app *app.Application) *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", handler.Home(app))
	mux.HandleFunc("/snippet/view", handler.SnippetView(app))
	mux.HandleFunc("/snippet/create", handler.SnippetCreate(app))
	return mux
}
