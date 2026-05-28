package app

import "log/slog"

func (a *App) initDeps() {
	a.unixServer = NewUnixServer(a.container.APIServer())
	a.tcpServer = NewTCPServer(a.container.APIServer(), a.container.TLSConfig(), a.container.Config().Server.TCPAddr)
	slog.SetDefault(a.container.Logger())
}
