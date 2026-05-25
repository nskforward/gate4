package tools

import (
	"context"
)

func MergeContext(a, b context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-a.Done():
			cancel()
			return
		case <-b.Done():
			cancel()
			return
		}
	}()

	return ctx, cancel
}
