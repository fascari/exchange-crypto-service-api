package createtransaction

import (
	"context"
	"errors"
	"fmt"

	accountdomain "exchange-crypto-service-api/internal/app/account/domain"
	exchangedomain "exchange-crypto-service-api/internal/app/exchange/domain"
	"exchange-crypto-service-api/internal/app/transaction/domain"
)

type (
	AccountRepository interface {
		FindByID(ctx context.Context, accountID uint) (accountdomain.Account, error)
		Update(ctx context.Context, account accountdomain.Account) error
	}

	ExchangeRepository interface {
		FindByID(ctx context.Context, exchangeID uint) (exchangedomain.Exchange, error)
	}

	TransactionRepository interface {
		Create(ctx context.Context, transaction domain.Transaction) error
	}

	UseCase struct {
		accountRepository     AccountRepository
		exchangeRepository    ExchangeRepository
		transactionRepository TransactionRepository
	}
)

func New(accountRepository AccountRepository, exchangeRepository ExchangeRepository, transactionRepository TransactionRepository) UseCase {
	return UseCase{
		accountRepository:     accountRepository,
		exchangeRepository:    exchangeRepository,
		transactionRepository: transactionRepository,
	}
}

func (uc UseCase) Execute(ctx context.Context, transactionType domain.TransactionType, accountID uint, amount float64) error {
	if transactionType == domain.Deposit {
		return uc.Deposit(ctx, accountID, amount)
	}

	if transactionType == domain.Withdrawal {
		return uc.Withdrawal(ctx, accountID, amount)
	}

	return errors.New("invalid transaction type")
}

func (uc UseCase) validateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

func (uc UseCase) validateTransferLimit(ctx context.Context, exchangeID uint, amount float64) error {
	exchange, err := uc.exchangeRepository.FindByID(ctx, exchangeID)
	if err != nil {
		return err
	}

	if amount > exchange.MaxTransferAmount {
		return fmt.Errorf("amount %.2f exceeds exchange maximum transfer limit of %.2f", amount, exchange.MaxTransferAmount)
	}
	return nil
}

func (uc UseCase) updateBalance(ctx context.Context, account accountdomain.Account, amount float64, transactionType domain.TransactionType) error {
	account.Balance = calcBalance(account.Balance, amount, transactionType)
	return uc.accountRepository.Update(ctx, account)
}

func (uc UseCase) createTransaction(ctx context.Context, accountID uint, amount float64, transactionType domain.TransactionType) error {
	transaction := domain.Transaction{
		AccountID: accountID,
		Type:      transactionType,
		Amount:    amount,
	}

	return uc.transactionRepository.Create(ctx, transaction)
}

func calcBalance(currentBalance, amount float64, transactionType domain.TransactionType) float64 {
	if transactionType == domain.Deposit {
		return currentBalance + amount
	}
	return currentBalance - amount
}
