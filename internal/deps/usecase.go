package deps

import (
	generatetokenusecase "exchange-crypto-service-api/internal/app/auth/usecase/generatetoken"
)

type UseCases struct {
	GenerateToken generatetokenusecase.UseCase
}

func initUseCases(repos Repositories) UseCases {
	return UseCases{
		GenerateToken: generatetokenusecase.New(repos.User),
	}
}
