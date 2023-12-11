package grpc

import (
	urlgrpc "github.com/romandnk/shortener/internal/server/grpc/url"
	urlpb "github.com/romandnk/shortener/internal/server/grpc/url/pb"
	"github.com/romandnk/shortener/internal/service"
	"github.com/romandnk/shortener/pkg/logger"
)

type HandlerGRPC struct {
	URLHandler urlpb.EventServiceServer
	service    *service.Services
	logger     logger.Logger
}

func NewHandlerGRPC(services *service.Services, logger logger.Logger) *HandlerGRPC {
	return &HandlerGRPC{
		URLHandler: urlgrpc.NewURLHandler(services.URL),
		service:    services,
		logger:     logger,
	}
}
