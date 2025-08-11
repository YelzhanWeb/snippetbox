package app

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/YelzhanWeb/snippetbox/internal/models"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Snippets *models.SnippetModel
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
