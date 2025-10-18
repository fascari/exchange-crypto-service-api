package modules

import (
	exchangerepo "exchange-crypto-service-api/internal/app/exchange/repository"

	"go.uber.org/fx"
)

var exchangeFactories = fx.Provide(
	// repositories
	exchangerepo.New,
)

var Exchange = fx.Options(
	exchangeFactories,
)
