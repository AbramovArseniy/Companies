// package cfg creates configuration for servers
package cfg

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

// Config describes server config
type Config struct {
	Address   string `env:"RUN_ADDRESS"`
	DBAddress string `env:"DATABASE_URI"`
}

// New creates Config from environment and flags
func New() *Config {
	var cfg Config

	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "set server listening address")
	flag.StringVar(&cfg.DBAddress, "d", "", "set the DB address")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Printf("env parse failed :%s", err)
	}

	return &cfg
}
