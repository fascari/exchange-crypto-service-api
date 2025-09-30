package domain

import (
	"time"
)

type Account struct {
	ID         uint
	UserID     uint
	ExchangeID uint
	Balance    float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
