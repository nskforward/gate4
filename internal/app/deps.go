package app

import (
	"github.com/nskforward/gate4/internal/config"
	handler "github.com/nskforward/gate4/internal/domain/handler/user"
	repository "github.com/nskforward/gate4/internal/domain/repository/user"
	usecases "github.com/nskforward/gate4/internal/domain/usecases/user"
	"github.com/nskforward/gate4/pkg/di"
)

func (app *App) initDeps() {
	di.Provide[config.Config](app.container, config.Load)
	di.Provide[usecases.UserRepository](app.container, repository.NewMemoryRepo)
	di.Provide[*usecases.UserUsecases](app.container, usecases.NewUserUsecases)
	di.Provide[*handler.UserHandler](app.container, handler.NewUserHandler)
}
