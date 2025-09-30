package repository

import (
	"exchange-crypto-service-api/internal/app/account/domain"

	"gorm.io/gorm"
)

type Account struct {
	ID         uint    `gorm:"primaryKey"`
	UserID     uint    `gorm:"not null"`
	ExchangeID uint    `gorm:"not null"`
	Balance    float64 `gorm:"not null;default:0"`
	gorm.Model
}

func (Account) TableName() string {
	return "accounts"
}

func fromDomain(account domain.Account) Account {
	return Account{
		ID:         account.ID,
		UserID:     account.UserID,
		ExchangeID: account.ExchangeID,
		Balance:    account.Balance,
	}
}

func (a Account) toDomain() domain.Account {
	return domain.Account{
		ID:         a.ID,
		UserID:     a.UserID,
		ExchangeID: a.ExchangeID,
		Balance:    a.Balance,
	}
}
