package server

import (
	"log/slog"
	"net/http"

	"github.com/nskforward/gate4/internal/api/http/server/handler"
	"github.com/nskforward/gate4/internal/api/http/server/interceptor"
	"github.com/nskforward/gate4/internal/domain/service"
)

func NewRouter(userService *service.UserService, tokenService *service.TokenService) *http.ServeMux {
	mux := http.NewServeMux()

	auth := interceptor.NewAuth(tokenService)

	mux.HandleFunc("/api/whoami", auth.Auth(handler.Whoami))
	mux.HandleFunc("/", auth.Auth(defaultRoute))

	return mux
}

func defaultRoute(w http.ResponseWriter, r *http.Request) {
	slog.Info("http api call", "method", r.Method, "path", r.RequestURI)
	http.Error(w, "route not found", 404)
}
