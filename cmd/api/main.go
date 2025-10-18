package main

import (
	"exchange-crypto-service-api/cmd/api/modules"
	"exchange-crypto-service-api/internal/bootstrap"
	"exchange-crypto-service-api/internal/jwt"
	"exchange-crypto-service-api/pkg/logger"

	"go.uber.org/fx"
)

func main() {
	logger.Init()
	jwt.Initialize()

	app := fx.New(
		bootstrap.Logger(),
		// Infrastructure
		bootstrap.Config,
		bootstrap.Database,
		bootstrap.Telemetry,
		bootstrap.Router,
		bootstrap.Server,
		// Business Modules
		modules.Health,
		modules.Exchange,
		modules.User,
		modules.Auth,
		modules.Account,
		modules.Transaction,
	)

	app.Run()
}
