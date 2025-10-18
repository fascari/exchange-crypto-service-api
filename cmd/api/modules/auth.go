package modules

import (
	"exchange-crypto-service-api/internal/app/auth/handler/tokengen"
	generatetokenusecase "exchange-crypto-service-api/internal/app/auth/usecase/generatetoken"
	userrepo "exchange-crypto-service-api/internal/app/user/repository"

	"go.uber.org/fx"
)

var authFactories = fx.Provide(
	// use cases
	generatetokenusecase.New,

	// handlers
	tokengen.NewHandler,
)

var authDependencies = fx.Provide(
	// repositories
	func(repo userrepo.Repository) generatetokenusecase.UserRepository {
		return repo
	},
)

var authInvokes = fx.Invoke(
	func(params RouterParams, h tokengen.Handler) {
		tokengen.RegisterEndpoint(params.Router, h)
	},
)

var Auth = fx.Options(
	authFactories,
	authDependencies,
	authInvokes,
)
