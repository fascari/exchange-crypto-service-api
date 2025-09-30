package finddailytransaction

import (
	"slices"
	"time"

	"exchange-crypto-service-api/internal/app/transaction/domain"
)

type (
	OutputPayload struct {
		Exchange    string      `json:"exchange"`
		DailyAmout  []DailyData `json:"daily_amout"`
		TotalAmount float64     `json:"total_amount"`
	}

	DailyData struct {
		Date   string  `json:"date"`
		Amount float64 `json:"amount"`
	}

	exchangeGroup map[string]exchangeData

	exchangeData struct {
		DailyData   []DailyData
		TotalAmount float64
	}
)

func ToOutputPayloads(transactions []domain.DailyTransaction) []OutputPayload {
	return convertToOutputPayloads(groupTransactionsByExchange(transactions))
}

func groupTransactionsByExchange(transactions []domain.DailyTransaction) exchangeGroup {
	group := make(exchangeGroup)

	for _, transaction := range transactions {
		addTransactionToExchange(group, transaction)
	}

	return group
}

func addTransactionToExchange(exchangeMap exchangeGroup, transaction domain.DailyTransaction) {
	exchange := transaction.Exchange
	data := exchangeMap[exchange]
	data.DailyData = append(data.DailyData, DailyData{
		Date:   transaction.Date.Format(time.DateOnly),
		Amount: transaction.TotalAmount,
	})
	data.TotalAmount += transaction.TotalAmount
	exchangeMap[exchange] = data
}

func convertToOutputPayloads(group exchangeGroup) []OutputPayload {
	exchanges := make([]string, 0, len(group))
	for exchange := range group {
		exchanges = append(exchanges, exchange)
	}
	slices.Sort(exchanges)

	payloads := make([]OutputPayload, 0, len(group))
	for _, exchange := range exchanges {
		data := group[exchange]
		payloads = append(payloads, OutputPayload{
			Exchange:    exchange,
			DailyAmout:  data.DailyData,
			TotalAmount: data.TotalAmount,
		})
	}

	return payloads
}
