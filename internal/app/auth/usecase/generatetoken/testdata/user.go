package testdata

import (
	"time"

	"exchange-crypto-service-api/internal/app/user/domain"
)

const InvalidUserID = uint(999)

func ValidUser() domain.User {
	return domain.User{
		ID:             1,
		Username:       "testuser",
		DocumentNumber: "123456789",
		DateOfBirth:    time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
	}
}
