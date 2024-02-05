package servicesroute

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

type HealthCheckRoutes struct{}

func NewHealthCheckRoutes(g *gin.RouterGroup) {
	r := &HealthCheckRoutes{}

	g.POST("/", r.CreateURLAlias)
	g.GET("/version", r.Version)
}

// Version
//
//	@Summary		Return a version of the application
//	@Description	Return a build date, release and git hash
//	@UUID			200
//	@Success		200		{object}	VersionResponse	"Version was successfully received"
//	@Router			services/version [get]
//	@Tags			services
func (r *HealthCheckRoutes) Version(ctx *gin.Context) {
	response := VersionResponse{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}
	ctx.JSON(http.StatusOK, response)
}

func (r *HealthCheckRoutes) Healthz(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func (r *HealthCheckRoutes) Readyz(ctx *gin.Context) {

}
