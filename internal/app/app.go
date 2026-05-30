package app

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/nskforward/gate4/internal/config"
	handler "github.com/nskforward/gate4/internal/domain/handler/user"
	"github.com/nskforward/gate4/pkg/di"
	"github.com/nskforward/gate4/pkg/pb"
)

type App struct {
	container   *di.Container
	userHandler *handler.UserHandler
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

	cfg, ok := di.Resolve[config.Config](app.container)
	if !ok {
		return errors.New("Config not registered")
	}

	fmt.Println("tcp-addr:", cfg.TCPAddr)

	userHandler, ok := di.Resolve[*handler.UserHandler](app.container)
	if !ok {
		return errors.New("UserHandler not registered")
	}
	users, err := userHandler.ListUsers(ctx, &pb.EmptyMessage{})
	if err != nil {
		return err
	}
	fmt.Println(users)
	return errors.New("not implemented")
}
