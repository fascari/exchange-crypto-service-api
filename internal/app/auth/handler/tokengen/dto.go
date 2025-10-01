package tokengen

import (
	"exchange-crypto-service-api/internal/jwt"
)

type (
	OutputPayload struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}
)

func toOutputPayload(response jwt.TokenResponse) OutputPayload {
	return OutputPayload{
		Token:     response.Token,
		ExpiresAt: response.ExpiresAt.Format("2006-01-02T15:04:05Z"),
	}
}
