package v1

import (
	"github.com/gin-gonic/gin"
	docs "github.com/romandnk/shortener/docs"
	"github.com/romandnk/shortener/internal/server/http/middleware"
	urlroute "github.com/romandnk/shortener/internal/server/http/v1/url"
	"github.com/romandnk/shortener/internal/service"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	engine   *gin.Engine
	services *service.Services
	mw       *middleware.MW
}

func NewHandler(services *service.Services, mw *middleware.MW) *Handler {
	return &Handler{
		services: services,
		mw:       mw,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	h.engine = router

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("/api/v1", h.mw.Logging())
	{
		// urls management group
		urls := api.Group("/urls")
		{
			urlroute.NewUrlRoutes(urls, h.services.URL)
		}
	}

	return router
}
