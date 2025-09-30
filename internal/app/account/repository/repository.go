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

func (r Repository) Update(ctx context.Context, account domain.Account) error {
	model := fromDomain(account)
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r Repository) FindByID(ctx context.Context, id uint) (domain.Account, error) {
	var model Account

	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Account{}, err
	}

	return model.toDomain(), nil
}
