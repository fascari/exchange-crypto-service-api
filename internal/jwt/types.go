package jwt

import "time"

type (
	Claims struct {
		UserID         uint   `json:"user_id"`
		Username       string `json:"username"`
		DocumentNumber string `json:"document_number"`
	}

	TokenRequest struct {
		UserID         uint
		Username       string
		DocumentNumber string
	}

	TokenResponse struct {
		Token     string
		ExpiresAt time.Time
	}

	Service struct {
		secret     string
		expiration time.Duration
	}
)
