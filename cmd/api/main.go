package main

import (
	"fmt"
	"log"
	"net/http"

	"exchange-crypto-service-api/internal/config"
	"exchange-crypto-service-api/internal/database"

	"github.com/gorilla/mux"
)

func main() {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatal("Failed to load database config:", err)
	}

	_, err = database.ConnectPostgres(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := mux.NewRouter()

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
