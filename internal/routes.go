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
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.HomeHandler())
	mux.HandleFunc("/snippet/view", app.SnippetViewHandler())
	mux.HandleFunc("/snippet/create", app.SnipperCreateHandler())

  // logRequest ↔ secureHeaders ↔ servemux ↔ application handler
  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}
