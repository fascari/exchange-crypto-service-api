package createaccount_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"exchange-crypto-service-api/internal/app/account/domain"
	"exchange-crypto-service-api/internal/app/account/usecase/createaccount"
	"exchange-crypto-service-api/internal/app/account/usecase/createaccount/mocks"
	"exchange-crypto-service-api/internal/app/account/usecase/createaccount/testdata"
	exchangedomain "exchange-crypto-service-api/internal/app/exchange/domain"
	userdomain "exchange-crypto-service-api/internal/app/user/domain"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute_Success(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*mocks.Repository, *mocks.UserRepository, *mocks.ExchangeRepository)
		inputAccount   domain.Account
		expectedResult domain.Account
	}{
		{
			name: "should create account successfully when user meets minimum age requirement",
			setupMocks: func(accountRepo *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidUser(), nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidExchange(), nil)
				accountRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("domain.Account")).Return(testdata.ValidAccount(), nil)
			},
			inputAccount:   testdata.InputAccount(),
			expectedResult: testdata.ValidAccount(),
		},
		{
			name: "should create account successfully when user birthday has not occurred this year (YearDay check)",
			setupMocks: func(accountRepo *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(2)).Return(testdata.UserWithFutureBirthday(), nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(2)).Return(testdata.ExchangeWithMinimumAge(24), nil)
				accountRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("domain.Account")).Return(testdata.AccountWithID(2), nil)
			},
			inputAccount:   testdata.InputAccountWithUserAndExchange(2, 2),
			expectedResult: testdata.AccountWithID(2),
		},
		{
			name: "should create account successfully when user has future birth date but calculateAge returns 0",
			setupMocks: func(accountRepo *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				futureUser := testdata.ValidUser()
				futureUser.DateOfBirth = time.Now().AddDate(1, 0, 0)

				userRepo.EXPECT().FindByID(mock.Anything, uint(5)).Return(futureUser, nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(5)).Return(testdata.ExchangeWithMinimumAge(0), nil)
				accountRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("domain.Account")).Return(testdata.AccountWithID(5), nil)
			},
			inputAccount:   testdata.InputAccountWithUserAndExchange(5, 5),
			expectedResult: testdata.AccountWithID(5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := setupUseCase(t, tt.setupMocks)

			result, err := useCase.Execute(context.Background(), tt.inputAccount)

			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestUseCase_Execute_Error(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*mocks.Repository, *mocks.UserRepository, *mocks.ExchangeRepository)
		inputAccount  domain.Account
		expectedError string
	}{
		{
			name: "should return error when user is not found",
			setupMocks: func(_ *mocks.Repository, userRepo *mocks.UserRepository, _ *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(userdomain.User{}, errors.New("user not found"))
			},
			inputAccount:  testdata.InputAccount(),
			expectedError: "user not found",
		},
		{
			name: "should return error when exchange is not found",
			setupMocks: func(_ *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidUser(), nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(exchangedomain.Exchange{}, errors.New("exchange not found"))
			},
			inputAccount:  testdata.InputAccount(),
			expectedError: "exchange not found",
		},
		{
			name: "should return error when user does not meet minimum age requirement",
			setupMocks: func(_ *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.UnderageUser(), nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidExchange(), nil)
			},
			inputAccount:  testdata.InputAccount(),
			expectedError: "user does not meet minimum age requirement",
		},
		{
			name: "should return error when account creation fails",
			setupMocks: func(accountRepo *mocks.Repository, userRepo *mocks.UserRepository, exchangeRepo *mocks.ExchangeRepository) {
				userRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidUser(), nil)
				exchangeRepo.EXPECT().FindByID(mock.Anything, uint(1)).Return(testdata.ValidExchange(), nil)
				accountRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("domain.Account")).Return(domain.Account{}, errors.New("database connection failed"))
			},
			inputAccount:  testdata.InputAccount(),
			expectedError: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := setupUseCase(t, tt.setupMocks)

			result, err := useCase.Execute(context.Background(), tt.inputAccount)

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
			require.Equal(t, domain.Account{}, result)
		})
	}
}

func setupUseCase(t *testing.T, setupMocks func(*mocks.Repository, *mocks.UserRepository, *mocks.ExchangeRepository)) createaccount.UseCase {
	mockAccountRepo := mocks.NewRepository(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockExchangeRepo := mocks.NewExchangeRepository(t)

	setupMocks(mockAccountRepo, mockUserRepo, mockExchangeRepo)

	return createaccount.New(mockAccountRepo, mockUserRepo, mockExchangeRepo)
}
