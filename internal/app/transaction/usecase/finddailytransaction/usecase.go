package finddailytransaction

import (
	"context"
	"errors"
	"time"

	"exchange-crypto-service-api/internal/app/transaction/domain"
)

type (
	Repository interface {
		FindDailyTransactions(ctx context.Context, startDate, endDate time.Time) ([]domain.DailyTransaction, error)
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

func (uc UseCase) Execute(ctx context.Context, startDate, endDate time.Time) ([]domain.DailyTransaction, error) {
	if err := validateDateRange(startDate, endDate); err != nil {
		return nil, err
	}

	return uc.repository.FindDailyTransactions(ctx, startDate, endDate)
}

func validateDateRange(startDate, endDate time.Time) error {
	if startDate.IsZero() {
		return errors.New("start date is required")
	}

	if endDate.IsZero() {
		return errors.New("end date is required")
	}

	if startDate.After(endDate) {
		return errors.New("start date must be before or equal to end date")
	}

	return nil
}
