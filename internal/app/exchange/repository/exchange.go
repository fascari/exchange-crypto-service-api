package repository

import (
	"exchange-crypto-service-api/internal/app/exchange/domain"

	"gorm.io/gorm"
)

type exchangeModel struct {
	ID                uint    `gorm:"primaryKey"`
	Name              string  `gorm:"not null"`
	MinimumAge        uint    `gorm:"not null"`
	MaxTransferAmount float64 `gorm:"column:maximum_transfer_amount;not null"`
	gorm.Model
}

func (exchangeModel) TableName() string {
	return "exchanges"
}

func (e exchangeModel) ToDomain() domain.Exchange {
	return domain.Exchange{
		ID:                e.ID,
		Name:              e.Name,
		MinimumAge:        e.MinimumAge,
		MaxTransferAmount: e.MaxTransferAmount,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}
