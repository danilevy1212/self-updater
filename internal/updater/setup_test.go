package updater

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"

	"github.com/danilevy1212/self-updater/internal/manifest"
	"github.com/danilevy1212/self-updater/internal/models"
)

func TestMain(m *testing.M) {
	ManifestFetcherFactory = func(ctx context.Context, applicationMeta models.ApplicationMeta, logger *zerolog.Logger) (manifest.ManifestFetcher, error) {
		return manifest.NewStaticFetcher(), nil
	}

	v := m.Run()
	os.Exit(v)
}
