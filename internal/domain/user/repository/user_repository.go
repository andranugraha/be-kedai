package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"
	errs "kedai/backend/be-kedai/internal/common/error"

	"gorm.io/gorm"
)

type UserRepository interface {
		SignUp(user *model.User) (*model.User, error)
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

func (r *userRepositoryImpl) SignUp(user *model.User) (*model.User, error) {
	err := r.db.Create(&user)
	if err.Error != nil {
		return nil, err.Error
	}

	if err.RowsAffected == 0 {
		return nil, errs.ErrUserAlreadyExist
	}

	return user, nil
}