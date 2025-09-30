package createaccount

import (
	"context"
	"time"

	accountpkg "exchange-crypto-service-api/internal/app/account"
	"exchange-crypto-service-api/internal/app/account/domain"
	exchangedomain "exchange-crypto-service-api/internal/app/exchange/domain"
	userdomain "exchange-crypto-service-api/internal/app/user/domain"
)

//go:generate mockery --all

type (
	Repository interface {
		Create(ctx context.Context, account domain.Account) (domain.Account, error)
	}

	UserRepository interface {
		FindByID(ctx context.Context, id uint) (userdomain.User, error)
	}

	ExchangeRepository interface {
		FindByID(ctx context.Context, id uint) (exchangedomain.Exchange, error)
	}

	UseCase struct {
		repository         Repository
		userRepository     UserRepository
		exchangeRepository ExchangeRepository
	}
)

func New(repository Repository, userRepository UserRepository, exchangeRepository ExchangeRepository) UseCase {
	return UseCase{
		repository:         repository,
		userRepository:     userRepository,
		exchangeRepository: exchangeRepository,
	}
}

func (u UseCase) Execute(ctx context.Context, account domain.Account) (domain.Account, error) {
	user, err := u.userRepository.FindByID(ctx, account.UserID)
	if err != nil {
		return domain.Account{}, err
	}

	exchange, err := u.exchangeRepository.FindByID(ctx, account.ExchangeID)
	if err != nil {
		return domain.Account{}, err
	}

	userAge := calculateAge(user.DateOfBirth)
	if userAge < exchange.MinimumAge {
		return domain.Account{}, accountpkg.NewErrInvalidMinimumAge(exchange.MinimumAge, userAge)
	}

	return u.repository.Create(ctx, account)
}

func calculateAge(birthDate time.Time) uint {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	if age < 0 {
		return 0
	}
	return uint(age)
}
