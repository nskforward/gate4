package main

import (
	"fmt"
	"log/slog"

	"github.com/nskforward/gate4/internal/apps"
	"github.com/nskforward/gate4/internal/domain/handler"
	"github.com/nskforward/gate4/internal/domain/repository/user"
	"github.com/nskforward/gate4/internal/domain/usecase"
	"github.com/nskforward/gate4/pkg/di"
)

func main() {
	c := di.NewContainer()
	di.Provide[usecase.UserRepo](c, user.NewMemoryRepo)
	di.Provide[*usecase.UserUsecase](c, usecase.NewUserUsecase)
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
