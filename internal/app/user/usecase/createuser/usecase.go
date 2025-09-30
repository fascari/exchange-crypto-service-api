package createuser

import (
	"context"

	"exchange-crypto-service-api/internal/app/user/domain"
)

type (
	Repository interface {
		Create(ctx context.Context, user domain.User) (domain.User, error)
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

func (u UseCase) Execute(ctx context.Context, user domain.User) (domain.User, error) {
	return u.repository.Create(ctx, user)
}
