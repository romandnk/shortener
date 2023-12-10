package postgresstorage

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/romandnk/shortener/internal/constant"
	"github.com/romandnk/shortener/internal/entity"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	"github.com/romandnk/shortener/pkg/storage/postgres"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestURLRepo_CreateURL(t *testing.T) {
	type input struct {
		sql   string
		args  []any
		error *pgconn.PgError
	}

	type mockBehaviour func(m pgxmock.PgxPoolIface, input input)

	testCases := []struct {
		name              string
		url               entity.URL
		mockBehaviour     mockBehaviour
		expectedExecError *pgconn.PgError
		expectedError     error
	}{
		{
			name: "OK",
			url: entity.URL{
				Original: "http://test.com",
				Alias:    "testtest11",
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, input input) {
				m.ExpectExec(regexp.QuoteMeta(input.sql)).
					WithArgs(input.args...).
					WillReturnResult(pgxmock.NewResult("insert", 1))
			},
		},
		{
			name: "original url already exists",
			url: entity.URL{
				Original: "http://test.com",
				Alias:    "testtest11",
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, input input) {
				m.ExpectExec(regexp.QuoteMeta(input.sql)).
					WithArgs(input.args...).
					WillReturnError(input.error)
			},
			expectedExecError: &pgconn.PgError{
				Code:   "23505",
				Detail: "Key (original)=(http://test.com) already exists.",
			},
			expectedError: storageerrors.ErrOriginalURLExists,
		},
		{
			name: "url alias already exists",
			url: entity.URL{
				Original: "http://test.com",
				Alias:    "testtest11",
			},
			mockBehaviour: func(m pgxmock.PgxPoolIface, input input) {
				m.ExpectExec(regexp.QuoteMeta(input.sql)).
					WithArgs(input.args...).
					WillReturnError(input.error)
			},
			expectedExecError: &pgconn.PgError{
				Code:   "23505",
				Detail: "Key (alias)=(testtest11) already exists.",
			},
			expectedError: storageerrors.ErrURLAliasExists,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			db := postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mock,
			}

			sql, args, _ := db.Builder.
				Insert(constant.URLSTable).
				Columns("original", "alias").
				Values(tc.url.Original, tc.url.Alias).
				ToSql()

			ctx := context.Background()

			in := input{
				sql:   sql,
				args:  args,
				error: tc.expectedExecError,
			}

			tc.mockBehaviour(mock, in)

			urlStorage := NewURLRepo(&db)

			err = urlStorage.CreateURL(ctx, tc.url)
			require.ErrorIs(t, err, tc.expectedError)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
}

func TestURLRepo_GetOriginalByAlias(t *testing.T) {
	type input struct {
		sql           string
		args          []any
		rows          *pgxmock.Rows
		expectedError error
	}

	type mockBehaviour func(m pgxmock.PgxPoolIface, input input)

	testCases := []struct {
		name               string
		inputAlias         string
		rows               *pgxmock.Rows
		mockBehaviour      mockBehaviour
		expectedOriginal   string
		expectedQueryError error
		expectedError      error
	}{
		{
			name:       "OK",
			inputAlias: "testtest11",
			rows:       pgxmock.NewRows([]string{"original"}).AddRow("http://google.com/"),
			mockBehaviour: func(m pgxmock.PgxPoolIface, input input) {
				m.ExpectQuery(regexp.QuoteMeta(input.sql)).
					WithArgs(input.args...).
					WillReturnRows(input.rows)
			},
			expectedOriginal: "http://google.com/",
		},
		{
			name:       "alias is not found",
			inputAlias: "testtest11",
			rows:       pgxmock.NewRows([]string{"original"}).AddRow("http://google.com/"),
			mockBehaviour: func(m pgxmock.PgxPoolIface, input input) {
				m.ExpectQuery(regexp.QuoteMeta(input.sql)).
					WithArgs(input.args...).
					WillReturnError(input.expectedError)
			},
			expectedQueryError: pgx.ErrNoRows,
			expectedError:      storageerrors.ErrURLAliasNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			db := postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    mock,
			}

			sql, args, _ := db.Builder.
				Select("original").
				From(constant.URLSTable).
				Where(squirrel.Eq{"alias": tc.inputAlias}).
				ToSql()

			ctx := context.Background()

			in := input{
				sql:           sql,
				args:          args,
				rows:          tc.rows,
				expectedError: tc.expectedQueryError,
			}

			tc.mockBehaviour(mock, in)

			urlStorage := NewURLRepo(&db)

			original, err := urlStorage.GetOriginalByAlias(ctx, tc.inputAlias)
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedOriginal, original)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
}
