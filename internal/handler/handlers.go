package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func Home(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.NotFound(w)
			return
		}

		snippets, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(w, err)
			return
		}

		data := app.NewTemplateData(r)
		data.Snippets = snippets

		app.Render(w, http.StatusOK, "home.tmpl.html", data)
	}
}

func SnippetView(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())

		id, err := strconv.Atoi(params.ByName("id"))
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

		data := app.NewTemplateData(r)
		data.Snippet = snippet

		app.Render(w, http.StatusOK, "view.tmpl.html", data)
	}
}

func SnippetCreate(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Create a new snippet..."))
	}

}

func SnippetCreatePost(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Checking if the request method is a POST is now superfluous and can be
		// removed, because this is done automatically by httprouter.
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7
		id, err := app.Snippets.Insert(title, content, expires)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		// Update the redirect path to use the new clean URL format.
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

	}
}
