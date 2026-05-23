package servers

import (
	"context"
	"sync"
)

type Pool struct {
	servers []RunFunc
}

type RunFunc func(context.Context) error

func (m *Pool) Add(f RunFunc) {
	if m.servers == nil {
		m.servers = make([]RunFunc, 0, 4)
	}
	m.servers = append(m.servers, f)
}

func (m *Pool) Run(ctx context.Context) error {
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errorc := make(chan error, len(m.servers))
	defer close(errorc)

	for _, s := range m.servers {
		wg.Add(1)
		go func(f RunFunc) {
			defer wg.Done()
			err := f(innerCtx)
			if err != nil {
				errorc <- err
			}
			cancel()
		}(s)
	}

	wg.Wait()

	select {
	case err := <-errorc:
		return err
	default:
		return nil
	}
}
