package modules

import (
	"exchange-crypto-service-api/internal/app/transaction/handler/createtransaction"
	createtrantuc "exchange-crypto-service-api/internal/app/transaction/usecase/createtransaction"
	"exchange-crypto-service-api/internal/deps"
	"exchange-crypto-service-api/internal/infra"
)

func Transaction(app infra.App, deps deps.Dependencies) {
	useCase := createtrantuc.New(deps.Repositories.Account, deps.Repositories.Exchange, deps.Repositories.Transaction)
	createtransaction.RegisterEndpoint(app.Router, createtransaction.NewHandler(useCase))
}
