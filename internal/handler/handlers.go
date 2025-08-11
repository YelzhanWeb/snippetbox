package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/YelzhanWeb/snippetbox/internal/app"
)

func Home(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w)
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.ServerError(w, err)
		}
	}
}

func SnippetView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFound(w)
			return
		}
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func SnippetCreate(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("Create a new snippet..."))
	}

}
