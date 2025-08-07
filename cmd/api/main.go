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
	Version     string = "development"
	Commit      string = "unknown"
	checkUpdate        = flag.Bool("check-update-boot", false, "Check for updates immediately and exit if updated")
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
	)
	ctx := context.Background()

	fmt.Printf("Current executable sha256: %s\n", am.DigestString())

	updater, err := updater.New(ctx, am)
	if err != nil {
		fmt.Println("Error creating updater:", err)
		return
	}

	server, err := server.New(ctx, am)
	if err != nil {
		fmt.Println("Error creating application:", err)
		return
	}

	server.RegisterGlobalMiddleware()
	server.RegisterRoutes()

	if *checkUpdate {
		updater.Run()
	}

	if _, err := updater.Start(); err != nil {
		fmt.Println("Error starting updater:", err)
		return
	}

	if err := server.Serve(server.Config.Port); err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	// TODO Check here that server was cancelled due to a newer version being available.
}
