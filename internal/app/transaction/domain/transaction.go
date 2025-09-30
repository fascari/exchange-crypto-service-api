package domain

import (
	"fmt"
	"strings"
)

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
)

type (
	TransactionType string

	Transaction struct {
		ID        uint
		AccountID uint
		Type      TransactionType
		Amount    float64
	}
)

func ParseTransactionType(value string) (TransactionType, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	transactionType := TransactionType(normalized)

	if transactionType == Deposit || transactionType == Withdrawal {
		return transactionType, nil
	}

	return "", fmt.Errorf("invalid transaction type: %s. Accepted values: DEPOSIT, WITHDRAWAL", value)
}
