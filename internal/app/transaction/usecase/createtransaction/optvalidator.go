package createtransaction

import (
	"fmt"

	accountdomain "exchange-crypto-service-api/internal/app/account/domain"
)

type (
	TransactionOption func(*transactionConfig)

	transactionConfig struct {
		validator func(account accountdomain.Account, amount float64) error
	}
)

func WithBalanceValidation() TransactionOption {
	return func(cfg *transactionConfig) {
		cfg.validator = func(account accountdomain.Account, amount float64) error {
			if account.Balance < amount {
				return fmt.Errorf("withdrawal amount %.2f exceeds account balance of %.2f", amount, account.Balance)
			}
			return nil
		}
	}
}
