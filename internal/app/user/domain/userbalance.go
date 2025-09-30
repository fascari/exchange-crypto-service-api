package domain

type (
	ExchangeBalance struct {
		ExchangeName string
		Balance      float64
	}

	UserBalance struct {
		UserID           uint
		Username         string
		ExchangeBalances []ExchangeBalance
		TotalBalance     float64
	}
)
