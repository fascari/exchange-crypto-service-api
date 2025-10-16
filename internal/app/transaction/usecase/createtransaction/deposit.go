package createtransaction

import (
	"context"

	"exchange-crypto-service-api/internal/app/transaction/domain"
)

func (uc UseCase) Deposit(ctx context.Context, accountID uint, amount float64) error {
	if err := validateAmount(amount); err != nil {
		return err
	}

	return uc.accountRepository.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		account, err := uc.accountRepository.FindByID(ctx, accountID)
		if err != nil {
			return err
		}

		if err := uc.validateTransferLimit(ctx, account.ExchangeID, amount); err != nil {
			return err
		}

		account.Balance = calcBalance(account.Balance, amount, domain.Deposit)

		if err := uc.accountRepository.Update(ctx, account); err != nil {
			return err
		}

		return uc.transactionRepository.Create(ctx, domain.Transaction{
			AccountID: accountID,
			Type:      domain.Deposit,
			Amount:    amount,
		})
	})
}
