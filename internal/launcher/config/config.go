package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	IsDev            bool   `env:"LAUNCHER_IS_DEV,default=false"`
	SessionDirectory string `env:"LAUNCHER_SESSION_FOLDER,default=self-updater"`
}

type ConfigFunc func(context.Context) (*Config, error)

var ConfigFetcher ConfigFunc = fetchFromEnvironment

func New(ctx context.Context) (*Config, error) {
	return ConfigFetcher(ctx)
}

func fetchFromEnvironment(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	cfg.SessionDirectory = filepath.Join(os.TempDir(), cfg.SessionDirectory)

	return &cfg, nil
}
