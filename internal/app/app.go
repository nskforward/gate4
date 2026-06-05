package app

import (
	"context"

	"github.com/nskforward/gate4/internal/api"
	grpcserver "github.com/nskforward/gate4/internal/api/grpc/server"
	httpserver "github.com/nskforward/gate4/internal/api/http/server"
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

	server := api.NewServer(
		grpcserver.NewUnixServer(),
		grpcserver.NewTCPServer(""),
		httpserver.NewHTTPServer(),
	)

	return server.Start(ctx)
}
