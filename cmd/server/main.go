package main

import (
	"github.com/nskforward/gate4/internal/app"
	"github.com/nskforward/gate4/pkg/console"
)

func main() {
	a := app.NewApp()

	err := a.Run()
	if err != nil {
		console.LogError("server app exited", err)
	}
}
