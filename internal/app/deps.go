package app

import (
	"net/http"
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

	di.Provide[*http.ServeMux](app.container, httpserver.NewRouter)
	di.Provide[*httpserver.HTTPServer](app.container, httpserver.NewHTTPServer)
	di.Provide[*grpcserver.UnixServer](app.container, grpcserver.NewUnixServer)
	di.Provide[*grpcserver.TCPServer](app.container, grpcserver.NewTCPServer)

	di.Provide[*api.Server](app.container, api.NewServer)
}
