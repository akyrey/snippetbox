package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/akyrey/snippetbox/internal"
)

func serverError(app *internal.Application, w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.Error(
		err.Error(),
		slog.String("method", r.Method),
		slog.String("uri", r.URL.RequestURI()),
		slog.String("trace", string(debug.Stack())),
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func notFound(w http.ResponseWriter) {
	clientError(w, http.StatusNotFound)
}
