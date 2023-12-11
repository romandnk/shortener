package urlgrpc

import (
	"context"
	"errors"
	urlpb "github.com/romandnk/shortener/internal/server/grpc/url/pb"
	"github.com/romandnk/shortener/internal/service"
	urlservice "github.com/romandnk/shortener/internal/service/url"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type URLHandler struct {
	url service.URL
	urlpb.UnimplementedEventServiceServer
}

func NewURLHandler(url service.URL) *URLHandler {
	return &URLHandler{url: url}
}

func (h *URLHandler) CreateURLAlias(ctx context.Context, req *urlpb.CreateURLAliasRequest) (*urlpb.CreateURLAliasResponse, error) {
	alias, err := h.url.CreateURLAlias(ctx, req.GetOriginal())
	if err != nil {
		code := codes.InvalidArgument
		if errors.Is(err, urlservice.ErrInternalError) {
			code = codes.Internal
		}
		return nil, status.Error(code, err.Error())
	}
	return &urlpb.CreateURLAliasResponse{
		Alias: alias,
	}, nil
}

func (h *URLHandler) GetOriginalByAlias(ctx context.Context, req *urlpb.GetOriginalByAliasRequest) (*urlpb.GetOriginalByAliasResponse, error) {
	original, err := h.url.GetOriginalByAlias(ctx, req.GetAlias())
	if err != nil {
		code := codes.InvalidArgument
		if errors.Is(err, urlservice.ErrInternalError) {
			code = codes.Internal
		}
		return nil, status.Error(code, err.Error())
	}
	return &urlpb.GetOriginalByAliasResponse{
		Original: original,
	}, nil
}
