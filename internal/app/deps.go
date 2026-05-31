package app

import (
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/domain/handler"
	"github.com/nskforward/gate4/internal/domain/repository"
	"github.com/nskforward/gate4/internal/domain/service"
	"github.com/nskforward/gate4/internal/transport"
	"github.com/nskforward/gate4/pkg/di"
)

func (app *App) initDeps() {
	di.Provide[config.Config](app.container, config.Load)
	di.Provide[service.UserRepository](app.container, repository.NewMemoryUserRepo)
	di.Provide[*service.UserService](app.container, service.NewUserService)
	di.Provide[*handler.GRPCUserHandler](app.container, handler.NewGRPCUserHandler)
	di.Provide[*transport.GRPCTransport](app.container, transport.NewGRPCTransport)
}
