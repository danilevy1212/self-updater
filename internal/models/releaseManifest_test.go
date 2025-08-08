package models

import (
	"encoding/json"
	"github.com/danilevy1212/self-updater/internal/models/fixtures"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonParse(t *testing.T) {
	var manifest ReleaseManifest
	err := json.Unmarshal(fixtures.ReleaseFixture, &manifest)
	assert.NoError(t, err, "should parse release fixture without error")

	assert.Equal(t, "v1.2.3", manifest.Latest)
	assert.Equal(t, 3, len(manifest.Versions), "should contain 3 versions")

	latest := manifest.Versions[0]
	assert.Equal(t, "v1.2.3", latest.Version)
	assert.Equal(t, 3, len(latest.Artifacts), "latest version should have 3 artifacts")

	for _, artifact := range latest.Artifacts {
		assert.NotEmpty(t, artifact.OS, "should have an OS")
		assert.NotEmpty(t, artifact.Arch, "should have the hardware architecture")
		assert.NotEmpty(t, artifact.Filename, "should have the filename")
		assert.NotEmpty(t, artifact.Digest, "should contain the digest")
		assert.NotEmpty(t, artifact.SignatureBase64, "should contain the base64 signature of the binary")
		assert.NotEmpty(t, artifact.URL, "should have a URL for the artifact")
	}
}
