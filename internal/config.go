package internal

import (
	"database/sql"
	"flag"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Addr string
}

func (c *Config) Parse() {
	flag.StringVar(&c.Addr, "addr", ":4000", "HTTP network address")

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
