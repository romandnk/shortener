package urlgrpc

import (
	"context"
	"errors"
	urlpb "github.com/romandnk/shortener/internal/server/grpc/url/pb"
	mock_service "github.com/romandnk/shortener/internal/service/mock"
	urlservice "github.com/romandnk/shortener/internal/service/url"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

func startGRPCServer() (*grpc.Server, *bufconn.Listener) {
	bufferSize := 1024 * 1024
	listener := bufconn.Listen(bufferSize)

	srv := grpc.NewServer()
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	return srv, listener
}

func getDialer(lis *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestHandlerGRPCCreateEvent(t *testing.T) {
	type args struct {
		input         string
		output        string
		expectedError error
	}

	type mockBehaviour func(m *mock_service.MockURL, args args)

	testCases := []struct {
		name          string
		input         urlpb.CreateURLAliasRequest
		args          args
		mock          mockBehaviour
		expectedAlias string
		expectedError error
	}{
		{
			name: "OK",
			input: urlpb.CreateURLAliasRequest{
				Original: "http://google.com",
			},
			args: args{
				input:  "http://google.com",
				output: "testtest11",
			},
			mock: func(m *mock_service.MockURL, args args) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			expectedAlias: "testtest11",
		},
		{
			name: "original url is empty",
			args: args{
				expectedError: urlservice.ErrEmptyOriginalURL,
			},
			mock: func(m *mock_service.MockURL, args args) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			expectedError: errors.New("rpc error: code = InvalidArgument desc = url cannot be empty"),
		},
		{
			name: "invalid original url format",
			args: args{
				expectedError: urlservice.ErrInvalidOriginalURL,
			},
			mock: func(m *mock_service.MockURL, args args) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			expectedError: errors.New("rpc error: code = InvalidArgument desc = invalid url format"),
		},
		{
			name: "original url exists",
			args: args{
				expectedError: storageerrors.ErrOriginalURLExists,
			},
			mock: func(m *mock_service.MockURL, args args) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			expectedError: errors.New("rpc error: code = InvalidArgument desc = original url already exists"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv, lis := startGRPCServer()
			defer srv.Stop()
			defer lis.Close()

			urlService := mock_service.NewMockURL(ctrl)
			handler := URLHandler{
				url: urlService,
			}

			urlpb.RegisterEventServiceServer(srv, handler)

			ctx := context.Background()

			conn, err := grpc.DialContext(ctx, "",
				grpc.WithContextDialer(getDialer(lis)),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()

			client := urlpb.NewEventServiceClient(conn)

			tc.mock(urlService, tc.args)

			res, err := client.CreateURLAlias(ctx, &tc.input)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedAlias, res.GetAlias())
		})
	}
}

func TestURLHandler_GetOriginalByAlias(t *testing.T) {
	type args struct {
		input         string
		output        string
		expectedError error
	}

	type mockBehaviour func(m *mock_service.MockURL, args args)

	testCases := []struct {
		name             string
		input            urlpb.GetOriginalByAliasRequest
		args             args
		mock             mockBehaviour
		expectedOriginal string
		expectedError    error
	}{
		{
			name: "OK",
			input: urlpb.GetOriginalByAliasRequest{
				Alias: "testtest11",
			},
			args: args{
				input:  "testtest11",
				output: "http://google.com",
			},
			mock: func(m *mock_service.MockURL, args args) {
				m.EXPECT().GetOriginalByAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			expectedOriginal: "http://google.com",
		},
		{
			name: "too short alias",
			input: urlpb.GetOriginalByAliasRequest{
				Alias: "testtest",
			},
			args: args{
				input:         "testtest",
				expectedError: urlservice.ErrInvalidAliasFormat,
			},
			mock: func(m *mock_service.MockURL, args args) {
				m.EXPECT().GetOriginalByAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			expectedError: errors.New("rpc error: code = InvalidArgument desc = unique id must be 10 symbols length"),
		},
		{
			name: "original url is not found",
			input: urlpb.GetOriginalByAliasRequest{
				Alias: "testtest12",
			},
			args: args{
				input:         "testtest12",
				expectedError: urlservice.ErrOriginalURLNotFound,
			},
			mock: func(m *mock_service.MockURL, args args) {
				m.EXPECT().GetOriginalByAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			expectedError: errors.New("rpc error: code = InvalidArgument desc = original url is not found"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srv, lis := startGRPCServer()
			defer srv.Stop()
			defer lis.Close()

			urlService := mock_service.NewMockURL(ctrl)
			handler := URLHandler{
				url: urlService,
			}

			urlpb.RegisterEventServiceServer(srv, handler)

			ctx := context.Background()

			conn, err := grpc.DialContext(ctx, "",
				grpc.WithContextDialer(getDialer(lis)),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)
			defer conn.Close()

			client := urlpb.NewEventServiceClient(conn)

			tc.mock(urlService, tc.args)

			res, err := client.GetOriginalByAlias(ctx, &tc.input)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedOriginal, res.GetOriginal())
		})
	}
}
