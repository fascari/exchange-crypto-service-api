package createtransaction

import (
	"context"

	"exchange-crypto-service-api/internal/app/transaction/domain"
)

func (uc UseCase) Deposit(ctx context.Context, accountID uint, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	account, err := uc.accountRepository.FindByID(ctx, accountID)
	if err != nil {
		return err
	}

	if err := uc.validateTransferLimit(ctx, account.ExchangeID, amount); err != nil {
		return err
	}

	if err := uc.updateBalance(ctx, account, amount, domain.Deposit); err != nil {
		return err
	}

	return uc.createTransaction(ctx, accountID, amount, domain.Deposit)
}
