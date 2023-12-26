package internal

import (
	"flag"
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
