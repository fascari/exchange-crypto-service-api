package modules

import (
	healthhandler "exchange-crypto-service-api/internal/app/health/handler"

	"go.uber.org/fx"
)

var healthInvokes = fx.Invoke(
	func(params RouterParams) {
		params.Router.HandleFunc("/health", healthhandler.HealthCheck).Methods("GET")
	},
)

var Health = fx.Options(
	healthInvokes,
)
