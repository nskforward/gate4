package app

import (
	"context"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nskforward/gate4/internal/di"
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

	slog.Info("application started",
		slog.String("tcp-addr", a.container.Config().Server.TCPAddr),
	)

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
