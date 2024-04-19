package internal

import (
	"database/sql"
	"flag"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Addr      string
	StaticDir string
}

func (c *Config) Parse() {
	flag.StringVar(&c.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&c.StaticDir, "static-dir", "./ui/static/", "Path to static assets")

	flag.Parse()
}

func (c *Config) OpenDb() (*sql.DB, error) {
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
