package app

import "log/slog"

func (a *App) initDeps() {
	slog.SetDefault(a.container.Logger())
	a.server = NewJointServer(a.container.APIServer(), a.container.TLSConfig(), a.container.Config().Server.TCPAddr)
}
