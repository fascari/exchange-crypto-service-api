package modules

import (
	"log"

	"exchange-crypto-service-api/internal/config"
	"exchange-crypto-service-api/internal/database"
	"exchange-crypto-service-api/internal/infra"

	"github.com/gorilla/mux"
)

func NewApp() infra.App {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatal("Failed to load database config:", err)
	}

	db, err := database.ConnectPostgres(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return infra.App{Router: mux.NewRouter(), DB: db}
}
