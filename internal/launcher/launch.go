package launcher

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func (l *Launcher) LaunchServer(ctx context.Context, serverPath string) (*exec.Cmd, error) {
	logger := l.Logger.With().
		Str("handler", "LaunchServer").
		Logger()

	launchArgs := []string{
		"--server",
		"--current-session-dir", l.SessionDirectory,
	}
	cmd := exec.CommandContext(ctx, serverPath, launchArgs...)
	cmd.Env = os.Environ()

	// Parent sees the logs of child
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	logger.Info().
		Str("path", serverPath).
		Str("args", strings.Join(launchArgs, " ")).
		Msg("launching server")

	if err := cmd.Start(); err != nil {
		logger.Error().
			Err(err).
			Str("path", serverPath).
			Str("args", strings.Join(launchArgs, " ")).
			Msg("failed to start server process")
		return nil, err
	}

	logger.Info().
		Str("path", serverPath).
		Str("args", strings.Join(launchArgs, " ")).
		Int("pid", cmd.Process.Pid).
		Msg("server process started successfully")

	return cmd, nil
}
