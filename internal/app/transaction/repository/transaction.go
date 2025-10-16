package repository

import (
	"time"

	"exchange-crypto-service-api/internal/app/transaction/domain"

	"gorm.io/gorm"
)

type transactionModel struct {
	ID              uint      `gorm:"primaryKey;column:id"`
	AccountID       uint      `gorm:"not null;column:account_id"`
	Type            string    `gorm:"not null;column:type"`
	Amount          float64   `gorm:"not null;column:amount"`
	PreviousBalance float64   `gorm:"column:previous_balance"`
	NewBalance      float64   `gorm:"column:new_balance"`
	TransactionID   string    `gorm:"column:transaction_id;uniqueIndex"`
	IdempotencyKey  string    `gorm:"column:idempotency_key;index"`
	CreatedAt       time.Time `gorm:"not null;column:created_at;autoCreateTime"`
	DeletedAt       gorm.DeletedAt
}

func (transactionModel) TableName() string {
	return "transactions"
}

func fromDomain(transaction domain.Transaction) transactionModel {
	return transactionModel{
		ID:              transaction.ID,
		AccountID:       transaction.AccountID,
		Amount:          transaction.Amount,
		Type:            string(transaction.Type),
		PreviousBalance: transaction.PreviousBalance,
		NewBalance:      transaction.NewBalance,
		TransactionID:   transaction.TransactionID,
		IdempotencyKey:  transaction.IdempotencyKey,
	}
}
