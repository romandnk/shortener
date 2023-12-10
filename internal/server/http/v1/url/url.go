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
}

// CreateShortURL
//
//	@Summary		Create short URL
//	@Description	Create new short URL from original.
//	@UUID			100
//	@Param			params	body		CreateShortURLRequest	true	"Required JSON body with original url"
//	@Success		200		{object}	CreateShortURLResponse	"Short URL was created successfully"
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
