package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init(verbose bool) {
	var w io.Writer = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}
	log = zerolog.New(w).With().Timestamp().Logger().Level(level)
}

func Get() *zerolog.Logger {
	return &log
}
