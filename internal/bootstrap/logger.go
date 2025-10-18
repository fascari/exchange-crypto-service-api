package bootstrap

import (
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

const verbose = "FX_VERBOSE"

func Logger() fx.Option {
	if os.Getenv(verbose) == "true" {
		return fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ConsoleLogger{W: os.Stdout}
		})
	}
	return fx.NopLogger
}
