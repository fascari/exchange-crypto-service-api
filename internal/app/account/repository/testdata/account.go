package testdata

import "exchange-crypto-service-api/internal/app/account/domain"

func ValidAccount() domain.Account {
	return createAccount(1, 1, 100.0)
}

func ZeroBalanceAccount() domain.Account {
	return createAccount(2, 1, 0.0)
}

func HighBalanceAccount() domain.Account {
	return createAccount(3, 2, 999999.99)
}

func AccountForUpdate() domain.Account {
	return createAccount(1, 1, 100.0)
}

func UpdatedAccount(id uint, newBalance float64) domain.Account {
	account := createAccount(1, 1, newBalance)
	account.ID = id
	return account
}

func InvalidAccount() domain.Account {
	return createAccount(0, 0, -100.0)
}

func createAccount(userID, exchangeID uint, balance float64) domain.Account {
	return domain.Account{
		UserID:     userID,
		ExchangeID: exchangeID,
		Balance:    balance,
	}
}
