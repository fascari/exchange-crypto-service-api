package testdata

import (
	"time"

	exchangedomain "exchange-crypto-service-api/internal/app/exchange/domain"
)

func ValidExchange() exchangedomain.Exchange {
	return exchangedomain.Exchange{
		ID:         1,
		Name:       "Test Exchange",
		MinimumAge: 18,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
