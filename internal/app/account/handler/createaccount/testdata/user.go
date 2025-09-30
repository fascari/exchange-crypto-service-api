package testdata

import (
	"time"

	userdomain "exchange-crypto-service-api/internal/app/user/domain"
)

func ValidUser() userdomain.User {
	return userdomain.User{
		ID:             1,
		Username:       "John Doe",
		DocumentNumber: "123456789",
		DateOfBirth:    time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func UnderageUser() userdomain.User {
	return userdomain.User{
		ID:             2,
		Username:       "Jane Smith",
		DocumentNumber: "987654321",
		DateOfBirth:    time.Date(2010, 1, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
