package modules

import (
	"exchange-crypto-service-api/internal/app/account/handler/createaccount"
	accountrepo "exchange-crypto-service-api/internal/app/account/repository"
	createaccountuc "exchange-crypto-service-api/internal/app/account/usecase/createaccount"
	exchangerepo "exchange-crypto-service-api/internal/app/exchange/repository"
	userrepo "exchange-crypto-service-api/internal/app/user/repository"

	"go.uber.org/fx"
)

var accountFactories = fx.Provide(
	// repositories
	accountrepo.New,

	// use cases
	createaccountuc.New,

	// handlers
	createaccount.NewHandler,
)

var accountDependencies = fx.Provide(
	// repositories
	func(repo accountrepo.Repository) createaccountuc.Repository {
		return repo
	},
	func(repo userrepo.Repository) createaccountuc.UserRepository {
		return repo
	},
	func(repo exchangerepo.Repository) createaccountuc.ExchangeRepository {
		return repo
	},
)

var accountInvokes = fx.Invoke(
	func(params RouterParams, h createaccount.Handler) {
		createaccount.RegisterEndpoint(params.APIRouter, h)
	},
)

var Account = fx.Options(
	accountFactories,
	accountDependencies,
	accountInvokes,
)
