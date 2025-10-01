package generatetoken

import (
	"context"

	"exchange-crypto-service-api/internal/app/user/domain"
	"exchange-crypto-service-api/internal/jwt"
)

//go:generate mockery --all

type (
	UserRepository interface {
		FindByID(ctx context.Context, id uint) (domain.User, error)
	}

	UseCase struct {
		userRepository UserRepository
	}
)

func New(userRepository UserRepository) UseCase {
	return UseCase{
		userRepository: userRepository,
	}
}

func (u UseCase) Execute(ctx context.Context, userID uint) (jwt.TokenResponse, error) {
	user, err := u.userRepository.FindByID(ctx, userID)
	if err != nil {
		return jwt.TokenResponse{}, err
	}

	tokenRequest := jwt.TokenRequest{
		UserID:         user.ID,
		Username:       user.Username,
		DocumentNumber: user.DocumentNumber,
	}

	return jwt.Instance().GenerateToken(tokenRequest)
}
