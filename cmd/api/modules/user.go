package modules

import (
	"exchange-crypto-service-api/internal/app/user/handler/createuser"
	"exchange-crypto-service-api/internal/app/user/handler/finduserbalance"
	userrepo "exchange-crypto-service-api/internal/app/user/repository"
	createuseruc "exchange-crypto-service-api/internal/app/user/usecase/createuser"
	finduserbalanceuc "exchange-crypto-service-api/internal/app/user/usecase/finduserbalance"

	"go.uber.org/fx"
)

var userFactories = fx.Provide(
	// repositories
	userrepo.New,

	// use cases
	createuseruc.New,
	finduserbalanceuc.New,

	// handlers
	createuser.NewHandler,
	finduserbalance.NewHandler,
)

var userDependencies = fx.Provide(
	// repositories
	func(repo userrepo.Repository) createuseruc.Repository {
		return repo
	},
	func(repo userrepo.Repository) finduserbalanceuc.Repository {
		return repo
	},
)

var userInvokes = fx.Invoke(
	func(params RouterParams, h createuser.Handler) {
		createuser.RegisterEndpoint(params.APIRouter, h)
	},
	func(params RouterParams, h finduserbalance.Handler) {
		finduserbalance.RegisterEndpoint(params.APIRouter, h)
	},
)

var User = fx.Options(
	userFactories,
	userDependencies,
	userInvokes,
)
