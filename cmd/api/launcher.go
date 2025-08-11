package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/danilevy1212/self-updater/internal/launcher"
	"github.com/danilevy1212/self-updater/internal/launcher/utils"
	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/models/exitcodes"
)

func runLauncher(ctx context.Context, am models.ApplicationMeta) {
	launcherOrchestrator, err := launcher.New(ctx, am)
	logger := launcherOrchestrator.Logger.With().
		Str("handler", "runLauncher").
		Logger()

	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create launcher orchestrator")

		return
	}

	// Copy this binary to temp file.
	err = launcherOrchestrator.CreateSessionDir()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to create session directory")

		return
	}
	currentName := getCurrentServerFileName(am)
	newName := getNewServerFileName(am)

	currentPath := filepath.Join(launcherOrchestrator.SessionDirectory, currentName)
	newPath := filepath.Join(launcherOrchestrator.SessionDirectory, newName)

	err = utils.CopyFile(am.ExecutablePath, currentPath)
	if err != nil {
		logger.Error().
			Err(err).
			Str("currentPath", currentPath).
			Msg("Failed to copy current binary to session directory")

		return
	}

	cmd, err := launcherOrchestrator.LaunchServer(ctx, currentPath)
	if err != nil {
		logger.Error().
			Err(err).
			Str("currentPath", currentPath).
			Msg("Failed to launch server process")

		return
	}

	for {
		logger.Info().
			Msg("Waiting for server to signal update ready")

		_ = cmd.Wait()
		code := cmd.ProcessState.ExitCode()

		if code != exitcodes.ExitUpdateReady {
			logger.Error().
				Int("exitCode", code).
				Msg("Server exited with unexpected code, not an update ready signal")

			return
		}

		if _, err := os.Stat(newPath); err != nil {
			if !os.IsNotExist(err) {
				logger.Error().
					Err(err).
					Str("newPath", newPath).
					Msg("Failed to stat new binary")

				return
			}
		}

		if err := os.Chmod(newPath, 0o700); err != nil {
			logger.Error().
				Err(err).
				Str("newPath", newPath).
				Msg("Failed to set permissions on new binary")

			return
		}

		// In windows, there is no atomic rename, we need to delete current first.
		if am.OS == "windows" {
			_ = os.Remove(currentPath)
		}

		// Swap!
		if err := os.Rename(newPath, currentPath); err != nil {
			logger.Error().
				Err(err).
				Str("currentPath", currentPath).
				Str("newPath", newPath).
				Msg("Failed to swap new binary into place")

			// NOTE  If we were doing backups, it would be here where we do the rollback
			return
		}

		cmd, err = launcherOrchestrator.LaunchServer(ctx, currentPath)
		if err != nil {
			logger.Error().
				Err(err).
				Str("currentPath", currentPath).
				Msg("Failed to relaunch server process after update")

			return
		}
	}
}
