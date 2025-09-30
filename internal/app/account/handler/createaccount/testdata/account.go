package testdata

import (
	"time"

	"exchange-crypto-service-api/internal/app/account/domain"
	"exchange-crypto-service-api/internal/app/account/handler/createaccount"
)

func ValidAccount() domain.Account {
	return domain.Account{
		ID:         1,
		UserID:     1,
		ExchangeID: 1,
		Balance:    100.0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func AccountWithBalance(balance float64) domain.Account {
	account := ValidAccount()
	account.Balance = balance
	return account
}

func ZeroBalanceAccount() domain.Account {
	return AccountWithBalance(0.0)
}

func ValidInputPayload() createaccount.InputPayload {
	return createaccount.InputPayload{
		UserID:     1,
		ExchangeID: 1,
		Balance:    100.0,
	}
}

func InputPayloadMissingUserID() createaccount.InputPayload {
	return createaccount.InputPayload{
		ExchangeID: 1,
		Balance:    100.0,
	}
}
