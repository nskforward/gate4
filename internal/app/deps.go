package app

import (
	"os"

	"github.com/nskforward/gate4/internal/api/grpc/server"
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

	di.Provide[*server.UserHandler](app.container, server.NewUserHandler)
	di.Provide[*server.TokenHandler](app.container, server.NewTokenHandler)

	di.Provide[*server.Server](app.container, server.NewServer)
}
