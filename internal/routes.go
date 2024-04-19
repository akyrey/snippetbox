package internal

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
)

type Application struct {
	Logger        *slog.Logger
	DB            *sql.DB
	TemplateCache map[string]*template.Template
}

func (app *Application) Routes(config Config) http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("GET /{$}", app.home)
	router.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	router.HandleFunc("GET /snippet/create", app.snippetCreate)
	router.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// logRequest ↔ secureHeaders ↔ servemux ↔ application handler
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
