package repository

import (
	"context"

	"exchange-crypto-service-api/internal/app/transaction/domain"

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

func (r Repository) Create(ctx context.Context, transaction domain.Transaction) error {
	model := fromDomain(transaction)
	return r.db.WithContext(ctx).Create(&model).Error
}
