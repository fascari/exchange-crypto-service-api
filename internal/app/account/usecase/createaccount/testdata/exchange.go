package testdata

import (
	exchangedomain "exchange-crypto-service-api/internal/app/exchange/domain"
)

func ValidExchange() exchangedomain.Exchange {
	return exchangedomain.Exchange{
		ID:         1,
		Name:       "Test Exchange",
		MinimumAge: 18,
	}
}

func ExchangeWithMinimumAge(minimumAge uint) exchangedomain.Exchange {
	exchange := ValidExchange()
	exchange.MinimumAge = minimumAge
	return exchange
}
