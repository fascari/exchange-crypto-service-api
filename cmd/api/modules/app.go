package modules

import (
	"exchange-crypto-service-api/internal/config"
	"exchange-crypto-service-api/internal/database"
	"exchange-crypto-service-api/internal/infra"
	"exchange-crypto-service-api/internal/jwt"
	"exchange-crypto-service-api/pkg/logger"
	"exchange-crypto-service-api/pkg/telemetry"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/gorm"
)

const serviceName = "exchange-crypto-service-api"

type RouterSetup struct {
	MainRouter *mux.Router
	Router     *mux.Router
}

func NewApp() infra.App {
	logger.Init()
	jwt.Initialize()

	routerSetup := setupRouter()

	return infra.App{
		DB:             connectDB(),
		TracerProvider: initTelemetry(),
		MainRouter:     routerSetup.MainRouter,
		Router:         routerSetup.Router,
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

func connectDB() *gorm.DB {
	db, err := database.ConnectPostgres(config.LoadDatabase())
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		panic(err)
	}
	return db
}

func setupRouter() RouterSetup {
	router := mux.NewRouter()
	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	setupMiddlewares(apiV1Router, serviceName)

	return RouterSetup{
		MainRouter: router,
		Router:     apiV1Router,
	}
}
