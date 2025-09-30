package repository

import (
	"context"
	"errors"

	"exchange-crypto-service-api/internal/app/exchange/domain"

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

func (r Repository) FindByID(ctx context.Context, id uint) (domain.Exchange, error) {
	var model exchangeModel

	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Exchange{}, errors.New("exchange not found")
		}
		return domain.Exchange{}, err
	}

	return model.ToDomain(), nil
}
