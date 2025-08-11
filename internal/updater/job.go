package updater

import (
	"context"
	"encoding/hex"
	"os"
	"time"

	"github.com/danilevy1212/self-updater/internal/audit"
	"github.com/danilevy1212/self-updater/internal/digest"
	"github.com/danilevy1212/self-updater/internal/downloader"
)

func (u *Updater) Run() {
	logger := u.Logger

	logger.Info().
		Str("version", u.Meta.Version).
		Str("commit", u.Meta.Commit).
		Str("digest", u.Meta.DigestString()).
		Str("public_key", string(u.Meta.AuthorsPublicKey)).
		Msg("Running updater job")

	logger.Info().
		Msg("fetching latest manifest")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	manifest, err := u.ManifestFetcher.FetchManifest(ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to fetch manifest")
		return
	}

	if manifest.PublicKey != string(u.Meta.AuthorsPublicKey) {
		logger.Error().
			Str("manifest_public_key", manifest.PublicKey).
			Str("application_public_key", string(u.Meta.AuthorsPublicKey)).
			Msg("Manifest public key does not match application public key")

		return
	}

	// NOTE  For the sake of time, I will only do this check, but it probably makes sense to make other checks here:
	//  - Is the current binary tampered? (Digest AND signature won't match)
	//  - Is the current version in the manifest (Only check this if we are not in DEV mode)
	// Should log out an error and stop in those cases.
	latestVersion := manifest.Latest
	if u.Meta.Version == latestVersion {
		logger.Info().
			Msg("No updates available. Current version is up to date.")

		return
	}

	matchingVersion, err := manifest.GetVersionInfo(latestVersion)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get version info from manifest")

		return
	}

	artifactForPlatform, err := matchingVersion.GetArtifactForPlatform(u.Meta.OS, u.Meta.Arch)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get artifact for platform from manifest")

		return
	}

	ctxDownload, cancelDownload := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelDownload()
	artifactFile, err := downloader.DownloadToTemporaryFile(
		ctxDownload,
		artifactForPlatform.URL,
		artifactForPlatform.Filename+".*",
	)

	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to download artifact file")

		return
	}

	cleanArtifactTmp := func() {
		_ = artifactFile.Close()
		_ = os.Remove(artifactFile.Name())
	}

	logger.Info().
		Str("artifact_file", artifactFile.Name()).
		Msg("Downloaded artifact file")

	artifactDigest, err := digest.DigestFile(artifactFile.Name())
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to calculate artifact file digest")

		cleanArtifactTmp()
		return
	}

	artifactDigestHex := hex.EncodeToString(artifactDigest)
	if artifactDigestHex != artifactForPlatform.Digest {
		logger.Error().
			Str("expected_digest", artifactForPlatform.Digest).
			Str("actual_digest", artifactDigestHex).
			Msg("Artifact file digest does not match expected digest")

		cleanArtifactTmp()
		return
	}

	isVerified, err := audit.VerifySignature(
		u.Meta.AuthorsPublicKey,
		artifactDigestHex,
		artifactForPlatform.SignatureBase64,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to verify artifact signature")

		cleanArtifactTmp()
		return
	}
	if !isVerified {
		logger.Error().
			Msg("Artifact signature verification failed. Artifact did not come from authors")

		cleanArtifactTmp()
		return
	}

	u.OnUpgradeReady(artifactFile)
}
