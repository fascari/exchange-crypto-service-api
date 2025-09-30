package modules

import (
	"exchange-crypto-service-api/internal/config"
	"exchange-crypto-service-api/internal/database"
	"exchange-crypto-service-api/internal/infra"
	"exchange-crypto-service-api/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func NewApp() infra.App {
	logger.Init()

	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load database config")
		panic(err)
	}

	db, err := database.ConnectPostgres(dbConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		panic(err)
	}

	return infra.App{Router: mux.NewRouter(), DB: db}
}
