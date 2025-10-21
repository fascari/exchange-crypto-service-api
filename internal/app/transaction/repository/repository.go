package repository

import (
	"context"
	"errors"
	"time"

	"exchange-crypto-service-api/internal/app/transaction"
	"exchange-crypto-service-api/internal/app/transaction/domain"
	"exchange-crypto-service-api/internal/database"
	"exchange-crypto-service-api/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) CheckIdempotency(ctx context.Context, accountID uint, idempotencyKey string) error {
	db := r.db
	if tx := database.TXFromContext(ctx); tx != nil {
		db = tx
	}

	// Check if transaction with this idempotency key already exists
	// The UNIQUE constraint (account_id, idempotency_key) provides the ultimate protection
	// against race conditions in distributed systems
	var existing transactionModel
	err := db.WithContext(ctx).
		Where("idempotency_key = ? AND account_id = ?", idempotencyKey, accountID).
		First(&existing).Error

	if err == nil {
		return transaction.ErrDuplicateIdempotencyKey{TransactionID: existing.TransactionID}
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

func (r Repository) Create(ctx context.Context, tr domain.Transaction) error {
	model := fromDomain(tr)

	db := r.db
	if tx := database.TXFromContext(ctx); tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Create(&model).Error
}

func (r Repository) FindDailyTransactions(ctx context.Context, startDate, endDate time.Time) ([]domain.DailyTransaction, error) {
	transactions := make([]domain.DailyTransaction, 0)

	err := telemetry.TraceRepository(ctx, "repository.find_daily_transactions",
		func(ctx context.Context) ([]attribute.KeyValue, error) {
			endDate = endDate.Add(24*time.Hour - time.Nanosecond)

			err := r.db.WithContext(ctx).
				Table("transactions t").
				Select(`e.name as exchange, 
					DATE(t.created_at) as date, 
					SUM(CASE WHEN t.type = 'DEPOSIT' THEN t.amount ELSE -t.amount END) as total_amount`).
				Joins("JOIN accounts a ON t.account_id = a.id").
				Joins("JOIN exchanges e ON a.exchange_id = e.id").
				Where("t.created_at BETWEEN ? AND ?", startDate, endDate).
				Where("t.deleted_at IS NULL AND a.deleted_at IS NULL AND e.deleted_at IS NULL").
				Group("e.name, DATE(t.created_at)").
				Order("DATE(t.created_at), e.name").
				Scan(&transactions).Error

			if err != nil {
				return nil, err
			}

			return []attribute.KeyValue{
				attribute.Int("result_count", len(transactions)),
			}, nil
		},
		attribute.String("start_date", startDate.Format(time.DateOnly)),
		attribute.String("end_date", endDate.Format(time.DateOnly)),
	)

	return transactions, err
}
