package grpc

import (
	"github.com/romandnk/shortener/config"
	urlpb "github.com/romandnk/shortener/internal/server/grpc/url/pb"
	"github.com/romandnk/shortener/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net"
	"strconv"
)

type ServerGRPC struct {
	srv     *grpc.Server
	cfg     config.GRPCServer
	handler *HandlerGRPC
}

func NewServerGRPC(handler *HandlerGRPC, logger logger.Logger, cfg config.GRPCServer) *ServerGRPC {
	serverOptions := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle,
			MaxConnectionAge:  cfg.MaxConnectionAge,
			Time:              cfg.Time,
			Timeout:           cfg.Timeout,
		}),
	}

	srv := grpc.NewServer(serverOptions...)

	return &ServerGRPC{
		srv:     srv,
		cfg:     cfg,
		handler: handler,
	}
}

func (s *ServerGRPC) Start() error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)))
	if err != nil {
		return err
	}

	urlpb.RegisterEventServiceServer(s.srv, s.handler.URLHandler)

	return s.srv.Serve(lsn)
}

func (s *ServerGRPC) Stop() {
	s.srv.GracefulStop()
}
