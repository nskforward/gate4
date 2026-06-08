package interceptor

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	reqID := uuid.NewString()
	nextCtx := context.WithValue(ctx, CtxReqID, reqID)

	start := time.Now()

	resp, err := handler(nextCtx, req)
	duration := time.Since(start)

	if err != nil {
		slog.Warn("grpc api call",
			slog.String("req-id", reqID),
			slog.String("method", info.FullMethod),
			slog.Int64("time-ms", duration.Milliseconds()),
			slog.String("reason", err.Error()),
		)
		return resp, err
	}

	slog.Info("grpc api call",
		slog.String("req-id", reqID),
		slog.String("method", info.FullMethod),
		slog.Int64("time-ms", duration.Milliseconds()),
	)

	return resp, err
}
