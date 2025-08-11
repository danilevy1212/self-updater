package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/danilevy1212/self-updater/internal/digest"
	"github.com/danilevy1212/self-updater/internal/models"
)

var (
	Version          string = "development"
	Commit           string = "unknown"
	asServer                = flag.Bool("server", false, "run server + updater process directly")
	sessionDirectory        = flag.String("current-session-dir", "", "directory to store session files")
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
		fmt.Println("Error creating current executable digest:", err)
		return
	}
	am := models.NewApplicationMeta(
		d,
		Version,
		Commit,
		currentExecutablePath,
	)
	ctx := context.Background()

	if *asServer {
		code := runServer(ctx, am)
		os.Exit(code)
	} else {
		runLauncher(ctx, am)
	}
}
