package redisstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/romandnk/shortener/internal/constant"
	"github.com/romandnk/shortener/internal/entity"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
)

type URLRepo struct {
	client *redis.Client
}

func NewURLRepo(client *redis.Client) *URLRepo {
	return &URLRepo{client: client}
}

func (r *URLRepo) CreateURL(ctx context.Context, url entity.URL) error {
	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		originExists, err := tx.SetNX(ctx, url.Original, url.Alias, constant.ZeroTTL).Result()
		if err != nil {
			return fmt.Errorf("URLRepo.CreateURL - tx.SetNX - 1: %v", err)
		}
		if !originExists {
			return storageerrors.ErrOriginalURLExists
		}

		aliasExists, err := tx.SetNX(ctx, url.Alias, url.Original, constant.ZeroTTL).Result()
		if err != nil {
			return fmt.Errorf("URLRepo.CreateURL - tx.SetNX - 2: %v", err)
		}
		if !aliasExists {
			return storageerrors.ErrURLAliasExists
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, storageerrors.ErrOriginalURLExists) || errors.Is(err, storageerrors.ErrURLAliasExists) {
			return err
		}
		return fmt.Errorf("URLRepo.CreateURL - r.client.Watch: %v", err)
	}

	return nil
}

func (r *URLRepo) GetOriginalByAlias(ctx context.Context, alias string) (string, error) {
	original, err := r.client.Get(ctx, alias).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", storageerrors.ErrURLAliasNotFound
		}
		return "", fmt.Errorf("URLRepo.GetOriginalByAlias - r.client.Get: %v", err)
	}
	return original, nil
}
