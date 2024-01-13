package internal

import (
	"database/sql"
	"log/slog"
	"net/http"
)

type Application struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func (app *Application) Routes(config Config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.HomeHandler())
	mux.HandleFunc("/snippet/view", app.SnippetViewHandler())
	mux.HandleFunc("/snippet/create", app.SnipperCreateHandler())

	return mux
}
