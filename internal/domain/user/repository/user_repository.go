package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
<<<<<<< HEAD
	"kedai/backend/be-kedai/internal/domain/user/model"
=======
	model "kedai/backend/be-kedai/internal/domain/user/model"
>>>>>>> a58ad6c9e3eb0eec59cef9daf55dd08bdeaae6b3

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(ID int) (*model.User, error)
<<<<<<< HEAD
=======
	GetByEmail(email string) (*model.User, error)
>>>>>>> a58ad6c9e3eb0eec59cef9daf55dd08bdeaae6b3
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
<<<<<<< HEAD
=======

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
>>>>>>> a58ad6c9e3eb0eec59cef9daf55dd08bdeaae6b3
