package urlroute

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	mock_service "github.com/romandnk/shortener/internal/service/mock"
	urlservice "github.com/romandnk/shortener/internal/service/url"
	storageerrors "github.com/romandnk/shortener/internal/storage/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlRoutes_CreateURLAlias(t *testing.T) {
	url := "/api/v1/urls"

	type argsUrl struct {
		input         string
		output        string
		expectedError error
	}

	type mockUrlBehaviour func(m *mock_service.MockURL, args argsUrl)

	testCases := []struct {
		name                 string
		argsUrl              argsUrl
		urlM                 mockUrlBehaviour
		requestBody          map[string]interface{}
		expectedResponseBody string
		expectedHTTPCode     int
	}{
		{
			name: "OK",
			argsUrl: argsUrl{
				input:  "https://google.com",
				output: "testtest12",
			},
			urlM: func(m *mock_service.MockURL, args argsUrl) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			requestBody: map[string]interface{}{
				"original_url": "https://google.com",
			},
			expectedResponseBody: `{"alias":"testtest12"}`,
			expectedHTTPCode:     http.StatusCreated,
		},
		{
			name: "original url is empty",
			argsUrl: argsUrl{
				expectedError: urlservice.ErrEmptyOriginalURL,
			},
			urlM: func(m *mock_service.MockURL, args argsUrl) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			requestBody: map[string]interface{}{
				"original_url": "",
			},
			expectedResponseBody: `{"message":"error creating short url","error":"url cannot be empty"}`,
			expectedHTTPCode:     http.StatusBadRequest,
		},
		{
			name: "invalid original url format",
			argsUrl: argsUrl{
				input:         "http//google.com",
				expectedError: urlservice.ErrInvalidOriginalURL,
			},
			urlM: func(m *mock_service.MockURL, args argsUrl) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			requestBody: map[string]interface{}{
				"original_url": "http//google.com",
			},
			expectedResponseBody: `{"message":"error creating short url","error":"invalid url format"}`,
			expectedHTTPCode:     http.StatusBadRequest,
		},
		{
			name: "invalid original url format",
			argsUrl: argsUrl{
				input:         "http://google.com",
				expectedError: storageerrors.ErrOriginalURLExists,
			},
			urlM: func(m *mock_service.MockURL, args argsUrl) {
				m.EXPECT().CreateURLAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			requestBody: map[string]interface{}{
				"original_url": "http://google.com",
			},
			expectedResponseBody: `{"message":"error creating short url","error":"original url already exists"}`,
			expectedHTTPCode:     http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlService := mock_service.NewMockURL(ctrl)

			if tc.urlM != nil {
				tc.urlM(urlService, tc.argsUrl)
			}

			urlR := UrlRoutes{
				url: urlService,
			}

			r := gin.Default()
			r.POST(url, urlR.CreateURLAlias)

			jsonBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedHTTPCode, w.Code)

			require.Equal(t, []byte(tc.expectedResponseBody), w.Body.Bytes())
		})
	}
}

func TestUrlRoutes_GetOriginalByAlias(t *testing.T) {
	url := "/api/v1/urls/"

	type argsAlias struct {
		input         string
		output        string
		expectedError error
	}

	type mockUrlBehaviour func(m *mock_service.MockURL, args argsAlias)

	testCases := []struct {
		name                 string
		argsUrl              argsAlias
		urlM                 mockUrlBehaviour
		pathParam            string
		expectedResponseBody string
		expectedHTTPCode     int
	}{
		{
			name: "OK",
			argsUrl: argsAlias{
				input:  "testtest12",
				output: "https://google.com",
			},
			urlM: func(m *mock_service.MockURL, args argsAlias) {
				m.EXPECT().GetOriginalByAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			pathParam:            "testtest12",
			expectedResponseBody: `{"original_url":"https://google.com"}`,
			expectedHTTPCode:     http.StatusOK,
		},
		{
			name: "too short alias",
			argsUrl: argsAlias{
				input:         "testtest",
				expectedError: urlservice.ErrInvalidAliasFormat,
			},
			urlM: func(m *mock_service.MockURL, args argsAlias) {
				m.EXPECT().GetOriginalByAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			pathParam:            "testtest",
			expectedResponseBody: `{"message":"error getting original url by alias","error":"unique id must be 10 symbols length"}`,
			expectedHTTPCode:     http.StatusBadRequest,
		},
		{
			name: "too short alias",
			argsUrl: argsAlias{
				input:         "testtest12",
				expectedError: urlservice.ErrOriginalURLNotFound,
			},
			urlM: func(m *mock_service.MockURL, args argsAlias) {
				m.EXPECT().GetOriginalByAlias(gomock.Any(), args.input).Return(args.output, args.expectedError)
			},
			pathParam:            "testtest12",
			expectedResponseBody: `{"message":"error getting original url by alias","error":"original url is not found"}`,
			expectedHTTPCode:     http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			urlService := mock_service.NewMockURL(ctrl)

			if tc.urlM != nil {
				tc.urlM(urlService, tc.argsUrl)
			}

			urlR := UrlRoutes{
				url: urlService,
			}

			r := gin.Default()
			r.GET(url+":alias", urlR.GetOriginalByAlias)

			w := httptest.NewRecorder()

			ctx := gin.CreateTestContextOnly(w, r)
			ctx.AddParam("alias", tc.pathParam)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+tc.pathParam, nil)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedHTTPCode, w.Code)

			require.Equal(t, []byte(tc.expectedResponseBody), w.Body.Bytes())
		})
	}
}
