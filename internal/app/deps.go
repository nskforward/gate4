package app

import (
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/domain/handler/grpc/server"
	"github.com/nskforward/gate4/internal/domain/repository"
	"github.com/nskforward/gate4/internal/domain/service"
	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/di"
)

func (app *App) initDeps() {
	di.Provide[config.Config](app.container, config.Load)
	di.Provide[service.UserRepository](app.container, repository.NewUserMemoryRepo)
	di.Provide[*service.UserService](app.container, service.NewUserService)
	di.Provide[*server.UserHandler](app.container, server.NewUserHandler)
	di.Provide[*transport.GrpcServer](app.container, transport.NewGrpcServer)
}
