package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/romandnk/shortener/config"
	"net"
	"strconv"
)

func New(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
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
