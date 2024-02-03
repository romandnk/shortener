package interceptor

import (
	"context"
	"errors"
	"github.com/romandnk/shortener/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func LoggingInterceptor(logger logger.Logger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		logErr := err
		if logErr == nil {
			logErr = errors.New("empty")
		}

		logger.Info("Request info GRPC",
			zap.String("method", info.FullMethod),
			zap.String("processing time", duration.String()),
			zap.String("errors", logErr.Error()),
		)

		return resp, err
	}
}
