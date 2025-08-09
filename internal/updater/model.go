package updater

import (
	"context"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"

	"github.com/danilevy1212/self-updater/internal/logger"
	"github.com/danilevy1212/self-updater/internal/manifest"
	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/updater/config"
)

type OnUpgradeReadyFunc func(newVersion *os.File, updater *Updater)

type JobID int

type Updater struct {
	Meta            models.ApplicationMeta
	Cron            *cron.Cron
	Config          *config.Config
	Logger          *zerolog.Logger
	ManifestFetcher manifest.ManifestFetcher
	OnUpgradeReady  OnUpgradeReadyFunc
}

func (u *Updater) Start() (JobID, error) {
	u.Logger.Info().
		Msg("Starting updater job")

	id, err := u.Cron.AddJob(u.Config.Schedule, u)
	if err != nil {
		return 0, err
	}

	u.Logger.Info().
		Int("job_id", int(id)).
		Str("schedule", u.Config.Schedule).
		Msg("Updater job scheduled")

	u.Cron.Start()

	u.Logger.Info().
		Msg("Updater job started")

	return JobID(id), nil
}

func New(ctx context.Context, am models.ApplicationMeta, onUpgradeReadyCallback OnUpgradeReadyFunc) (*Updater, error) {
	conf, err := config.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cr := cron.New()

	l := logger.New(conf.IsDev).
		With().
		Str("app", "updater").
		Logger()
	mfl := l.With().Str("service", "manifest_fetcher").Logger()

	mf, err := ManifestFetcherFactory(ctx, am, &mfl)
	if err != nil {
		return nil, fmt.Errorf("failed to create manifest fetcher: %w", err)
	}

	return &Updater{
		Meta:            am,
		Config:          conf,
		Cron:            cr,
		Logger:          &l,
		ManifestFetcher: mf,
		OnUpgradeReady:  onUpgradeReadyCallback,
	}, nil
}
