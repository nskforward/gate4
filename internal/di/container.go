package di

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"log/slog"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/console"
	"github.com/nskforward/gate4/pkg/ssl"
	"google.golang.org/grpc"
)

type Container struct {
	config        *config.Config
	logger        *slog.Logger
	userStore     users.Store
	keychainStore *keychain.Store
	clientPool    *brokers.ClientPool
	tlsConfig     *tls.Config
	caCert        *x509.Certificate
	caKey         crypto.PrivateKey
	serverCert    *x509.Certificate
	serverKey     crypto.PrivateKey
	apiServer     *api.Server
	unixServer    *grpc.Server
	tcpServer     *grpc.Server
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Config() *config.Config {
	if c.config == nil {
		c.config = config.Load()
	}
	return c.config
}

func (c *Container) Logger() *slog.Logger {
	if c.logger == nil {
		c.logger = initLogger()
	}
	return c.logger
}

func (c *Container) TLSConfig() *tls.Config {
	if c.tlsConfig == nil {
		c.tlsConfig = initTLSConfig(c.ServerCert(), c.ServerKey(), c.CACert())
	}
	return c.tlsConfig
}

func (c *Container) CAKey() crypto.PrivateKey {
	if c.caKey == nil {
		key, err := ssl.LoadPrivateKey(c.Config().CA.Key)
		if err != nil {
			console.LogFatal("cannot load CA private key", err)
		}
		c.caKey = key
	}
	return c.caKey
}

func (c *Container) CACert() *x509.Certificate {
	if c.caCert == nil {
		cert, err := ssl.LoadCertificate(c.Config().CA.Cert)
		if err != nil {
			console.LogFatal("cannot load CA certificate", err)
		}
		c.caCert = cert
	}
	return c.caCert
}

func (c *Container) ServerKey() crypto.PrivateKey {
	if c.serverKey == nil {
		key, err := ssl.LoadPrivateKey(c.Config().Server.SSL.Key)
		if err != nil {
			console.LogFatal("cannot load server private key", err)
		}
		c.serverKey = key
	}
	return c.serverKey
}

func (c *Container) ServerCert() *x509.Certificate {
	if c.serverCert == nil {
		cert, err := ssl.LoadCertificate(c.Config().Server.SSL.Cert)
		if err != nil {
			console.LogFatal("cannot load server certificate", err)
		}
		c.serverCert = cert
	}
	return c.serverCert
}

func (c *Container) UserSore() users.Store {
	if c.userStore == nil {
		storage, err := users.NewFileStorage()
		if err != nil {
			console.LogFatal("cannot init UserStore", err)
		}
		c.userStore = storage
	}
	return c.userStore
}

func (c *Container) KeychainStore() *keychain.Store {
	if c.keychainStore == nil {
		c.keychainStore = keychain.NewStore(c.CAKey(), c.CACert())
	}
	return c.keychainStore
}

func (c *Container) ClientPool() *brokers.ClientPool {
	if c.clientPool == nil {
		c.clientPool = brokers.NewClientPool()
	}
	return c.clientPool
}

func (c *Container) APIServer() *api.Server {
	if c.apiServer == nil {
		c.apiServer = api.NewServer(c.UserSore(), c.KeychainStore(), c.ClientPool())
	}
	return c.apiServer
}
