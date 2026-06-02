package app

import (
	"context"

	"github.com/nskforward/gate4/internal/api/grpc/server"
	"github.com/nskforward/gate4/internal/infra"
	"github.com/nskforward/gate4/pkg/di"
)

type App struct {
	container *di.Container
}

func NewApp() *App {
	c := di.NewContainer()
	a := &App{
		container: c,
	}
	a.initDeps()
	return a
}

func (app *App) Start(ctx context.Context) error {
	infra.InitLogger()

	grpcServer, err := di.Resolve[*server.Server](app.container)
	if err != nil {
		return err
	}

	return grpcServer.Start(ctx)
}
