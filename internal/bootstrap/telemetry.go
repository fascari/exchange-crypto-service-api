package bootstrap

import (
	"context"

	"exchange-crypto-service-api/pkg/telemetry"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var Telemetry = fx.Module("telemetry",
	fx.Provide(NewTracerProvider),
)

func NewTracerProvider(lc fx.Lifecycle) (*trace.TracerProvider, error) {
	tp, err := telemetry.InitTracer(serviceName)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize otel tracer")
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("shutting down tracer provider")
			return tp.Shutdown(ctx)
		},
	})

	return tp, nil
}
