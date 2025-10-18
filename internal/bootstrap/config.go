package bootstrap

import (
	"exchange-crypto-service-api/internal/config"

	"go.uber.org/fx"
)

var Config = fx.Module("config",
	fx.Provide(
		config.LoadDatabase,
		config.LoadJWT,
		config.LoadRateLimiter,
	),
)
