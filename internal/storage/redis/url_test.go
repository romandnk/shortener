package redisstorage

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/romandnk/shortener/internal/constant"
	"github.com/romandnk/shortener/internal/entity"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestURLRepo_CreateURL(t *testing.T) {
	type input struct {
		keyOne           string
		valueOne         string
		expectedErrorOne error
		keyTwo           string
		valueTwo         string
		expectedErrorTwo error
	}

	type mockBehaviour func(m redismock.ClientMock, input input)

	testCases := []struct {
		name          string
		url           entity.URL
		input         input
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name: "OK",
			url: entity.URL{
				Original: "http://test.com",
				Alias:    "testtest11",
			},
			input: input{
				keyOne:   "http://test.com",
				valueOne: "testtest11",
				keyTwo:   "testtest11",
				valueTwo: "http://test.com",
			},
			mockBehaviour: func(m redismock.ClientMock, input input) {
				m.ExpectSetNX(input.keyOne, input.valueOne, constant.ZeroTTL).SetVal(true)
				m.ExpectSetNX(input.keyTwo, input.valueTwo, constant.ZeroTTL).SetVal(true)
			},
		},
		{
			name: "original url already exists",
			url: entity.URL{
				Original: "http://test.com",
				Alias:    "testtest11",
			},
			input: input{
				keyOne:   "http://test.com",
				valueOne: "testtest11",
			},
			mockBehaviour: func(m redismock.ClientMock, input input) {
				m.ExpectSetNX(input.keyOne, input.valueOne, constant.ZeroTTL).SetVal(false)
			},
			expectedError: storageerrors.ErrOriginalURLExists,
		},
		{
			name: "url alias already exists",
			url: entity.URL{
				Original: "http://test.com",
				Alias:    "testtest11",
			},
			input: input{
				keyOne:   "http://test.com",
				valueOne: "testtest11",
				keyTwo:   "testtest11",
				valueTwo: "http://test.com",
			},
			mockBehaviour: func(m redismock.ClientMock, input input) {
				m.ExpectSetNX(input.keyOne, input.valueOne, constant.ZeroTTL).SetVal(true)
				m.ExpectSetNX(input.keyTwo, input.valueTwo, constant.ZeroTTL).SetVal(false)
			},
			expectedError: storageerrors.ErrURLAliasExists,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			defer db.Close()

			ctx := context.Background()

			tc.mockBehaviour(mock, tc.input)

			urlStorage := NewURLRepo(db)

			err := urlStorage.CreateURL(ctx, tc.url)
			require.ErrorIs(t, err, tc.expectedError)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
}

func TestURLRepo_GetOriginalByAlias(t *testing.T) {
	type input struct {
		key           string
		value         string
		expectedError error
	}

	type mockBehaviour func(m redismock.ClientMock, input input)

	testCases := []struct {
		name             string
		inputAlias       string
		input            input
		mockBehaviour    mockBehaviour
		expectedOriginal string
		expectedError    error
	}{
		{
			name:       "OK",
			inputAlias: "testtest11",
			input: input{
				key:   "testtest11",
				value: "http://test.com",
			},
			mockBehaviour: func(m redismock.ClientMock, input input) {
				m.ExpectGet(input.key).SetVal(input.value)
			},
			expectedOriginal: "http://test.com",
		},
		{
			name:       "OK",
			inputAlias: "testtest11",
			input: input{
				key:           "testtest11",
				expectedError: redis.Nil,
			},
			mockBehaviour: func(m redismock.ClientMock, input input) {
				m.ExpectGet(input.key).SetErr(input.expectedError)
			},
			expectedError: storageerrors.ErrURLAliasNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			defer db.Close()

			ctx := context.Background()

			tc.mockBehaviour(mock, tc.input)

			urlStorage := NewURLRepo(db)

			original, err := urlStorage.GetOriginalByAlias(ctx, tc.inputAlias)
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedOriginal, original)

			require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
		})
	}
}
