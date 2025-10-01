package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Database struct {
		MaxLifetimeConnections time.Duration
		Conn                   DBConnection
		MaxIdleConnections     int
		MaxOpenConnections     int
	}

	DBConnection struct {
		Host         string
		Port         string
		Username     string
		Password     string
		DatabaseName string
		SslMode      string
		Schema       string
	}
)

func loadDatabaseConfig() (Database, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return Database{}, err
	}

	return Database{
		Conn: DBConnection{
			Host:         viper.GetString("DB_URL"),
			Port:         viper.GetString("DB_PORT"),
			Username:     viper.GetString("DB_APP_USER"),
			Password:     viper.GetString("DB_APP_PASSWORD"),
			DatabaseName: viper.GetString("DB_NAME"),
			SslMode:      viper.GetString("DB_SSL_MODE"),
			Schema:       viper.GetString("DB_SCHEMA"),
		},
		MaxIdleConnections:     viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
		MaxOpenConnections:     viper.GetInt("DB_MAX_OPEN_CONNECTIONS"),
		MaxLifetimeConnections: viper.GetDuration("DB_MAX_LIFETIME_CONNECTIONS") * time.Minute,
	}, nil
}
