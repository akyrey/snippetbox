package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/akyrey/snippetbox/internal"
)

type Application struct {
	logger *slog.Logger
}

func main() {
	config := internal.Config{}
	config.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	app := &Application{
		logger: logger,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	logger.Info("starting server", slog.String("addr", config.Addr))

	err := http.ListenAndServe(config.Addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
