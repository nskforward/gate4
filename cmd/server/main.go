package main

import (
	"context"
	"log/slog"

	"github.com/nskforward/gate4/internal/app"
)

func main() {
	app := app.NewApp()
	err := app.Start(context.Background())
	if err != nil {
		slog.Error("cannot start app", "reason", err.Error())
	}
}
