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
		SSL        struct {
			Cert string
			Key  string
		}
	}
	Admin struct {
		ListenAddr string
		SSL        struct {
			CACert string
			CAKey  string
		}
	}
	Finam struct {
		APIAddr string
	}
	FileStorageDir string
}

func Load() Config {
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
	cfg.Finam.APIAddr = buildDefaultStr(finamAddr, "api.finam.ru:443")
	cfg.FileStorageDir = buildStoreDir(storeDir)
	cfg.Admin.SSL.CAKey = buildDefaultStr(caKey, "")
	cfg.Admin.SSL.CACert = buildDefaultStr(caCert, "")
	cfg.Gateway.SSL.Key = buildDefaultStr(serverKey, "")
	cfg.Gateway.SSL.Cert = buildDefaultStr(serverCert, "")

	return cfg
}

func (cfg Config) Log() {
	slog.Info("config initialized",
		slog.Group("gateway",
			"addr", cfg.Gateway.ListenAddr,
			slog.Group("ssl",
				"key", cfg.Gateway.SSL.Key,
				"cert", cfg.Gateway.SSL.Cert,
			),
		),
		slog.Group("admin",
			"addr", cfg.Admin.ListenAddr,
			"ca-key", cfg.Admin.SSL.CAKey,
			"ca-cert", cfg.Admin.SSL.CACert,
		),
		"finam-api-addr", cfg.Finam.APIAddr,
		"file-storage-dir", cfg.FileStorageDir,
	)
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
