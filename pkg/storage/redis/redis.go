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

type Redis struct {
	Client *redis.Client
}

func New(ctx context.Context, cfg Config) (*Redis, error) {
	r := &Redis{}

	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Password: cfg.Password,
		DB:       0,
	})

	ping := rdb.Ping(ctx)
	if err := ping.Err(); err != nil {
		return r, err
	}

	r.Client = rdb

	return r, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
