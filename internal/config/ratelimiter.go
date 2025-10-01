package config

import (
	"time"

	"github.com/spf13/viper"
)

type RateLimiter struct {
	RequestsPerSecond float64
	BurstSize         int
	CleanupInterval   time.Duration
}

func loadRateLimiterConfig() (RateLimiter, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return RateLimiter{}, err
	}

	return RateLimiter{
		RequestsPerSecond: viper.GetFloat64("RATE_LIMITER_REQUESTS_PER_SECOND"),
		BurstSize:         viper.GetInt("RATE_LIMITER_BURST_SIZE"),
		CleanupInterval:   viper.GetDuration("RATE_LIMITER_CLEANUP") * time.Minute,
	}, nil
}

func defaultRateLimiterConfig() RateLimiter {
	return RateLimiter{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute * 5,
	}
}
