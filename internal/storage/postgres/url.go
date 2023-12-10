package postgresstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/romandnk/shortener/internal/constant"
	"github.com/romandnk/shortener/internal/entity"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	"github.com/romandnk/shortener/pkg/storage/postgres"
	"strings"
)

type URLRepo struct {
	*postgres.Postgres
}

func NewURLRepo(db *postgres.Postgres) *URLRepo {
	return &URLRepo{db}
}

func (r *URLRepo) CreateURL(ctx context.Context, url entity.URL) error {
	sql, args, _ := r.Builder.
		Insert(constant.URLSTable).
		Columns("original", "alias").
		Values(url.Original, url.Alias).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Detail, "Key (original)") {
					return storageerrors.ErrOriginalURLExists
				}
				if strings.Contains(pgErr.Detail, "Key (alias)") {
					return storageerrors.ErrURLAliasExists
				}
			}
		}
		return fmt.Errorf("URLRepo.CreateShortURL - r.Pool.Exec: %v", err)
	}

	return nil
}

func (r *URLRepo) GetOriginalByAlias(ctx context.Context, alias string) (string, error) {
	sql, args, _ := r.Builder.
		Select("original").
		From(constant.URLSTable).
		Where(squirrel.Eq{"alias": alias}).
		ToSql()

	var original string
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&original)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return original, storageerrors.ErrURLAliasNotFound
		}
		return original, fmt.Errorf("URLRepo.GetOriginalByAlias - r.Pool.Query: %v", err)
	}

	return original, nil
}
