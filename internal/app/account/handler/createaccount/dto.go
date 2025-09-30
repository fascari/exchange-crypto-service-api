package createaccount

import (
	"exchange-crypto-service-api/internal/app/account/domain"
	"exchange-crypto-service-api/pkg/validator"
)

type (
	InputPayload struct {
		UserID     uint    `json:"user_id" validate:"required"`
		ExchangeID uint    `json:"exchange_id" validate:"required"`
		Balance    float64 `json:"balance" validate:"min=0"`
	}

	OutputPayload struct {
		ID         uint    `json:"id"`
		UserID     uint    `json:"user_id"`
		ExchangeID uint    `json:"exchange_id"`
		Balance    float64 `json:"balance"`
	}
)

func (p InputPayload) Validate() error {
	return validator.Validate(p)
}

func (p InputPayload) ToDomain() domain.Account {
	return domain.Account{
		UserID:     p.UserID,
		ExchangeID: p.ExchangeID,
		Balance:    p.Balance,
	}
}

func ToOutputPayload(account domain.Account) OutputPayload {
	return OutputPayload{
		ID:         account.ID,
		UserID:     account.UserID,
		ExchangeID: account.ExchangeID,
		Balance:    account.Balance,
	}
}
