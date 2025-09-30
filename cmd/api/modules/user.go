package modules

import (
	"exchange-crypto-service-api/internal/app/user/handler/createuser"
	"exchange-crypto-service-api/internal/app/user/handler/finduserbalance"
	createuseruc "exchange-crypto-service-api/internal/app/user/usecase/createuser"
	finduserbalanceuc "exchange-crypto-service-api/internal/app/user/usecase/finduserbalance"
	"exchange-crypto-service-api/internal/deps"
	"exchange-crypto-service-api/internal/infra"
)

func User(app infra.App, dep deps.Dependencies) {
	repository := dep.Repositories.User

	createUserUC := createuseruc.New(repository)
	createuser.RegisterEndpoint(app.Router, createuser.NewHandler(createUserUC))

	findUserBalanceUC := finduserbalanceuc.New(repository)
	finduserbalance.RegisterEndpoint(app.Router, finduserbalance.NewHandler(findUserBalanceUC))
}
