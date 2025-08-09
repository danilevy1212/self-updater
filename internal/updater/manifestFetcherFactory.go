package updater

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/danilevy1212/self-updater/internal/manifest"
	"github.com/danilevy1212/self-updater/internal/models"
)

type ManifestFetcherFactoryFunc func(ctx context.Context, applicationMeta models.ApplicationMeta, logger *zerolog.Logger) (manifest.ManifestFetcher, error)

func createGithubManifestFetcher(ctx context.Context, applicationMeta models.ApplicationMeta, logger *zerolog.Logger) (manifest.ManifestFetcher, error) {
	factoryCtx := manifest.SetGithubManifestFetcherLogger(ctx, logger)
	ghmf, err := manifest.NewGithubManifestFetcher(factoryCtx, applicationMeta)
	if err != nil {
		return nil, err
	}
	return ghmf, nil
}

var ManifestFetcherFactory ManifestFetcherFactoryFunc = createGithubManifestFetcher
