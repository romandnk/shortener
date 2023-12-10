package urlservice

import (
	"context"
	"github.com/romandnk/shortener/internal/entity"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	mock_storage "github.com/romandnk/shortener/internal/storage/mock"
	mock_generate "github.com/romandnk/shortener/pkg/generator/mock"
	mock_logger "github.com/romandnk/shortener/pkg/logger/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"net/url"
	"testing"
)

func TestURLService_CreateShortURL(t *testing.T) {
	u, err := url.ParseRequestURI("http://localhost:9090/api/v1/")
	require.NoError(t, err)

	type loggerArgs struct {
		msg  string
		args []any
	}

	type urlArgs struct {
		ctx   context.Context
		url   entity.URL
		error error
	}

	type generatorArgs struct {
		expectedRandomString string
		error                error
	}

	type loggerBehaviour func(m *mock_logger.MockLogger, args loggerArgs)
	type generatorBehaviour func(m *mock_generate.MockGenerator, args generatorArgs)
	type repoBehaviour func(m *mock_storage.MockURL, args urlArgs)

	testCases := []struct {
		name               string
		inputOriginal      string
		loggerArgs         loggerArgs
		loggerMock         loggerBehaviour
		urlArgs            urlArgs
		generatorArgs      generatorArgs
		generatorBehaviour generatorBehaviour
		urlMock            repoBehaviour
		expectedShort      string
		expectedError      error
	}{
		{
			name:          "OK",
			inputOriginal: "http://google.com/",
			loggerArgs: loggerArgs{
				msg:  "URLService.CreateShortURL - alias was created successfully",
				args: []any{zap.String("alias", "abcdefghig")},
			},
			loggerMock: func(m *mock_logger.MockLogger, args loggerArgs) {
				m.EXPECT().Info(args.msg, args.args)
			},
			urlArgs: urlArgs{
				ctx: context.Background(),
				url: entity.URL{
					Original: "http://google.com/",
					Alias:    "abcdefghig",
				},
			},
			generatorArgs: generatorArgs{
				expectedRandomString: "abcdefghig",
			},
			generatorBehaviour: func(m *mock_generate.MockGenerator, args generatorArgs) {
				m.EXPECT().Random().Return(args.expectedRandomString, args.error)
			},
			urlMock: func(m *mock_storage.MockURL, args urlArgs) {
				m.EXPECT().CreateURL(args.ctx, args.url).Return(args.error)
			},
			expectedShort: "http://localhost:9090/api/v1/abcdefghig",
		},
		{
			name: "empty original url",
			loggerArgs: loggerArgs{
				msg:  "URLService.CreateShortURL",
				args: []any{zap.String("error", ErrEmptyOriginalURL.Error())},
			},
			loggerMock: func(m *mock_logger.MockLogger, args loggerArgs) {
				m.EXPECT().Error(args.msg, args.args)
			},
			expectedShort: "",
			expectedError: ErrEmptyOriginalURL,
		},
		{
			name:          "original url too long",
			inputOriginal: string(make([]byte, 3000)),
			loggerArgs: loggerArgs{
				msg:  "URLService.CreateShortURL",
				args: []any{zap.String("error", ErrOriginalURLTooLong.Error())},
			},
			loggerMock: func(m *mock_logger.MockLogger, args loggerArgs) {
				m.EXPECT().Error(args.msg, args.args)
			},
			expectedShort: "",
			expectedError: ErrOriginalURLTooLong,
		},
		{
			name:          "original url already exists",
			inputOriginal: "http://google.com/",
			loggerArgs: loggerArgs{
				msg: "URLService.CreateShortURL",
				args: []any{
					zap.String("original", "http://google.com/"),
					zap.String("error", storageerrors.ErrOriginalURLExists.Error()),
				},
			},
			loggerMock: func(m *mock_logger.MockLogger, args loggerArgs) {
				m.EXPECT().Error(args.msg, args.args)
			},
			generatorArgs: generatorArgs{
				expectedRandomString: "abcdefghig",
			},
			generatorBehaviour: func(m *mock_generate.MockGenerator, args generatorArgs) {
				m.EXPECT().Random().Return(args.expectedRandomString, args.error)
			},
			urlArgs: urlArgs{
				ctx: context.Background(),
				url: entity.URL{
					Original: "http://google.com/",
					Alias:    "abcdefghig",
				},
				error: storageerrors.ErrOriginalURLExists,
			},
			urlMock: func(m *mock_storage.MockURL, args urlArgs) {
				m.EXPECT().CreateURL(args.ctx, args.url).Return(args.error)
			},
			expectedShort: "",
			expectedError: storageerrors.ErrOriginalURLExists,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			urlStorage := mock_storage.NewMockURL(ctrl)
			generator := mock_generate.NewMockGenerator(ctrl)
			log := mock_logger.NewMockLogger(ctrl)

			urlService := URLService{
				generator: generator,
				baseURL:   *u,
				url:       urlStorage,
				logger:    log,
			}

			if tc.loggerMock != nil {
				tc.loggerMock(log, tc.loggerArgs)
			}
			if tc.generatorBehaviour != nil {
				tc.generatorBehaviour(generator, tc.generatorArgs)
			}
			if tc.urlMock != nil {
				tc.urlMock(urlStorage, tc.urlArgs)
			}

			output, err := urlService.CreateShortURL(ctx, tc.inputOriginal)
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedShort, output)
		})
	}
}

func TestURLService_GetOriginalByAlias(t *testing.T) {
	type loggerArgs struct {
		msg  string
		args []any
	}

	type urlArgs struct {
		ctx      context.Context
		alias    string
		original string
		error    error
	}

	type loggerBehaviour func(m *mock_logger.MockLogger, args loggerArgs)
	type repoBehaviour func(m *mock_storage.MockURL, args urlArgs)

	testCases := []struct {
		name             string
		inputAlias       string
		loggerArgs       loggerArgs
		loggerMock       loggerBehaviour
		urlArgs          urlArgs
		urlMock          repoBehaviour
		expectedOriginal string
		expectedError    error
	}{
		{
			name:       "OK",
			inputAlias: "abcdefghig",
			loggerArgs: loggerArgs{
				msg:  "URLService.GetOriginalByAlias - alias was received successfully",
				args: []any{zap.String("alias", "abcdefghig")},
			},
			loggerMock: func(m *mock_logger.MockLogger, args loggerArgs) {
				m.EXPECT().Info(args.msg, args.args)
			},
			urlArgs: urlArgs{
				ctx:      context.Background(),
				alias:    "abcdefghig",
				original: "http://google.com/",
			},
			urlMock: func(m *mock_storage.MockURL, args urlArgs) {
				m.EXPECT().GetOriginalByAlias(args.ctx, args.alias).Return(args.original, args.error)
			},
			expectedOriginal: "http://google.com/",
		},
		{
			name: "empty alias",
			loggerArgs: loggerArgs{
				msg:  "URLService.GetOriginalByAlias",
				args: []any{zap.String("error", ErrEmptyURLAlias.Error())},
			},
			loggerMock: func(m *mock_logger.MockLogger, args loggerArgs) {
				m.EXPECT().Error(args.msg, args.args)
			},
			expectedOriginal: "",
			expectedError:    ErrEmptyURLAlias,
		},
		{
			name:       "original url is not found",
			inputAlias: "abcdefghig",
			loggerArgs: loggerArgs{
				msg: "URLService.GetOriginalByAlias",
				args: []any{
					zap.String("alias", "abcdefghig"),
					zap.String("error", storageerrors.ErrURLAliasNotFound.Error()),
				},
			},
			loggerMock: func(m *mock_logger.MockLogger, args loggerArgs) {
				m.EXPECT().Error(args.msg, args.args)
			},
			urlArgs: urlArgs{
				ctx:   context.Background(),
				alias: "abcdefghig",
				error: storageerrors.ErrURLAliasNotFound,
			},
			urlMock: func(m *mock_storage.MockURL, args urlArgs) {
				m.EXPECT().GetOriginalByAlias(args.ctx, args.alias).Return(args.original, args.error)
			},
			expectedError: ErrOriginalURLNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			urlStorage := mock_storage.NewMockURL(ctrl)
			log := mock_logger.NewMockLogger(ctrl)

			urlService := URLService{
				url:    urlStorage,
				logger: log,
			}

			if tc.loggerMock != nil {
				tc.loggerMock(log, tc.loggerArgs)
			}
			if tc.urlMock != nil {
				tc.urlMock(urlStorage, tc.urlArgs)
			}

			original, err := urlService.GetOriginalByAlias(ctx, tc.inputAlias)
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedOriginal, original)
		})
	}
}
