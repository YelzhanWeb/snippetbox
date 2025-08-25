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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
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
		form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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

func UserSignup(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = userSignupForm{}
		app.Render(w, http.StatusOK, "signup.tmpl.html", data)
	}
}

func UserSignupPost(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form userSignupForm

		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
		form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

		if !form.Valid() {
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
			return
		}

		err = app.Users.Insert(form.Name, form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateEmail) {
				form.AddFieldError("email", "Email address is already in use")
				data := app.NewTemplateData(r)
				data.Form = form
				app.Render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		app.SessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	}
}

func UserLogin(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := app.NewTemplateData(r)
		data.Form = userLoginForm{}
		app.Render(w, http.StatusOK, "login.tmpl.html", data)
	}
}

func UserLoginPost(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form userLoginForm
		err := app.DecodePostForm(r, &form)
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
		form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
		form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

		if !form.Valid() {
			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		}

		id, err := app.Users.Authenticate(form.Email, form.Password)
		if err != nil {
			if errors.Is(err, models.ErrInvalidCredentials) {
				form.AddNonFieldError("Email or password is incorrect")
				data := app.NewTemplateData(r)
				data.Form = form
				app.Render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		err = app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.SessionManager.Put(r.Context(), "authenticatedUserID", id)

		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	}
}

func UserLogoutPost(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := app.SessionManager.RenewToken(r.Context())
		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.SessionManager.Remove(r.Context(), "authenticatedID")

		app.SessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
