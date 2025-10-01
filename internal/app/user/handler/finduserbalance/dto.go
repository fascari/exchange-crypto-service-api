package finduserbalance

import (
	"exchange-crypto-service-api/internal/app/user/domain"
)

type (
	OutputPayload struct {
		UserID           uint                    `json:"user_id,omitempty"`
		Username         string                  `json:"username,omitempty"`
		ExchangeBalances []ExchangeBalanceOutput `json:"exchange_balances,omitempty"`
		TotalBalance     float64                 `json:"total_balance,omitempty"`
	}

	ExchangeBalanceOutput struct {
		ExchangeName string  `json:"exchange_name,omitempty"`
		Balance      float64 `json:"balance,omitempty"`
	}
)

func ToOutputPayload(ub domain.UserBalance) any {
	return OutputPayload{
		UserID:           ub.UserID,
		Username:         ub.Username,
		ExchangeBalances: toExchangeOutputs(ub.ExchangeBalances),
		TotalBalance:     ub.TotalBalance,
	}
}

func toExchangeOutputs(balances []domain.ExchangeBalance) []ExchangeBalanceOutput {
	result := make([]ExchangeBalanceOutput, len(balances))
	for i, eb := range balances {
		result[i] = ExchangeBalanceOutput{
			ExchangeName: eb.ExchangeName,
			Balance:      eb.Balance,
		}
	}
	return result
}
