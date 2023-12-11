package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	urlservice "github.com/romandnk/shortener/internal/service/url"
	"github.com/romandnk/shortener/internal/storage"
	"github.com/romandnk/shortener/pkg/generator"
	"github.com/romandnk/shortener/pkg/logger"
)

type URL interface {
	CreateURLAlias(ctx context.Context, original string) (string, error)
	GetOriginalByAlias(ctx context.Context, alias string) (string, error)
}

type Services struct {
	URL URL
}

type Dependencies struct {
	Generator generator.Generator
	Repo      *storage.Storage
	Logger    logger.Logger
}

func NewServices(dep Dependencies) *Services {
	return &Services{
		URL: urlservice.NewURLService(dep.Generator, dep.Repo.URL, dep.Logger),
	}
}
