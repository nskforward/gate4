package main

import (
	"fmt"
	"log/slog"

	"github.com/nskforward/gate4/internal/apps"
	handler "github.com/nskforward/gate4/internal/domain/handler/user"
	repository "github.com/nskforward/gate4/internal/domain/repository/user"
	usecases "github.com/nskforward/gate4/internal/domain/usecases/user"
	"github.com/nskforward/gate4/pkg/di"
)

func main() {
	c := di.NewContainer()
	di.Provide[usecases.UserRepository](c, repository.NewMemoryRepo)
	di.Provide[*usecases.UserUsecases](c, usecases.NewUserUsecases)
	di.Provide[*handler.UserHandler](c, handler.NewUserHandler)
	di.Provide[*apps.App](c, apps.NewApp)

	app := di.Resolve[*apps.App](c)

	users, err := app.Users()
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	err = app.Start()
	if err != nil {
		slog.Error("cannot start app", "reason", err.Error())
	}
}
