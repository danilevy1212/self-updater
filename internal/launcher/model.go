package launcher

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/danilevy1212/self-updater/internal/launcher/config"
	"github.com/danilevy1212/self-updater/internal/logger"
	"github.com/danilevy1212/self-updater/internal/models"
)

type Launcher struct {
	Meta             models.ApplicationMeta
	Logger           *zerolog.Logger
	Config           *config.Config
	SessionDirectory string
}

func New(ctx context.Context, am models.ApplicationMeta) (*Launcher, error) {
	conf, err := config.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	l := logger.New(conf.IsDev).
		With().
		Str("app", "launcher").
		Logger()

	sessionDir := filepath.Join(conf.SessionDirectory, uuid.NewString())

	return &Launcher{
		Meta:             am,
		Logger:           &l,
		Config:           conf,
		SessionDirectory: sessionDir,
	}, nil
}
