package config

import (
	"time"

	"github.com/spf13/viper"
)

type JWT struct {
	Secret          string
	ExpirationHours time.Duration
}

func loadJWTConfig() (JWT, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return JWT{}, err
	}

	return JWT{
		Secret:          viper.GetString("JWT_SECRET"),
		ExpirationHours: viper.GetDuration("JWT_EXPIRATION_HOURS") * time.Hour,
	}, nil
}

func defaultJWTConfig() JWT {
	return JWT{
		Secret:          "default-secret-change-me",
		ExpirationHours: 24 * time.Hour,
	}
}
