package updater

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/danilevy1212/self-updater/internal/assets"
	"github.com/danilevy1212/self-updater/internal/audit"
	"github.com/danilevy1212/self-updater/internal/digest"
	"github.com/danilevy1212/self-updater/internal/downloader"
	"github.com/danilevy1212/self-updater/internal/manifest"
	"github.com/danilevy1212/self-updater/internal/models"
)

type ErrorFetcher struct{}

func (ef *ErrorFetcher) FetchManifest(ctx context.Context) (*models.ReleaseManifest, error) {
	return nil, errors.New("failed to fetch manifest")
}

func Test_Updater_Run(t *testing.T) {
	t.Run("should return if application public key doesn't match manifests", func(t *testing.T) {
		up, _ := New(context.Background(), models.ApplicationMeta{AuthorsPublicKey: []byte(`wrong`)}, func(newVersion *os.File) {
			assert.Fail(t, "should not call callback when public key doesn't match")
		})

		var buf bytes.Buffer
		logger := zerolog.New(&buf)
		up.Logger = &logger

		up.Run()

		assert.Contains(t, buf.String(), "Manifest public key does not match application public key")
	})

	t.Run("should return if application version matchest latest from manifest", func(t *testing.T) {
		up, _ := New(context.Background(), models.ApplicationMeta{
			AuthorsPublicKey: assets.PublicKeyPEM,
			Version:          "v1.2.3",
		}, func(newVersion *os.File) {
			assert.Fail(t, "should not call callback when version matches latest from manifest")
		})

		var buf bytes.Buffer
		logger := zerolog.New(&buf)
		up.Logger = &logger

		up.Run()
		assert.Contains(t, buf.String(), "No updates available. Current version is up to date.")
	})

	t.Run("should return if manifest fails to download", func(t *testing.T) {
		old := ManifestFetcherFactory
		defer func() {
			ManifestFetcherFactory = old
		}()
		ManifestFetcherFactory = func(ctx context.Context, applicationMeta models.ApplicationMeta, logger *zerolog.Logger) (manifest.ManifestFetcher, error) {

			return &ErrorFetcher{}, nil
		}

		up, _ := New(context.Background(), models.ApplicationMeta{
			AuthorsPublicKey: assets.PublicKeyPEM,
			Version:          "v1.2.2",
		}, func(newVersion *os.File) {
			assert.Fail(t, "should not call callback when manifest fetch fails")
		})

		var buf bytes.Buffer
		logger := zerolog.New(&buf)
		assert.NotNil(t, logger, "Logger should not be nil")
		up.Logger = &logger

		up.Run()
		assert.Contains(t, buf.String(), "Failed to fetch manifest")
	})

	t.Run("should return if artifact fetch fails", func(t *testing.T) {
		old := downloader.DownloadToTemporaryFile
		defer func() {
			downloader.DownloadToTemporaryFile = old
		}()
		downloader.DownloadToTemporaryFile = func(ctx context.Context, url, _ string) (*os.File, error) {
			return nil, errors.New("failed to fetch manifest")
		}

		up, _ := New(context.Background(), models.ApplicationMeta{
			AuthorsPublicKey: assets.PublicKeyPEM,
			Version:          "v1.2.2",
		}, func(newVersion *os.File) {
			assert.Fail(t, "should not call callback when manifest fetch fails")
		})

		var buf bytes.Buffer
		logger := zerolog.New(&buf)
		up.Logger = &logger

		up.Run()
		assert.Contains(t, buf.String(), "Failed to get artifact for platform from manifest")
	})

	t.Run("should return if artifact digest does not match", func(t *testing.T) {
		oldDownload := downloader.DownloadToTemporaryFile
		var fileName string
		defer func() {
			downloader.DownloadToTemporaryFile = oldDownload
		}()
		downloader.DownloadToTemporaryFile = func(ctx context.Context, url, _ string) (*os.File, error) {
			file, _ := os.CreateTemp("", "artifact")
			fileName = file.Name()
			return file, nil
		}

		up, _ := New(context.Background(), models.ApplicationMeta{
			AuthorsPublicKey: assets.PublicKeyPEM,
			Version:          "v1.2.2",
			OS:               "linux",
			Arch:             "amd64",
		}, func(newVersion *os.File) {
			assert.Fail(t, "should not call callback when artifact digest does not match")
		})

		var buf bytes.Buffer
		logger := zerolog.New(&buf)
		up.Logger = &logger

		up.Run()
		assert.Contains(t, buf.String(), "Artifact file digest does not match expected digest")
		assert.NotEmpty(t, fileName, "Temporary artifact file should have been created")
		assert.NoFileExists(t, fileName, "Temporary artifact file should not exist")
	})

	t.Run("should return if artifact signature verification fails", func(t *testing.T) {
		oldDownload := downloader.DownloadToTemporaryFile
		oldVerify := audit.VerifySignature
		oldDigest := digest.DigestFile
		var fileName string
		defer func() {
			downloader.DownloadToTemporaryFile = oldDownload
			audit.VerifySignature = oldVerify
			digest.DigestFile = oldDigest
		}()
		downloader.DownloadToTemporaryFile = func(ctx context.Context, url, _ string) (*os.File, error) {
			file, _ := os.CreateTemp("", "artifact")
			fileName = file.Name()
			return file, nil
		}
		digest.DigestFile = func(filePath string) ([]byte, error) {
			res, _ := hex.DecodeString("aaaa3333")

			return res, nil
		}
		audit.VerifySignature = func(publicKeyPEM []byte, digestHex, signatureBase64 string) (bool, error) {
			return false, nil
		}

		up, _ := New(context.Background(), models.ApplicationMeta{
			AuthorsPublicKey: assets.PublicKeyPEM,
			Version:          "v1.2.2",
			OS:               "linux",
			Arch:             "amd64",
		}, func(newVersion *os.File) {
			assert.Fail(t, "should not call callback when artifact signature verification fails")
		})

		var buf bytes.Buffer
		logger := zerolog.New(&buf)
		up.Logger = &logger

		up.Run()
		assert.Contains(t, buf.String(), "Artifact signature verification failed. Artifact did not come from authors")
		assert.NotEmpty(t, fileName, "Temporary artifact file should have been created")
		assert.NoFileExists(t, fileName, "Temporary artifact file should not exist")
	})

	t.Run("should call callback with new version file", func(t *testing.T) {
		oldDownload := downloader.DownloadToTemporaryFile
		oldVerify := audit.VerifySignature
		oldDigest := digest.DigestFile
		var fileName string
		defer func() {
			downloader.DownloadToTemporaryFile = oldDownload
			audit.VerifySignature = oldVerify
			digest.DigestFile = oldDigest
		}()
		downloader.DownloadToTemporaryFile = func(ctx context.Context, url, _ string) (*os.File, error) {
			file, _ := os.CreateTemp("", "artifact")
			fileName = file.Name()
			return file, nil
		}
		digest.DigestFile = func(filePath string) ([]byte, error) {
			res, _ := hex.DecodeString("aaaa3333")

			return res, nil
		}
		audit.VerifySignature = func(publicKeyPEM []byte, digestHex, signatureBase64 string) (bool, error) {
			return true, nil
		}

		up, _ := New(context.Background(), models.ApplicationMeta{
			AuthorsPublicKey: assets.PublicKeyPEM,
			Version:          "v1.2.2",
			OS:               "linux",
			Arch:             "amd64",
		}, func(newVersion *os.File) {
			assert.NotNil(t, newVersion, "Callback should be called with new version file")
			assert.Equal(t, fileName, newVersion.Name(), "Callback should receive the correct new version file")
			assert.FileExists(t, fileName, "New version file should exist")

			_ = newVersion.Close()
			_ = os.Remove(fileName)
		})

		var buf bytes.Buffer
		logger := zerolog.New(&buf)
		up.Logger = &logger

		up.Run()
		assert.Contains(t, buf.String(), "Downloaded artifact file")
	})
}
