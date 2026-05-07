package config

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	Gateway struct {
		ListenAddr string
	}
	Admin struct {
		ListenAddr string
	}
	SSL struct {
		CA struct {
			CertPath string
			KeyPath  string
		}
		Server struct {
			CertPath string
			KeyPath  string
		}
	}
	FinamAddr string
	StoreDir  string
}

func Load() (Config, error) {
	gatewayListenAddr := flag.String("gateway-addr", os.Getenv("GATE4_GATEWAY_ADDR"), "gateway address to listen")
	adminListenAddr := flag.String("admin-addr", os.Getenv("GATE4_ADMIN_ADDR"), "admin address to listen")

	finamAddr := flag.String("addr-finam", os.Getenv("GATE4_FINAM_ADDR"), "finam address to connect")
	storeDir := flag.String("store-dir", os.Getenv("GATE4_STORE_DIR"), "path to store dir")

	caKey := flag.String("ssl-ca-key", os.Getenv("GATE4_SSL_CA_KEY"), "path to file of CA key")
	caCert := flag.String("ssl-ca-cert", os.Getenv("GATE4_SSL_CA_CERT"), "path to file of CA cert")

	serverKey := flag.String("ssl-server-key", os.Getenv("GATE4_SSL_SERVER_KEY"), "path to file of SERVER key")
	serverCert := flag.String("ssl-server-cert", os.Getenv("GATE4_SSL_SERVER_CERT"), "path to file of SERVER cert")

	flag.Parse()

	var cfg Config

	cfg.Gateway.ListenAddr = buildDefaultStr(gatewayListenAddr, ":4000")
	cfg.Admin.ListenAddr = buildDefaultStr(adminListenAddr, "127.0.0.1:4001")
	cfg.FinamAddr = buildDefaultStr(finamAddr, "api.finam.ru:443")
	cfg.StoreDir = buildStoreDir(storeDir)
	cfg.SSL.CA.KeyPath = buildDefaultStr(caKey, "")
	cfg.SSL.CA.CertPath = buildDefaultStr(caCert, "")
	cfg.SSL.Server.KeyPath = buildDefaultStr(serverKey, "")
	cfg.SSL.Server.CertPath = buildDefaultStr(serverCert, "")

	return cfg, nil
}

func (cfg Config) LogParams() {
	slog.Info("config initialized", "gateway-addr", cfg.Gateway.ListenAddr, "admin-addr", cfg.Admin.ListenAddr, "finam-addr", cfg.FinamAddr, "store-dir", cfg.StoreDir)
}

func buildDefaultStr(in *string, def string) string {
	if in == nil || *in == "" {
		return def
	}
	return *in
}

func buildStoreDir(in *string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if in == nil || *in == "" {
		return wd
	}
	path := *in
	if path[0] == '/' {
		return path
	}
	return filepath.Join(wd, path)
}
