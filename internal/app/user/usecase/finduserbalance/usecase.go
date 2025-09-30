package finduserbalance

import (
	"context"

	"exchange-crypto-service-api/internal/app/user/domain"
)

type (
	Repository interface {
		FindUserBalances(ctx context.Context, userID uint) (domain.UserBalance, error)
	}

	UseCase struct {
		repository Repository
	}
)

func New(repository Repository) UseCase {
	return UseCase{
		repository: repository,
	}
}

func (u UseCase) Execute(ctx context.Context, userID uint) (domain.UserBalance, error) {
	return u.repository.FindUserBalances(ctx, userID)
}
