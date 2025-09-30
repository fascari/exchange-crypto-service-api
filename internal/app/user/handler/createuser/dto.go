package createuser

import (
	"time"

	"exchange-crypto-service-api/internal/app/user/domain"
	"exchange-crypto-service-api/pkg/validator"
)

type (
	InputPayload struct {
		Username       string `json:"username" validate:"required"`
		DateOfBirth    string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
		DocumentNumber string `json:"document_number" validate:"required"`
	}

	OutputPayload struct {
		ID             uint   `json:"id"`
		Username       string `json:"username"`
		DateOfBirth    string `json:"date_of_birth"`
		DocumentNumber string `json:"document_number"`
	}
)

func (p InputPayload) Validate() error {
	return validator.Validate(p)
}

func (p InputPayload) ToDomain() (domain.User, error) {
	dateOfBirth, err := time.Parse("2006-01-02", p.DateOfBirth)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Username:       p.Username,
		DateOfBirth:    dateOfBirth,
		DocumentNumber: p.DocumentNumber,
	}, nil
}

func ToOutputPayload(user domain.User) OutputPayload {
	return OutputPayload{
		ID:             user.ID,
		Username:       user.Username,
		DateOfBirth:    user.DateOfBirth.Format("2006-01-02"),
		DocumentNumber: user.DocumentNumber,
	}
}
