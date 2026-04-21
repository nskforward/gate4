package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr string
}

func Load() (Config, error) {
	addr := flag.String("addr", os.Getenv("GATE4_ADDR"), "address to listen")
	flag.Parse()

	var cfg Config

	if addr != nil && *addr != "" {
		cfg.Addr = *addr
	} else {
		cfg.Addr = ":4000"
	}

	return cfg, nil
}
