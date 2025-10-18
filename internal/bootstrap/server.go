package bootstrap

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

const (
	serverPort  = ":8080"
	serviceName = "exchange-crypto-service-api"
)

var Server = fx.Module("server",
	fx.Invoke(NewHTTPServer),
)

type ServerParams struct {
	fx.In

	Lifecycle  fx.Lifecycle
	MainRouter *mux.Router `name:"main"`
}

func NewHTTPServer(params ServerParams) {
	server := &http.Server{
		Addr:    serverPort,
		Handler: params.MainRouter,
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				log.Info().Msgf("server starting on %s", serverPort)
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Msg("server error")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("shutting down server")
			return server.Shutdown(ctx)
		},
	})
}
