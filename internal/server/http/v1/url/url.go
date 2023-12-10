package urlroute

import (
	"errors"
	"github.com/gin-gonic/gin"
	httpresponse "github.com/romandnk/shortener/internal/server/http/v1/response"
	"github.com/romandnk/shortener/internal/service"
	urlservice "github.com/romandnk/shortener/internal/service/url"
	"net/http"
)

type UrlRoutes struct {
	url service.URL
}

func NewUrlRoutes(g *gin.RouterGroup, url service.URL) {
	r := &UrlRoutes{
		url: url,
	}

	g.POST("/", r.CreateShortURL)
	g.GET("/:alias", r.GetOriginalByAlias)
}

// CreateShortURL
//
//	@Summary		Create short URL
//	@Description	Create new short URL from original.
//	@UUID			100
//	@Param			params	body		CreateShortURLRequest	true	"Required JSON body with original url"
//	@Success		201		{object}	CreateShortURLResponse	"Short URL was created successfully"
//	@Failure		400		{object}	httpresponse.Response	"Invalid input data"
//	@Failure		500		{object}	httpresponse.Response	"Internal error"
//	@Router			/urls [post]
//	@Tags			URL
func (r *UrlRoutes) CreateShortURL(ctx *gin.Context) {
	var params CreateShortURLRequest

	if err := ctx.BindJSON(&params); err != nil {
		httpresponse.SentErrorResponse(ctx, http.StatusBadRequest, "error binding json body", err)
		return
	}

	shortURL, err := r.url.CreateShortURL(ctx, params.OriginalURL)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, urlservice.ErrInternalError) {
			code = http.StatusInternalServerError
		}
		httpresponse.SentErrorResponse(ctx, code, "error creating short url", err)
		return
	}

	resp := CreateShortURLResponse{ShortURL: shortURL}

	ctx.JSON(http.StatusCreated, resp)
}

// GetOriginalByAlias
//
//	@Summary		Get original URL
//	@Description	Get original URL by its alias.
//	@UUID			101
//	@Param			alias	path		string						true	"Required path param with original url alias"
//	@Success		200		{object}	GetOriginalByAliasResponse	"Short URL was received successfully"
//	@Failure		400		{object}	httpresponse.Response		"Invalid input data"
//	@Failure		500		{object}	httpresponse.Response		"Internal error"
//	@Router			/urls/:alias [get]
//	@Tags			URL
func (r *UrlRoutes) GetOriginalByAlias(ctx *gin.Context) {
	alias := ctx.Param("alias")

	original, err := r.url.GetOriginalByAlias(ctx, alias)
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, urlservice.ErrInternalError) {
			code = http.StatusInternalServerError
		}
		httpresponse.SentErrorResponse(ctx, code, "error getting original url by alias", err)
		return
	}

	resp := GetOriginalByAliasResponse{OriginalURL: original}

	ctx.JSON(http.StatusOK, resp)
}
