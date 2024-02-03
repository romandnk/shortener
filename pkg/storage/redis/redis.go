package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
)

type Config struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
}

func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Password: cfg.Password,
		DB:       0,
	})

	ping := rdb.Ping(ctx)
	if err := ping.Err(); err != nil {
		return rdb, err
	}

	return rdb, nil
}
