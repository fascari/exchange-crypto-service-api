package repository

import (
	"context"

	"exchange-crypto-service-api/internal/app/account/domain"
	"exchange-crypto-service-api/internal/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := database.WithTX(ctx, tx)
		return fn(ctxWithTx)
	})
}

func (r Repository) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	model := fromDomain(account)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Account{}, err
	}

	return model.toDomain(), nil
}

func (r Repository) Update(ctx context.Context, account domain.Account) error {
	model := fromDomain(account)

	db := r.db
	if tx := database.TXFromContext(ctx); tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Save(&model).Error
}

func (r Repository) FindByID(ctx context.Context, id uint) (domain.Account, error) {
	var model Account

	db := r.db
	if tx := database.TXFromContext(ctx); tx != nil {
		db = tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate})
	}

	if err := db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Account{}, err
	}

	return model.toDomain(), nil
}
