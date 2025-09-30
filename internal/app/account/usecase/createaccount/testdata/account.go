package testdata

import (
	"exchange-crypto-service-api/internal/app/account/domain"
)

func ValidAccount() domain.Account {
	return domain.Account{
		ID:         1,
		UserID:     1,
		ExchangeID: 1,
		Balance:    100.0,
	}
}

func AccountWithID(id uint) domain.Account {
	account := ValidAccount()
	account.ID = id
	return account
}

func InputAccount() domain.Account {
	return domain.Account{
		UserID:     1,
		ExchangeID: 1,
		Balance:    100.0,
	}
}

func InputAccountWithUserAndExchange(userID, exchangeID uint) domain.Account {
	account := InputAccount()
	account.UserID = userID
	account.ExchangeID = exchangeID
	return account
}
