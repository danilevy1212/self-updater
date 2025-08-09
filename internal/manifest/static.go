package manifest

import (
	"context"
	"encoding/json"

	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/models/fixtures"
)

var StaticManifestRaw = fixtures.ReleaseFixture

type StaticFetcher struct{}

func (StaticFetcher) FetchManifest(ctx context.Context) (*models.ReleaseManifest, error) {
	var manifest models.ReleaseManifest
	if err := json.Unmarshal(StaticManifestRaw, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func NewStaticFetcher() ManifestFetcher {
	return &StaticFetcher{}
}
