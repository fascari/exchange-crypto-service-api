package testdata

import (
	"time"

	userdomain "exchange-crypto-service-api/internal/app/user/domain"
)

const (
	defaultUsername       = "John Doe"
	defaultDocumentNumber = "123456789"
)

func createUser(id uint, username, documentNumber string, birthDate time.Time) userdomain.User {
	return userdomain.User{
		ID:             id,
		Username:       username,
		DocumentNumber: documentNumber,
		DateOfBirth:    birthDate,
	}
}

func ValidUser() userdomain.User {
	return createUser(1, defaultUsername, defaultDocumentNumber, time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC))
}

func UnderageUser() userdomain.User {
	return createUser(2, "Jane Smith", "987654321", time.Date(2010, 1, 15, 0, 0, 0, 0, time.UTC))
}

func UserWithFutureBirthday() userdomain.User {
	return createUser(1, defaultUsername, defaultDocumentNumber, time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC))
}
