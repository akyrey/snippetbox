package internal

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Application struct {
	Logger        *slog.Logger
	DB            *sql.DB
	TemplateCache map[string]*template.Template
}

func (app *Application) Routes(config Config) http.Handler {
	router := httprouter.New()

	// router.MethodNotAllowed could also be assigned some custom handler
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// logRequest ↔ secureHeaders ↔ servemux ↔ application handler
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
