package manifest

import (
	"context"

	"github.com/danilevy1212/self-updater/internal/models"
)

type ManifestFetcher interface {
	FetchManifest(ctx context.Context) (*models.ReleaseManifest, error)
}
