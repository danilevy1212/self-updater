package launcher

import "os"

func (l *Launcher) CreateSessionDir() error {
	logger := l.Logger.With().
		Str("handler", "CreateSessionDir").
		Logger()

	// Restrict session directory to the current user
	// This prevents other *different* OS users from writing here,
	// which blocks TOCTOU path-swap attacks in that threat model.
	//
	// If an attacker is running as the same OS user (same UID), they can
	// still modify files in this directory...
	// Protecting against that would require running the launcher/updater as
	// a dedicated user account, which is out of scope for this assignment.
	if err := os.MkdirAll(l.SessionDirectory, 0o700); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to create session directory")

		return err
	}

	return nil
}
