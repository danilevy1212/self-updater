package manifest

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"

	"github.com/danilevy1212/self-updater/internal/audit"
	"github.com/danilevy1212/self-updater/internal/digest"
	"github.com/danilevy1212/self-updater/internal/downloader"
	"github.com/danilevy1212/self-updater/internal/models"
)

type ctxKeyLogger struct{}

var loggerKey = ctxKeyLogger{}

type GithubManifestFetcher struct {
	ApplicationMeta models.ApplicationMeta
	Logger          *zerolog.Logger
}

func getLogger(ctx context.Context) (*zerolog.Logger, error) {
	l, ok := ctx.Value(loggerKey).(*zerolog.Logger)
	if !ok {
		return nil, errors.New("context does not have a logger")
	}

	return l, nil
}

func SetGithubManifestFetcherLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func NewGithubManifestFetcher(ctx context.Context, applicationMeta models.ApplicationMeta) (*GithubManifestFetcher, error) {
	logger, err := getLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("missing a logger: %w", err)
	}

	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &GithubManifestFetcher{
		ApplicationMeta: applicationMeta,
		Logger:          logger,
	}, nil
}

func (ghf *GithubManifestFetcher) FetchManifest(ctx context.Context) (*models.ReleaseManifest, error) {
	logger := ghf.Logger
	meta := ghf.ApplicationMeta

	releaseJSONURL := fmt.Sprintf(
		"https://%s/%s/%s/releases/latest/download/release.json",
		meta.SourceInfo.Host,
		meta.SourceInfo.Owner,
		meta.SourceInfo.Name,
	)
	releaseJSONSignatureURL := releaseJSONURL + ".sig.base64"

	logger.Info().
		Str("release_json_url", releaseJSONURL).
		Str("release_json_signature_url", releaseJSONSignatureURL).
		Msg("Fetching manifest from GitHub")

	var (
		wg                    sync.WaitGroup
		manifestFile, sigFile *os.File
		manifestErr, sigErr   error
	)

	// Fetch both files in parallel, cancel all if one fails
	downloadCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg.Add(2)
	go func() {
		defer wg.Done()
		manifestFile, manifestErr = downloader.DownloadToTemporaryFile(downloadCtx, releaseJSONURL, "release.json.*")
		if manifestErr != nil {
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		sigFile, sigErr = downloader.DownloadToTemporaryFile(downloadCtx, releaseJSONSignatureURL, "release.json.sig.base64.*")
		if sigErr != nil {
			cancel()
		}
	}()
	logger.Info().
		Msg("Waiting for manifest and signature downloads to complete")
	wg.Wait()

	if manifestErr != nil || sigErr != nil {
		errs := []error{}

		if manifestErr != nil {
			logger.Error().
				Err(manifestErr).
				Msg("Failed to download manifest file")

			errs = append(errs, manifestErr)
		}

		if sigErr != nil {
			logger.Error().
				Err(sigErr).
				Msg("Failed to download manifest signature file")

			errs = append(errs, sigErr)
		}

		if manifestFile != nil {
			_ = manifestFile.Close()
			_ = os.Remove(manifestFile.Name())
		}

		if sigFile != nil {
			_ = sigFile.Close()
			_ = os.Remove(sigFile.Name())
		}

		return nil, fmt.Errorf("failed to download manifest files: %w", errors.Join(errs...))
	}

	defer func() {
		_ = sigFile.Close()
		_ = os.Remove(sigFile.Name())
		_ = manifestFile.Close()
		_ = os.Remove(manifestFile.Name())
	}()

	publicKey := meta.AuthorsPublicKey
	sigFileContents, err := io.ReadAll(sigFile)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to read signature file")

		return nil, fmt.Errorf("failed to read signature file: %w", err)
	}

	manifestDigestRaw, err := digest.DigestFile(manifestFile.Name())
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to compute manifest file digest")

		return nil, fmt.Errorf("failed to compute manifest file digest: %w", err)
	}

	// A little silly back-and-forth I have to do in the name of re-usability.
	isVerified, err := audit.VerifySignature(publicKey, hex.EncodeToString(manifestDigestRaw), string(sigFileContents))
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to verify manifest signature")

		return nil, fmt.Errorf("failed to verify manifest signature: %w", err)
	}

	if !isVerified {
		logger.Error().
			Msg("Manifest signature verification failed. Fetched manifest did not come from authors")

		return nil, errors.New("manifest signature verification failed: fetched manifest did not come from authors")
	}

	var result models.ReleaseManifest
	err = json.NewDecoder(manifestFile).Decode(&result)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to unmarshal manifest JSON")

		return nil, fmt.Errorf("failed to unmarshal manifest JSON: %w", err)
	}

	return &result, nil
}
