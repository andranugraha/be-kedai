package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	model "kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(ID int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
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

func (r *userRepositoryImpl) GetByID(ID int) (*model.User, error) {
	var user model.User

	err := r.db.Where("user_id = ?", ID).Preload("Profile").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.Where("email = ?", email).Preload("Profile").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &user, nil
}
