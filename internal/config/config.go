package config

import (
	"flag"
	"log/slog"
	"os"
)

type Config struct {
	ListenAddr string
	FinamAddr  string
}

func Load() (Config, error) {
	listenAddr := flag.String("addr-listen", os.Getenv("GATE4_ADDR"), "address to listen")
	finamAddr := flag.String("addr-finam", os.Getenv("FINAM_ADDR"), "finam address")
	flag.Parse()

	var cfg Config

	if listenAddr != nil && *listenAddr != "" {
		cfg.ListenAddr = *listenAddr
	} else {
		cfg.ListenAddr = ":4000"
	}

	if finamAddr != nil && *finamAddr != "" {
		cfg.FinamAddr = *finamAddr
	} else {
		cfg.ListenAddr = "api.finam.ru:443"
	}

	return cfg, nil
}

func (cfg Config) LogParams(logger *slog.Logger) {
	slog.Info("config initialized", "addr-listen", cfg.ListenAddr, "finam-addr", cfg.FinamAddr)
}
