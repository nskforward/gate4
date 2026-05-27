package app

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
	a := &App{
		container: di.NewContainer(),
	}
	a.initDeps()
	return a
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errorc := make(chan error, 2)
	defer close(errorc)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := a.unixServer.Serve()
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := a.tcpServer.Serve()
		if err != nil && err != grpc.ErrServerStopped {
			errorc <- err
		}
	}()

	console.LogInfo("application started")

	select {
	case err := <-errorc:
		a.unixServer.Close()
		a.tcpServer.Close()
		return err

	case <-ctx.Done():
		a.unixServer.Close()
		a.tcpServer.Close()
		return nil
	}
}
