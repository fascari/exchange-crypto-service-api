package createtransaction

import (
	"context"
	"errors"
	"fmt"

	accountdomain "exchange-crypto-service-api/internal/app/account/domain"
	exchangedomain "exchange-crypto-service-api/internal/app/exchange/domain"
	"exchange-crypto-service-api/internal/app/transaction/domain"

	"github.com/google/uuid"
)

type (
	AccountRepository interface {
		FindByID(ctx context.Context, accountID uint) (accountdomain.Account, error)
		Update(ctx context.Context, account accountdomain.Account) error
		ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	}

	ExchangeRepository interface {
		FindByID(ctx context.Context, exchangeID uint) (exchangedomain.Exchange, error)
	}

	TransactionRepository interface {
		Create(ctx context.Context, transaction domain.Transaction) error
		CheckIdempotency(ctx context.Context, accountID uint, idempotencyKey string) error
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

func (uc UseCase) Execute(ctx context.Context, transactionType domain.TransactionType, accountID uint, amount float64, idempotencyKey string) (string, error) {
	if transactionType != domain.Deposit && transactionType != domain.Withdrawal {
		return "", fmt.Errorf("invalid transaction type %s", transactionType)
	}

	return uc.processTransaction(ctx, accountID, amount, transactionType, idempotencyKey, validations(transactionType)...)
}

func (uc UseCase) processTransaction(
	ctx context.Context,
	accountID uint,
	amount float64,
	transactionType domain.TransactionType,
	idempotencyKey string,
	opts ...TransactionOption,
) (string, error) {
	if err := validateAmount(amount); err != nil {
		return "", err
	}

	if err := validateIdempotencyKey(idempotencyKey); err != nil {
		return "", err
	}

	cfg := &transactionConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	transactionID := generateTransactionID()

	err := uc.accountRepository.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		return uc.executeTransaction(ctx, accountID, amount, transactionType, transactionID, idempotencyKey, cfg)
	})

	return transactionID, err
}

func (uc UseCase) executeTransaction(
	ctx context.Context,
	accountID uint,
	amount float64,
	transactionType domain.TransactionType,
	transactionID string,
	idempotencyKey string,
	cfg *transactionConfig,
) error {
	if err := uc.checkIdempotency(ctx, accountID, idempotencyKey); err != nil {
		return err
	}

	account, err := uc.accountRepository.FindByID(ctx, accountID)
	if err != nil {
		return err
	}

	if err := uc.validateTransaction(ctx, account, amount, cfg); err != nil {
		return err
	}

	if err := uc.updateAccountBalance(ctx, &account, amount, transactionType); err != nil {
		return err
	}

	previousBalance := account.Balance - calcBalance(account.Balance, -amount, transactionType)

	return uc.transactionRepository.Create(ctx, domain.Transaction{
		AccountID:       accountID,
		Type:            transactionType,
		Amount:          amount,
		PreviousBalance: previousBalance,
		NewBalance:      account.Balance,
		TransactionID:   transactionID,
		IdempotencyKey:  idempotencyKey,
	})
}

func (uc UseCase) checkIdempotency(ctx context.Context, accountID uint, idempotencyKey string) error {
	if idempotencyKey == "" {
		return nil
	}
	return uc.transactionRepository.CheckIdempotency(ctx, accountID, idempotencyKey)
}

func (uc UseCase) validateTransaction(ctx context.Context, account accountdomain.Account, amount float64, cfg *transactionConfig) error {
	if cfg.validator != nil {
		if err := cfg.validator(account, amount); err != nil {
			return err
		}
	}
	return uc.validateTransferLimit(ctx, account.ExchangeID, amount)
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

func (uc UseCase) updateAccountBalance(ctx context.Context, account *accountdomain.Account, amount float64, transactionType domain.TransactionType) error {
	account.Balance = calcBalance(account.Balance, amount, transactionType)
	return uc.accountRepository.Update(ctx, *account)
}

func generateTransactionID() string {
	return uuid.Must(uuid.NewV7()).String()
}

func validateIdempotencyKey(idempotencyKey string) error {
	if idempotencyKey == "" {
		return nil
	}

	if _, err := uuid.Parse(idempotencyKey); err != nil {
		return errors.New("invalid idempotencyKey format: must be a valid UUID")
	}

	return nil
}

func validateAmount(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

func calcBalance(currentBalance, amount float64, transactionType domain.TransactionType) float64 {
	if transactionType == domain.Deposit {
		return currentBalance + amount
	}
	return currentBalance - amount
}

func validations(transactionType domain.TransactionType) []TransactionOption {
	if transactionType == domain.Withdrawal {
		return []TransactionOption{WithBalanceValidation()}
	}
	return nil
}
