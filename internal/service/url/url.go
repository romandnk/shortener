package urlservice

import (
	"context"
	"errors"
	"github.com/romandnk/shortener/internal/constant"
	"github.com/romandnk/shortener/internal/entity"
	"github.com/romandnk/shortener/internal/storage"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	"github.com/romandnk/shortener/pkg/generate"
	"github.com/romandnk/shortener/pkg/logger"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"unicode/utf8"
)

type URLService struct {
	baseURL url.URL
	url     storage.URL
	logger  logger.Logger
}

func NewURLService(baseURL url.URL, url storage.URL, logger logger.Logger) *URLService {
	return &URLService{
		baseURL: baseURL,
		url:     url,
		logger:  logger,
	}
}

func (s *URLService) CreateShortURL(ctx context.Context, original string) (string, error) {
	original = strings.TrimSpace(original)
	if original == "" {
		s.logger.Error("URLService.CreateShortURL", zap.String("error", ErrEmptyOriginalURL.Error()))
		return "", ErrEmptyOriginalURL
	}

	if utf8.RuneCountInString(original) > 2048 {
		s.logger.Error("URLService.CreateShortURL", zap.String("error", ErrOriginalURLTooLong.Error()))
		return "", ErrOriginalURLTooLong
	}

	_, err := url.ParseRequestURI(original)
	if err != nil {
		s.logger.Error("URLService.CreateShortURL", zap.String("error", err.Error()))
		return "", ErrInvalidOriginalURL
	}

	alias, err := generate.Random()
	if err != nil {
		s.logger.Error("URLService.CreateShortURL - generate.Random()", zap.String("error", err.Error()))
		return "", ErrInternalError
	}

	URL := entity.URL{
		Origin: original,
		Alias:  alias,
	}

	err = s.url.CreateURL(ctx, URL)
	if err != nil {
		if errors.Is(err, storageerrors.ErrOriginalURLExists) || errors.Is(err, storageerrors.ErrURLAliasExists) {
			return "", err
		}
		s.logger.Error("URLService.CreateShortURL - s.url.CreateURL", zap.String("error", err.Error()))
		return "", ErrInternalError
	}

	s.logger.Info("URLService.CreateShortURL - alias was created successfully")

	shortURL := s.baseURL.JoinPath(alias)

	return shortURL.String(), nil
}

func (s *URLService) GetOriginalByAlias(ctx context.Context, alias string) (string, error) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		s.logger.Error("URLService.GetOriginalByAlias", zap.String("error", ErrEmptyURLAlias.Error()))
		return "", ErrEmptyURLAlias
	}

	if utf8.RuneCountInString(alias) != constant.AliasLength {
		s.logger.Error("URLService.GetOriginalByAlias", zap.String("error", ErrInvalidAliasFormat.Error()))
		return "", ErrInvalidAliasFormat
	}

	original, err := s.url.GetOriginalByAlias(ctx, alias)
	if err != nil {
		if errors.Is(err, storageerrors.ErrURLAliasNotFound) {
			return "", ErrOriginalURLNotFound
		}
		s.logger.Error("URLService.GetOriginalByAlias - s.url.GetOriginalByAlias", zap.String("error", err.Error()))
		return "", ErrInternalError
	}

	return original, nil
}
