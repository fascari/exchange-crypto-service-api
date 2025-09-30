package modules

import (
	"exchange-crypto-service-api/internal/app/transaction/handler/createtransaction"
	createtrantuc "exchange-crypto-service-api/internal/app/transaction/usecase/createtransaction"
	"exchange-crypto-service-api/internal/deps"
	"exchange-crypto-service-api/internal/infra"
)

func Transaction(app infra.App, dep deps.Dependencies) {
	useCase := createtrantuc.New(dep.Repositories.Account, dep.Repositories.Exchange, dep.Repositories.Transaction)
	createtransaction.RegisterEndpoint(app.Router, createtransaction.NewHandler(useCase))
}
