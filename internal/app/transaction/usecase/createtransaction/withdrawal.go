package createtransaction

import (
	"context"
	"fmt"

	accountdomain "exchange-crypto-service-api/internal/app/account/domain"
	"exchange-crypto-service-api/internal/app/transaction/domain"
)

func (uc UseCase) Withdrawal(ctx context.Context, accountID uint, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	account, err := uc.accountRepository.FindByID(ctx, accountID)
	if err != nil {
		return err
	}

	if err := validateBalance(account, amount); err != nil {
		return err
	}

	if err := uc.validateTransferLimit(ctx, account.ExchangeID, amount); err != nil {
		return err
	}

	if err := uc.updateBalance(ctx, account, amount, domain.Withdrawal); err != nil {
		return err
	}

	return uc.createTransaction(ctx, accountID, amount, domain.Withdrawal)
}

func validateBalance(account accountdomain.Account, amount float64) error {
	if account.Balance < amount {
		return fmt.Errorf("withdrawal amount %.2f exceeds account balance of %.2f", amount, account.Balance)
	}
	return nil
}
