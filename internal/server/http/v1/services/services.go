package servicesroute

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

type HealthCheckRoutes struct {
	ok *atomic.Bool
}

func NewHealthCheckRoutes(g *gin.RouterGroup, ok *atomic.Bool) {
	r := &HealthCheckRoutes{
		ok: ok,
	}

	g.GET("/liveness", r.LivenessProbe)
	g.GET("/readiness", r.ReadinessProbe)
	g.GET("/version", r.Version)
}

// Version
//
//	@Summary		Return a version of the application
//	@Description	Return a build date, release and git hash
//	@UUID			200
//	@Success		200	{object}	VersionResponse	"Version was successfully received"
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

// LivenessProbe
//
//	@Summary		K8s checks liveness probe
//	@Description	Return a status okay
//	@UUID			200
//	@Success		200	nil	"Liveness probe is got over"
//	@Router			services/liveness [get]
//	@Tags			services
func (r *HealthCheckRoutes) LivenessProbe(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

// ReadinessProbe
//
//	@Summary		K8s checks readiness probe
//	@Description	Return a status okay if everything is initialized
//	@UUID			200
//	@Success		200	nil	"Readiness probe is got over"
//	@Router			services/readiness [get]
//	@Tags			services
func (r *HealthCheckRoutes) ReadinessProbe(ctx *gin.Context) {
	if r.ok == nil || !r.ok.Load() {
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}

	ctx.AbortWithStatus(http.StatusOK)
}
