package finam

import (
	"context"
	"errors"
)

var (
	ErrPeerAway     = errors.New("peer is away")
	ErrUnauthorized = errors.New("unauthorized")
)

func joinContext(a, b context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(a)
	go func() {
		select {
		case <-ctx.Done():
			return

		case <-b.Done():
			cancel()
		}
	}()
	return ctx, cancel
}
