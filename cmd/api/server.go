package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/danilevy1212/self-updater/internal/models"
	"github.com/danilevy1212/self-updater/internal/models/exitcodes"
	"github.com/danilevy1212/self-updater/internal/server"
	"github.com/danilevy1212/self-updater/internal/updater"
)

func runServer(ctx context.Context, am models.ApplicationMeta) int {
	if *sessionDirectory == "" {
		fmt.Println("Session directory is not set. Please provide a valid session directory using --current-session-dir flag.")
		return exitcodes.ExitFatal
	}

	app, err := server.New(ctx, am)
	if err != nil {
		fmt.Println("Error creating server:", err)
		return exitcodes.ExitFatal
	}

	var exitCode atomic.Int32
	exitCode.Store(int32(exitcodes.ExitOK))

	updater, err := updater.New(ctx, am, func(newVersion *os.File) {
		defer newVersion.Close()
		newPath := filepath.Join(*sessionDirectory, getNewServerFileName(am))

		if err := os.Rename(newVersion.Name(), newPath); err != nil {
			fmt.Println("rename new version:", err)
			exitCode.Store(int32(exitcodes.ExitFatal))
		} else {
			exitCode.Store(int32(exitcodes.ExitUpdateReady))
		}

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		_ = app.Shutdown(c)
	})
	if err != nil {
		fmt.Println("Error creating updater:", err)
		return exitcodes.ExitFatal
	}

	if updater.Config.RunAtBoot {
		updater.Run()
	}

	if _, err := updater.Start(); err != nil {
		fmt.Println("Error starting updater:", err)
		return exitcodes.ExitFatal
	}

	app.RegisterGlobalMiddleware()
	app.RegisterRoutes()

	if err := app.Serve(app.Config.Port); err != nil {
		fmt.Println("Error starting server:", err)
		return exitcodes.ExitFatal
	}

	return int(exitCode.Load())
}
