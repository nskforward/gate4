package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/infra"
	"github.com/nskforward/gate4/internal/transport"
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
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	infra.InitLogger()

	grpcServer, err := di.Resolve[*transport.GRPCTransport](app.container)
	if err != nil {
		return err
	}

	return grpcServer.Start(ctx)
}
