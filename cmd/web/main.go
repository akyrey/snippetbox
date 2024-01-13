package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/akyrey/snippetbox/internal"
)

func main() {
	config := internal.Config{}
	config.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	app := &internal.Application{
		Logger: logger,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", HomeHandler(app))
	mux.HandleFunc("/snippet/view", SnippetViewHandler(app))
	mux.HandleFunc("/snippet/create", SnipperCreateHandler(app))

	logger.Info("starting server", slog.String("addr", config.Addr))

	err := http.ListenAndServe(config.Addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
