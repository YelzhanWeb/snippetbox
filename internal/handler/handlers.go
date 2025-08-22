package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/models"
	"github.com/YelzhanWeb/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

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
		data := app.NewTemplateData(r)
		data.Form = snippetCreateForm{
			Expires: 365,
		}

		app.Render(w, http.StatusOK, "create.tmpl.html", data)
	}

}

func SnippetCreatePost(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var form snippetCreateForm

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

		if !form.Valid() {
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
			return
		}
		id, err := app.Snippets.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}
