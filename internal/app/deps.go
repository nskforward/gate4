package app

import (
	"fmt"
	"os"

	"github.com/nskforward/gate4/internal/api"
	grpcserver "github.com/nskforward/gate4/internal/api/grpc/server"
	"github.com/nskforward/gate4/internal/api/grpc/server/handler"
	httpserver "github.com/nskforward/gate4/internal/api/http/server"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/domain/repository"
	"github.com/nskforward/gate4/internal/domain/repository/fs"
	"github.com/nskforward/gate4/internal/domain/service"
	"github.com/nskforward/gate4/pkg/di"
)

func (app *App) initDeps() {
	di.Provide[config.Config](app.container, func() config.Config {
		return config.Load(os.Args)
	})

	di.Provide[repository.UserRepository](app.container, fs.NewUserRepo)
	di.Provide[repository.TokenRepository](app.container, fs.NewTokenRepo)

	di.Provide[*service.UserService](app.container, service.NewUserService)
	di.Provide[*service.TokenService](app.container, service.NewTokenService)

	di.Provide[*handler.UserHandler](app.container, handler.NewUserHandler)
	di.Provide[*handler.TokenHandler](app.container, handler.NewTokenHandler)

	di.Provide[*httpserver.HTTPServer](app.container, httpserver.NewHTTPServer)

	di.Provide[*grpcserver.UnixServer](app.container, grpcserver.NewUnixServer)
	di.Provide[*grpcserver.TCPServer](app.container, func() (*grpcserver.TCPServer, error) {
		cfg, err := di.Resolve[config.Config](app.container)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve config: %w", err)
		}
		tlsConfig, err := grpcserver.NewTLSConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("cannot get tls config: %w", err)
		}
		return grpcserver.NewTCPServer(cfg.TCPAddr, tlsConfig), nil
	})

	di.Provide[*api.Server](app.container, func() (*api.Server, error) {

		unixServer, err := di.Resolve[*grpcserver.UnixServer](app.container)
		if err != nil {
			return nil, err
		}

		tcpServer, err := di.Resolve[*grpcserver.TCPServer](app.container)
		if err != nil {
			return nil, err
		}

		httpServer, err := di.Resolve[*httpserver.HTTPServer](app.container)
		if err != nil {
			return nil, err
		}

		return api.NewServer(unixServer, tcpServer, httpServer), nil
	})
}
