package infra

import (
	"context"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

type App struct {
	Router         *mux.Router
	DB             *gorm.DB
	TracerProvider *trace.TracerProvider
}

func (a *App) Start() {
	server := startHTTPServer(a)
	waitForShutdownSignal()
	shutdown(server, a)
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.TracerProvider != nil {
		return a.TracerProvider.Shutdown(ctx)
	}
	return nil
}
