package config

import (
	"flag"
	"log/slog"
	"os"
)

type Config struct {
	GRPC struct {
		Addr string
	}

	HTTP struct {
		Addr string
	}

	SSL struct {
		CA struct {
			Cert string
			Key  string
		}
		Server struct {
			Cert string
			Key  string
		}
	}

	Finam struct {
		Addr string
	}

	FileStorageDir string
}

func Load(args []string) Config {
	flags := flag.NewFlagSet("app", flag.ExitOnError)

	grpcAddr := flags.String("grpc-addr", os.Getenv("GATE4_GRPC_ADDR"), "tcp address of grpc server to listen")
	httpAddr := flags.String("http-addr", os.Getenv("GATE4_HTTP_ADDR"), "tcp address of http server to listen")

	finamAddr := flags.String("addr-finam", os.Getenv("GATE4_FINAM_ADDR"), "finam address to connect")
	storeDir := flags.String("store-dir", os.Getenv("GATE4_STORE_DIR"), "path to store dir")

	caKey := flags.String("ssl-ca-key", os.Getenv("GATE4_SSL_CA_KEY"), "path to file of CA key")
	caCert := flags.String("ssl-ca-cert", os.Getenv("GATE4_SSL_CA_CERT"), "path to file of CA cert")

	serverKey := flags.String("ssl-key", os.Getenv("GATE4_SSL_KEY"), "server ssl key path")
	serverCert := flags.String("ssl-cert", os.Getenv("GATE4_SSL_CERT"), "server ssl cert path")

	err := flags.Parse(args)
	if err != nil {
		slog.Error("config cannot parse flags", "reason", err.Error())
		os.Exit(1)
	}

	var cfg Config

	cfg.HTTP.Addr = useDefault(*httpAddr, ":443")
	cfg.GRPC.Addr = useDefault(*grpcAddr, ":4443")

	cfg.Finam.Addr = useDefault(*finamAddr, "api.finam.ru:443")
	cfg.FileStorageDir = useDefault(*storeDir, "data/storage")
	cfg.SSL.CA.Key = useDefault(*caKey, "data/ssl/ca.key")
	cfg.SSL.CA.Cert = useDefault(*caCert, "data/ssl/ca.crt")
	cfg.SSL.Server.Key = useDefault(*serverKey, "data/ssl/server.key")
	cfg.SSL.Server.Cert = useDefault(*serverCert, "data/ssl/server.crt")

	return cfg
}

func useDefault(in string, def string) string {
	if in == "" {
		return def
	}
	return in
}
