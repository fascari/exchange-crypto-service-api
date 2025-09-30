package repository

import (
	"time"

	"exchange-crypto-service-api/internal/app/user/domain"

	"gorm.io/gorm"
)

type userModel struct {
	Username       string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	DateOfBirth    time.Time `gorm:"type:date;not null" json:"date_of_birth"`
	DocumentNumber string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	gorm.Model
}

func (userModel) TableName() string {
	return "users"
}

func fromDomain(user domain.User) userModel {
	return userModel{
		Username:       user.Username,
		DateOfBirth:    user.DateOfBirth,
		DocumentNumber: user.DocumentNumber,
	}
}

func (u userModel) toDomain() domain.User {
	return domain.User{
		ID:             u.ID,
		Username:       u.Username,
		DateOfBirth:    u.DateOfBirth,
		DocumentNumber: u.DocumentNumber,
	}
}
