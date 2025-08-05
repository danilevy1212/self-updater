package logger

import (
	"io"
	"os"
	"testing"
	"github.com/gkampitakis/go-snaps/snaps"
)

func Test_logger_New(t *testing.T) {
	t.Run("should log out pretty logs pretty = true", func(t *testing.T) {
		// Capture Stdout
		oldOut := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		l := New(true)

		if l == nil {
			t.Fatalf("Expected logger to be created, got nil")
		}

		l.Info().Msg("I'm supposed to be pretty")

		w.Close()

		msg, err := io.ReadAll(r)
		os.Stdout = oldOut

		if err != nil {
			t.Fatalf("Failed to read from pipe: %v", err)
		}

		snaps.MatchSnapshot(t, string(msg))
	})

	t.Run("should print out JSON logs pretty = false", func(t *testing.T) {
		// Capture Stdout
		oldOut := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		l := New(false)

		if l == nil {
			t.Fatalf("Expected logger to be created, got nil")
		}

		l.Info().Msg("I'm supposed to be JSON")
		w.Close()

		msg, err := io.ReadAll(r)
		os.Stdout = oldOut

		if err != nil {
			t.Fatalf("Failed to read from pipe: %v", err)
		}

		snaps.MatchSnapshot(t, string(msg))
	})

}
