package server

import (
	"net/http"

	ap "github.com/YelzhanWeb/snippetbox/internal/app"
	"github.com/YelzhanWeb/snippetbox/internal/handler"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func Routes(app *ap.Application) http.Handler {
	// mux := http.NewServeMux()
	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// mux.HandleFunc("/", handler.Home(app))
	// mux.HandleFunc("/snippet/view", handler.SnippetView(app))
	// mux.HandleFunc("/snippet/create", handler.SnippetCreate(app))

	// standard := alice.New(app.RecoverPanic, app.LogRequest, ap.SecureHeaders)
	// return standard.Then(mux)

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	dynamic := alice.New(app.SessionManager.LoadAndSave)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(handler.Home(app)))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(handler.SnippetView(app)))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(handler.SnippetCreate(app)))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(handler.SnippetCreatePost(app)))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(handler.UserSignup(app)))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(handler.UserSignupPost(app)))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(handler.UserLogin(app)))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(handler.UserLoginPost(app)))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(handler.UserLogoutPost(app)))

	standard := alice.New(app.RecoverPanic, app.LogRequest, ap.SecureHeaders)

	return standard.Then(router)
}
