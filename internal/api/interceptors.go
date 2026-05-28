package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nskforward/gate4/pkg/console"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	if err != nil {
		console.LogDebug(fmt.Sprint("failed api call ", info.FullMethod, " ", duration, " ", err))
	} else {
		console.LogDebug(fmt.Sprint("successful api call ", info.FullMethod, " ", duration))
	}

	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			console.LogError("panic recovered", errors.New(fmt.Sprint(r)))
			err = status.Errorf(codes.Internal, "Internal Server Error")
		}
	}()

	return handler(ctx, req)
}
