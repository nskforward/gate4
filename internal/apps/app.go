package apps

import (
	"context"
	"errors"

	handler "github.com/nskforward/gate4/internal/domain/handler/user"
	"github.com/nskforward/gate4/pkg/pb"
)

type App struct {
	userHandler *handler.UserHandler
}

func NewApp(userHandler *handler.UserHandler) *App {
	return &App{
		userHandler: userHandler,
	}
}

func (app *App) Users() (*pb.UserList, error) {
	return app.userHandler.ListUsers(context.Background(), &pb.EmptyMessage{})
}

func (app *App) Start() error {
	return errors.New("not implemented")
}

type Starter interface {
	Start() error
}
