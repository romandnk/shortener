package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/romandnk/shortener/pkg/grpcserver"
	"github.com/romandnk/shortener/pkg/httpserver"
	zaplogger "github.com/romandnk/shortener/pkg/logger/zap"
	"github.com/romandnk/shortener/pkg/storage/postgres"
	"github.com/romandnk/shortener/pkg/storage/redis"
	"go.uber.org/fx"
)

var Module = fx.Module("config", fx.Provide(NewConfig))

const configPath string = "./config/config.yaml"

type Config struct {
	//fx.Out     `yaml:"-"`
	ZapLogger  zaplogger.Config  `yaml:"zap_logger"`
	Postgres   postgres.Config   `yaml:"postgres"`
	Redis      redis.Config      `yaml:"redis"`
	HTTPServer httpserver.Config `yaml:"http_server"`
	GRPCServer grpcserver.Config `yaml:"grpc_server"`
	DBType     string            `yaml:"db_type"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return &cfg, err
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("error updating env: %w", err)
	}

	return &cfg, nil
}
