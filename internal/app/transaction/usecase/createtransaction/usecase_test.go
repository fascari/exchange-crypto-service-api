//go:build integration

package createtransaction_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"exchange-crypto-service-api/internal/app/account/domain"
	"exchange-crypto-service-api/internal/app/account/repository"
	exchangeRepo "exchange-crypto-service-api/internal/app/exchange/repository"
	transactionDomain "exchange-crypto-service-api/internal/app/transaction/domain"
	transactionRepo "exchange-crypto-service-api/internal/app/transaction/repository"
	"exchange-crypto-service-api/internal/app/transaction/usecase/createtransaction"
	"exchange-crypto-service-api/internal/app/transaction/usecase/createtransaction/testdata"
	"exchange-crypto-service-api/internal/testing/integration"

	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	integration.Suite
	accountRepo  repository.Repository
	exchangeRepo exchangeRepo.Repository
	transRepo    transactionRepo.Repository
	usecase      createtransaction.UseCase
}

func (s *TransactionTestSuite) SetupSuite() {
	s.Suite.WithFixtures("./testdata/fixtures").
		WithSkipShortTests()
}

func (s *TransactionTestSuite) SetupTest() {
	s.Suite.SetupTest(s.T())
	s.accountRepo = repository.New(s.DB)
	s.exchangeRepo = exchangeRepo.New(s.DB)
	s.transRepo = transactionRepo.New(s.DB)
	s.usecase = createtransaction.New(s.accountRepo, s.exchangeRepo, s.transRepo)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}

func (s *TransactionTestSuite) TestConcurrentDeposits() {
	ctx, cancel := context.WithTimeout(context.Background(), testdata.ShortTimeout)
	defer cancel()

	account := testdata.AccountWithBalance()
	expectedBalance := testdata.ExpectedBalanceAfterDeposits(account.Balance, testdata.ConcurrentDeposits, testdata.MediumAmount)

	errors := s.ExecuteConcurrent(testdata.ConcurrentDeposits, func() {
		_, _ = s.usecase.Execute(ctx, transactionDomain.Deposit, account.ID, testdata.MediumAmount, s.NewUUID())
	})

	s.AssertNoErrors(errors)
	s.assertBalance(ctx, account.ID, expectedBalance, "Balance should match expected value")
}

func (s *TransactionTestSuite) TestConcurrentWithdrawals() {
	ctx, cancel := context.WithTimeout(context.Background(), testdata.ShortTimeout)
	defer cancel()

	account := testdata.AccountWithBalance()
	var successCount atomic.Int32

	errors := s.ExecuteConcurrent(testdata.ConcurrentWithdrawals, func() {
		_, err := s.usecase.Execute(ctx, transactionDomain.Withdrawal, account.ID, 50.0, s.NewUUID())
		if err == nil {
			successCount.Add(1)
		}
	})

	s.AssertNoErrors(errors)
	s.Require().Equal(int32(testdata.ConcurrentWithdrawals), successCount.Load(), "All withdrawals should succeed")

	finalAccount, err := s.accountRepo.FindByID(ctx, account.ID)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(finalAccount.Balance, 0.0, "Balance should never be negative!")

	s.assertBalanceDelta(ctx, account.ID, 0.0, 0.01, "Balance should be 0")
}

func (s *TransactionTestSuite) TestIdempotencyKeyPreventsDoubleSpend() {
	ctx, cancel := context.WithTimeout(context.Background(), testdata.ShortTimeout)
	defer cancel()

	account := testdata.AccountWithBalance()
	idempotencyKey := s.NewUUID()
	var successCount atomic.Int32

	errors := s.ExecuteConcurrent(testdata.ConcurrentDeposits, func() {
		_, err := s.usecase.Execute(ctx, transactionDomain.Deposit, account.ID, testdata.LargeAmount, idempotencyKey)
		if err == nil {
			successCount.Add(1)
		}
	})

	s.AssertNoErrors(errors)
	s.Require().Equal(int32(1), successCount.Load(), "Only 1 transaction should succeed with same idempotency key")
	s.assertBalance(ctx, account.ID, account.Balance+testdata.LargeAmount, "Balance should reflect only one deposit")
}

func (s *TransactionTestSuite) TestMixedConcurrentOperations() {
	ctx, cancel := context.WithTimeout(context.Background(), testdata.MediumTimeout)
	defer cancel()

	account := testdata.AccountWithBalance()
	var wg sync.WaitGroup
	var successCount atomic.Int32

	for i := 0; i < 10; i++ {
		wg.Go(func() {
			_, err := s.usecase.Execute(ctx, transactionDomain.Deposit, account.ID, testdata.MediumAmount, s.NewUUID())
			if err == nil {
				successCount.Add(1)
			}
		})
		wg.Go(func() {
			_, err := s.usecase.Execute(ctx, transactionDomain.Withdrawal, account.ID, testdata.MediumAmount, s.NewUUID())
			if err == nil {
				successCount.Add(1)
			}
		})
	}
	wg.Wait()

	s.Require().Equal(int32(20), successCount.Load(), "All 20 operations (10 deposits + 10 withdrawals) should succeed")

	s.assertBalance(ctx, account.ID, account.Balance, "Balance should be unchanged")
}

func (s *TransactionTestSuite) TestRaceConditionPrevention() {
	ctx, cancel := context.WithTimeout(context.Background(), testdata.MediumTimeout)
	defer cancel()

	account := testdata.AccountEmpty()
	expectedBalance := float64(testdata.ConcurrentDeposits) * testdata.SmallAmount

	errors := s.ExecuteConcurrent(testdata.ConcurrentDeposits, func() {
		_, _ = s.usecase.Execute(ctx, transactionDomain.Deposit, account.ID, testdata.SmallAmount, s.NewUUID())
	})

	s.AssertNoErrors(errors)
	s.assertBalance(ctx, account.ID, expectedBalance, "Balance should match all deposits (no lost updates)")
}

func (s *TransactionTestSuite) TestTransactionsErrors() {
	testCases := []struct {
		name               string
		transactionType    transactionDomain.TransactionType
		Account            domain.Account
		amount             float64
		shouldCheckBalance bool
	}{
		{
			name:               "Deposit with invalid account ID",
			transactionType:    transactionDomain.Deposit,
			Account:            testdata.InvalidAccount(),
			amount:             testdata.MediumAmount,
			shouldCheckBalance: false,
		},
		{
			name:               "Withdrawal with insufficient balance",
			transactionType:    transactionDomain.Withdrawal,
			Account:            testdata.AccountWithBalance(),
			amount:             testdata.ExcessiveAmount(testdata.AccountWithBalance().Balance),
			shouldCheckBalance: true,
		},
		{
			name:               "Deposit with zero amount",
			transactionType:    transactionDomain.Deposit,
			Account:            testdata.AccountWithBalance(),
			amount:             0.0,
			shouldCheckBalance: true,
		},
		{
			name:               "Deposit with negative amount",
			transactionType:    transactionDomain.Deposit,
			Account:            testdata.AccountWithBalance(),
			amount:             -100.0,
			shouldCheckBalance: true,
		},
		{
			name:               "Withdrawal from empty account",
			transactionType:    transactionDomain.Withdrawal,
			Account:            testdata.AccountEmpty(),
			amount:             testdata.SmallAmount,
			shouldCheckBalance: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), testdata.ShortTimeout)
			defer cancel()

			account := tc.Account

			_, err := s.usecase.Execute(ctx, tc.transactionType, account.ID, tc.amount, s.NewUUID())

			s.Require().Error(err)

			if tc.shouldCheckBalance {
				s.assertBalance(ctx, account.ID, account.Balance, "Balance should remain unchanged")
			}
		})
	}
}

func (s *TransactionTestSuite) assertBalance(ctx context.Context, accountID uint, expectedBalance float64, msgAndArgs ...interface{}) {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	s.Require().NoError(err)
	s.Require().Equal(expectedBalance, account.Balance, msgAndArgs...)
}

func (s *TransactionTestSuite) assertBalanceDelta(ctx context.Context, accountID uint, expectedBalance, delta float64, msgAndArgs ...interface{}) {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	s.Require().NoError(err)
	s.Require().InDelta(expectedBalance, account.Balance, delta, msgAndArgs...)
}
