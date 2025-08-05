package server

import (
	"github.com/danilevy1212/self-updater/internal/logger"
	"github.com/gin-gonic/gin"
)

func (a *Application) RegisterGlobalMiddleware() {
	r := a.Router

	// GLOBAL
	// Recover from panics
	r.Use(gin.Recovery())

	// Zerolog logger
	l := logger.New(a.Config.IsDev).
		With().
		Str("app", "server").
		Logger()

	r.Use(logger.NewMiddleware(&l))
}
