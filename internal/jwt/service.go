package jwt

import (
	"errors"
	"sync"
	"time"

	"exchange-crypto-service-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	instance Service
	once     sync.Once
)

func Initialize() {
	once.Do(func() {
		cfg := config.LoadJWT()
		instance = Service{
			secret:     cfg.Secret,
			expiration: cfg.ExpirationHours,
		}
	})
}

func Instance() Service {
	return instance
}

func (s Service) ValidateToken(tokenString string) (Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, errors.New("invalid token claims")
	}

	return Claims{
		UserID:         uint(claims["user_id"].(float64)),
		Username:       claims["username"].(string),
		DocumentNumber: claims["document_number"].(string),
	}, nil
}

func (s Service) GenerateToken(request TokenRequest) (TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(s.expiration)

	claims := jwt.MapClaims{
		"user_id":         request.UserID,
		"username":        request.Username,
		"document_number": request.DocumentNumber,
		"exp":             expiresAt.Unix(),
		"iat":             now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}
