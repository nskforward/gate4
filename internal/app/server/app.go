package server

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nskforward/gate4/internal/di"
	"github.com/nskforward/gate4/pkg/console"
	"google.golang.org/grpc"
)

type App struct {
	container  *di.Container
	unixServer *UnixServer
	tcpServer  *TCPServer
}

func NewApp() *App {
	app := &App{
		container: di.NewContainer(),
	}
	app.initDeps()
	console.LogInfo("server deps successfully initialized")
	return app
}

func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errorc := make(chan error, 2)
	defer close(errorc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := app.unixServer.Serve()
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := app.tcpServer.Serve()
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()

	select {
	case err := <-errorc:
		app.unixServer.Close()
		app.tcpServer.Close()
		return err

	case <-ctx.Done():
		app.unixServer.Close()
		app.tcpServer.Close()
		return nil
	}
}
