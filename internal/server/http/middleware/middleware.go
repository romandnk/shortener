package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/shortener/pkg/logger"
	"go.uber.org/zap"
	"time"
)

type MW struct {
	logger logger.Logger
}

func New(logger logger.Logger) *MW {
	return &MW{
		logger: logger,
	}
}

func (m *MW) Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		duration := time.Since(start)

		info := requestInformation(ctx.Request, duration)

		code := ctx.Writer.Status()

		msg := "HTTP requests"

		m.logger.Info(msg,
			zap.String("client ip", info.ClientIP),
			zap.String("date", info.Date),
			zap.String("method", info.Method),
			zap.String("method path", info.Path),
			zap.String("HTTP version", info.HTTPVersion),
			zap.Int("status code", code),
			zap.String("processing time", info.Latency),
			zap.String("user agent", info.UserAgent),
		)
	}
}
