package account

import (
	"exchange-crypto-service-api/pkg/apperror"
)

const ErrCodeInvalidMinimumAge = "error_invalid_minimum_age"

func NewErrInvalidMinimumAge(requiredAge, currentAge uint) error {
	return apperror.New(ErrCodeInvalidMinimumAge,
		"user does not meet minimum age requirement for this exchange: required %d, current %d", requiredAge, currentAge)
}
