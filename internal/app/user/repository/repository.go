package repository

import (
	"context"
	"errors"

	"exchange-crypto-service-api/internal/app/user/domain"

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
