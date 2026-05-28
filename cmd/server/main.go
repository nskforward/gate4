package main

import (
	"log/slog"

	"github.com/nskforward/gate4/internal/app"
)

func main() {
	a := app.NewApp()

	err := a.Run()
	if err != nil {
		slog.Error("server app exited", "reason", err.Error())
	}
}
