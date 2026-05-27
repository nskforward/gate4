package server

func (app *App) initDeps() {
	app.unixServer = NewUnixServer(app.container.APIServer())
	app.tcpServer = NewTCPServer(app.container.APIServer(), app.container.TLSConfig(), app.container.Config().Server.TCPAddr)
}
