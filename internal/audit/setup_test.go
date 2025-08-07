package audit

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/models/fixtures"
)

var releaseFixture models.ReleaseManifest

func TestMain(m *testing.M) {
	error := json.Unmarshal(fixtures.ReleaseFixture, &releaseFixture)
	if error != nil {
		panic("failed to parse release fixture: " + error.Error())
	}

	v := m.Run()

	os.Exit(v)
}
