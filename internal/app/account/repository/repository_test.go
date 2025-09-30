package repository_test

import (
	"context"
	"testing"

	"exchange-crypto-service-api/internal/app/account/domain"
	"exchange-crypto-service-api/internal/app/account/repository"
	"exchange-crypto-service-api/internal/app/account/repository/testdata"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&repository.Account{})
	require.NoError(t, err)

	return db
}

func TestRepository_Create_Success(t *testing.T) {
	tests := []struct {
		name         string
		inputAccount domain.Account
	}{
		{
			name:         "should create account successfully with valid data",
			inputAccount: testdata.ValidAccount(),
		},
		{
			name:         "should create account successfully with zero balance",
			inputAccount: testdata.ZeroBalanceAccount(),
		},
		{
			name:         "should create account successfully with high balance",
			inputAccount: testdata.HighBalanceAccount(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := repository.New(db)

			result, err := repo.Create(context.Background(), tt.inputAccount)

			require.NoError(t, err)
			require.NotZero(t, result.ID)
			require.Equal(t, tt.inputAccount.UserID, result.UserID)
			require.Equal(t, tt.inputAccount.ExchangeID, result.ExchangeID)
			require.Equal(t, tt.inputAccount.Balance, result.Balance)

			var count int64
			db.Model(&repository.Account{}).Count(&count)
			require.Greater(t, count, int64(0))
		})
	}
}

func TestRepository_Update_Success(t *testing.T) {
	tests := []struct {
		name           string
		initialAccount domain.Account
		newBalance     float64
	}{
		{
			name:           "should update account balance successfully",
			initialAccount: testdata.AccountForUpdate(),
			newBalance:     200.0,
		},
		{
			name:           "should update account to zero balance",
			initialAccount: testdata.AccountForUpdate(),
			newBalance:     0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := repository.New(db)

			created, err := repo.Create(context.Background(), tt.initialAccount)
			require.NoError(t, err)

			updatedAccount := testdata.UpdatedAccount(created.ID, tt.newBalance)
			err = repo.Update(context.Background(), updatedAccount)

			require.NoError(t, err)

			found, err := repo.FindByID(context.Background(), created.ID)
			require.NoError(t, err)
			require.Equal(t, tt.newBalance, found.Balance)
		})
	}
}

func TestRepository_FindByID_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)

	inputAccount := testdata.ValidAccount()

	created, err := repo.Create(context.Background(), inputAccount)
	require.NoError(t, err)

	found, err := repo.FindByID(context.Background(), created.ID)

	require.NoError(t, err)
	require.Equal(t, created, found)
}

func TestRepository_Create_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)

	_, err := repo.Create(context.Background(), testdata.InvalidAccount())

	if err != nil {
		require.Error(t, err)
	}
}

func TestRepository_FindByID_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)

	_, err := repo.FindByID(context.Background(), 999)

	require.Error(t, err)
	require.Contains(t, err.Error(), "record not found")
}

func TestRepository_Update_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.New(db)

	nonExistentAccount := testdata.UpdatedAccount(999, 100.0)

	err := repo.Update(context.Background(), nonExistentAccount)

	require.NoError(t, err)
}
