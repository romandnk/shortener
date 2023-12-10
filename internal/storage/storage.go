package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/romandnk/shortener/internal/entity"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	postgresstorage "github.com/romandnk/shortener/internal/storage/postgres"
	redisstorage "github.com/romandnk/shortener/internal/storage/redis"
	"github.com/romandnk/shortener/pkg/storage/postgres"
)

type URL interface {
	CreateURL(ctx context.Context, url entity.URL) error
	GetOriginalByAlias(ctx context.Context, alias string) (string, error)
}

type Storage struct {
	URL URL
}

func NewStorage(db any) (*Storage, error) {
	var storage Storage

	switch v := db.(type) {
	case *postgres.Postgres:
		storage = Storage{
			URL: postgresstorage.NewURLRepo(v),
		}
	case *redis.Client:
		storage = Storage{
			URL: redisstorage.NewURLRepo(v),
		}
	default:
		return &storage, storageerrors.ErrInvalidDB
	}

	return &storage, nil
}
