package modules

import (
	"exchange-crypto-service-api/cmd/api/modules/dependency"
	"exchange-crypto-service-api/internal/app/user/handler/createuser"
	createuseruc "exchange-crypto-service-api/internal/app/user/usecase/createuser"
	"exchange-crypto-service-api/internal/infra"
)

func User(app infra.App, deps dependency.Dependencies) {
	useCase := createuseruc.New(deps.Repositories.User)
	createuser.RegisterEndpoint(app.Router, createuser.NewHandler(useCase))
}
