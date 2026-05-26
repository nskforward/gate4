package finam

import (
	"context"
	"time"

	"google.golang.org/grpc/metadata"
)

type Token struct {
	Created time.Time
	Expires time.Time
	Value   string
}

func (t *Token) Expired() bool {
	return time.Since(t.Expires) >= 0
}

func (t *Token) Context(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "Authorization", t.Value)
}
