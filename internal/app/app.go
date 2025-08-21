package app

import (
	"html/template"
	"log"

	"github.com/YelzhanWeb/snippetbox/internal/models"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Snippets      *models.SnippetModel
	TemplateCache map[string]*template.Template
}
