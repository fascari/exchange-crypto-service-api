package modules

import (
	"exchange-crypto-service-api/internal/app/auth/handler/tokengen"
	"exchange-crypto-service-api/internal/app/auth/usecase/generatetoken"
	"exchange-crypto-service-api/internal/deps"
	"exchange-crypto-service-api/internal/infra"
)

func Auth(app infra.App, dep deps.Dependencies) {
	useCase := generatetoken.New(dep.Repositories.User)
	tokengen.RegisterEndpoint(app.MainRouter, tokengen.NewHandler(useCase))
}
