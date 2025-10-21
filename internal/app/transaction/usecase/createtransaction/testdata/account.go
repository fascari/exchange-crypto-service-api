package testdata

import "exchange-crypto-service-api/internal/app/account/domain"

func AccountWithBalance() domain.Account {
	return domain.Account{
		ID:      uint(1),
		UserID:  1,
		Balance: 1000.0,
	}
}

func AccountEmpty() domain.Account {
	return domain.Account{
		ID:      uint(2),
		UserID:  2,
		Balance: 0.0,
	}
}

func InvalidAccount() domain.Account {
	return domain.Account{
		ID:      uint(9999),
		UserID:  0,
		Balance: 0.0,
	}
}

func ExpectedBalanceAfterDeposits(initialBalance float64, count int, amount float64) float64 {
	return initialBalance + (float64(count) * amount)
}
