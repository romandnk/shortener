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
	"github.com/romandnk/shortener/pkg/logger"
	zaplogger "github.com/romandnk/shortener/pkg/logger/zap"
	"github.com/romandnk/shortener/pkg/storage/postgres"
	"github.com/romandnk/shortener/pkg/storage/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
)

func NewApp() fx.Option {
	return fx.Options(
		MutualContextModule(),
		config.Module,
		LoggerModule(),
		StartCheckModule(),
		PostgresModule(),
		RedisModule(),
		ShortURLGeneratorModule(),
		storage.Module,
		middleware.Module,
		service.Module,
		v1.Module,
		HTTPServerModule(),
		GRPCServerModule(),

		CheckInitializedModules(),
	)
}

func LoggerModule() fx.Option {
	return fx.Module("logger",
		fx.Provide(
			fx.Annotate(
				zaplogger.NewLogger,
				fx.As(new(logger.Logger))),
			func(cfg *config.Config) zaplogger.Config {
				return cfg.ZapLogger
			},
		),
	)
}

func StartCheckModule() fx.Option {
	return fx.Module("start with atomic",
		fx.Provide(func() *atomic.Bool {
			return &atomic.Bool{}
		}),
	)
}

func PostgresModule() fx.Option {
	return fx.Module("postgres",
		fx.Provide(
			func(cfg *config.Config) postgres.Config {
				return cfg.Postgres
			},
			postgres.New,
		),
		fx.Invoke(func(lc fx.Lifecycle, pg *postgres.Postgres) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					pg.Close()
					return nil
				},
			})
		}),
	)
}

func RedisModule() fx.Option {
	return fx.Module("redis",
		fx.Provide(
			func(cfg *config.Config) redis.Config {
				return cfg.Redis
			},
			redis.New,
		),
		fx.Invoke(func(lc fx.Lifecycle, rdb *redis.Redis) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return rdb.Close()
				},
			})
		}),
	)
}

func MutualContextModule() fx.Option {
	return fx.Module("context",
		fx.Provide(func() (context.Context, context.CancelFunc) {
			ctx, cancel := signal.NotifyContext(context.Background(),
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGHUP,
			)
			return ctx, cancel
		}))
}

func ShortURLGeneratorModule() fx.Option {
	return fx.Module("generator",
		fx.Provide(
			fx.Annotate(
				generator.NewGen,
				fx.As(new(generator.Generator)),
			),
			func() int {
				return constant.AliasLength
			},
		),
	)
}

func HTTPServerModule() fx.Option {
	return fx.Module("http server",
		fx.Provide(
			func(cfg *config.Config) httpserver.Config {
				return cfg.HTTPServer
			},
			httpserver.NewServer,
		),
		fx.Invoke(func(lc fx.Lifecycle, srv *httpserver.Server, cfg httpserver.Config, logger logger.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := srv.Start(); err != nil {
							logger.Error("error starting HTTP server",
								zap.String("address", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return srv.Stop(ctx)
				},
			})
		}),
	)
}

func GRPCServerModule() fx.Option {
	return fx.Module("grpc server",
		fx.Provide(
			func(cfg *config.Config) grpcserver.Config {
				return cfg.GRPCServer
			},
			func(logger logger.Logger) []grpc.ServerOption {
				return []grpc.ServerOption{
					grpc.UnaryInterceptor(interceptor.LoggingInterceptor(logger)),
				}
			},
			grpcserver.NewServer,
		),
		fx.Invoke(
			func(srv *grpcserver.Server, services *service.Services) {
				urlgrpc.Register(srv.Srv, services.URL)
			},
			func(lc fx.Lifecycle, srv *grpcserver.Server, cfg grpcserver.Config, logger logger.Logger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							if err := srv.Start(); err != nil {
								logger.Error("error starting GRPC server",
									zap.String("address", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))))
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						srv.Stop()
						return nil
					},
				})
			}),
	)
}

func CheckInitializedModules() fx.Option {
	return fx.Module("check modules",
		fx.Invoke(
			func(cfg *config.Config) {},
			func(logger logger.Logger) {},
			func(logger logger.Logger) {},
			func(ok *atomic.Bool, pg *postgres.Postgres) {
				if pg.Pool != nil {
					ok.Store(true)
				}
			},
			func(storage *storage.Storage) {},
			func(mw *middleware.MW) {},
			func(service *service.Services) {},
			func(h http.Handler) {},
			func(srv *httpserver.Server) {},
			func(srv *grpcserver.Server) {},
		),
	)
}
