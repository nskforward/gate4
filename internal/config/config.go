package config

import (
	"flag"
	"os"
)

type Config struct {
	TCPAddr string
	SSL     struct {
		Cert string
		Key  string
	}

	CA struct {
		Cert string
		Key  string
	}

	Finam struct {
		APIAddr string
	}

	FileStorageDir string
}

func Load() Config {
	tcpAddr := flag.String("tcp-addr", os.Getenv("GATE4_TCP_ADDR"), "server tcp address to listen")
	finamAddr := flag.String("addr-finam", os.Getenv("GATE4_FINAM_ADDR"), "finam address to connect")
	storeDir := flag.String("store-dir", os.Getenv("GATE4_STORE_DIR"), "path to store dir")

	caKey := flag.String("ssl-ca-key", os.Getenv("GATE4_SSL_CA_KEY"), "path to file of CA key")
	caCert := flag.String("ssl-ca-cert", os.Getenv("GATE4_SSL_CA_CERT"), "path to file of CA cert")

	serverKey := flag.String("ssl-key", os.Getenv("GATE4_SSL_KEY"), "server ssl key path")
	serverCert := flag.String("ssl-cert", os.Getenv("GATE4_SSL_CERT"), "server ssl cert path")

	flag.Parse()

	var cfg Config

	cfg.TCPAddr = useDefault(*tcpAddr, ":443")
	cfg.Finam.APIAddr = useDefault(*finamAddr, "api.finam.ru:443")
	cfg.FileStorageDir = useDefault(*storeDir, "data/storage")
	cfg.CA.Key = useDefault(*caKey, "data/ssl/ca.key")
	cfg.CA.Cert = useDefault(*caCert, "data/ssl/ca.cert")
	cfg.SSL.Key = useDefault(*serverKey, "data/ssl/server.key")
	cfg.SSL.Cert = useDefault(*serverCert, "data/ssl/server.cert")

	return cfg
}

func useDefault(in string, def string) string {
	if in == "" {
		return def
	}
	return in
}
