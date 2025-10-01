package config

import (
	"time"

	"github.com/spf13/viper"
)

type RateLimiterConfig struct {
	RequestsPerSecond float64
	BurstSize         int
	CleanupInterval   time.Duration
}

func LoadRateLimiterConfig() (RateLimiterConfig, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return RateLimiterConfig{}, err
	}

	return RateLimiterConfig{
		RequestsPerSecond: viper.GetFloat64("RATE_LIMITER_REQUESTS_PER_SECOND"),
		BurstSize:         viper.GetInt("RATE_LIMITER_BURST_SIZE"),
		CleanupInterval:   viper.GetDuration("RATE_LIMITER_CLEANUP") * time.Minute,
	}, nil
}

func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute * 5,
	}
}
