package modules

import (
	accountrepo "exchange-crypto-service-api/internal/app/account/repository"
	exchangerepo "exchange-crypto-service-api/internal/app/exchange/repository"
	"exchange-crypto-service-api/internal/app/transaction/handler/createtransaction"
	"exchange-crypto-service-api/internal/app/transaction/handler/finddailytransaction"
	transrepo "exchange-crypto-service-api/internal/app/transaction/repository"
	createtrantuc "exchange-crypto-service-api/internal/app/transaction/usecase/createtransaction"
	finddailytransactionuc "exchange-crypto-service-api/internal/app/transaction/usecase/finddailytransaction"

	"go.uber.org/fx"
)

var transactionFactories = fx.Provide(
	// repositories
	transrepo.New,

	// use cases
	createtrantuc.New,
	finddailytransactionuc.New,

	// handlers
	createtransaction.NewHandler,
	finddailytransaction.NewHandler,
)

var transactionDependencies = fx.Provide(
	// repositories
	func(repo accountrepo.Repository) createtrantuc.AccountRepository {
		return repo
	},
	func(repo exchangerepo.Repository) createtrantuc.ExchangeRepository {
		return repo
	},
	func(repo transrepo.Repository) createtrantuc.TransactionRepository {
		return repo
	},
	func(repo transrepo.Repository) finddailytransactionuc.Repository {
		return repo
	},
)

var transactionInvokes = fx.Invoke(
	func(params RouterParams, h createtransaction.Handler) {
		createtransaction.RegisterEndpoint(params.APIRouter, h)
	},
	func(params RouterParams, h finddailytransaction.Handler) {
		finddailytransaction.RegisterEndpoint(params.APIRouter, h)
	},
)

var Transaction = fx.Options(
	transactionFactories,
	transactionDependencies,
	transactionInvokes,
)
