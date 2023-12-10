package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath string = "./config/config.yaml"

type Config struct {
	ZapLogger  ZapLogger  `yaml:"zap_logger"`
	Postgres   Postgres   `yaml:"postgres"`
	Redis      Redis      `yaml:"redis"`
	HTTPServer HTTPServer `yaml:"http_server"`
	BaseURL    string     `yaml:"base_url"`
	DBType     string     `yaml:"db_type"`
}

type ZapLogger struct {
	Test             bool     `yaml:"test"`
	Level            string   `yaml:"level"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB"`
	SSLMode  string `yaml:"ssl_mode" env:"POSTGRES_SSLMODE"`
	MaxConns int32  `yaml:"max_conns"`
	MinConns int32  `yaml:"min_conns"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
}

type HTTPServer struct {
	Host            string        `env:"HTTP_SERVER_HOST" env-required:"true"`
	Port            int           `env:"HTTP_SERVER_PORT" env-required:"true"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"3s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"5s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return &cfg, nil
}
