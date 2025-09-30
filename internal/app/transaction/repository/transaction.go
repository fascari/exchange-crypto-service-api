package repository

import (
	"time"

	"exchange-crypto-service-api/internal/app/transaction/domain"

	"gorm.io/gorm"
)

type transactionModel struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	AccountID uint      `gorm:"not null;column:account_id"`
	Type      string    `gorm:"not null;column:type"`
	Amount    float64   `gorm:"not null;column:amount"`
	CreatedAt time.Time `gorm:"not null;column:created_at;autoCreateTime"`
	DeletedAt gorm.DeletedAt
}

func (transactionModel) TableName() string {
	return "transactions"
}

func (t transactionModel) toDomain() domain.Transaction {
	return domain.Transaction{
		ID:        t.ID,
		AccountID: t.AccountID,
		Amount:    t.Amount,
		Type:      domain.TransactionType(t.Type),
	}
}

func fromDomain(transaction domain.Transaction) transactionModel {
	return transactionModel{
		ID:        transaction.ID,
		AccountID: transaction.AccountID,
		Amount:    transaction.Amount,
		Type:      string(transaction.Type),
	}
}
