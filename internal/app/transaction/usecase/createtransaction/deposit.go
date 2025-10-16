package createtransaction

import (
	"context"

	"exchange-crypto-service-api/internal/app/transaction/domain"
)

func (uc UseCase) Deposit(ctx context.Context, accountID uint, amount float64, idempotencyKey string) (string, error) {
	if err := validateAmount(amount); err != nil {
		return "", err
	}

	if err := validateIdempotencyKey(idempotencyKey); err != nil {
		return "", err
	}

	transactionID := generateTransactionID()

	err := uc.accountRepository.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		account, err := uc.accountRepository.FindByID(ctx, accountID)
		if err != nil {
			return err
		}

		if err := uc.validateTransferLimit(ctx, account.ExchangeID, amount); err != nil {
			return err
		}

		previousBalance := account.Balance
		newBalance := calcBalance(account.Balance, amount, domain.Deposit)
		account.Balance = newBalance

		if err := uc.accountRepository.Update(ctx, account); err != nil {
			return err
		}

		transaction := domain.Transaction{
			AccountID:       accountID,
			Type:            domain.Deposit,
			Amount:          amount,
			PreviousBalance: previousBalance,
			NewBalance:      newBalance,
			TransactionID:   transactionID,
			IdempotencyKey:  idempotencyKey,
		}

		return uc.transactionRepository.Create(ctx, transaction)
	})

	return transactionID, err
}
