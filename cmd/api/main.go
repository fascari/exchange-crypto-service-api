package main

import (
	"fmt"
	"log"
	"net/http"

	"exchange-crypto-service-api/cmd/api/modules"
	"exchange-crypto-service-api/internal/deps"
)

func main() {
	app := modules.NewApp()
	dependencies := deps.New(app)

	modules.User(app, dependencies)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
