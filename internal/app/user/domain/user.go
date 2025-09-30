package domain

import (
	"time"
)

type User struct {
	ID             uint
	Username       string
	DateOfBirth    time.Time
	DocumentNumber string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
