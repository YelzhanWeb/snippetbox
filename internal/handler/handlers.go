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
			http.NotFound(w, r)
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/home.tmpl",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.ErrorLog.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}

func SnippetView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	}
}

func SnippetCreate(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("Create a new snippet..."))
	}

}
