package api

import (
	"context"
	"sync"
	"time"
)

type ServerInterface interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Server struct {
	children []ServerInterface
}

func NewServer(children ...ServerInterface) *Server {
	return &Server{
		children: children,
	}
}

func (s *Server) Start(ctx context.Context) error {
	errorc := make(chan error, len(s.children))
	defer close(errorc)

	for _, child := range s.children {
		go func(ser ServerInterface) {
			err := ser.Start(ctx)
			if err != nil {
				errorc <- err
			}
		}(child)
	}

	select {

	case <-ctx.Done():
		stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return s.Stop(stopCtx)

	case err := <-errorc:
		stopCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		s.Stop(stopCtx)
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	errorc := make(chan error, len(s.children))

	var wg sync.WaitGroup

	for _, child := range s.children {
		wg.Add(1)
		go func(ser ServerInterface) {
			defer wg.Done()
			err := ser.Stop(ctx)
			if err != nil {
				errorc <- err
			}
		}(child)
	}

	wg.Wait()

	close(errorc)

	return <-errorc
}
