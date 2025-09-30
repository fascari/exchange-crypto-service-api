package modules

import (
	"exchange-crypto-service-api/internal/app/account/handler/createaccount"
	createaccountuc "exchange-crypto-service-api/internal/app/account/usecase/createaccount"
	"exchange-crypto-service-api/internal/deps"
	"exchange-crypto-service-api/internal/infra"
)

func Account(app infra.App, deps deps.Dependencies) {
	useCase := createaccountuc.New(deps.Repositories.Account, deps.Repositories.User, deps.Repositories.Exchange)
	createaccount.RegisterEndpoint(app.Router, createaccount.NewHandler(useCase))
}
