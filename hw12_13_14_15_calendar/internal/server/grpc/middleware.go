package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Lanworm/OTUS_GO/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
)

func LoggingRequest(logger *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		logger.ServerLog(fmt.Sprintf(
			"[%s] %s",
			time.Now().Format(time.RFC3339),
			info.FullMethod,
		))

		return handler(ctx, req)
	}
}
