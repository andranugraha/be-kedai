package repository

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/user/entity"
	errs "kedai/backend/be-kedai/internal/domain/user/error"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(ID int) (*entity.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

type UserRConfig struct {
	DB *gorm.DB
}

func NewUserRepository(cfg *UserRConfig) UserRepository {
	return &userRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *userRepositoryImpl) GetByID(ID int) (*entity.User, error) {
	var user entity.User

	err := r.db.Where("user_id = ?", ID).Preload("Profile").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &user, nil
}
