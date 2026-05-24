package tools

import (
	"context"
	"errors"
)

func MergeContext(server, peer context.Context) (context.Context, context.CancelCauseFunc) {
	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-server.Done():
			cancel(errors.New("server away"))
			return
		case <-peer.Done():
			cancel(errors.New("peer away"))
			return
		}
	}()

	return ctx, cancel
}
