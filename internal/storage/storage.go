package storage

import (
	"context"
	"github.com/romandnk/shortener/internal/entity"
	postgresstorage "github.com/romandnk/shortener/internal/storage/postgres"
	"github.com/romandnk/shortener/pkg/storage/postgres"
	"go.uber.org/fx"
)

//go:generate mockgen -source=storage.go -destination=mock/mock.go storage

var Module = fx.Module("storage", fx.Provide(NewStorage))

type URL interface {
	CreateURL(ctx context.Context, url entity.URL) error
	GetOriginalByAlias(ctx context.Context, alias string) (string, error)
}

type Storage struct {
	URL URL
}

func NewStorage(db *postgres.Postgres) (*Storage, error) {
	var storage Storage

	//switch v := db.(type) {
	//case *postgres.Postgres:
	storage = Storage{
		URL: postgresstorage.NewURLRepo(db),
	}
	//case *redis.Redis:
	//	storage = Storage{
	//		URL: redisstorage.NewURLRepo(v),
	//	}
	//default:
	//	return &storage, storageerrors.ErrInvalidDB
	//}

	return &storage, nil
}
