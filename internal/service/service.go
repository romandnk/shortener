package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	urlservice "github.com/romandnk/shortener/internal/service/url"
	"github.com/romandnk/shortener/internal/storage"
	"github.com/romandnk/shortener/pkg/generator"
	"github.com/romandnk/shortener/pkg/logger"
	"go.uber.org/fx"
)

var Module = fx.Module("services", fx.Provide(NewServices))

type URL interface {
	CreateURLAlias(ctx context.Context, original string) (string, error)
	GetOriginalByAlias(ctx context.Context, alias string) (string, error)
}

type Services struct {
	URL URL
}

func NewServices(generator generator.Generator, repo *storage.Storage, logger logger.Logger) *Services {
	return &Services{
		URL: urlservice.NewURLService(generator, repo.URL, logger),
	}
}
