package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserWishlistRepository interface {
	GetUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error)
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

func (r *userWishlistRepositoryImpl) GetUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error) {
	var res model.UserWishlist

	err := r.db.Where("user_id = ? AND product_id = ?", userWishlist.UserID, userWishlist.ProductID).First(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}
