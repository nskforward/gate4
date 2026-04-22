package config

import (
	"flag"
	"log/slog"
	"os"
)

type Config struct {
	Gateway struct {
		ListenAddr string
	}
	Admin struct {
		ListenAddr string
	}
	FinamAddr string
}

func Load() (Config, error) {
	gatewayListenAddr := flag.String("gateway-addr", os.Getenv("GATE4_GATEWAY_ADDR"), "gateway address to listen")
	adminListenAddr := flag.String("admin-addr", os.Getenv("GATE4_ADMIN_ADDR"), "admin address to listen")

	finamAddr := flag.String("addr-finam", os.Getenv("GATE4_FINAM_ADDR"), "finam address to connect")
	flag.Parse()

	var cfg Config

	if gatewayListenAddr != nil && *gatewayListenAddr != "" {
		cfg.Gateway.ListenAddr = *gatewayListenAddr
	} else {
		cfg.Gateway.ListenAddr = ":4000"
	}

	if adminListenAddr != nil && *adminListenAddr != "" {
		cfg.Admin.ListenAddr = *adminListenAddr
	} else {
		cfg.Admin.ListenAddr = ":4001"
	}

	if finamAddr != nil && *finamAddr != "" {
		cfg.FinamAddr = *finamAddr
	} else {
		cfg.FinamAddr = "api.finam.ru:443"
	}

	return cfg, nil
}

func (cfg Config) LogParams(logger *slog.Logger) {
	slog.Info("config initialized", "gateway-addr", cfg.Gateway.ListenAddr, "admin-addr", cfg.Admin.ListenAddr, "finam-addr", cfg.FinamAddr)
}
