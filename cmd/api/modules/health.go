package modules

import (
	"exchange-crypto-service-api/internal/app/health/handler"
	"exchange-crypto-service-api/internal/infra"
)

func Health(app infra.App) {
	app.Router.HandleFunc("/health", handler.HealthCheck).Methods("GET")
}
