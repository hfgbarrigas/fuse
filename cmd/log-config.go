// Package cmd is the entry point for cobra cli
package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

// ConfigLog configures logzero log according to cli flags
func ConfigLog() {
	// set info log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if logStackTraces {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	if prettyLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if logVerbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
