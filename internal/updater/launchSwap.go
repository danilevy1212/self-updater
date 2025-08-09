package updater

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/danilevy1212/self-updater/internal/updater/utils"
)

func (u *Updater) LaunchSwap(newVersionPath string) error {
	logger := u.Logger.With().Str("handler", "LaunchSwap").Logger()
	finalDestination := u.Meta.ExecutablePath
	finalDestinationDir := filepath.Dir(finalDestination)
	finalFileName := filepath.Base(finalDestination)

	// Create a backup, which we will run as our swapper
	backupFilePath := filepath.Join(finalDestinationDir, "backup."+finalFileName)
	backupFile, err := os.Create(backupFilePath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create backup file")

		return err
	}
	backupFile.Close()

	err = utils.CopyFile(finalDestination, backupFilePath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to copy current executable to backup file")
		_ = os.Remove(backupFile.Name())

		return err
	}
	// Ensure executable on POSIX
	if u.Meta.OS != "windows" {
		_ = os.Chmod(backupFilePath, 0o755)
	}

	// Build swapper command with your flags
	args := []string{
		"-swapping",
		"-new-version-path", newVersionPath,
		"-original-executable-location", finalDestination,
		"-original-executable-backup-path", backupFilePath,
	}
	cmd := exec.Command(backupFilePath, args...)
	cmd.Env = os.Environ()

	// Detach
	utils.SetDetach(cmd)

	if err := cmd.Start(); err != nil {
		logger.Error().Err(err).Msg("failed to launch swapper process")
		return err
	}

	logger.Info().
		Str("swapper", backupFilePath).
		Str("new_version", newVersionPath).
		Msg("swapper launched")

	return nil
}
