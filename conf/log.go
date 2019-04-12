package conf

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(debug bool) zerolog.Logger {
	if debug {
		w := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = os.Stderr
		})
		return zerolog.New(w).Level(zerolog.DebugLevel).With().Timestamp().Caller().Logger()
	}
	return zerolog.New(os.Stderr).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
}
