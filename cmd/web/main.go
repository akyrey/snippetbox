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

	logger.Info("starting server", slog.String("addr", config.Addr))

	err := http.ListenAndServe(config.Addr, app.Routes(config))
	logger.Error(err.Error())
	os.Exit(1)
}
