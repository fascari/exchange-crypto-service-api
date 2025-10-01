package config

import (
	"github.com/rs/zerolog/log"
)

func LoadDatabase() Database {
	return loadConfigWithPanic(loadDatabaseConfig, "failed to load database config")
}

func LoadJWT() JWT {
	return loadConfigWithDefaults(loadJWTConfig, defaultJWTConfig, "failed to load JWT config, using defaults")
}

func LoadRateLimiter() RateLimiter {
	return loadConfigWithDefaults(loadRateLimiterConfig, defaultRateLimiterConfig, "failed to load rate limiter config, using defaults")
}

func loadConfigWithPanic[T any](loader func() (T, error), errorMsg string) T {
	config, err := loader()
	if err != nil {
		log.Error().Err(err).Msg(errorMsg)
		panic(err)
	}
	return config
}

func loadConfigWithDefaults[T any](loader func() (T, error), defaultProvider func() T, errorMsg string) T {
	config, err := loader()
	if err != nil {
		log.Warn().Err(err).Msg(errorMsg)
		config = defaultProvider()
	}
	return config
}
