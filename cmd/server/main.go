package main

import (
	"log/slog"

	"github.com/nskforward/gate4/internal/apps"
	"github.com/nskforward/gate4/pkg/di"
)

func main() {
	c := di.NewContainer()
	di.Provide[*apps.App](c, apps.NewApp)

	app := di.Resolve[*apps.App](c)

	err := app.Start()
	if err != nil {
		slog.Error("cannot start app", "reason", err.Error())
	}
}
