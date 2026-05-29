package apps

import "errors"

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Start() error {
	return errors.New("not implemented")
}

type Starter interface {
	Start() error
}
