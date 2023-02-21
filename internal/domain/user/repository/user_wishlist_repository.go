package repository

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserWishlistRepository interface {
	AddUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error)
}

type userWishlistRepositoryImpl struct {
	db *gorm.DB
}

type UserWishlistRConfig struct {
	DB *gorm.DB
}

func NewUserWishlistRepository(cfg *UserWishlistRConfig) UserWishlistRepository {
	return &userWishlistRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *userWishlistRepositoryImpl) AddUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error) {
	err := r.db.Create(userWishlist).Error
	if err != nil {
		if errs.IsDuplicateKeyError(err) {
			return nil, errs.ErrProductInWishlist
		}
		return nil, err
	}

	return userWishlist, nil
}
