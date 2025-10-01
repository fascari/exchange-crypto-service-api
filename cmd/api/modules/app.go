package modules

import (
	"exchange-crypto-service-api/internal/config"
	"exchange-crypto-service-api/internal/database"
	"exchange-crypto-service-api/internal/infra"
	"exchange-crypto-service-api/pkg/logger"
	"exchange-crypto-service-api/pkg/telemetry"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

const serviceName = "exchange-crypto-service-api"

func NewApp() infra.App {
	logger.Init()
	tp := initTelemetry()

	return infra.App{
		Router:         setupRouter(),
		DB:             connectDB(),
		TracerProvider: tp,
	}
}

func initTelemetry() *trace.TracerProvider {
	tp, err := telemetry.InitTracer(serviceName)
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize otel tracer")
		panic(err)
	}
	return tp
}

func loadDBConfig() config.Database {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to load database config")
		panic(err)
	}
	return dbConfig
}

func connectDB() *gorm.DB {
	db, err := database.ConnectPostgres(loadDBConfig())
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		panic(err)
	}
	return db
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	setupMiddlewares(router, serviceName)
	return router
}
