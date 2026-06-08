package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/nskforward/gate4/internal/api/http/server/interceptor"
	"github.com/nskforward/gate4/internal/config"
	"github.com/nskforward/gate4/internal/domain/service"
)

type HTTPServer struct {
	transport *http.Server
	auth      *interceptor.Auth
}

func NewHTTPServer(cfg config.Config, userService *service.UserService, tokenService *service.TokenService) (*HTTPServer, error) {
	tlsConfig, err := NewTLSConfig(cfg)
	if err != nil {
		return nil, err
	}

	auth := interceptor.NewAuth(userService, tokenService)

	defaultRoute := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("http api call", "method", r.Method, "path", r.RequestURI)
		http.Error(w, "route not found", 404)
	}

	router := http.NewServeMux()
	router.HandleFunc("/", auth.Auth(defaultRoute))

	return &HTTPServer{
		transport: &http.Server{
			Addr:      cfg.HTTP.Addr,
			Handler:   router,
			TLSConfig: tlsConfig,
		},
		auth: auth,
	}, nil
}

func (s *HTTPServer) Start(ctx context.Context) error {
	err := s.transport.ListenAndServeTLS("", "")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.transport.Shutdown(context.Background())
		close(stopped)
	}()

	select {
	case <-stopped:

	case <-ctx.Done():
		s.transport.Close()
	}
	return nil
}
