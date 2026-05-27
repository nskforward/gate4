package di

import (
	"crypto/tls"
	"log/slog"

	"github.com/nskforward/gate4/internal/api"
	"github.com/nskforward/gate4/internal/brokers"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/keychain"
	"github.com/nskforward/gate4/internal/users"
	"github.com/nskforward/gate4/pkg/console"
	"google.golang.org/grpc"
)

type Container struct {
	config        *config.Config
	logger        *slog.Logger
	userStore     users.Store
	keychainStore *keychain.Store
	clientPool    *brokers.ClientPool
	tlsConfig     *tls.Config
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
		cfg := c.Config()
		c.tlsConfig = initTLSConfig(cfg)
	}
	return c.tlsConfig
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
		cfg := c.Config()
		store, err := keychain.NewStore(cfg.CA.Key, cfg.CA.Cert)
		if err != nil {
			console.LogFatal("cannot init KeychainStore", err)
		}
		c.keychainStore = store
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

func (c *Container) UnixServer() *grpc.Server {
	if c.unixServer == nil {
		c.unixServer = NewUnixServer(c.APIServer())
	}
	return c.unixServer
}

func (c *Container) TCPServer() *grpc.Server {
	if c.tcpServer == nil {
		c.tcpServer = NewTCPServer(c.APIServer(), c.TLSConfig())
	}
	return c.tcpServer
}
