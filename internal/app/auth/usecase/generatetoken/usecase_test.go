package generatetoken_test

import (
	"context"
	"errors"
	"testing"

	"exchange-crypto-service-api/internal/app/auth/usecase/generatetoken"
	"exchange-crypto-service-api/internal/app/auth/usecase/generatetoken/mocks"
	"exchange-crypto-service-api/internal/app/auth/usecase/generatetoken/testdata"
	"exchange-crypto-service-api/internal/app/user/domain"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewUserRepository(t)
	useCase := generatetoken.New(userRepo)

	userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidUser(), nil)

	result, err := useCase.Execute(ctx, 1)

	require.NoError(t, err)
	require.NotEmpty(t, result.Token)
	require.NotZero(t, result.ExpiresAt)
}

func TestUseCase_Execute_Error(t *testing.T) {
	ctx := context.Background()
	userRepo := mocks.NewUserRepository(t)
	useCase := generatetoken.New(userRepo)

	userRepo.EXPECT().FindByID(mock.Anything, testdata.InvalidUserID).Return(domain.User{}, errors.New("user not found"))

	result, err := useCase.Execute(ctx, testdata.InvalidUserID)

	require.Empty(t, result)
	require.Error(t, err)
	require.Equal(t, "user not found", err.Error())
}
