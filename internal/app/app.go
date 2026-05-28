package app

import (
	"log/slog"

	"github.com/nskforward/gate4/internal/di"
)

type App struct {
	container *di.Container
	server    *JointServer
}

func NewApp() *App {
	a := &App{
		container: di.NewContainer(),
	}
	a.initDeps()
	return a
}

func (a *App) Run() error {
	slog.Info("application started",
		slog.String("tcp-addr", a.container.Config().Server.TCPAddr),
	)
	return a.server.Serve()
}
