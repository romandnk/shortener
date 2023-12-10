package service

import (
	"context"
	urlservice "github.com/romandnk/shortener/internal/service/url"
	"github.com/romandnk/shortener/internal/storage"
	"github.com/romandnk/shortener/pkg/logger"
	"net/url"
)

type URL interface {
	CreateShortURL(ctx context.Context, original string) (string, error)
	GetShortByOrigin(ctx context.Context, origin string) (string, error)
}

type Services struct {
	URL URL
}

type Dependencies struct {
	BaseURL url.URL
	Repo    *storage.Storage
	Logger  logger.Logger
}

func NewServices(dep Dependencies) *Services {
	return &Services{
		URL: urlservice.NewURLService(dep.BaseURL, dep.Repo.URL, dep.Logger),
	}
}
