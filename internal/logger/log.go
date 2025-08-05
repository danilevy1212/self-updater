package logger

import (
	"io"
	"os"
	"time"
	"github.com/rs/zerolog"
)

func New(
	pretty bool,
) *zerolog.Logger {
	var w io.Writer

	if pretty {
		w = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		w = os.Stdout
	}

	l := zerolog.New(w).With().Timestamp().Logger()

	return &l
}
