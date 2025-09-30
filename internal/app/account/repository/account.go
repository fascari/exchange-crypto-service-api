package repository

import (
	"exchange-crypto-service-api/internal/app/account/domain"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID     uint `gorm:"not null" json:"user_id"`
	ExchangeID uint `gorm:"not null" json:"exchange_id"`
}

func (Account) TableName() string {
	return "accounts"
}

func fromDomain(account domain.Account) Account {
	return Account{
		UserID:     account.UserID,
		ExchangeID: account.ExchangeID,
	}
}

func (a Account) toDomain() domain.Account {
	return domain.Account{
		ID:         a.ID,
		UserID:     a.UserID,
		ExchangeID: a.ExchangeID,
	}
}
