package modules

import (
	"exchange-crypto-service-api/internal/config"
	"exchange-crypto-service-api/internal/database"
	"exchange-crypto-service-api/internal/infra"
	"exchange-crypto-service-api/internal/middleware"
	"exchange-crypto-service-api/pkg/logger"
	"exchange-crypto-service-api/pkg/telemetry"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

const serviceName = "exchange-crypto-service-api"

func NewApp() infra.App {
	logger.Init()

	tp, err := telemetry.InitTracer(serviceName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize tracer")
		panic(err)
	}

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

	router := mux.NewRouter()

	middleware.SetupOTEL(router, serviceName)

	return infra.App{
		Router:         router,
		DB:             db,
		TracerProvider: tp,
	}
}
