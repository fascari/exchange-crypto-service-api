package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TraceRepository(ctx context.Context, operation string, fn func(ctx context.Context) ([]attribute.KeyValue, error), attrs ...attribute.KeyValue) error {
	ctx, span := otel.Tracer("repository").Start(ctx, operation)
	defer span.End()

	span.SetAttributes(attrs...)

	additionalAttrs, err := fn(ctx)
	if len(additionalAttrs) > 0 {
		span.SetAttributes(additionalAttrs...)
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
