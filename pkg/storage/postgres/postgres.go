package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romandnk/shortener/config"
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type Postgres struct {
	Builder squirrel.StatementBuilderType
	Pool    PgxPool
}

func New(ctx context.Context, cfg config.Postgres) (*Postgres, error) {
	pg := Postgres{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	pgxConf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return &pg, err
	}

	pgxConf.MaxConns = cfg.MaxConns
	pgxConf.MinConns = cfg.MinConns

	db, err := pgxpool.NewWithConfig(ctx, pgxConf)
	if err != nil {
		return &pg, fmt.Errorf("error creating new pgx pool: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return &pg, fmt.Errorf("error connecting pgx pool: %w", err)
	}

	pg.Pool = db

	return &pg, nil
}
