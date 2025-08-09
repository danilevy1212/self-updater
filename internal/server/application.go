package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/server/config"
)

type Application struct {
	Server *http.Server
	Router *gin.Engine
	Config *config.Config
	Meta   models.ApplicationMeta
}

func (a *Application) Serve(port uint) error {
	a.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.Port),
		Handler: a.Router,
	}

	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	if a.Server == nil {
		return nil
	}

	return a.Server.Shutdown(ctx)
}

func New(ctx context.Context, meta models.ApplicationMeta) (*Application, error) {
	c, err := config.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if !c.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}
	r.RemoveExtraSlash = true

	return &Application{
		Config: c,
		Meta:   meta,
		Router: r,
	}, nil
}
