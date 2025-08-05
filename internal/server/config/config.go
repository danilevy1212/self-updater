package config

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	IsDev bool `env:"SERVER_IS_DEV,default=false"`
	Port  uint `env:"SERVER_PORT,default=3000"`
}

type ConfigFunc func(context.Context) (*Config, error)

var ConfigFetcher ConfigFunc = fetchFromEnvironment

func New(ctx context.Context) (*Config, error) {
	return ConfigFetcher(ctx)
}

// TODO  Unit test this shit
func fetchFromEnvironment(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}
	return &cfg, nil
}
