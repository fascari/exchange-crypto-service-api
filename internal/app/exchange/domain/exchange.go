package domain

import (
	"time"
)

type Exchange struct {
	ID                uint
	Name              string
	MinimumAge        uint
	MaxTransferAmount float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
