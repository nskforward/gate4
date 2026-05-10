package streams

import "context"

func MergedContext(ctx1, ctx2 context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx1)
	go func() {
		defer cancel()
		select {
		case <-ctx1.Done():
		case <-ctx2.Done():
		}
	}()
	return ctx
}
