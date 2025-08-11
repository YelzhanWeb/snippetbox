package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/models"
)

func Home(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w)
			return
		}

		_, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(w, err)
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

		snippet, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/pages/view.tmpl",
		}
		// Parse the template files...
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", snippet)
		if err != nil {
			app.ServerError(w, err)
		}
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
