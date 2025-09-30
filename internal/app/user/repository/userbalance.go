package repository

import (
	"time"

	"exchange-crypto-service-api/internal/app/user/domain"
)

type userExchangeBalanceModel struct {
	UserID       uint      `gorm:"column:user_id"`
	Username     string    `gorm:"column:username"`
	ExchangeName string    `gorm:"column:exchange_name"`
	Balance      float64   `gorm:"column:balance"`
	TotalBalance float64   `gorm:"column:total_balance"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (userExchangeBalanceModel) TableName() string {
	return "user_exchange_balances"
}

func (u userExchangeBalanceModel) toDomain() domain.ExchangeBalance {
	return domain.ExchangeBalance{
		ExchangeName: u.ExchangeName,
		Balance:      u.Balance,
	}
}

func toDomain(models []userExchangeBalanceModel) domain.UserBalance {
	exchangeBalances := make([]domain.ExchangeBalance, 0, len(models))
	for _, model := range models {
		exchangeBalances = append(exchangeBalances, model.toDomain())
	}

	balanceModel := models[0]
	return domain.UserBalance{
		UserID:           balanceModel.UserID,
		Username:         balanceModel.Username,
		ExchangeBalances: exchangeBalances,
		TotalBalance:     balanceModel.TotalBalance,
	}
}
