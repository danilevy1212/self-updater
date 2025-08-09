package manifest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danilevy1212/self-updater/internal/assets"
	"github.com/danilevy1212/self-updater/internal/downloader"
	"github.com/danilevy1212/self-updater/internal/logger"
	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/models/fixtures"
)

func TestSetGithubManifestFetcherLogger(t *testing.T) {
	logger := logger.New(true)
	ctx := context.Background()

	newCtx := SetGithubManifestFetcherLogger(ctx, logger)
	ctxLogger := newCtx.Value(loggerKey)

	assert.NotNil(t, newCtx, "new context should not be nil")
	assert.NotNil(t, ctxLogger, "context should contain a logger")
	assert.Equal(t, logger, ctxLogger, "context logger should match the one set")
}

func TestNewGithubManifestFetcher(t *testing.T) {
	t.Run("should create a new manifest fetcher with a logger", func(t *testing.T) {
		ctx := SetGithubManifestFetcherLogger(context.Background(), logger.New(true))
		applicationMeta := models.ApplicationMeta{}

		fetcher, err := NewGithubManifestFetcher(ctx, applicationMeta)

		assert.NoError(t, err, "should not return an error when creating a new fetcher")
		assert.NotNil(t, fetcher, "fetcher should not be nil")
		assert.Equal(t, applicationMeta, fetcher.ApplicationMeta, "application meta should match")
		assert.NotNil(t, fetcher.Logger, "fetcher should have a logger")
	})

	t.Run("should return an error if logger is missing", func(t *testing.T) {
		ctx := context.Background()
		applicationMeta := models.ApplicationMeta{}

		fetcher, err := NewGithubManifestFetcher(ctx, applicationMeta)

		assert.Error(t, err, "should return an error when logger is missing")
		assert.Nil(t, fetcher, "fetcher should be nil when there is an error")
	})

	t.Run("should return an error if logger is nil", func(t *testing.T) {
		ctx := SetGithubManifestFetcherLogger(context.Background(), nil)
		applicationMeta := models.ApplicationMeta{}

		fetcher, err := NewGithubManifestFetcher(ctx, applicationMeta)

		assert.Error(t, err, "should return an error when logger is nil")
		assert.Nil(t, fetcher, "fetcher should be nil when there is an error")
	})
}

func Test_GithubManifestFetcher_FetchManifest(t *testing.T) {
	ctx := SetGithubManifestFetcherLogger(context.Background(), logger.New(true))

	t.Run("should cancel one download if the other fails", func(t *testing.T) {
		originalDownloader := downloader.DownloadToTemporaryFile
		defer func() {
			downloader.DownloadToTemporaryFile = originalDownloader
		}()

		// expected URLs
		manifestURL := "https://github.com/acme/widget/releases/latest/download/release.json"
		sigURL := manifestURL + ".sig.base64"

		var (
			manifestCalled int32
			sigCalled      int32
			sawCanceled    atomic.Bool
		)

		downloader.DownloadToTemporaryFile = func(dctx context.Context, url, pattern string) (*os.File, error) {
			switch url {
			case manifestURL:
				atomic.AddInt32(&manifestCalled, 1)
				return nil, errors.New("manifest boom")
			case sigURL:
				atomic.AddInt32(&sigCalled, 1)
				<-dctx.Done()
				sawCanceled.Store(errors.Is(dctx.Err(), context.Canceled))
				return nil, dctx.Err()
			default:
				t.Fatalf("unexpected url: %s", url)
				return nil, nil
			}
		}

		meta := models.ApplicationMeta{
			SourceInfo: models.SourceInfo{
				Host:  "github.com",
				Owner: "acme",
				Name:  "widget",
			},
			AuthorsPublicKey: assets.PublicKeyPEM,
		}

		fetcher, err := NewGithubManifestFetcher(ctx, meta)
		assert.NoError(t, err)

		_, err = fetcher.FetchManifest(context.Background())
		assert.Error(t, err)

		assert.Equal(t, int32(1), atomic.LoadInt32(&manifestCalled))
		assert.Equal(t, int32(1), atomic.LoadInt32(&sigCalled))
		assert.True(t, sawCanceled.Load(), "sig download should observe context cancellation")
	})

	t.Run("should verify digest with signature", func(t *testing.T) {
		originalDownloader := downloader.DownloadToTemporaryFile
		defer func() {
			downloader.DownloadToTemporaryFile = originalDownloader
		}()

		// expected URLs
		manifestURL := "https://github.com/acme/widget/releases/latest/download/release.json"
		sigURL := manifestURL + ".sig.base64"

		var (
			manifestCalled            int32
			sigCalled                 int32
			manifestFilePath          string
			manifestSignatureFilePath string
		)

		downloader.DownloadToTemporaryFile = func(dctx context.Context, url, pattern string) (*os.File, error) {
			switch url {
			case manifestURL:
				atomic.AddInt32(&manifestCalled, 1)
				manifestFile, err := os.CreateTemp("", "release.json.*")
				if err != nil {
					t.Fatalf("failed to create temp manifest file: %v", err)
				}
				_, err = io.Copy(manifestFile, bytes.NewReader(fixtures.ReleaseFixture))
				if err != nil {
					_ = manifestFile.Close()
					t.Fatalf("failed to write to temp manifest file: %v", err)
				}
				if _, err := manifestFile.Seek(0, io.SeekStart); err != nil {
					_ = manifestFile.Close()
					t.Fatalf("failed to seek to start of temp manifest file: %v", err)
				}
				return manifestFile, nil
			case sigURL:
				atomic.AddInt32(&sigCalled, 1)
				manifestSignatureFile, err := os.CreateTemp("", "release.json.sig.base64*")
				if err != nil {
					t.Fatalf("failed to create temp signature file: %v", err)
				}
				_, err = io.Copy(manifestSignatureFile, bytes.NewReader(fixtures.ReleaseSignatureFixture))
				if err != nil {
					_ = manifestSignatureFile.Close()
					t.Fatalf("failed to write to temp signature file: %v", err)
				}
				if _, err := manifestSignatureFile.Seek(0, io.SeekStart); err != nil {
					_ = manifestSignatureFile.Close()
					t.Fatalf("failed to seek to start of temp signature file: %v", err)
				}
				return manifestSignatureFile, nil
			default:
				t.Fatalf("unexpected url: %s", url)
				return nil, nil
			}
		}
		ctx := SetGithubManifestFetcherLogger(context.Background(), logger.New(true))
		meta := models.ApplicationMeta{
			SourceInfo: models.SourceInfo{
				Host:  "github.com",
				Owner: "acme",
				Name:  "widget",
			},
			AuthorsPublicKey: assets.PublicKeyPEM, // must match how the signature was produced
		}
		fetcher, err := NewGithubManifestFetcher(ctx, meta)
		assert.NoError(t, err)

		got, err := fetcher.FetchManifest(context.Background())
		assert.NoError(t, err, "happy path should verify and parse")
		assert.NotNil(t, got)
		assert.Equal(t, int32(1), atomic.LoadInt32(&manifestCalled), "should have called manifest download once")
		assert.Equal(t, int32(1), atomic.LoadInt32(&sigCalled), "shoulld have called signature download once")
		assert.NoFileExists(t, manifestFilePath, "manifest file should have been cleaned-up")
		assert.NoFileExists(t, manifestSignatureFilePath, "signature file should have been cleaned-up")
	})

	t.Run("should error if digest doesn't match signature", func(t *testing.T) {
		originalDownloader := downloader.DownloadToTemporaryFile
		defer func() {
			downloader.DownloadToTemporaryFile = originalDownloader
		}()

		// expected URLs
		manifestURL := "https://github.com/acme/widget/releases/latest/download/release.json"
		sigURL := manifestURL + ".sig.base64"

		var (
			manifestCalled            int32
			sigCalled                 int32
			manifestFilePath          string
			manifestSignatureFilePath string
		)

		downloader.DownloadToTemporaryFile = func(dctx context.Context, url, pattern string) (*os.File, error) {
			switch url {
			case manifestURL:
				atomic.AddInt32(&manifestCalled, 1)
				manifestFile, err := os.CreateTemp("", "release.json.*")
				if err != nil {
					t.Fatalf("failed to create temp manifest file: %v", err)
				}
				_, err = io.Copy(manifestFile, bytes.NewReader(fixtures.ReleaseFixture))
				if err != nil {
					_ = manifestFile.Close()
					t.Fatalf("failed to write to temp manifest file: %v", err)
				}
				if _, err := manifestFile.Seek(0, io.SeekStart); err != nil {
					_ = manifestFile.Close()
					t.Fatalf("failed to seek to start of temp manifest file: %v", err)
				}
				manifestFilePath = manifestFile.Name()
				return manifestFile, nil
			case sigURL:
				atomic.AddInt32(&sigCalled, 1)
				manifestSignatureFile, err := os.CreateTemp("", "release.json.sig.base64*")
				if err != nil {
					t.Fatalf("failed to create temp signature file: %v", err)
				}
				_, err = io.Copy(
					manifestSignatureFile,
					bytes.NewReader(
						[]byte(`NCqAUE2dp7828kzmHzImiGF8AdRdo/Sr+ZJ0FlQWOJJUZm4Qf4eyO/52+nSfg/fs81sZ28rC6m+hrah9kivkBA==`),
					),
				)
				if err != nil {
					_ = manifestSignatureFile.Close()
					t.Fatalf("failed to write to temp signature file: %v", err)
				}
				if _, err := manifestSignatureFile.Seek(0, io.SeekStart); err != nil {
					_ = manifestSignatureFile.Close()
					t.Fatalf("failed to seek to start of temp signature file: %v", err)
				}
				manifestSignatureFilePath = manifestSignatureFile.Name()
				return manifestSignatureFile, nil
			default:
				t.Fatalf("unexpected url: %s", url)
				return nil, nil
			}
		}
		ctx := SetGithubManifestFetcherLogger(context.Background(), logger.New(true))
		meta := models.ApplicationMeta{
			SourceInfo: models.SourceInfo{
				Host:  "github.com",
				Owner: "acme",
				Name:  "widget",
			},
			AuthorsPublicKey: assets.PublicKeyPEM, // must match how the signature was produced
		}
		fetcher, err := NewGithubManifestFetcher(ctx, meta)
		assert.NoError(t, err)

		got, err := fetcher.FetchManifest(context.Background())
		assert.Error(t, err, "should return an error when digest does not match signature")
		assert.ErrorContains(t, err, "fetched manifest did not come from authors", "error should indicate digest mismatch")
		assert.Nil(t, got, "should return nil manifest on error")
		assert.Equal(t, int32(1), atomic.LoadInt32(&manifestCalled), "should have called manifest download once")
		assert.Equal(t, int32(1), atomic.LoadInt32(&sigCalled), "shoulld have called signature download once")
		assert.NoFileExists(t, manifestFilePath, "manifest file should have been cleaned-up")
		assert.NoFileExists(t, manifestSignatureFilePath, "signature file should have been cleaned-up")
	})
}
