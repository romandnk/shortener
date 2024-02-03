package app

import (
	"context"
	"github.com/romandnk/shortener/config"
	"github.com/romandnk/shortener/internal/constant"
	"github.com/romandnk/shortener/internal/server/grpc/interceptor"
	urlgrpc "github.com/romandnk/shortener/internal/server/grpc/url"
	"github.com/romandnk/shortener/internal/server/http/middleware"
	v1 "github.com/romandnk/shortener/internal/server/http/v1"
	"github.com/romandnk/shortener/internal/service"
	"github.com/romandnk/shortener/internal/storage"
	"github.com/romandnk/shortener/pkg/generator"
	"github.com/romandnk/shortener/pkg/grpcserver"
	"github.com/romandnk/shortener/pkg/httpserver"
	zaplogger "github.com/romandnk/shortener/pkg/logger/zap"
	"github.com/romandnk/shortener/pkg/storage/postgres"
	"github.com/romandnk/shortener/pkg/storage/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

//	@title			URL shortener project
//	@version		1.0
//	@description	Swagger API for Golang Project URL Shortener.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API [Roman] Support
//	@license.name	romandnk
//	@license.url	https://github.com/romandnk/shortener

// @BasePath	/api/v1/

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	// initializing config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
	}

	// initializing zap logger
	logger, err := zaplogger.NewLogger(cfg.ZapLogger)
	if err != nil {
		log.Fatalf("error initializing zap logger: %s", err.Error())
	}

	logger.Info("using zap logger")

	// initializing storage
	var st *storage.Storage
	switch cfg.DBType {
	case constant.POSTGRES:
		pg, err := postgres.New(ctx, cfg.Postgres)
		if err != nil {
			logger.Fatal("error initializing postgres db", zap.Error(err))
		}
		defer pg.Pool.Close()
		st, err = storage.NewStorage(pg)
		if err != nil {
			pg.Pool.Close()
			logger.Fatal("error creating storage", zap.Error(err))
		}

		logger.Info("using postgres storage",
			zap.String("host", cfg.Postgres.Host),
			zap.Int("port", cfg.Postgres.Port),
		)
	case constant.REDIS:
		r, err := redis.New(ctx, cfg.Redis)
		if err != nil {
			logger.Fatal("error initializing redis db", zap.Error(err))
		}
		defer r.Close()
		st, err = storage.NewStorage(r)
		if err != nil {
			r.Close()
			logger.Fatal("error creating storage", zap.Error(err))
		}

		logger.Info("using redis storage",
			zap.String("host", cfg.Redis.Host),
			zap.Int("port", cfg.Redis.Port),
		)
	default:
		logger.Fatal("invalid db type", zap.String("db type", cfg.DBType))
	}

	// initializing services
	dep := service.Dependencies{
		Generator: generator.NewGen(constant.AliasLength),
		Repo:      st,
		Logger:    logger,
	}
	services := service.NewServices(dep)

	// initializing middlewares
	mw := middleware.New(logger)

	// initializing handlers
	HTTPHandler := v1.NewHandler(services, mw)

	// initializing servers
	HTTPServer := httpserver.NewServer(cfg.HTTPServer, HTTPHandler.InitRoutes())
	GRPCServer := grpcserver.NewServer(cfg.GRPCServer,
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor(logger)),
	)

	// register grpc
	urlgrpc.Register(GRPCServer.Srv, services.URL)

	go func() {
		<-ctx.Done()

		if err := HTTPServer.Stop(ctx); err != nil {
			logger.Error("error stopping HTTP server",
				zap.String("host", cfg.HTTPServer.Host),
				zap.Int("port", cfg.HTTPServer.Port),
			)
		}

		GRPCServer.Stop()

		logger.Info("app is stopped")
	}()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		logger.Info("http server is running...",
			zap.String("address", net.JoinHostPort(cfg.HTTPServer.Host, strconv.Itoa(cfg.HTTPServer.Port))))
		if err := HTTPServer.Start(); err != nil {
			logger.Error("error starting HTTP server",
				zap.String("address", net.JoinHostPort(cfg.HTTPServer.Host, strconv.Itoa(cfg.HTTPServer.Port))))
		}
	}()
	go func() {
		defer wg.Done()
		logger.Info("grpc server is running...",
			zap.String("address", net.JoinHostPort(cfg.GRPCServer.Host, strconv.Itoa(cfg.GRPCServer.Port))))
		if err := GRPCServer.Start(); err != nil {
			logger.Error("error starting GRPC server",
				zap.String("address", net.JoinHostPort(cfg.GRPCServer.Host, strconv.Itoa(cfg.GRPCServer.Port))))
		}
	}()
	wg.Wait()
}
