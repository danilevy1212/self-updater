package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/danilevy1212/self-updater/internal/digest"
	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/server"
	"github.com/danilevy1212/self-updater/internal/updater"
)

var (
	Version                      string = "development"
	Commit                       string = "unknown"
	swapping                            = flag.Bool("swapping", false, "runs the swapping procedure instead of the server")
	newVersionPath                      = flag.String("new-version-path", "", "the new version to swap with")
	originalExecutableLocation          = flag.String("original-executable-location", "", "the original path that must be swapped to")
	originalExecutableBackupPath        = flag.String("original-executable-backup-path", "", "the original executable backup path, used for swapping")
)

func main() {
	flag.Parse()

	currentExecutablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	d, err := digest.DigestFile(currentExecutablePath)
	if err != nil {
		fmt.Println("Error creating file digest:", err)
		return
	}
	am := models.NewApplicationMeta(
		d,
		Version,
		Commit,
		currentExecutablePath,
	)
	ctx := context.Background()

	if err != nil {
		fmt.Println("Error creating updater:", err)
		return
	}

	server, err := server.New(ctx, am)
	if err != nil {
		fmt.Println("Error creating application:", err)
		return
	}

	updater, err := updater.New(ctx, am, func(newVersion *os.File, u *updater.Updater) {
		defer newVersion.Close()
		logger := u.Logger.With().
			Str("handler", "OnUpgradeReadyCallback").
			Logger()

		if err = server.Shutdown(ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("Error shutting down server before upgrade")

			return
		}

		logger.Info().
			Msg("Server shutdown successfully, proceeding with upgrade")

		if err := u.LaunchSwap(newVersion.Name()); err != nil {
			logger.Error().
				Err(err).
				Msg("Error launching swap process")
		} else {
			logger.Info().
				Msg("Swap process launched successfully, exiting current process")

			os.Exit(0)
		}
	})
	if logger := updater.Logger; *swapping {
		if *newVersionPath == "" {
			logger.Error().
				Msg("A new-version flag must be provided for swapping")
			os.Exit(1)
		}

		if *originalExecutableBackupPath == "" {
			logger.Error().
				Msg("An original-executable-backup-path flag must be provided for swapping")
			os.Exit(1)
		}

		if *originalExecutableLocation == "" {
			logger.Error().
				Msg("An original-executable-location flag must be provided for swapping")
			os.Exit(1)
		}

		if err = updater.Swap(*newVersionPath, *originalExecutableBackupPath, *originalExecutableLocation); err != nil {
			logger.Error().
				Err(err).
				Msg("Error swapping versions")

			os.Exit(1)
		}

		return
	}

	if updater.Config.RunAtBoot {
		updater.Run()
	}

	if _, err := updater.Start(); err != nil {
		fmt.Println("Error starting updater:", err)
		return
	}

	server.RegisterGlobalMiddleware()
	server.RegisterRoutes()

	if err := server.Serve(server.Config.Port); err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
