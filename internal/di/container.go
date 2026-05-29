package di

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"log/slog"
	"os"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/domain/users"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/storage"
	"github.com/nskforward/gate4/pkg/ssl"
	"google.golang.org/grpc"
)

type Container struct {
	config        *config.Config
	logger        *slog.Logger
	objectStorage storage.ObjectStorage
	userStorage   *users.UserStorage
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
			slog.Error("cannot load CA private key", "reason", err)
			os.Exit(1)
		}
		c.caKey = key
	}
	return c.caKey
}

func (c *Container) CACert() *x509.Certificate {
	if c.caCert == nil {
		cert, err := ssl.LoadCertificate(c.Config().CA.Cert)
		if err != nil {
			slog.Error("cannot load CA certificate", "reason", err)
			os.Exit(1)
		}
		c.caCert = cert
	}
	return c.caCert
}

func (c *Container) ServerKey() crypto.PrivateKey {
	if c.serverKey == nil {
		key, err := ssl.LoadPrivateKey(c.Config().Server.SSL.Key)
		if err != nil {
			slog.Error("cannot load server private key", "reason", err)
			os.Exit(1)
		}
		c.serverKey = key
	}
	return c.serverKey
}

func (c *Container) ServerCert() *x509.Certificate {
	if c.serverCert == nil {
		cert, err := ssl.LoadCertificate(c.Config().Server.SSL.Cert)
		if err != nil {
			slog.Error("cannot load server certificate", "reason", err)
			os.Exit(1)
		}
		c.serverCert = cert
	}
	return c.serverCert
}

func (c *Container) ObjectStorage() storage.ObjectStorage {
	if c.objectStorage == nil {
		c.objectStorage = storage.NewFileObjectStorage(c.Config().FileStorageDir)
	}
	return c.objectStorage
}

func (c *Container) UserStorage() *users.UserStorage {
	if c.userStorage == nil {
		c.userStorage = users.NewUserStorage(c.ObjectStorage())
	}
	return c.userStorage
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
		c.apiServer = api.NewServer(c.UserStorage(), c.KeychainStore(), c.ClientPool())
	}
	return c.apiServer
}
