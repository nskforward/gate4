package main

import (
	"github.com/nskforward/gate4/internal/app/server"
	"github.com/nskforward/gate4/pkg/console"
)

func main() {
	app := server.NewApp()

	err := app.Run()
	if err != nil {
		console.LogError("server app exited", err)
	}
}
