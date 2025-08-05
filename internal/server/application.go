package server

import (
	"context"
	"fmt"
	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/server/config"
	"github.com/gin-gonic/gin"
)

type Application struct {
	Router *gin.Engine
	Config *config.Config
	Meta   models.ApplicationMeta
}

func (a *Application) Serve(port uint) error {
	return a.Router.Run(fmt.Sprintf(":%d", a.Config.Port))
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

	err = r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	r.RemoveExtraSlash = true

	return &Application{
		Config: c,
		Meta:   meta,
		Router: r,
	}, nil
}
