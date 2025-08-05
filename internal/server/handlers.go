package server

import (
	"net/http"

	"github.com/danilevy1212/self-updater/internal/logger"
	"github.com/gin-gonic/gin"
)

func (a *Application) HealthCheck(ctx *gin.Context) {
	log := logger.FromContext(ctx.Request.Context()).
		With().
		Str("handler", "HealthCheck").
		Logger()

	log.Info().
		Msg("Health check endpoint hit")

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"sha256":  a.Meta.Digest,
		"version": a.Meta.Version,
		"commit":  a.Meta.Commit,
	})
}
