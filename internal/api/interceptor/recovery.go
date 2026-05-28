package interceptor

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			reqID := "-"
			v, ok := ctx.Value(CtxReqID).(string)
			if ok {
				reqID = v
			}

			slog.Error("panic recovered",
				slog.String("request-id", reqID),
				slog.String("method", info.FullMethod),
				slog.String("reason", fmt.Sprint(r)),
			)

			desc := "internal server error"
			if ok {
				desc = fmt.Sprintf("request-id:%s", reqID)
			}

			err = status.Errorf(codes.Internal, desc)
		}
	}()

	return handler(ctx, req)
}
