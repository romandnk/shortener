package redisstorage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/romandnk/shortener/internal/entity"
)

type URLRepo struct {
	client *redis.Client
}

func NewURLRepo(client *redis.Client) *URLRepo {
	return &URLRepo{client: client}
}

func (r *URLRepo) CreateURL(ctx context.Context, url entity.URL) error {
	return nil
}
func (r *URLRepo) GetOriginalByAlias(ctx context.Context, alias string) (string, error) {
	return "", nil
}
