package repository

import (
	"context"

	"exchange-crypto-service-api/internal/app/account/domain"

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

func (r Repository) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	model := fromDomain(account)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Account{}, err
	}

	return model.toDomain(), nil
}
