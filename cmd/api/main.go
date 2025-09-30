package main

import (
	"exchange-crypto-service-api/cmd/api/modules"
	"exchange-crypto-service-api/internal/deps"
	"exchange-crypto-service-api/internal/infra"
)

func main() {
	app := setupApp()
	app.Start()
}

func setupApp() infra.App {
	app := modules.NewApp()
	dependencies := deps.New(app)

	modules.User(app, dependencies)
	modules.Account(app, dependencies)
	modules.Transaction(app, dependencies)

	return app
}
