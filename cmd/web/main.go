package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/akyrey/snippetbox/internal"
)

func main() {
	config := internal.Config{}
	config.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err := openDb()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &internal.Application{
		Logger: logger,
		DB:     db,
	}

	logger.Info("starting server", slog.String("addr", config.Addr))

	err = http.ListenAndServe(config.Addr, app.Routes(config))
	logger.Error(err.Error())
	os.Exit(1)
}

func openDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
