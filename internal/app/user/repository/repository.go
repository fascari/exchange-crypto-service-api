package repository

import (
	"context"
	"errors"

	"exchange-crypto-service-api/internal/app/user/domain"
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

func (r Repository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	model := fromDomain(user)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.User{}, err
	}

	return model.toDomain(), nil
}

func (r Repository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var user userModel

	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user.toDomain(), nil
}

func (r Repository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var user userModel

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user.toDomain(), nil
}

func (r Repository) FindUserBalances(ctx context.Context, userID uint) (domain.UserBalance, error) {
	var models []userExchangeBalanceModel
	var result domain.UserBalance

	err := telemetry.TraceRepository(ctx, "repository.find_user_balances",
		func(ctx context.Context) ([]attribute.KeyValue, error) {
			err := r.db.WithContext(ctx).
				Table("user_exchange_balances").
				Select(`user_id, username, exchange_name, balance, 
							  SUM(balance) OVER (PARTITION BY user_id) as total_balance, 
							  created_at, updated_at `).
				Where("user_id = ?", userID).
				Find(&models).Error

			if err != nil {
				return nil, err
			}

			if len(models) == 0 {
				return nil, nil
			}

			result = toDomain(models)

			return []attribute.KeyValue{
				attribute.Int("balances_count", len(models)),
			}, nil
		},
		attribute.Int("user_id", int(userID)),
	)

	return result, err
}
